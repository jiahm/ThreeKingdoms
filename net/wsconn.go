/**
* @Author : jiahongming
* @Description :
* @Time : 2022/4/4 4:07 PM
* @File : wsconn
* @Software: GoLand
**/
package net

type ReqBody struct {
	Seq   int64       `json:"seq" :"seq"`
	Name  string      `json:"name" :"name"`
	Msg   interface{} `json:"msg" :"msg"`
	Proxy string      `json:"proxy"`
}

type RspBody struct {
	Seq  int64       `json:"seq" :"seq"`
	Name string      `json:"name" :"name"`
	Code int         `json:"code"`
	Msg  interface{} `json:"msg" :"msg"`
}

type WsMsgReq struct {
	Body *ReqBody
	Conn WSConn
}

type WsMsgRsp struct {
	Body *RspBody
}

//理解为 request请求 请求会有参数 请求中放参数 请求中取参数
type WSConn interface {
	SetProperty(key string, value interface{})
	GetProPerty(key string) (interface{}, error)
	RemoveProperty(key string)
	Addr() string
	Push(name string, data interface{})
}
