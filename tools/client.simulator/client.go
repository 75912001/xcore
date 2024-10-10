package main

import (
	xnetpacket "xcore/lib/net/packet"
	xnettcp "xcore/lib/net/tcp"
)

type defaultClient struct {
	*xnettcp.Client
}

func (p *defaultClient) OnConnect(remote xnettcp.IRemote) error {
	return nil
}
func (p *defaultClient) OnCheckPacketLength(length uint32) error {
	return nil
}
func (p *defaultClient) OnCheckPacketLimit(remote xnettcp.IRemote) error {
	return nil
}
func (p *defaultClient) OnUnmarshalPacket(remote xnettcp.IRemote, data []byte) (xnetpacket.IPacket, error) {
	return nil, nil
}
func (p *defaultClient) OnPacket(packet xnetpacket.IPacket) error {
	return nil
}
func (p *defaultClient) OnDisconnect(remote xnettcp.IRemote) error {
	// todo menglc

	return nil
}
