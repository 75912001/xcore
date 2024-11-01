package packet

// PacketTransparent 透传数据包
type PacketTransparent struct {
	Header  *Header // 包头
	RawData []byte  // 原始数据(包头+包体)
}

// NewPacketTransparent 新建-透传数据包
func NewPacketTransparent() *PacketTransparent {
	return &PacketTransparent{}
}

func (p *PacketTransparent) WithHeader(header *Header) *PacketTransparent {
	p.Header = header
	return p
}

func (p *PacketTransparent) Marshal() (data []byte, err error) {
	return p.RawData, nil
}

// IsPassThrough 是否透传
// 是: Packet.PassThroughBody 可用
// 否: Packet.options 可用
//func (p *Packet) IsPassThrough() bool {
//	return p.PassThroughData != nil
//}
