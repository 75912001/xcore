package gateway

import (
	xconstants "xcore/lib/constants"
	xerror "xcore/lib/error"
	xnetpacket "xcore/lib/net/packet"
	xnettcp "xcore/lib/net/tcp"
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
	if length < 24 || xconstants.PacketLengthDefault < length {
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
	header := xnetpacket.NewDefaultHeader(0, 0, 0, 0, 0)
	header.Unpack(data)
	// todo menglc 判断消息是否禁用

	packet := xnetpacket.NewDefaultPacket(header, nil)
	packet.IMessage = GMessage.Find(header.MessageID)
	if packet.IMessage == nil {
		return nil, xerror.MessageIDNonExistent
	}
	pb, err := packet.IMessage.Unmarshal(data[24:])
	if err != nil {
		return nil, err
	}
	packet.PBMessage = pb
	return packet, nil
}
func (p *Service) OnPacket(packet xnetpacket.IPacket) error {
	defaultPacket, ok := packet.(*xnetpacket.DefaultPacket)
	if !ok {
		return xerror.TypeMismatch
	}

	defaultPacket.IMessage.Set(1, 2, 3, 4, 5)
	return defaultPacket.IMessage.Execute()
}
func (p *Service) OnDisconnect(remote xnettcp.IRemote) error {
	p.Log.Tracef("OnDisconnect: %v", remote)
	return nil
}
