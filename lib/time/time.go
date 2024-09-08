// 时间
// 程序运行过程中,会使用时间.计算时间.

package time

import (
	"time"
	"xcore/lib/util"
)

// Mgr 时间管理器
type Mgr struct {
	timestampSecond       int64        // 上一次调用Update更新的时间戳-秒
	timestampMillisecond  int64        // 上一次调用Update更新的时间戳-毫秒
	time                  time.Time    // 上一次调用Update更新的时间
	timestampSecondOffset int64        // 时间戳偏移量-秒
	utcSwitch             util.ISwitch // UTC 时间开关
}

func NewMgr() *Mgr {
	return &Mgr{
		utcSwitch: util.NewDefaultSwitch(false),
	}
}

// AbleUTC 使用UTC时间
func (p *Mgr) AbleUTC() {
	p.utcSwitch.Enable()
}

// DisableUTC 不使用UTC时间
func (p *Mgr) DisableUTC() {
	p.utcSwitch.Disable()
}

// NowTime 获取当前时间
func (p *Mgr) NowTime() time.Time {
	if p.utcSwitch.IsEnabled() {
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
