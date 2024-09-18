package util

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
	xerror "xcore/lib/error"
	xruntime "xcore/lib/runtime"
)

// Split2Uint32 拆分字符串, 返回 uint32 类型的 slice
// e.g.: "1,2,3" => []uint32{1, 2, 3}
func Split2Uint32(s string, sep string) (u32Slice []uint32, err error) {
	if 0 == len(s) {
		return u32Slice, nil
	}

	slice := strings.Split(s, sep)
	var u64 uint64
	for _, v := range slice {
		if u64, err = strconv.ParseUint(v, 10, 32); err != nil {
			return u32Slice, errors.WithMessage(err, xruntime.Location())
		}
		u32Slice = append(u32Slice, uint32(u64))
	}
	return u32Slice, nil
}

// Split2Map 拆分字符串, 返回key为uint32类型、val为int64类型的map
// e.g.: "1,2;3,4" => map[uint32]int64{1:2, 3:4}
func Split2Map(s string, sep1 string, sep2 string) (map[uint32]int64, error) {
	slice := strings.Split(s, sep1)
	m := make(map[uint32]int64)
	var err error
	for _, v := range slice {
		if 0 == len(v) {
			continue
		}
		sliceAttr := strings.Split(v, sep2)
		if len(sliceAttr) != 2 {
			return nil, errors.WithMessage(xerror.Param, xruntime.Location())
		}
		var idUint64 uint64
		var valInt64 int64
		if idUint64, err = strconv.ParseUint(sliceAttr[0], 10, 32); err != nil {
			return nil, errors.WithMessage(err, xruntime.Location())
		}
		if valInt64, err = strconv.ParseInt(sliceAttr[1], 10, 32); err != nil {
			return nil, errors.WithMessage(err, xruntime.Location())
		}
		m[uint32(idUint64)] = valInt64
	}
	return m, nil
}
