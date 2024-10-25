package service

import (
	"runtime"
	"runtime/debug"
	"time"
	xcallback "xcore/lib/callback"
	xconstants "xcore/lib/constants"
	xlog "xcore/lib/log"
	xtimer "xcore/lib/timer"
)

func StateTimerPrint(timer xtimer.ITimer, l xlog.ILog) {
	defaultCallBack := xcallback.NewCallBack(timeOut, timer, l)
	_ = timer.AddSecond(defaultCallBack, time.Now().Unix()+xconstants.ServiceInfoTimeOutSec)
}

// 服务信息 打印
func timeOut(arg ...interface{}) error {
	s := debug.GCStats{}
	debug.ReadGCStats(&s)
	l := arg[1].(xlog.ILog)
	l.Infof("goroutineCnt:%v, numGC:%d, lastGC:%v, GCPauseTotal:%v",
		runtime.NumGoroutine(), s.NumGC, s.LastGC, s.PauseTotal)
	//xlog.PrintfInfo("goroutineCnt:%v, numGC:%d, lastGC:%v, GCPauseTotal:%v",
	//	runtime.NumGoroutine(), s.NumGC, s.LastGC, s.PauseTotal)
	StateTimerPrint(arg[0].(xtimer.ITimer), l)
	return nil
}
