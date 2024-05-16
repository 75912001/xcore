package log

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
	libtime "xcore/lib/time"
)

// todo menglc 完成覆盖率 [100%] 测试

func TestGetInstance(t *testing.T) {
	instance1 := GetInstance()
	instance2 := GetInstance()

	if instance1 != instance2 {
		t.Errorf("GetInstance() returned different instances")
	}
}

func TestIsEnable(t *testing.T) {
	// 当 instance 为 nil 时，isEnable 应返回 false
	instance = nil
	if isEnable() != false {
		t.Errorf("Expected isEnable to return false, but it returned true")
	}
	// 当 instance 不为 nil 时， instance.logChan == nil, 此时 isEnable 应返回 false
	instance = &mgr{}
	if isEnable() != false {
		t.Errorf("Expected isEnable to return false, but it returned true")
	}
	// 启动, 使 instance.logChan 不为 nil, 此时 isEnable 应返回 true
	err := instance.start()
	if err != nil {
		t.Errorf("start() returned error: %v", err)
	}
	if isEnable() != true {
		t.Errorf("Expected isEnable to return true, but it returned false")
	}

	// Reset instance to nil after the test
	instance = nil
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

func TestNewOptions(t *testing.T) {
	NewOptions()
}

func TestPrintErr(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		preFunc  func()
		postFunc func()
	}{
		{
			name:  "日志启用",
			input: []interface{}{"test"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelFatal)) // LevelFatal is less than LevelError
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
		{
			name:  "日志-未启用",
			input: []interface{}{"test"},
			preFunc: func() {
				GetInstance()
				instance = nil
			},
			postFunc: func() {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.preFunc()
			PrintErr(tt.input...)
			tt.postFunc()
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

func Test_entry_Debug(t *testing.T) {
	tests := []struct {
		name         string
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
		preFunc      func()
		postFunc     func()
	}{
		{
			name:         "normal-日志等级 < 等级",
			level:        LevelDebug,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelInfo))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
		{
			name:         "normal- 等级 <= 日志等级",
			level:        LevelDebug,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelDebug))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &entry{
				level:        tt.level,
				time:         tt.time,
				callerInfo:   tt.callerInfo,
				message:      tt.message,
				ctx:          tt.ctx,
				extendFields: tt.extendFields,
			}
			tt.preFunc()
			p.Debug("this is message:", 123, "xxx", p)
			tt.postFunc()
		})
	}
}

func Test_entry_Debugf(t *testing.T) {
	tests := []struct {
		name         string
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
		preFunc      func()
		postFunc     func()
	}{
		{
			name:         "normal-日志等级 < 等级",
			level:        LevelDebug,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelInfo))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
		{
			name:         "normal- 等级 <= 日志等级",
			level:        LevelDebug,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelDebug))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &entry{
				level:        tt.level,
				time:         tt.time,
				callerInfo:   tt.callerInfo,
				message:      tt.message,
				ctx:          tt.ctx,
				extendFields: tt.extendFields,
			}
			tt.preFunc()
			p.Debugf("this is message:%v", 123)
			tt.postFunc()
		})
	}
}

func Test_entry_Error(t *testing.T) {
	tests := []struct {
		name         string
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
		preFunc      func()
		postFunc     func()
	}{
		{
			name:         "normal-日志等级 < trace 等级",
			level:        LevelError,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelFatal))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
		{
			name:         "normal- trace 等级 <= 日志等级",
			level:        LevelError,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelError))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &entry{
				level:        tt.level,
				time:         tt.time,
				callerInfo:   tt.callerInfo,
				message:      tt.message,
				ctx:          tt.ctx,
				extendFields: tt.extendFields,
			}
			tt.preFunc()
			p.Error("this is message:", 123)
			tt.postFunc()
		})
	}
}

func Test_entry_Errorf(t *testing.T) {
	tests := []struct {
		name         string
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
		preFunc      func()
		postFunc     func()
	}{
		{
			name:         "normal-日志等级 < trace 等级",
			level:        LevelError,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelFatal))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
		{
			name:         "normal- trace 等级 <= 日志等级",
			level:        LevelError,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelError))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &entry{
				level:        tt.level,
				time:         tt.time,
				callerInfo:   tt.callerInfo,
				message:      tt.message,
				ctx:          tt.ctx,
				extendFields: tt.extendFields,
			}
			tt.preFunc()
			p.Errorf("this is message:%v", 123)
			tt.postFunc()
		})
	}
}

func Test_entry_Fatal(t *testing.T) {
	tests := []struct {
		name         string
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
		preFunc      func()
		postFunc     func()
	}{
		{
			name:         "normal-日志等级 < trace 等级",
			level:        LevelFatal,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelOff))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
		{
			name:         "normal- trace 等级 <= 日志等级",
			level:        LevelFatal,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelFatal))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &entry{
				level:        tt.level,
				time:         tt.time,
				callerInfo:   tt.callerInfo,
				message:      tt.message,
				ctx:          tt.ctx,
				extendFields: tt.extendFields,
			}
			tt.preFunc()
			p.Fatal("this is message:", 123)
			tt.postFunc()
		})
	}
}

func Test_entry_Fatalf(t *testing.T) {
	tests := []struct {
		name         string
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
		preFunc      func()
		postFunc     func()
	}{
		{
			name:         "normal-日志等级 < trace 等级",
			level:        LevelFatal,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelOff))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
		{
			name:         "normal- trace 等级 <= 日志等级",
			level:        LevelFatal,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelFatal))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &entry{
				level:        tt.level,
				time:         tt.time,
				callerInfo:   tt.callerInfo,
				message:      tt.message,
				ctx:          tt.ctx,
				extendFields: tt.extendFields,
			}
			tt.preFunc()
			p.Fatalf("this is message:%v", 123)
			tt.postFunc()
		})
	}
}

func Test_entry_Info(t *testing.T) {
	tests := []struct {
		name         string
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
		preFunc      func()
		postFunc     func()
	}{
		{
			name:         "normal-日志等级 < 等级",
			level:        LevelInfo,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelWarn))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
		{
			name:         "normal- 等级 <= 日志等级",
			level:        LevelInfo,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelInfo))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &entry{
				level:        tt.level,
				time:         tt.time,
				callerInfo:   tt.callerInfo,
				message:      tt.message,
				ctx:          tt.ctx,
				extendFields: tt.extendFields,
			}
			tt.preFunc()
			p.Info("this is message:", 123, "xxx", p)
			tt.postFunc()
		})
	}
}

func Test_entry_Infof(t *testing.T) {
	tests := []struct {
		name         string
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
		preFunc      func()
		postFunc     func()
	}{
		{
			name:         "normal-日志等级 < 等级",
			level:        LevelInfo,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelWarn))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
		{
			name:         "normal- 等级 <= 日志等级",
			level:        LevelInfo,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelInfo))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &entry{
				level:        tt.level,
				time:         tt.time,
				callerInfo:   tt.callerInfo,
				message:      tt.message,
				ctx:          tt.ctx,
				extendFields: tt.extendFields,
			}
			tt.preFunc()
			p.Infof("this is message:%v", 123)
			tt.postFunc()
		})
	}
}

func Test_entry_Trace(t *testing.T) {
	tests := []struct {
		name         string
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
		preFunc      func()
		postFunc     func()
	}{
		{
			name:         "normal-日志等级 < trace 等级",
			level:        LevelTrace,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelDebug))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
		{
			name:         "normal- trace 等级 <= 日志等级",
			level:        LevelTrace,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelTrace))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &entry{
				level:        tt.level,
				time:         tt.time,
				callerInfo:   tt.callerInfo,
				message:      tt.message,
				ctx:          tt.ctx,
				extendFields: tt.extendFields,
			}
			tt.preFunc()
			p.Trace("this is message:", 123, "xxx", p)
			tt.postFunc()
		})
	}
}

func Test_entry_Tracef(t *testing.T) {
	tests := []struct {
		name         string
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
		preFunc      func()
		postFunc     func()
	}{
		{
			name:         "normal-日志等级 < trace 等级",
			level:        LevelTrace,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelDebug))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
		{
			name:         "normal- trace 等级 <= 日志等级",
			level:        LevelTrace,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelTrace))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &entry{
				level:        tt.level,
				time:         tt.time,
				callerInfo:   tt.callerInfo,
				message:      tt.message,
				ctx:          tt.ctx,
				extendFields: tt.extendFields,
			}
			tt.preFunc()
			p.Tracef("this is message:%v", 123)
			tt.postFunc()
		})
	}
}

func Test_entry_Warn(t *testing.T) {
	tests := []struct {
		name         string
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
		preFunc      func()
		postFunc     func()
	}{
		{
			name:         "normal-日志等级 < trace 等级",
			level:        LevelWarn,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelError))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
		{
			name:         "normal- trace 等级 <= 日志等级",
			level:        LevelWarn,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelWarn))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &entry{
				level:        tt.level,
				time:         tt.time,
				callerInfo:   tt.callerInfo,
				message:      tt.message,
				ctx:          tt.ctx,
				extendFields: tt.extendFields,
			}
			tt.preFunc()
			p.Warn("this is message:", 123)
			tt.postFunc()
		})
	}
}

func Test_entry_Warnf(t *testing.T) {
	tests := []struct {
		name         string
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
		preFunc      func()
		postFunc     func()
	}{
		{
			name:         "normal-日志等级 < trace 等级",
			level:        LevelWarn,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelError))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
		{
			name:         "normal- trace 等级 <= 日志等级",
			level:        LevelWarn,
			time:         time.Now(),
			callerInfo:   "callerInfo",
			message:      "normal-LevelFatal-message",
			ctx:          context.TODO(),
			extendFields: extendFields{"key-1", "value-1", "key-2", "value-2"},
			preFunc: func() {
				GetInstance()
				instance = &mgr{}
				instance.start(NewOptions().WithLevel(LevelWarn))
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &entry{
				level:        tt.level,
				time:         tt.time,
				callerInfo:   tt.callerInfo,
				message:      tt.message,
				ctx:          tt.ctx,
				extendFields: tt.extendFields,
			}
			tt.preFunc()
			p.Warnf("this is message:%v", 123)
			tt.postFunc()
		})
	}
}

func Test_entry_WithContext(t *testing.T) {
	instance = new(mgr)
	_ = GetInstance().Start()
	_ = newEntry().withContext(context.Background())
	GetInstance().Stop()
	return
}

func Test_entry_WithExtendField(t *testing.T) {
	instance = new(mgr)
	_ = GetInstance().Start()
	_ = newEntry().withExtendField("key", "value")
	GetInstance().Stop()
	return
}

func Test_entry_WithExtendFields(t *testing.T) {
	instance = new(mgr)
	_ = GetInstance().Start()
	fields := extendFields{"key-1", "value-1", "key-2", "value-2"}
	_ = newEntry().withExtendFields(fields)
	GetInstance().Stop()
	return
}

func Test_entry_WithMessage(t *testing.T) {
	instance = new(mgr)
	_ = GetInstance().Start()
	_ = newEntry().withMessage("message")
	GetInstance().Stop()
	return
}

func Test_entry_formatMessage(t *testing.T) {
	tests := []struct {
		name     string
		preFunc  func()
		postFunc func()
	}{
		{
			name: "normal",
			preFunc: func() {
				instance = new(mgr)
				_ = GetInstance().Start()
				ctx := context.Background()
				ctx = context.WithValue(ctx, TraceIDKey, "traceID-1001")
				ctx = context.WithValue(ctx, UserIDKey, 1001)
				_ = newEntry().withContext(ctx).formatMessage()
			},
			postFunc: func() {
				GetInstance().Stop()
			},
		},
		{
			name: "normal",
			preFunc: func() {
				instance = new(mgr)
				_ = GetInstance().Start()
				ctx := context.Background()
				ctx = context.WithValue(ctx, TraceIDKey, "traceID-1001")
				_ = newEntry().withContext(ctx).
					withExtendFields(extendFields{
						UserIDKey, 1001,
						"field-1-key", "field-1-value",
						2001, "field-3-value",
						"field-2-key", 9002,
					}).
					formatMessage()
			},
			postFunc: func() {
				GetInstance().Stop()
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

func Test_entry_log(t *testing.T) {
	type fields struct {
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
	}
	type args struct {
		level int
		skip  int
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
			p := &entry{
				level:        tt.fields.level,
				time:         tt.fields.time,
				callerInfo:   tt.fields.callerInfo,
				message:      tt.fields.message,
				ctx:          tt.fields.ctx,
				extendFields: tt.fields.extendFields,
			}
			p.log(tt.args.level, tt.args.skip, tt.args.v...)
		})
	}
}

func Test_entry_logf(t *testing.T) {
	type fields struct {
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
	}
	type args struct {
		level  int
		skip   int
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
			p := &entry{
				level:        tt.fields.level,
				time:         tt.fields.time,
				callerInfo:   tt.fields.callerInfo,
				message:      tt.fields.message,
				ctx:          tt.fields.ctx,
				extendFields: tt.fields.extendFields,
			}
			p.logf(tt.args.level, tt.args.skip, tt.args.format, tt.args.v...)
		})
	}
}

func Test_entry_reset(t *testing.T) {
	type fields struct {
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &entry{
				level:        tt.fields.level,
				time:         tt.fields.time,
				callerInfo:   tt.fields.callerInfo,
				message:      tt.fields.message,
				ctx:          tt.fields.ctx,
				extendFields: tt.fields.extendFields,
			}
			p.reset()
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

func Test_mgr_Debug(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		pool            *sync.Pool
		timeMgr         *libtime.Mgr
	}
	type args struct {
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Debug(tt.args.level, tt.args.v...)
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
		pool            *sync.Pool
		timeMgr         *libtime.Mgr
	}
	type args struct {
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Debugf(tt.args.level, tt.args.format, tt.args.v...)
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
		pool            *sync.Pool
		timeMgr         *libtime.Mgr
	}
	type args struct {
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Error(tt.args.level, tt.args.v...)
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
		pool            *sync.Pool
		timeMgr         *libtime.Mgr
	}
	type args struct {
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Errorf(tt.args.level, tt.args.format, tt.args.v...)
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
		pool            *sync.Pool
		timeMgr         *libtime.Mgr
	}
	type args struct {
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Fatal(tt.args.level, tt.args.v...)
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
		pool            *sync.Pool
		timeMgr         *libtime.Mgr
	}
	type args struct {
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Fatalf(tt.args.level, tt.args.format, tt.args.v...)
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
		pool            *sync.Pool
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			if got := p.GetLevel(); got != tt.want {
				t.Errorf("GetLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mgr_Info(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		pool            *sync.Pool
		timeMgr         *libtime.Mgr
	}
	type args struct {
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Info(tt.args.level, tt.args.v...)
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
		pool            *sync.Pool
		timeMgr         *libtime.Mgr
	}
	type args struct {
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Infof(tt.args.level, tt.args.format, tt.args.v...)
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
		pool            *sync.Pool
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			if err := p.SetLevel(tt.args.level); (err != nil) != tt.wantErr {
				t.Errorf("SetLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mgr_Start(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		pool            *sync.Pool
		timeMgr         *libtime.Mgr
	}
	type args struct {
		opts []*options
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			if err := p.start(tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("start() error = %v, wantErr %v", err, tt.wantErr)
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
		pool            *sync.Pool
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			if err := p.Stop(); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		pool            *sync.Pool
		timeMgr         *libtime.Mgr
	}
	type args struct {
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Trace(tt.args.level, tt.args.v...)
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
		pool            *sync.Pool
		timeMgr         *libtime.Mgr
	}
	type args struct {
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Tracef(tt.args.level, tt.args.format, tt.args.v...)
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
		pool            *sync.Pool
		timeMgr         *libtime.Mgr
	}
	type args struct {
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Warn(tt.args.level, tt.args.v...)
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
		pool            *sync.Pool
		timeMgr         *libtime.Mgr
	}
	type args struct {
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			p.Warnf(tt.args.level, tt.args.format, tt.args.v...)
		})
	}
}

func Test_mgr_WithContext(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		pool            *sync.Pool
		timeMgr         *libtime.Mgr
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			if got := p.WithContext(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("withContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mgr_WithField(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		pool            *sync.Pool
		timeMgr         *libtime.Mgr
	}
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			if got := p.WithField(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mgr_WithFields(t *testing.T) {
	type fields struct {
		options         *options
		loggerSlice     [LevelOn]*log.Logger
		logChan         chan *entry
		waitGroupOutPut sync.WaitGroup
		logDuration     int
		openFiles       []*os.File
		pool            *sync.Pool
		timeMgr         *libtime.Mgr
	}
	type args struct {
		f extendFields
	}
	tests := []struct {
		name   string
		fields fields
		args   args
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			if got := p.WithFields(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithFields() = %v, want %v", got, tt.want)
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
		pool            *sync.Pool
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			p.doLog()
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
		pool            *sync.Pool
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			p.fireHooks(tt.args.entry)
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
		pool            *sync.Pool
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			if got := p.getLogDuration(tt.args.sec); got != tt.want {
				t.Errorf("getLogDuration() = %v, want %v", got, tt.want)
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
		pool            *sync.Pool
		timeMgr         *libtime.Mgr
	}
	type args struct {
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			p.log(tt.args.level, tt.args.v...)
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
		pool            *sync.Pool
		timeMgr         *libtime.Mgr
	}
	type args struct {
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			p.logf(tt.args.level, tt.args.format, tt.args.v...)
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
		pool            *sync.Pool
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
				pool:            tt.fields.pool,
				timeMgr:         tt.fields.timeMgr,
			}
			if err := p.newWriters(); (err != nil) != tt.wantErr {
				t.Errorf("newWriters() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newEntry(t *testing.T) {
	tests := []struct {
		name string
		want *entry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newEntry(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEntry() = %v, want %v", got, tt.want)
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

func Test_options_AddHooks(t *testing.T) {
	type fields struct {
		level          *int
		absPath        *string
		isReportCaller *bool
		namePrefix     *string
		isWriteFile    *bool
		enablePool     *bool
		hooks          LevelHooks
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
				level:          tt.fields.level,
				absPath:        tt.fields.absPath,
				isReportCaller: tt.fields.isReportCaller,
				namePrefix:     tt.fields.namePrefix,
				isWriteFile:    tt.fields.isWriteFile,
				enablePool:     tt.fields.enablePool,
				hooks:          tt.fields.hooks,
			}
			if got := p.AddHooks(tt.args.hook); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddHooks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_options_IsEnablePool(t *testing.T) {
	type fields struct {
		level          *int
		absPath        *string
		isReportCaller *bool
		namePrefix     *string
		isWriteFile    *bool
		enablePool     *bool
		hooks          LevelHooks
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
				level:          tt.fields.level,
				absPath:        tt.fields.absPath,
				isReportCaller: tt.fields.isReportCaller,
				namePrefix:     tt.fields.namePrefix,
				isWriteFile:    tt.fields.isWriteFile,
				enablePool:     tt.fields.enablePool,
				hooks:          tt.fields.hooks,
			}
			if got := p.IsEnablePool(); got != tt.want {
				t.Errorf("IsEnablePool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_options_WithAbsPath(t *testing.T) {
	type fields struct {
		level          *int
		absPath        *string
		isReportCaller *bool
		namePrefix     *string
		isWriteFile    *bool
		enablePool     *bool
		hooks          LevelHooks
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
				level:          tt.fields.level,
				absPath:        tt.fields.absPath,
				isReportCaller: tt.fields.isReportCaller,
				namePrefix:     tt.fields.namePrefix,
				isWriteFile:    tt.fields.isWriteFile,
				enablePool:     tt.fields.enablePool,
				hooks:          tt.fields.hooks,
			}
			if got := p.WithAbsPath(tt.args.absPath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithAbsPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_options_WithEnablePool(t *testing.T) {
	type fields struct {
		level          *int
		absPath        *string
		isReportCaller *bool
		namePrefix     *string
		isWriteFile    *bool
		enablePool     *bool
		hooks          LevelHooks
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
				level:          tt.fields.level,
				absPath:        tt.fields.absPath,
				isReportCaller: tt.fields.isReportCaller,
				namePrefix:     tt.fields.namePrefix,
				isWriteFile:    tt.fields.isWriteFile,
				enablePool:     tt.fields.enablePool,
				hooks:          tt.fields.hooks,
			}
			if got := p.WithEnablePool(tt.args.enablePool); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithEnablePool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_options_WithHooks(t *testing.T) {
	type fields struct {
		level          *int
		absPath        *string
		isReportCaller *bool
		namePrefix     *string
		isWriteFile    *bool
		enablePool     *bool
		hooks          LevelHooks
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
				level:          tt.fields.level,
				absPath:        tt.fields.absPath,
				isReportCaller: tt.fields.isReportCaller,
				namePrefix:     tt.fields.namePrefix,
				isWriteFile:    tt.fields.isWriteFile,
				enablePool:     tt.fields.enablePool,
				hooks:          tt.fields.hooks,
			}
			if got := p.WithHooks(tt.args.hooks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithHooks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_options_WithIsReportCaller(t *testing.T) {
	type fields struct {
		level          *int
		absPath        *string
		isReportCaller *bool
		namePrefix     *string
		isWriteFile    *bool
		enablePool     *bool
		hooks          LevelHooks
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
				level:          tt.fields.level,
				absPath:        tt.fields.absPath,
				isReportCaller: tt.fields.isReportCaller,
				namePrefix:     tt.fields.namePrefix,
				isWriteFile:    tt.fields.isWriteFile,
				enablePool:     tt.fields.enablePool,
				hooks:          tt.fields.hooks,
			}
			if got := p.WithIsReportCaller(tt.args.isReportCaller); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithIsReportCaller() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_options_WithIsWriteFile(t *testing.T) {
	type fields struct {
		level          *int
		absPath        *string
		isReportCaller *bool
		namePrefix     *string
		isWriteFile    *bool
		enablePool     *bool
		hooks          LevelHooks
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
				level:          tt.fields.level,
				absPath:        tt.fields.absPath,
				isReportCaller: tt.fields.isReportCaller,
				namePrefix:     tt.fields.namePrefix,
				isWriteFile:    tt.fields.isWriteFile,
				enablePool:     tt.fields.enablePool,
				hooks:          tt.fields.hooks,
			}
			if got := p.WithIsWriteFile(tt.args.isWriteFile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithIsWriteFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_options_WithLevel(t *testing.T) {
	type fields struct {
		level          *int
		absPath        *string
		isReportCaller *bool
		namePrefix     *string
		isWriteFile    *bool
		enablePool     *bool
		hooks          LevelHooks
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
				level:          tt.fields.level,
				absPath:        tt.fields.absPath,
				isReportCaller: tt.fields.isReportCaller,
				namePrefix:     tt.fields.namePrefix,
				isWriteFile:    tt.fields.isWriteFile,
				enablePool:     tt.fields.enablePool,
				hooks:          tt.fields.hooks,
			}
			if got := p.WithLevel(tt.args.level); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_options_WithNamePrefix(t *testing.T) {
	type fields struct {
		level          *int
		absPath        *string
		isReportCaller *bool
		namePrefix     *string
		isWriteFile    *bool
		enablePool     *bool
		hooks          LevelHooks
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
				level:          tt.fields.level,
				absPath:        tt.fields.absPath,
				isReportCaller: tt.fields.isReportCaller,
				namePrefix:     tt.fields.namePrefix,
				isWriteFile:    tt.fields.isWriteFile,
				enablePool:     tt.fields.enablePool,
				hooks:          tt.fields.hooks,
			}
			if got := p.WithNamePrefix(tt.args.namePrefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithNamePrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
