package connector

import (
	"common"
	"container/list"
	. "msg_proto"
	"sync"
	"time"
)

const (
	MAX_TOKEN_LIFE = 1 * 60 * 1000 // token的生命周期为1min
)

type TokenManager struct {
	tokenList *list.List // tonken列表
	lock      *sync.RWMutex
	lifeTimer *common.Timer // 生命周期定时器
}

func NewTokenManager() *TokenManager {
	tmgr := &TokenManager{
		tokenList: list.New(),
		lock:      new(sync.RWMutex),
	}
	tmgr.lifeTimer = common.NewTimer(time.Minute*MAX_TOKEN_LIFE, 0, true, tmgr.Check)
	return tmgr
}

// 启动Token管理器
func (this *TokenManager) Begin() {
	this.lifeTimer.Start()
}

// 检查token是否过期
func (this *TokenManager) Check(t *common.Timer, arg ...interface{}) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	cur_t := common.GetTime()
	e := this.tokenList.Front()
	for {
		if e == nil {
			break
		}
		next := e.Next()
		if e.Value.(*Token).GetConnTime()-cur_t >= int64(MAX_TOKEN_LIFE) {
			this.tokenList.Remove(e)
		}
		e = next
	}
	return true
}

// 添加token
func (this *TokenManager) Add(token *Token) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.tokenList.PushBack(token)
}

// 删除token
func (this *TokenManager) Remove(token *Token) {
	this.lock.Lock()
	defer this.lock.Unlock()
	var elem *list.Element
	for e := this.tokenList.Front(); e != nil; e.Next() {
		if e.Value.(*Token) == token {
			elem = e
			break
		}
	}
	this.tokenList.Remove(elem)
}
