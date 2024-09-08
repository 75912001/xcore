package message

import (
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"reflect"
	"runtime"
	"strings"
	xconstants "xcore/lib/constants"
	xruntime "xcore/lib/runtime"
)

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

// 生成函数名称
func getFuncName(i interface{}, seps ...rune) string {
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
