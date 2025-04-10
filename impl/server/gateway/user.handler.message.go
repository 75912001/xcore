package gateway

import (
	"fmt"
	xprotobufgateway "github.com/75912001/xcore/impl/protobuf/gateway"
	xnettcp "github.com/75912001/xcore/lib/net/tcp"
	xnetpacket "github.com/75912001/xcore/lib/packet"
	xruntime "github.com/75912001/xcore/lib/runtime"
	"github.com/pkg/errors"
	"math/rand"
)

//func UserOnlineMsg(ctx context.Context, header packet.IHeader, message proto.Message, obj interface{}) error {
//	// todo menglc 处理用户上线
//	fmt.Println(ctx, header, message, obj, xruntime.Location())
//	return nil
//}

func UserOnlineMsg(args ...interface{}) error {
	user := args[0].(*User)
	defaultPacket := args[1].(*xnetpacket.Packet)
	pb := defaultPacket.PBMessage.(*xprotobufgateway.UserOnlineMsgReq)
	fmt.Println(user, defaultPacket, pb, xruntime.Location())
	// todo menglc 处理用户上线
	// 在用户数据中寻找,是否有该用户,处于非活跃状态,则激活
	// 在用户数据中寻找,是否有该用户,处于活跃状态,则断开连接

	// 返回消息
	res := &xprotobufgateway.UserOnlineMsgRes{
		Uid: 668599,
	}
	if err := xnettcp.Send(user.connect.IRemote, res, xprotobufgateway.UserOnlineMsgRes_CMD, 0, 0); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}

func UserHeartbeatMsg(args ...interface{}) error {
	user := args[0].(*User)
	defaultPacket := args[1].(*xnetpacket.Packet)
	pb := defaultPacket.PBMessage.(*xprotobufgateway.UserHeartbeatMsgReq)
	fmt.Println(user, defaultPacket, pb, xruntime.Location())

	if user.connect.heartbeatRandom == 0 { // 第一次心跳,不验证
		user.connect.heartbeatRandom = rand.Uint64()
	} else { // 验证本次收到的,是否是上次发送给用户的
		if user.connect.heartbeatRandom != pb.Random {
			// todo menglc 发送 错误码, 并断开连接
			//return xxx
		} else {
			user.connect.heartbeatRandom = rand.Uint64()
		}
	}
	// 返回消息
	res := &xprotobufgateway.UserHeartbeatMsgRes{
		Random: user.connect.heartbeatRandom,
	}
	if err := xnettcp.Send(user.connect.IRemote, res, xprotobufgateway.UserHeartbeatMsgRes_CMD, 0, 0); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}
