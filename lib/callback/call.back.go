package callback

import "xcore/lib/parameters"

type ICallBack interface {
	Execute() error        // 执行回调
	parameters.IParameters // 参数
}

type defaultCallBack struct {
	onFunction func(arg ...interface{}) error // 回调函数
	parameters.IParameters
}

func NewDefaultCallBack(onFunction func(arg ...interface{}) error, arg ...interface{}) ICallBack {
	par := parameters.NewDefaultParameters()
	par.Set(arg...)
	return &defaultCallBack{
		onFunction:  onFunction,
		IParameters: par,
	}
}

func (p *defaultCallBack) Execute() error {
	if p.onFunction == nil {
		return nil
	}
	return p.onFunction(p.Get())
}
