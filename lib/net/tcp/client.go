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
	options *ClientOptions
}

// Connect 连接
//
//	每个连接有 一个 发送协程, 一个 接收协程
func (p *Client) Connect(ctx context.Context, opts ...*ClientOptions) error {
	p.options = mergeClientOptions(opts...)
	if err := clientConfigure(p.options); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	p.Event = &DefaultEvent{
		eventChan: p.options.eventChan,
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp4", *p.options.serverAddress)
	if nil != err {
		return errors.WithMessage(err, xruntime.Location())
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if nil != err {
		return errors.WithMessage(err, xruntime.Location())
	}
	p.Remote.Conn = conn
	p.Remote.sendChan = make(chan interface{}, *p.options.sendChanCapacity)
	//p.Remote.DefaultPacket = p.options.packet
	p.Remote.start(&p.options.connOptions, p.Event, p.Handler)
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
