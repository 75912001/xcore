package timer

import (
	"context"
	xutil "xcore/lib/callback"
	xswitch "xcore/lib/switch"
)

type ITimerSecond interface {
	AddSecond(callBackFunc xutil.ICallBack, expire int64) *second
	DelSecond(second *second)
}

type ITimerMillisecond interface {
	AddMillisecond(callBackFunc xutil.ICallBack, expireMillisecond int64) *millisecond
	DelMillisecond(millisecond *millisecond)
}

type ITimer interface {
	Start(ctx context.Context, opts ...*options) error
	Stop()
	ITimerSecond
	ITimerMillisecond
}

type EventTimerSecond struct {
	xswitch.ISwitch
	xutil.ICallBack
}

type EventTimerMillisecond struct {
	xswitch.ISwitch
	xutil.ICallBack
}
