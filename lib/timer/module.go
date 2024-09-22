package timer

import (
	"context"
	xutil "xcore/lib/util"
)

type ITimer interface {
	Start(ctx context.Context, opts ...*options) error
	Stop()
	AddMillisecond(callBackFunc xutil.ICallBack, expireMillisecond int64) *Millisecond
	DelMillisecond(millisecond *Millisecond)
	AddSecond(callBackFunc xutil.ICallBack, expire int64) *Second
	DelSecond(second *Second)
}
