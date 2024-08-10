package log

import (
    "github.com/pkg/errors"
    libruntime "xcore/lib/runtime"
)

// Hook 钩子
// 钩子是一个接口，它定义了x个日志级别的钩子函数
// e.g.: 修改entry的内容, 将 entry 发送至其他目标地
type Hook interface {
    Levels() []int           // 需要hook的等级列表
    Fire(entry *entry) error // 执行的方法
}

type LevelHookMap map[int][]Hook // key: 日志等级, value: 钩子

// add 添加钩子
func (p LevelHookMap) add(hook Hook) {
    for _, level := range hook.Levels() {
        p[level] = append(p[level], hook)
    }
}

// fire 处理钩子
func (p LevelHookMap) fire(entry *entry) error {
    for _, hook := range p[entry.level] {
        if err := hook.Fire(entry); err != nil {
            return errors.WithMessage(err, libruntime.Location())
        }
    }
    return nil
}
