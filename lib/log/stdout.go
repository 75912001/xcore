package log

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
	libconstants "xcore/lib/constants"
)

var stdOut = log.New(os.Stdout, "", 0)

// PrintInfo 输出到os.Stdout
func PrintInfo(v ...interface{}) {
	if isEnable() { // 日志已启用,使用日志打印
		mgrInstance.log(mgrInstance.newEntry(), LevelInfo, v...)
	} else {
		funcName := libconstants.Unknown
		pc, _, line, ok := runtime.Caller(calldepth1)
		if ok {
			funcName = runtime.FuncForPC(pc).Name()
		}
		formatAndPrint(stdOut, LevelInfo, line, funcName, v...)
	}
}

// PrintfInfo 输出到os.Stdout
func PrintfInfo(format string, v ...interface{}) {
	if isEnable() { // 日志已启用,需要放入日志 channel 中
		mgrInstance.logf(mgrInstance.newEntry(), LevelInfo, format, v...)
	} else {
		funcName := libconstants.Unknown
		pc, _, line, ok := runtime.Caller(calldepth1)
		if ok {
			funcName = runtime.FuncForPC(pc).Name()
		}
		formatAndPrint(stdOut, LevelInfo, line, funcName, v...)
	}
}

func formatAndPrint(logger *log.Logger, level int, line int, funcName string, v ...interface{}) {
	var buf bytes.Buffer
	buf.Grow(bufferCapacity)
	// 格式为  [时间][日志级别][TID:xxx][UID:xxx][堆栈信息][{extendFields-key:extendFields:val}...{}][自定义内容]
	buf.WriteString(fmt.Sprint("[", time.Now().Format(logTimeFormat), "]"))
	buf.WriteString(fmt.Sprint("[", levelDesc[level], "]"))
	buf.WriteString(traceIDKeyString0)
	buf.WriteString(userIDKeyString0)
	buf.WriteString(fmt.Sprint("[", fmt.Sprintf(callerInfoFormat, line, funcName), "]"))
	buf.WriteString("[]")
	buf.WriteString(fmt.Sprintf("%s%s%s", "[", fmt.Sprint(v...), "]"))
	_ = logger.Output(calldepth3, buf.String())
}
