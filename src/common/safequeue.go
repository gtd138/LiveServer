package common

import (
	"sync"
)

// 安全队列
type Queue struct {
	front, rear *Node         // 头尾节点
	Count       int           // 长度
	lock        *sync.RWMutex // 读写锁
}

func NewQueue() *Queue {
	q := new(Queue)
	q.front = new(Node)
	q.rear = q.front
	q.lock = new(sync.RWMutex)
	return q
}

// 插入队列
func (this *Queue) EnQueue(e interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	q := new(Node)
	q.Elem = e
	this.rear.Next = q
	this.rear = q
	this.Count++
}

// 出队
func (this *Queue) DeQueue() (e interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.IsEmpty() {
		return
	}
	q := this.front.Next
	e = q.Elem
	if this.front.Next == this.rear {
		this.rear = this.front
	}
	this.front.Next = q.Next
	this.Count--
	return
}

// 全部出队
func (this *Queue) DeQueueAll() (es []interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for {
		if this.IsEmpty() {
			break
		}
		q := this.front.Next
		e := q.Elem
		if this.front.Next == this.rear {
			this.rear = this.front
		}
		this.front.Next = q.Next
		this.Count--
		es = append(es, e)
	}
	return
}

// 是否为空
func (this *Queue) IsEmpty() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return (this.front == this.rear)
}

// 获取队头
func (this *Queue) TopQueue() (e interface{}) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if this.IsEmpty() {
		return
	}
	return this.front.Next.Elem
}

// 清空
func (this *Queue) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.front.Next = nil
	this.rear = this.front
}
