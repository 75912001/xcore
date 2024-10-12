package packet

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	xerror "xcore/lib/error"
	xnetmessage "xcore/lib/net/message"
	xruntime "xcore/lib/runtime"
)

// DefaultPacket 默认数据包
type DefaultPacket struct {
	DefaultHeader   *DefaultHeader // 包头
	PBMessage       proto.Message  // 消息
	PassThroughData []byte         // 包体数据(不带包头)
	CTX             context.Context
	IMessage        xnetmessage.IMessage
}

// NewDefaultPacket 新建数据包
func NewDefaultPacket(header *DefaultHeader, pb proto.Message) *DefaultPacket {
	return &DefaultPacket{
		DefaultHeader: header,
		PBMessage:     pb,
	}
}

func (p *DefaultPacket) Marshal() (data []byte, err error) {
	if p.PBMessage == nil {
		return nil, xerror.NotImplemented
	}
	data, err = proto.Marshal(p.PBMessage)
	if err != nil {
		return nil, errors.WithMessage(err, xruntime.Location())
	}
	p.DefaultHeader.PacketLength = 24 + uint32(len(data))
	buf := make([]byte, p.DefaultHeader.PacketLength)
	p.DefaultHeader.Pack(buf)
	copy(buf[24:p.DefaultHeader.PacketLength], data)
	return buf, nil
}

func (p *DefaultPacket) Unmarshal(data []byte) (header IHeader, message proto.Message, err error) {
	// todo menglc 反序列化
	return nil, nil, xerror.NotImplemented
}

// IsPassThrough 是否透传
// 是: DefaultPacket.PassThroughBody 可用
// 否: DefaultPacket.options 可用
//func (p *DefaultPacket) IsPassThrough() bool {
//	return p.PassThroughData != nil
//}
