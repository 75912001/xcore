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
func PrintInfo(v ...interface{}) {
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

// PrintfInfo 输出到os.Stdout
func PrintfInfo(format string, v ...interface{}) {
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
