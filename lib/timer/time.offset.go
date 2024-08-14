package timer

import "time"

var secondOffset int64 // 时间偏移量-秒

// SetSecondOffset 设置时间偏移量-秒
func SetSecondOffset(offset int64) {
	secondOffset = offset
}

// ShadowTimestamp 叠加偏移量的时间戳-秒
func ShadowTimestamp() int64 {
	return time.Now().Unix() + secondOffset
}
