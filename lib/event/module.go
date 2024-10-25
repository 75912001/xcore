package event

import xcallback "xcore/lib/callback"

type IEvent interface {
	Trigger() error
}

type defaultEvent struct {
	ID uint32
	xcallback.ICallBack
}

func (p *defaultEvent) Trigger() error {

	return nil
}

func NewDefaultEvent(ID uint32, onFunction func(arg ...interface{}) error, arg ...interface{}) IEvent {
	return &defaultEvent{
		ID:        ID,
		ICallBack: xcallback.NewCallBack(onFunction, arg...),
	}
}
