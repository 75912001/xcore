package tcp

// IHandler 处理 接口
type IHandler interface {
	OnConnect(remote *DefaultRemote) error                                        // 处理-链接成功-对方链接过来
	OnCheckPacketLength(length uint32) error                                      // 处理-检查长度是否合法(包头中)
	OnCheckPacketLimit(remote *DefaultRemote) error                               // 处理-限流
	OnUnmarshalPacket(remote *DefaultRemote, data []byte) (*DefaultPacket, error) // 处理-数据包-反序列化 // data:数据 [NOTE] 如果保存该参数 则 需要copy
	OnPacket(packet *DefaultPacket) error                                         // 处理-数据包
	OnDisconnect(remote *DefaultRemote) error                                     // 处理-断开链接
}