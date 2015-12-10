package server

// 基础服务器
import (
	"common"
	"config/message_config"
	"config/server_config"
	. "framework/network"
	"log"
	"net"
	"net/rpc"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	LOGIC_INTERVAL        = 10                                     // 逻辑循环时间间隔，单位：毫秒
	SEND_REQUEST_INTERVAL = 20                                     // 发送请求的时间间隔，单位：毫秒
	SERVER_CONFIG         = "LiveServer\\bin\\config\\server.json" // 配置文件路径
	SERVER_DIR            = "liveserver"                           // 服务器所在目录
)

type BaseServer struct {
	// 相关模块
	*Connection                   // 连接
	*SessionManager               //会话管理
	*MessageManager               // 消息管理器
	*ProxyManager                 // 代理管理器
	rpcService      interface{}   // rpc服务
	logicTimer      *common.Timer // 逻辑循环的定时器
	sendTimer       *common.Timer // 发送请求定时器
	IServer                       // 服务器接口

	// Server自身的属性
	Type       string //服务器类型
	ID         int    // 服务器ID
	Addr       string // 分开的地址
	Port       string // 服务器端口
	ClientPort string // 用于客户端连接的端口
	TCPAddr    string // tcp地址，addr:port
	RPCAddr    string // rpc地址，addr:port
}

func NewBaseServer(server_type string) *BaseServer {
	return &BaseServer{
		Connection:     NewConnection(),
		SessionManager: NewSessionManager(),
		MessageManager: NewMessageManager(),
		ProxyManager:   NewProxyManager(),
		Type:           server_type,
	}
}

// 初始化服务器
func (this *BaseServer) Init() {
	// 读取配置
	this.LoadServerConfig()
	// 设置ID
	this.SetServerID()
	// 设置地址
	this.SetupAddr()
	// 进行连接远程服务器
	go this.StartConnectRemoteServer()
}

// 启动服务器
func (this *BaseServer) Run() {
	// 服务器以定时器为主循环
	this.logicTimer = common.NewTimer(time.Millisecond*LOGIC_INTERVAL, 0, false, this.MainLoop)
	// 发送请求定时器
	this.sendTimer = common.NewTimer(time.Millisecond*SEND_REQUEST_INTERVAL, 0, true, this.IServer.SendRequest)
	// 启动定时器
	this.logicTimer.Start()
	this.sendTimer.Start()
	//// 退出服务器循环直接退出进程
	//os.Exit(1)
}

// 服务器主循环，消息处理，逻辑暂时放在一起
func (this *BaseServer) MainLoop(timer *common.Timer, args ...interface{}) bool {
	// 处理请求
	this.IServer.HandleRequest()
	// 逻辑循环
	this.IServer.LogicLoop()
	runtime.Gosched()
	return true
}

// 启动TCP服务
func (this *BaseServer) StartTCPService(address, port string) {
	listener := this.CreateTCPListener(address, port)
	this.TCPAddr = address + ":" + port
	if listener == nil {
		return
	}
	bOk := this.AddListener("tcp", listener)
	if !bOk {
		return
	}
	// 添加监听
	this.AddTCPListenerCallback(this.IServer.TcpListenerCallback)
	// 监听为循环，需要开goroutine进行维护
	go this.RunTCPListen(listener)
}

// 注册rpc
func (this *BaseServer) RegisterRPC(rpc_handle interface{}) {
	this.rpcService = rpc_handle
	rpc.Register(rpc_handle)
}

// 启动RPC服务
func (this *BaseServer) StartRPCService(address, port string) {
	listener := this.CreateTCPListener(address, port)
	this.RPCAddr = address + ":" + port
	if listener == nil {
		return
	}
	bOk := this.AddListener("rpc", listener)
	if !bOk {
		return
	}
	// 监听为循环，需要开goroutine进行维护
	this.RunRPCListen(listener)
}

// tcp监听回调，此函数已用goroutine维护(此方法可以被重写)
func (this *BaseServer) TcpListenerCallback(conn net.Conn) {
	conn_time := time.Now()
	se := this.CreateSession(conn, &conn_time)
	if se == nil {
		conn.Close()
		this.DestorySession(se)
		return
	}
	this.AddSession(se)
	// 创建连接通道
	channel := this.BuildConnChannel(conn)
	// 接收连接消息请求
	this.AcceptTCPConnRequest(channel)
}

// 获取基本服务器
func (this *BaseServer) GetBaseServer() *BaseServer {
	return this
}

// 放置接受的消息
func (this *BaseServer) PushRevMessage(msg_list []*MessageObject) {
	for _, v := range msg_list {
		this.ReceiveMessage(v.Sid, v.Byte_Msg_body)
	}
}

// 放置发送的消息(后端发送到前端的消息)
func (this *BaseServer) PushSendMessage(msg_list []*MessageObject) {
	for _, v := range msg_list {
		this.MessageManager.Push(SEND_MSG, v)
	}
}

// 绑定代理
func (this *BaseServer) BindProxy(sid, server string, server_id int) {
	this.CreateProxy(sid, server, server_id)
}

// 读取服务器配置文件
func (this *BaseServer) LoadServerConfig() {
	// 全局配置只读一次
	if server_conf.IsReadConfig() {
		return
	}
	server_conf.ReadConfig(true)
	conf := server_conf.GetSingleton()
	conf_dir := common.GetDir()
	if conf_dir == "" {
		println("读取服务器配置文件失败!")
		os.Exit(1)
	}
	conf_slice := strings.Split(conf_dir, "\\")
	var index int = -1
	for i, v := range conf_slice {
		if strings.ToLower(v) == SERVER_DIR {
			index = i
			break
		}
	}
	if index == -1 {
		println("请把服务器拷贝到liveserver下，读取服务器配置文件失败!")
		os.Exit(1)
	}
	var conf_path string
	for i := 0; i < index; i++ {
		conf_path += conf_slice[i] + "\\"
	}
	conf_path += SERVER_CONFIG
	common.ReadJson(conf_path, conf)

	// 转换
	conf.ConvertToMap()
}

// 读取服务器启动参数
func (this *BaseServer) SetServerID() {
	args := os.Args
	if len(args) < 3 {
		this.ID = 1
		//println("ID = ", this.ID)
		return
	}
	id, err := strconv.ParseInt(args[2], 10, 0)
	if err != nil || id == 0 {
		this.ID = 1
	} else {
		this.ID = int(id)
	}
}

// 设置服务器地址
func (this *BaseServer) SetupAddr() {
	if !server_conf.IsReadConfig() {
		println("没有读取服务器配置！")
		os.Exit(1)
		return
	}
	//println("type = ", this.Type, " id = ", this.ID)
	if this.Type == "" || this.ID == 0 {
		println("没有设置服务器类型或者服务器ID")
		os.Exit(1)
		return
	}

	conf := server_conf.GetSingleton()
	conf_elem := conf.GetConfig(this.Type, this.ID)
	if conf_elem == nil {
		println("没有配置此服务器！")
		os.Exit(1)
		return
	}
	this.Addr = conf_elem.Host
	this.Port = conf_elem.Port
	this.RPCAddr = this.Addr + ":" + this.Port
	if conf_elem.Fronted {
		this.ClientPort = conf_elem.ClientPort
		this.TCPAddr = this.Addr + ":" + this.ClientPort
	}
}

// 进行RPC连接
func (this *BaseServer) StartConnectRemoteServer() {
	if !server_conf.IsReadConfig() {
		log.Println("没有读取服务器配置！")
		return
	}
	conf := server_conf.GetSingleton()
	for k := range conf.ConfigMap {
		if k == this.Type {
			continue
		}
		for i := range conf.ConfigMap[k] {
			server_conf := conf.ConfigMap[k][i]
			log.Println("开始连接", k, " ", server_conf.Id)
			this.ConnectRPCServer(server_conf.Host, server_conf.Port, server_conf.Type, server_conf.Id)
		}
	}
	log.Println("连接远程服务器成功!")
}

// 处理客户端请求
func (this *BaseServer) HandleRequest() {
	msg_list := this.MessageManager.PopAll(REV_MSG)
	msg_callback := msg_conf.MessageCallback()
	for _, v := range msg_list {
		msg_callback.Handle(v.ID, v.Sid, v.Proto_Msg_body)
	}
}

// rpc调用
func (this *BaseServer) RPCCall(server_type string, server_id int, fun string, arg_1 interface{}, arg_2 interface{}) {
	list := this.FindPRCClientsByDetail(server_type, server_id)
	for _, v := range list {
		go v.Call(fun, arg_1, arg_2)
	}
}

// 同步rpc调用
func (this *BaseServer) RPCSyncCall(server_type string, server_id int, fun string, arg_1 interface{}, arg_2 interface{}) {
	list := this.FindPRCClientsByDetail(server_type, server_id)
	for _, v := range list {
		v.Call(fun, arg_1, arg_2)
	}
}
