package tcp

import (
	"github.com/pkg/errors"
	xerror "xcore/lib/error"
	xlog "xcore/lib/log"
	xpacket "xcore/lib/packet"
	xruntime "xcore/lib/runtime"
)

type Event struct {
	eventChan chan<- interface{} // 待处理的事件
}

func NewEvent(eventChan chan<- interface{}) *Event {
	return &Event{
		eventChan: eventChan,
	}
}

// Connect 连接
func (p *Event) Connect(handler IHandler, remote IRemote) error {
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
func (p *Event) Disconnect(handler IHandler, remote IRemote) error {
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
func (p *Event) Packet(handler IHandler, remote IRemote, packet xpacket.IPacket) error {
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
