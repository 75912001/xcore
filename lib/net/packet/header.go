package packet

type IHeader interface {
	//GetPacketLength() uint32 // 总长度
	//GetMessageID() uint32    // 消息ID
	//SetPacketLength(uint32)
	//SetMessageID(uint32)
	//SetKey(uint64)
	//SetSessionID(uint32)
	//SetResultID(uint32)

	Pack([]byte)   // 将 成员变量 -> data 中
	Unpack([]byte) // 将 data 数据 -> 成员变量中
}
