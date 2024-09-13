package tcp

import (
	xerror "xcore/lib/error"
)

type DefaultHandlerClient struct {
}

func (p *DefaultHandlerClient) OnConnect(_ *DefaultRemote) error {
	return nil
}
func (p *DefaultHandlerClient) OnCheckPacketLength(_ uint32) error {
	return nil
}
func (p *DefaultHandlerClient) OnCheckPacketLimit(_ *DefaultRemote) error {
	return nil
}
func (p *DefaultHandlerClient) OnUnmarshalPacket(remote *DefaultRemote, data []byte) (*Packet, error) {
	return nil, xerror.NotImplemented
}
func (p *DefaultHandlerClient) OnPacket(packet *Packet) error {
	return xerror.NotImplemented
}
func (p *DefaultHandlerClient) OnDisconnect(remote *DefaultRemote) error {
	return xerror.NotImplemented
	// 下面是断开连接后需要做的事情
	//if remote.IsConnect() {
	//	remote.stop()
	//}
	//return nil
}
