package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	xconstants "xcore/lib/constants"
)

var stdErr = log.New(os.Stderr, "", 0)

// PrintErr 输出到os.Stderr
func PrintErr(v ...interface{}) {
	if isEnable() { // 日志已启用,使用日志打印
		mgrInstance.log(mgrInstance.newEntry(), LevelError, v...)
	} else {
		funcName := xconstants.Unknown
		pc, _, line, ok := runtime.Caller(calldepth1)
		if ok {
			funcName = runtime.FuncForPC(pc).Name()
		}
		formatAndPrint(stdErr, LevelError, line, funcName, v...)
	}
}

// PrintfErr 输出到os.Stderr
func PrintfErr(format string, v ...interface{}) {
	if isEnable() { // 日志已启用,使用日志打印
		mgrInstance.logf(mgrInstance.newEntry(), LevelError, format, v...)
	} else {
		funcName := xconstants.Unknown
		pc, _, line, ok := runtime.Caller(calldepth1)
		if ok {
			funcName = runtime.FuncForPC(pc).Name()
		}
		formatAndPrint(stdErr, LevelError, line, funcName, fmt.Sprintf(format, v...))
	}
}
