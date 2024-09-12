package tcp

import (
	"context"
	"github.com/pkg/errors"
	"net"
	"runtime/debug"
	"time"
	xconstants "xcore/lib/constants"
	xerror "xcore/lib/error"
	xlog "xcore/lib/log"
	xruntime "xcore/lib/runtime"
	xutil "xcore/lib/util"
)

// Server 己方作为服务端
type Server struct {
	Handler  IHandler // 需要的话, 可以设置
	Event    IEvent
	options  *ServerOptions
	listener *net.TCPListener //监听
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
func (p *Server) Start(_ context.Context, opts ...*ServerOptions) error {
	p.options = mergeServerOptions(opts...)
	if err := serverConfigure(p.options); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	p.Event = &DefaultEvent{
		eventChan: p.options.eventChan,
	}
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
func (p *Server) Stop() {
	if p.listener != nil {
		err := p.listener.Close()
		if err != nil {
			xlog.PrintfErr("listener close err:%v", err)
		}
		p.listener = nil
	}
}

func (p *Server) OnConnect(remote *DefaultRemote) error {
	return p.options.OnConnect(remote)
}

func (p *Server) OnCheckPacketLength(length uint32) error {
	if p.options.OnCheckPacketLength == nil {
		return nil
	}
	return p.options.OnCheckPacketLength(length)
}

func (p *Server) OnCheckPacketLimit(remote *DefaultRemote) error {
	if p.options.OnCheckPacketLimit == nil {
		return nil
	}
	return p.options.OnCheckPacketLimit(remote)
}

func (p *Server) OnUnmarshalPacket(remote *DefaultRemote, data []byte) (*Packet, error) {
	return p.options.OnUnmarshalPacket(remote, data)
}

func (p *Server) OnPacket(parsePacket *Packet) error {
	return p.options.OnPacket(parsePacket)
}

func (p *Server) OnDisconnect(remote *DefaultRemote) error {
	if err := p.options.OnDisconnect(remote); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if remote.IsConn() {
		remote.stop()
	}
	return nil
}

// ActiveDisconnect 逻辑层 主动 断开连接
func (p *Server) ActiveDisconnect(remote *DefaultRemote) error {
	if remote == nil || !remote.IsConn() {
		return errors.WithMessage(xerror.Link, xruntime.Location())
	}
	remote.ActiveDisconnection = true
	if err := p.OnDisconnect(remote); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}

func (p *Server) handleConn(conn *net.TCPConn) {
	remote := &DefaultRemote{
		Conn:     conn,
		sendChan: make(chan interface{}, *p.options.sendChanCapacity),
		Packet:   p.options.packet,
	}
	if err := p.Event.Connect(remote); err != nil {
		xlog.PrintfErr("Event.Connect err:%v", err)
		return
	}
	remote.start(&p.options.connOptions, p.Event)
}
