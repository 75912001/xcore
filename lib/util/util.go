package util

import (
	"net"
	"unsafe"
)

func IsLittleEndian() bool {
	var i int32 = 0x01020304
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)
	b := *pb
	return b == 0x04
}

// IsNetErrorTemporary checks if a network error is temporary.
// [NOTE] 不建议使用
func IsNetErrorTemporary(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && netErr.Temporary()
}
