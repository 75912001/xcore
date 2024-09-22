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
	"xcore/impl/common"
	gatewayhandler "xcore/impl/service/gateway/handler"
	xbench "xcore/lib/bench"
	xconstants "xcore/lib/constants"
	xerror "xcore/lib/error"
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
	common.GLog, err = xlog.NewMgr(xlog.NewOptions().
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
			common.GLog.Infof(xconstants.GoroutineDone)
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

	if len(*p.BenchMgr.Json.ServiceNet.Addr) != 0 { // 网络服务
		switch *p.BenchMgr.Json.ServiceNet.Type {
		case "tcp": // 启动 TCP 服务
			if err := xnettcp.GetServer().Start(ctx, xnettcp.NewServerOptions().
				SetListenAddress(*p.BenchMgr.Json.ServiceNet.Addr).
				SetEventChan(p.BusChannel).
				SetSendChanCapacity(*p.BenchMgr.Json.Base.SendChanCapacity).
				SetPacket(xnettcp.NewDefaultPacket()).
				SetHandler(gatewayhandler.NewServer())); err != nil {
				return errors.WithMessage(err, xruntime.Location())
			}
		case "udp":
			return errors.WithMessage(xerror.NotImplemented, xruntime.Location())
		default:
			return errors.WithMessage(xerror.NotImplemented, xruntime.Location())
		}
	}

	// todo menglc 开启 jaeger
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
	// todo menglc 开启 MongoDB
	// 连接mongodb
	//subZoneMongodb := sub_bench.GMgr.Json.ZoneMongoDB
	//if err := xrmongodb.GetInstance().Connect(context.TODO(),
	//	xrmongodb.NewOptions().
	//		SetAddrs(subZoneMongodb.Addrs).
	//		SetUserName(subZoneMongodb.User).
	//		SetPW(subZoneMongodb.Pwd).
	//		SetDBName(subZoneMongodb.DBName).
	//		SetMaxPoolSize(subZoneMongodb.MaxPoolSize).
	//		SetMinPoolSize(subZoneMongodb.MinPoolSize).
	//		SetTimeoutDuration(subZoneMongodb.TimeoutDuration).
	//		SetMaxConnIdleTime(subZoneMongodb.MaxConnIdleTime).
	//		SetMaxConnecting(subZoneMongodb.MaxConnecting),
	//); err != nil {
	//	return errors.WithMessagef(err, "%v %v", sub_bench.GMgr.Json.ZoneMongoDB, xrutil.GetCodeLocation(1).String())
	//} else {
	//	xrmongodb.GetInstance().SwitchedDatabase(mdb.GenDatabaseName(server_del.GMgr.ZoneID))
	//
	//	mdb_sys.Collection = xrmongodb.GetInstance().SwitchedCollection(mdb_sys.CollectionName)
	//	mdb_user.Collection = xrmongodb.GetInstance().SwitchedCollection(mdb_user.CollectionName)
	//	mdb_inventory.Collection = xrmongodb.GetInstance().SwitchedCollection(mdb_inventory.CollectionName)
	//
	//	//加载系统 event 数据
	//	var record common_proto.DBInventory
	//	opts := options.FindOne().SetProjection(bson.M{"_id": 0, "uid": 0})
	//	if _, err = mdb_inventory.FindOne(context.Background(), xrmongodb.GetInstance(), 0, &record, opts); err != nil {
	//		return errors.WithMessagef(err, "%v mongodb col_event FindOne %v", dk.event, xrutil.GetCodeLocation(1).String())
	//	} else {
	//		xrlog.GetInstance().Tracef("%v %+v", dk.event, &record)
	//		world.GEventMgr.LoadFromDB(0, record.Events)
	//	}
	//}
	// todo menglc 开启 Redis
	// 连接redis
	//if err = xrredis.GetInstance().Connect(context.Background(),
	//	xrredis.NewOptions().
	//		SetAddrs(sub_bench.GMgr.Json.ZoneRedis.Addrs).
	//		SetPW(sub_bench.GMgr.Json.ZoneRedis.Password),
	//); err != nil {
	//	return errors.WithMessagef(err, "redis connect addrs:%v %v",
	//		sub_bench.GMgr.Json.ZoneRedis.Addrs, xrutil.GetCodeLocation(1).String())
	//}
	// redis订阅
	//go func() {
	//	defer func() {
	//		if xrutil.IsRelease() {
	//			if err := recover(); err != nil {
	//				xrlog.GetInstance().Fatal(dk.GoroutinePanic, err, debug.Stack())
	//			}
	//		}
	//		xrlog.GetInstance().Fatal(dk.GoroutineDone)
	//	}()
	//	sub := xrredis.GetInstance().Client.Subscribe(context.Background(),
	//		redis_pub_sub.GenChannelNotice(server_del.GMgr.ZoneID))
	//	// sub.Channel() 返回go channel，可以循环读取redis服务器发过来的消息
	//	for message := range sub.Channel() {
	//		xrlog.GetInstance().Debugf("redis Subscribe %v", dk.Notice)
	//		server_del.GMgr.BusChannel <- message
	//	}
	//}()
	// todo menglc 开启 NATS
	//{
	//	// 连接 global NATS
	//	var urls string
	//	for k := range sub_bench.GMgr.Json.GlobalNATS.Addrs {
	//		urls += sub_bench.GMgr.Json.GlobalNATS.Addrs[k] + ","
	//	}
	//	urls = strings.TrimRight(urls, ",")
	//
	//	var option nats.Option
	//	if sub_bench.GMgr.Json.GlobalNATS.User != nil && sub_bench.GMgr.Json.GlobalNATS.Password != nil {
	//		option = nats.UserInfo(*sub_bench.GMgr.Json.GlobalNATS.User, *sub_bench.GMgr.Json.GlobalNATS.Password)
	//	}
	//
	//	if err = server_del.GGlobalNats.Connect(urls, option); err != nil {
	//		return errors.WithMessagef(err, "nats connect addrs:%v urls:%v %v",
	//			sub_bench.GMgr.Json.GlobalNATS.Addrs, urls, xrutil.GetCodeLocation(1).String())
	//	}
	//	server_del.GGlobalNats.Sub = mq_nats.GenGlobalSub(server_del.GMgr.ZoneID, server_del.GMgr.ServiceName,
	//		server_del.GMgr.ServiceID)
	//	server_del.GGlobalNats.Subscription, err = server_del.GGlobalNats.Subscribe(server_del.GGlobalNats.Sub, &world.GMQPbFunMgr, server_del.GMgr.BusChannel)
	//	if err != nil {
	//		return errors.WithMessagef(err, "GGlobalNats Subscribe sub:%v %v",
	//			server_del.GGlobalNats.Sub, xrutil.GetCodeLocation(1).String())
	//	}
	//}
	{
		//// 连接zone NATS
		//var urls string
		//for k := range sub_bench.GMgr.Json.ZoneNATS.Addrs {
		//	urls += sub_bench.GMgr.Json.ZoneNATS.Addrs[k] + ","
		//}
		//urls = strings.TrimRight(urls, ",")
		//
		//var option nats.Option
		//if sub_bench.GMgr.Json.ZoneNATS.User != nil && sub_bench.GMgr.Json.ZoneNATS.Password != nil {
		//	option = nats.UserInfo(*sub_bench.GMgr.Json.ZoneNATS.User, *sub_bench.GMgr.Json.ZoneNATS.Password)
		//}
		//
		//if err = server_del.GZoneNats.Connect(urls, option); err != nil {
		//	return errors.WithMessagef(err, "nats connect addrs:%v urls:%v %v",
		//		sub_bench.GMgr.Json.ZoneNATS.Addrs, urls, xrutil.GetCodeLocation(1).String())
		//}
		//server_del.GZoneNats.Sub = mq_nats.GenZoneSub(server_del.GMgr.ZoneID, server_del.GMgr.ServiceName,
		//	server_del.GMgr.ServiceID)
		//server_del.GZoneNats.Subscription, err = server_del.GZoneNats.Subscribe(server_del.GZoneNats.Sub, &world.GMQPbFunMgr, server_del.GMgr.BusChannel)
		//if err != nil {
		//	return errors.WithMessagef(err, "GZoneNats Subscribe sub:%v err:%v",
		//		server_del.GZoneNats.Sub, xrutil.GetCodeLocation(1).String())
		//}
	}

	runtime.GC()
	return nil
}
