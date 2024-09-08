// 日志
// 使用系统log,自带锁
// 使用协程操作io输出日志
// release 每小时自动创建新的日志文件
// debug 每天自动创建新的日志文件

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
	xconstants "xcore/lib/constants"
	xerror "xcore/lib/error"
	xruntime "xcore/lib/runtime"
	xtime "xcore/lib/time"
)

var (
	mgrInstance *mgr
)

// 是否 启用
func isEnable() bool {
	if mgrInstance == nil {
		return false
	}
	if mgrInstance.logChan == nil {
		return false
	}
	return true
}

// NewMgr 创建日志管理器
func NewMgr(opts ...*options) (*mgr, error) {
	if isEnable() {
		return mgrInstance, nil
	}
	m := &mgr{}
	err := m.handleOptions(opts...)
	if err != nil {
		return nil, err
	}
	err = m.start()
	if err != nil {
		mgrInstance = nil
		return nil, err
	} else {
		mgrInstance = m
	}
	return m, nil
}

// 日志管理器
type mgr struct {
	options         *options
	loggerSlice     [LevelOn]*log.Logger // 日志实例 [note]:使用时,注意协程安全
	logChan         chan *entry          // 日志写入通道
	waitGroupOutPut sync.WaitGroup       // 同步锁 用于日志退出时,等待完全输出
	logDuration     int                  // 日志分割刻度,变化时,使用新的日志文件 按天或者小时  e.g.: 20210819 或 2021081901
	openFiles       []*os.File           // 当前打开的文件
	timeMgr         *xtime.Mgr
}

func (p *mgr) handleOptions(opts ...*options) error {
	p.options = NewOptions().merge(opts...)
	if err := p.options.configure(); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}

// 开始
func (p *mgr) start() error {
	// 初始化logger
	for i := LevelOff; i < LevelOn; i++ {
		p.loggerSlice[i] = log.New(os.Stdout, "", 0)
	}
	p.logChan = make(chan *entry, logChannelEntryCapacity)
	p.timeMgr = &xtime.Mgr{}
	// 初始化各级别的日志输出
	if err := newWriters(p); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	p.waitGroupOutPut.Add(1)
	go func() {
		defer func() {
			if xruntime.IsRelease() {
				if err := recover(); err != nil {
					PrintErr(xconstants.GoroutinePanic, err, string(debug.Stack()))
				}
			}
			p.waitGroupOutPut.Done()
			PrintInfo(xconstants.GoroutineDone)
		}()
		doLog(p)
	}()
	return nil
}

// GetLevel 获取日志等级
func (p *mgr) GetLevel() uint32 {
	return *p.options.level
}

// getLogDuration 取得日志刻度
func (p *mgr) getLogDuration(sec int64) int {
	var logFormat string
	if xruntime.IsRelease() {
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
		v.outString = formatLogData(v)
		p.fireHooks(v)
		// 检查自动切换日志
		if p.logDuration != p.getLogDuration(v.time.Unix()) {
			if err := newWriters(p); err != nil {
				PrintfErr("log duration changed, init writers failed, err:%v", err)
				p.options.entryPoolOptions.put(v)
				continue
			}
		}
		if *p.options.isWriteFile {
			p.loggerSlice[v.level].Print(v.outString)
		}
		p.options.entryPoolOptions.put(v)
	}
	// goroutine 退出,再设置chan为nil, (如果没有退出就设置为nil, 读chan == nil  会 block)
	p.logChan = nil
}

// SetLevel 设置日志等级
func (p *mgr) SetLevel(level uint32) error {
	if level < LevelOff || LevelOn < level {
		return errors.WithMessage(xerror.LogLevel, xruntime.Location())
	}
	p.options.WithLevel(level)
	return nil
}

// newWriters 初始化各级别的日志输出
func newWriters(p *mgr) error {
	// 检查是否要关闭文件
	for i := range p.openFiles {
		if err := p.openFiles[i].Close(); err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
	}
	second := p.timeMgr.NowTime().Unix()
	logDuration := p.getLogDuration(second)
	normalWriter, err := newNormalFileWriter(*p.options.absPath, *p.options.namePrefix, logDuration)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	errorWriter, err := newErrorFileWriter(*p.options.absPath, *p.options.namePrefix, logDuration)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
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
func Stop() error {
	if mgrInstance.logChan != nil {
		// close chan, for range 读完chan会退出.
		close(mgrInstance.logChan)
		// 等待logChan 的for range 退出.
		mgrInstance.waitGroupOutPut.Wait()
	}
	// 检查是否要关闭文件
	if len(mgrInstance.openFiles) > 0 {
		for i := range mgrInstance.openFiles {
			_ = mgrInstance.openFiles[i].Close()
		}
		mgrInstance.openFiles = mgrInstance.openFiles[0:0]
	}
	return nil
}

// fireHooks 处理钩子
func (p *mgr) fireHooks(entry *entry) {
	if 0 == len(p.options.hookMap) {
		return
	}
	err := p.options.hookMap.fire(entry)
	if err != nil {
		PrintfErr("failed to fire hook. err:%v", err)
	}
}

func (p *mgr) newEntry() *entry {
	entry := p.options.entryPoolOptions.newEntryFunc()
	return entry
}

// log 记录日志
func (p *mgr) log(entry *entry, level uint32, v ...interface{}) {
	withLevel(entry, level)
	withTime(entry, p.timeMgr.NowTime())
	withMessage(entry, fmt.Sprintln(v...))
	if *p.options.isReportCaller {
		pc, _, line, ok := runtime.Caller(calldepth2)
		funcName := xconstants.Unknown
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
func (p *mgr) logf(entry *entry, level uint32, format string, v ...interface{}) {
	withLevel(entry, level)
	withTime(entry, p.timeMgr.NowTime())
	withMessage(entry, fmt.Sprintf(format, v...))
	if *p.options.isReportCaller {
		pc, _, line, ok := runtime.Caller(calldepth2)
		funcName := xconstants.Unknown
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
func Trace(v ...interface{}) {
	if mgrInstance.GetLevel() < LevelTrace {
		return
	}
	mgrInstance.log(mgrInstance.newEntry(), LevelTrace, v...)
}

func TraceWithEntry(entry *entry, v ...interface{}) {
	if mgrInstance.GetLevel() < LevelTrace {
		return
	}
	mgrInstance.log(entry, LevelTrace, v...)
}

// Tracef 踪迹日志
func Tracef(format string, v ...interface{}) {
	if mgrInstance.GetLevel() < LevelTrace {
		return
	}
	mgrInstance.logf(mgrInstance.newEntry(), LevelTrace, format, v...)
}

func TracefWithEntry(entry *entry, format string, v ...interface{}) {
	if mgrInstance.GetLevel() < LevelTrace {
		return
	}
	mgrInstance.logf(entry, LevelTrace, format, v...)
}

// Debug 调试日志
func Debug(v ...interface{}) {
	if mgrInstance.GetLevel() < LevelDebug {
		return
	}
	mgrInstance.log(mgrInstance.newEntry(), LevelDebug, v...)
}

// DebugLazy 调试日志-惰性
//
//	等级满足之后才会计算
func DebugLazy(vFunc func() []interface{}) {
	if mgrInstance.GetLevel() < LevelDebug {
		return
	}
	v := vFunc()
	mgrInstance.log(mgrInstance.newEntry(), LevelDebug, v...)
}

// Debugf 调试日志
func Debugf(format string, v ...interface{}) {
	if mgrInstance.GetLevel() < LevelDebug {
		return
	}
	mgrInstance.logf(mgrInstance.newEntry(), LevelDebug, format, v...)
}

// Info 信息日志
func Info(v ...interface{}) {
	if mgrInstance.GetLevel() < LevelInfo {
		return
	}
	mgrInstance.log(mgrInstance.newEntry(), LevelInfo, v...)
}

// Infof 信息日志
func Infof(format string, v ...interface{}) {
	if mgrInstance.GetLevel() < LevelInfo {
		return
	}
	mgrInstance.logf(mgrInstance.newEntry(), LevelInfo, format, v...)
}

// Warn 警告日志
func Warn(v ...interface{}) {
	if mgrInstance.GetLevel() < LevelWarn {
		return
	}
	mgrInstance.log(mgrInstance.newEntry(), LevelWarn, v...)
}

// Warnf 警告日志
func Warnf(format string, v ...interface{}) {
	if mgrInstance.GetLevel() < LevelWarn {
		return
	}
	mgrInstance.logf(mgrInstance.newEntry(), LevelWarn, format, v...)
}

// Error 错误日志
func Error(v ...interface{}) {
	if mgrInstance.GetLevel() < LevelError {
		return
	}
	mgrInstance.log(mgrInstance.newEntry(), LevelError, v...)
}

// Errorf 错误日志
func Errorf(format string, v ...interface{}) {
	if mgrInstance.GetLevel() < LevelError {
		return
	}
	mgrInstance.logf(mgrInstance.newEntry(), LevelError, format, v...)
}

// Fatal 致命日志
func Fatal(v ...interface{}) {
	if mgrInstance.GetLevel() < LevelFatal {
		return
	}
	mgrInstance.log(mgrInstance.newEntry(), LevelFatal, v...)
}

// Fatalf 致命日志
func Fatalf(format string, v ...interface{}) {
	if mgrInstance.GetLevel() < LevelFatal {
		return
	}
	mgrInstance.logf(mgrInstance.newEntry(), LevelFatal, format, v...)
}
