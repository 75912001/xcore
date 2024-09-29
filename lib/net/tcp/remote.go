package tcp

type IRemote interface {
	ISend
	IsConnect() bool
	Start(tcpOptions *connOptions, event IEvent, handler IHandler)
	Stop()
	GetIP() string
}
