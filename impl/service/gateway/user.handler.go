package gateway

import (
	"fmt"
	"github.com/pkg/errors"
	xcommonservice "xcore/impl/common"
	xutil "xcore/lib/control"
	xerror "xcore/lib/error"
	xnettcp "xcore/lib/net/tcp"
	packet2 "xcore/lib/packet"
	xruntime "xcore/lib/runtime"
)

//	type Server struct {
//		*xnettcp.DefaultHandlerServer
//	}
//
//	func NewService() *Server {
//		def := &Server{
//			DefaultHandlerServer: xnettcp.NewDefaultHandlerServer(),
//		}
//		return def
//	}
func userLoginTimeout(arg ...interface{}) error {
	fmt.Printf("cbSecond:%v\n", arg...)
	return nil
}
func (p *Service) OnConnect(remote xnettcp.IRemote) error {
	p.Log.Tracef("OnConnect: %v", remote)
	u := newUser(remote)
	remote.(*xnettcp.Remote).Object = u
	gUserMgr.add(u, u.connect.IRemote)
	// 用户登录超时
	p.Timer.AddSecond(
		xutil.NewCallBack(
			func(args ...interface{}) error {
				u := args[0].(*User)
				if !u.timeoutValid {
					return nil
				}
				if u.login { // 已经登录
					return nil
				}
				u.exit()
				return nil
			},
			u,
		),
		p.TimeMgr.ShadowTimestamp()+UserLoginTimeOut,
	)
	// 用户心跳超时
	p.Timer.AddSecond(
		xutil.NewCallBack(
			func(args ...interface{}) error {
				u := args[0].(*User)
				if !u.timeoutValid {
					return nil
				}
				if u.login { // 已经登录
					return nil
				}
				u.exit()
				return nil
			},
			u,
		),
		p.TimeMgr.ShadowTimestamp()+UserHeartbeatInterval,
	)
	return nil
}

func (p *Service) OnCheckPacketLength(length uint32) error {
	if length < packet2.HeaderSize || *p.BenchMgr.Json.Base.PacketLengthMax < length {
		return xerror.Length
	}
	return nil
}

func (p *Service) OnCheckPacketLimit(remote xnettcp.IRemote) error {
	if false {
		return xerror.Quantity
	}
	return nil
}

func (p *Service) OnUnmarshalPacket(remote xnettcp.IRemote, data []byte) (packet2.IPacket, error) {
	header := packet2.NewHeader()
	header.Unpack(data)
	// todo menglc 判断消息是否禁用

	switch xcommonservice.GetServiceTypeByMessageID(header.MessageID) {
	case xcommonservice.GatewayMessage:
		packet := packet2.NewPacket().WithHeader(header)
		packet.IMessage = GMessage.Find(header.MessageID)
		if packet.IMessage == nil {
			return nil, errors.WithMessage(xerror.NotExist, xruntime.Location())
		}
		pb, err := packet.IMessage.Unmarshal(data[packet2.HeaderSize:])
		if err != nil {
			return nil, errors.WithMessage(err, xruntime.Location())
		}
		packet.PBMessage = pb
		return packet, nil
	case xcommonservice.LoginMessage, xcommonservice.LogicMessage:
		packet := packet2.NewPacketPassThrough().WithHeader(header)
		packet.RawData = make([]byte, len(data))
		copy(packet.RawData, data)
		return packet, nil
	default:
		return nil, errors.WithMessage(xerror.NotImplemented, xruntime.Location())
	}
}

func (p *Service) OnPacket(remote xnettcp.IRemote, packet packet2.IPacket) error {
	defaultRemote := remote.(*xnettcp.Remote)
	user := defaultRemote.Object.(*User)
	switch packet.(type) {
	case *packet2.Packet:
		defaultPacket, ok := packet.(*packet2.Packet)
		if !ok {
			return xerror.Mismatch
		}
		defaultPacket.IMessage.Override(user, defaultPacket)
		return defaultPacket.IMessage.Execute()
	case *packet2.PacketPassThrough:
		packetTransparent, ok := packet.(*packet2.PacketPassThrough)
		if !ok {
			return xerror.Mismatch
		}
		switch xcommonservice.GetServiceTypeByMessageID(packetTransparent.Header.MessageID) {
		case xcommonservice.LoginMessage:
			if user.LoginService != nil {
				return errors.WithMessage(xerror.Duplicate, xruntime.Location())
			}
			loginService := p.LoginServiceMgr.GetLoginService()
			if loginService == nil {
				return errors.WithMessage(xerror.Unavailable, xruntime.Location())
			}
			user.LoginService = loginService
			// 将消息转发到 login service
			err := user.LoginService.IRemote.Send(packetTransparent)
			if err != nil {
				return errors.WithMessage(err, xruntime.Location())
			}
			return errors.WithMessage(xerror.NotImplemented, xruntime.Location())
		case xcommonservice.LogicMessage:
			// todo menglc 处理 logic message
			return errors.WithMessage(xerror.NotImplemented, xruntime.Location())
		default:
			return errors.WithMessage(xerror.NotImplemented, xruntime.Location())
		}
	default:
		return errors.WithMessage(xerror.NotImplemented, xruntime.Location())
	}
}

func (p *Service) OnDisconnect(remote xnettcp.IRemote) error {
	p.Log.Tracef("OnDisconnect: %v", remote)
	gUserMgr.remove(remote)
	return nil
}
