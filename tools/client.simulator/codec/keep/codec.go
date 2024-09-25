package keep

import (
	"bytes"
	"io"

	"tevat.nd.org/basecode/goost/encoding/binary"
	"tevat.nd.org/toolchain/simulator/codec/model"
)

const (
	MsgTypeError     uint16 = 10
	MsgTypeSyn       uint16 = 11
	MsgTypeSynAck    uint16 = 12
	MsgTypeKeepalive uint16 = 15

	MsgTypeUser uint16 = 1000
)

const NetVersion = 0x12

type Codec struct{}

func NewCodec() *Codec {
	return &Codec{}
}

func (c *Codec) Encode(w io.Writer, data any) error {
	hLen := 4 + 2
	header := make([]byte, hLen)
	l := binary.Size(data) + hLen
	binary.LittleEndian.PutUint32(header, uint32(l))
	switch data.(type) {
	case model.MsgSyn:
		binary.LittleEndian.PutUint16(header[4:], MsgTypeSyn)
	case model.MsgSynAck:
		binary.LittleEndian.PutUint16(header[4:], MsgTypeSynAck)
	case model.MsgError:
		binary.LittleEndian.PutUint16(header[4:], MsgTypeError)
	case model.ResponseWithHeader:
		binary.LittleEndian.PutUint16(header[4:], MsgTypeUser)
	case *model.Request:
		binary.LittleEndian.PutUint16(header[4:], MsgTypeUser)
	}
	buf := bytes.NewBuffer(make([]byte, 0, l))
	if _, err := buf.Write(header); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, data); err != nil {
		return err
	}
	_, err := io.Copy(w, buf)
	return err
}

func (*Codec) Decode(r io.Reader) (any, error) {
	hLen := 4 + 2
	header := make([]byte, hLen)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, err
	}
	msgType := binary.LittleEndian.Uint16(header[4:])
	switch msgType {
	case MsgTypeSynAck:
		var syncAck model.MsgSynAck
		if err := binary.Read(r, binary.LittleEndian, &syncAck); err != nil {
			return nil, err
		}
		return syncAck, nil
	case MsgTypeSyn:
		var syn model.MsgSyn
		if err := binary.Read(r, binary.LittleEndian, &syn); err != nil {
			return nil, err
		}
		return syn, nil
	case MsgTypeKeepalive:
		var keepalive model.MsgKeepalive
		if err := binary.Read(r, binary.LittleEndian, &keepalive); err != nil {
			return nil, err
		}
		return keepalive, nil
	case MsgTypeUser:
		var res model.ResponseWithHeader
		if err := binary.Read(r, binary.LittleEndian, &res); err != nil {
			return nil, err
		}
		return res.Response, nil
	}
	return nil, nil
}
