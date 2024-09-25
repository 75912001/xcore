package packet

type IHeader interface {
	GetLength() uint32 // 总长度
	GetCmd() uint32    // 消息ID
}
