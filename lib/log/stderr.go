package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
	xconstants "xcore/lib/constants"
)

var stdErr = log.New(os.Stderr, "", 0)

// PrintErr 输出到os.Stderr
func PrintErr(l ILog, v ...interface{}) {
	if l != nil { // 日志已启用,使用日志打印
		l.Error(v...)
	} else {
		funcName := xconstants.Unknown
		pc, _, line, ok := runtime.Caller(calldepth1)
		if ok {
			funcName = runtime.FuncForPC(pc).Name()
		}
		element := newEntry().
			withLevel(LevelError).
			withTime(time.Now()).
			withCallerInfo(fmt.Sprintf(callerInfoFormat, line, funcName)).
			withMessage(fmt.Sprint(v...))
		formatLogData(element)
		_ = stdErr.Output(calldepth2, element.outString)
	}
}

// PrintfErr 输出到os.Stderr
func PrintfErr(l ILog, format string, v ...interface{}) {
	if l != nil { // 日志已启用,使用日志打印
		l.Errorf(format, v...)
	} else {
		funcName := xconstants.Unknown
		pc, _, line, ok := runtime.Caller(calldepth1)
		if ok {
			funcName = runtime.FuncForPC(pc).Name()
		}
		element := newEntry().
			withLevel(LevelError).
			withTime(time.Now()).
			withCallerInfo(fmt.Sprintf(callerInfoFormat, line, funcName)).
			withMessage(fmt.Sprintf(format, v...))
		formatLogData(element)
		_ = stdErr.Output(calldepth2, element.outString)
	}
}
