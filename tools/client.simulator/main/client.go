package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	xerror "xcore/lib/error"
	xnetpacket "xcore/lib/net/packet"
	xnettcp "xcore/lib/net/tcp"
	xruntime "xcore/lib/runtime"
)

var client *defaultClient // 客户端

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
	header := xnetpacket.NewHeader()
	header.Unpack(data)

	// todo menglc 判断消息是否禁用

	packet := xnetpacket.NewPacket().WithHeader(header)
	packet.IMessage = GMessage.Find(header.MessageID)
	if packet.IMessage == nil {
		return nil, xerror.NotExist
	}
	pb, err := packet.IMessage.Unmarshal(data[xnetpacket.HeaderSize:])
	if err != nil {
		return nil, err
	}
	packet.PBMessage = pb
	return packet, nil
}

func (p *defaultClient) OnPacket(remote xnettcp.IRemote, packet xnetpacket.IPacket) error {
	defaultPacket, ok := packet.(*xnetpacket.Packet)
	if !ok {
		return xerror.Mismatch
	}
	var msgName string = reflect.TypeOf(defaultPacket.PBMessage).Elem().Name()
	var strHeader string
	var strPBMessage string
	{
		fmt.Println()
		fmt.Printf("\033[32mMessage Name: %s\033[0m\n", msgName)
		type HexDefaultHeader struct {
			PacketLength uint32
			MessageID    string
			SessionID    uint32
			ResultID     string
			Key          uint64
		}
		header := defaultPacket.Header
		hexHeader := &HexDefaultHeader{
			PacketLength: header.Length,
			MessageID:    fmt.Sprintf("0x%x", header.MessageID),
			SessionID:    header.SessionID,
			ResultID:     fmt.Sprintf("0x%x", header.ResultID),
			Key:          header.Key,
		}
		headerJson, err := json.MarshalIndent(hexHeader, "", "  ")
		if err != nil {
			fmt.Printf("\033[31mJSON marshaling failed: %s\033[0m", err)
		}
		fmt.Printf("\033[32mHeader: %s\033[0m\n", headerJson)
		strHeader = string(headerJson)
	}
	{
		pbMessageJson, err := json.MarshalIndent(defaultPacket.PBMessage, "", "  ")
		if err != nil {
			fmt.Printf("\033[31mJSON marshaling failed: %s\033[0m", err)
		}
		fmt.Printf("\033[32mMessage: %s\033[0m\n", pbMessageJson)
		strPBMessage = string(pbMessageJson)
	}
	log.Infof("\n======recv message======\n%s\nHeader: %s\nMessage: %s", msgName, strHeader, strPBMessage)
	return nil
}
func (p *defaultClient) OnDisconnect(remote xnettcp.IRemote) error {
	ColorPrintf(Red, "%v\n", xruntime.Location())
	panic(fmt.Errorf("OnDisconnect:%v", remote))
	return nil
}
