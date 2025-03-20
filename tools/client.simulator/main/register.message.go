package main

import (
	xprotobufgateway "github.com/75912001/xcore/impl/protobuf/gateway"
	xprotobuflogin "github.com/75912001/xcore/impl/protobuf/login"
	xcallback "github.com/75912001/xcore/lib/control"
	xmessage "github.com/75912001/xcore/lib/message"
	"google.golang.org/protobuf/proto"
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
