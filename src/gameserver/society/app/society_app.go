package society

// 社会服
import (
	. "framework/server"
)

type Society struct {
	*BackendServer
}

func NewSociety(server_type string) *Society {
	instance := &Society{
		BackendServer: NewBackendServer(server_type),
	}
	instance.IServer = instance
	// 注册RPC
	instance.RegisterRPC(&SocietyRPCService{&RPCService{ServerInterface: instance}})
	// 初始化
	instance.Init()
	return instance
}
