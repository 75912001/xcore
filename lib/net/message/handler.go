package message

import (
    "context"
    "google.golang.org/protobuf/proto"
    "xcore/lib/net/packet"
)

// Handler 处理函数
type Handler func(ctx context.Context, header packet.IHeader, message proto.Message, obj interface{}) error
