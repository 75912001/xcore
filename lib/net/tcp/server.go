package tcp

import (
	"context"
	"github.com/pkg/errors"
	"net"
	"runtime/debug"
	"time"
	xconstants "xcore/lib/constants"
	xlog "xcore/lib/log"
	xruntime "xcore/lib/runtime"
	xutil "xcore/lib/util"
)

// 己方作为服务端
type server struct {
	IEvent   IEvent
	IHandler IHandler
	//xnetpacket.IPacket
	listener *net.TCPListener //监听
	options  *serverOptions
}

// NewServer 新建服务
func NewServer(handler IHandler) *server {
	return &server{
		IEvent:   nil,
		IHandler: handler,
		//IPacket:  packet,
		listener: nil,
		options:  nil,
	}
}

// 网络 错误 暂时
func netErrorTemporary(tempDelay time.Duration) (newTempDelay time.Duration) {
	if tempDelay == 0 {
		tempDelay = 5 * time.Millisecond
	} else {
		tempDelay *= 2
	}
	if max := 1 * time.Second; tempDelay > max {
		tempDelay = max
	}
	return tempDelay
}

// Start 运行服务
func (p *server) Start(_ context.Context, opts ...*serverOptions) error {
	p.options = mergeServerOptions(opts...)
	if err := serverConfigure(p.options); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	p.IEvent = NewDefaultEvent(p.options.eventChan)
	tcpAddr, err := net.ResolveTCPAddr("tcp", *p.options.listenAddress)
	if nil != err {
		return errors.WithMessage(err, xruntime.Location())
	}
	p.listener, err = net.ListenTCP("tcp", tcpAddr)
	if nil != err {
		return errors.WithMessage(err, xruntime.Location())
	}
	go func() {
		defer func() {
			if xruntime.IsRelease() {
				if err := recover(); err != nil {
					xlog.PrintErr(xconstants.GoroutinePanic, err, debug.Stack())
				}
			}
			xlog.PrintInfo(xconstants.GoroutineDone)
		}()
		var tempDelay time.Duration
		for {
			conn, err := p.listener.AcceptTCP()
			if nil != err {
				if xutil.IsNetErrorTemporary(err) {
					tempDelay = netErrorTemporary(tempDelay)
					xlog.PrintfErr("listen.AcceptTCP, IsNetErrorTemporary, tempDelay:%v, err:%v", tempDelay, err)
					time.Sleep(tempDelay)
					continue
				}
				xlog.PrintfErr("listen.AcceptTCP, err:%v", err)
				return
			}
			tempDelay = 0
			go p.handleConn(conn)
		}
	}()
	return nil
}

// Stop 停止 AcceptTCP
func (p *server) Stop() {
	if p.listener != nil {
		err := p.listener.Close()
		if err != nil {
			xlog.PrintfErr("listener close err:%v", err)
		}
		p.listener = nil
	}
}

// Disconnect 逻辑层 主动 断开连接
//func (p *server) Disconnect(remote IRemote) error {
//	if remote == nil || !remote.IsConnect() {
//		return errors.WithMessage(xerror.Link, xruntime.Location())
//	}
//	if remote.IsConnect() {
//		remote.Stop()
//	}
//	return nil
//}

func (p *server) handleConn(conn *net.TCPConn) {
	remote := NewDefaultRemote(conn, make(chan interface{}, *p.options.sendChanCapacity))
	if err := p.IEvent.Connect(p.IHandler, remote); err != nil {
		xlog.PrintfErr("event.Connect err:%v", err)
		return
	}
	remote.Start(&p.options.connOptions, p.IEvent, p.IHandler)
}
