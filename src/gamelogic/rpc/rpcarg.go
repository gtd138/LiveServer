package game_rpc

import (
	"gopkg.in/mgo.v2/bson"
	"msg_proto"
)

// 连接数量
type ConnState struct {
	Name string // 服务器名
	Type string // 服务器类型
	Id   int    // 服务器Id
	Num  int    // 连接数
}

// 用户信息
type UserLoginInfo struct {
	Sid       string
	LoginInfo msg_proto.RequestLogin // 登录信息
	UserId    bson.ObjectId          // 数据库ID
	Error     msg_proto.Error        // 错误码
}
