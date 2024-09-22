package tcp

import (
	"github.com/pkg/errors"
	xerror "xcore/lib/error"
	xlog "xcore/lib/log"
	xruntime "xcore/lib/runtime"
)

type IEvent interface {
	Connect(remote *DefaultRemote) error    // 链接 放入 事件中
	Disconnect(remote *DefaultRemote) error // 断开链接 放入 事件中
	Packet(packet *DefaultPacket) error     // 数据包 放入 事件中
}

type DefaultEvent struct {
	eventChan chan<- interface{} // 待处理的事件
}

// Connect 连接
func (p *DefaultEvent) Connect(remote *DefaultRemote) error {
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
func (p *DefaultEvent) Disconnect(remote *DefaultRemote) error {
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
func (p *DefaultEvent) Packet(packet *DefaultPacket) error {
	select {
	case p.eventChan <- &EventPacket{
		Packet: packet,
	}:
	default:
		xlog.PrintfErr("push DefaultPacket failed with eventChan full. remote:%v", packet.Remote)
		return errors.WithMessage(xerror.ChannelFull, xruntime.Location())
	}
	return nil
}

// EventDisconnect 事件-断开链接
type EventDisconnect struct {
	Remote *DefaultRemote
}

// EventConnect 事件-链接成功
type EventConnect struct {
	Remote *DefaultRemote
}

// EventPacket 事件-数据包
type EventPacket struct {
	Packet *DefaultPacket
}
