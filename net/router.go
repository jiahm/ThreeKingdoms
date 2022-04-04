/**
* @Author : jiahongming
* @Description :
* @Time : 2022/4/3 10:24 PM
* @File : router
* @Software: GoLand
**/
package net

import "strings"

type HandlerFunc func(req *WsMsgReq, rsp *WsMsgRsp)

//account login||logout
type group struct {
	prefix     string
	handlerMap map[string]HandlerFunc
}

func
func (g *group) exec(name string, req *WsMsgReq, rsp *WsMsgRsp) {
	h := g.handlerMap[name]
	if h != nil {
		h(req, rsp)
	}
}

type router struct {
	group []*group
}

func (r *router) Run(req *WsMsgReq, rsp *WsMsgRsp) {
	//req.Body.Name路径  登陆业务 account.login （account组标识）login 路由标识
	strs := strings.Split(req.Body.Name, ".")
	prefix := ""
	name := ""
	if len(strs) == 2 {
		prefix = strs[0]
		name = strs[1]
	}
	for _, g := range r.group {
		if g.prefix == prefix {
			g.exec(name, req, rsp)
		}
	}
}
