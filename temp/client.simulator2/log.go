package main

import (
	"fmt"
	xlog "xcore/lib/log"
	xruntime "xcore/lib/runtime"
)

var glog xlog.ILog

func logCallBackFunc(level uint32, outString string) {
	if xruntime.IsDebug() {
		fmt.Println(level, outString)
	}
	return
}
