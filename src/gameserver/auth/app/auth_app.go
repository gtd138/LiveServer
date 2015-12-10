package auth

// 第三方平台验证服
import (
	"common"
	. "framework/database"
	. "framework/server"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	AUTH_DB        = "auth" // 验证数据库名
	AUTH_DB_TABLE  = "user" // 数据表名
	UPDATE_DURTIME = 15     // 分钟更新一下数据
)

type AuthData struct {
	Id_      bson.ObjectId `bson:"_id"` // 数据ID
	User     string        // 用户名
	Password string        // 密码
	UserID   bson.ObjectId // 关联的用户数据ID
}

type Auth struct {
	*BackendServer
	db        *DB // 数据库
	dataTimer *common.Timer
	userMap   *common.BeeMap // 用户列表，map[username]=AuthData
}

func NewAuth(server_type string) *Auth {
	instance := &Auth{
		BackendServer: NewBackendServer(server_type),
		db:            NewDB(AUTH_DB),
		userMap:       common.NewBeeMap(),
	}
	instance.IServer = instance
	// 注册RPC
	instance.RegisterRPC(&AuthRPCService{&RPCService{ServerInterface: instance}})
	// 初始化
	instance.Init()
	return instance
}

func (this *Auth) Init() {
	this.BackendServer.Init()
	// 添加初始化完成后事件
	this.db.InitCallback = this.DBFinInitCallback
	// 初始化数据库
	go this.db.Init()
}

// 数据库模块启动完成回调
func (this *Auth) DBFinInitCallback() {
	this.dataTimer = common.NewTimer(time.Minute*UPDATE_DURTIME, 0, true, this.UpdateData)
	this.dataTimer.Start()
	// 第一次就开始更新
	this.UpdateData(this.dataTimer)
}

// 定时更新数据
func (this *Auth) UpdateData(t *common.Timer, arg ...interface{}) bool {
	var data_list []AuthData
	this.db.FindAll(AUTH_DB_TABLE, &data_list)
	for _, v := range data_list {
		this.userMap.Set(v.User, &v)
	}
	return true
}
