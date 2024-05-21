package log

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
	libconstants "xcore/lib/constants"
	liberror "xcore/lib/error"
	libruntime "xcore/lib/runtime"
	libtime "xcore/lib/time"
)

// NewMgr 创建日志管理器
func NewMgr(opts ...*options) (*mgr, error) {
	element := new(mgr)
	err := withOptions(element, opts...)
	if err != nil {
		return nil, err
	}
	err = start(element)
	if err != nil {
		stdInstance = nil
		return nil, err
	} else {
		stdInstance = element
	}
	return element, nil
}

// 日志管理器
type mgr struct {
	options         *options
	loggerSlice     [LevelOn]*log.Logger // 日志实例 [note]:使用时,注意协程安全
	logChan         chan *entry          // 日志写入通道
	waitGroupOutPut sync.WaitGroup       // 同步锁 用于日志退出时,等待完全输出
	logDuration     int                  // 日志分割刻度,变化时,使用新的日志文件 按天或者小时  e.g.:20210819或2021081901
	openFiles       []*os.File           // 当前打开的文件
	timeMgr         *libtime.Mgr
}

func withOptions(p *mgr, opts ...*options) error {
	p.options = mergeOptions(opts...)
	if err := configure(p.options); err != nil {
		return errors.WithMessage(err, libruntime.Location())
	}
	return nil
}

// 开始
func start(p *mgr) error {
	// 初始化logger
	for i := LevelOff; i < LevelOn; i++ {
		p.loggerSlice[i] = log.New(os.Stdout, "", 0)
	}
	p.logChan = make(chan *entry, logChannelCapacity)
	p.timeMgr = &libtime.Mgr{}
	// 初始化各级别的日志输出
	if err := newWriters(p); err != nil {
		return errors.WithMessage(err, libruntime.Location())
	}
	// 内存池
	if p.options.IsEnablePool() {
		p.options.entryPoolOptions.pool = &sync.Pool{
			New: func() interface{} {
				element := &entry{}
				return element
			},
		}
		p.options.entryPoolOptions.newEntryFunc = func() *entry {
			return p.options.entryPoolOptions.pool.Get().(*entry)
		}
	} else {
		p.options.entryPoolOptions.newEntryFunc = func() *entry {
			element := &entry{}
			return element
		}
	}
	p.waitGroupOutPut.Add(1)
	go func() {
		defer func() {
			if libruntime.IsRelease() {
				if err := recover(); err != nil {
					PrintErr(libconstants.GoroutinePanic, err, string(debug.Stack()))
				}
			}
			p.waitGroupOutPut.Done()
			PrintInfo(libconstants.GoroutineDone)
		}()
		doLog(p)
	}()
	return nil
}

// GetLevel 获取日志等级
func (p *mgr) GetLevel() int {
	return *p.options.level
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
func doLog(p *mgr) {
	for v := range p.logChan {
		p.fireHooks(v)
		// 检查自动切换日志
		if p.logDuration != p.getLogDuration(v.time.Unix()) {
			if err := newWriters(p); err != nil {
				PrintfErr("log duration changed, init writers failed, err:%v", err)
				if p.options.IsEnablePool() {
					reset(v)
					p.options.entryPoolOptions.pool.Put(v)
				}
				continue
			}
		}
		if *p.options.isWriteFile {
			p.loggerSlice[v.level].Print(formatLogData(v))
		}
		if p.options.IsEnablePool() {
			reset(v)
			p.options.entryPoolOptions.pool.Put(v)
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
func newWriters(p *mgr) error {
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

func (p *mgr) NewEntry() *entry {
	entry := p.options.entryPoolOptions.newEntryFunc()
	return entry
}

// log 记录日志
func (p *mgr) log(entry *entry, level int, v ...interface{}) {
	withLevel(entry, level)
	withTime(entry, p.timeMgr.NowTime())
	withMessage(entry, fmt.Sprintln(v...))
	if *p.options.isReportCaller {
		pc, _, line, ok := runtime.Caller(calldepth2)
		funcName := libconstants.Unknown
		if !ok {
			line = 0
		} else {
			funcName = runtime.FuncForPC(pc).Name()
		}
		withCallerInfo(entry, fmt.Sprintf(callerInfoFormat, line, funcName))
	}
	p.logChan <- entry
}

// logf 记录日志
func (p *mgr) logf(entry *entry, level int, format string, v ...interface{}) {
	withLevel(entry, level)
	withTime(entry, p.timeMgr.NowTime())
	withMessage(entry, fmt.Sprintf(format, v...))
	if *p.options.isReportCaller {
		pc, _, line, ok := runtime.Caller(calldepth2)
		funcName := libconstants.Unknown
		if !ok {
			line = 0
		} else {
			funcName = runtime.FuncForPC(pc).Name()
		}
		withCallerInfo(entry, fmt.Sprintf(callerInfoFormat, line, funcName))
	}
	p.logChan <- entry
}

// Trace 踪迹日志
func (p *mgr) Trace(v ...interface{}) {
	if p.GetLevel() < LevelTrace {
		return
	}
	p.log(p.NewEntry(), LevelTrace, v...)
}

func (p *mgr) TraceWithEntry(entry *entry, v ...interface{}) {
	if p.GetLevel() < LevelTrace {
		return
	}
	p.log(entry, LevelTrace, v...)
}

// Tracef 踪迹日志
func (p *mgr) Tracef(format string, v ...interface{}) {
	if p.GetLevel() < LevelTrace {
		return
	}
	p.logf(p.NewEntry(), LevelTrace, format, v...)
}

func (p *mgr) TracefWithEntry(entry *entry, format string, v ...interface{}) {
	if p.GetLevel() < LevelTrace {
		return
	}
	p.logf(entry, LevelTrace, format, v...)
}

// Debug 调试日志
func (p *mgr) Debug(v ...interface{}) {
	if p.GetLevel() < LevelDebug {
		return
	}
	p.log(p.NewEntry(), LevelDebug, v...)
}

// DebugLazy 调试日志-惰性
//
//	等级满足之后才会计算
func (p *mgr) DebugLazy(vFunc func() []interface{}) {
	if p.GetLevel() < LevelDebug {
		return
	}
	v := vFunc()
	p.log(p.NewEntry(), LevelDebug, v...)
}

// Debugf 调试日志
func (p *mgr) Debugf(format string, v ...interface{}) {
	if p.GetLevel() < LevelDebug {
		return
	}
	p.logf(p.NewEntry(), LevelDebug, format, v...)
}

// Info 信息日志
func (p *mgr) Info(v ...interface{}) {
	if p.GetLevel() < LevelInfo {
		return
	}
	p.log(p.NewEntry(), LevelInfo, v...)
}

// Infof 信息日志
func (p *mgr) Infof(format string, v ...interface{}) {
	if p.GetLevel() < LevelInfo {
		return
	}
	p.logf(p.NewEntry(), LevelInfo, format, v...)
}

// Warn 警告日志
func (p *mgr) Warn(v ...interface{}) {
	if p.GetLevel() < LevelWarn {
		return
	}
	p.log(p.NewEntry(), LevelWarn, v...)
}

// Warnf 警告日志
func (p *mgr) Warnf(format string, v ...interface{}) {
	if p.GetLevel() < LevelWarn {
		return
	}
	p.logf(p.NewEntry(), LevelWarn, format, v...)
}

// Error 错误日志
func (p *mgr) Error(v ...interface{}) {
	if p.GetLevel() < LevelError {
		return
	}
	p.log(p.NewEntry(), LevelError, v...)
}

// Errorf 错误日志
func (p *mgr) Errorf(format string, v ...interface{}) {
	if p.GetLevel() < LevelError {
		return
	}
	p.logf(p.NewEntry(), LevelError, format, v...)
}

// Fatal 致命日志
func (p *mgr) Fatal(v ...interface{}) {
	if p.GetLevel() < LevelFatal {
		return
	}
	p.log(p.NewEntry(), LevelFatal, v...)
}

// Fatalf 致命日志
func (p *mgr) Fatalf(format string, v ...interface{}) {
	if p.GetLevel() < LevelFatal {
		return
	}
	p.logf(p.NewEntry(), LevelFatal, format, v...)
}
