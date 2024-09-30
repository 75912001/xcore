package packet

import (
	"context"
	"google.golang.org/protobuf/proto"
	xerror "xcore/lib/error"
	xnetmessage "xcore/lib/net/message"
)

// defaultPacket 数据包
type defaultPacket struct {
	IHeader                // 包头
	proto.Message          // 解析出的数据
	PassThroughData []byte // 包体数据(不带包头)
	xnetmessage.IMessage
	CTX context.Context
}

// NewDefaultPacket 新建数据包
func NewDefaultPacket() IPacket {
	return &defaultPacket{}
}

func (p *defaultPacket) Marshal() (data []byte, err error) {
	// todo menglc 序列化
	return nil, xerror.NotImplemented
}

func (p *defaultPacket) Unmarshal(data []byte) (header IHeader, message proto.Message, err error) {
	// todo menglc 反序列化
	return nil, nil, xerror.NotImplemented
}

// IsPassThrough 是否透传
// 是: defaultPacket.PassThroughBody 可用
// 否: defaultPacket.options 可用
//func (p *defaultPacket) IsPassThrough() bool {
//	return p.PassThroughData != nil
//}
