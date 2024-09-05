package service

import (
	"github.com/pkg/errors"
	"os"
	"path"
	"path/filepath"
	xetcd "xcore/lib/etcd"
	xruntime "xcore/lib/runtime"
	xutil "xcore/lib/util"
)

// options contains options to configure a server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type options struct {
	Path      *string // 路径
	BenchPath *string // 配置文件路径

	NatsHandler    OnNatsHandler    // Nats 处理函数
	DefaultHandler OnDefaultHandler // default 处理函数 // todo menglc 改名 -> BusHandler

	EtcdHandler            xetcd.OnFunc // etcd 处理函数
	EtcdWatchServicePrefix *string      // etcd 关注 服务 前缀
	EtcdWatchCommandPrefix *string      // etcd 关注 命令 前缀
	EtcdWatchGMPrefix      *string      // etcd 关注 GM 前缀
}

// NewOptions 新的Options
func NewOptions() *options {
	ops := new(options)
	return ops
}

func (p *options) WithPath(path string) *options {
	p.Path = &path
	return p
}
func (p *options) WithBenchPath(benchPath string) *options {
	p.BenchPath = &benchPath
	return p
}
func (p *options) WithNatsHandler(natsHandler OnNatsHandler) *options {
	p.NatsHandler = natsHandler
	return p
}
func (p *options) WithDefaultHandler(defaultHandler OnDefaultHandler) *options {
	p.DefaultHandler = defaultHandler
	return p
}
func (p *options) WithEtcdHandler(etcdHandler xetcd.OnFunc) *options {
	p.EtcdHandler = etcdHandler
	return p
}
func (p *options) WithEtcdWatchServicePrefix(etcdWatchServicePrefix string) *options {
	p.EtcdWatchServicePrefix = &etcdWatchServicePrefix
	return p
}
func (p *options) WithEtcdWatchCommandPrefix(etcdWatchCommandPrefix string) *options {
	p.EtcdWatchCommandPrefix = &etcdWatchCommandPrefix
	return p
}
func (p *options) WithEtcdWatchGMPrefix(etcdWatchGMPrefix string) *options {
	p.EtcdWatchGMPrefix = &etcdWatchGMPrefix
	return p
}

// mergeOptions combines the given *options into a single *options in a last one wins fashion.
// The specified options are merged with the existing options on the server, with the specified options taking
// precedence.
func mergeOptions(opts ...*options) *options {
	so := NewOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.Path != nil {
			so.Path = opt.Path
		}
		if opt.BenchPath != nil {
			so.BenchPath = opt.BenchPath
		}
		if opt.NatsHandler != nil {
			so.NatsHandler = opt.NatsHandler
		}
		if opt.DefaultHandler != nil {
			so.DefaultHandler = opt.DefaultHandler
		}
		if opt.EtcdHandler != nil {
			so.EtcdHandler = opt.EtcdHandler
		}
		if opt.EtcdWatchServicePrefix != nil {
			so.EtcdWatchServicePrefix = opt.EtcdWatchServicePrefix
		}
		if opt.EtcdWatchCommandPrefix != nil {
			so.EtcdWatchCommandPrefix = opt.EtcdWatchCommandPrefix
		}
		if opt.EtcdWatchGMPrefix != nil {
			so.EtcdWatchGMPrefix = opt.EtcdWatchGMPrefix
		}
	}
	return so
}

func GetCurrentPath() (currentPath string, err error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}

// 配置
func configure(opts *options) error {
	if opts.Path == nil {
		pathValue, err := xutil.GetCurrentPath()
		if err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
		opts.Path = &pathValue
	}
	if opts.BenchPath == nil {
		benchPath := path.Join(*opts.Path, "bench.json")
		opts.BenchPath = &benchPath
	}
	return nil
}
