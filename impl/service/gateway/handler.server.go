package gateway

import (
	"github.com/pkg/errors"
	xconstants "xcore/lib/constants"
	xerror "xcore/lib/error"
	xnetpacket "xcore/lib/net/packet"
	xnettcp "xcore/lib/net/tcp"
	xruntime "xcore/lib/runtime"
)

//type Server struct {
//	*xnettcp.DefaultHandlerServer
//}
//
//func NewServer() *Server {
//	def := &Server{
//		DefaultHandlerServer: xnettcp.NewDefaultHandlerServer(),
//	}
//	return def
//}

func (p *Service) OnConnect(remote xnettcp.IRemote) error {
	p.Log.Tracef("OnConnect: %v", remote)
	return nil
}

func (p *Service) OnCheckPacketLength(length uint32) error {
	if length < xnetpacket.DefaultHeaderSize || xconstants.PacketLengthDefault < length {
		return xerror.PacketHeaderLength
	}
	return nil
}

func (p *Service) OnCheckPacketLimit(remote xnettcp.IRemote) error {
	if false {
		return xerror.PacketQuantityLimit
	}
	return nil
}

func (p *Service) OnUnmarshalPacket(remote xnettcp.IRemote, data []byte) (xnetpacket.IPacket, error) {
	header := xnetpacket.NewDefaultHeader()
	header.Unpack(data)

	// todo menglc 判断消息是否禁用

	// todo menglc 判断消息是否需要转发
	if 0x10000 <= header.MessageID && header.MessageID <= 0x1ffff { // login
		// todo menglc login
		return nil, errors.WithMessage(xerror.NotImplemented, xruntime.Location())
	} else if 0x20000 <= header.MessageID && header.MessageID <= 0x2ffff { // gateway
		packet := xnetpacket.NewDefaultPacket().WithDefaultHeader(header)
		packet.IMessage = GMessage.Find(header.MessageID)
		if packet.IMessage == nil {
			return nil, errors.WithMessage(xerror.MessageIDNonExistent, xruntime.Location())
		}
		pb, err := packet.IMessage.Unmarshal(data[xnetpacket.DefaultHeaderSize:])
		if err != nil {
			return nil, errors.WithMessage(err, xruntime.Location())
		}
		packet.PBMessage = pb
		return packet, nil
	} else if 0x30000 <= header.MessageID && header.MessageID <= 0x3ffff { // logic
		return nil, errors.WithMessage(xerror.NotImplemented, xruntime.Location())
	} else {
		return nil, errors.WithMessage(xerror.NotImplemented, xruntime.Location())
	}
}

func (p *Service) OnPacket(remote xnettcp.IRemote, packet xnetpacket.IPacket) error {
	defaultPacket, ok := packet.(*xnetpacket.DefaultPacket)
	if !ok {
		return xerror.TypeMismatch
	}
	defaultPacket.IMessage.Set(remote, defaultPacket)
	return defaultPacket.IMessage.Execute()
}

func (p *Service) OnDisconnect(remote xnettcp.IRemote) error {
	p.Log.Tracef("OnDisconnect: %v", remote)
	return nil
}
