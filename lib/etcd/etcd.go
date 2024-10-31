package etcd

import (
	"context"
	etcdclientv3 "go.etcd.io/etcd/client/v3"
	xcontrol "xcore/lib/control"
)

type IEtcd interface {
	Start(ctx context.Context) (err error)
	Stop() (err error)

	Put(key string, value string) (*etcdclientv3.PutResponse, error)
}

type Event struct {
	ICallBack xcontrol.ICallBack
}

type CallbackFun func(arg ...interface{}) error
