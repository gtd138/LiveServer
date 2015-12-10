package connector

import (
	. "framework/server"
	"msg_proto"
)

type ConnectorRPCService struct {
	*RPCService
}

// 添加客户端的Token
func (this *ConnectorRPCService) AddClientToken(arg *msg_proto.Token, bOk *bool) error {
	server := this.ServerInterface.(*Connector)
	token := *arg
	server.TokenManager.Add(&token)
	*bOk = true
	return nil
}
