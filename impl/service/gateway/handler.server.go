package gateway

import (
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
	return nil
}
func (p *Service) OnCheckPacketLength(length uint32) error {
	return nil
}
func (p *Service) OnCheckPacketLimit(remote xnettcp.IRemote) error {
	return nil
}
func (p *Service) OnUnmarshalPacket(remote xnettcp.IRemote, data []byte) (*xnettcp.DefaultPacket, error) {
	return nil, nil
}
func (p *Service) OnPacket(packet *xnettcp.DefaultPacket) error {
	return nil
}
func (p *Service) OnDisconnect(remote xnettcp.IRemote) error {
	// todo menglc

	if err := p.DefaultService.OnDisconnect(remote); err != nil {
		return err
	}
	return nil
}
