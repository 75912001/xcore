package service

type IHandler interface {
	// 事件处理
	Handle() error
}
