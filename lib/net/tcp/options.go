package tcp

type ConnOptions struct {
	readBuffer  *int // readBuffer sets the size of the operating system's receive buffer associated with the connection. [default]: 系统默认
	writeBuffer *int // writeBuffer sets the size of the operating system's transmit buffer associated with the connection. [default]: 系统默认
}
