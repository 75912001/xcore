package tcp

import (
	xerror "github.com/75912001/xcore/lib/error"
	xruntime "github.com/75912001/xcore/lib/runtime"
	"github.com/pkg/errors"
)

// serviceOptions contains options to configure a Service instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type serviceOptions struct {
	listenAddress    *string            // 127.0.0.1:8787
	eventChan        chan<- interface{} // 待处理的事件
	sendChanCapacity *uint32            // 发送 channel 大小
	connOptions      ConnOptions
}

// NewServerOptions 新的ServerOptions
func NewServerOptions() *serviceOptions {
	return new(serviceOptions)
}

func (p *serviceOptions) SetListenAddress(listenAddress string) *serviceOptions {
	p.listenAddress = &listenAddress
	return p
}

func (p *serviceOptions) SetEventChan(eventChan chan<- interface{}) *serviceOptions {
	p.eventChan = eventChan
	return p
}

func (p *serviceOptions) SetSendChanCapacity(sendChanCapacity uint32) *serviceOptions {
	p.sendChanCapacity = &sendChanCapacity
	return p
}

func (p *serviceOptions) SetReadBuffer(readBuffer int) *serviceOptions {
	p.connOptions.ReadBuffer = &readBuffer
	return p
}

func (p *serviceOptions) SetWriteBuffer(writeBuffer int) *serviceOptions {
	p.connOptions.WriteBuffer = &writeBuffer
	return p
}

// mergeServiceOptions combines the given *serviceOptions into a single *serviceOptions in a last one wins fashion.
// The specified options are merged with the existing options on the Service, with the specified options taking
// precedence.
func mergeServiceOptions(opts ...*serviceOptions) *serviceOptions {
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
func serviceConfigure(opts *serviceOptions) error {
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
