package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"math/rand"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
	xbench "xcore/lib/bench"
	xcallback "xcore/lib/callback"
	xconstants "xcore/lib/constants"
	xerror "xcore/lib/error"
	xetcd "xcore/lib/etcd"
	xlog "xcore/lib/log"
	xnettcp "xcore/lib/net/tcp"
	xpprof "xcore/lib/pprof"
	xruntime "xcore/lib/runtime"
	xtime "xcore/lib/time"
	xtimer "xcore/lib/timer"
	xutil "xcore/lib/util"
)

type DefaultService struct {
	BenchMgr xbench.Mgr
	BenchSub xbench.IBenchSub

	GroupID        uint32 // 组ID
	Name           string // 名称
	ID             uint32 // ID
	ExecutablePath string // 执行程序路径 // 程序所在路径(如为link,则为link所在的路径)

	Log     xlog.ILog
	TimeMgr *xtime.Mgr
	Timer   xtimer.ITimer

	Etcd    xetcd.IEtcd
	EtcdKey string // etcd key

	BusChannel          chan interface{} // 总线
	BusChannelWaitGroup sync.WaitGroup   // 总线等待

	QuitChan chan struct{} // 退出信号, 用于关闭服务
}

func NewDefaultService() *DefaultService {
	val := &DefaultService{
		TimeMgr:  xtime.NewMgr(),
		QuitChan: make(chan struct{}),
	}
	// 程序所在路径(如为link,则为link所在的路径)
	if executablePath, err := xruntime.GetExecutablePath(); err != nil {
		xlog.PrintErr(err, xruntime.Location())
		return nil
	} else {
		val.ExecutablePath = executablePath
	}
	return val
}

//func (p *DefaultService) Start(ctx context.Context) (err error) {
//	return xerror.NotImplemented
//}
//
//func (p *DefaultService) PreStop() error {
//	return xerror.NotImplemented
//}
//
//func (p *DefaultService) Stop() (err error) {
//	return xerror.NotImplemented
//}

func (p *DefaultService) Start(ctx context.Context, handler xnettcp.IHandler, logCallbackFunc xlog.CallBackFunc, etcdCallbackFun xetcd.CallbackFun) (err error) {
	rand.Seed(time.Now().UnixNano())
	p.TimeMgr.Update()
	// 小端
	if !xutil.IsLittleEndian() {
		return errors.Errorf("system is bigEndian! %v", xruntime.Location())
	}
	// 开启UUID随机
	uuid.EnableRandPool()
	// 服务配置文件
	benchPath := path.Join(p.ExecutablePath, fmt.Sprintf("%v.%v.%v.%v",
		p.GroupID, p.Name, p.ID, xconstants.ServiceConfigFile))
	content, err := os.ReadFile(benchPath)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	benchJson := string(content)
	// 加载服务配置文件-root部分
	err = p.BenchMgr.RootJson.Parse(benchJson)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	// 加载服务配置文件-公共部分
	err = p.BenchMgr.Json.Parse(benchJson)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if false { // 从etcd获取配置项 todo menglc
		client, err := clientv3.New(
			clientv3.Config{
				Endpoints:   p.BenchMgr.RootJson.Etcd.Addrs,
				DialTimeout: 5 * time.Second, // todo menglc 确定用途?
			},
		)
		if err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
		kv := clientv3.NewKV(client)
		key := fmt.Sprintf("/%v/%v/%v/%v/%v",
			*p.BenchMgr.Json.Base.ProjectName, xconstants.EtcdWatchMsgTypeServiceBench, p.GroupID, p.Name, p.ID)
		getResponse, err := kv.Get(ctx, key, clientv3.WithPrefix())
		if err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
		if len(getResponse.Kvs) != 1 {
			return errors.WithMessagef(xerror.Config, "%v %v %v", key, getResponse.Kvs, xruntime.Location())
		}
		benchJson = string(getResponse.Kvs[0].Value)
		xlog.PrintfInfo(benchJson)
	}
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	switch *p.BenchMgr.Json.Base.RunMode {
	case 0:
		xruntime.SetRunMode(xruntime.RunModeRelease)
	case 1:
		xruntime.SetRunMode(xruntime.RunModeDebug)
	default:
		return errors.Errorf("runMode err:%v %v", *p.BenchMgr.Json.Base.RunMode, xruntime.Location())
	}
	// GoMaxProcess
	previous := runtime.GOMAXPROCS(*p.BenchMgr.Json.Base.GoMaxProcess)
	xlog.PrintfInfo("go max process new:%v, previous setting:%v",
		*p.BenchMgr.Json.Base.GoMaxProcess, previous)
	// 日志
	p.Log, err = xlog.NewMgr(xlog.NewOptions().
		WithLevel(*p.BenchMgr.Json.Base.LogLevel).
		WithAbsPath(*p.BenchMgr.Json.Base.LogAbsPath).
		WithNamePrefix(fmt.Sprintf("%v.%v.%v", p.GroupID, p.Name, p.ID)).
		WithLevelCallBack(logCallbackFunc, xlog.LevelFatal, xlog.LevelError, xlog.LevelWarn),
	)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	// 加载服务配置文件-子项部分
	if p.BenchSub != nil {
		err = p.BenchSub.Unmarshal(benchJson)
		if err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
	}
	// eventChan
	p.BusChannel = make(chan interface{}, *p.BenchMgr.Json.Base.BusChannelCapacity)
	go func() {
		defer func() {
			p.BusChannelWaitGroup.Done()
			// 主事件 channel 报错 不 recover
			p.Log.Infof(xconstants.GoroutineDone)
		}()
		p.BusChannelWaitGroup.Add(1)
		_ = p.Handle()
	}()
	// 是否开启http采集分析
	if p.BenchMgr.Json.Base.PprofHttpPort != nil {
		xpprof.StartHTTPprof(fmt.Sprintf("0.0.0.0:%d", *p.BenchMgr.Json.Base.PprofHttpPort))
	}
	// 全局定时器
	if p.BenchMgr.Json.Timer.ScanSecondDuration != nil || p.BenchMgr.Json.Timer.ScanMillisecondDuration != nil {
		p.Timer = xtimer.NewTimer()
		err = p.Timer.Start(ctx,
			xtimer.NewOptions().
				WithScanSecondDuration(*p.BenchMgr.Json.Timer.ScanSecondDuration).
				WithScanMillisecondDuration(*p.BenchMgr.Json.Timer.ScanMillisecondDuration).
				WithOutgoingTimerOutChan(p.BusChannel),
		)
		if err != nil {
			return errors.Errorf("timer Start err:%v %v", err, xruntime.Location())
		}
	}
	// etcd
	p.EtcdKey = xetcd.GenKey(*p.BenchMgr.Json.Base.ProjectName, xconstants.EtcdWatchMsgTypeService, p.GroupID, p.Name, p.ID)
	cmdMin, err := xutil.HexStringToUint32(*p.BenchMgr.Json.Base.CmdMin)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	cmdMax, err := xutil.HexStringToUint32(*p.BenchMgr.Json.Base.CmdMax)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	defaultEtcd := xetcd.NewDefaultEtcd(
		xetcd.NewOptions().
			WithAddrs(p.BenchMgr.RootJson.Etcd.Addrs).
			WithTTL(*p.BenchMgr.RootJson.Etcd.TTL).
			WithWatchKeyPrefix(xetcd.GenPrefixKey(*p.BenchMgr.Json.Base.ProjectName)).
			WithKey(p.EtcdKey).
			WithValue(
				&xetcd.ValueJson{
					ServiceNet:    &p.BenchMgr.Json.ServiceNet,
					Version:       *p.BenchMgr.Json.Base.Version,
					AvailableLoad: *p.BenchMgr.Json.Base.AvailableLoad,
					SecondOffset:  0,
					CmdMin:        cmdMin,
					CmdMax:        cmdMax,
				},
			).
			WithEventChan(p.BusChannel),
	)
	defaultEtcd.CallbackFun = etcdCallbackFun
	p.Etcd = defaultEtcd
	if err = p.Etcd.Start(ctx); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	// 续租
	err = defaultEtcd.KeepAlive(ctx)
	if err != nil {
		return errors.WithMessagef(err, xruntime.Location())
	}
	// etcd-定时上报
	p.Timer.AddSecond(xcallback.NewCallBack(EtcdReportFunction, p), p.TimeMgr.ShadowTimestamp()+xconstants.EtcdReportIntervalSecondDefault)

	// 网络服务
	if len(*p.BenchMgr.Json.ServiceNet.Addr) != 0 {
		switch *p.BenchMgr.Json.ServiceNet.Type {
		case "tcp": // 启动 TCP 服务
			if err = xnettcp.NewServer(handler).Start(ctx,
				xnettcp.NewServerOptions().
					SetListenAddress(*p.BenchMgr.Json.ServiceNet.Addr).
					SetEventChan(p.BusChannel).
					SetSendChanCapacity(*p.BenchMgr.Json.Base.SendChanCapacity),
			); err != nil {
				return errors.WithMessage(err, xruntime.Location())
			}
		case "udp":
			return errors.WithMessage(xerror.NotImplemented, xruntime.Location())
		default:
			return errors.WithMessage(xerror.NotImplemented, xruntime.Location())
		}
	}
	StateTimerPrint(p.Timer, p.Log)
	return nil
}

func (p *DefaultService) Stop() (err error) {
	return nil
}
