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

	OnCheckPacketLength OnCheckPacketLength // 检查长度,包头中
	OnCheckPacketLimit  OnCheckPacketLimit  // 检查包限制,限流
	OnConnect           OnConnect           // 连接成功
	OnUnmarshalPacket   OnUnmarshalPacket   // 解析数据包
	OnPacket            OnPacket            // 数据包
	OnDisconnect        OnDisconnect        // 断开连接
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

func (p *ServerOptions) SetOnCheckLength(onCheckLength OnCheckPacketLength) *ServerOptions {
	p.OnCheckPacketLength = onCheckLength
	return p
}

func (p *ServerOptions) SetOnCheckPacketLimit(onCheckPacketLimit OnCheckPacketLimit) *ServerOptions {
	p.OnCheckPacketLimit = onCheckPacketLimit
	return p
}

func (p *ServerOptions) SetOnConnect(onConnFromClient OnConnect) *ServerOptions {
	p.OnConnect = onConnFromClient
	return p
}

func (p *ServerOptions) SetOnUnmarshalPacket(onUnmarshalPacket OnUnmarshalPacket) *ServerOptions {
	p.OnUnmarshalPacket = onUnmarshalPacket
	return p
}

func (p *ServerOptions) SetOnPacket(onPacket OnPacket) *ServerOptions {
	p.OnPacket = onPacket
	return p
}

func (p *ServerOptions) SetOnDisconnect(onDisconnect OnDisconnect) *ServerOptions {
	p.OnDisconnect = onDisconnect
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
		if opt.OnCheckPacketLength != nil {
			newOptions.SetOnCheckLength(opt.OnCheckPacketLength)
		}
		if opt.OnCheckPacketLimit != nil {
			newOptions.SetOnCheckPacketLimit(opt.OnCheckPacketLimit)
		}
		if opt.OnConnect != nil {
			newOptions.SetOnConnect(opt.OnConnect)
		}
		if opt.OnUnmarshalPacket != nil {
			newOptions.SetOnUnmarshalPacket(opt.OnUnmarshalPacket)
		}
		if opt.OnPacket != nil {
			newOptions.SetOnPacket(opt.OnPacket)
		}
		if opt.OnDisconnect != nil {
			newOptions.SetOnDisconnect(opt.OnDisconnect)
		}
		if opt.connOptions.readBuffer != nil {
			newOptions.SetReadBuffer(*opt.connOptions.readBuffer)
		}
		if opt.connOptions.writeBuffer != nil {
			newOptions.SetWriteBuffer(*opt.connOptions.writeBuffer)
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
	if opts.OnConnect == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.OnUnmarshalPacket == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.OnPacket == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.OnDisconnect == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	return nil
}
