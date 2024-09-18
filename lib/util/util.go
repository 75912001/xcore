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
func If(condition bool, trueVal interface{}, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// IsDuplicateUint32 是否有重复uint32
func IsDuplicateUint32(uint32Slice []uint32) bool {
	set := make(map[uint32]struct{})
	for _, v := range uint32Slice {
		if _, ok := set[v]; ok {
			return true
		}
		set[v] = struct{}{}
	}
	return false
}
