package timer

import (
	"context"
	xcontrol "xcore/lib/control"
)

type ITimerSecond interface {
	AddSecond(callBackFunc xcontrol.ICallBack, expire int64) *second
	DelSecond(second *second)
}

type ITimerMillisecond interface {
	AddMillisecond(callBackFunc xcontrol.ICallBack, expireMillisecond int64) *millisecond
	DelMillisecond(millisecond *millisecond)
}

type ITimer interface {
	Start(ctx context.Context, opts ...*options) error
	Stop()
	ITimerSecond
	ITimerMillisecond
}

type EventTimerSecond struct {
	ISwitch   xcontrol.ISwitchButton
	ICallBack xcontrol.ICallBack
}

type EventTimerMillisecond struct {
	ISwitch   xcontrol.ISwitchButton
	ICallBack xcontrol.ICallBack
}
