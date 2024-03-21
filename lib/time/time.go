package time

import (
	"time"
)

// Mgr 时间管理器
type Mgr struct {
	timestampSecond       int64     //上一次调用Update更新的时间戳-秒
	timestampMillisecond  int64     //上一次调用Update更新的时间戳-毫秒
	time                  time.Time //上一次调用Update更新的时间
	timestampSecondOffset int64     //时间戳偏移量-秒
	utcAble               bool      //是否使用UTC时间
}

func (p *Mgr) AbleUTC() {
	p.utcAble = true
}

func (p *Mgr) DisableUTC() {
	p.utcAble = false
}

// NowTime 获取当前时间
func (p *Mgr) NowTime() time.Time {
	if p.utcAble {
		return time.Now().UTC()
	}
	return time.Now()
}

// Update 更新时间管理器中的,当前时间
func (p *Mgr) Update() {
	p.time = p.NowTime()
	p.timestampSecond = p.time.Unix()
	p.timestampMillisecond = p.time.UnixMilli()
}

// ShadowTimestampSecond 叠加偏移量的时间戳-秒
func (p *Mgr) ShadowTimestampSecond() int64 {
	return p.timestampSecond + p.timestampSecondOffset
}

// SetTimestampSecondOffset 设置 时间戳偏移量-秒
func (p *Mgr) SetTimestampSecondOffset(offset int64) {
	p.timestampSecondOffset = offset
}
