package packet

import (
	"encoding/binary"
	xutil "xcore/lib/util"
)

// 包头
type defaultHeader struct {
	PacketLength uint32 // 总包长度,包含包头＋包体长度
	MessageID    uint32 // 命令
	SessionID    uint32 // 会话id
	ResultID     uint32 // 结果id
	Key          uint64
}

// NewDefaultHeader 新建包头
func NewDefaultHeader() IHeader {
	return &defaultHeader{}
}

func (p *defaultHeader) Pack(data []byte) {
	xutil.PackUint32(p.PacketLength, data[0:])
	xutil.PackUint32(p.MessageID, data[4:])
	xutil.PackUint32(p.SessionID, data[8:])
	xutil.PackUint32(p.ResultID, data[12:])
	xutil.PackUint64(p.Key, data[16:])
}

func (p *defaultHeader) Unpack(data []byte) {
	p.PacketLength = binary.LittleEndian.Uint32(data[0:4])
	p.MessageID = binary.LittleEndian.Uint32(data[4:8])
	p.SessionID = binary.LittleEndian.Uint32(data[8:12])
	p.ResultID = binary.LittleEndian.Uint32(data[12:16])
	p.Key = binary.LittleEndian.Uint64(data[16:24])
}

//func (p *defaultHeader) GetPacketLength() uint32 {
//	return p.PacketLength
//}
//
//func (p *defaultHeader) GetMessageID() uint32 {
//	return p.MessageID
//}
//
//func (p *defaultHeader) SetPacketLength(value uint32) {
//	p.PacketLength = value
//}
//func (p *defaultHeader) SetMessageID(value uint32) {
//	p.MessageID = value
//}
//func (p *defaultHeader) SetKey(value uint64) {
//	p.Key = value
//}
//func (p *defaultHeader) SetSessionID(value uint32) {
//	p.SessionID = value
//}
//func (p *defaultHeader) SetResultID(value uint32) {
//	p.ResultID = value
//}
