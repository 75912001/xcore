package tcp

// IHandler 处理 接口
type IHandler interface {
	OnConnect(remote IRemote) error                                        // 处理-链接成功-对方链接过来
	OnCheckPacketLength(length uint32) error                               // 处理-检查长度是否合法 [read 协程中调用] (包头中)
	OnCheckPacketLimit(remote IRemote) error                               // 处理-限流 [read 协程中调用]
	OnUnmarshalPacket(remote IRemote, data []byte) (*DefaultPacket, error) // 处理-数据包-反序列化 [read 协程中调用]// data:数据 [NOTE] 如果保存该参数 则 需要copy
	OnPacket(packet *DefaultPacket) error                                  // 处理-数据包
	OnDisconnect(remote IRemote) error                                     // 处理-断开链接
}
