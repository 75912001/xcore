package standard

import (
	"bytes"
	"io"
	"sync"

	"tevat.nd.org/basecode/goost/encoding/binary"
	"tevat.nd.org/toolchain/simulator/codec/model"
)

const (
	headerLen         = 4
	defaultBufferSize = 1024
)

type Codec struct {
	bufferPool sync.Pool
}

func NewCodec() *Codec {
	return &Codec{
		bufferPool: sync.Pool{
			New: func() any {
				return bytes.NewBuffer(make([]byte, 0, defaultBufferSize))
			},
		},
	}
}

func (c *Codec) Encode(w io.Writer, data any) error {
	buf := c.bufferPool.Get().(*bytes.Buffer)
	defer c.bufferPool.Put(buf)
	buf.Reset()

	l := binary.Size(data) + headerLen

	if err := binary.Write(buf, binary.LittleEndian, uint32(l)); err != nil {
		return err
	}

	if err := binary.Write(buf, binary.LittleEndian, data); err != nil {
		return err
	}

	_, err := io.Copy(w, buf)

	return err
}

func (*Codec) Decode(r io.Reader) (any, error) {
	var h uint32

	if err := binary.Read(r, binary.LittleEndian, &h); err != nil {
		return nil, err
	}

	var res model.Response
	err := binary.Read(r, binary.LittleEndian, &res)

	return res, err
}
