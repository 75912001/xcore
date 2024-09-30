package send

import (
	"github.com/pkg/errors"
	xerror "xcore/lib/error"
	xnetpacket "xcore/lib/net/packet"
	xruntime "xcore/lib/runtime"
)

type ISend interface {
	Send(packet xnetpacket.IPacket) error
}

type defaultSend struct {
}

func NewDefaultSend() ISend {
	return &defaultSend{}
}

func (p *defaultSend) Send(packet xnetpacket.IPacket) error {
	return errors.WithMessage(xerror.NotImplemented, xruntime.Location())
}
