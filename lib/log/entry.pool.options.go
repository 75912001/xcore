package log

import (
	"sync"
	xswitch "xcore/lib/xswitch"
)

// entry的内存池选项
type entryPoolOptions struct {
	poolSwitch   xswitch.ISwitch // 内存池开关 [default]: true
	pool         *sync.Pool      // 内存池 [default]: &sync.Pool{New: func() interface{} { return newEntry() }}
	newEntryFunc func() *entry   // 创建 entry 的方法 [default]: func() *entry { return p.pool.Get().(*entry) }
}

// newEntryPoolOptions 新的entryPoolOptions
func newEntryPoolOptions() *entryPoolOptions {
	pool := &sync.Pool{
		New: func() interface{} {
			return newEntry()
		},
	}
	opt := &entryPoolOptions{
		poolSwitch: xswitch.NewDefaultSwitch(true),
		pool:       pool,
		newEntryFunc: func() *entry {
			return pool.Get().(*entry)
		},
	}
	return opt
}

func (p *entryPoolOptions) merge(opts ...*entryPoolOptions) *entryPoolOptions {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.poolSwitch.IsEnabled() {
			p.poolSwitch.Enable()
		} else {
			p.poolSwitch.Disable()
		}
		if opt.pool != nil {
			p.pool = opt.pool
		}
		if opt.newEntryFunc != nil {
			p.newEntryFunc = opt.newEntryFunc
		}
	}
	return p
}

// 配置
func (p *entryPoolOptions) configure() error {
	if p.poolSwitch.IsEnabled() {
		p.pool = &sync.Pool{
			New: func() interface{} {
				return newEntry()
			},
		}
		p.newEntryFunc = func() *entry {
			return p.pool.Get().(*entry)
		}
	} else {
		p.newEntryFunc = func() *entry {
			return newEntry()
		}
	}
	return nil
}

// 将内存放回池中
func (p *entryPoolOptions) put(value *entry) {
	if p.poolSwitch.IsEnabled() {
		value.reset()
		p.pool.Put(value)
	}
}
