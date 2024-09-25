package gateway

import (
	"fmt"
	xruntime "xcore/lib/runtime"
)

func logCallBackFunc(level uint32, outString string) {
	if xruntime.IsDebug() {
		fmt.Println(level, outString)
	}
	return
}
