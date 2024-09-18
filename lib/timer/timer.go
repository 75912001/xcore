// Package timer 定时器
// 优先级: 到期时间,加入顺序
package timer

import (
	"container/list"
	"context"
	"github.com/pkg/errors"
	"runtime/debug"
	"sync"
	"time"
	"xcore/lib/constants"
	"xcore/lib/log"
	"xcore/lib/runtime"
)

// Mgr 定时器管理器
type Mgr struct {
	opts            *options
	secondSlice     [cycleSize]list.List // 时间轮-数组. 秒,数据
	millisecondList list.List            // 毫秒,数据
	cancelFunc      context.CancelFunc
	waitGroup       sync.WaitGroup   // Stop 等待信号
	milliSecondChan chan interface{} // 毫秒, channel
	secondChan      chan interface{} // 秒, channel
}

func NewMgr() *Mgr {
	return &Mgr{}
}

// OnFun 回调定时器函数(使用协程回调)
type OnFun func(arg interface{})

// 每秒更新
func (p *Mgr) funcSecond(ctx context.Context) {
	defer func() {
		if runtime.IsRelease() {
			if err := recover(); err != nil {
				log.PrintErr(constants.GoroutinePanic, err, string(debug.Stack()))
			}
		}
		p.waitGroup.Done()
		log.PrintInfo(constants.GoroutineDone)
	}()
	idleDelay := time.NewTimer(*p.opts.scanSecondDuration)
	defer func() {
		idleDelay.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			log.PrintInfo(constants.GoroutineDone)
			return
		case v := <-p.secondChan:
			s := v.(*Second)
			p.pushBackCycle(s, searchCycleIdxIteration(s.expire), true)
		case <-idleDelay.C:
			idleDelay.Reset(*p.opts.scanSecondDuration)
			p.scanSecond(ShadowTimestamp())
		}
	}
}

// 每 Millisecond 个毫秒更新
func (p *Mgr) funcMillisecond(ctx context.Context) {
	defer func() {
		if runtime.IsRelease() {
			if err := recover(); err != nil {
				log.PrintErr(constants.GoroutinePanic, err, string(debug.Stack()))
			}
		}
		p.waitGroup.Done()
		log.PrintInfo(constants.GoroutineDone)
	}()
	scanMillisecondDuration := *p.opts.scanMillisecondDuration
	scanMillisecond := scanMillisecondDuration / time.Millisecond
	idleDelay := time.NewTimer(scanMillisecondDuration)
	defer func() {
		idleDelay.Stop()
	}()
	nextMillisecond := time.Duration(time.Now().UnixMilli()) + scanMillisecond

	for {
		select {
		case <-ctx.Done():
			log.PrintInfo(constants.GoroutineDone)
			return
		case v := <-p.milliSecondChan:
			p.millisecondList.PushBack(v)
			moveLastElementToProperPosition(&p.millisecondList)
		case <-idleDelay.C:
			nowMillisecond := time.Now().UnixMilli()
			reset := scanMillisecondDuration - (time.Duration(nowMillisecond)-nextMillisecond)*time.Millisecond
			idleDelay.Reset(reset)

			nextMillisecond += scanMillisecond
			p.scanMillisecond(nowMillisecond)
		}
	}
}

// 移动最后一个元素到合适的位置,移动到大于他的元素的前面[实现按照时间排序,加入顺序排序]
// e.g.: 1,2,2,3,4,4,3 => 1,2,2,3,3,4,4 [将最后一个元素移动到4的前面]
func moveLastElementToProperPosition(l *list.List) {
	lastElement := l.Back() // 获取最后一个元素
	target := lastElement.Value.(*Millisecond)
	var element *list.Element
	for element = lastElement.Prev(); element != nil; element = element.Prev() {
		current := element.Value.(*Millisecond)
		if current.expire <= target.expire {
			l.MoveAfter(lastElement, element)
			return
		}
	}
	if element == nil {
		// 如果没有找到比目标小或等于的元素，将目标元素移动到列表的前面
		l.MoveToFront(lastElement)
	}
}

// Start
// [NOTE] 处理定时器相关数据,必须与该 timeoutChan 线性处理.如:在同一个 goroutine select 中处理数据
func (p *Mgr) Start(ctx context.Context, opts ...*options) error {
	p.opts = &options{}
	p.opts = p.opts.merge(opts...)
	if err := p.opts.configure(); err != nil {
		return errors.WithMessage(err, runtime.Location())
	}

	ctxWithCancel, cancelFunc := context.WithCancel(ctx)
	p.cancelFunc = cancelFunc

	if p.opts.scanSecondDuration != nil {
		p.secondChan = make(chan interface{}, 100)
		p.waitGroup.Add(1)

		go p.funcSecond(ctxWithCancel)
	}
	if p.opts.scanMillisecondDuration != nil {
		p.milliSecondChan = make(chan interface{}, 100)
		p.waitGroup.Add(1)

		go p.funcMillisecond(ctxWithCancel)
	}
	return nil
}

// Stop 停止服务
func (p *Mgr) Stop() {
	if p.cancelFunc != nil {
		p.cancelFunc()
		// 等待 Second, milliSecond goroutine退出.
		p.waitGroup.Wait()
		p.cancelFunc = nil
	}
}

// AddMillisecond 添加毫秒级定时器
//
//	参数:
//		cb: 回调函数
//		Arg: 回调参数
//		expireMillisecond: 过期毫秒数
//	返回值:
//		毫秒定时器
func (p *Mgr) AddMillisecond(cb OnFun, arg interface{}, expireMillisecond int64) *Millisecond {
	t := &Millisecond{
		Arg:      arg,
		Function: cb,
		expire:   expireMillisecond,
		valid:    true,
	}
	p.milliSecondChan <- t
	return t
}

// 扫描毫秒级定时器
//
//	参数:
//		ms: 到期毫秒数
func (p *Mgr) scanMillisecond(ms int64) {
	var next *list.Element
	for e := p.millisecondList.Front(); e != nil; e = next {
		timerMillisecond := e.Value.(*Millisecond)
		if !timerMillisecond.isValid() {
			next = e.Next()
			p.millisecondList.Remove(e)
			continue
		}
		if timerMillisecond.expire <= ms {
			p.opts.outgoingTimeoutChan <- timerMillisecond
			next = e.Next()
			p.millisecondList.Remove(e)
			continue
		}
		break
	}
}

// AddSecond 添加秒级定时器
//
//	参数:
//		cb: 回调函数
//		Arg: 回调参数
//		expire: 过期秒数
//	返回值:
//		秒定时器
func (p *Mgr) AddSecond(cb OnFun, arg interface{}, expire int64) *Second {
	t := &Second{
		Millisecond{
			Arg:      arg,
			Function: cb,
			expire:   expire,
			valid:    true,
		},
	}
	p.secondChan <- t
	return t
}

// 将秒级定时器,添加到轮转IDX的末尾.
//
//		参数:
//			timerSecond: 秒定时器
//			cycleIdx: 轮序号
//	     needMove: 是否需要移动到合适的位置
func (p *Mgr) pushBackCycle(timerSecond *Second, cycleIdx int, needMove bool) {
	p.secondSlice[cycleIdx].PushBack(timerSecond)
	if needMove {
		moveLastElementToProperPositionSecond(&p.secondSlice[cycleIdx])
	}
}

// 移动最后一个元素到合适的位置,移动到大于他的元素的前面[实现按照时间排序,加入顺序排序]
// e.g.: 1,2,2,3,4,4,3 => 1,2,2,3,3,4,4 [将最后一个元素移动到4的前面]
func moveLastElementToProperPositionSecond(l *list.List) {
	lastElement := l.Back() // 获取最后一个元素
	target := lastElement.Value.(*Second)
	var element *list.Element
	for element = lastElement.Prev(); element != nil; element = element.Prev() {
		current := element.Value.(*Second)
		if current.expire <= target.expire {
			l.MoveAfter(lastElement, element)
			return
		}
	}
	if element == nil {
		// 如果没有找到比目标小或等于的元素，将目标元素移动到列表的前面
		l.MoveToFront(lastElement)
	}
}

// 扫描秒级定时器
//
//	timestamp: 到期时间戳
func (p *Mgr) scanSecond(timestamp int64) {
	var next *list.Element
	cycle0 := &p.secondSlice[0]
	for e := cycle0.Front(); nil != e; e = next {
		t := e.Value.(*Second)
		if !t.isValid() {
			next = e.Next()
			cycle0.Remove(e)
			continue
		}
		if t.expire <= timestamp {
			p.opts.outgoingTimeoutChan <- t
			next = e.Next()
			cycle0.Remove(e)
			continue
		}
		break
	}
	// 更新时间轮,从序号为1的数组开始
	for idx := 1; idx < cycleSize; idx++ {
		if 0 != p.secondSlice[idx-1].Len() { // 如果(idx-1)的cycle中还有元素,则不需要(idx-1)之后的cycle向前移动
			break
		}
		c := &p.secondSlice[idx]
		for e := c.Front(); e != nil; e = next {
			t := e.Value.(*Second)
			if !t.isValid() {
				next = e.Next()
				c.Remove(e)
				continue
			}
			if t.expire <= timestamp {
				p.opts.outgoingTimeoutChan <- t
				next = e.Next()
				c.Remove(e)
				continue
			}
			if newIdx := findPrevCycleIdx(t.expire-timestamp, idx); idx != newIdx {
				next = e.Next()
				c.Remove(e)
				p.pushBackCycle(t, newIdx, false)
				continue
			}
			break
		}
	}
}
