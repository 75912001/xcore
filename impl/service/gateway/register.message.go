package gateway

import (
	"google.golang.org/protobuf/proto"
	xprotobufgateway "xcore/impl/protobuf/gateway"
	xcallback "xcore/lib/callback"
	xnetmessage "xcore/lib/net/message"
)

var GMessage xnetmessage.Mgr

func init() {
	// todo menglc [优化] 通过配置文件配置,自动生成
	GMessage.Register(xprotobufgateway.UserOnlineMsgReq_CMD,
		xnetmessage.NewOptions().
			WithHandler(xcallback.NewDefaultCallBack(UserOnlineMsg)).
			WithNewProtoMessage(func() proto.Message { return new(xprotobufgateway.UserOnlineMsgReq) }),
	)
	GMessage.Register(xprotobufgateway.UserOnlineMsgRes_CMD,
		xnetmessage.NewOptions().
			WithHandler(xcallback.NewDefaultCallBack(nil)).
			WithNewProtoMessage(func() proto.Message { return new(xprotobufgateway.UserOnlineMsgRes) }),
	)
}
