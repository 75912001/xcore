package tcp

import (
	"github.com/pkg/errors"
	xerror "xcore/lib/error"
	xnetpacket "xcore/lib/net/packet"
	xruntime "xcore/lib/runtime"
)

// serverOptions contains options to configure a server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type serverOptions struct {
	listenAddress    *string            // 127.0.0.1:8787
	eventChan        chan<- interface{} // 待处理的事件
	sendChanCapacity *uint32            // 发送 channel 大小
	packet           xnetpacket.IPacket
	connOptions      connOptions
	handler          IHandler
}

// NewServerOptions 新的ServerOptions
func NewServerOptions() *serverOptions {
	return new(serverOptions)
}

func (p *serverOptions) SetListenAddress(listenAddress string) *serverOptions {
	p.listenAddress = &listenAddress
	return p
}

func (p *serverOptions) SetEventChan(eventChan chan<- interface{}) *serverOptions {
	p.eventChan = eventChan
	return p
}

func (p *serverOptions) SetSendChanCapacity(sendChanCapacity uint32) *serverOptions {
	p.sendChanCapacity = &sendChanCapacity
	return p
}

func (p *serverOptions) SetPacket(packet xnetpacket.IPacket) *serverOptions {
	p.packet = packet
	return p
}

func (p *serverOptions) SetReadBuffer(readBuffer int) *serverOptions {
	p.connOptions.readBuffer = &readBuffer
	return p
}

func (p *serverOptions) SetWriteBuffer(writeBuffer int) *serverOptions {
	p.connOptions.writeBuffer = &writeBuffer
	return p
}

func (p *serverOptions) SetHandler(handler IHandler) *serverOptions {
	p.handler = handler
	return p
}

// mergeServerOptions combines the given *serverOptions into a single *serverOptions in a last one wins fashion.
// The specified options are merged with the existing options on the server, with the specified options taking
// precedence.
func mergeServerOptions(opts ...*serverOptions) *serverOptions {
	newOptions := NewServerOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.listenAddress != nil {
			newOptions.SetListenAddress(*opt.listenAddress)
		}
		if opt.eventChan != nil {
			newOptions.SetEventChan(opt.eventChan)
		}
		if opt.sendChanCapacity != nil {
			newOptions.SetSendChanCapacity(*opt.sendChanCapacity)
		}
		if opt.packet != nil {
			newOptions.SetPacket(opt.packet)
		}
		if opt.connOptions.readBuffer != nil {
			newOptions.SetReadBuffer(*opt.connOptions.readBuffer)
		}
		if opt.connOptions.writeBuffer != nil {
			newOptions.SetWriteBuffer(*opt.connOptions.writeBuffer)
		}
		if opt.handler != nil {
			newOptions.SetHandler(opt.handler)
		}
	}
	return newOptions
}

// 配置
func serverConfigure(opts *serverOptions) error {
	if opts.listenAddress == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.eventChan == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.sendChanCapacity == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.packet == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.handler == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	return nil
}
