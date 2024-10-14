package etcd

import (
	"github.com/pkg/errors"
	"time"
	xcallback "xcore/lib/callback"
	xerror "xcore/lib/error"
	xruntime "xcore/lib/runtime"
)

var (
	grantLeaseRetryDuration     = time.Second * 3 // 授权租约 重试 间隔时长
	grantLeaseMaxRetriesDefault = 600             // 授权租约 最大 重试次数
	dialTimeoutDefault          = time.Second * 5 // dialTimeout is the timeout for failing to establish a connection.
)

// KV key-value pair
type KV struct {
	Key   string
	Value string
}

type options struct {
	addrs                []string       // 地址
	ttl                  *int64         // Time To Live, etcd内部会按照 ttl/3 的时间(最小1秒),保持连接
	grantLeaseMaxRetries *int           // 授权租约 最大 重试次数 [default:600]
	kvSlice              []KV           // 事件
	dialTimeout          *time.Duration // dialTimeout is the timeout for failing to establish a connection. [default:time.Second*5]
	ICallBack            xcallback.ICallBack
	eventChan            chan<- interface{} // 传出 channel
}

// NewOptions 新的Options
func NewOptions() *options {
	return &options{}
}

func (p *options) WithAddrs(addrs []string) *options {
	p.addrs = p.addrs[0:0]
	p.addrs = append(p.addrs, addrs...)
	return p
}

func (p *options) WithTTL(ttl int64) *options {
	p.ttl = &ttl
	return p
}

func (p *options) WithGrantLeaseMaxRetries(retries int) *options {
	p.grantLeaseMaxRetries = &retries
	return p
}

func (p *options) WithKV(kvSlice []KV) *options {
	p.kvSlice = p.kvSlice[0:0]
	p.kvSlice = append(p.kvSlice, kvSlice...)
	return p
}

func (p *options) WithDialTimeout(dialTimeout time.Duration) *options {
	p.dialTimeout = &dialTimeout
	return p
}

func (p *options) WithICallBack(cb xcallback.ICallBack) *options {
	p.ICallBack = cb
	return p
}

func (p *options) WithEventChan(eventChan chan<- interface{}) *options {
	p.eventChan = eventChan
	return p
}

func mergeOptions(opts ...*options) *options {
	no := NewOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if len(opt.addrs) != 0 {
			no.WithAddrs(opt.addrs)
		}
		if opt.ttl != nil {
			no.WithTTL(*opt.ttl)
		}
		if opt.grantLeaseMaxRetries != nil {
			no.WithGrantLeaseMaxRetries(*opt.grantLeaseMaxRetries)
		}
		if len(opt.kvSlice) != 0 {
			no.WithKV(opt.kvSlice)
		}
		if opt.dialTimeout != nil {
			no.WithDialTimeout(*opt.dialTimeout)
		}
		if opt.ICallBack != nil {
			no.WithICallBack(opt.ICallBack)
		}
		if opt.eventChan != nil {
			no.WithEventChan(opt.eventChan)
		}
	}
	return no
}

// 配置
func configure(opts *options) error {
	if len(opts.addrs) == 0 {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.ttl == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.grantLeaseMaxRetries == nil {
		var v = grantLeaseMaxRetriesDefault
		opts.grantLeaseMaxRetries = &v
	}
	//if len(opts.kvSlice) == 0 { // todo menglc
	//	return errors.WithMessage(xerror.Param, xruntime.Location())
	//}
	if opts.dialTimeout == nil {
		opts.WithDialTimeout(dialTimeoutDefault)
	}
	//if opts.ICallBack == nil { // todo menglc
	//	return errors.WithMessage(xerror.Param, xruntime.Location())
	//}
	if opts.eventChan == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	return nil
}
