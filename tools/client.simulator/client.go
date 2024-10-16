package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	xerror "xcore/lib/error"
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
	header := xnetpacket.NewDefaultHeader()
	header.Unpack(data)

	// todo menglc 判断消息是否禁用

	packet := xnetpacket.NewDefaultPacket().WithDefaultHeader(header)
	packet.IMessage = GMessage.Find(header.MessageID)
	if packet.IMessage == nil {
		return nil, xerror.MessageIDNonExistent
	}
	pb, err := packet.IMessage.Unmarshal(data[xnetpacket.DefaultHeaderSize:])
	if err != nil {
		return nil, err
	}
	packet.PBMessage = pb
	return packet, nil
}

func (p *defaultClient) OnPacket(remote xnettcp.IRemote, packet xnetpacket.IPacket) error {
	defaultPacket, ok := packet.(*xnetpacket.DefaultPacket)
	if !ok {
		return xerror.TypeMismatch
	}
	{
		fmt.Println()
		fmt.Printf("\033[32mMessage Name: %s\033[0m\n", reflect.TypeOf(defaultPacket.PBMessage).Elem().Name())
		type HexDefaultHeader struct {
			PacketLength uint32
			MessageID    string
			SessionID    uint32
			ResultID     string
			Key          uint64
		}
		header := defaultPacket.DefaultHeader
		hexHeader := &HexDefaultHeader{
			PacketLength: header.PacketLength,
			MessageID:    fmt.Sprintf("0x%x", header.MessageID),
			SessionID:    header.SessionID,
			ResultID:     fmt.Sprintf("0x%x", header.ResultID),
			Key:          header.Key,
		}
		headerJson, err := json.MarshalIndent(hexHeader, "", "  ")
		if err != nil {
			fmt.Printf("\033[31mJSON marshaling failed: %s\033[0m", err)
		}
		//fmt.Printf("\nHeader: %s\n", headerJson)
		fmt.Printf("\033[32mHeader: %s\033[0m\n", headerJson)
	}
	{
		pbMessageJson, err := json.MarshalIndent(defaultPacket.PBMessage, "", "  ")
		if err != nil {
			fmt.Printf("\033[31mJSON marshaling failed: %s\033[0m", err)
		}
		fmt.Printf("\033[32mMessage: %s\033[0m\n", pbMessageJson)
	}
	return nil
}
func (p *defaultClient) OnDisconnect(remote xnettcp.IRemote) error {
	// todo menglc

	return nil
}
