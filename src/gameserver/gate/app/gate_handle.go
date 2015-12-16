package gate

// gate服的消息请求处理
import (
	// "code.google.com/p/goprotobuf/proto"
	"common"
	"gamelogic/rpc"
	"msg_proto"

	"github.com/golang/protobuf/proto"
)

// 处理登录请求
func (this *Gate) RequestLogin(sid string, message proto.Message) {
	rev := message.(*msg_proto.RequestLogin)
	arg := &game_rpc.UserLoginInfo{
		Sid:       sid,
		LoginInfo: *rev,
		Error:     msg_proto.Error_Init_None,
	}
	// 获取账号验证结果
	go func(server *Gate) {
		result := &game_rpc.UserLoginInfo{}
		// 使用同步RPC
		this.RPCSyncCall("auth", 1, "AuthTest", arg, result)
		Type, Id, ip, port := server.GetFitConnector()
		switch result.Error {
		case msg_proto.Error_UnKnow_User:
			fallthrough
		case msg_proto.Error_Mismatch_Password:
			msg := &msg_proto.CmdResult{
				ErrorCode: &(result.Error),
			}
			server.SendMessage(msg_proto.MsgCmd_CmdResult_S, result.Sid, msg, false)
		case msg_proto.Error_Sucess:
			token := &msg_proto.Token{
				Sid:      proto.String(result.Sid),
				Ip:       proto.String(ip),
				Port:     proto.String(port),
				ConnTime: proto.Int64(common.GetTime()),
				Pid:      proto.String(result.UserId.Hex()),
			}
			var bOk *bool
			// 异步RPC
			server.RPCCall(Type, Id, "AddClientToken", token, bOk)
			// 发送消息给客户端
			server.SendMessage(msg_proto.MsgCmd_LoginToken_S, result.Sid, token, false)
		}
	}(this)
}
