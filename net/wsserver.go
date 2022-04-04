/**
* @Author : jiahongming
* @Description :
* @Time : 2022/4/4 4:24 PM
* @File : wsserver
* @Software: GoLand
**/
package net

import (
	"ThreeKingdoms/utils"
	"encoding/json"
	"fmt"
	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

//websocket服务
type wsServer struct {
	wsConn       *websocket.Conn
	router       *router
	outChan      chan *WsMsgRsp         //写队列
	Seq          int64                  //序列
	property     map[string]interface{} //属性
	propertyLock sync.RWMutex
}

func NewWsServer(wsConn *websocket.Conn) *wsServer {
	return &wsServer{
		wsConn:   wsConn,
		outChan:  make(chan *WsMsgRsp, 1000),
		property: make(map[string]interface{}),
		Seq:      0,
	}
}

func (w *wsServer) Router(router *router) {
	w.router = router
}

func (w *wsServer) SetProperty(key string, value interface{}) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	w.property[key] = value
}
func (w *wsServer) GetProPerty(key string) (interface{}, error) {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()
	return w.property[key], nil
}
func (w *wsServer) RemoveProperty(key string) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	delete(w.property, key)
}
func (w *wsServer) Addr() string {
	return w.wsConn.RemoteAddr().String()
}
func (w *wsServer) Push(name string, data interface{}) {
	rsp := &WsMsgRsp{Body: &RspBody{Name: name, Msg: data, Seq: 0}}
	w.outChan <- rsp
}

//通道一旦建立，那么收发消息 就要一直监听才行
func (w *wsServer) Start() {
	//启动读写数据的处理逻辑
	go w.readMsgLoop()
	go w.writeMsgLoop()
}

func (w *wsServer) writeMsgLoop() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
			w.Close()
		}
	}()
	for {
		select {
		case msg := <-w.outChan:
			fmt.Println(msg)
			//w.wsConn.WriteMessage()
		}
	}
}

func (w *wsServer) readMsgLoop() {
	// 先读到客户端发送过来的数据，然后进行处理，然后回消息
	//经过路由 实际处理程序
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
			w.Close()
		}
	}()
	for {
		_, data, err := w.wsConn.ReadMessage()
		if err != nil {
			log.Println("收消息出现错过：", err)
			break
		}
		fmt.Println(data)
		//收到消息 解析消息 前端发送过来的消息就是json格式
		//1.data进行解压 unzip
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("解压数据出错，非法格式:", err)
			continue
		}
		//2.前端的消息 加密消息 进行解密
		secretKey, err := w.GetProPerty("secretKey")
		if err == nil {
			//有加密
			key := secretKey.(string)
			//客户端传过来的数据是加密的 需要解密
			d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
			if err != nil {
				log.Println("数据格式有误，解密失败：", err)
				//出错后 发起握手

			} else {
				data = d
			}
		}
		//3.data转为body
		body := &ReqBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			log.Println("数据格式有误，非法格式：", err)
		} else {
			//获取到前端传递的数据了，拿上这些数据 去具体的业务进行处理
			req := &WsMsgReq{Conn: w, Body: body}
			rsp := &WsMsgRsp{Body: &RspBody{Name: body.Name, Seq: req.Body.Seq}}
			w.router.Run(req, rsp)
			w.outChan <- rsp
		}
	}
	//rsp := &WsMsgRsp{Body: &RspBody{Name: name, Msg: data, Seq: 0}}
	//w.outChan <- rsp
	w.Close()
}

func (w *wsServer) Close() {
	_ = w.wsConn.Close()
}
