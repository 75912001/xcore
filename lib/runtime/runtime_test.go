package runtime

import (
	"fmt"
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
				GRunMode = RunModeRelease
			},
		},
		{
			name: constants.Normal,
			preFunc: func() {
				GRunMode = RunModeDebug
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

//func TestIsLinux(t *testing.T) {
//	tests := []struct {
//		name string
//		want bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := IsLinux(); got != tt.want {
//				t.Errorf("IsLinux() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

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
				GRunMode = RunModeRelease
			},
		},
		{
			name: constants.Normal,
			preFunc: func() {
				GRunMode = RunModeDebug
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

//func TestIsWindows(t *testing.T) {
//	tests := []struct {
//		name string
//		want bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := IsWindows(); got != tt.want {
//				t.Errorf("IsWindows() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

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
