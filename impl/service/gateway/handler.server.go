package gateway

import (
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
	return nil
}
func (p *Service) OnCheckPacketLimit(remote xnettcp.IRemote) error {
	return nil
}
func (p *Service) OnUnmarshalPacket(remote xnettcp.IRemote, data []byte) (xnetpacket.IPacket, error) {
	return nil, nil
}
func (p *Service) OnPacket(packet xnetpacket.IPacket) error {
	return nil
}
func (p *Service) OnDisconnect(remote xnettcp.IRemote) error {
	p.Log.Tracef("OnDisconnect: %v", remote)
	return nil
}
