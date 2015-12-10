package network

// 会话组件
import (
	"common"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net"
	"time"
)

// 暂定10000个会话
const MAX_SESSION = 10000

// 会话
type Session struct {
	net.Conn            // 通用性连接
	SId      string     // 会话id
	bVerify  bool       // 是否已经验证
	ConnTime *time.Time // 上次链接时间
}

func (this *Session) Reset() {
	this.Conn = nil
	this.SId = ""
	this.ConnTime = nil
}

// 会话管理器
type SessionManager struct {
	*common.Pool                // 用于回收session对象
	SessionMap   *common.BeeMap // 会话列表，map[sid]=session
	Count        int            // 当前会话数量
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		SessionMap: common.NewBeeMap(),
		Pool:       common.NewPool(MAX_SESSION),
	}
}

// 随机一个会话id
func (this *SessionManager) randSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

// 创建会话
func (this *SessionManager) CreateSession(conn net.Conn, conn_time *time.Time) *Session {
	if this.Count >= MAX_SESSION {
		println("会话达到最大数量！")
		return nil
	}
	obj, bOk := this.Borrow()
	if !bOk {
		obj = new(Session)
	}
	se := obj.(*Session)
	se.Conn = conn
	se.ConnTime = conn_time
	se.SId = this.randSessionId()
	return se
}

// 添加会话
func (this *SessionManager) AddSession(se *Session) {
	if !this.SessionMap.Check(se.SId) {
		println("已存在此会话，id = ", se.SId)
		return
	}
	this.Count++
	this.SessionMap.Set(se.SId, se)
}

// 删除会话
func (this *SessionManager) RemoveSession(SId string) {
	var se *Session
	if !this.SessionMap.Check(SId) {
		println("不存在此会话，id = ", SId)
		return
	}
	se = this.SessionMap.Get(SId).(*Session)
	this.SessionMap.Delete(SId)
	// 归还给pool
	this.DestorySession(se)
	this.Count--
}

// 销毁会话
func (this *SessionManager) DestorySession(se *Session) {
	this.GiveBack(se)
}

// 通过Conn查找session
func (this *SessionManager) FindSessionByConn(conn net.Conn) (se *Session) {
	this.SessionMap.Foreach(func(k, v interface{}) bool {
		session := v.(*Session)
		if session.Conn == conn {
			se = session
			return true
		}
		return false
	})
	return
}

// 通过Sid查找session
func (this *SessionManager) FindSessionBySID(sid string) (se *Session) {
	if !this.SessionMap.Check(sid) {
		return
	}
	se = this.SessionMap.Get(sid).(*Session)
	return
}
