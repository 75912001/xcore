package message

import (
	"google.golang.org/protobuf/proto"
	xcallback "xcore/lib/callback"
)

type IMessage interface {
	xcallback.ICallBack
	Marshal(message proto.Message) (data []byte, err error)
	Unmarshal(data []byte) (message proto.Message, err error)
	JsonUnmarshal(data []byte) (message proto.Message, err error)
}
