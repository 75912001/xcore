package error

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"testing"
	"xcore/lib/constants"
)

//go:generate go test -v -gcflags=all=-l -coverprofile=coverage.out
//go:generate go tool cover -html=coverage.out -o coverage.html

func TestObjectError(t *testing.T) {
	type fields struct {
		code         uint32
		name         string
		desc         string
		extraMessage string
		extraError   error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: constants.Normal,
			fields: fields{
				code: Success.code,
				name: Success.name,
				desc: Success.desc,
			},
			want: "",
		},
		{
			name: constants.Normal,
			fields: fields{
				code: Link.code,
				name: Link.name,
				desc: Link.desc,
			},
			want: fmt.Sprintf("name:%v code:%v %#x description:%v extraMessage:%v extraError:%v",
				Link.name, Link.code, Link.code, Link.desc, Link.extraMessage, Link.extraError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &object{
				code:         tt.fields.code,
				name:         tt.fields.name,
				desc:         tt.fields.desc,
				extraMessage: tt.fields.extraMessage,
				extraError:   tt.fields.extraError,
			}
			if got := p.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectWithExtra(t *testing.T) {
	obj := object{
		code:         0xffff,
		name:         "0xffff-name",
		desc:         "0xffff-desc",
		extraMessage: "0xffff-message",
		extraError:   errors.New("0xffff-error"),
	}
	tests := []struct {
		name   string
		fields object
		want   *object
	}{
		{
			name:   constants.Normal,
			fields: obj,
			want: newObject(obj.code).WithName(obj.name).WithDesc(obj.desc).
				WithExtraError(obj.extraError).WithExtraMessage(obj.extraMessage),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &object{
				code:         tt.fields.code,
				name:         tt.fields.name,
				desc:         tt.fields.desc,
				extraMessage: tt.fields.extraMessage,
				extraError:   tt.fields.extraError,
			}
			if !reflect.DeepEqual(p, tt.want) {
				t.Errorf("WithExtraError() = %v, want %v", p, tt.want)
			}
		})
	}
}

func TestCreateObject(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// 在这里处理 panic 的情况
			t.Logf("这里应该panic %v", r)
		}
	}()

	_ = NewError(Unknown.code).WithName(Unknown.name).WithDesc(Unknown.desc)
	t.Errorf("期望的 panic 没有出现")
}
