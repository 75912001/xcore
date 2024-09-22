package handler

import (
	"google.golang.org/protobuf/proto"
	protobufgateway "xcore/impl/protobuf/gateway"
	xnetmessage "xcore/lib/net/message"
)

var GMessage xnetmessage.Mgr

func init() {
	// todo menglc [优化] 通过配置文件配置,自动生成
	GMessage.Register(protobufgateway.UserOnlineMsg_CMD,
		xnetmessage.NewMessage().WithHandler(UserOnlineMsg).
			WithNewProtoMessage(func() proto.Message { return new(protobufgateway.UserOnlineMsg) }))
}
