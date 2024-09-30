package remote

import (
	xnetconnect "xcore/lib/net/connect"
	xnetevent "xcore/lib/net/event"
	xnethandler "xcore/lib/net/handler"
	xnetsend "xcore/lib/net/send"
)

type IRemote interface {
	xnetsend.ISend
	IsConnect() bool
	Start(tcpOptions *xnetconnect.ConnOptions, event xnetevent.IEvent, handler xnethandler.IHandler)
	Stop()
	GetIP() string
}
