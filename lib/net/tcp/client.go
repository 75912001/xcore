package tcp

import (
	"context"
	"github.com/pkg/errors"
	"net"
	xerror "xcore/lib/error"
	xruntime "xcore/lib/runtime"
)

// Client 客户端
type Client struct {
	Handler IHandler
	Event   IEvent
	Remote  DefaultRemote
}

// Connect 连接
//
//	每个连接有 一个 发送协程, 一个 接收协程
func (p *Client) Connect(ctx context.Context, opts ...*ClientOptions) error {
	newOpts := mergeClientOptions(opts...)
	if err := clientConfigure(newOpts); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	p.Handler = newOpts.handler
	p.Event = &DefaultEvent{
		eventChan: newOpts.eventChan,
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp4", *newOpts.serverAddress)
	if nil != err {
		return errors.WithMessage(err, xruntime.Location())
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if nil != err {
		return errors.WithMessage(err, xruntime.Location())
	}
	p.Remote.Conn = conn
	p.Remote.sendChan = make(chan interface{}, *newOpts.sendChanCapacity)
	p.Remote.Packet = newOpts.packet
	p.Remote.start(&newOpts.connOptions, p.Event, p.Handler)
	return nil
}

// ActiveDisconnect 主动断开连接
func (p *Client) ActiveDisconnect() error {
	if !p.Remote.IsConnect() {
		return errors.WithMessage(xerror.Link, xruntime.Location())
	}
	p.Remote.ActiveDisconnection = true
	if err := p.Handler.OnDisconnect(&p.Remote); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}
