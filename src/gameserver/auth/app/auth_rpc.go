package auth

import (
	. "framework/server"
	"gamelogic/rpc"
	"msg_proto"
)

type AuthRPCService struct {
	*RPCService
}

// 验证用户密码
func (this *AuthRPCService) AuthTest(arg *game_rpc.UserLoginInfo, result *game_rpc.UserLoginInfo) error {
	server := this.ServerInterface.(*Auth)
	result = &(*arg)
	user := arg.LoginInfo.GetUsername()
	if server.userMap.Check(user) {
		pw := server.userMap.Get(user).(*AuthData)
		// 验证密码成功
		if pw.Password == arg.LoginInfo.GetPassword() {
			result.Error = msg_proto.Error_Sucess
			result.UserId = server.userMap.Get(user).(*AuthData).UserID
		} else {
			// 密码不匹配
			result.Error = msg_proto.Error_Mismatch_Password
		}
	} else {
		// 不存在用户
		result.Error = msg_proto.Error_UnKnow_User
	}
	return nil
}
