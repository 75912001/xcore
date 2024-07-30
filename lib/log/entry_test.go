package log

import (
	"context"
	"fmt"
	"github.com/agiledragon/gomonkey"
	"github.com/pkg/errors"
	"log"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
	libtime "xcore/lib/time"
)

//go:generate go test -v -gcflags=all=-l -coverprofile=coverage.out

func Test_entry_reset(t *testing.T) {
	// 创建一个entry实例并设置一些初始值
	e := &entry{
		level:        1,
		time:         time.Now(),
		callerInfo:   "callerInfo",
		message:      "message",
		ctx:          context.Background(),
		extendFields: extendFields{"key", "value"},
	}

	// 调用reset方法
	reset(e)

	// 检查所有字段是否已经被重置为预期的默认值
	if e.level != LevelOff {
		t.Errorf("Expected level to be %d, but got %d", LevelOff, e.level)
	}
	if !e.time.IsZero() {
		//t.Errorf("Expected time to be zero, but got %v", e.time)
	}
	if e.callerInfo != "" {
		t.Errorf("Expected callerInfo to be empty, but got %s", e.callerInfo)
	}
	if e.message != "" {
		t.Errorf("Expected message to be empty, but got %s", e.message)
	}
	if e.ctx != nil {
		t.Errorf("Expected ctx to be nil, but got %v", e.ctx)
	}
	if len(e.extendFields) != 0 {
		t.Errorf("Expected extendFields to be empty, but got %v", e.extendFields)
	}
}

func Test_entry_withLevel(t *testing.T) {
	// 创建一个entry实例
	e := &entry{}

	// 设置一个预期的日志级别
	expectedLevel := LevelError

	// 调用withLevel方法
	withLevel(e, expectedLevel)

	// 检查日志级别是否已经被设置为预期的值
	if e.level != expectedLevel {
		t.Errorf("Expected level to be %d, but got %d", expectedLevel, e.level)
	}
}

func Test_entry_withTime(t *testing.T) {
	// 创建一个entry实例
	e := &entry{}

	// 设置一个预期的时间
	expectedTime := time.Now()

	// 调用withTime方法
	withTime(e, expectedTime)

	// 检查时间是否已经被设置为预期的值
	if !e.time.Equal(expectedTime) {
		t.Errorf("Expected time to be %v, but got %v", expectedTime, e.time)
	}
}

func Test_entry_withCallerInfo(t *testing.T) {
	// 创建一个entry实例
	e := &entry{}

	// 设置一个预期的调用者信息
	expectedCallerInfo := "callerInfo"

	// 调用withCallerInfo方法
	withCallerInfo(e, expectedCallerInfo)

	// 检查调用者信息是否已经被设置为预期的值
	if e.callerInfo != expectedCallerInfo {
		t.Errorf("Expected callerInfo to be %s, but got %s", expectedCallerInfo, e.callerInfo)
	}
}

func Test_entry_WithContext(t *testing.T) {
	// 创建一个entry实例
	e := &entry{}

	// 设置一个预期的上下文
	expectedCtx := context.WithValue(context.Background(), "key", "value")

	// 调用WithContext方法
	e.WithContext(expectedCtx)

	// 检查上下文是否已经被设置为预期的值
	if e.ctx != expectedCtx {
		t.Errorf("Expected ctx to be %v, but got %v", expectedCtx, e.ctx)
	}
}

func Test_entry_WithExtendField(t *testing.T) {
	// 创建一个entry实例
	e := &entry{}

	// 设置一个预期的扩展字段
	expectedKey := "key"
	expectedValue := "value"

	// 调用WithExtendField方法
	e.WithExtendField(expectedKey, expectedValue)

	// 检查扩展字段是否已经被设置为预期的值
	if len(e.extendFields) != 2 || e.extendFields[0] != expectedKey || e.extendFields[1] != expectedValue {
		t.Errorf("Expected extendFields to be [%s, %s], but got %v", expectedKey, expectedValue, e.extendFields)
	}
}

func Test_entry_WithExtendFields(t *testing.T) {
	// 创建一个entry实例
	e := &entry{}

	// 设置一个预期的扩展字段
	expectedFields := extendFields{"key1", "value1", "key2", "value2"}

	// 调用WithExtendFields方法
	e.WithExtendFields(expectedFields)

	// 检查扩展字段是否已经被设置为预期的值
	if len(e.extendFields) != len(expectedFields) {
		t.Errorf("Expected length of extendFields to be %d, but got %d", len(expectedFields), len(e.extendFields))
	} else {
		for i, v := range e.extendFields {
			if v != expectedFields[i] {
				t.Errorf("Expected extendFields to be %v, but got %v", expectedFields, e.extendFields)
				break
			}
		}
	}
}

func Test_entry_withMessage(t *testing.T) {
	// 创建一个entry实例
	e := &entry{}

	// 设置一个预期的消息
	expectedMessage := "test message"

	// 调用withMessage方法
	withMessage(e, expectedMessage)

	// 检查消息是否已经被设置为预期的值
	if e.message != expectedMessage {
		t.Errorf("Expected message to be %s, but got %s", expectedMessage, e.message)
	}
}

func Test_entry_formatLogData(t *testing.T) {
	var uid uint64 = 1001
	tests := []struct {
		name     string
		preFunc  func()
		postFunc func()
	}{
		{
			name: "normal",
			preFunc: func() {
				// 创建一个entry实例
				e := &entry{}

				// 设置预期的日志级别、时间、调用者信息、消息、上下文和扩展字段
				expectedLevel := LevelInfo
				expectedTime := time.Now()
				expectedCallerInfo := "callerInfo"
				expectedMessage := "test message"
				expectedCtx := context.WithValue(context.Background(), TraceIDKey, "TraceIDKey-Value")
				expectedCtx = context.WithValue(expectedCtx, UserIDKey, uid)
				expectedFields := extendFields{"key1", "value1", "key2", "value2", uid, "uid-value"}

				// 调用各种方法来设置日志级别、时间、调用者信息、消息、上下文和扩展字段
				withLevel(e, expectedLevel)
				withTime(e, expectedTime)
				withCallerInfo(e, expectedCallerInfo)
				withMessage(e, expectedMessage).
					WithContext(expectedCtx).
					WithExtendFields(expectedFields)

				// 调用formatLogData方法来格式化日志数据
				formattedLogData := formatLogData(e)

				// 检查返回的字符串是否符合预期
				expectedLogData := fmt.Sprintf("[%v][INF][%v:%v][%v:%v][callerInfo][{key1:value1}{key2:value2}{%v:uid-value}]test message",
					expectedTime.Format(logTimeFormat), TraceIDKey, "TraceIDKey-Value", UserIDKey, uid, uid)
				if formattedLogData != expectedLogData {
					t.Errorf("Expected log data to be %s, but got %s", expectedLogData, formattedLogData)
				}
			},
			postFunc: func() {
			},
		},
		{
			name: "normal",
			preFunc: func() {
				// 创建一个entry实例
				e := &entry{}

				// 设置预期的日志级别、时间、调用者信息、消息、上下文和扩展字段
				expectedLevel := LevelInfo
				expectedTime := time.Now()
				expectedCallerInfo := "callerInfo"
				expectedMessage := "test message"
				expectedCtx := context.WithValue(context.Background(), TraceIDKey, "TraceIDKey-Value")
				//expectedCtx = context.WithValue(expectedCtx, UserIDKey, uid)
				expectedFields := extendFields{"key1", "value1", UserIDKey, uid}

				// 调用各种方法来设置日志级别、时间、调用者信息、消息、上下文和扩展字段
				withLevel(e, expectedLevel)
				withTime(e, expectedTime)
				withCallerInfo(e, expectedCallerInfo)
				withMessage(e, expectedMessage)
				e.WithContext(expectedCtx).
					WithExtendFields(expectedFields)

				// 调用formatLogData方法来格式化日志数据
				formattedLogData := formatLogData(e)

				// 检查返回的字符串是否符合预期
				expectedLogData := fmt.Sprintf("[%v][INF][%v:%v][%v:%v][callerInfo][{key1:value1}{%v:%v}]test message",
					expectedTime.Format(logTimeFormat), TraceIDKey, "TraceIDKey-Value", UserIDKey, uid, UserIDKey, uid)
				if formattedLogData != expectedLogData {
					t.Errorf("Expected log data to be %s, but got %s", expectedLogData, formattedLogData)
				}
			},
			postFunc: func() {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.preFunc()
			tt.postFunc()
		})
	}
}

type myHook struct {
	fired   bool
	fireErr error
}

func (h *myHook) Levels() []int {
	return []int{LevelTrace, LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal}
}
func (h *myHook) Fire(entry *entry) error {
	fmt.Printf("hook fire %+v", entry)
	h.fired = true
	return h.fireErr
}

func TestLevelHooks_add(t *testing.T) {
	levelHooks := make(LevelHooks)
	myHook := new(myHook)
	levelHooks.add(myHook)
	if len(levelHooks) != 6 {
		t.Errorf("Expected 6 hook, but got %d", len(levelHooks))
	}
}

func TestLevelHooks_fire(t *testing.T) {
	levelHooks := make(LevelHooks)
	myHook := new(myHook)
	levelHooks.add(myHook)
	if len(levelHooks) != 6 {
		t.Errorf("Expected 6 hook, but got %d", len(levelHooks))
	}
	{ // return nil
		err := levelHooks.fire(
			&entry{
				level:        LevelInfo,
				time:         time.Now(),
				callerInfo:   "callerInfo",
				message:      "message",
				ctx:          nil,
				extendFields: nil,
			},
		)
		if err != nil {
			t.Errorf("fire() returned error: %v", err)
		}
		if !myHook.fired {
			t.Errorf("Expected hook to be fired, but it was not")
		}
	}
	{ // return error
		myHook.fireErr = fmt.Errorf("fire error")
		err := levelHooks.fire(
			&entry{
				level:        LevelInfo,
				time:         time.Now(),
				callerInfo:   "callerInfo",
				message:      "message",
				ctx:          nil,
				extendFields: nil,
			},
		)
		if err == nil {
			t.Error("fire() returned error: nil")
		}
	}
}

func TestNewMgr(t *testing.T) {
	type args struct {
		opts    []*options
		patches *gomonkey.Patches
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		preFunc  func(args *args)
		postFunc func(args *args)
	}{
		{
			name: "Test NewMgr with default options",
			args: args{
				//opts: []*options{NewOptions()},
			},
			wantErr: false,
		},
		{
			name: "Test NewMgr with error",
			args: args{
				//opts: []*options{NewOptions()},
			},
			wantErr: true,
			preFunc: func(args *args) {

				args.patches = gomonkey.ApplyFunc(withOptions, func(_ *mgr, _ ...*options) error {
					return errors.New("forced error")
				})
			},
			postFunc: func(args *args) {
				args.patches.Reset()
			},
		},
		{
			name: "Test NewMgr with error",
			args: args{
				//opts: []*options{NewOptions()},
			},
			wantErr: true,
			preFunc: func(args *args) {
				args.patches = gomonkey.ApplyFunc(start, func(_ *mgr) error {
					return errors.New("forced error")
				})
			},
			postFunc: func(args *args) {
				args.patches.Reset()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preFunc != nil {
				tt.preFunc(&tt.args)
			}
			_, err := NewMgr(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMgr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.postFunc != nil {
				tt.postFunc(&tt.args)
			}
		})
	}
}

func Test_withOptions(t *testing.T) {
	type args struct {
		p       *mgr
		opts    []*options
		patches *gomonkey.Patches
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		preFunc  func(args *args)
		postFunc func(args *args)
	}{
		{
			name:    "normal",
			wantErr: false,
			preFunc: func(args *args) {
				mgr, err := NewMgr()
				if err != nil {
					t.Fatal(err)
				}
				args.p = mgr
			},
		},
		{
			name:    "configure forced error",
			wantErr: true,
			preFunc: func(args *args) {
				mgr, err := NewMgr()
				if err != nil {
					t.Fatal(err)
				}
				args.p = mgr
				args.patches = gomonkey.ApplyFunc(configure, func(_ *options) error {
					return errors.New("forced error")
				})
			},
			postFunc: func(args *args) {
				args.patches.Reset()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preFunc != nil {
				tt.preFunc(&tt.args)
			}
			if err := withOptions(tt.args.p, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("withOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.postFunc != nil {
				tt.postFunc(&tt.args)
			}
		})
	}
}

func Test_mgr_start(t *testing.T) {
	type args struct {
		p       *mgr
		patches *gomonkey.Patches
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		preFunc  func(args *args)
		postFunc func(args *args)
	}{
		{
			name: "normal",
			preFunc: func(args *args) {
				element := new(mgr)
				args.p = element
				withOptions(element)
			},
			postFunc: func(args *args) {
				stdInstance = args.p
				args.p.Debug("this is debug log")
				args.p.Stop()
			},
		},
		{
			name: "normal-enablePool-false",
			preFunc: func(args *args) {
				element := new(mgr)
				args.p = element
				withOptions(element, NewOptions().WithEnablePool(false))
			},
			postFunc: func(args *args) {
				stdInstance = args.p
				args.p.Debug("this is debug log")
				args.p.Stop()
			},
		},
		{
			name:    "forced error - newWriters",
			wantErr: true,
			preFunc: func(args *args) {
				element := new(mgr)
				args.p = element
				withOptions(element)
				args.patches = gomonkey.ApplyFunc(newWriters, func(_ *mgr) error {
					return errors.New("forced error")
				})
			},
			postFunc: func(args *args) {
				stdInstance = args.p
				args.p.Stop()
				args.patches.Reset()
			},
		},
		{
			name: "forced panic - doLog",
			preFunc: func(args *args) {
				element := new(mgr)
				args.p = element
				withOptions(element)
				args.patches = gomonkey.ApplyFunc(doLog, func(_ *mgr) {
					args.p.logChan = nil
					panic("forced panic")
				})
			},
			postFunc: func(args *args) {
				//stdInstance = args.p
				args.p.Stop()
				args.patches.Reset()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preFunc != nil {
				tt.preFunc(&tt.args)
			}
			if err := start(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("start() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.postFunc != nil {
				tt.postFunc(&tt.args)
			}
		})
	}
}

func Test_mgr_GetLevel(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			if got := p.GetLevel(); got != tt.want {
				t.Errorf("GetLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mgr_getLogDuration(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	type args struct {
		sec int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			if got := p.getLogDuration(tt.args.sec); got != tt.want {
				t.Errorf("getLogDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mgr_doLog(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			doLog(p)
		})
	}
}

func Test_mgr_SetLevel(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	type args struct {
		level int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			if err := p.SetLevel(tt.args.level); (err != nil) != tt.wantErr {
				t.Errorf("SetLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mgr_newWriters(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			if err := newWriters(p); (err != nil) != tt.wantErr {
				t.Errorf("newWriters() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mgr_Stop(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			if err := p.Stop(); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mgr_fireHooks(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	type args struct {
		entry *entry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			p.fireHooks(tt.args.entry)
		})
	}
}

func Test_mgr_NewEntry(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	tests := []struct {
		name   string
		fields fields
		want   *entry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			if got := p.NewEntry(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mgr_log(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	type args struct {
		entry *entry
		level int
		v     []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			p.log(tt.args.entry, tt.args.level, tt.args.v...)
		})
	}
}

func Test_mgr_logf(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	type args struct {
		entry  *entry
		level  int
		format string
		v      []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			p.logf(tt.args.entry, tt.args.level, tt.args.format, tt.args.v...)
		})
	}
}

func Test_mgr_Trace(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	type args struct {
		v []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Trace(tt.args.v...)
		})
	}
}

func Test_mgr_TraceWithEntry(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	type args struct {
		entry *entry
		v     []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			p.TraceWithEntry(tt.args.entry, tt.args.v...)
		})
	}
}

func Test_mgr_Tracef(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	type args struct {
		format string
		v      []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Tracef(tt.args.format, tt.args.v...)
		})
	}
}

func Test_mgr_TracefWithEntry(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	type args struct {
		entry  *entry
		format string
		v      []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			p.TracefWithEntry(tt.args.entry, tt.args.format, tt.args.v...)
		})
	}
}

func Test_mgr_Debug(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	type args struct {
		v []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Debug(tt.args.v...)
		})
	}
}

func Test_mgr_Debugf(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	type args struct {
		format string
		v      []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Debugf(tt.args.format, tt.args.v...)
		})
	}
}

func Test_mgr_DebugLazy(t *testing.T) {
	// 创建一个mgr实例
	m, err := NewMgr()
	if err != nil {
		t.Fatalf("Failed to create mgr: %v", err)
	}

	// 设置日志级别为Debug
	err = m.SetLevel(LevelDebug)
	if err != nil {
		t.Fatalf("Failed to set log level: %v", err)
	}

	// 调用DebugLazy函数
	m.DebugLazy(func() []interface{} {
		// 这里是一些可能会消耗性能的计算
		// 这些计算只有在日志级别满足条件时才会执行
		return []interface{}{"debug message", "expensive computation result"}
	})
	m.DebugLazy(func() []interface{} {
		return []interface{}{fmt.Sprintf("%v:%v", "id", 1001)}
	})

	// 由于我们无法直接检查DebugLazy函数的输出，所以这个测试主要是为了确保DebugLazy函数可以被正确调用，而不会引发任何错误或异常。
	m.Stop()
}

func Test_mgr_Info(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	type args struct {
		v []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Info(tt.args.v...)
		})
	}
}

func Test_mgr_Infof(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	type args struct {
		format string
		v      []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Infof(tt.args.format, tt.args.v...)
		})
	}
}

func Test_mgr_Warn(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	type args struct {
		v []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Warn(tt.args.v...)
		})
	}
}

func Test_mgr_Warnf(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	type args struct {
		format string
		v      []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Warnf(tt.args.format, tt.args.v...)
		})
	}
}

func Test_mgr_Error(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	type args struct {
		v []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Error(tt.args.v...)
		})
	}
}

func Test_mgr_Errorf(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	type args struct {
		format string
		v      []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Errorf(tt.args.format, tt.args.v...)
		})
	}
}

func Test_mgr_Fatal(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	type args struct {
		v []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Fatal(tt.args.v...)
		})
	}
}

func Test_mgr_Fatalf(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		timeMgr         *libtime.Mgr
	}
	type args struct {
		format string
		v      []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mgr{
				options:         tt.fields.options,
				loggerSlice:     tt.fields.loggerSlice,
				logChan:         tt.fields.logChan,
				waitGroupOutPut: tt.fields.waitGroupOutPut,
				logDuration:     tt.fields.logDuration,
				openFiles:       tt.fields.openFiles,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Fatalf(tt.args.format, tt.args.v...)
		})
	}
}

func TestNewOptions(t *testing.T) {
	tests := []struct {
		name string
		want *options
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_options_WithLevel(t *testing.T) {
	type fields struct {
		level            *int
		absPath          *string
		isReportCaller   *bool
		namePrefix       *string
		isWriteFile      *bool
		entryPoolOptions *entryPoolOptions
		hooks            LevelHooks
	}
	type args struct {
		level int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *options
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &options{
				level:            tt.fields.level,
				absPath:          tt.fields.absPath,
				isReportCaller:   tt.fields.isReportCaller,
				namePrefix:       tt.fields.namePrefix,
				isWriteFile:      tt.fields.isWriteFile,
				entryPoolOptions: tt.fields.entryPoolOptions,
				hooks:            tt.fields.hooks,
			}
			if got := p.WithLevel(tt.args.level); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_options_WithAbsPath(t *testing.T) {
	type fields struct {
		level            *int
		absPath          *string
		isReportCaller   *bool
		namePrefix       *string
		isWriteFile      *bool
		entryPoolOptions *entryPoolOptions
		hooks            LevelHooks
	}
	type args struct {
		absPath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *options
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &options{
				level:            tt.fields.level,
				absPath:          tt.fields.absPath,
				isReportCaller:   tt.fields.isReportCaller,
				namePrefix:       tt.fields.namePrefix,
				isWriteFile:      tt.fields.isWriteFile,
				entryPoolOptions: tt.fields.entryPoolOptions,
				hooks:            tt.fields.hooks,
			}
			if got := p.WithAbsPath(tt.args.absPath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithAbsPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_options_WithIsReportCaller(t *testing.T) {
	type fields struct {
		level            *int
		absPath          *string
		isReportCaller   *bool
		namePrefix       *string
		isWriteFile      *bool
		entryPoolOptions *entryPoolOptions
		hooks            LevelHooks
	}
	type args struct {
		isReportCaller bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *options
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &options{
				level:            tt.fields.level,
				absPath:          tt.fields.absPath,
				isReportCaller:   tt.fields.isReportCaller,
				namePrefix:       tt.fields.namePrefix,
				isWriteFile:      tt.fields.isWriteFile,
				entryPoolOptions: tt.fields.entryPoolOptions,
				hooks:            tt.fields.hooks,
			}
			if got := p.WithIsReportCaller(tt.args.isReportCaller); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithIsReportCaller() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_options_WithNamePrefix(t *testing.T) {
	type fields struct {
		level            *int
		absPath          *string
		isReportCaller   *bool
		namePrefix       *string
		isWriteFile      *bool
		entryPoolOptions *entryPoolOptions
		hooks            LevelHooks
	}
	type args struct {
		namePrefix string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *options
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &options{
				level:            tt.fields.level,
				absPath:          tt.fields.absPath,
				isReportCaller:   tt.fields.isReportCaller,
				namePrefix:       tt.fields.namePrefix,
				isWriteFile:      tt.fields.isWriteFile,
				entryPoolOptions: tt.fields.entryPoolOptions,
				hooks:            tt.fields.hooks,
			}
			if got := p.WithNamePrefix(tt.args.namePrefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithNamePrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_options_WithHooks(t *testing.T) {
	type fields struct {
		level            *int
		absPath          *string
		isReportCaller   *bool
		namePrefix       *string
		isWriteFile      *bool
		entryPoolOptions *entryPoolOptions
		hooks            LevelHooks
	}
	type args struct {
		hooks LevelHooks
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *options
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &options{
				level:            tt.fields.level,
				absPath:          tt.fields.absPath,
				isReportCaller:   tt.fields.isReportCaller,
				namePrefix:       tt.fields.namePrefix,
				isWriteFile:      tt.fields.isWriteFile,
				entryPoolOptions: tt.fields.entryPoolOptions,
				hooks:            tt.fields.hooks,
			}
			if got := p.WithHooks(tt.args.hooks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithHooks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_options_WithIsWriteFile(t *testing.T) {
	type fields struct {
		level            *int
		absPath          *string
		isReportCaller   *bool
		namePrefix       *string
		isWriteFile      *bool
		entryPoolOptions *entryPoolOptions
		hooks            LevelHooks
	}
	type args struct {
		isWriteFile bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *options
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &options{
				level:            tt.fields.level,
				absPath:          tt.fields.absPath,
				isReportCaller:   tt.fields.isReportCaller,
				namePrefix:       tt.fields.namePrefix,
				isWriteFile:      tt.fields.isWriteFile,
				entryPoolOptions: tt.fields.entryPoolOptions,
				hooks:            tt.fields.hooks,
			}
			if got := p.WithIsWriteFile(tt.args.isWriteFile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithIsWriteFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_options_WithEnablePool(t *testing.T) {
	type fields struct {
		level            *int
		absPath          *string
		isReportCaller   *bool
		namePrefix       *string
		isWriteFile      *bool
		entryPoolOptions *entryPoolOptions
		hooks            LevelHooks
	}
	type args struct {
		enablePool bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *options
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &options{
				level:            tt.fields.level,
				absPath:          tt.fields.absPath,
				isReportCaller:   tt.fields.isReportCaller,
				namePrefix:       tt.fields.namePrefix,
				isWriteFile:      tt.fields.isWriteFile,
				entryPoolOptions: tt.fields.entryPoolOptions,
				hooks:            tt.fields.hooks,
			}
			if got := p.WithEnablePool(tt.args.enablePool); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithEnablePool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_options_IsEnablePool(t *testing.T) {
	type fields struct {
		level            *int
		absPath          *string
		isReportCaller   *bool
		namePrefix       *string
		isWriteFile      *bool
		entryPoolOptions *entryPoolOptions
		hooks            LevelHooks
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &options{
				level:            tt.fields.level,
				absPath:          tt.fields.absPath,
				isReportCaller:   tt.fields.isReportCaller,
				namePrefix:       tt.fields.namePrefix,
				isWriteFile:      tt.fields.isWriteFile,
				entryPoolOptions: tt.fields.entryPoolOptions,
				hooks:            tt.fields.hooks,
			}
			if got := p.IsEnablePool(); got != tt.want {
				t.Errorf("IsEnablePool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_options_AddHooks(t *testing.T) {
	type fields struct {
		level            *int
		absPath          *string
		isReportCaller   *bool
		namePrefix       *string
		isWriteFile      *bool
		entryPoolOptions *entryPoolOptions
		hooks            LevelHooks
	}
	type args struct {
		hook Hook
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *options
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &options{
				level:            tt.fields.level,
				absPath:          tt.fields.absPath,
				isReportCaller:   tt.fields.isReportCaller,
				namePrefix:       tt.fields.namePrefix,
				isWriteFile:      tt.fields.isWriteFile,
				entryPoolOptions: tt.fields.entryPoolOptions,
				hooks:            tt.fields.hooks,
			}
			if got := p.AddHooks(tt.args.hook); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddHooks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeOptions(t *testing.T) {
	type args struct {
		opts []*options
	}
	tests := []struct {
		name string
		args args
		want *options
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeOptions(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configure(t *testing.T) {
	type args struct {
		opts *options
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := configure(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("configure() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPrintErr(t *testing.T) {
	type args struct {
		v []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintErr(tt.args.v...)
		})
	}
}

func TestPrintfErr(t *testing.T) {
	type args struct {
		format string
		v      []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintfErr(tt.args.format, tt.args.v...)
		})
	}
}

func Test_isEnable(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEnable(); got != tt.want {
				t.Errorf("isEnable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrintInfo(t *testing.T) {
	type args struct {
		v []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintInfo(tt.args.v...)
		})
	}
}

func TestPrintfInfo(t *testing.T) {
	type args struct {
		format string
		v      []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintfInfo(tt.args.format, tt.args.v...)
		})
	}
}

func Test_formatAndPrint(t *testing.T) {
	type args struct {
		level    int
		line     int
		funcName string
		v        []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatAndPrint(tt.args.level, tt.args.line, tt.args.funcName, tt.args.v...)
		})
	}
}

func Test_newNormalFileWriter(t *testing.T) {
	type args struct {
		filePath    string
		namePrefix  string
		logDuration int
	}
	tests := []struct {
		name    string
		args    args
		want    *os.File
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newNormalFileWriter(tt.args.filePath, tt.args.namePrefix, tt.args.logDuration)
			if (err != nil) != tt.wantErr {
				t.Errorf("newNormalFileWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newNormalFileWriter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newErrorFileWriter(t *testing.T) {
	type args struct {
		filePath    string
		namePrefix  string
		logDuration int
	}
	tests := []struct {
		name    string
		args    args
		want    *os.File
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newErrorFileWriter(tt.args.filePath, tt.args.namePrefix, tt.args.logDuration)
			if (err != nil) != tt.wantErr {
				t.Errorf("newErrorFileWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newErrorFileWriter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newFileWriter(t *testing.T) {
	type args struct {
		filePath     string
		namePrefix   string
		logDuration  int
		fileBaseName string
	}
	tests := []struct {
		name    string
		args    args
		want    *os.File
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newFileWriter(tt.args.filePath, tt.args.namePrefix, tt.args.logDuration, tt.args.fileBaseName)
			if (err != nil) != tt.wantErr {
				t.Errorf("newFileWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newFileWriter() got = %v, want %v", got, tt.want)
			}
		})
	}
}
