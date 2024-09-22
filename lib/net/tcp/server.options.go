package tcp

import (
	"github.com/pkg/errors"
	xerror "xcore/lib/error"
	xnetpacket "xcore/lib/net/packet"
	xruntime "xcore/lib/runtime"
)

// ServerOptions contains options to configure a server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type ServerOptions struct {
	listenAddress    *string            // 127.0.0.1:8787
	eventChan        chan<- interface{} // 待处理的事件
	sendChanCapacity *uint32            // 发送 channel 大小
	packet           xnetpacket.IPacket
	connOptions      ConnOptions
	handler          IHandler
}

// NewServerOptions 新的ServerOptions
func NewServerOptions() *ServerOptions {
	return new(ServerOptions)
}

func (p *ServerOptions) SetListenAddress(listenAddress string) *ServerOptions {
	p.listenAddress = &listenAddress
	return p
}

func (p *ServerOptions) SetEventChan(eventChan chan<- interface{}) *ServerOptions {
	p.eventChan = eventChan
	return p
}

func (p *ServerOptions) SetSendChanCapacity(sendChanCapacity uint32) *ServerOptions {
	p.sendChanCapacity = &sendChanCapacity
	return p
}

func (p *ServerOptions) SetPacket(packet xnetpacket.IPacket) *ServerOptions {
	p.packet = packet
	return p
}

func (p *ServerOptions) SetReadBuffer(readBuffer int) *ServerOptions {
	p.connOptions.readBuffer = &readBuffer
	return p
}

func (p *ServerOptions) SetWriteBuffer(writeBuffer int) *ServerOptions {
	p.connOptions.writeBuffer = &writeBuffer
	return p
}

func (p *ServerOptions) SetHandler(handler IHandler) *ServerOptions {
	p.handler = handler
	return p
}

// mergeServerOptions combines the given *ServerOptions into a single *ServerOptions in a last one wins fashion.
// The specified options are merged with the existing options on the server, with the specified options taking
// precedence.
func mergeServerOptions(opts ...*ServerOptions) *ServerOptions {
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
func serverConfigure(opts *ServerOptions) error {
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
