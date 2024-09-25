package tcp

import (
	"github.com/pkg/errors"
	xerror "xcore/lib/error"
	xlog "xcore/lib/log"
	xruntime "xcore/lib/runtime"
)

type IEvent interface {
	Connect(remote IRemote) error                       // 链接 放入 事件中
	Disconnect(remote IRemote) error                    // 断开链接 放入 事件中
	Packet(remote IRemote, packet *DefaultPacket) error // 数据包 放入 事件中
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
func (p *defaultEvent) Connect(remote IRemote) error {
	select {
	case p.eventChan <- &EventConnect{
		Remote: remote,
	}:
	default:
		xlog.PrintfErr("push EventConnect failed with eventChan full. remote:%v", remote)
		return errors.WithMessage(xerror.ChannelFull, xruntime.Location())
	}
	return nil
}

// Disconnect 断开链接
func (p *defaultEvent) Disconnect(remote IRemote) error {
	select {
	case p.eventChan <- &EventDisconnect{
		Remote: remote,
	}:
	default:
		xlog.PrintfErr("push EventDisconnect failed with eventChan full. remote:%v", remote)
		return errors.WithMessage(xerror.ChannelFull, xruntime.Location())
	}
	return nil
}

// Packet 数据包
func (p *defaultEvent) Packet(remote IRemote, packet *DefaultPacket) error {
	select {
	case p.eventChan <- &EventPacket{
		Remote: remote,
		Packet: packet,
	}:
	default:
		xlog.PrintfErr("push EventPacket failed with eventChan full. remote:%v packet:%v", remote, packet)
		return errors.WithMessage(xerror.ChannelFull, xruntime.Location())
	}
	return nil
}

// EventDisconnect 事件-断开链接
type EventDisconnect struct {
	Remote IRemote
}

// EventConnect 事件-链接成功
type EventConnect struct {
	Remote IRemote
}

// EventPacket 事件-数据包
type EventPacket struct {
	Remote IRemote
	Packet *DefaultPacket
}
