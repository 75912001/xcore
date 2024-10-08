package tcp

import (
	xnetpacket "xcore/lib/net/packet"
)

type ISend interface {
	Send(packet xnetpacket.IPacket) error
}
