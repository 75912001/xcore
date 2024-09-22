package tcp

import (
	xerror "xcore/lib/error"
)

type DefaultHandlerServer struct {
}

func NewDefaultHandlerServer() *DefaultHandlerServer {
	return &DefaultHandlerServer{}
}

func (p *DefaultHandlerServer) OnConnect(remote *DefaultRemote) error {
	return xerror.NotImplemented
}
func (p *DefaultHandlerServer) OnCheckPacketLength(length uint32) error {
	return xerror.NotImplemented
}
func (p *DefaultHandlerServer) OnCheckPacketLimit(remote *DefaultRemote) error {
	return xerror.NotImplemented
}
func (p *DefaultHandlerServer) OnUnmarshalPacket(remote *DefaultRemote, data []byte) (*DefaultPacket, error) {
	return nil, xerror.NotImplemented
}
func (p *DefaultHandlerServer) OnPacket(packet *DefaultPacket) error {
	return xerror.NotImplemented
}

// OnDisconnect [NOTE] 需要外部实现调用
func (p *DefaultHandlerServer) OnDisconnect(remote *DefaultRemote) error {
	if remote.IsConnect() {
		remote.Stop()
	}
	return nil
}
