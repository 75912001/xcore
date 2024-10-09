package callback

import "xcore/lib/parameters"

type ICallBack interface {
	Execute() error        // 执行回调
	parameters.IParameters // 参数
}
