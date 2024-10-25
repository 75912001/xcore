package timer

import (
	xcontrol "xcore/lib/control"
)

// 秒级定时器
type second struct {
	ISwitch   xcontrol.ISwitchButton // 有效(false:不执行,扫描时自动删除)
	ICallBack xcontrol.ICallBack     // 回调
	expire    int64                  // 过期时间
}

func (p *second) reset() {
	p.ISwitch.Off()
	p.ICallBack = nil
	p.expire = 0
}
