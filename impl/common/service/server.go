package service

//
//import (
//	"context"
//	"dawn-server/impl/common"
//	"dawn-server/impl/common/db_retry"
//	"dawn-server/impl/common/etcd"
//	"dawn-server/impl/common/mq_nats"
//	"dawn-server/impl/common/statistics/st_mgr"
//	xretcd "dawn-server/impl/xr/lib/etcd"
//	xrjaeger "dawn-server/impl/xr/lib/jaeger"
//	xrlog "dawn-server/impl/xr/lib/log"
//	xrtcp "dawn-server/impl/xr/lib/tcp"
//	xrutil "dawn-server/impl/xr/lib/util"
//	"fmt"
//	"github.com/pkg/errors"
//	"runtime"
//	"sync"
//	"sync/atomic"
//	"time"
//)
//
//const StatusRunning = 0  // 服务状态：运行中
//const StatusStopping = 1 // 服务状态：关闭中
//
//var GGlobalNats mq_nats.Mgr
//var GZoneNats mq_nats.Mgr
//var GTcpService xrtcp.Server
//
//var GBusChannelWaitGroup sync.WaitGroup
//var GBusChannelCheckChan = make(chan struct{}, 1)
//
//var GDBRetryMgr db_retry.DBRetryMgr //world
//
//var GServerStatus uint32
//
//var GST st_mgr.STMgr
//var GQuitChan = make(chan bool)
//
//// IsServerStopping 服务是否关闭中
//func IsServerStopping() bool {
//	return atomic.LoadUint32(&GServerStatus) == StatusStopping
//}
//
//// IsServerRunning 服务是否运行中
//func IsServerRunning() bool {
//	return atomic.LoadUint32(&GServerStatus) == StatusRunning
//}
//
//// SetServerStopping 设置为关闭中
//func SetServerStopping() {
//	atomic.StoreUint32(&GServerStatus, StatusStopping)
//}
//
//// GetDstNATS 获取数据目的地
//func GetDstNATS(zoneID uint32) *mq_nats.Mgr {
//	if common.IsZoneSet(GMgr.ZoneID, zoneID) {
//		return &GZoneNats
//	}
//	return &GGlobalNats
//}
//
//// GetSubj  获取 subj
//func GetSubj(zoneID uint32, serviceName string, worldID uint32) (subj string) {
//	if GMgr.ZoneID == zoneID {
//		return mq_nats.GenZoneSub(zoneID, serviceName, worldID)
//	}
//	return mq_nats.GenGlobalSub(zoneID, serviceName, worldID)
//}
//
//func (p *Mgr) PostInit(ctx context.Context, opts ...*options) error {
//	var err error
//
//	p.Opt = mergeOptions(opts...)
//	err = configure(p.Opt)
//	if err != nil {
//		return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
//	}
//
//	// 启动Etcd
//	p.Bench.BenchRootJson.Etcd.Key = fmt.Sprintf("/%v/%v/%v/%v/%v",
//		common.ProjectName, common.EtcdWatchMsgTypeService, p.ZoneID, p.ServiceName, p.ServiceID)
//	err = etcd.EtcdStart(&p.Bench.BenchRootJson.Etcd, p.BusChannel, p.Opt.EtcdHandler)
//	if err != nil {
//		return errors.Errorf("Etcd start err:%v %v", err, xrutil.GetCodeLocation(1).String())
//	}
//
//	// etcd 关注 服务 首次启动服务需要拉取一次
//	if p.Opt.EtcdWatchServicePrefix != nil {
//		if err = xretcd.GetInstance().WatchPrefixIntoChan(*p.Opt.EtcdWatchServicePrefix); err != nil {
//			return errors.Errorf("EtcdWatchPrefix err:%v %v", err, xrutil.GetCodeLocation(1).String())
//		}
//		if err = xretcd.GetInstance().GetPrefixIntoChan(*p.Opt.EtcdWatchServicePrefix); err != nil {
//			return errors.Errorf("EtcdGetPrefix err:%v %v", err, xrutil.GetCodeLocation(1).String())
//		}
//	}
//
//	// etcd 关注 命令
//	if p.Opt.EtcdWatchCommandPrefix != nil {
//		if err = xretcd.GetInstance().WatchPrefixIntoChan(*p.Opt.EtcdWatchCommandPrefix); err != nil {
//			return errors.Errorf("EtcdWatchPrefix err:%v %v", err, xrutil.GetCodeLocation(1).String())
//		}
//	}
//
//	// etcd 关注 GM
//	if p.Opt.EtcdWatchGMPrefix != nil {
//		if err = xretcd.GetInstance().WatchPrefixIntoChan(*p.Opt.EtcdWatchGMPrefix); err != nil {
//			return errors.Errorf("EtcdWatchPrefix err:%v %v", err, xrutil.GetCodeLocation(1).String())
//		}
//		if err = xretcd.GetInstance().GetPrefixIntoChan(*p.Opt.EtcdWatchGMPrefix); err != nil {
//			return errors.Errorf("EtcdGetPrefix err:%v %v", err, xrutil.GetCodeLocation(1).String())
//		}
//	}
//
//	runtime.GC()
//	return nil
//}
//
//func (p *Mgr) Stop() error {
//	// 定时检查事件总线是否消费完成
//	go checkGBusChannel()
//
//	// 等待GEventChan处理结束
//	p.BusChannelWaitGroup.Wait()
//
//	if p.Bench.Json.Timer.ScanSecondDuration != nil || p.Bench.Json.Timer.ScanMillisecondDuration != nil {
//		p.Timer.Stop()
//		xrlog.GetInstance().Warn("GTimer stop")
//	}
//	if len(p.Bench.Json.Jaeger.Addr) != 0 { // 链路追踪jaeger
//		_ = xrjaeger.GetInstance().Stop()
//		xrlog.GetInstance().Warn("GJaeger stop")
//	}
//	if xretcd.IsEnable() {
//		_ = xretcd.GetInstance().Stop()
//		xrlog.GetInstance().Warn("GEtcd stop")
//	}
//
//	return nil
//}
//
//func checkGBusChannel() {
//	xrlog.GetInstance().Warn("start checkGBusChannel timer")
//
//	idleDuration := 500 * time.Millisecond
//	idleDelay := time.NewTimer(idleDuration)
//	defer func() {
//		idleDelay.Stop()
//	}()
//
//	for {
//		select {
//		case <-idleDelay.C:
//			idleDelay.Reset(idleDuration)
//			GBusChannelCheckChan <- struct{}{}
//			xrlog.GetInstance().Warn("send to GBusChannelCheckChan")
//		}
//	}
//}
