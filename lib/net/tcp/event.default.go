package tcp

import (
	xconstants "github.com/75912001/xcore/lib/constants"
	xlog "github.com/75912001/xcore/lib/log"
	xcommon "github.com/75912001/xcore/lib/net/common"
	xpacket "github.com/75912001/xcore/lib/packet"
	xruntime "github.com/75912001/xcore/lib/runtime"
	xutil "github.com/75912001/xcore/lib/util"
	"github.com/pkg/errors"
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
	err := xutil.PushEventWithTimeout(p.eventChan,
		&xcommon.Connect{
			IHandler: handler,
			IRemote:  remote,
		},
		xconstants.BusAddTimeoutDuration)
	if err != nil {
		xlog.PrintfErr("push Connect failed with eventChan full. remote:%v", remote)
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}

// Disconnect 断开链接
func (p *Event) Disconnect(handler xcommon.IHandler, remote xcommon.IRemote) error {
	err := xutil.PushEventWithTimeout(p.eventChan,
		&xcommon.Disconnect{
			IHandler: handler,
			IRemote:  remote,
		},
		xconstants.BusAddTimeoutDuration)
	if err != nil {
		xlog.PrintfErr("push Disconnect failed with eventChan full. remote:%v", remote)
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}

// Packet 数据包
func (p *Event) Packet(handler xcommon.IHandler, remote xcommon.IRemote, packet xpacket.IPacket) error {
	err := xutil.PushEventWithTimeout(p.eventChan,
		&xcommon.Packet{
			IHandler: handler,
			IRemote:  remote,
			IPacket:  packet,
		},
		xconstants.BusAddTimeoutDuration)
	if err != nil {
		xlog.PrintfErr("push Packet failed with eventChan full. remote:%v packet:%v", remote, packet)
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}
