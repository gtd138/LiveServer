package connector

import (
	"common"
	"container/list"
	"framework/network"
	"math"
	"sync"
	"time"
)

const (
	CHECK_DUTIME = 1 // 定时检查时间，为1分钟
)

type TimeoutCallback func(se *network.Session)

type SessionCheck struct {
	seList      *list.List      // 待验证的会话列表
	seCheckLock *sync.RWMutex   // 验证会话锁
	checkTimer  *common.Timer   // 定时检查
	SeTimeout   TimeoutCallback // 超时回调
}

func NewSessionCheck() *SessionCheck {
	sc := &SessionCheck{
		seList:      list.New(),
		seCheckLock: new(sync.RWMutex),
	}
	sc.checkTimer = common.NewTimer(time.Minute*CHECK_DUTIME, 0, true, sc.Check)
	return sc
}

// 开始检测
func (this *SessionCheck) Begin() {
	this.checkTimer.Start()
}

// 检查会话
func (this *SessionCheck) Check(t *common.Timer, arg ...interface{}) bool {
	this.seCheckLock.Lock()
	defer this.seCheckLock.Unlock()
	e := this.seList.Front()
	cur_t := float64(time.Now().UnixNano())
	for {
		if e == nil {
			break
		}
		se := e.Value.(*network.Session)
		next := e.Next()
		dt := (float64(se.ConnTime.UnixNano()) - cur_t) / math.Pow(float64(10), float64(9))
		if dt/float64(60) >= float64(CHECK_DUTIME) {
			this.seList.Remove(e)
			if this.SeTimeout != nil {
				this.SeTimeout(se)
			}
		}
		e = next
	}
	return true
}

// 添加待测会话
func (this *SessionCheck) Add(se *network.Session) {
	this.seCheckLock.Lock()
	defer this.seCheckLock.Unlock()
	this.seList.PushBack(se)
}

// 移除会话
func (this *SessionCheck) Remove(se *network.Session) {
	this.seCheckLock.Lock()
	defer this.seCheckLock.Unlock()
	var elem *list.Element
	l := this.seList
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value.(*network.Session) == se {
			elem = e
			break
		}
	}
	this.seList.Remove(elem)
}
