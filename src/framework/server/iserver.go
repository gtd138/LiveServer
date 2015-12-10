package server

import (
	"common"
	"net"
)

// 服务器接口
type IServer interface {
	Init()                                                 // 初始化
	Run()                                                  // 启动
	LogicLoop()                                            // 具体服务器的逻辑循环
	SendRequest(t *common.Timer, args ...interface{}) bool // 发送请求
	HandleRequest()                                        // 处理请求
	GetBaseServer() *BaseServer                            // 获取基本服务器
	TcpListenerCallback(conn net.Conn)                     // tcp监听回调
}
