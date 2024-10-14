package etcd

import "context"

type IEtcd interface {
	Start(ctx context.Context) (err error)
	Stop() (err error)
}
