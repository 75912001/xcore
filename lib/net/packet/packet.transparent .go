package packet

import (
	"google.golang.org/protobuf/proto"
	xerror "xcore/lib/error"
)

// PacketTransparent 透传数据包
type PacketTransparent struct {
	DefaultHeader *DefaultHeader // 包头
	RawData       []byte         // 原始数据(包头+包体)
}

// NewPacketTransparent 新建-透传数据包
func NewPacketTransparent() *PacketTransparent {
	return &PacketTransparent{}
}

func (p *PacketTransparent) WithDefaultHeader(header *DefaultHeader) *PacketTransparent {
	p.DefaultHeader = header
	return p
}

func (p *PacketTransparent) Marshal() (data []byte, err error) {
	return p.RawData, nil
}

func (p *PacketTransparent) Unmarshal(data []byte) (header IHeader, message proto.Message, err error) {
	// todo menglc 反序列化
	return nil, nil, xerror.NotImplemented
}

// IsPassThrough 是否透传
// 是: DefaultPacket.PassThroughBody 可用
// 否: DefaultPacket.options 可用
//func (p *DefaultPacket) IsPassThrough() bool {
//	return p.PassThroughData != nil
//}
