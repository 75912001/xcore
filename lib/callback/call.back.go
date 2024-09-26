package callback

import "xcore/lib/parameters"

type ICallBack interface {
	Execute() error        // 执行回调
	parameters.IParameters // 参数
}

type callBack struct {
	onFunction func(arg ...interface{}) error // 回调函数
	parameters.IParameters
}

func NewCallBack(onFunction func(arg ...interface{}) error, arg ...interface{}) ICallBack {
	par := parameters.NewParameters()
	par.Set(arg...)
	return &callBack{
		onFunction:  onFunction,
		IParameters: par,
	}
}

func (p *callBack) Execute() error {
	return p.onFunction(p.Get())
}
