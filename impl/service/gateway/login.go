package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	xcommonservice "xcore/impl/common"
	xerror "xcore/lib/error"
	xnettcp "xcore/lib/net/tcp"
	packet2 "xcore/lib/packet"
	xruntime "xcore/lib/runtime"
)

type LoginService struct {
	*xnettcp.Client
}

func NewLoginService() *LoginService {
	return &LoginService{}
}

func (p *LoginService) OnConnect(remote xnettcp.IRemote) error {
	return nil
}

func (p *LoginService) OnCheckPacketLength(length uint32) error {
	return nil
}

func (p *LoginService) OnCheckPacketLimit(remote xnettcp.IRemote) error {
	return nil
}

func (p *LoginService) OnUnmarshalPacket(remote xnettcp.IRemote, data []byte) (packet2.IPacket, error) {
	header := packet2.NewHeader()
	header.Unpack(data)
	// todo menglc 判断消息是否禁用
	switch xcommonservice.GetServiceTypeByMessageID(header.MessageID) {
	case xcommonservice.GatewayMessage:
		packet := packet2.NewPacket().WithHeader(header)
		packet.IMessage = GMessage.Find(header.MessageID)
		if packet.IMessage == nil {
			return nil, errors.WithMessage(xerror.NotExist, xruntime.Location())
		}
		pb, err := packet.IMessage.Unmarshal(data[packet2.HeaderSize:])
		if err != nil {
			return nil, errors.WithMessage(err, xruntime.Location())
		}
		packet.PBMessage = pb
		return packet, nil
	case xcommonservice.LoginMessage:
		packet := packet2.NewPacketPassThrough().WithHeader(header)
		packet.RawData = make([]byte, len(data))
		copy(packet.RawData, data)
		return packet, nil
	default:
		return nil, errors.WithMessage(xerror.NotImplemented, xruntime.Location())
	}
}

func (p *LoginService) OnPacket(remote xnettcp.IRemote, packet packet2.IPacket) error {
	defaultPacket, ok := packet.(*packet2.Packet)
	if !ok {
		return xerror.Mismatch
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
func (p *LoginService) OnDisconnect(remote xnettcp.IRemote) error {
	// todo menglc

	return nil
}
