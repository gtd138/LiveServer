package database

// 数据库缓存服
import (
	. "framework/server"
)

type DataBase struct {
	*BackendServer
}

func NewDataBase(server_type string) *DataBase {
	instance := &DataBase{
		BackendServer: NewBackendServer(server_type),
	}
	instance.IServer = instance
	// 注册RPC
	instance.RegisterRPC(&DataBaseRPCService{&RPCService{ServerInterface: instance}})
	// 初始化
	instance.Init()
	return instance
}
