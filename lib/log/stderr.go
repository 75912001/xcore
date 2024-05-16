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
	if IsEnable() { // 日志已启用,使用日志打印
		GetInstance().log(LevelError, v...)
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
	if IsEnable() { // 日志已启用,使用日志打印
		GetInstance().logf(LevelError, format, v...)
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
