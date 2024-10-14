package gateway

import (
	"github.com/pkg/errors"
	xcommonservice "xcore/impl/common/service"
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
	packet := xnetpacket.NewDefaultPacket().WithDefaultHeader(header)
	switch xcommonservice.GetServiceTypeByMessageID(header.MessageID) {
	case xcommonservice.LoginMessage, xcommonservice.LogicMessage:
		packet.RawData = make([]byte, len(data))
		copy(packet.RawData, data)
		return packet, nil
	case xcommonservice.GatewayMessage:
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
	default:
		return nil, errors.WithMessage(xerror.NotImplemented, xruntime.Location())
	}
}

func (p *Service) OnPacket(remote xnettcp.IRemote, packet xnetpacket.IPacket) error {
	defaultPacket, ok := packet.(*xnetpacket.DefaultPacket)
	if !ok {
		return xerror.TypeMismatch
	}
	switch xcommonservice.GetServiceTypeByMessageID(defaultPacket.DefaultHeader.MessageID) {
	case xcommonservice.LoginMessage:
		// todo menglc 找到对应的 login service

		// 将消息转发到 login service
		return errors.WithMessage(xerror.NotImplemented, xruntime.Location())
	case xcommonservice.LogicMessage:
		// todo menglc 处理 logic message
		return errors.WithMessage(xerror.NotImplemented, xruntime.Location())
	case xcommonservice.GatewayMessage:
		defaultPacket.IMessage.Set(remote, defaultPacket)
		return defaultPacket.IMessage.Execute()
	default:
		return errors.WithMessage(xerror.NotImplemented, xruntime.Location())
	}
}

func (p *Service) OnDisconnect(remote xnettcp.IRemote) error {
	p.Log.Tracef("OnDisconnect: %v", remote)
	return nil
}
