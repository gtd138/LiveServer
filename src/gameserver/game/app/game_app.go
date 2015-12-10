package game

// 游戏服
import (
	. "framework/server"
)

type Game struct {
	*BackendServer
}

func NewGame(server_type string) *Game {
	instance := &Game{
		BackendServer: NewBackendServer(server_type),
	}
	instance.IServer = instance
	// 注册RPC
	instance.RegisterRPC(&GameRPCService{&RPCService{ServerInterface: instance}})
	// 初始化
	instance.Init()
	return instance
}
