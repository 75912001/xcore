package service

import (
	"runtime"
	"runtime/debug"
	"time"
	xutil "xcore/lib/callback"
	xconstants "xcore/lib/constants"
	xlog "xcore/lib/log"
	xtimer "xcore/lib/timer"
)

func StateTimerPrint(timer xtimer.ITimer) {
	defaultCallBack := xutil.NewDefaultCallBack(timeOut, timer)
	_ = timer.AddSecond(defaultCallBack, time.Now().Unix()+xconstants.ServiceInfoTimeOutSec)
}

// 服务信息 打印
func timeOut(arg ...interface{}) error {
	s := debug.GCStats{}
	debug.ReadGCStats(&s)
	xlog.PrintfInfo("goroutineCnt:%v, numGC:%d, lastGC:%v, GCPauseTotal:%v",
		runtime.NumGoroutine(), s.NumGC, s.LastGC, s.PauseTotal)
	StateTimerPrint(arg[0].(xtimer.ITimer))
	return nil
}
