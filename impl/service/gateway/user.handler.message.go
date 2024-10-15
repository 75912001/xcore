package gateway

import (
	"fmt"
	"github.com/pkg/errors"
	xprotobufgateway "xcore/impl/protobuf/gateway"
	xnetpacket "xcore/lib/net/packet"
	xnettcp "xcore/lib/net/tcp"
	xruntime "xcore/lib/runtime"
)

//func UserOnlineMsg(ctx context.Context, header packet.IHeader, message proto.Message, obj interface{}) error {
//	// todo menglc 处理用户上线
//	fmt.Println(ctx, header, message, obj, xruntime.Location())
//	return nil
//}

func UserOnlineMsg(args ...interface{}) error {
	remote := args[0].(xnettcp.IRemote)
	defaultPacket := args[1].(*xnetpacket.DefaultPacket)
	pb := defaultPacket.PBMessage.(*xprotobufgateway.UserOnlineMsgReq)
	fmt.Println(remote, defaultPacket, pb, xruntime.Location())

	// todo menglc 处理用户上线

	// 返回消息
	res := &xprotobufgateway.UserOnlineMsgRes{
		Uid: 668599,
	}
	if err := xnettcp.Send(remote, res, xprotobufgateway.UserOnlineMsgRes_CMD, 0, 0); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	remote.Stop()
	return nil
}
