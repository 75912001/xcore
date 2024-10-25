package event

type IEvent interface {
	Trigger() error
}

type defaultEvent struct {
	ID uint32
	callback.ICallBack
}

func (p *defaultEvent) Trigger() error {

	return nil
}

func NewDefaultEvent(ID uint32, onFunction func(arg ...interface{}) error, arg ...interface{}) IEvent {
	return &defaultEvent{
		ID:        ID,
		ICallBack: callback.NewCallBack(onFunction, arg...),
	}
}
