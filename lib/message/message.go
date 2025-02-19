package message

import (
	"google.golang.org/protobuf/proto"
	xcontrol "xcore/lib/control"
)

type IMessage interface {
	xcontrol.ICallBack
	Marshal(message proto.Message) (data []byte, err error)
	Unmarshal(data []byte) (message proto.Message, err error)
	JsonUnmarshal(data []byte) (message proto.Message, err error)
}
