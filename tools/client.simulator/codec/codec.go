package codec

import "io"

// Encoder interface
type Encoder interface {
	Encode(io.Writer, any) error
}

// Decoder interface
type Decoder interface {
	Decode(io.Reader) (any, error)
}

type Coder interface {
	Encoder
	Decoder
}
