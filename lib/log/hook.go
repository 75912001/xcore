package log

import (
	"github.com/pkg/errors"
	xruntime "xcore/lib/runtime"
)

// IHook 钩子
// 钩子是一个接口，它定义了x个日志级别的钩子函数
// e.g.: 修改entry的内容, 将 entry 发送至其他目标地
type IHook interface {
	Levels() []uint32            // 需要 IHook 的等级列表
	Fire(outString string) error // 执行的方法 (outString: 输出的字符串)
}

type LevelHookMap map[uint32][]IHook // key: 日志等级, value: 钩子

// add 添加钩子
func (p LevelHookMap) add(hook IHook) {
	for _, level := range hook.Levels() {
		p[level] = append(p[level], hook)
	}
}

// fire 处理钩子
func (p LevelHookMap) fire(entry *entry) error {
	for _, hook := range p[entry.level] {
		if err := hook.Fire(entry.outString); err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
	}
	return nil
}
