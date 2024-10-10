package service

import (
	"runtime"
	"runtime/debug"
	"time"
	xutil "xcore/lib/callback"
	xlog "xcore/lib/log"
	xtimer "xcore/lib/timer"
)

func StateStart(timer xtimer.ITimer) {
	defaultCallBack := xutil.NewDefaultCallBack(timeOut, timer)
	timer.AddSecond(defaultCallBack, time.Now().Unix()+1) //xconstants.ServiceInfoTimeOutSec)
}

// 服务信息 打印
func timeOut(arg ...interface{}) error {
	s := debug.GCStats{}
	debug.ReadGCStats(&s)

	xlog.PrintfInfo("goroutineCnt:%v, numGC:%d, lastGC:%v, GCPauseTotal:%v",
		runtime.NumGoroutine(), s.NumGC, s.LastGC, s.PauseTotal)

	StateStart(arg[0].(xtimer.ITimer))
	return nil
}
