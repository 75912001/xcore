package tcp

import (
	xnetpacket "xcore/lib/net/packet"
)

type IEvent interface {
	Connect(handler IHandler, remote IRemote) error                           // 链接 放入 事件中
	Disconnect(handler IHandler, remote IRemote) error                        // 断开链接 放入 事件中
	Packet(handler IHandler, remote IRemote, packet xnetpacket.IPacket) error // 数据包 放入 事件中
}

// EventDisconnect 事件-断开链接
type EventDisconnect struct {
	IHandler
	IRemote
}

// EventConnect 事件-链接成功
type EventConnect struct {
	IHandler
	IRemote
}

// EventPacket 事件-数据包
type EventPacket struct {
	IHandler
	IRemote
	xnetpacket.IPacket
}
