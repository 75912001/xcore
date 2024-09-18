package util

import (
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
