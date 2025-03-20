package gateway

import (
	xprotobufgateway "github.com/75912001/xcore/impl/protobuf/gateway"
	xcallback "github.com/75912001/xcore/lib/control"
	"github.com/75912001/xcore/lib/message"
	"google.golang.org/protobuf/proto"
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
