package tcp

import (
	xconstants "github.com/75912001/xcore/lib/constants"
	xerror "github.com/75912001/xcore/lib/error"
	xlog "github.com/75912001/xcore/lib/log"
	xcommon "github.com/75912001/xcore/lib/net/common"
	xpacket "github.com/75912001/xcore/lib/packet"
	xruntime "github.com/75912001/xcore/lib/runtime"
	"github.com/pkg/errors"
	"time"
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
func (p *Event) Connect(handler xcommon.IHandler, remote xcommon.IRemote) error {
	select {
	case p.eventChan <- &xcommon.Connect{
		IHandler: handler,
		IRemote:  remote,
	}:
	case <-time.After(xconstants.BusAddTimeoutDuration):
		xlog.PrintfErr("push Connect failed with eventChan full. remote:%v", remote)
		return errors.WithMessage(xerror.ChannelFull, xruntime.Location())
	default:
		xlog.PrintfErr("push Connect failed with eventChan full. remote:%v", remote)
		return errors.WithMessage(xerror.ChannelFull, xruntime.Location())
	}
	return nil
}

// Disconnect 断开链接
func (p *Event) Disconnect(handler xcommon.IHandler, remote xcommon.IRemote) error {
	select {
	case p.eventChan <- &xcommon.Disconnect{
		IHandler: handler,
		IRemote:  remote,
	}:
	case <-time.After(xconstants.BusAddTimeoutDuration):
		xlog.PrintfErr("push Disconnect failed with eventChan full. remote:%v", remote)
		return errors.WithMessage(xerror.ChannelFull, xruntime.Location())
	default:
		xlog.PrintfErr("push Disconnect failed with eventChan full. remote:%v", remote)
		return errors.WithMessage(xerror.ChannelFull, xruntime.Location())
	}
	return nil
}

// Packet 数据包
func (p *Event) Packet(handler xcommon.IHandler, remote xcommon.IRemote, packet xpacket.IPacket) error {
	select {
	case p.eventChan <- &xcommon.Packet{
		IHandler: handler,
		IRemote:  remote,
		IPacket:  packet,
	}:
	case <-time.After(xconstants.BusAddTimeoutDuration):
		xlog.PrintfErr("push Packet failed with eventChan full. remote:%v packet:%v", remote, packet)
		return errors.WithMessage(xerror.ChannelFull, xruntime.Location())
	default:
		xlog.PrintfErr("push Packet failed with eventChan full. remote:%v packet:%v", remote, packet)
		return errors.WithMessage(xerror.ChannelFull, xruntime.Location())
	}
	return nil
}
