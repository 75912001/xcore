package tcp

import (
	"context"
	"github.com/pkg/errors"
	"net"
	xruntime "xcore/lib/runtime"
)

// Client 客户端
type Client struct {
	IEvent   IEvent
	IHandler IHandler
	IRemote  IRemote
}

func NewClient(handler IHandler) *Client {
	return &Client{
		IEvent:   nil,
		IHandler: handler,
		IRemote:  nil,
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
	p.IEvent = NewEvent(newOpts.eventChan)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", *newOpts.serverAddress)
	if nil != err {
		return errors.WithMessage(err, xruntime.Location())
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if nil != err {
		return errors.WithMessage(err, xruntime.Location())
	}
	p.IRemote = NewRemote(conn, make(chan interface{}, *newOpts.sendChanCapacity))
	p.IRemote.Start(&newOpts.connOptions, p.IEvent, p.IHandler)
	return nil
}
