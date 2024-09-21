package util

type ICallBackFunc interface {
	CallBack() error
	GetArg() interface{}
}

// CallBackFunc 回调函数
type CallBackFunc func(arg interface{})

type defaultCallBackFunc struct {
	Arg      interface{}  // 参数
	Function CallBackFunc // 回调函数
}

func NewDefaultCallBackFunc(fun CallBackFunc, arg interface{}) ICallBackFunc {
	return &defaultCallBackFunc{
		Arg:      arg,
		Function: fun,
	}
}

func (p *defaultCallBackFunc) CallBack() error {
	p.Function(p.GetArg())
	return nil
}

func (p *defaultCallBackFunc) GetArg() interface{} {
	return p.Arg
}
