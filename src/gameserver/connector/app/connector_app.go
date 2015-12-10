package connector

// 连接服
import (
	"framework/network"
	. "framework/server"
	"gamelogic/rpc"
	"net"
	"strconv"
)

type Connector struct {
	*FrontendServer
	*SessionCheck // 检查会话
	*TokenManager
}

func NewConnector(server_type string) *Connector {
	instance := &Connector{
		FrontendServer: NewFrontendServer(server_type),
		SessionCheck:   NewSessionCheck(),
		TokenManager:   NewTokenManager(),
	}
	instance.IServer = instance
	// 注册RPC
	instance.RegisterRPC(&ConnectorRPCService{&RPCService{ServerInterface: instance}})
	// 初始化
	instance.Init()
	return instance
}

// 重写初始化函数
func (this *Connector) Init() {
	this.FrontendServer.Init()
	this.SessionCheck.SeTimeout = this.DeleteTimeoutSession
}

// 重写启动函数
func (this *Connector) Run() {
	this.SessionCheck.Begin()
	this.TokenManager.Begin()
	this.FrontendServer.Run()
}

// 重写TCP监听
func (this *Connector) TcpListenerCallback(conn net.Conn) {
	this.FrontendServer.TcpListenerCallback(conn)
	// 更新一下gete服的会话数
	this.UpdateConnNum()
	se := this.FindSessionByConn(conn)
	if se != nil {
		this.SessionCheck.Add(se)
	}
}

// 删除超时会话
func (this *Connector) DeleteTimeoutSession(se *network.Session) {
	this.CloseSession(se)
	// 更新一下gete服的会话数
	this.UpdateConnNum()
}

// 更新连接数
func (this *Connector) UpdateConnNum() {
	var bOk bool
	state := &game_rpc.ConnState{this.Type + strconv.Itoa(this.ID), this.Type, this.ID, this.SessionManager.Count}
	this.RPCCall("gate", 1, "SetConnectorConnNum", state, &bOk)
}
