package network

// 用于处于RPC调用
type IRPCHandler interface {
	RouteMessage([]*MessageObject, *bool) error
}
