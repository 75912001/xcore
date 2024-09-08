package service

import (
	"fmt"
	xlog "xcore/lib/log"
	xruntime "xcore/lib/runtime"
)

// 终端钩子
type terminalHook struct {
}

// 创建终端钩子
func newTerminalHook() xlog.IHook {
	return &terminalHook{}
}

func (p *terminalHook) Levels() []uint32 {
	return []uint32{xlog.LevelFatal, xlog.LevelError, xlog.LevelWarn}
}

func (p *terminalHook) Fire(outString string) error {
	if xruntime.IsDebug() {
		fmt.Println(outString)
	}
	return nil
}
