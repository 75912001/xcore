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
func (p *DefaultHandlerClient) OnUnmarshalPacket(remote *DefaultRemote, data []byte) (*DefaultPacket, error) {
	return nil, xerror.NotImplemented
}
func (p *DefaultHandlerClient) OnPacket(packet *DefaultPacket) error {
	return xerror.NotImplemented
}

// OnDisconnect [NOTE] 需要外部实现调用
func (p *DefaultHandlerClient) OnDisconnect(remote *DefaultRemote) error {
	if remote.IsConnect() {
		remote.Stop()
	}
	return nil
}
