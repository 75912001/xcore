package message

import (
	"google.golang.org/protobuf/proto"
	xcallback "xcore/lib/callback"
)

type IMessage interface {
	xcallback.ICallBack
	Unmarshal(data []byte) (message proto.Message, err error)
}
