package service

import (
	"github.com/pkg/errors"
	"path"
	xconstants "xcore/lib/constants"
	xerror "xcore/lib/error"
	xlog "xcore/lib/log"
	xnettcp "xcore/lib/net/tcp"
	xruntime "xcore/lib/runtime"
)

// options contains options to configure a server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type options struct {
	BenchPath *string // 配置文件路径 [default]: 可执行程序所在目录的 bench.json
	xnettcp.IHandler
}

// NewOptions 新的Options
func NewOptions() *options {
	ops := new(options)
	return ops
}

func (p *options) WithBenchPath(benchPath string) *options {
	p.BenchPath = &benchPath
	return p
}

func (p *options) WithHandler(handler xnettcp.IHandler) *options {
	p.IHandler = handler
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
		if opt.BenchPath != nil {
			so.BenchPath = opt.BenchPath
		}
		if opt.IHandler != nil {
			so.IHandler = opt.IHandler
		}
	}
	return so
}

// 配置
func configure(opts *options) error {
	if opts.BenchPath == nil {
		benchPath, err := xruntime.GetExecutablePath()
		if err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
		benchPath = path.Join(benchPath, xconstants.ServiceConfigFile)
		opts.BenchPath = &benchPath
	}
	if opts.IHandler == nil {
		xlog.PrintErr(errors.WithMessage(xerror.Param, xruntime.Location()))
	}
	return nil
}
