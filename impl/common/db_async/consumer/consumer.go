package consumer

import (
	"context"
	"dawn-server/impl/common/bench"
	"dawn-server/impl/common/db_batch"
	"dawn-server/impl/common/db_retry"
	"dawn-server/impl/common/dk"
	xrlog "dawn-server/impl/xr/lib/log"
	"dawn-server/impl/xr/lib/mongodb"
	"dawn-server/impl/xr/lib/util"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

const ModelAsyncOne uint32 = 1   // 异步执行单条
const ModelAsyncBatch uint32 = 2 // 异步执行合并写

var (
	AsyncConsumerInstance *AsyncConsumer
	asyncConsumerOnce     sync.Once
)

// AsyncConsumer DB异步操作消费者
type AsyncConsumer struct {
	dbMgr       *mongodb.Mgr
	dbRetryMgr  *db_retry.DBRetryMgr
	asyncModel  uint32             // 1:单条执行 2:合并执行
	batchMaxMs  uint32             // 执行周期 单位毫秒
	batchMaxCnt uint32             // 合并写最大数量
	workChan    <-chan interface{} // DB操作通道
	waitGroup   sync.WaitGroup     // 等待异步处理
}

func NewDbAsyncConsumer(dm *mongodb.Mgr, drm *db_retry.DBRetryMgr, bench *bench.DBAsync, workerChan <-chan interface{}) *AsyncConsumer {
	asyncConsumerOnce.Do(func() {
		AsyncConsumerInstance = &AsyncConsumer{
			dbMgr:       dm,
			dbRetryMgr:  drm,
			asyncModel:  bench.Model,
			batchMaxMs:  bench.BulkWriteMs,
			batchMaxCnt: bench.BulkWriteMax,
			workChan:    workerChan,
		}
	})
	return AsyncConsumerInstance
}

// IsAsyncOne 是否异步执行单条
func (p *AsyncConsumer) IsAsyncOne() bool {
	return atomic.LoadUint32(&p.asyncModel) == ModelAsyncOne
}

// IsAsyncBatch 是否异步合并执行
func (p *AsyncConsumer) IsAsyncBatch() bool {
	return atomic.LoadUint32(&p.asyncModel) == ModelAsyncBatch
}

// SwitchMode 切换模式
func (p *AsyncConsumer) SwitchMode(mode uint32) {
	atomic.StoreUint32(&p.asyncModel, mode)
	xrlog.PrintfInfo("DbAsyncConsumer switch mode to:%d", mode)
}

// Start 启动
func (p *AsyncConsumer) Start() {
	p.waitGroup.Add(1)
	go p.watch()
}

// watch 监听消费通道
func (p *AsyncConsumer) watch() {
	defer func() {
		if util.IsRelease() {
			if err := recover(); err != nil {
				xrlog.PrintErr(dk.GoroutinePanic, err, debug.Stack())
			}
		}
		p.waitGroup.Done()
		xrlog.PrintErr(dk.GoroutineDone)
	}()

	funcList := make([]*mongodb.FunctionArg, 0, p.batchMaxCnt)
	dbBatchMgr := db_batch.NewDbBatchMgr(p.dbMgr, p.dbRetryMgr)
	ctx := context.Background()

	// 添加定时器 定时合并 防止数据操作未达到合并上限 长期得不到处理
	idleDuration := time.Duration(p.batchMaxMs) * time.Millisecond
	idleDelay := time.NewTimer(idleDuration)
	defer func() {
		idleDelay.Stop()
	}()

	for {
		select {
		case <-idleDelay.C:
			funcLen := len(funcList)
			if funcLen > 0 {
				p.AsyncBatch(ctx, funcList, dbBatchMgr)
				funcList = funcList[0:0]
				xrlog.PrintfInfo("DbAsyncConsumer exec FunctionArg with timer, funcList len:%d", funcLen)
			}
			idleDelay.Reset(idleDuration)
		case v, ok := <-p.workChan:
			if ok {
				switch t := v.(type) {
				case *mongodb.FunctionArg:
					if p.IsAsyncOne() { // 单条执行
						// 检查是否有未处理的操作（从合并执行切换为单条执行时）
						if len(funcList) > 0 {
							xrlog.PrintfInfo("DbAsyncConsumer exec AsyncOne, funcList len:%d", len(funcList))
							p.AsyncBatch(ctx, funcList, dbBatchMgr)
							funcList = funcList[0:0]
						}
						p.AsyncOne(t)
					} else if p.IsAsyncBatch() { // 合并执行
						funcList = append(funcList, t)
						if len(funcList) >= int(p.batchMaxCnt) {
							xrlog.PrintfInfo("DbAsyncConsumer exec FunctionArg with len:%d >= %d", len(funcList), p.batchMaxCnt)
							p.AsyncBatch(ctx, funcList, dbBatchMgr)
							funcList = funcList[0:0]
						}
					} else {
						xrlog.PrintfErr("DbAsyncConsumer model err:%d", p.asyncModel)
					}
				default:
					xrlog.PrintfErr("Async work channel unknown type:%v", t)
				}
			} else {
				funcLen := len(funcList)
				if funcLen > 0 {
					p.AsyncBatch(ctx, funcList, dbBatchMgr)
					funcList = funcList[0:0]
				}
				xrlog.PrintfInfo("DbAsyncConsumer work channel closed, funcList length:%d", funcLen)
				return
			}
		}
	}
}

// AsyncOne 单条执行
func (p *AsyncConsumer) AsyncOne(funcArg *mongodb.FunctionArg) {
	if p.dbMgr.IsAbnormal() { //容灾
		xrlog.PrintfErr("funcArg:%+v", funcArg)
		if err := p.dbRetryMgr.Save(context.Background(), funcArg); err != nil {
			//设置服务停止
			xrlog.PrintfErr("err:%v", err)
			p.dbRetryMgr.QuitServer()
			return
		}
	} else {
		if _, err := funcArg.Invoke(); err != nil { //处理数据存储失败
			//容灾处理...
			xrlog.PrintfErr("funcArg:%+v", funcArg)
			p.dbMgr.SetAbnormal(true)
			err = p.dbRetryMgr.Save(context.Background(), funcArg)
			if err != nil {
				//设置服务停止
				xrlog.PrintfErr("err:%v", err)
				p.dbRetryMgr.QuitServer()
			}
		}
	}
}

// AsyncBatch 批量bulk write执行
func (p *AsyncConsumer) AsyncBatch(ctx context.Context, funcArgList []*mongodb.FunctionArg, dbBatchMgr *db_batch.DbBatchMgr) {
	if p.dbMgr.IsAbnormal() { //容灾
		for i := range funcArgList {
			funcArg := funcArgList[i]
			if funcArg == nil {
				continue
			}
			xrlog.PrintfErr("funcArg:%+v", funcArg)
			if err := p.dbRetryMgr.Save(context.Background(), funcArg); err != nil {
				//设置服务停止
				xrlog.PrintfErr("err:%v", err)
				p.dbRetryMgr.QuitServer()
			}
		}
	} else {
		startTime := time.Now()
		dbBatchMgr.SetData(funcArgList, p.batchMaxCnt)
		xrlog.PrintfInfo("DbBatch exec data len:%d", len(funcArgList))
		err := dbBatchMgr.Execute(ctx)
		if err != nil {
			xrlog.PrintfErr("DbBatchSave err:%v", err)
		}
		dt := time.Now().Sub(startTime).Milliseconds()
		xrlog.PrintfInfo("DbBatch exec data over, cost %d ms", dt)
	}
}

// Stop 关闭
func (p *AsyncConsumer) Stop() {
	xrlog.PrintInfo("DbAsyncConsumer stopping")

	//等待workChan执行结束
	p.waitGroup.Wait()

	xrlog.PrintInfo("DbAsyncConsumer Stop done")
}
