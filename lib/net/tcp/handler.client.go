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
	// todo menglc 实现解包
	return nil, xerror.NotImplemented
}
func (p *DefaultHandlerClient) OnPacket(packet *DefaultPacket) error {
	// todo menglc 实现处理数据包
	return xerror.NotImplemented
}
func (p *DefaultHandlerClient) OnDisconnect(remote *DefaultRemote) error {
	// todo menglc 实现断开连接
	return xerror.NotImplemented
	// 下面是断开连接后需要做的事情
	//if remote.IsConnect() {
	//	remote.Stop()
	//}
	//return nil
}
