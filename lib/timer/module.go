package timer

import (
	"context"
	xutil "xcore/lib/callback"
)

type ITimer interface {
	Start(ctx context.Context, opts ...*options) error
	Stop()
	AddMillisecond(callBackFunc xutil.ICallBack, expireMillisecond int64) *millisecond
	DelMillisecond(millisecond *millisecond)
	AddSecond(callBackFunc xutil.ICallBack, expire int64) *second
	DelSecond(second *second)
}

type EventTimerSecond struct {
	xutil.ICallBack
}

type EventTimerMillisecond struct {
	xutil.ICallBack
}
