package common

// 对象池实现

// pool回收物
type IPoolObject interface {
	Reset()
}

type Pool struct {
	*Stack
	max_size int
}

func NewPool(size int) *Pool {
	return &Pool{
		Stack:    NewStack(),
		max_size: size,
	}
}

// 从对象池中取出已有对象，没有对象时，需要手动new对象
func (this *Pool) Borrow() (pool_obj IPoolObject, bOk bool) {
	bOk = false
	if !this.IsEmpty() {
		pool_obj = this.Pop().(IPoolObject)
		pool_obj.Reset()
		bOk = true
	}
	return
}

// 归还对象
func (this *Pool) GiveBack(pool_obj IPoolObject) {
	if this.Count >= this.max_size {
		return
	}
	this.Push(pool_obj)
}
