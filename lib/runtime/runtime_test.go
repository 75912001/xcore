package runtime

import (
	"fmt"
	"github.com/agiledragon/gomonkey"
	"os"
	"runtime"
	"testing"
	"xcore/lib/constants"
)

func TestIsDebug(t *testing.T) {
	tests := []struct {
		name    string
		want    bool
		preFunc func()
	}{
		{
			name: constants.Normal,
			want: false,
			preFunc: func() {
				programRunMode = RunModeRelease
			},
		},
		{
			name: constants.Normal,
			preFunc: func() {
				programRunMode = RunModeDebug
			},
			want: true,
		},
	}
	for _, tt := range tests {
		tt.preFunc()
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDebug(); got != tt.want {
				t.Errorf("IsDebug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRelease(t *testing.T) {
	tests := []struct {
		name    string
		want    bool
		preFunc func()
	}{
		{
			name: constants.Normal,
			want: true,
			preFunc: func() {
				programRunMode = RunModeRelease
			},
		},
		{
			name: constants.Normal,
			preFunc: func() {
				programRunMode = RunModeDebug
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt.preFunc()
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRelease(); got != tt.want {
				t.Errorf("IsRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocation(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: constants.Normal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc, fileName, line, _ := runtime.Caller(0)
			funcName := runtime.FuncForPC(pc).Name()
			tt.want = fmt.Sprintf("file:%v line:%v func:%v", fileName, line+3, funcName)
			if got := Location(); got != tt.want {
				t.Errorf("Location() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_codeLocation_String(t *testing.T) {
	type fields struct {
		fileName string
		funcName string
		line     int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: constants.Normal,
			fields: fields{
				fileName: "fileName",
				funcName: "funcName",
				line:     100,
			},
			want: fmt.Sprintf("file:%v line:%v func:%v", "fileName", 100, "funcName"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &codeLocation{
				fileName: tt.fields.fileName,
				funcName: tt.fields.funcName,
				line:     tt.fields.line,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsWindows(t *testing.T) {
	t.Logf("IsWindows:%v", IsWindows())
}

func TestIsLinux(t *testing.T) {
	t.Logf("IsLinux:%v", IsLinux())
}

func TestGetExecutablePath(t *testing.T) {
	type fields struct {
		patches *gomonkey.Patches
	}
	tests := []struct {
		name     string
		fields   fields
		wantErr  bool
		preFunc  func(*fields)
		postFunc func(*fields)
	}{
		{
			name: "normal case",
			fields: fields{
				patches: gomonkey.NewPatches(),
			},
			preFunc: func(args *fields) {
				args.patches.ApplyFunc(os.Executable, func() (string, error) {
					return "/path/to/executable", nil
				})
			},
			postFunc: func(args *fields) {
				args.patches.Reset()
			},
			wantErr: false,
		},
		{
			name: "os.Executable returns an error",
			fields: fields{
				patches: gomonkey.NewPatches(),
			},
			preFunc: func(args *fields) {
				args.patches.ApplyFunc(os.Executable, func() (string, error) {
					return "", fmt.Errorf("forced error")
				})
			},
			postFunc: func(args *fields) {
				args.patches.Reset()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt.preFunc(&tt.fields)
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetExecutablePath()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExecutablePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		tt.postFunc(&tt.fields)
	}
}

func TestSetRunMode(t *testing.T) {
	tests := []struct {
		name string
		mode runMode
		want runMode
	}{
		{
			name: "Set RunMode to Debug",
			mode: RunModeDebug,
			want: RunModeDebug,
		},
		{
			name: "Set RunMode to Release",
			mode: RunModeRelease,
			want: RunModeRelease,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetRunMode(tt.mode)
			if programRunMode != tt.want {
				t.Errorf("SetRunMode() = %v, want %v", programRunMode, tt.want)
			}
		})
	}
}
