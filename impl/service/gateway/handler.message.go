package gateway

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"xcore/lib/net/packet"
	xruntime "xcore/lib/runtime"
)

func UserOnlineMsg(ctx context.Context, header packet.IHeader, message proto.Message, obj interface{}) error {
	// todo menglc 处理用户上线
	fmt.Println(ctx, header, message, obj, xruntime.Location())
	return nil
}
