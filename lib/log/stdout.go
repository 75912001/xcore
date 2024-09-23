package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
	xconstants "xcore/lib/constants"
)

var stdOut = log.New(os.Stdout, "", 0)

// PrintInfo 输出到os.Stdout
func PrintInfo(l ILog, v ...interface{}) {
	if l != nil { // 日志已启用,使用日志打印
		l.Info(v...)
	} else {
		funcName := xconstants.Unknown
		pc, _, line, ok := runtime.Caller(calldepth1)
		if ok {
			funcName = runtime.FuncForPC(pc).Name()
		}
		element := newEntry().
			withLevel(LevelInfo).
			withTime(time.Now()).
			withCallerInfo(fmt.Sprintf(callerInfoFormat, line, funcName)).
			withMessage(fmt.Sprint(v...))
		formatLogData(element)
		_ = stdOut.Output(calldepth2, element.outString)
	}
}

// PrintfInfo 输出到os.Stdout
func PrintfInfo(l ILog, format string, v ...interface{}) {
	if l != nil { // 日志已启用,使用日志打印
		l.Infof(format, v...)
	} else {
		funcName := xconstants.Unknown
		pc, _, line, ok := runtime.Caller(calldepth1)
		if ok {
			funcName = runtime.FuncForPC(pc).Name()
		}
		element := newEntry().
			withLevel(LevelInfo).
			withTime(time.Now()).
			withCallerInfo(fmt.Sprintf(callerInfoFormat, line, funcName)).
			withMessage(fmt.Sprintf(format, v...))
		formatLogData(element)
		_ = stdOut.Output(calldepth2, element.outString)
	}
}
