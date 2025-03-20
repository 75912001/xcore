package tcp

import (
	xpacket "github.com/75912001/xcore/lib/packet"
)

type IEvent interface {
	Connect(handler IHandler, remote IRemote) error                        // 链接 放入 事件中
	Disconnect(handler IHandler, remote IRemote) error                     // 断开链接 放入 事件中
	Packet(handler IHandler, remote IRemote, packet xpacket.IPacket) error // 数据包 放入 事件中
}

// Disconnect 事件数据-断开链接
type Disconnect struct {
	IHandler IHandler
	IRemote  IRemote
}

// Connect 事件数据-链接成功
type Connect struct {
	IHandler IHandler
	IRemote  IRemote
}

// Packet 事件数据-数据包
type Packet struct {
	IHandler IHandler
	IRemote  IRemote
	IPacket  xpacket.IPacket
}
