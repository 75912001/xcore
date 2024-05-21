package log

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"sync"
	libruntime "xcore/lib/runtime"
)

// 创建entry的内存池选项
type entryPoolOptions struct {
	pool         *sync.Pool
	enablePool   *bool         // 使用内存池 default: true
	newEntryFunc func() *entry // 创建 entry 的方法
}

// options contains options to configure a server stdInstance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type options struct {
	level            *int    // 日志等级允许的最小等级 default: LevelOn
	absPath          *string // 日志绝对路径 default: 当前执行的程序-绝对路径,指向启动当前进程的可执行文件-目录路径. e.g.:absPath/log
	isReportCaller   *bool   // 是否打印调用信息 default: true
	namePrefix       *string // 日志名 前缀 default: 当前执行的程序名称
	isWriteFile      *bool   // 是否写文件 default: true
	entryPoolOptions *entryPoolOptions
	hooks            LevelHooks // 各日志级别对应的钩子
}

// NewOptions 新的Options
func NewOptions() *options {
	ops := new(options)
	ops.entryPoolOptions = &entryPoolOptions{}
	ops.hooks = make(LevelHooks)
	return ops
}

func (p *options) WithLevel(level int) *options {
	p.level = &level
	return p
}

func (p *options) WithAbsPath(absPath string) *options {
	p.absPath = &absPath
	return p
}

func (p *options) WithIsReportCaller(isReportCaller bool) *options {
	p.isReportCaller = &isReportCaller
	return p
}

func (p *options) WithNamePrefix(namePrefix string) *options {
	p.namePrefix = &namePrefix
	return p
}

func (p *options) WithHooks(hooks LevelHooks) *options {
	p.hooks = hooks
	return p
}

func (p *options) WithIsWriteFile(isWriteFile bool) *options {
	p.isWriteFile = &isWriteFile
	return p
}

func (p *options) WithEnablePool(enablePool bool) *options {
	p.entryPoolOptions.enablePool = &enablePool
	return p
}

func (p *options) IsEnablePool() bool {
	return *p.entryPoolOptions.enablePool
}

// AddHooks 添加钩子
func (p *options) AddHooks(hook Hook) *options {
	p.hooks.add(hook)
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
		if opt.level != nil {
			so.level = opt.level
		}
		if opt.absPath != nil {
			so.absPath = opt.absPath
		}
		if opt.isReportCaller != nil {
			so.isReportCaller = opt.isReportCaller
		}
		if opt.namePrefix != nil {
			so.namePrefix = opt.namePrefix
		}
		if opt.hooks != nil {
			so.hooks = opt.hooks
		}
		if opt.isWriteFile != nil {
			so.isWriteFile = opt.isWriteFile
		}
		if opt.entryPoolOptions.enablePool != nil {
			so.entryPoolOptions.enablePool = opt.entryPoolOptions.enablePool
		}
	}
	return so
}

// 配置
func configure(opts *options) error {
	if opts.level == nil {
		var level = LevelOn
		opts.level = &level
	}
	if opts.absPath == nil {
		executablePath, err := libruntime.GetExecutablePath()
		if err != nil {
			return errors.WithMessage(err, libruntime.Location())
		}
		executablePath = filepath.Join(executablePath, "log")
		opts.absPath = &executablePath
	}
	if err := os.MkdirAll(*opts.absPath, os.ModePerm); err != nil {
		return errors.WithMessage(err, libruntime.Location())
	}
	if opts.isReportCaller == nil {
		var reportCaller = true
		opts.isReportCaller = &reportCaller
	}
	if opts.namePrefix == nil {
		executableName, err := libruntime.GetExecutableName()
		if err != nil {
			return errors.WithMessage(err, libruntime.Location())
		}
		opts.namePrefix = &executableName
	}
	if opts.isWriteFile == nil {
		var writeFile = true
		opts.isWriteFile = &writeFile
	}
	if opts.entryPoolOptions == nil {
		opts.entryPoolOptions = &entryPoolOptions{}
	}
	if opts.entryPoolOptions.enablePool == nil {
		var enablePool = true
		opts.entryPoolOptions.enablePool = &enablePool
	}
	return nil
}
