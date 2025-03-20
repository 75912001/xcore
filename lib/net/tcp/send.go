package tcp

import (
	xpacket "github.com/75912001/xcore/lib/packet"
)

type ISend interface {
	Send(packet xpacket.IPacket) error
}
