package log

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"
)

//日志条目

// 日志数据字段,扩展字段
type extendFields []interface{}

// 日志数据信息
type entry struct {
	level        uint32    // 本条目的日志级别
	time         time.Time // 生成日志的时间
	callerInfo   string    // 调用堆栈信息
	message      string    // 日志消息
	ctx          context.Context
	extendFields extendFields // [string,interface{}] key,value;key,value...
}

func newEntry() *entry {
	return &entry{}
}

func (p *entry) reset() {
	p.level = LevelOff
	p.callerInfo = ""
	p.message = ""
	p.ctx = nil
	p.extendFields = nil
}

func withLevel(p *entry, level uint32) *entry {
	p.level = level
	return p
}

func withTime(p *entry, nowTime time.Time) *entry {
	p.time = nowTime
	return p
}

func withCallerInfo(p *entry, callerInfo string) *entry {
	p.callerInfo = callerInfo
	return p
}

func withMessage(p *entry, message string) *entry {
	p.message = message
	return p
}

func (p *entry) withContext(ctx context.Context) *entry {
	p.ctx = ctx
	return p
}

func (p *entry) withExtendField(key string, value interface{}) *entry {
	if p.extendFields == nil {
		p.extendFields = make(extendFields, 0, 4)
	}
	p.extendFields = append(p.extendFields, key, value)
	return p
}

func (p *entry) withExtendFields(fields extendFields) *entry {
	if p.extendFields == nil {
		fieldsSize := len(fields)
		p.extendFields = make(extendFields, 0, fieldsSize)
	}
	p.extendFields = append(p.extendFields, fields...)
	return p
}

// 格式化日志数据
func formatLogData(p *entry) string {
	// 格式为  [时间][日志级别][TID:xxx][UID:xxx][堆栈信息][{extendFields-key:extendFields:val}...{}][日志消息]
	var buf bytes.Buffer
	buf.Grow(bufferCapacity)
	// 时间
	buf.WriteString(fmt.Sprint("[", p.time.Format(logTimeFormat), "]"))
	// 日志级别
	buf.WriteString(fmt.Sprint("[", levelDesc[p.level], "]"))
	// TraceID
	if p.ctx != nil { // 处理 ctx 中的 traceID
		traceIdVal := p.ctx.Value(traceIDKey)
		if traceIdVal != nil {
			buf.WriteString(fmt.Sprint("[", traceIDKey, ":", traceIdVal.(string), "]"))
		}
	} else { // 没有 ctx , 则 traceID 为0
		buf.WriteString(fmt.Sprint("[", traceIDKey, ":0]"))
	}
	// UID 优先从 ctx 查找,其次查找 field 当 UID 不存在时也需要占位 值为0
	var uid uint64
	if p.ctx != nil {
		uidVal := p.ctx.Value(userIDKey)
		if uidVal != nil {
			uid, _ = uidVal.(uint64)
		}
	}
	if 0 == uid { //没有找到UID,从fields中查找,找到第一个
		for idx, v := range p.extendFields {
			str, ok := v.(string)
			if ok && str == userIDKey { //找到
				uid, _ = p.extendFields[idx+1].(uint64)
				break
			}
		}
	}
	buf.WriteString(fmt.Sprint("[", userIDKey, ":", strconv.FormatUint(uid, 10), "]"))
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
	// 日志消息
	buf.WriteString(p.message)
	return buf.String()
}
