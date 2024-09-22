package util

type ICallBack interface {
	CallBackFunc() error
	SetArg(arg interface{})
	GetArg() interface{}
}

// CallBackFunc 回调函数
type CallBackFunc func(arg interface{}) error

type defaultCallBack struct {
	Arg      interface{}  // 参数
	Function CallBackFunc // 回调函数
}

func NewDefaultCallBack(callBackFunc CallBackFunc, arg interface{}) ICallBack {
	return &defaultCallBack{
		Arg:      arg,
		Function: callBackFunc,
	}
}

func (p *defaultCallBack) CallBackFunc() error {
	return p.Function(p.GetArg())
}

func (p *defaultCallBack) SetArg(arg interface{}) {
	p.Arg = arg
}

func (p *defaultCallBack) GetArg() interface{} {
	return p.Arg
}
