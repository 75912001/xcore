package timer

import xutil "xcore/lib/util"

// Millisecond 毫秒级定时器
type Millisecond struct {
	xutil.ICallBack       // 回调
	xutil.ISwitch         // 有效(false:不执行,扫描时自动删除)
	expire          int64 // 过期时间
}

func (p *Millisecond) reset() {
	p.ICallBack = nil
	p.ISwitch.Disable()
	p.expire = 0
}
