package tcp

import (
	xpacket "xcore/lib/packet"
)

type ISend interface {
	Send(packet xpacket.IPacket) error
}
