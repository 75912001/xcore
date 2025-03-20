package pprof

import (
	xerror "github.com/75912001/xcore/lib/error"
	xlog "github.com/75912001/xcore/lib/log"
	xruntime "github.com/75912001/xcore/lib/runtime"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
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
