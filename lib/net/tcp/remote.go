package tcp

type IRemote interface {
	ISend
	IsConnect() bool
	Start(tcpOptions *ConnOptions, event IEvent, handler IHandler)
	Stop()
	GetIP() string
	GetDisconnectReason() DisconnectReason
	SetDisconnectReason(reason DisconnectReason)
}
