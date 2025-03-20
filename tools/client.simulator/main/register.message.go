package main

import (
	"google.golang.org/protobuf/proto"
	xprotobufgateway "xcore/impl/protobuf/gateway"
	xprotobuflogin "xcore/impl/protobuf/login"
	xcallback "xcore/lib/control"
	xmessage "xcore/lib/message"
)

var GMessage xmessage.Mgr

func init() {
	// todo menglc [优化] 通过配置文件配置,自动生成
	GMessage.Register(xprotobufgateway.UserOnlineMsgReq_CMD,
		xmessage.NewOptions().
			WithHandler(xcallback.NewCallBack(nil)).
			WithNewProtoMessage(func() proto.Message { return new(xprotobufgateway.UserOnlineMsgReq) }),
	)
	GMessage.Register(xprotobufgateway.UserOnlineMsgRes_CMD,
		xmessage.NewOptions().
			WithHandler(xcallback.NewCallBack(nil)).
			WithNewProtoMessage(func() proto.Message { return new(xprotobufgateway.UserOnlineMsgRes) }),
	)

	GMessage.Register(xprotobuflogin.LoginMsgReq_CMD,
		xmessage.NewOptions().
			WithHandler(xcallback.NewCallBack(nil)).
			WithNewProtoMessage(func() proto.Message { return new(xprotobuflogin.LoginMsgReq) }),
	)

	GMessage.Register(xprotobuflogin.LoginMsgRes_CMD,
		xmessage.NewOptions().
			WithHandler(xcallback.NewCallBack(nil)).
			WithNewProtoMessage(func() proto.Message { return new(xprotobuflogin.LoginMsgRes) }),
	)
}
