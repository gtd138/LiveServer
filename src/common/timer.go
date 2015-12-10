package common

import "time"

// 返回值用于停止定时器，返回false停止定时器
type TimerCallBack func(*Timer, ...interface{}) bool

type Timer struct {
	*time.Ticker               // Tiker
	interval     time.Duration // 触发间隔
	bSync        bool          // 是否在goroutine中运行
	bStop        bool          // 定时器停止标志位
	callback     TimerCallBack // 回调
	args         []interface{} // 回调的参数
	MaxCount     int           // 最大触发次数
	Count        int           // 当前触发次数
}

// 定时器
// interval:时间间隔以time.Nanosecond为单位
// count:触发最大次数，其中0为无限触发
// bSync:是否同步执行，是：则创建goroutine执行，否：会阻塞当前goroutine
// fun:触发回调
// args:回调传进的参数
func NewTimer(interval time.Duration, count int, bSync bool, fun TimerCallBack, args ...interface{}) *Timer {
	return &Timer{
		interval: interval,
		bSync:    bSync,
		callback: fun,
		args:     args,
		MaxCount: count,
	}
}

// 启动定时器
func (this *Timer) Start() {
	if this.bSync {
		go this.tick()
	} else {
		this.tick()
	}
}

// 停止定时器
func (this *Timer) Stop() {
	this.bStop = true
}

// 判断是否要停止
func (this *Timer) isStop() bool {
	if this.bStop {
		return true
	}

	// 达到最大计数，则停止
	if this.MaxCount <= 0 {
		return false
	} else {
		this.Count++
		if this.Count >= this.MaxCount {
			return true
		}
	}

	return false
}

// 定时器循环
func (this *Timer) tick() {
	this.Ticker = time.NewTicker(this.interval)
	this.bStop = false
TICK:
	for {
		select {
		case <-this.Ticker.C:
			if this.isStop() {
				break TICK
			}
			this.bStop = !this.callback(this, this.args...)
		}
	}
	this.Ticker.Stop()
}
