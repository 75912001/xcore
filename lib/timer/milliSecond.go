package timer

import (
	"xcore/lib/callback"
	xswitch "xcore/lib/switch"
)

// 毫秒级定时器
type millisecond struct {
	xswitch.ISwitch          // 有效(false:不执行,扫描时自动删除)
	callback.ICallBack       // 回调
	expire             int64 // 过期时间
}

func (p *millisecond) reset() {
	p.ISwitch.Disable()
	p.ICallBack = nil
	p.expire = 0
}
