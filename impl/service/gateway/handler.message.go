package gateway

import (
	"fmt"
	xruntime "xcore/lib/runtime"
)

//func UserOnlineMsg(ctx context.Context, header packet.IHeader, message proto.Message, obj interface{}) error {
//	// todo menglc 处理用户上线
//	fmt.Println(ctx, header, message, obj, xruntime.Location())
//	return nil
//}

func UserOnlineMsg(args ...interface{}) error {
	// todo menglc 处理用户上线
	fmt.Println(args, xruntime.Location())
	return nil
}
