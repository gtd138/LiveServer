package network

// 消息模块
import (
	"code.google.com/p/goprotobuf/proto"
	"common"
	"config/message_config"
	"math"
	"msg_proto"
	//"reflect"
)

const (
	SEND_MSG = iota
	REV_MSG
)

type ByteMessage []byte

// 消息实体
type MessageObject struct {
	Sid            string           // 对应的会话ID
	ID             msg_proto.MsgCmd // 消息编号
	IsBc           bool             // 是否为广播消息
	Byte_Msg_body  ByteMessage      // Byte消息体
	Proto_Msg_body proto.Message    // proto消息体
}

func (this *MessageObject) Reset() {
	this.Sid = ""
	this.ID = msg_proto.MsgCmd_None
	this.Byte_Msg_body = this.Byte_Msg_body[:0]
	this.Proto_Msg_body.Reset()
}

// 消息管理器
type MessageManager struct {
	//msgRegistry *common.BeeMap //消息原型注册表，map[name]Message
	//msgIDMap    *common.BeeMap // 消息ID与名字映射表，map[name]id
	sendChannel *common.Queue // 发消息通道
	revChannel  *common.Queue // 接收到消息通道
}

func NewMessageManager() *MessageManager {
	msgMgr := &MessageManager{
		//msgRegistry: common.NewBeeMap(),
		//msgIDMap:    common.NewBeeMap(),
		sendChannel: common.NewQueue(),
		revChannel:  common.NewQueue(),
	}
	//msgMgr.RegisterAllMessage()
	// 注册所有消息
	msg_conf.RegisterMsg()
	return msgMgr
}

// 消息通道操作
func (this *MessageManager) Push(control_type int, msg_pack *MessageObject) {
	// 添加到发消息通道
	if control_type == SEND_MSG {
		this.sendChannel.EnQueue(msg_pack)
	} else {
		// 添加到接受消息
		this.revChannel.EnQueue(msg_pack)
	}
}

// 弹出所有消息
func (this *MessageManager) PopAll(control_type int) (msg_que []*MessageObject) {
	var temp []interface{}
	var count int
	if control_type == SEND_MSG {
		count = this.sendChannel.Count
		temp = this.sendChannel.DeQueueAll()
	} else {
		count = this.revChannel.Count
		temp = this.revChannel.DeQueueAll()
	}
	for i := 0; i < count; i++ {
		msg_que = append(msg_que, temp[i].(*MessageObject))
	}
	return
}

// 消息编码
func (this *MessageManager) Encode(handle msg_proto.MsgCmd, msg_body proto.Message) (encode_msg ByteMessage) {
	// 消息组成[255 00 00 22 22..]，其中255为消息总长度含自身，00 00位代表消息的handle，22 22为消息本体
	// 其中handle最大为9999
	// 编码消息handle
	if _, ok := msg_conf.MessageMap[handle]; !ok {
		println("不存在消息handle！....SendMessage...1")
		return
	}
	msgId := int(handle)
	var byte_handle []byte
	if msgId < 100 {
		byte_handle = []byte{0, byte(msgId)}
	} else {
		hight := byte(math.Floor((float64(msgId)) / 100))
		low := byte(msgId - int(hight)*100)
		byte_handle = []byte{hight, low}
	}

	// 消息体编码
	byte_msg_body, err := proto.Marshal(msg_body)
	if err != nil {
		println(err, " SendMessage....2")
		return
	}

	msg_len := len(byte_handle) + len(byte_msg_body) + 1
	encode_msg = append(encode_msg, byte(msg_len))
	encode_msg = append(encode_msg, byte_handle...)
	encode_msg = append(encode_msg, byte_msg_body...)
	return
}

// 消息解码
func (this *MessageManager) Decode(byte_msg ByteMessage) (msg_body proto.Message, msg_id msg_proto.MsgCmd, bOk bool) {
	// 消息组成[255 00 00 22 22..]，其中255为消息总长度含自身，00 00位代表消息的handle，22 22为消息本体
	// 其中handle最大为9999
	if len(byte_msg) < 5 {
		println("消息长度过短，长度小于等于6，消息解码失败！........MessageManager.decode........1")
		bOk = false
		return
	}

	msg_id, bOk = this.GetMessageHandle(byte_msg)
	if !bOk {
		println("解析消息名失败, 不存在此消息！")
		return
	}

	msg_prototype := msg_conf.MessageMap[msg_id].MsgObj

	defer func() {
		if r := recover(); r != nil {
			println("消息解码失败！返回正常流程！", r, "......MessageManager.decode........3")
		}
	}()

	// 生成新消息实例
	msg_body = proto.Clone(msg_prototype)
	byte_msg_body := byte_msg[3:]
	if len(byte_msg_body) == 0 {
		bOk = false
	} else {
		proto.Unmarshal(byte_msg_body, msg_body)
		bOk = true
	}
	return msg_body, msg_id, bOk
}

// 获取消息名字
func (this *MessageManager) GetMessageHandle(m ByteMessage) (msg_head msg_proto.MsgCmd, ok bool) {
	ok = false
	if len(m) < 1 {
		return
	}
	// 消息头
	byte_msg := m[1:3]
	msg_head = msg_proto.MsgCmd(byte_msg[0]*100 + byte_msg[1])
	_, ok = msg_conf.MessageMap[msg_head]
	return
}

// 发送消息
func (this *MessageManager) SendMessage(msg_id msg_proto.MsgCmd, sid string, msg_body proto.Message, isbc bool) {
	byte_msg := this.Encode(msg_id, msg_body)
	msg_obj := &MessageObject{
		Sid:           sid,
		ID:            msg_id,
		Byte_Msg_body: byte_msg,
		IsBc:          isbc,
	}
	this.Push(SEND_MSG, msg_obj)
}

// 接收消息
func (this *MessageManager) ReceiveMessage(sid string, byte_msg ByteMessage) {
	msg, msg_id, bOk := this.Decode(byte_msg)
	if !bOk {
		println("接收消息失败！")
		return
	}
	msg_obj := &MessageObject{
		Sid:            sid,
		ID:             msg_id,
		Proto_Msg_body: msg,
	}
	this.Push(REV_MSG, msg_obj)
}

// 创建byte消息Obj
func (this *MessageManager) CreateByteMessage(sid string, byte_msg ByteMessage) (msg_obj *MessageObject) {
	msg_obj = new(MessageObject)
	msg_obj.Sid = sid
	msg_obj.Byte_Msg_body = byte_msg
	return
}
