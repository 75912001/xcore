package tcp

import (
	"github.com/pkg/errors"
	xerror "xcore/lib/error"
	xnetpacket "xcore/lib/net/packet"
	xruntime "xcore/lib/runtime"
)

type defaultSend struct {
}

func NewDefaultSend() ISend {
	return &defaultSend{}
}

func (p *defaultSend) Send(packet xnetpacket.IPacket) error {
	return errors.WithMessage(xerror.NotImplemented, xruntime.Location())
}
