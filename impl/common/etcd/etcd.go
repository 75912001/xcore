package etcd

import (
	"context"
	"dawn-server/impl/common/bench"
	xretcd "dawn-server/impl/xr/lib/etcd"
	xrutil "dawn-server/impl/xr/lib/util"
	"encoding/json"
	"path"
	"time"
	xbench "xcore/lib/bench"

	"github.com/pkg/errors"
)

// todo menglc 放到单独 生成,解析 etcd 数据的功能实现
// EtcdValueJson etcd 通讯的数据,由服务中的数据生成,定时更新->etcd->服务
type EtcdValueJson struct {
	ServiceNet    xbench.ServiceNet `json:"serviceNet"`    // 有:直接使用. 没有:使用 benchJson.ServiceNet
	Version       string            `json:"version"`       // 有:直接使用. 没有:使用 base.version 生成
	AvailableLoad uint32            `json:"availableLoad"` // 剩余可用负载, 可用资源数
	SecondOffset  int32             `json:"secondOffset"`  // 服务 时间(秒)偏移量
}

//Key   string        `json:"key"`   // [required] common.ProjectName/common.EtcdWatchMsgTypeService/groupID/serviceName/serviceID
//Value EtcdValueJson `json:"value"` // [required]

// Parse
// e.g.:/${projectName}/${EtcdWatchMsgType}/${groupID}/${serviceName}/${serviceID}
func Parse(key string) (msgType string, groupID string, serviceName string, serviceID string) {
	serviceID = path.Base(key)

	key = path.Dir(key)
	serviceName = path.Base(key)

	key = path.Dir(key)
	groupID = path.Base(key)

	key = path.Dir(key)
	msgType = path.Base(key)
	return msgType, groupID, serviceName, serviceID
}

// EtcdStart 启动Etcd
func EtcdStart(conf *bench.Etcd, BusChannel chan interface{}, onFunc xretcd.OnFunc) error {
	if len(conf.Addrs) != 0 {
		etcdValue, err := json.Marshal(conf.Value)
		if err != nil {
			return errors.WithMessagef(err, xrutil.GetCodeLocation(1).String())
		}
		var kvSlice []xretcd.KV
		kvSlice = append(kvSlice, xretcd.KV{
			Key:   conf.Key,
			Value: string(etcdValue),
		})
		err = xretcd.GetInstance().Start(context.TODO(),
			xretcd.NewOptions().
				SetAddr(conf.Addrs).
				SetTTL(conf.TTL).
				SetDialTimeout(5*time.Second).
				SetKV(kvSlice).SetOnFunc(onFunc).
				SetEventChan(BusChannel),
		)
		if err != nil {
			return errors.WithMessagef(err, xrutil.GetCodeLocation(1).String())
		}

		// 续租
		err = xretcd.GetInstance().KeepAlive(context.TODO())
		if err != nil {
			return errors.WithMessagef(err, xrutil.GetCodeLocation(1).String())
		}
	}

	return nil
}
