package gate

// 网关服
import (
	"common"
	"config/message_config"
	"config/server_config"
	. "framework/server"
	"gamelogic/rpc"
	"msg_proto"
)

type Gate struct {
	*FrontendServer
	CotorStateMap *common.BeeMap // 连接服状态，map[name]=*ConnState
}

func NewGate(server_type string) *Gate {
	instance := &Gate{
		FrontendServer: NewFrontendServer(server_type),
		CotorStateMap:  common.NewBeeMap(),
	}
	instance.IServer = instance
	// 注册RPC
	instance.RegisterRPC(&GateRPCService{&RPCService{ServerInterface: instance}})
	// 初始化
	instance.Init()
	return instance
}

func (this *Gate) Init() {
	this.RegisterMsgHandle()
	this.BaseServer.Init()
}

// 注册消息回调
func (this *Gate) RegisterMsgHandle() {
	reg_func := msg_conf.MessageCallback()
	reg_func.Register(msg_proto.MsgCmd_RequestLogin_C, this.RequestLogin)
}

// 获取最佳连接服
func (this *Gate) GetFitConnector() (Type string, Id int, IP, Port string) {
	conn_num := -1
	var state *game_rpc.ConnState
	this.CotorStateMap.Foreach(func(k, v interface{}) bool {
		if conn_num == -1 {
			state = v.(*game_rpc.ConnState)
			conn_num = state.Num
		}
		if conn_num > v.(*game_rpc.ConnState).Num {
			conn_num = v.(*game_rpc.ConnState).Num
			state = v.(*game_rpc.ConnState)
		}
		return false
	})
	if state == nil {
		println("获取适合的连接服失败！")
		return
	}
	Type = state.Type
	Id = state.Id
	config := server_conf.GetSingleton().GetConfig(state.Type, state.Id)
	IP = config.Host
	Port = config.Port
	return
}
