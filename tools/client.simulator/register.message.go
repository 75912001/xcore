package main

import (
	"google.golang.org/protobuf/proto"
	xprotobufgateway "xcore/impl/protobuf/gateway"
	xprotobuflogin "xcore/impl/protobuf/login"
	xcallback "xcore/lib/callback"
	xnetmessage "xcore/lib/net/message"
)

var GMessage xnetmessage.Mgr

func init() {
	// todo menglc [优化] 通过配置文件配置,自动生成
	GMessage.Register(xprotobufgateway.UserOnlineMsgReq_CMD,
		xnetmessage.NewOptions().
			WithHandler(xcallback.NewCallBack(nil)).
			WithNewProtoMessage(func() proto.Message { return new(xprotobufgateway.UserOnlineMsgReq) }),
	)
	GMessage.Register(xprotobufgateway.UserOnlineMsgRes_CMD,
		xnetmessage.NewOptions().
			WithHandler(xcallback.NewCallBack(nil)).
			WithNewProtoMessage(func() proto.Message { return new(xprotobufgateway.UserOnlineMsgRes) }),
	)

	GMessage.Register(xprotobuflogin.LoginMsgReq_CMD,
		xnetmessage.NewOptions().
			WithHandler(xcallback.NewCallBack(nil)).
			WithNewProtoMessage(func() proto.Message { return new(xprotobuflogin.LoginMsgReq) }),
	)

	GMessage.Register(xprotobuflogin.LoginMsgRes_CMD,
		xnetmessage.NewOptions().
			WithHandler(xcallback.NewCallBack(nil)).
			WithNewProtoMessage(func() proto.Message { return new(xprotobuflogin.LoginMsgRes) }),
	)
}
