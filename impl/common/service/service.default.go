package service

import (
	"sync"
	"xcore/impl/common/bench"
	xerror "xcore/lib/error"
	xtime "xcore/lib/time"
	xtimer "xcore/lib/timer"
)

type DefaultService struct {
	BenchMgr bench.Mgr
	BenchSub bench.IBenchSub
	GroupID  uint32 // 组ID
	Name     string // 名称
	ID       uint32 // ID

	TimeMgr  xtime.Mgr
	TimerMgr xtimer.Mgr
	Opts     *options

	BusChannel          chan interface{} //  总线 channel
	BusChannelWaitGroup sync.WaitGroup
}

func NewDefaultService() *DefaultService {
	return &DefaultService{}
}

func (p *DefaultService) Start() (err error) {
	return xerror.NotImplemented
}

func (p *DefaultService) Stop() (err error) {
	return xerror.NotImplemented
}
