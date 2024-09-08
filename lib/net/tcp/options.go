package tcp

type ConnOptions struct {
	readBuffer  *int // readBuffer sets the size of the operating system's receive buffer associated with the connection.
	writeBuffer *int // writeBuffer sets the size of the operating system's transmit buffer associated with the connection.
}
