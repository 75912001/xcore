package main

import (
	"context"
	"fmt"
	xlog "xcore/lib/log"
	xruntime "xcore/lib/runtime"
)

func logCallBackFunc(level uint32, outString string) {
	if xruntime.IsDebug() {
		fmt.Println(level, outString)
	}
	return
}

func exampleLog() {
	if false {
		xruntime.SetRunMode(xruntime.RunModeDebug)
		fmt.Println("============================================================")
		xlog.PrintInfo("print info")
		xlog.PrintfInfo("print info %s", "format")
		xlog.PrintErr("print err")
		xlog.PrintfErr("print err %s", "format")
		fmt.Println("============================================================")
		var log xlog.ILog
		log, err := xlog.NewMgr(xlog.NewOptions().
			WithLevelCallBack(logCallBackFunc, xlog.LevelFatal, xlog.LevelError, xlog.LevelWarn),
		)
		if err != nil {
			panic(err)
		}
		log.Fatal("fatal")
		log.Fatalf("fatal %s", "format")
		log.Error("error")
		log.Errorf("error %s", "format")
		log.Warn("warn")
		log.Warnf("warn %s", "format")
		log.Info("info")
		log.Infof("info %s", "format")
		log.Debug("debug")
		{
			log.DebugLazy(func() []interface{} {
				return []interface{}{fmt.Sprintf("%v %v", "This is a complex log message", "msg")}
			})
		}
		log.Debugf("debug %s", "format")
		{
			log.DebugfLazy(func() (string, []interface{}) {
				return "format %v %v", []interface{}{"This is a complex log message", 111}
			})
		}
		log.Trace("trace")
		{
			ctx := context.Background()
			ctx = context.WithValue(ctx, xlog.TraceIDKey, "TraceIDKey.value1")
			ctx = context.WithValue(ctx, xlog.UserIDKey, uint64(668))
			log.TraceExtend(ctx, xlog.ExtendFields{"key1", "value1", 1001, 1}, "trace")
		}
		{
			ctx := context.Background()
			ctx = context.WithValue(ctx, xlog.TraceIDKey, "TraceIDKey.value1")
			log.TraceExtend(ctx, xlog.ExtendFields{"key1", "value1", 1001, 1, xlog.UserIDKey, uint64(7200)}, "trace")
		}
		log.Tracef("trace %s", "format")
		{
			ctx := context.Background()
			ctx = context.WithValue(ctx, xlog.TraceIDKey, "TraceIDKey.value1")
			ctx = context.WithValue(ctx, xlog.UserIDKey, uint64(668))
			log.TracefExtend(ctx, xlog.ExtendFields{"key1", "value1", 1001, 1}, "trace %s", "format")
		}
		err = log.Stop()
		if err != nil {
			panic(err)
		}
	}
	return
}
