package server

import (
	. "framework/network"
)

type BindProxyArg struct {
	Sid        string
	ServerType string
	ServerID   int
}

// rpc服务
type RPCService struct {
	ServerInterface IServer // 传入的服务器接口
}

// 接收到路由的消息(后端使用)
func (this *RPCService) RouteMessage(msg_list []*MessageObject, pbOk *bool) error {
	this.ServerInterface.GetBaseServer().PushRevMessage(msg_list)
	*pbOk = true
	return nil
}

// 接收后端消息(前端)
func (this *RPCService) ReceiveMessage(msg_list []*MessageObject, pbOk *bool) error {
	this.ServerInterface.GetBaseServer().PushSendMessage(msg_list)
	*pbOk = true
	return nil
}

// 绑定代理
func (this *RPCService) BindProxy(arg *BindProxyArg, bOk *bool) error {
	this.ServerInterface.GetBaseServer().BindProxy(arg.Sid, arg.ServerType, arg.ServerID)
	*bOk = true
	return nil
}
