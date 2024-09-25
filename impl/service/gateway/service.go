package gateway

import (
	"runtime"
	commonservice "xcore/impl/common/service"
)

type Service struct {
	*commonservice.DefaultService
}

func NewService(defaultService *commonservice.DefaultService) *Service {
	return &Service{
		DefaultService: defaultService,
	}
}

func (p *Service) Start() (err error) {
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

	runtime.GC()

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
	//world.PreShutdown()
	//_ = server_del.GMgr.Stop()

	return err
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
