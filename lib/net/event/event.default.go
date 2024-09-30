package event

import (
	"github.com/pkg/errors"
	xerror "xcore/lib/error"
	xlog "xcore/lib/log"
	xnethandler "xcore/lib/net/handler"
	xnetpacket "xcore/lib/net/packet"
	xnetremote "xcore/lib/net/remote"
	xruntime "xcore/lib/runtime"
)

type defaultEvent struct {
	eventChan chan<- interface{} // 待处理的事件
}

func NewDefaultEvent(eventChan chan<- interface{}) IEvent {
	return &defaultEvent{
		eventChan: eventChan,
	}
}

// Connect 连接
func (p *defaultEvent) Connect(handler xnethandler.IHandler, remote xnetremote.IRemote) error {
	select {
	case p.eventChan <- &Connect{
		IHandler: handler,
		IRemote:  remote,
	}:
	default:
		xlog.PrintfErr("push Connect failed with eventChan full. remote:%v", remote)
		return errors.WithMessage(xerror.ChannelFull, xruntime.Location())
	}
	return nil
}

// Disconnect 断开链接
func (p *defaultEvent) Disconnect(handler xnethandler.IHandler, remote xnetremote.IRemote) error {
	select {
	case p.eventChan <- &Disconnect{
		IHandler: handler,
		IRemote:  remote,
	}:
	default:
		xlog.PrintfErr("push Disconnect failed with eventChan full. remote:%v", remote)
		return errors.WithMessage(xerror.ChannelFull, xruntime.Location())
	}
	return nil
}

// Packet 数据包
func (p *defaultEvent) Packet(handler xnethandler.IHandler, remote xnetremote.IRemote, packet xnetpacket.IPacket) error {
	select {
	case p.eventChan <- &Packet{
		IHandler: handler,
		IRemote:  remote,
		IPacket:  packet,
	}:
	default:
		xlog.PrintfErr("push Packet failed with eventChan full. remote:%v packet:%v", remote, packet)
		return errors.WithMessage(xerror.ChannelFull, xruntime.Location())
	}
	return nil
}
