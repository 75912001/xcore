package gateway

import (
	xnettcp "xcore/lib/net/tcp"
)

type Server struct {
	*xnettcp.DefaultHandlerServer
}

func NewServer() *Server {
	return &Server{
		DefaultHandlerServer: xnettcp.NewDefaultHandlerServer(),
	}
}

func (p *Server) OnConnect(remote *xnettcp.DefaultRemote) error {
	return nil
}
func (p *Server) OnCheckPacketLength(length uint32) error {
	return nil
}
func (p *Server) OnCheckPacketLimit(remote *xnettcp.DefaultRemote) error {
	return nil
}
func (p *Server) OnUnmarshalPacket(remote *xnettcp.DefaultRemote, data []byte) (*xnettcp.DefaultPacket, error) {
	return nil, nil
}
func (p *Server) OnPacket(packet *xnettcp.DefaultPacket) error {
	return nil
}
func (p *Server) OnDisconnect(remote *xnettcp.DefaultRemote) error {
	// todo menglc

	if err := p.DefaultHandlerServer.OnDisconnect(remote); err != nil {
		return err
	}
	return nil
}
