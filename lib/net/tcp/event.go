package tcp

import (
	xnetpacket "xcore/lib/net/packet"
)

type IEvent interface {
	Connect(handler IHandler, remote IRemote) error                           // 链接 放入 事件中
	Disconnect(handler IHandler, remote IRemote) error                        // 断开链接 放入 事件中
	Packet(handler IHandler, remote IRemote, packet xnetpacket.IPacket) error // 数据包 放入 事件中
}

// Disconnect 事件数据-断开链接
type Disconnect struct {
	IHandler
	IRemote
}

// Connect 事件数据-链接成功
type Connect struct {
	IHandler
	IRemote
}

// Packet 事件数据-数据包
type Packet struct {
	IHandler
	IRemote
	xnetpacket.IPacket
}
