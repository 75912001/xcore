package timer

import (
	"context"
	xcallback "xcore/lib/callback"
	xswitch "xcore/lib/xswitch"
)

type ITimerSecond interface {
	AddSecond(callBackFunc xcallback.ICallBack, expire int64) *second
	DelSecond(second *second)
}

type ITimerMillisecond interface {
	AddMillisecond(callBackFunc xcallback.ICallBack, expireMillisecond int64) *millisecond
	DelMillisecond(millisecond *millisecond)
}

type ITimer interface {
	Start(ctx context.Context, opts ...*options) error
	Stop()
	ITimerSecond
	ITimerMillisecond
}

type EventTimerSecond struct {
	ISwitch   xswitch.ISwitch
	ICallBack xcallback.ICallBack
}

type EventTimerMillisecond struct {
	ISwitch   xswitch.ISwitch
	ICallBack xcallback.ICallBack
}
