package callback

type ICallBack interface {
	Function() error
	SetArg(arg interface{})
	GetArg() interface{}
}

// OnFunction 回调函数
type OnFunction func(arg interface{}) error

type callBack struct {
	arg        interface{} // 参数
	onFunction OnFunction  // 回调函数
}

func NewCallBack(onFunction OnFunction, arg interface{}) ICallBack {
	return &callBack{
		arg:        arg,
		onFunction: onFunction,
	}
}

func (p *callBack) Function() error {
	return p.onFunction(p.GetArg())
}

func (p *callBack) SetArg(arg interface{}) {
	p.arg = arg
}

func (p *callBack) GetArg() interface{} {
	return p.arg
}
