package log

import (
	xutil "xcore/lib/util"
)

// IHook 钩子
// 钩子是一个接口，它定义了x个日志级别的钩子函数
// e.g.: 修改entry的内容, 将 entry 发送至其他目标地
type IHook interface {
	Fire(outString string) error // 执行的方法 (outString: 输出的字符串)
	xutil.ICallBackFunc
}

type levelHook struct {
	levelMap map[uint32]struct{}
	xutil.ICallBackFunc
}

func NewLevelHook(fun xutil.CallBackFunc) *levelHook {
	return &levelHook{
		levelMap:      make(map[uint32]struct{}),
		ICallBackFunc: xutil.NewDefaultCallBackFunc(fun, nil),
	}
}

func (p *levelHook) AddLevel(levelSlice []uint32) {
	for _, level := range levelSlice {
		p.levelMap[level] = struct{}{}
	}
}

func (p *levelHook) Fire(outString string) error {
	p.SetArg(outString)
	return p.CallBack()
}

//func (p *levelHook) CallBack() error {
//	_, ok := p.levelMap[level]
//	return ok
//}
//
//func (p *levelHook) GetArg() interface{} {
//	p.levelMap[level] = struct{}{}
//}

//// fire 处理钩子
//func (p levelHookMap) fire(entry *entry) error {
//	for _, hook := range p[entry.level] {
//		if err := hook.Fire(entry.outString); err != nil {
//			return errors.WithMessage(err, xruntime.Location())
//		}
//	}
//	return nil
//}
