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

// CreateObject 创建错误码对象,初始化程序的时候创建,创建失败会 panic.
func CreateObject(code uint32) *object {
	newError := newObject(code)
	e := checkDuplication(newError)
	if e != nil {
		panic(
			errors.WithMessage(e, fmt.Sprintf("create error object duplicates %v %#x", code, code)),
		)
	}
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
	errMap[err.code] = struct{}{}
	return nil
}

func newObject(code uint32) *object {
	return &object{
		code: code,
	}
}
