package server

import (
	xcontrol "github.com/75912001/xcore/lib/control"
	xetcd "github.com/75912001/xcore/lib/etcd"
)

// etcdReportFunction etcd-上报
func etcdReportFunction(args ...interface{}) error {
	defaultServer := args[0].(*Server)
	defer func() {
		defaultServer.Timer.AddSecond(
			xcontrol.NewCallBack(etcdReportFunction, defaultServer),
			defaultServer.TimeMgr.ShadowTimestamp()+xetcd.ReportIntervalSecondDefault,
		)
	}()
	valueJson := &xetcd.ValueJson{
		ServerNet:     &defaultServer.BenchMgr.Json.ServerNet,
		Version:       *defaultServer.BenchMgr.Json.Base.Version,
		AvailableLoad: *defaultServer.BenchMgr.Json.Base.AvailableLoad,
		SecondOffset:  0,
	}
	value := xetcd.ValueJson2String(valueJson)
	if _, err := defaultServer.Etcd.Put(defaultServer.EtcdKey, value); err != nil {
		defaultServer.Log.Errorf("etcdReportFunction Put key:%v val:%v err:%v", defaultServer.EtcdKey, value, err)
	}
	return nil
}
