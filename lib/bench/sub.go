package bench

import (
	"encoding/json"
	"github.com/pkg/errors"
	"time"
	xruntime "xcore/lib/runtime"
)

type Sub struct {
	*Jaeger  `json:"jaeger"`  // todo
	*MongoDB `json:"mongoDB"` // todo
	*Redis   `json:"redis"`   // todo
	*NATS    `json:"nats"`    // todo
}

func (p *Sub) Unmarshal(strJson string) error {
	if err := json.Unmarshal([]byte(strJson), &p); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}

func (p *Sub) GetJaeger() *Jaeger {
	return p.Jaeger
}

func (p *Sub) GetMongoDB() *MongoDB {
	return p.MongoDB
}

func (p *Sub) GetRedis() *Redis {
	return p.Redis
}

func (p *Sub) GetNATS() *NATS {
	return p.NATS
}

type Jaeger struct {
	Addrs []string `json:"addrs"`
}

type MongoDB struct {
	Addrs           []string       `json:"addrs"`
	User            *string        `json:"user"`
	Password        *string        `json:"password"`
	DBName          *string        `json:"dbName"`          // 数据库名称 [default]: todo
	MaxPoolSize     *uint64        `json:"maxPoolSize"`     // 连接池最大数量,该数量应该与并发数量匹配 [default]: todo
	MinPoolSize     *uint64        `json:"minPoolSize"`     // 池最小数量 [default]: todo
	TimeoutDuration *time.Duration `json:"timeoutDuration"` // 操作超时时间 [default]: todo
	MaxConnIdleTime *time.Duration `json:"maxConnIdleTime"` // 指定连接在连接池中保持空闲的最长时间 [default]: todo
	MaxConnecting   *uint64        `json:"maxConnecting"`   // 指定连接池可以同时建立的最大连接数 [default]: todo
	DBAsync         *DBAsync       `json:"dbAsync"`         // DB异步消费配置
}

type Redis struct {
	Addrs    []string `json:"addrs"`
	Password *string  `json:"password"`
}

type NATS struct {
	Addrs    []string `json:"addrs"`
	User     *string  `json:"user"`     // 用户 default: todo
	Password *string  `json:"password"` // 密码 default: todo
}

type DBAsync struct {
	ChanCnt              *uint32 `json:"chanCnt"`              // DB异步消费chan数量. 为0或者没有则不开启异步消费
	Model                *uint32 `json:"model"`                // DB异步消费模型 [default] todo
	BulkWriteMax         *uint32 `json:"bulkWriteMax"`         // DB合并写 单个集合最大合批数量 [default] todo
	BulkWriteMillisecond *uint32 `json:"bulkWriteMillisecond"` // DB合并写周期  单位毫秒 [default] todo
}
