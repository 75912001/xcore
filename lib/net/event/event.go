package event

import (
	xnethandler "xcore/lib/net/handler"
	xnetpacket "xcore/lib/net/packet"
	xnetremote "xcore/lib/net/remote"
)

type IEvent interface {
	Connect(handler xnethandler.IHandler, remote xnetremote.IRemote) error                           // 链接 放入 事件中
	Disconnect(handler xnethandler.IHandler, remote xnetremote.IRemote) error                        // 断开链接 放入 事件中
	Packet(handler xnethandler.IHandler, remote xnetremote.IRemote, packet xnetpacket.IPacket) error // 数据包 放入 事件中
}

// Disconnect 事件数据-断开链接
type Disconnect struct {
	xnethandler.IHandler
	xnetremote.IRemote
}

// Connect 事件数据-链接成功
type Connect struct {
	xnethandler.IHandler
	xnetremote.IRemote
}

// Packet 事件数据-数据包
type Packet struct {
	xnethandler.IHandler
	xnetremote.IRemote
	xnetpacket.IPacket
}
