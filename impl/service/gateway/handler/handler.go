package handler

import (
    "context"
    "fmt"
    "google.golang.org/protobuf/proto"
    protobufgateway "xcore/impl/protobuf/gateway"
    xnetmessage "xcore/lib/net/message"
    "xcore/lib/net/packet"
    xruntime "xcore/lib/runtime"
    xutil "xcore/lib/util"
)

var GMessage xnetmessage.Mgr

func init() {
    GMessage.Register(protobufgateway.UserOnlineMsg_CMD,
        xnetmessage.NewMessage().WithHandler(UserOnlineMsg).
            WithNewProtoMessage(func() proto.Message { return new(protobufgateway.UserOnlineMsg) }).
            WithStateSwitch(xutil.NewDefaultSwitch(true)))
}

func UserOnlineMsg(ctx context.Context, header packet.IHeader, message proto.Message, obj interface{}) error {
    // todo menglc 处理用户上线
    fmt.Println(ctx, header, message, obj, xruntime.Location())
    return nil
}
