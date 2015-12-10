package network

import (
	"common"
	"net"
	"sync"
)

// 连接通道
type ConnChannel struct {
	conn net.Conn // 连接
	//rev_chan  *common.Queue // 接收通道
	send_chan *common.Queue // 发送通道
	bclosed   bool          // 是否要关闭
	slock     *sync.RWMutex // 状态锁
}

func NewConnChannel(conn net.Conn) *ConnChannel {
	return &ConnChannel{
		conn:    conn,
		bclosed: true,
		//rev_chan:  common.NewQueue(),
		send_chan: common.NewQueue(),
		slock:     new(sync.RWMutex),
	}
}

func (this *ConnChannel) Reset() {
	this.conn = nil
	//this.rev_chan.Clear()
	this.send_chan.Clear()
	this.bclosed = false
}

// 写字节
func (this *ConnChannel) WriteBytes(bytes []byte) {
	this.send_chan.EnQueue(bytes)
}

// 读取字节
func (this *ConnChannel) ReadBytes() [][]byte {
	var bytes [][]byte
	r := this.send_chan.DeQueueAll()
	for i := 0; i < len(r); i++ {
		bytes = append(bytes, r[i].([]byte))
	}
	this.send_chan.Clear()
	return bytes
}

// 设置通道是否关闭
func (this *ConnChannel) SetClose(bclose bool) {
	this.slock.Lock()
	defer this.slock.Unlock()
	this.bclosed = bclose
}

// 指出通道是否关闭
func (this *ConnChannel) IsClose() (bclose bool) {
	this.slock.RLock()
	defer this.slock.RUnlock()
	bclose = this.bclosed
	return
}

// 获取连接
func (this *ConnChannel) GetConn() net.Conn {
	return this.conn
}

// 真正写消息
func (this *ConnChannel) ConnWrite() {
	bytes := this.ReadBytes()
	for _, v := range bytes {
		this.conn.Write(v)
	}
}
