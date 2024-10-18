package service

import (
	xcallback "xcore/lib/callback"
	xconstants "xcore/lib/constants"
	xetcd "xcore/lib/etcd"
)

// EtcdReportFunction etcd-上报
func EtcdReportFunction(args ...interface{}) error {
	defaultService := args[0].(*DefaultService)
	valueJson := &xetcd.ValueJson{
		ServiceNet:    &defaultService.BenchMgr.Json.ServiceNet,
		Version:       *defaultService.BenchMgr.Json.Base.Version,
		AvailableLoad: *defaultService.BenchMgr.Json.Base.AvailableLoad,
		SecondOffset:  0,
	}
	value := xetcd.ValueJson2String(valueJson)
	if _, err := defaultService.Etcd.Put(defaultService.EtcdKey, value); err != nil {
		defaultService.Log.Errorf("EtcdReportFunction Put key:%v val:%v err:%v", defaultService.EtcdKey, value, err)
	}
	defaultService.Timer.AddSecond(xcallback.NewDefaultCallBack(EtcdReportFunction, defaultService), defaultService.TimeMgr.ShadowTimestamp()+xconstants.EtcdReportIntervalSecondDefault)
	return nil
}
