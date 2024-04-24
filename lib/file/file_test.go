package file

import (
	"fmt"
	"github.com/agiledragon/gomonkey"
	"github.com/pkg/errors"
	"os"
	"path"
	"reflect"
	"testing"
	libruntime "xcore/lib/runtime"
)

// 测试文件名称
var __testFileName__ = "__testFileName__"

// 测试文件内容
const __testFileContent__ string = `
this is test file content.
`

// 当前执行的程序-绝对路径
var __testExecutablePath__ string

// 测试-路径-文件
var __testPathFile__ string

func init() {
	__testExecutablePath__, _ = libruntime.GetExecutablePath()
	__testPathFile__ = path.Join(__testExecutablePath__, __testFileName__)
}

func TestCreateDirectory(t *testing.T) {
	type args struct {
		path    string
		patches *gomonkey.Patches
	}
	// 用例
	type useCase struct {
		name     string
		args     args
		wantErr  bool
		preFunc  func(*args)
		postFunc func(*args)
	}

	tests := []useCase{
		{
			name: "normal-没有目录",
			args: args{
				path: path.Join(__testExecutablePath__, "__test_dir1__", "__test_dir2__"),
			},
			preFunc: func(args *args) {
				_ = os.RemoveAll(args.path)
			},
			wantErr: false,
			postFunc: func(args *args) {
				_ = os.RemoveAll(path.Dir(args.path))
			},
		},
		{
			name: "normal-已有目录",
			args: args{
				path: path.Join(__testExecutablePath__, "__test_dir1__", "__test_dir2__"),
			},
			preFunc: func(args *args) {
				_ = CreateDirectory(args.path)
			},
			wantErr: false,
			postFunc: func(args *args) {
				_ = os.RemoveAll(path.Dir(args.path))
			},
		},
		{
			name: "normal-没有目录 - os.MkdirAll error",
			args: args{
				path: path.Join(__testExecutablePath__, "__test_dir1__", "__test_dir2__"),
			},
			preFunc: func(args *args) {
				_ = os.RemoveAll(args.path)
				args.patches = gomonkey.ApplyFunc(os.MkdirAll, func(path string, perm os.FileMode) error {
					return os.ErrExist
				})
			},
			wantErr: true,
			postFunc: func(args *args) {
				args.patches.Reset()
				_ = os.RemoveAll(path.Dir(args.path))
			},
		},
	}
	for _, tt := range tests {
		tt.preFunc(&tt.args)
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateDirectory(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("CreateDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		tt.postFunc(&tt.args)
	}
}

func TestNewOptions(t *testing.T) {
	opts := NewOptions()

	tests := []struct {
		name string
		want *options
	}{
		{
			name: "normal",
			want: opts,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	type args struct {
		pathFile string
		data     []byte
		opts     *options
		patches  *gomonkey.Patches
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		preFunc  func(*args)
		postFunc func(*args)
	}{
		{
			name: "normal-append",
			args: args{
				pathFile: __testPathFile__,
				data:     []byte(__testFileContent__),
				opts:     NewOptions().Append(),
			},
			wantErr: false,
			preFunc: func(args *args) {
				_ = os.RemoveAll(args.pathFile)
			},
			postFunc: func(args *args) {
				_ = os.RemoveAll(args.pathFile)
			},
		},
		{
			name: "normal-overwrite",
			args: args{
				pathFile: __testPathFile__,
				data:     []byte(__testFileContent__),
				opts:     NewOptions().Overwrite(),
			},
			wantErr: false,
			preFunc: func(args *args) {
				_ = os.RemoveAll(args.pathFile)
			},
			postFunc: func(args *args) {
				_ = os.RemoveAll(args.pathFile)
			},
		},
		{
			name: "normal-os.OpenFile-error",
			args: args{
				pathFile: __testPathFile__,
				data:     []byte(__testFileContent__),
				opts:     NewOptions().Overwrite(),
			},
			preFunc: func(args *args) {
				_ = os.RemoveAll(args.pathFile)
				args.patches = gomonkey.ApplyFunc(os.OpenFile, func(name string, flag int, perm os.FileMode) (*os.File, error) {
					return nil, os.ErrNotExist
				})
			},
			wantErr: true,
			postFunc: func(args *args) {
				args.patches.Reset()
				_ = os.RemoveAll(args.pathFile)
			},
		},
		{
			name: "normal-file.Write-error",
			args: args{
				pathFile: __testPathFile__,
				data:     []byte(__testFileContent__),
				opts:     NewOptions().Overwrite(),
			},
			preFunc: func(args *args) {
				_ = os.RemoveAll(args.pathFile)
				file := &os.File{}
				args.patches = gomonkey.ApplyMethod(reflect.TypeOf(file), "Write", func(_ *os.File, b []byte) (int, error) {
					return 0, errors.New("forced error")
				})
			},
			wantErr: true,
			postFunc: func(args *args) {
				args.patches.Reset()
				_ = os.RemoveAll(args.pathFile)
			},
		},
	}
	for _, tt := range tests {
		tt.preFunc(&tt.args)
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteFile(tt.args.pathFile, tt.args.data, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("WriteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		tt.postFunc(&tt.args)
	}
}

func Test_options_Append(t *testing.T) {
	opt := NewOptions()
	opt.overwrite = false
	opt.append = true

	tests := []struct {
		name string
		want *options
	}{
		{
			name: "normal",
			want: opt,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fo := NewOptions()
			if got := fo.Append(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Append() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_options_Overwrite(t *testing.T) {
	opt := NewOptions()
	opt.overwrite = true
	opt.append = false

	tests := []struct {
		name string
		want *options
	}{
		{
			name: "normal",
			want: opt,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fo := NewOptions()
			if got := fo.Overwrite(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Overwrite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathFileExists(t *testing.T) {
	type args struct {
		pathFile string
		patches  *gomonkey.Patches
	}
	tests := []struct {
		name     string
		args     args
		want     bool
		preFunc  func(*args)
		postFunc func(*args)
	}{
		{
			name: "file exists",
			args: args{
				pathFile: __testPathFile__,
			},
			preFunc: func(args *args) {
				_ = os.RemoveAll(args.pathFile)
				file, _ := os.Create(args.pathFile)
				file.Close()
			},
			postFunc: func(args *args) {
				_ = os.RemoveAll(args.pathFile)
			},
			want: true,
		},
		{
			name: "file does not exist",
			args: args{
				pathFile: __testPathFile__,
			},
			preFunc: func(args *args) {
				_ = os.RemoveAll(args.pathFile)
			},
			postFunc: func(args *args) {
				_ = os.RemoveAll(args.pathFile)
			},
			want: false,
		},
		{
			name: "os.Stat returns an error that is not os.IsNotExist",
			args: args{
				pathFile: __testPathFile__,
				patches:  gomonkey.NewPatches(),
			},
			preFunc: func(args *args) {
				args.patches.ApplyFunc(os.Stat, func(_ string) (os.FileInfo, error) {
					// 返回一个既不是nil也不是os.IsNotExist的错误
					return nil, fmt.Errorf("forced error")
				})
			},
			postFunc: func(args *args) {
				args.patches.Reset()
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt.preFunc(&tt.args)
		t.Run(tt.name, func(t *testing.T) {
			if got := PathFileExists(tt.args.pathFile); got != tt.want {
				t.Errorf("PathFileExists() = %v, want %v", got, tt.want)
			}
		})
		tt.postFunc(&tt.args)
	}
}
