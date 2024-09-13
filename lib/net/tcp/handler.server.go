package tcp

import (
	xerror "xcore/lib/error"
)

type DefaultHandlerServer struct {
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
func (p *DefaultHandlerServer) OnDisconnect(remote *DefaultRemote) error {
	return xerror.NotImplemented
	// 下面是断开连接后需要做的事情
	//if remote.IsConnect() {
	//	remote.stop()
	//}
	//return nil
}
