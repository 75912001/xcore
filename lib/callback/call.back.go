package callback

import xcontrol "xcore/lib/control"

type ICallBack interface {
	Execute() error      // 执行回调
	xcontrol.IParameters // 参数
}
