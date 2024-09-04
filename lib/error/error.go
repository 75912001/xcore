// 错误码
// 错误码是程序运行过程中发生错误时返回给调用者的错误代码。
// 错误码类型: 整数类型
// 错误码范围: 0-0xffffffff
// 系统错误码: [0x0,0xffff] e.g.: system.code.go
// 业务错误码: [0x10000000,0x1fffffff]
// 用户错误码: [0x20000000,0x2fffffff]
// 示例
//  0x0: 成功
//  0xf001: 链接错误
//  0x10000001: 资源数量不足

package error

import (
	"fmt"
	"github.com/pkg/errors"
)

// 对象
type object struct {
	code         uint32 // 码
	name         string // 名称
	desc         string // 描述 description
	extraMessage string // 附加信息
	extraError   error  // 附加错误
}

// 错误信息
func (p *object) Error() string {
	if Success.code == p.code {
		return ""
	}
	return fmt.Sprintf("name:%v code:%v %#x description:%v extraMessage:%v extraError:%v",
		p.name, p.code, p.code, p.desc, p.extraMessage, p.extraError)
}

func (p *object) WithExtraMessage(extraMessage string) *object {
	p.extraMessage = extraMessage
	return p
}

func (p *object) WithExtraError(extraError error) *object {
	p.extraError = extraError
	return p
}

func (p *object) WithName(name string) *object {
	p.name = name
	return p
}

func (p *object) WithDesc(desc string) *object {
	p.desc = desc
	return p
}

// NewError 创建错误码,初始化程序的时候创建,创建失败会 panic.
func NewError(code uint32) *object {
	newError := newObject(code)
	e := checkDuplication(newError)
	if e != nil {
		panic(
			errors.WithMessage(e, fmt.Sprintf("new error object duplicates %v %#x", code, code)),
		)
	}
	errMap[code] = struct{}{}
	return newError
}

// 错误信息
// 用来确保 错误码-唯一性
var errMap = make(map[uint32]struct{})

// 检查重复情况
func checkDuplication(err *object) error {
	if _, ok := errMap[err.code]; ok { // 重复
		return errors.WithStack(errors.Errorf("duplicate err, code:%v %#x", err.code, err.code))
	}
	return nil
}

func newObject(code uint32) *object {
	return &object{
		code: code,
	}
}
