package lobby

// 大厅服
import (
	. "framework/database"
	. "framework/server"
)

type Lobby struct {
	*BackendServer
	db *DB // 数据库
}

func NewLobby(server_type string) *Lobby {
	instance := &Lobby{
		BackendServer: NewBackendServer(server_type),
		db:            NewDB("game"),
	}
	instance.IServer = instance
	// 注册RPC
	instance.RegisterRPC(&LobbyRPCService{&RPCService{ServerInterface: instance}})
	// 初始化
	instance.Init()
	return instance
}

func (this *Lobby) Init() {
	this.BackendServer.Init()
	// 初始化数据库
	go this.db.Init()
}
