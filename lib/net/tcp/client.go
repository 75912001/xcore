package tcp

import (
	"context"
	"github.com/pkg/errors"
	"net"
	xerror "xcore/lib/error"
	xnetpacket "xcore/lib/net/packet"
	xruntime "xcore/lib/runtime"
)

// Client 客户端
type Client struct {
	IEvent
	IHandler
	IRemote
	xnetpacket.IPacket
}

// Connect 连接
//
//	每个连接有 一个 发送协程, 一个 接收协程
func (p *Client) Connect(ctx context.Context, packet xnetpacket.IPacket, opts ...*clientOptions) error {
	newOpts := mergeClientOptions(opts...)
	if err := clientConfigure(newOpts); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	p.IEvent = newDefaultEvent(newOpts.eventChan)
	p.IHandler = NewDefaultHandlerClient()
	tcpAddr, err := net.ResolveTCPAddr("tcp4", *newOpts.serverAddress)
	if nil != err {
		return errors.WithMessage(err, xruntime.Location())
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if nil != err {
		return errors.WithMessage(err, xruntime.Location())
	}
	defaultRemote := NewDefaultRemote(conn, make(chan interface{}, *newOpts.sendChanCapacity))
	defaultRemote.start(&newOpts.connOptions, p.IEvent, p.IHandler)
	p.IRemote = defaultRemote
	p.IPacket = packet
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
