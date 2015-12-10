package common

import (
	"sync"
)

// 安全栈
type Stack struct {
	lock  *sync.RWMutex // 锁
	top   *Node         // 栈顶
	Count int           // 栈计数
}

func NewStack() *Stack {
	return &Stack{
		lock: new(sync.RWMutex),
	}
}

func (this *Stack) Push(elem interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	n := new(Node)
	n.Elem = elem
	n.Next = this.top
	this.top = n
	this.Count++
}

func (this *Stack) Pop() (elem interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.IsEmpty() {
		println("栈为空")
		return
	}
	n := this.top
	elem = n.Elem
	this.top = n.Next
	this.Count--
	return
}

func (this *Stack) IsEmpty() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return (this.top == nil)
}

func (this *Stack) StackTop() (elem interface{}) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if this.IsEmpty() {
		println("栈为空")
		return
	}
	elem = this.top.Elem
	return
}

func (this *Stack) PopAll() (elems []interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for {
		if this.top == nil {
			break
		}
		n := this.top
		e := n.Elem
		this.top = n.Next
		elems = append(elems, e)
		this.Count--
	}
	return
}
