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
	Event  IEvent
	Remote IRemote
	Packet xnetpacket.IPacket
}

// Connect 连接
//
//	每个连接有 一个 发送协程, 一个 接收协程
func (p *Client) Connect(ctx context.Context, opts ...*clientOptions) error {
	newOpts := mergeClientOptions(opts...)
	if err := clientConfigure(newOpts); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	p.Event = newDefaultEvent(newOpts.eventChan)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", *newOpts.serverAddress)
	if nil != err {
		return errors.WithMessage(err, xruntime.Location())
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if nil != err {
		return errors.WithMessage(err, xruntime.Location())
	}
	defaultRemote := NewDefaultRemote(conn, make(chan interface{}, *newOpts.sendChanCapacity), newOpts.handler)
	defaultRemote.start(&newOpts.connOptions, p.Event)
	p.Remote = defaultRemote
	p.Packet = newOpts.packet
	return nil
}

// ActiveDisconnect 主动断开连接
func (p *Client) ActiveDisconnect() error {
	if !p.Remote.IsConnect() {
		return errors.WithMessage(xerror.Link, xruntime.Location())
	}
	p.Remote.SetActiveDisconnection(true)
	if err := p.Remote.OnDisconnect(p.Remote); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}
