package server

import (
	"common"
	"config/message_config"
	. "framework/network"
	//"msg_proto"
	"net"
	//"strings"
	"time"
)

const (
	ROUTE_REV_FUNC         = "RouteMessage"
	ROUTE_REQUEST_INTERVAL = 20 // 路由请求的时间间隔，单位：毫秒
)

// 路由规则
type RouterPolicy func([]*ConnClient, *MessageObject)

// 前端服务器
type FrontendServer struct {
	*BaseServer
	routerPolicy map[string]RouterPolicy // 路由规则列表，格式：map[server_type]RouterPolicy
	routeTimer   *common.Timer           // 路由请求定时器
}

func NewFrontendServer(server_type string) *FrontendServer {
	instance := &FrontendServer{
		routerPolicy: make(map[string]RouterPolicy),
	}
	instance.BaseServer = NewBaseServer(server_type)
	return instance
}

// 初始化
func (this *FrontendServer) Init() {
	this.BaseServer.Init()
	// 添加tcp请求监听
	this.AddTCPRequestCallback(this.tcpRequestCallback)
	// 启动tpc服务
	this.StartTCPService(this.Addr, this.ClientPort)
	this.StartRPCService(this.Addr, this.Port)
	// 路由请求定时器
	this.routeTimer = common.NewTimer(time.Millisecond*ROUTE_REQUEST_INTERVAL, 0, true, this.RouteFlush)
}

func (this *FrontendServer) Run() {
	// 启动路由请求定时器
	this.routeTimer.Start()
	// 必须最后才能调用，以进入逻辑循环
	this.BaseServer.Run()
}

// tcp请求回调，用于处理通道接收到的消息
func (this *FrontendServer) tcpRequestCallback(conn net.Conn, rev [][]byte) {
	se := this.FindSessionByConn(conn)
	if se == nil {
		println("会话不存在！")
		return
	}
	for _, request := range rev {
		msg_obj := this.CreateByteMessage(se.SId, request)
		this.RouteRequest(msg_obj)
	}
}

// 添加路由规则
func (this *FrontendServer) AddRouterPolicy(server_type string, policy RouterPolicy) {
	if _, bOk := this.routerPolicy[server_type]; bOk {
		return
	}
	this.routerPolicy[server_type] = policy
}

// 执行路由
func (this *FrontendServer) DoRoute(server_type string, client_list []*ConnClient, msg_obj *MessageObject) {
	var policy RouterPolicy
	var bOk bool
	if policy, bOk = this.routerPolicy[server_type]; !bOk {
		this.doDefaultRoute(client_list, msg_obj)
	} else {
		policy(client_list, msg_obj)
	}
}

// 执行默认路由
func (this *FrontendServer) doDefaultRoute(client_list []*ConnClient, msg_obj *MessageObject) {
	// 把信息全部存到符合条件的所有客户端
	for i := range client_list {
		client_list[i].PushRequest(msg_obj)
	}
}

// 路由请求，仅仅把请求分类
func (this *FrontendServer) RouteRequest(msg_obj *MessageObject) {
	// 获取消息Handle
	msg_id, bOk := this.GetMessageHandle(msg_obj.Byte_Msg_body)
	//server_type := strings.Split(handle, "_")[0]
	var server_type string
	if !bOk {
		return
	}
	server_type, bOk = msg_conf.GetMsgServerType(msg_id)
	if !bOk {
		return
	}
	// 查找远程客户端列表
	rpc_client_list := this.FindRPCClient(server_type)
	if len(rpc_client_list) <= 0 {
		return
	}
	// 本地消息则放入接消息通道
	if server_type == this.Type {
		this.ReceiveMessage(msg_obj.Sid, msg_obj.Byte_Msg_body)
	} else {
		// 路由消息
		this.DoRoute(server_type, rpc_client_list, msg_obj)
	}
}

// 转发所有请求，真正转发所有消息，需要定时转发
func (this *FrontendServer) RouteFlush(t *common.Timer, args ...interface{}) bool {
	rpc_clients := this.GetRPCClients()
	for i := range rpc_clients {
		rpc_clients[i].RemoteCall(ROUTE_REV_FUNC)
	}
	return true
}

// 逻辑循环
func (this *FrontendServer) LogicLoop() {
}

// 发送消息到客户端
func (this *FrontendServer) SendRequest(t *common.Timer, args ...interface{}) bool {
	channel_lock := this.GetChannelLock()
	channel_lock.Lock()
	defer channel_lock.Unlock()
	channel_list := this.GetAllConnChannel()
	msg_list := this.MessageManager.PopAll(SEND_MSG)
	for _, v := range msg_list {
		for j, channel := range channel_list {
			if v.IsBc {
				channel_list[j].WriteBytes(v.Byte_Msg_body)
			} else {
				client_conn := this.FindSessionBySID(v.Sid)
				if client_conn == channel_list[j].GetConn() {
					channel.WriteBytes(v.Byte_Msg_body)
				}
			}
		}
	}

	for i := range channel_list {
		channel_list[i].ConnWrite()
	}
	return true
}

// 关闭会话
func (this *FrontendServer) CloseSession(se *Session) {
	this.CloseConnChannel(se.Conn)
	this.RemoveSession(se.SId)
}
