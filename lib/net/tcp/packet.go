package tcp

import (
	"context"
	"google.golang.org/protobuf/proto"
	xerror "xcore/lib/error"
	xnetmessage "xcore/lib/net/message"
	"xcore/lib/net/packet"
)

// DefaultPacket 数据包
type DefaultPacket struct {
	Remote          *DefaultRemote
	Header          packet.IHeader // 包头
	Message         proto.Message  // 解析出的数据
	PassThroughData []byte         // 包体数据(不带包头)
	Entity          *xnetmessage.Message
	CTX             context.Context
}

func (p *DefaultPacket) Marshal() (data []byte, err error) {
	return nil, xerror.NotImplemented
}

func (p *DefaultPacket) Unmarshal(data []byte) (header packet.IHeader, message proto.Message, err error) {
	return nil, nil, xerror.NotImplemented
}

// IsPassThrough 是否透传
// 是: DefaultPacket.PassThroughBody 可用
// 否: DefaultPacket.Message 可用
func (p *DefaultPacket) IsPassThrough() bool {
	return p.PassThroughData != nil
}
