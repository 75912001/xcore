package gateway

import (
	"fmt"
	xruntime "github.com/75912001/xcore/lib/runtime"
)

func logCallBackFunc(level uint32, outString string) {
	if xruntime.IsDebug() {
		fmt.Println(level, outString)
	}
	return
}
