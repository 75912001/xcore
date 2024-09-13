package packet

import (
	"google.golang.org/protobuf/proto"
)

// IPacket 接口-数据包
type IPacket interface {
	// Marshal 序列化
	Marshal() (data []byte, err error)
	// Unmarshal 反序列化
	Unmarshal(data []byte) (header IHeader, message proto.Message, err error)
}
