package packet

import (
	"encoding/binary"
	xutil "xcore/lib/util"
)

// DefaultHeader 默认包头
type DefaultHeader struct {
	PacketLength uint32 // 总包长度,包含包头＋包体长度
	MessageID    uint32 // 命令
	SessionID    uint32 // 会话id
	ResultID     uint32 // 结果id
	Key          uint64
}

// NewDefaultHeader 新建包头
func NewDefaultHeader(packetLength uint32,
	messageID uint32,
	sessionID uint32,
	resultID uint32,
	key uint64) *DefaultHeader {
	return &DefaultHeader{
		PacketLength: packetLength,
		MessageID:    messageID,
		SessionID:    sessionID,
		ResultID:     resultID,
		Key:          key,
	}
}

func (p *DefaultHeader) Pack(data []byte) {
	xutil.PackUint32(p.PacketLength, data[0:])
	xutil.PackUint32(p.MessageID, data[4:])
	xutil.PackUint32(p.SessionID, data[8:])
	xutil.PackUint32(p.ResultID, data[12:])
	xutil.PackUint64(p.Key, data[16:])
}

func (p *DefaultHeader) Unpack(data []byte) {
	p.PacketLength = binary.LittleEndian.Uint32(data[0:4])
	p.MessageID = binary.LittleEndian.Uint32(data[4:8])
	p.SessionID = binary.LittleEndian.Uint32(data[8:12])
	p.ResultID = binary.LittleEndian.Uint32(data[12:16])
	p.Key = binary.LittleEndian.Uint64(data[16:24])
}
