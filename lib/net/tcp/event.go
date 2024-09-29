package tcp

import (
	"github.com/pkg/errors"
	xerror "xcore/lib/error"
	xlog "xcore/lib/log"
	xruntime "xcore/lib/runtime"
)

type IEvent interface {
	Connect(handler IHandler, remote IRemote) error                       // 链接 放入 事件中
	Disconnect(handler IHandler, remote IRemote) error                    // 断开链接 放入 事件中
	Packet(handler IHandler, remote IRemote, packet *DefaultPacket) error // 数据包 放入 事件中
}

type defaultEvent struct {
	eventChan chan<- interface{} // 待处理的事件
}

func newDefaultEvent(eventChan chan<- interface{}) IEvent {
	return &defaultEvent{
		eventChan: eventChan,
	}
}

// Connect 连接
func (p *defaultEvent) Connect(handler IHandler, remote IRemote) error {
	select {
	case p.eventChan <- &EventConnect{
		IHandler: handler,
		IRemote:  remote,
	}:
	default:
		xlog.PrintfErr("push EventConnect failed with eventChan full. remote:%v", remote)
		return errors.WithMessage(xerror.ChannelFull, xruntime.Location())
	}
	return nil
}

// Disconnect 断开链接
func (p *defaultEvent) Disconnect(handler IHandler, remote IRemote) error {
	select {
	case p.eventChan <- &EventDisconnect{
		IHandler: handler,
		IRemote:  remote,
	}:
	default:
		xlog.PrintfErr("push EventDisconnect failed with eventChan full. remote:%v", remote)
		return errors.WithMessage(xerror.ChannelFull, xruntime.Location())
	}
	return nil
}

// Packet 数据包
func (p *defaultEvent) Packet(handler IHandler, remote IRemote, packet *DefaultPacket) error {
	select {
	case p.eventChan <- &EventPacket{
		IHandler: handler,
		IRemote:  remote,
		Packet:   packet,
	}:
	default:
		xlog.PrintfErr("push EventPacket failed with eventChan full. remote:%v packet:%v", remote, packet)
		return errors.WithMessage(xerror.ChannelFull, xruntime.Location())
	}
	return nil
}

// EventDisconnect 事件-断开链接
type EventDisconnect struct {
	IHandler
	IRemote
}

// EventConnect 事件-链接成功
type EventConnect struct {
	IHandler
	IRemote
}

// EventPacket 事件-数据包
type EventPacket struct {
	IHandler
	IRemote
	Packet *DefaultPacket
}
