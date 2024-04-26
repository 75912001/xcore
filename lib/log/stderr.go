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

var stdErr = log.New(os.Stderr, "", 0)

// PrintErr 输出到os.Stderr
func PrintErr(v ...interface{}) {
	if GetInstance().GetLevel() < LevelError {
		return
	}
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
		var buf bytes.Buffer
		buf.Grow(bufferCapacity)
		// 格式为  [时间][日志级别][TraceID:xxx][UID:xxx][堆栈信息][{extendFields-key:extendFields:val}...{}][自定义内容]
		buf.WriteString(fmt.Sprint("[", time.Now().Format(logTimeFormat), "]"))
		buf.WriteString(fmt.Sprint("[", levelDesc[LevelError], "]"))
		buf.WriteString("[TraceID:0]")
		buf.WriteString("[UID:0]")
		buf.WriteString(fmt.Sprint("[", fmt.Sprintf(callerInfoFormat, line, funcName), "]"))
		buf.WriteString("[]")
		buf.WriteString("[")
		buf.WriteString(fmt.Sprint(v...))
		buf.WriteString("]")
		_ = stdErr.Output(calldepth2, buf.String())
	}
}

// PrintfErr 输出到os.Stderr
func PrintfErr(format string, v ...interface{}) {
	if GetInstance().GetLevel() < LevelError {
		return
	}
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
		var buf bytes.Buffer
		buf.Grow(bufferCapacity)
		// 格式为  [时间][日志级别][TraceID:xxx][UID:xxx][堆栈信息][{extendFields-key:extendFields:val}...{}][自定义内容]
		buf.WriteString(fmt.Sprint("[", time.Now().Format(logTimeFormat), "]"))
		buf.WriteString(fmt.Sprint("[", levelDesc[LevelError], "]"))
		buf.WriteString("[TraceID:0]")
		buf.WriteString("[UID:0]")
		buf.WriteString(fmt.Sprint("[", fmt.Sprintf(callerInfoFormat, line, funcName), "]"))
		buf.WriteString("[]")
		buf.WriteString("[")
		buf.WriteString(fmt.Sprint(v...))
		buf.WriteString("]")
		_ = stdErr.Output(calldepth2, buf.String())
	}
}
