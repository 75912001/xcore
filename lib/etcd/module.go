package etcd

import (
	"context"
	xcallback "xcore/lib/callback"
)

type IEtcd interface {
	Start(ctx context.Context) (err error)
	Stop() (err error)
}

type Event struct {
	ICallBack xcallback.ICallBack
}

type CallbackFun func(arg ...interface{}) error
