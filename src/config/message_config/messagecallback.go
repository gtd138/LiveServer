package msg_conf

import (
	"code.google.com/p/goprotobuf/proto"
	. "msg_proto"
)

// 单件
var instance *messageCallback

// 消息处理函数
type MessageHandleFunc func(sid string, message proto.Message)

type messageCallback struct {
	funcMap map[MsgCmd]MessageHandleFunc // 消息请求回调
}

// 获取消息回调单件
func MessageCallback() *messageCallback {
	if instance == nil {
		instance = &messageCallback{
			funcMap: make(map[MsgCmd]MessageHandleFunc),
		}
	}
	return instance
}

func (this *messageCallback) Register(msg_handle MsgCmd, fun MessageHandleFunc) {
	if _, bOk := this.funcMap[msg_handle]; bOk {
		println("重复注册消息回调！")
		return
	}
	this.funcMap[msg_handle] = fun
}

// 回调消息
func (this *messageCallback) Handle(msg_handle MsgCmd, sid string, message proto.Message) {
	if _, bOk := this.funcMap[msg_handle]; !bOk {
		println("不存在此消息回调")
		return
	}
	this.funcMap[msg_handle](sid, message)
}
