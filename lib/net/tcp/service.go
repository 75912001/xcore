package tcp

import (
	"context"
	xerror "github.com/75912001/xcore/lib/error"
	xlog "github.com/75912001/xcore/lib/log"
	xruntime "github.com/75912001/xcore/lib/runtime"
	xutil "github.com/75912001/xcore/lib/util"
	"github.com/pkg/errors"
	"net"
	"runtime/debug"
	"time"
)

// Service 服务端
type Service struct {
	IEvent   IEvent
	IHandler IHandler
	listener *net.TCPListener //监听
	options  *serviceOptions
}

// NewService 新建服务
func NewService(handler IHandler) *Service {
	return &Service{
		IEvent:   nil,
		IHandler: handler,
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
func (p *Service) Start(_ context.Context, opts ...*serviceOptions) error {
	p.options = mergeServiceOptions(opts...)
	if err := serviceConfigure(p.options); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	p.IEvent = NewEvent(p.options.eventChan)
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
					xlog.PrintErr(xerror.GoroutinePanic, err, debug.Stack())
				}
			}
			xlog.PrintInfo(xerror.GoroutineDone)
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
func (p *Service) Stop() {
	if p.listener != nil {
		err := p.listener.Close()
		if err != nil {
			xlog.PrintfErr("listener close err:%v", err)
		}
		p.listener = nil
	}
}

func (p *Service) handleConn(conn *net.TCPConn) {
	remote := NewRemote(conn, make(chan interface{}, *p.options.sendChanCapacity))
	if err := p.IEvent.Connect(p.IHandler, remote); err != nil {
		xlog.PrintfErr("event.Connect err:%v", err)
		return
	}
	remote.Start(&p.options.connOptions, p.IEvent, p.IHandler)
}
