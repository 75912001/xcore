package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
	xbench "xcore/lib/bench"
	xconstants "xcore/lib/constants"
	xerror "xcore/lib/error"
	xlog "xcore/lib/log"
	xpprof "xcore/lib/pprof"
	xruntime "xcore/lib/runtime"
	xtime "xcore/lib/time"
	xtimer "xcore/lib/timer"
	xutil "xcore/lib/util"
)

type DefaultService struct {
	BenchMgr xbench.Mgr
	BenchSub xbench.IBenchSub

	GroupID uint32 // 组ID
	Name    string // 名称
	ID      uint32 // ID

	TimeMgr  xtime.Mgr
	TimerMgr xtimer.Mgr
	Opts     *options

	EtcdMgr etcdMgr

	BusChannel          chan interface{}
	BusChannelWaitGroup sync.WaitGroup
}

func NewDefaultService() *DefaultService {
	return &DefaultService{}
}

func (p *DefaultService) WithBenchSub(benchSub xbench.IBenchSub) *DefaultService {
	p.BenchSub = benchSub
	return p
}

func (p *DefaultService) Start() (err error) {
	return xerror.NotImplemented
}

func (p *DefaultService) Stop() (err error) {
	return xerror.NotImplemented
}

func (p *DefaultService) PreStart(ctx context.Context, opts ...*options) error {
	rand.Seed(time.Now().UnixNano())
	p.TimeMgr.Update()
	// 小端
	if !xutil.IsLittleEndian() {
		return errors.Errorf("system is bigEndian! %v", xruntime.Location())
	}
	// 开启UUID随机
	uuid.EnableRandPool()
	// 处理服务参数选项
	p.Opts = mergeOptions(opts...)
	err := configure(p.Opts)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	content, err := os.ReadFile(*p.Opts.BenchPath)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	benchJson := string(content)
	// 加载服务配置文件-root部分
	err = p.BenchMgr.RootJson.Parse(benchJson)
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
			p.BenchMgr.Json.Base.ProjectName, xconstants.EtcdWatchMsgTypeServiceBench, p.GroupID, p.Name, p.ID)
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
	// 加载服务配置文件-公共部分
	err = p.BenchMgr.Json.Parse(benchJson)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	// GoMaxProcess
	previous := runtime.GOMAXPROCS(*p.BenchMgr.Json.Base.GoMaxProcess)
	xlog.PrintfInfo("go max process new:%v, previous setting:%v",
		*p.BenchMgr.Json.Base.GoMaxProcess, previous)
	// 日志
	var log xlog.ILog
	log, err = xlog.NewMgr(xlog.NewOptions().
		WithLevel(*p.BenchMgr.Json.Base.LogLevel).
		WithAbsPath(*p.BenchMgr.Json.Base.LogAbsPath).
		WithNamePrefix(fmt.Sprintf("%v.%v.%v", p.GroupID, p.Name, p.ID)).
		WithLevelCallBack(logCallBackFunc, xlog.LevelFatal, xlog.LevelError, xlog.LevelWarn),
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
			log.Infof(xconstants.GoroutineDone)
		}()
		p.BusChannelWaitGroup.Add(1)
		p.HandlerBus()
	}()
	// 是否开启http采集分析
	if p.BenchMgr.Json.Base.PprofHttpPort != nil {
		xpprof.StartHTTPprof(fmt.Sprintf("0.0.0.0:%d", *p.BenchMgr.Json.Base.PprofHttpPort))
	}
	// 全局定时器
	if p.BenchMgr.Json.Timer.ScanSecondDuration != nil || p.BenchMgr.Json.Timer.ScanMillisecondDuration != nil {
		err = p.TimerMgr.Start(ctx,
			xtimer.NewOptions().
				WithScanSecondDuration(*p.BenchMgr.Json.Timer.ScanSecondDuration).
				WithScanMillisecondDuration(*p.BenchMgr.Json.Timer.ScanMillisecondDuration).
				WithOutgoingTimerOutChan(p.BusChannel),
		)
		if err != nil {
			return errors.Errorf("timer Start err:%v %v", err, xruntime.Location())
		}
	}

	// todo menglc 开启 jaeger

	// todo menglc 开启 MongoDB

	// todo menglc 开启 Redis

	// todo menglc 开启 NATS

	// 链路追踪jaeger
	//if len(p.BenchMgr.Json.Jaeger.Addr) != 0 {
	//
	//}
	//if len(p.Bench.Json.Jaeger.Addr) != 0 {
	//	err = xrjaeger.GetInstance().Start(context.TODO(),
	//		xrjaeger.NewOptions().
	//			SetEndpoint(p.Bench.Json.Jaeger.Addr).
	//			SetServiceName(fmt.Sprintf("%v.%v.%v", p.ZoneID, p.ServiceName, p.ServiceID)).
	//			SetServiceVersion(fmt.Sprintf("%v", p.Bench.Json.Base.Version)),
	//	)
	//	if err != nil {
	//		return errors.Errorf("Jaeger connect %+v err:%v %v", p.Bench.Json.Jaeger.Addr, err, xrutil.GetCodeLocation(1).String())
	//	}
	//}
	runtime.GC()
	return nil
}
