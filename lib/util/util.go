package util

import (
	"net"
	"strings"
	"unsafe"
)

// IsLittleEndian 是否小端
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

// IsNetErrorTimeout checks if a network error is a timeout.
func IsNetErrorTimeout(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}

// IsErrNetClosing checks if a network error is due to a closed connection.
func IsErrNetClosing(err error) bool {
	return err != nil && strings.Contains(err.Error(), "use of closed network connection")
}

// If 三目运算符
// [NOTE] 传递的实参,会在调用时计算所有参数
// e.g.: If(true, 1, 2) => 1
func If(condition bool, trueVal interface{}, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// IsDuplicate 是否有重复
// e.g.: [1, 2, 3, 1] => true
func IsDuplicate(slice []interface{}, equals func(a, b interface{}) bool) bool {
	set := make(map[interface{}]struct{})
	for _, v := range slice {
		for k := range set {
			if equals(k, v) {
				return true
			}
		}
		set[v] = struct{}{}
	}
	return false
}
