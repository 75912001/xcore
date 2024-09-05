package bench

import (
	"encoding/json"
	"github.com/pkg/errors"
	"time"
	xruntime "xcore/lib/runtime"
)

// 配置-子项,用于各服务具体配置

type IBenchSub interface {
	Unmarshal(strJson string) error
}

type DefaultBenchSub struct {
	Jaeger  `json:"jaeger"`
	MongoDB `json:"mongoDB"`
	Redis   `json:"redis"`
	NATS    `json:"nats"`
}

func (p *DefaultBenchSub) Unmarshal(strJson string) error {
	if err := json.Unmarshal([]byte(strJson), &p); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}

type Jaeger struct {
	Addrs []string `json:"addrs"`
}

type MongoDB struct {
	Addrs           []string       `json:"addrs"`
	User            *string        `json:"user"`
	Password        *string        `json:"password"`
	DBName          *string        `json:"dbName"`          // 数据库名称 [default]: common.MongodbDatabaseNameDefault
	MaxPoolSize     *uint64        `json:"maxPoolSize"`     // 连接池最大数量,该数量应该与并发数量匹配 [default]: common.MongodbMaxPoolSizeDefault
	MinPoolSize     *uint64        `json:"minPoolSize"`     // 池最小数量 [default]: common.MongodbMinPoolSizeDefault
	TimeoutDuration *time.Duration `json:"timeoutDuration"` // 操作超时时间 [default]: common.MongodbTimeoutDurationDefault
	MaxConnIdleTime *time.Duration `json:"maxConnIdleTime"` // 指定连接在连接池中保持空闲的最长时间 [default]: common.MongodbMaxConnIdleTimeDefault
	MaxConnecting   *uint64        `json:"maxConnecting"`   // 指定连接池可以同时建立的最大连接数 [default]: common.MongodbMaxConnectingDefault
	DBAsync         *DBAsync       `json:"dbAsync"`         // DB异步消费配置
}

type Redis struct {
	Addrs    []string `json:"addrs"`
	Password *string  `json:"password"`
}

type NATS struct {
	Addrs    []string `json:"addrs"`
	User     *string  `json:"user"`     // 用户 default: nil
	Password *string  `json:"password"` // 密码 default: nil
}

type DBAsync struct {
	ChanCnt              *uint32 `json:"chanCnt"`              // DB异步消费chan数量. 为0或者没有则不开启异步消费
	Model                *uint32 `json:"model"`                // DB异步消费模型 [default] consumer.ModelAsyncOne
	BulkWriteMax         *uint32 `json:"bulkWriteMax"`         // DB合并写 单个集合最大合批数量 [default] common.BulkWriteMax
	BulkWriteMillisecond *uint32 `json:"bulkWriteMillisecond"` // DB合并写周期  单位毫秒 [default] common.BatchExecMaxMilliSecond
}
