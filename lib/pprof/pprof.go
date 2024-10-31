package pprof

import (
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
	xerror "xcore/lib/error"
	xlog "xcore/lib/log"
	xruntime "xcore/lib/runtime"
)

// StartHTTPprof 开启http采集分析
//
//	参数:
//		addr: "0.0.0.0:8090"
func StartHTTPprof(addr string) {
	go func() {
		defer func() {
			if xruntime.IsRelease() {
				if err := recover(); err != nil {
					xlog.PrintErr(xerror.GoroutinePanic, err, xruntime.Location(), debug.Stack())
				}
			}
			xlog.PrintInfo(xerror.GoroutineDone)
		}()
		if err := http.ListenAndServe(addr, nil); err != nil {
			xlog.PrintErr(xerror.Fail, addr, err)
		}
	}()
}
