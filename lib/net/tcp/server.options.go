package tcp

import (
	"github.com/pkg/errors"
	xerror "xcore/lib/error"
	xnetconnect "xcore/lib/net/connect"
	xruntime "xcore/lib/runtime"
)

// serverOptions contains options to configure a server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type serverOptions struct {
	listenAddress    *string            // 127.0.0.1:8787
	eventChan        chan<- interface{} // 待处理的事件
	sendChanCapacity *uint32            // 发送 channel 大小
	connOptions      xnetconnect.ConnOptions
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

func (p *serverOptions) SetReadBuffer(readBuffer int) *serverOptions {
	p.connOptions.ReadBuffer = &readBuffer
	return p
}

func (p *serverOptions) SetWriteBuffer(writeBuffer int) *serverOptions {
	p.connOptions.WriteBuffer = &writeBuffer
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
		if opt.connOptions.ReadBuffer != nil {
			newOptions.SetReadBuffer(*opt.connOptions.ReadBuffer)
		}
		if opt.connOptions.WriteBuffer != nil {
			newOptions.SetWriteBuffer(*opt.connOptions.WriteBuffer)
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
	return nil
}
