package tcp

import (
	"github.com/pkg/errors"
	xerror "xcore/lib/error"
	xnetpacket "xcore/lib/net/packet"
	xruntime "xcore/lib/runtime"
)

// clientOptions contains options to configure a server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type clientOptions struct {
	serverAddress    *string            // 服务端的地址 e.g.:127.0.0.1:8787
	eventChan        chan<- interface{} // 外部传递的事件处理管道.连接的事件会放入该管道,以供外部处理
	sendChanCapacity *uint32            // 发送管道容量
	packet           xnetpacket.IPacket
	connOptions      connOptions
	handler          IHandler
}

// NewClientOptions 新的ClientOptions
func NewClientOptions() *clientOptions {
	return new(clientOptions)
}

func (p *clientOptions) SetReadBuffer(readBuffer int) *clientOptions {
	p.connOptions.readBuffer = &readBuffer
	return p
}

func (p *clientOptions) SetWriteBuffer(writeBuffer int) *clientOptions {
	p.connOptions.writeBuffer = &writeBuffer
	return p
}

func (p *clientOptions) SetAddress(address string) *clientOptions {
	p.serverAddress = &address
	return p
}

func (p *clientOptions) SetEventChan(eventChan chan<- interface{}) *clientOptions {
	p.eventChan = eventChan
	return p
}

func (p *clientOptions) SetSendChanCapacity(sendChanCapacity uint32) *clientOptions {
	p.sendChanCapacity = &sendChanCapacity
	return p
}

func (p *clientOptions) SetPacket(packet xnetpacket.IPacket) *clientOptions {
	p.packet = packet
	return p
}

func (p *clientOptions) SetHandler(handler IHandler) *clientOptions {
	p.handler = handler
	return p
}

// mergeClientOptions combines the given *clientOptions into a single *clientOptions in a last one wins fashion.
// The specified options are merged with the existing options on the server, with the specified options taking
// precedence.
func mergeClientOptions(opts ...*clientOptions) *clientOptions {
	newOptions := NewClientOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.serverAddress != nil {
			newOptions.SetAddress(*opt.serverAddress)
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
func clientConfigure(opts *clientOptions) error {
	if opts.serverAddress == nil {
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
