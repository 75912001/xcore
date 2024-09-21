package timer

import xutil "xcore/lib/util"

// Millisecond 毫秒级定时器
type Millisecond struct {
	xutil.ICallBackFunc       // 回调
	xutil.ISwitch             // 有效(false:不执行,扫描时自动删除)
	expire              int64 // 过期时间
}

func (p *Millisecond) reset() {
	p.ICallBackFunc = nil
	p.ISwitch.Enable()
	p.expire = 0
}

// DelMillisecond 删除毫秒级定时器
//
//	[NOTE] 必须与该 outgoingTimeoutChan 线性处理.如:在同一个 goroutine select 中处理数据
//	参数:
//		毫秒定时器
func DelMillisecond(t *Millisecond) {
	t.Disable()
}
