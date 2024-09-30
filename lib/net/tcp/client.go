package tcp

import (
	"context"
	"github.com/pkg/errors"
	"net"
	xerror "xcore/lib/error"
	xnetevent "xcore/lib/net/event"
	xnethandler "xcore/lib/net/handler"
	xnetpacket "xcore/lib/net/packet"
	xnetremote "xcore/lib/net/remote"
	xruntime "xcore/lib/runtime"
)

// Client 客户端
type Client struct {
	xnetevent.IEvent
	xnethandler.IHandler
	xnetremote.IRemote
	xnetpacket.IPacket
}

func NewClient(packet xnetpacket.IPacket, handler xnethandler.IHandler) *Client {
	return &Client{
		IEvent:   nil,
		IHandler: handler,
		IRemote:  nil,
		IPacket:  packet,
	}
}

// Connect 连接
//
//	每个连接有 一个 发送协程, 一个 接收协程
func (p *Client) Connect(ctx context.Context, opts ...*clientOptions) error {
	newOpts := mergeClientOptions(opts...)
	if err := clientConfigure(newOpts); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	p.IEvent = xnetevent.NewDefaultEvent(newOpts.eventChan)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", *newOpts.serverAddress)
	if nil != err {
		return errors.WithMessage(err, xruntime.Location())
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if nil != err {
		return errors.WithMessage(err, xruntime.Location())
	}
	p.IRemote = NewDefaultRemote(conn, make(chan interface{}, *newOpts.sendChanCapacity))
	p.IRemote.Start(&newOpts.connOptions, p.IEvent, p.IHandler)
	return nil
}

// Disconnect 主动断开连接
func (p *Client) Disconnect() error {
	if !p.IRemote.IsConnect() {
		return errors.WithMessage(xerror.Link, xruntime.Location())
	}
	p.IRemote.Stop()
	return nil
}
