package gate

import (
	. "framework/server"
	"gamelogic/rpc"
)

type GateRPCService struct {
	*RPCService
}

// 设置连接数量
func (this *GateRPCService) SetConnectorConnNum(arg *game_rpc.ConnState, bOk *bool) error {
	*bOk = true
	server := this.ServerInterface.(*Gate)
	// 先拷贝一份
	state := *arg
	server.CotorStateMap.Set(state.Name, &state)
	return nil
}
