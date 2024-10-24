package util

import (
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"net"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unsafe"
	xconstants "xcore/lib/constants"
	xruntime "xcore/lib/runtime"
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

// IsNetErrClosing checks if a network error is due to a closed connection.
func IsNetErrClosing(err error) bool {
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

// GetFuncName 获取函数名称
func GetFuncName(i interface{}, seps ...rune) string {
	if i == nil {
		return xconstants.Nil
	}
	funcName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	fields := strings.FieldsFunc(funcName, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})
	if size := len(fields); size > 0 {
		return fields[size-1]
	}
	return xconstants.Unknown
}

// MutableCopy 深拷贝
func MutableCopy(src proto.Message, dst proto.Message) error {
	data, err := proto.Marshal(src)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	err = proto.Unmarshal(data, dst)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}

func PackUint32(v uint32, buf []byte) {
	buf[0] = byte(v)
	buf[1] = byte(v >> 8)
	buf[2] = byte(v >> 16)
	buf[3] = byte(v >> 24)
}

func PackUint64(v uint64, buf []byte) {
	buf[0] = byte(v)
	buf[1] = byte(v >> 8)
	buf[2] = byte(v >> 16)
	buf[3] = byte(v >> 24)
	buf[4] = byte(v >> 32)
	buf[5] = byte(v >> 40)
	buf[6] = byte(v >> 48)
	buf[7] = byte(v >> 56)
}

func HexStringToUint32(hexStr string) (uint32, error) {
	// Remove the "0x" prefix if it exists
	if len(hexStr) > 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}
	// Parse the hex string to a uint32
	value, err := strconv.ParseUint(hexStr, 16, 32)
	if err != nil {
		return 0, err
	}
	return uint32(value), nil
}
