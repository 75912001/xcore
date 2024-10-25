package timer

import (
	xcallback "xcore/lib/callback"
	xcontrol "xcore/lib/control"
)

// 毫秒级定时器
type millisecond struct {
	ISwitch   xcontrol.ISwitchButton // 有效(false:不执行,扫描时自动删除)
	ICallBack xcallback.ICallBack    // 回调
	expire    int64                  // 过期时间
}

func (p *millisecond) reset() {
	p.ISwitch.Off()
	p.ICallBack = nil
	p.expire = 0
}
