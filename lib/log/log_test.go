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

func TestGetInstance(t *testing.T) {
	instance1 := GetInstance()
	instance2 := GetInstance()

	if instance1 != instance2 {
		t.Errorf("GetInstance() returned different instances")
	}
}

func TestIsEnable(t *testing.T) {
	// 当 instance 为 nil 时，IsEnable 应返回 false
	instance = nil
	if IsEnable() != false {
		t.Errorf("Expected IsEnable to return false, but it returned true")
	}
	// 当 instance 不为 nil 时， instance.logChan == nil, 此时 IsEnable 应返回 false
	instance = &mgr{}
	if IsEnable() != false {
		t.Errorf("Expected IsEnable to return false, but it returned true")
	}
	// 启动, 使 instance.logChan 不为 nil, 此时 IsEnable 应返回 true
	err := instance.Start()
	if err != nil {
		t.Errorf("Start() returned error: %v", err)
	}
	if IsEnable() != true {
		t.Errorf("Expected IsEnable to return true, but it returned false")
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
	type fields struct {
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
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
			p := &entry{
				level:        tt.fields.level,
				time:         tt.fields.time,
				callerInfo:   tt.fields.callerInfo,
				message:      tt.fields.message,
				ctx:          tt.fields.ctx,
				extendFields: tt.fields.extendFields,
			}
			p.Debug(tt.args.v...)
		})
	}
}

func Test_entry_Debugf(t *testing.T) {
	type fields struct {
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
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
			p := &entry{
				level:        tt.fields.level,
				time:         tt.fields.time,
				callerInfo:   tt.fields.callerInfo,
				message:      tt.fields.message,
				ctx:          tt.fields.ctx,
				extendFields: tt.fields.extendFields,
			}
			p.Debugf(tt.args.format, tt.args.v...)
		})
	}
}

func Test_entry_Error(t *testing.T) {
	type fields struct {
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
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
			p := &entry{
				level:        tt.fields.level,
				time:         tt.fields.time,
				callerInfo:   tt.fields.callerInfo,
				message:      tt.fields.message,
				ctx:          tt.fields.ctx,
				extendFields: tt.fields.extendFields,
			}
			p.Error(tt.args.v...)
		})
	}
}

func Test_entry_Errorf(t *testing.T) {
	type fields struct {
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
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
			p := &entry{
				level:        tt.fields.level,
				time:         tt.fields.time,
				callerInfo:   tt.fields.callerInfo,
				message:      tt.fields.message,
				ctx:          tt.fields.ctx,
				extendFields: tt.fields.extendFields,
			}
			p.Errorf(tt.args.format, tt.args.v...)
		})
	}
}

func Test_entry_Fatal(t *testing.T) {
	type fields struct {
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
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
			p := &entry{
				level:        tt.fields.level,
				time:         tt.fields.time,
				callerInfo:   tt.fields.callerInfo,
				message:      tt.fields.message,
				ctx:          tt.fields.ctx,
				extendFields: tt.fields.extendFields,
			}
			p.Fatal(tt.args.v...)
		})
	}
}

func Test_entry_Fatalf(t *testing.T) {
	type fields struct {
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
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
			p := &entry{
				level:        tt.fields.level,
				time:         tt.fields.time,
				callerInfo:   tt.fields.callerInfo,
				message:      tt.fields.message,
				ctx:          tt.fields.ctx,
				extendFields: tt.fields.extendFields,
			}
			p.Fatalf(tt.args.format, tt.args.v...)
		})
	}
}

func Test_entry_Info(t *testing.T) {
	type fields struct {
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
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
			p := &entry{
				level:        tt.fields.level,
				time:         tt.fields.time,
				callerInfo:   tt.fields.callerInfo,
				message:      tt.fields.message,
				ctx:          tt.fields.ctx,
				extendFields: tt.fields.extendFields,
			}
			p.Info(tt.args.v...)
		})
	}
}

func Test_entry_Infof(t *testing.T) {
	type fields struct {
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
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
			p := &entry{
				level:        tt.fields.level,
				time:         tt.fields.time,
				callerInfo:   tt.fields.callerInfo,
				message:      tt.fields.message,
				ctx:          tt.fields.ctx,
				extendFields: tt.fields.extendFields,
			}
			p.Infof(tt.args.format, tt.args.v...)
		})
	}
}

func Test_entry_Trace(t *testing.T) {
	type fields struct {
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
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
			p := &entry{
				level:        tt.fields.level,
				time:         tt.fields.time,
				callerInfo:   tt.fields.callerInfo,
				message:      tt.fields.message,
				ctx:          tt.fields.ctx,
				extendFields: tt.fields.extendFields,
			}
			p.Trace(tt.args.v...)
		})
	}
}

func Test_entry_Tracef(t *testing.T) {
	type fields struct {
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
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
			p := &entry{
				level:        tt.fields.level,
				time:         tt.fields.time,
				callerInfo:   tt.fields.callerInfo,
				message:      tt.fields.message,
				ctx:          tt.fields.ctx,
				extendFields: tt.fields.extendFields,
			}
			p.Tracef(tt.args.format, tt.args.v...)
		})
	}
}

func Test_entry_Warn(t *testing.T) {
	type fields struct {
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
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
			p := &entry{
				level:        tt.fields.level,
				time:         tt.fields.time,
				callerInfo:   tt.fields.callerInfo,
				message:      tt.fields.message,
				ctx:          tt.fields.ctx,
				extendFields: tt.fields.extendFields,
			}
			p.Warn(tt.args.v...)
		})
	}
}

func Test_entry_Warnf(t *testing.T) {
	type fields struct {
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
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
			p := &entry{
				level:        tt.fields.level,
				time:         tt.fields.time,
				callerInfo:   tt.fields.callerInfo,
				message:      tt.fields.message,
				ctx:          tt.fields.ctx,
				extendFields: tt.fields.extendFields,
			}
			p.Warnf(tt.args.format, tt.args.v...)
		})
	}
}

func Test_entry_WithContext(t *testing.T) {
	type fields struct {
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
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
			p := &entry{
				level:        tt.fields.level,
				time:         tt.fields.time,
				callerInfo:   tt.fields.callerInfo,
				message:      tt.fields.message,
				ctx:          tt.fields.ctx,
				extendFields: tt.fields.extendFields,
			}
			if got := p.WithContext(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entry_WithExtendField(t *testing.T) {
	type fields struct {
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
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
			p := &entry{
				level:        tt.fields.level,
				time:         tt.fields.time,
				callerInfo:   tt.fields.callerInfo,
				message:      tt.fields.message,
				ctx:          tt.fields.ctx,
				extendFields: tt.fields.extendFields,
			}
			if got := p.WithExtendField(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithExtendField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entry_WithExtendFields(t *testing.T) {
	type fields struct {
		level        int
		time         time.Time
		callerInfo   string
		message      string
		ctx          context.Context
		extendFields extendFields
	}
	type args struct {
		fields extendFields
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
			p := &entry{
				level:        tt.fields.level,
				time:         tt.fields.time,
				callerInfo:   tt.fields.callerInfo,
				message:      tt.fields.message,
				ctx:          tt.fields.ctx,
				extendFields: tt.fields.extendFields,
			}
			if got := p.WithExtendFields(tt.args.fields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithExtendFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entry_formatMessage(t *testing.T) {
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
		want   string
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
			if got := p.formatMessage(); got != tt.want {
				t.Errorf("formatMessage() = %v, want %v", got, tt.want)
			}
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
			if err := p.Start(tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
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
				t.Errorf("WithContext() = %v, want %v", got, tt.want)
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
