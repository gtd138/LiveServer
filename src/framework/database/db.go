package db

import (
	"common"
	. "config/db_config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"strings"
	"time"
)

const (
	DBPATH = "LiveServer\\bin\\config\\db.json"
	DB_DIR = "liveserver"

	DURTIME    = 5  // 重连间隔5s
	RETRY_TIME = 10 // 重连10次

	DB_PROCESS_DURTIME = 15 // 数据库操作处理为15分钟一次
)

// 数据队列元素
type DBQueueElem struct {
	Table string        // 表名
	Id_   bson.ObjectId `bson:"_id"` // 数据ID
	Data  interface{}   // 数据
}

type DBInitFinCallback func()

// 数据库模块
type DB struct {
	*DBConfig                   // 数据库配置
	*mgo.Session                // 数据会话
	*mgo.Database               // 连接的数据库
	insertQueue   *common.Queue // 插入数据队列
	fixQueue      *common.Queue // 修改数据队列
	delQueue      *common.Queue // 删除数据队列
	retryTimer    *common.Timer // 重连定时器
	processTimer  *common.Timer // 数据库操作定时器

	// 数据库的一些属性
	Host string // 数据库地址
	Name string // 数据库名

	// 回调
	InitCallback DBInitFinCallback // 初始化结束回调
}

func NewDB(name string) *DB {
	return &DB{
		DBConfig:    &DBConfig{},
		insertQueue: common.NewQueue(),
		fixQueue:    common.NewQueue(),
		delQueue:    common.NewQueue(),
		Name:        name,
	}
}

// 初始化数据库，请用goroutine
func (this *DB) Init() {
	this.loadDBConfig()
	this.connectDB()
	this.setupDB()
	if this.InitCallback != nil {
		this.InitCallback()
	}
}

// 读取配置文件
func (this *DB) loadDBConfig() {
	conf_dir := common.GetDir()
	if conf_dir == "" {
		println("读数据库配置文件失败!")
		os.Exit(1)
	}
	conf_slice := strings.Split(conf_dir, "\\")
	var index int = -1
	for i, v := range conf_slice {
		if strings.ToLower(v) == DB_DIR {
			index = i
			break
		}
	}
	if index == -1 {
		println("请把服务器拷贝到liveserver下，读取数据库配置文件失败!")
		os.Exit(1)
	}
	var conf_path string
	for i := 0; i < index; i++ {
		conf_path += conf_slice[i] + "\\"
	}
	conf_path += DBPATH
	common.ReadJson(conf_path, this.DBConfig)
}

// 连接数据库
func (this *DB) connectDB() {
	this.Host = this.GetConfig(this.Name).Host
	// 开启goroutine进行连接数据库
	this.retryTimer = common.NewTimer(time.Second*DURTIME, 0, false, this.retry)
	log.Println("开始连接数据库...")
	this.retryTimer.Start()
}

// 重连DB
func (this *DB) retry(t *common.Timer, args ...interface{}) bool {
	var err error
	this.Session, err = mgo.Dial(this.Host)
	log.Println("连接数据库次数 = ", t.Count+1)
	if err != nil {
		if t.Count >= RETRY_TIME {
			log.Println("连接数据库次数过多，连接数据库失败！")
			return false
		}
		return true
	}
	log.Println("连接数据库成功！")
	return false
}

// 设置DB
func (this *DB) setupDB() {
	this.Session.SetMode(mgo.Monotonic, true)
	this.Database = this.Session.DB(this.Name)
	// 通过定时器驱动数据模块
	this.processTimer = common.NewTimer(time.Minute*DB_PROCESS_DURTIME, 0, true, this.process)
	this.processTimer.Start()
}

// 更改使用的数据库(IP相同)
func (this *DB) ChangeDB(db_name string) {
	this.Database = this.Session.DB(db_name)
}

// 获取数据库某张表所有数值
// table:表名
// result:数据结构
func (this *DB) FindAll(table string, result interface{}) (err error) {
	err = this.Database.C(table).Find(&bson.M{}).All(result)
	return
}

// note:以下操作，不会马上进行实际的数据库操作
// 插入数据
func (this *DB) Insert(table string, data interface{}) {
	e := &DBQueueElem{
		Table: table,
		Data:  data,
	}
	this.insertQueue.EnQueue(e)
}

// 更新数据
func (this *DB) Update(table string, id bson.ObjectId, data interface{}) {
	e := &DBQueueElem{
		Table: table,
		Id_:   id,
		Data:  data,
	}
	this.fixQueue.EnQueue(e)
}

// 删除数据
func (this *DB) Remove(table string, id bson.ObjectId) {
	e := &DBQueueElem{
		Table: table,
		Id_:   id,
	}
	this.delQueue.EnQueue(e)
}

// 以下为实际进行数据操作
func (this *DB) process(t *common.Timer, args ...interface{}) bool {
	// 插入
	q := this.insertQueue.DeQueueAll()
	for i := 0; i < len(q); i++ {
		e := q[i].(*DBQueueElem)
		this.Database.C(e.Table).Insert(e.Data)
	}
	// 修改
	q = q[:0]
	q = this.fixQueue.DeQueueAll()
	for i := 0; i < len(q); i++ {
		e := q[i].(*DBQueueElem)
		this.Database.C(e.Table).Update(bson.M{"_id": bson.ObjectIdHex(e.Id_.Hex())}, e.Data)
	}
	// 删除
	q = q[:0]
	q = this.delQueue.DeQueueAll()
	for i := 0; i < len(q); i++ {
		e := q[i].(*DBQueueElem)
		this.Database.C(e.Table).Remove(bson.M{"_id": bson.ObjectIdHex(e.Id_.Hex())})
	}
	return true
}
