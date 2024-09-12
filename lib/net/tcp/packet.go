package tcp

import (
	"context"
	"google.golang.org/protobuf/proto"
	xnetmessage "xcore/lib/net/message"
	"xcore/lib/net/packet"
)

// Packet 数据包
type Packet struct {
	Remote          *Remote
	Header          packet.IHeader // 包头
	Message         proto.Message  // 解析出的数据
	PassThroughData []byte         // 包体数据(不带包头)
	Entity          *xnetmessage.Message
	CTX             context.Context
}

// OnCheckPacketLength 检查长度是否合法(包头中)
type OnCheckPacketLength func(length uint32) error

// OnCheckPacketLimit 限流
//
//	返回:
//		nil:不限流
type OnCheckPacketLimit func(remote *Remote) error

// OnUnmarshalPacket [multithreading] 反序列化数据包
// data:数据 [NOTE] 如果保存该参数 则 需要copy
// 返回值: 包头, 消息, 包体数据(不带包头)
// [NOTE] 多协程调用
type OnUnmarshalPacket func(remote *Remote, data []byte) (*Packet, error)

// OnPacket 处理数据包
type OnPacket func(parsePacket *Packet) error

// EventDisconnect 事件-断开链接
type EventDisconnect struct {
	Remote *Remote
}

// OnDisconnect 处理断开链接
type OnDisconnect func(remote *Remote) error

// EventConnect 事件-链接成功
type EventConnect struct {
	Remote *Remote
}

// OnConnect 处理-链接成功
type OnConnect func(remote *Remote) error
