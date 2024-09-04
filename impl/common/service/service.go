package service

import (
    "sync"
    "xcore/impl/common/bench"
    xerror "xcore/lib/error"
    xtime "xcore/lib/time"
    xtimer "xcore/lib/timer"
)

var GService IService

type IService interface {
    Start() (err error)
    Stop() (err error)
}

type DefaultService struct {
    BenchMgr bench.Mgr

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
