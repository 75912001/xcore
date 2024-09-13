package packet

import (
	"context"
	"google.golang.org/protobuf/proto"
	xerror "xcore/lib/error"
)

// IPacket 接口-数据包
type IPacket interface {
	// Marshal 序列化
	Marshal() (data []byte, err error)
	// Unmarshal 反序列化
	Unmarshal(data []byte) (header IHeader, message proto.Message, err error)
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type DefaultPacket struct {
	CTX             context.Context // ctx
	Header          IHeader         // 包头
	Message         proto.Message   // 数据
	PassThroughBody []byte          // 透传数据
}

func NewDefaultPacket() *DefaultPacket {
	return &DefaultPacket{}
}

func (p *DefaultPacket) Marshal() (data []byte, err error) {
	return nil, xerror.NotImplemented
}

func (p *DefaultPacket) Unmarshal(data []byte) (header IHeader, message proto.Message, err error) {
	return nil, nil, xerror.NotImplemented
}

func (p *DefaultPacket) WithCTX(ctx context.Context) *DefaultPacket {
	p.CTX = ctx
	return p
}

func (p *DefaultPacket) WithHeader(header IHeader) *DefaultPacket {
	p.Header = header
	return p
}

func (p *DefaultPacket) WithMessage(message proto.Message) *DefaultPacket {
	p.Message = message
	return p
}

func (p *DefaultPacket) WithPassThroughBody(passThroughBody []byte) *DefaultPacket {
	p.PassThroughBody = passThroughBody
	return p
}

// IsPassThrough 是否透传
// 是: DefaultPacket.PassThroughBody 可用
// 否: DefaultPacket.Message 可用
func (p *DefaultPacket) IsPassThrough() bool {
	return p.PassThroughBody != nil
}
