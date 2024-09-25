package util

type ICallBack interface {
	Function() error
	SetArg(arg interface{})
	GetArg() interface{}
}

// CallbackFunction 回调函数
type CallbackFunction func(arg interface{}) error

type defaultCallBack struct {
	Arg      interface{}      // 参数
	Callback CallbackFunction // 回调函数
}

func NewDefaultCallBack(callbackFunction CallbackFunction, arg interface{}) ICallBack {
	return &defaultCallBack{
		Arg:      arg,
		Callback: callbackFunction,
	}
}

func (p *defaultCallBack) Function() error {
	return p.Callback(p.GetArg())
}

func (p *defaultCallBack) SetArg(arg interface{}) {
	p.Arg = arg
}

func (p *defaultCallBack) GetArg() interface{} {
	return p.Arg
}
