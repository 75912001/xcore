package log

import (
	"context"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
	libconstants "xcore/lib/constants"
	liberror "xcore/lib/error"
	libruntime "xcore/lib/runtime"
	libtime "xcore/lib/time"
)

var (
	instance *mgr
	once     sync.Once
)

// GetInstance 获取
func GetInstance() *mgr {
	once.Do(func() {
		instance = new(mgr)
	})
	return instance
}

// IsEnable 是否 开启
func IsEnable() bool {
	if instance == nil {
		return false
	}
	return GetInstance().logChan != nil
}

// mgr 日志管理器
type mgr struct {
	options         *options
	loggerSlice     [LevelOn]*log.Logger // 日志实例 note:此处非协程安全
	logChan         chan *entry          // 日志写入通道
	waitGroupOutPut sync.WaitGroup       // 同步锁 用于日志退出时,等待完全输出
	logDuration     int                  // 日志分割刻度,变化时,使用新的日志文件 按天或者小时  e.g.:20210819或2021081901
	openFiles       []*os.File           // 当前打开的文件
	pool            *sync.Pool
	timeMgr         *libtime.Mgr
}

// GetLevel 获取日志等级
func (p *mgr) GetLevel() int {
	return *p.options.level
}

// Start 开始
func (p *mgr) Start(opts ...*options) error {
	p.options = mergeOptions(opts...)
	if err := configure(p.options); err != nil {
		return errors.WithMessage(err, libruntime.Location())
	}
	// 初始化logger
	for i := LevelOff; i < LevelOn; i++ {
		p.loggerSlice[i] = log.New(os.Stdout, "", 0)
	}
	p.logChan = make(chan *entry, logChannelCapacity)
	p.timeMgr = &libtime.Mgr{}
	// 初始化各级别的日志输出
	if err := p.newWriters(); err != nil {
		return errors.WithMessage(err, libruntime.Location())
	}
	// 内存池
	if p.options.IsEnablePool() {
		p.pool = &sync.Pool{
			New: func() interface{} {
				return new(entry)
			},
		}
	}
	p.waitGroupOutPut.Add(1)
	go func() {
		defer func() {
			if libruntime.IsRelease() {
				if err := recover(); err != nil {
					PrintErr(libconstants.GoroutinePanic, err, debug.Stack())
				}
			}
			p.waitGroupOutPut.Done()
			PrintInfo(libconstants.GoroutineDone)
		}()
		p.doLog()
	}()
	return nil
}

// getLogDuration 取得日志刻度
func (p *mgr) getLogDuration(sec int64) int {
	var logFormat string
	if libruntime.IsRelease() {
		logFormat = "2006010215" //年月日小时
	} else {
		logFormat = "20060102" //年月日
	}
	durationStr := time.Unix(sec, 0).Format(logFormat)
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		PrintfErr("strconv.Atoi sec:%v durationStr:%v err:%v", sec, durationStr, err)
	}
	return duration
}

// doLog 处理日志
func (p *mgr) doLog() {
	for v := range p.logChan {
		p.fireHooks(v)
		// 检查自动切换日志
		if p.logDuration != p.getLogDuration(v.time.Unix()) {
			if err := p.newWriters(); err != nil {
				PrintfErr("log duration changed, init writers failed, err:%v", err)
				if p.options.IsEnablePool() {
					v.reset()
					p.pool.Put(v)
				}
				continue
			}
		}
		if *p.options.isWriteFile {
			p.loggerSlice[v.level].Print(v.formatMessage())
		}
		if p.options.IsEnablePool() {
			v.reset()
			p.pool.Put(v)
		}
	}
	// goroutine 退出,再设置chan为nil, (如果没有退出就设置为nil, 读chan == nil  会 block)
	p.logChan = nil
}

// SetLevel 设置日志等级
func (p *mgr) SetLevel(level int) error {
	if level < LevelOff || LevelOn < level {
		return errors.WithMessage(liberror.Level, libruntime.Location())
	}
	p.options.WithLevel(level)
	return nil
}

// newWriters 初始化各级别的日志输出
func (p *mgr) newWriters() error {
	// 检查是否要关闭文件
	for i := range p.openFiles {
		if err := p.openFiles[i].Close(); err != nil {
			return errors.WithMessage(err, libruntime.Location())
		}
	}
	second := p.timeMgr.NowTime().Unix()
	logDuration := p.getLogDuration(second)
	normalWriter, err := newNormalFileWriter(*p.options.absPath, *p.options.namePrefix, logDuration)
	if err != nil {
		return errors.WithMessage(err, libruntime.Location())
	}
	errorWriter, err := newErrorFileWriter(*p.options.absPath, *p.options.namePrefix, logDuration)
	if err != nil {
		return errors.WithMessage(err, libruntime.Location())
	}
	p.logDuration = logDuration
	allWriter := io.MultiWriter(normalWriter, errorWriter)
	// 标准输出,标准错误重定向
	stdOut.SetOutput(normalWriter)
	stdErr.SetOutput(allWriter)

	p.loggerSlice[LevelFatal].SetOutput(allWriter)
	p.loggerSlice[LevelError].SetOutput(allWriter)
	p.loggerSlice[LevelWarn].SetOutput(allWriter)
	p.loggerSlice[LevelInfo].SetOutput(normalWriter)
	p.loggerSlice[LevelDebug].SetOutput(normalWriter)
	p.loggerSlice[LevelTrace].SetOutput(normalWriter)

	// 记录打开的文件
	p.openFiles = p.openFiles[0:0]
	p.openFiles = append(p.openFiles, normalWriter)
	p.openFiles = append(p.openFiles, errorWriter)

	return nil
}

// Stop 停止
func (p *mgr) Stop() error {
	if p.logChan != nil {
		// close chan, for range 读完chan会退出.
		close(p.logChan)
		// 等待logChan 的for range 退出.
		p.waitGroupOutPut.Wait()
	}

	// 检查是否要关闭文件
	if len(p.openFiles) > 0 {
		for i := range p.openFiles {
			_ = p.openFiles[i].Close()
		}
		p.openFiles = p.openFiles[0:0]
	}
	return nil
}

// fireHooks 处理钩子
func (p *mgr) fireHooks(entry *entry) {
	if 0 == len(p.options.hooks) {
		return
	}

	err := p.options.hooks.fire(entry)
	if err != nil {
		PrintfErr("failed to fire hook. err:%v", err)
	}
}

// WithField 由field创建日志信息 默认大小2(cap:2*2=4)
func (p *mgr) WithField(key string, value interface{}) *entry {
	entry := newEntry()
	entry.extendFields = make(extendFields, 0, 4)
	return entry.WithExtendField(key, value)
}

// WithFields 由fields创建日志信息 默认大小4(cap:4*2=8)
func (p *mgr) WithFields(f extendFields) *entry {
	entry := newEntry()
	entry.extendFields = make(extendFields, 0, 8)
	return entry.WithExtendFields(f)
}

// WithContext 由ctx创建日志信息
func (p *mgr) WithContext(ctx context.Context) *entry {
	entry := newEntry()
	return entry.WithContext(ctx)
}

// log 记录日志
func (p *mgr) log(level int, v ...interface{}) {
	entry := newEntry()
	entry.log(level, calldepth3, v...)
}

// logf 记录日志
func (p *mgr) logf(level int, format string, v ...interface{}) {
	entry := newEntry()
	entry.logf(level, calldepth3, format, v...)
}

// Trace 踪迹日志
func (p *mgr) Trace(level int, v ...interface{}) {
	if *p.options.level < LevelTrace {
		return
	}
	p.log(level, v...)
}

// Tracef 踪迹日志
func (p *mgr) Tracef(level int, format string, v ...interface{}) {
	if *p.options.level < LevelTrace {
		return
	}
	p.logf(level, format, v...)
}

// Debug 调试日志
func (p *mgr) Debug(level int, v ...interface{}) {
	if *p.options.level < LevelDebug {
		return
	}
	p.log(level, v...)
}

// Debugf 调试日志
func (p *mgr) Debugf(level int, format string, v ...interface{}) {
	if *p.options.level < LevelDebug {
		return
	}
	p.logf(level, format, v...)
}

// Info 信息日志
func (p *mgr) Info(level int, v ...interface{}) {
	if *p.options.level < LevelInfo {
		return
	}
	p.log(level, v...)
}

// Infof 信息日志
func (p *mgr) Infof(level int, format string, v ...interface{}) {
	if *p.options.level < LevelInfo {
		return
	}
	p.logf(level, format, v...)
}

// Warn 警告日志
func (p *mgr) Warn(level int, v ...interface{}) {
	if *p.options.level < LevelWarn {
		return
	}
	p.log(level, v...)
}

// Warnf 警告日志
func (p *mgr) Warnf(level int, format string, v ...interface{}) {
	if *p.options.level < LevelWarn {
		return
	}
	p.logf(level, format, v...)
}

// Error 错误日志
func (p *mgr) Error(level int, v ...interface{}) {
	if *p.options.level < LevelError {
		return
	}
	p.log(level, v...)
}

// Errorf 错误日志
func (p *mgr) Errorf(level int, format string, v ...interface{}) {
	if *p.options.level < LevelError {
		return
	}
	p.logf(level, format, v...)
}

// Fatal 致命日志
func (p *mgr) Fatal(level int, v ...interface{}) {
	if *p.options.level < LevelFatal {
		return
	}
	p.log(level, v...)
}

// Fatalf 致命日志
func (p *mgr) Fatalf(level int, format string, v ...interface{}) {
	if *p.options.level < LevelFatal {
		return
	}
	p.logf(level, format, v...)
}
