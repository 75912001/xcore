package packet

import (
	"encoding/binary"
	xutil "xcore/lib/util"
)

const DefaultHeaderSize uint32 = 24

// Header 包头
type Header struct {
	Length    uint32 // 总包长度,包含包头＋包体长度
	MessageID uint32 // 命令
	SessionID uint32 // 会话id
	ResultID  uint32 // 结果id
	Key       uint64
}

// NewDefaultHeader 新建包头
func NewDefaultHeader() *Header {
	return &Header{}
}

func (p *Header) WithLength(length uint32) *Header {
	p.Length = length
	return p
}

func (p *Header) WithMessageID(messageID uint32) *Header {
	p.MessageID = messageID
	return p
}

func (p *Header) WithSessionID(sessionID uint32) *Header {
	p.SessionID = sessionID
	return p
}

func (p *Header) WithResultID(resultID uint32) *Header {
	p.ResultID = resultID
	return p
}

func (p *Header) WithKey(key uint64) *Header {
	p.Key = key
	return p
}

func (p *Header) Pack(data []byte) {
	xutil.PackUint32(p.Length, data[0:])
	xutil.PackUint32(p.MessageID, data[4:])
	xutil.PackUint32(p.SessionID, data[8:])
	xutil.PackUint32(p.ResultID, data[12:])
	xutil.PackUint64(p.Key, data[16:])
}

func (p *Header) Unpack(data []byte) {
	p.Length = binary.LittleEndian.Uint32(data[0:4])
	p.MessageID = binary.LittleEndian.Uint32(data[4:8])
	p.SessionID = binary.LittleEndian.Uint32(data[8:12])
	p.ResultID = binary.LittleEndian.Uint32(data[12:16])
	p.Key = binary.LittleEndian.Uint64(data[16:DefaultHeaderSize])
}
