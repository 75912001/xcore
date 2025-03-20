package main

import (
	"fmt"
	xlog "github.com/75912001/xcore/lib/log"
	xruntime "github.com/75912001/xcore/lib/runtime"
)

var log xlog.ILog

func logCallBackFunc(level uint32, outString string) {
	if xruntime.IsDebug() {
		fmt.Println(level, outString)
	}
	return
}
