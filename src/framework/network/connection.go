package network

// 连接组件
import (
	"common"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

const (
	MAX_CONN_BUFF      = 4096
	RPC_RETRY_INTERVAL = 5  // rpc重连间隔，单位秒
	RPC_RETRY_TIME     = 10 // rpc重连次数
)

// 监听回调，用于接收新连接
type ConnectionListenerCallback func(conn net.Conn)
type TCPRequestCallback func(net.Conn, [][]byte)

type Connection struct {
	*common.Pool                                 // 连接通道池
	connChanList     []*ConnChannel              // 连接通道列表
	buildChannelLock *sync.RWMutex               // 创建通道锁
	listenerMap      map[string]*net.TCPListener // TCP监听，map[标识]listener
	tcpListenerCB    ConnectionListenerCallback  // TCP监听回调
	tcpRequestCB     TCPRequestCallback          // tcp请求处理回调
	rpcClientMap     []*ConnClient               // 连接远程服务器的客户端
}

func NewConnection() *Connection {
	return &Connection{
		Pool:             common.NewPool(MAX_SESSION),
		buildChannelLock: new(sync.RWMutex),
		listenerMap:      make(map[string]*net.TCPListener),
	}
}

// 创建网络监听，使用TCP协议
func (this *Connection) CreateTCPListener(addr, port string) *net.TCPListener {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr+":"+port)
	if err != nil {
		println("创建监听失败!(1)")
		return nil
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		println("创建监听失败!(2)")
		return nil
	}
	return listener
}

// 添加监听
func (this *Connection) AddListener(label string, listener *net.TCPListener) bool {
	_, bOk := this.listenerMap[label]
	if bOk {
		return !bOk
	}
	this.listenerMap[label] = listener
	return bOk
}

// 启动Tcp监听
func (this *Connection) RunTCPListen(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			println(err)
			continue
		}
		if this.tcpListenerCB == nil {
			continue
		}
		// 开goroutine处理监听回调
		go this.tcpListenerCB(conn)
	}
}

// 启动RPC监听
func (this *Connection) RunRPCListen(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(conn)
	}
}

// 添加TCP监听回调
func (this *Connection) AddTCPListenerCallback(cb ConnectionListenerCallback) {
	this.tcpListenerCB = cb
}

// 添加TCP请求回调
func (this *Connection) AddTCPRequestCallback(cb TCPRequestCallback) {
	this.tcpRequestCB = cb
}

// 设置连接通道
func (this *Connection) SetupConnChannel(conn net.Conn) {
	this.connChanList = append(this.connChanList, NewConnChannel(conn))
}

// 建立RPC连接
// 最后一个参数为重试连接，默认5s连接一次
func (this *Connection) ConnectRPCServer(addr, port, server_type string, server_id int) {
	fullAddr := addr + ":" + port
	conn_client := NewConnClient(server_type, server_id)
	timer := *common.NewTimer(time.Second*RPC_RETRY_INTERVAL, 0, false, this.startConnectRemoteServer, conn_client, fullAddr)
	log.Println("开始连接远程服务器 server = ", server_type, ", id =", server_id)
	timer.Start()
}

// 重连服务器
func (this *Connection) startConnectRemoteServer(t *common.Timer, arg ...interface{}) bool {
	t.Count = t.Count + 1
	log.Println("尝试连接次数 = ", t.Count)
	if len(arg) <= 0 {
		return false
	}
	client := arg[0].(*ConnClient)
	fullAddr := arg[1].(string)
	var err error
	client.Client, err = rpc.Dial("tcp", fullAddr)
	if err == nil {
		this.rpcClientMap = append(this.rpcClientMap, client)
		log.Println("连接成功, server = ", client.connServerType, ", id =", client.connServerID)
		return false
	} else {
		if t.Count >= RPC_RETRY_TIME {
			log.Println("连接数达到最大值，断开连接...")
			return false
		}
	}
	return (err != nil)
}

// 创建连接通道
func (this *Connection) BuildConnChannel(conn net.Conn) *ConnChannel {
	this.buildChannelLock.Lock()
	defer this.buildChannelLock.Unlock()
	obj, bOk := this.Borrow()
	var channel *ConnChannel
	if !bOk {
		channel = NewConnChannel(conn)
	} else {
		obj.Reset()
		channel = obj.(*ConnChannel)
	}
	this.connChanList = append(this.connChanList, channel)
	return channel
}

// 销毁连接通道
func (this *Connection) DestoryConnChannel(conn net.Conn) {
	this.buildChannelLock.Lock()
	defer this.buildChannelLock.Unlock()
	var index int
	for i := range this.connChanList {
		if this.connChanList[i].conn == conn {
			index = i
			break
		}
	}
	channel := this.connChanList[index]
	this.connChanList = append(this.connChanList[0:index], this.connChanList[index+1:]...)
	channel.conn.Close()
	this.GiveBack(channel)
}

// 关闭连接
func (this *Connection) CloseConnChannel(conn net.Conn) {
	this.buildChannelLock.Lock()
	defer this.buildChannelLock.Unlock()
	for i := range this.connChanList {
		if this.connChanList[i].conn == conn {
			this.connChanList[i].SetClose(true)
			break
		}
	}
}

// 接收连接请求
func (this *Connection) AcceptTCPConnRequest(conn_chan *ConnChannel) {
	println("接收连接请求......")
	defer this.DestoryConnChannel(conn_chan.conn)
	var request []byte
	for {
		// 通道关闭则关闭连接
		if conn_chan.IsClose() {
			break
		}

		request = make([]byte, MAX_CONN_BUFF)
		length, err := conn_chan.conn.Read(request)
		if err != nil {
			println(err)
			continue
		}

		if length <= 0 {
			continue
		}

		// 把消息添加到相关队列
		clip := this.cutConnRequest(request)
		if len(clip) <= 0 {
			continue
		}
		// 处理请求
		this.tcpRequestCB(conn_chan.conn, clip)
	}
}

// 剪切请求数组
func (this *Connection) cutConnRequest(request []byte) (clip_reuqest [][]byte) {
	request_len := len(request)
	if request_len <= 0 {
		println("请求长度为0!")
		return
	}
	for {
		request_len = len(request)
		if request_len <= 0 {
			break
		}
		// 请求头一般为消息长度
		req_head := request[0]
		if int(req_head) > request_len {
			println("消息实际长度过短，可能被截断！")
			break
		}
		cut := request[:req_head]
		clip_reuqest = append(clip_reuqest, cut)
		request = request[req_head:]
	}
	return
}

// 获取通道
func (this *Connection) GetAllConnChannel() []*ConnChannel {
	return this.connChanList
}

// 获取通道锁
func (this *Connection) GetChannelLock() *sync.RWMutex {
	return this.buildChannelLock
}

// 查找rpc端
func (this *Connection) FindRPCClient(server_type string) (client_list []*ConnClient) {
	for i, v := range this.rpcClientMap {
		if v.connServerType == server_type {
			client_list = append(client_list, this.rpcClientMap[i])
		}
	}
	return
}

// 获取rpc客户端列表
func (this *Connection) GetRPCClients() []*ConnClient {
	return this.rpcClientMap
}

// 获取rpc客户端列表，通过名称以及id
// id:为小于0时，默认返回跟其类型一致的端
func (this *Connection) FindPRCClientsByDetail(server_type string, id int) (clients []*ConnClient) {
	list := this.FindRPCClient(server_type)
	if id > 0 {
		for i := 0; i < len(list); i++ {
			if list[i].connServerID == id {
				clients = append(clients, list[i])
			}
		}
	} else {
		clients = list
	}
	return
}
