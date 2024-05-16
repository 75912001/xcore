package log

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	"strconv"
	"time"
	libconstants "xcore/lib/constants"
)

//日志条目

// 日志数据字段,扩展字段
type extendFields []interface{}

// 日志数据信息
type entry struct {
	mgr          *mgr      //日志管理器
	level        int       //本条目的日志级别
	time         time.Time //生成日志的时间
	callerInfo   string    //调用堆栈信息
	message      string    //日志消息
	ctx          context.Context
	extendFields extendFields //[string,interface{}] key,value;key,value...
}

func (p *entry) reset() {
	p.mgr = nil
	p.level = LevelOff
	p.callerInfo = ""
	p.message = ""
	p.ctx = nil
	p.extendFields = nil
}

// 由ctx创建Entry
func (p *entry) withContext(ctx context.Context) *entry {
	p.ctx = ctx
	return p
}

// 由field创建Entry
func (p *entry) withExtendField(key string, value interface{}) *entry {
	if p.extendFields == nil {
		p.extendFields = make(extendFields, 0, 4)
	}
	p.extendFields = append(p.extendFields, key, value)
	return p
}

// 由多个field创建Entry
func (p *entry) withExtendFields(fields extendFields) *entry {
	if p.extendFields == nil {
		p.extendFields = make(extendFields, 0, 8)
	}
	p.extendFields = append(p.extendFields, fields...)
	return p
}

func (p *entry) withMessage(message string) *entry {
	p.message = message
	return p
}

// 格式化日志信息
func (p *entry) formatMessage() string {
	// 格式为  [时间][日志级别][TraceID:xxx][UID:xxx][堆栈信息][{extendFields-key:extendFields:val}...{}][自定义内容]
	var buf bytes.Buffer
	buf.Grow(bufferCapacity)
	// 时间
	buf.WriteString(fmt.Sprint("[", p.time.Format(logTimeFormat), "]"))
	// 日志级别
	buf.WriteString(fmt.Sprint("[", levelDesc[p.level], "]"))
	// 处理 ctx TraceID
	if p.ctx != nil {
		traceIdVal := p.ctx.Value(TraceIDKey)
		if traceIdVal != nil {
			buf.WriteString(fmt.Sprint("[", TraceIDKey, ":", traceIdVal.(string), "]"))
		}
	}
	// UID 优先从 ctx 查找,其次查找 field 当 UID 不存在时也需要占位 值为0
	var uid uint64
	if p.ctx != nil {
		uidVal := p.ctx.Value(UserIDKey)
		if uidVal != nil {
			uid, _ = uidVal.(uint64)
		}
	}
	if 0 == uid { //没有找到UID,从field中查找
		for _, v := range p.extendFields {
			str, ok := v.(string)
			if ok && str == UserIDKey { //找到
				uid, _ = v.(uint64)
				break
			}
		}
	}
	buf.WriteString(fmt.Sprint("[", UserIDKey, ":", strconv.FormatUint(uid, 10), "]"))
	// 堆栈
	buf.WriteString(fmt.Sprint("[", p.callerInfo, "]"))
	// 处理fields
	buf.WriteString(fmt.Sprint("["))
	for idx, v := range p.extendFields {
		if idx%2 == 0 { //key
			buf.WriteString("{")
			str, ok := v.(string)
			if ok {
				buf.WriteString(str)
			} else {
				buf.WriteString(fmt.Sprint(v))
			}
			buf.WriteString(":")
		} else { //val
			str, ok := v.(string)
			if ok {
				buf.WriteString(str)
			} else {
				buf.WriteString(fmt.Sprint(v))
			}
			buf.WriteString("}")
		}
	}
	buf.WriteString(fmt.Sprint("]"))
	// 自定义内容
	buf.WriteString(p.message)
	return buf.String()
}

// 记录日志
func (p *entry) log(level int, skip int, v ...interface{}) {
	p.level = level
	p.time = p.mgr.timeMgr.NowTime()
	if *p.mgr.options.isReportCaller {
		pc, _, line, ok := runtime.Caller(skip)
		funcName := libconstants.Unknown
		if !ok {
			line = 0
		} else {
			funcName = runtime.FuncForPC(pc).Name()
		}
		p.callerInfo = fmt.Sprintf(callerInfoFormat, line, funcName)
	}
	p.message = fmt.Sprintln(v...)
	p.mgr.logChan <- p
}

// 记录日志
func (p *entry) logf(level int, skip int, format string, v ...interface{}) {
	p.level = level
	p.time = p.mgr.timeMgr.NowTime()
	if *p.mgr.options.isReportCaller {
		pc, _, line, ok := runtime.Caller(skip)
		funcName := libconstants.Unknown
		if !ok {
			line = 0
		} else {
			funcName = runtime.FuncForPC(pc).Name()
		}
		p.callerInfo = fmt.Sprintf(callerInfoFormat, line, funcName)
	}
	p.message = fmt.Sprintf(format, v...)
	p.mgr.logChan <- p
}

// Trace 追踪日志
func (p *entry) Trace(v ...interface{}) {
	if p.mgr.GetLevel() < LevelTrace {
		return
	}
	p.log(LevelTrace, calldepth2, v...)
}

// Tracef 追踪日志
func (p *entry) Tracef(format string, v ...interface{}) {
	if p.mgr.GetLevel() < LevelTrace {
		return
	}
	p.logf(LevelTrace, calldepth2, format, v...)
}

// Debug 调试日志
func (p *entry) Debug(v ...interface{}) {
	if p.mgr.GetLevel() < LevelDebug {
		return
	}
	p.log(LevelDebug, calldepth2, v...)
}

// Debugf 调试日志
func (p *entry) Debugf(format string, v ...interface{}) {
	if p.mgr.GetLevel() < LevelDebug {
		return
	}
	p.logf(LevelDebug, calldepth2, format, v...)
}

// Info 信息日志
func (p *entry) Info(v ...interface{}) {
	if p.mgr.GetLevel() < LevelInfo {
		return
	}
	p.log(LevelInfo, calldepth2, v...)
}

// Infof 信息日志
func (p *entry) Infof(format string, v ...interface{}) {
	if p.mgr.GetLevel() < LevelInfo {
		return
	}
	p.logf(LevelInfo, calldepth2, format, v...)
}

// Warn 警告日志
func (p *entry) Warn(v ...interface{}) {
	if p.mgr.GetLevel() < LevelWarn {
		return
	}
	p.log(LevelWarn, calldepth2, v...)
}

// Warnf 警告日志
func (p *entry) Warnf(format string, v ...interface{}) {
	if p.mgr.GetLevel() < LevelWarn {
		return
	}
	p.logf(LevelWarn, calldepth2, format, v...)
}

// Error 错误日志
func (p *entry) Error(v ...interface{}) {
	if p.mgr.GetLevel() < LevelError {
		return
	}
	p.log(LevelError, calldepth2, v...)
}

// Errorf 错误日志
func (p *entry) Errorf(format string, v ...interface{}) {
	if p.mgr.GetLevel() < LevelError {
		return
	}
	p.logf(LevelError, calldepth2, format, v...)
}

// Fatal 致命日志
func (p *entry) Fatal(v ...interface{}) {
	if p.mgr.GetLevel() < LevelFatal {
		return
	}
	p.log(LevelFatal, calldepth2, v...)
}

// Fatalf 致命日志
func (p *entry) Fatalf(format string, v ...interface{}) {
	if p.mgr.GetLevel() < LevelFatal {
		return
	}
	p.logf(LevelFatal, calldepth2, format, v...)
}
