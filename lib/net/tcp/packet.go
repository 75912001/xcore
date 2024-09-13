package tcp

import (
	"context"
	"google.golang.org/protobuf/proto"
	xerror "xcore/lib/error"
	xnetmessage "xcore/lib/net/message"
	"xcore/lib/net/packet"
)

// Packet 数据包
type Packet struct {
	Remote          *DefaultRemote
	Header          packet.IHeader // 包头
	Message         proto.Message  // 解析出的数据
	PassThroughData []byte         // 包体数据(不带包头)
	Entity          *xnetmessage.Message
	CTX             context.Context
}

func (p *Packet) Marshal() (data []byte, err error) {
	return nil, xerror.NotImplemented
}

func (p *Packet) Unmarshal(data []byte) (header packet.IHeader, message proto.Message, err error) {
	return nil, nil, xerror.NotImplemented
}

// EventDisconnect 事件-断开链接
type EventDisconnect struct {
	Remote *DefaultRemote
}

// EventConnect 事件-链接成功
type EventConnect struct {
	Remote *DefaultRemote
}
