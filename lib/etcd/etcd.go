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

// CallbackFun 回调函数
// key := arg[0].(string)
// valueJson := arg[1].(*xetcd.ValueJson)
type CallbackFun func(arg ...interface{}) error
