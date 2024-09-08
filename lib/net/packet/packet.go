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
type defaultPacket struct {
	CTX             context.Context // ctx
	Header          IHeader         // 包头
	Message         proto.Message   // 数据
	PassThroughBody []byte          // 透传数据
}

func NewDefaultPacket() *defaultPacket {
	return &defaultPacket{}
}

func (p *defaultPacket) Marshal() (data []byte, err error) {
	return nil, xerror.NotImplemented
}

func (p *defaultPacket) Unmarshal(data []byte) (header IHeader, message proto.Message, err error) {
	return nil, nil, xerror.NotImplemented
}

func (p *defaultPacket) WithCTX(ctx context.Context) *defaultPacket {
	p.CTX = ctx
	return p
}

func (p *defaultPacket) WithHeader(header IHeader) *defaultPacket {
	p.Header = header
	return p
}

func (p *defaultPacket) WithMessage(message proto.Message) *defaultPacket {
	p.Message = message
	return p
}

func (p *defaultPacket) WithPassThroughBody(passThroughBody []byte) *defaultPacket {
	p.PassThroughBody = passThroughBody
	return p
}

// IsPassThrough 是否透传
// 是: defaultPacket.PassThroughBody 可用
// 否: defaultPacket.Message 可用
func (p *defaultPacket) IsPassThrough() bool {
	return p.PassThroughBody != nil
}
