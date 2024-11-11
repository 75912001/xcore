package gateway

import (
	"context"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	xruntime "xcore/lib/runtime"
	xservice "xcore/lib/service"
)

var gservice *Service

type Service struct {
	*xservice.Service
	LoginServiceMgr *LoginServiceMgr
}

func NewService(defaultService *xservice.Service) *Service {
	gservice = &Service{
		Service:         defaultService,
		LoginServiceMgr: NewLoginServiceMgr(),
	}
	return gservice
}

func (p *Service) Start(ctx context.Context) (err error) {
	if err = p.Service.Start(ctx, p, logCallBackFunc, EtcdKeyValue); err != nil {
		return errors.WithMessagef(err, xruntime.Location())
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

	// 退出服务
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
EXIT:
	for {
		select {
		case <-p.QuitChan:
			p.Log.Warn("service will shutdown in a few seconds")
			_ = p.PreStop()
			_ = p.Stop()
			break EXIT // 退出循环
		case s := <-sigChan:
			p.Log.Warnf("service got signal: %s, shutting down...", s)
			close(p.QuitChan)
		}
	}
	return nil
	// 数据统计
	//if err = server_del.GST.Start(sub_bench.GMgr.Json.STKafka.Addrs, st_mgr.GameID,
	//	sub_bench.GMgr.Json.StServerID, sub_bench.GMgr.Json.StServerType); err != nil {
	//	return errors.WithMessagef(err, "GSTKafka Start %v %v",
	//		sub_bench.GMgr.Json.STKafka.Addrs, xrutil.GetCodeLocation(1).String())
	//}

	// DB容灾处理
	//server_del.GDBRetryMgr.Start(server_del.GMgr.ZoneID, server_del.GMgr.ServiceID,
	//	xrredis.GetInstance(), xrmongodb.GetInstance(), server_del.GQuitChan)
	////server.GDBRetryMgr.Trigger()
	//{ //数据库容灾, 从redis中获取需要容灾的数据,更新至mongodb中.
	//	err = server_del.GDBRetryMgr.LoadAndExecute(context.Background())
	//	if err != nil {
	//		return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	//	}
	//	err = server_del.GDBRetryMgr.RedisDel(context.Background())
	//	if err != nil {
	//		return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	//	}
	//}

	// 加载策划配置表
	//if err = config.Load(sub_bench.GMgr.Json.Config.Path); err != nil {
	//	return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	//}

	// 启动 mongodb async模式
	//if sub_bench.GMgr.Json.DBAsync.ChanCnt > 0 {
	//	db_async.GDBAsyncChan = make(chan interface{}, sub_bench.GMgr.Json.DBAsync.ChanCnt)
	//
	//	dbAsyncProducer := producer.NewDbAsyncProducer(db_async.GDBAsyncChan)
	//	dbAsyncProducer.Start()
	//	dbAsyncConsumer := consumer.NewDbAsyncConsumer(xrmongodb.GetInstance(), &server_del.GDBRetryMgr, &sub_bench.GMgr.Json.DBAsync, db_async.GDBAsyncChan)
	//	dbAsyncConsumer.Start()
	//	defer func() {
	//		// 先关闭生产者
	//		dbAsyncProducer.Stop()
	//		xrlog.GetInstance().Warn("dbAsyncProducer stop")
	//		dbAsyncConsumer.Stop()
	//		xrlog.GetInstance().Warn("dbAsyncConsumer stop")
	//	}()
	//}

	// world服定时器
	//serverTimer := new(handler.ServerTimer)
	//serverTimer.Start()
	//defer func() {
	//	serverTimer.Stop()
	//	xrlog.GetInstance().Warn("serverTimer stop")
	//}()

	//err = server_del.GMgr.PostInit(context.TODO(),
	//	server_del.NewOptions().
	//		SetEtcdHandler(handler.OnEventEtcd).
	//		SetEtcdWatchServicePrefix(fmt.Sprintf("/%v/%v/", common.ProjectName, common.EtcdWatchMsgTypeService)).
	//		SetEtcdWatchCommandPrefix(fmt.Sprintf("/%v/%v/%v/%v/",
	//			common.ProjectName, common.EtcdWatchMsgTypeCommand,
	//			server_del.GMgr.ZoneID,
	//			server_del.GMgr.ServiceName)),
	//	//SetEtcdWatchGMPrefix(fmt.Sprintf("/%v/%v/%v/%v/%v",
	//	//	common.ProjectName, common.EtcdWatchMsgTypeGM,
	//	//	server.GMgr.ZoneID,
	//	//	server.GMgr.ServiceName,
	//	//	server.GMgr.ServiceID)),
	//) // 需要监听除本zone以外的world服务信息
	//if err != nil {
	//	return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	//}

	// 信息
	//state.Start()

	//runtime.GC()

	//// 退出服务
	//sigChan := make(chan os.Signal, 1)
	//signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	//
	//select {
	//case <-server_del.GQuitChan:
	//	xrlog.GetInstance().Warn("GServer will shutdown in a few seconds")
	//case s := <-sigChan:
	//	xrlog.GetInstance().Warnf("GServer got signal: %s, shutting down...", s)
	//}
	//
	//world.PreStop()
	//_ = server_del.GMgr.Stop()

	//return err
}

//func (p *Service) Stop() (err error) {
//	server_del.GDBRetryMgr.Stop()
//	xrlog.GetInstance().Warn("GDBRetryMgr stop")
//
//	if xrredis.IsEnable() {
//		_ = xrredis.GetInstance().Disconnect()
//		xrlog.GetInstance().Warn("GZoneRedis stop")
//	}
//
//	if xrmongodb.GetInstance().GetClient() != nil {
//		err = xrmongodb.GetInstance().Disconnect(context.Background())
//		if err != nil {
//			xrlog.GetInstance().Fatal(xrmongodb.ErrorKeyDisconnectFailure, err)
//		}
//		xrlog.GetInstance().Warn("GZoneMongoDB stop")
//	}
//
//	server_del.GST.Stop()
//	xrlog.GetInstance().Warn("GST stop")
//
//	if err := server_del.GZoneNats.Unsubscribe(server_del.GZoneNats.Subscription); err != nil {
//		xrlog.GetInstance().Errorf("err:%v", err)
//	}
//	server_del.GZoneNats.Close()
//	xrlog.GetInstance().Warn("GZoneNats stop")
//
//	if err := server_del.GGlobalNats.Unsubscribe(server_del.GGlobalNats.Subscription); err != nil {
//		xrlog.GetInstance().Errorf("err:%v", err)
//	}
//	server_del.GGlobalNats.Close()
//	xrlog.GetInstance().Warn("GGlobalNats stop")
//
//	xrlog.GetInstance().Warn("verify chat stop")
//	_ = verify_chat.GetInstance().Stop()
//
//	xrlog.PrintErr("GLog stop")
//	_ = xrlog.GetInstance().Stop()
//	return nil
//}

func (p *Service) PreStop() error {
	_ = p.Etcd.Stop()
	return nil
}
func (p *Service) Stop() (err error) {
	_ = p.Service.Stop()
	return nil
}
