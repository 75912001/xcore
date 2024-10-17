package gateway

import (
	xetcd "xcore/lib/etcd"
	xlog "xcore/lib/log"
)

// EtcdKeyValue etcd 刷新 key value
func EtcdKeyValue(arg ...interface{}) error {
	key := arg[0].(string)
	valueJson := arg[1].(*xetcd.ValueJson)
	xlog.PrintInfo("etcd key:%v, value:%v", key, valueJson)
	return nil
}
