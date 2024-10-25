package service

import (
	"github.com/pkg/errors"
	xconstants "xcore/lib/constants"
	xcallback "xcore/lib/control"
	xetcd "xcore/lib/etcd"
	xruntime "xcore/lib/runtime"
	xutil "xcore/lib/util"
)

// EtcdReportFunction etcd-上报
func EtcdReportFunction(args ...interface{}) error {
	defaultService := args[0].(*DefaultService)
	defer func() {
		defaultService.Timer.AddSecond(xcallback.NewCallBack(EtcdReportFunction, defaultService), defaultService.TimeMgr.ShadowTimestamp()+xconstants.EtcdReportIntervalSecondDefault)
	}()
	cmdMin, err := xutil.HexStringToUint32(*defaultService.BenchMgr.Json.Base.CmdMin)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	cmdMax, err := xutil.HexStringToUint32(*defaultService.BenchMgr.Json.Base.CmdMax)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	valueJson := &xetcd.ValueJson{
		ServiceNet:    &defaultService.BenchMgr.Json.ServiceNet,
		Version:       *defaultService.BenchMgr.Json.Base.Version,
		AvailableLoad: *defaultService.BenchMgr.Json.Base.AvailableLoad,
		SecondOffset:  0,
		CmdMin:        cmdMin,
		CmdMax:        cmdMax,
	}
	value := xetcd.ValueJson2String(valueJson)
	if _, err := defaultService.Etcd.Put(defaultService.EtcdKey, value); err != nil {
		defaultService.Log.Errorf("EtcdReportFunction Put key:%v val:%v err:%v", defaultService.EtcdKey, value, err)
	}
	return nil
}
