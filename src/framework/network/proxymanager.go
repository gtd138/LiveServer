package network

import (
	"common"
)

const (
	MAX_PROXY = 9999
)

// 代理，用于后端服务器
type Proxy struct {
	Sid          string // 会话Id
	FrontendType string //连接的前端服务类型
	FrontendID   int    // 连接的前端服务器ID
}

func (this *Proxy) Reset() {
	this.Sid = ""
	this.FrontendType = ""
	this.FrontendID = 0
}

// 代理管理器
type ProxyManager struct {
	*common.Pool
	proxyMap *common.BeeMap // 代理Map,proxyMap[sid]=Proxy
}

func NewProxyManager() *ProxyManager {
	return &ProxyManager{
		Pool:     common.NewPool(MAX_PROXY),
		proxyMap: common.NewBeeMap(),
	}
}

// 创建代理
func (this *ProxyManager) CreateProxy(sid, frontendType string, frontendID int) *Proxy {
	obj, bOk := this.Borrow()
	var p *Proxy
	if !bOk {
		p = new(Proxy)
	} else {
		obj.Reset()
		p = obj.(*Proxy)
	}
	p.Sid = sid
	p.FrontendType = frontendType
	p.FrontendID = frontendID
	this.proxyMap.Set(sid, p)
	return p
}

// 销毁代理
func (this *ProxyManager) DestroyProxy(sid string) {
	if !this.proxyMap.Check(sid) {
		return
	}
	p := this.proxyMap.Get(sid).(*Proxy)
	this.proxyMap.Delete(sid)
	this.GiveBack(p)
}

// 查找代理
func (this *ProxyManager) FindProxy(sid string) *Proxy {
	if !this.proxyMap.Check(sid) {
		return nil
	}
	return this.proxyMap.Get(sid).(*Proxy)
}
