package model

import "tevat.nd.org/basecode/goost/encoding/binary"

type Request struct {
	ID       uint32
	Method   binary.BytesWithUint16Len
	Metadata binary.BytesWithUint16Len
	Msg      binary.BytesWithUint32Len
}

type Response struct {
	ID       uint32
	Method   binary.BytesWithUint16Len
	Result   bool
	Metadata binary.BytesWithUint16Len
	Msg      binary.BytesWithUint32Len
}

type MsgError struct {
	ErrorCode int32
}

type MsgSyn struct {
	Version uint16
	Token   uint32
	LastID  uint32
}

type MsgSynAck struct {
	Version uint16
	Token   uint32
}

type MsgKeepalive struct {
	LastID uint32
}

type ResponseWithHeader struct {
	Size     uint16
	ID       uint32
	Response Response
}
