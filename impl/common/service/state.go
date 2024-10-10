package service

import (
	"runtime"
	"runtime/debug"
	xutil "xcore/lib/callback"
	xlog "xcore/lib/log"
	xtime "xcore/lib/time"
	xtimer "xcore/lib/timer"
)

func StateStart(timer xtimer.ITimer, timeMgr *xtime.Mgr) {
	defaultCallBack := xutil.NewDefaultCallBack(timeOut, timer, timeMgr)
	timer.AddSecond(defaultCallBack, timeMgr.ShadowTimestamp()+1) //xconstants.ServiceInfoTimeOutSec)
}

// 服务信息 打印
func timeOut(arg ...interface{}) error {
	s := debug.GCStats{}
	debug.ReadGCStats(&s)

	xlog.PrintfInfo("goroutineCnt:%v, numGC:%d, lastGC:%v, GCPauseTotal:%v",
		runtime.NumGoroutine(), s.NumGC, s.LastGC, s.PauseTotal)

	StateStart(arg[0].(xtimer.ITimer), arg[1].(*xtime.Mgr))
	return nil
}
