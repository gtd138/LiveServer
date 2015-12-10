package network

import (
	"net/rpc"
	"sync"
)

type ConnClient struct {
	*rpc.Client                        // rpc客户端
	connServerType    string           // 连接的服务端类型名
	connServerID      int              // 连接的服务器ID
	rpcRequestChannel []*MessageObject // rpc请求通道
	lock              *sync.RWMutex    // 读取请求锁
}

func NewConnClient(server_type string, server_id int) *ConnClient {
	return &ConnClient{
		//Client:         client,
		connServerType: server_type,
		connServerID:   server_id,
		lock:           new(sync.RWMutex),
	}
}

// 放置请求
func (this *ConnClient) PushRequest(request *MessageObject) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.rpcRequestChannel = append(this.rpcRequestChannel, request)
}

// 把所有请求发送到远程服务器
func (this *ConnClient) RemoteCall(funcname string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	request_channel := this.rpcRequestChannel[0:]
	var bOk bool
	go this.Client.Call(funcname, request_channel, &bOk)
}

// 获取远程服务器类型
func (this *ConnClient) GetRemoteServerType() string {
	return this.connServerType
}

// 获取远程服务器ID
func (this *ConnClient) GetRemoteServerID() int {
	return this.connServerID
}
