package master

// 监控服
import (
	. "framework/server"
)

type Master struct {
	*BackendServer
}

func NewMaster(server_type string) *Master {
	instance := &Master{
		BackendServer: NewBackendServer(server_type),
	}
	instance.IServer = instance
	// 注册RPC
	instance.RegisterRPC(&MasterRPCService{&RPCService{ServerInterface: instance}})
	// 初始化
	instance.Init()
	return instance
}
