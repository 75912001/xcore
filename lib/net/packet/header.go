package packet

type IHeader interface {
	GetLength() uint32
	GetCmd() uint32
}
