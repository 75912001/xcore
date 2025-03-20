package main

import (
	"encoding/json"
	"fmt"
	xerror "github.com/75912001/xcore/lib/error"
	xnettcp "github.com/75912001/xcore/lib/net/tcp"
	xpacket "github.com/75912001/xcore/lib/packet"
	xruntime "github.com/75912001/xcore/lib/runtime"
	"reflect"
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

func (p *defaultClient) OnUnmarshalPacket(remote xnettcp.IRemote, data []byte) (xpacket.IPacket, error) {
	header := xpacket.NewHeader()
	header.Unpack(data)

	// todo menglc 判断消息是否禁用

	packet := xpacket.NewPacket().WithHeader(header)
	packet.IMessage = GMessage.Find(header.MessageID)
	if packet.IMessage == nil {
		return nil, xerror.NotExist
	}
	pb, err := packet.IMessage.Unmarshal(data[xpacket.HeaderSize:])
	if err != nil {
		return nil, err
	}
	packet.PBMessage = pb
	return packet, nil
}

func (p *defaultClient) OnPacket(remote xnettcp.IRemote, packet xpacket.IPacket) error {
	defaultPacket, ok := packet.(*xpacket.Packet)
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
