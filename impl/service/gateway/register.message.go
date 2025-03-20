package gateway

import (
	"google.golang.org/protobuf/proto"
	xprotobufgateway "xcore/impl/protobuf/gateway"
	xcallback "xcore/lib/control"
	"xcore/lib/message"
)

var GMessage message.Mgr

func init() {
	// todo menglc [优化] 通过配置文件配置,自动生成
	GMessage.Register(xprotobufgateway.UserOnlineMsgReq_CMD,
		message.NewOptions().
			WithHandler(xcallback.NewCallBack(UserOnlineMsg)).
			WithNewProtoMessage(func() proto.Message { return new(xprotobufgateway.UserOnlineMsgReq) }),
	)
	GMessage.Register(xprotobufgateway.UserHeartbeatMsgReq_CMD,
		message.NewOptions().
			WithHandler(xcallback.NewCallBack(UserHeartbeatMsg)).
			WithNewProtoMessage(func() proto.Message { return new(xprotobufgateway.UserHeartbeatMsgReq) }),
	)
}
