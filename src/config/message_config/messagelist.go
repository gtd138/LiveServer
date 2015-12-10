package msg_conf

import (
	"code.google.com/p/goprotobuf/proto"
	. "msg_proto"
)

// 消息定义
type MsgDef struct {
	ServerType string        // 发送的服务器类型，客户端的则为"client"
	MsgObj     proto.Message // 消息原型
}

var MessageMap map[MsgCmd]MsgDef

// 获取消息对应的服务器类型
func GetMsgServerType(msg_id MsgCmd) (server_type string, bOk bool) {
	msg_obj, ok := MessageMap[msg_id]
	bOk = ok
	if !ok {
		return
	}
	server_type = msg_obj.ServerType
	return
}

// 注册消息
func RegisterMsg() {
	MessageMap = make(map[MsgCmd]MsgDef)

	// 以下为消息注册表
	// 服务器结果
	MessageMap[MsgCmd_CmdResult_S] = MsgDef{"client", &CmdResult{}}
	MessageMap[MsgCmd_LoginToken_S] = MsgDef{"client", &Token{}}

	// 客户端请求
	MessageMap[MsgCmd_RequestLogin_C] = MsgDef{"gate", &RequestLogin{}}
}
