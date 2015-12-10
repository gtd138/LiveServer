package server

import (
	"common"
	. "framework/network"
)

const (
	RECEIVER_MESSAGE = "ReceiveMessage"
)

// 后端服务器
type BackendServer struct {
	*BaseServer
	//*ProxyManager // 代理模块
}

func NewBackendServer(server_type string) *BackendServer {
	instance := new(BackendServer)
	instance.BaseServer = NewBaseServer(server_type)
	return instance
}

// 初始化
func (this *BackendServer) Init() {
	this.BaseServer.Init()
	this.StartRPCService(this.Addr, this.Port)
}

// 逻辑循环
func (this *BackendServer) LogicLoop() {
}

// 发送请求
func (this *BackendServer) SendRequest(t *common.Timer, args ...interface{}) bool {
	// 把接收的消息抽取出来
	msg_que := this.MessageManager.PopAll(SEND_MSG)
	rpc_client_list := this.GetRPCClients()
	// 分派消息
	for i, v := range msg_que {
		proxy := this.FindProxy(v.Sid)
		if proxy == nil {
			continue
		}
		for j, client := range rpc_client_list {
			if client.GetRemoteServerType() == proxy.FrontendType {
				if v.IsBc {
					rpc_client_list[j].PushRequest(msg_que[i])
				} else {
					if proxy.FrontendID == client.GetRemoteServerID() {
						rpc_client_list[j].PushRequest(msg_que[i])
						break
					}
				}
			}
		}
	}
	// 发送消息
	for i := range rpc_client_list {
		rpc_client_list[i].RemoteCall(RECEIVER_MESSAGE)
	}
	return true
}
