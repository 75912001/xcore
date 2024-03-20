package time

import (
	"sync"
	"time"
)

var (
	instance *mgr
	once     sync.Once
)

// 获取实例
func getInstance() *mgr {
	once.Do(func() {
		instance = &mgr{}
	})
	return instance
}

// mgr 时间管理器
type mgr struct {
	TimestampSecond       int64     //上一次调用Update更新的时间戳-秒
	TimestampMillisecond  int64     //上一次调用Update更新的时间戳-毫秒
	Time                  time.Time //上一次调用Update更新的时间
	TimestampSecondOffset int64     //时间戳偏移量-秒
	utcAble               bool      //是否使用UTC时间
}

func AbleUTC() {
	getInstance().utcAble = true
}

func DisableUTC() {
	getInstance().utcAble = false
}

func NowTime() time.Time {
	if getInstance().utcAble {
		return time.Now().UTC()
	}
	return time.Now()
}

// Update 更新时间管理器中的,当前时间
func Update() {
	getInstance().Time = NowTime()
	getInstance().TimestampSecond = getInstance().Time.Unix()
	getInstance().TimestampMillisecond = getInstance().Time.UnixMilli()
}

// ShadowTimestampSecond 叠加偏移量的时间戳-秒
func ShadowTimestampSecond() int64 {
	return getInstance().TimestampSecond + getInstance().TimestampSecondOffset
}

// SetTimestampSecondOffset 设置 时间戳偏移量-秒
func SetTimestampSecondOffset(offset int64) {
	getInstance().TimestampSecondOffset = offset
}
