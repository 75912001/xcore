package control

type CallBack struct {
	onFunction func(...interface{}) error // 回调函数
	IParameters
}

func NewCallBack(onFunction func(...interface{}) error, arg ...interface{}) ICallBack {
	par := NewParameters()
	par.Set(arg...)
	return &CallBack{
		onFunction:  onFunction,
		IParameters: par,
	}
}

func (p *CallBack) Execute() error {
	if p.onFunction == nil {
		return nil
	}
	return p.onFunction(p.IParameters.Get()...)
}
