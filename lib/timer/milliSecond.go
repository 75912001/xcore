package timer

import (
	"xcore/lib/callback"
	xutil "xcore/lib/util"
)

// 毫秒级定时器
type millisecond struct {
	callback.ICallBack       // 回调
	xutil.ISwitch            // 有效(false:不执行,扫描时自动删除)
	expire             int64 // 过期时间
}

func (p *millisecond) reset() {
	p.ICallBack = nil
	p.ISwitch.Disable()
	p.expire = 0
}
