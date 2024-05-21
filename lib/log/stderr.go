package log

import (
	"log"
	"os"
	"runtime"
	libconstants "xcore/lib/constants"
)

var stdErr = log.New(os.Stderr, "", 0)

// PrintErr 输出到os.Stderr
func PrintErr(v ...interface{}) {
	if isEnable() { // 日志已启用,使用日志打印
		stdInstance.log(stdInstance.NewEntry(), LevelError, v...)
	} else {
		pc, _, line, ok := runtime.Caller(calldepth1)
		funcName := libconstants.Unknown
		if !ok {
			line = 0
		} else {
			funcName = runtime.FuncForPC(pc).Name()
		}
		formatAndPrint(LevelError, line, funcName, v...)
	}
}

// PrintfErr 输出到os.Stderr
func PrintfErr(format string, v ...interface{}) {
	if isEnable() { // 日志已启用,使用日志打印
		stdInstance.logf(stdInstance.NewEntry(), LevelError, format, v...)
	} else {
		pc, _, line, ok := runtime.Caller(calldepth1)
		funcName := libconstants.Unknown
		if !ok {
			line = 0
		} else {
			funcName = runtime.FuncForPC(pc).Name()
		}
		formatAndPrint(LevelError, line, funcName, v...)
	}
}
