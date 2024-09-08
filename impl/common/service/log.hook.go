package service

import (
	"fmt"
	xlog "xcore/lib/log"
)

// TerminalHook 终端钩子
type TerminalHook struct {
}

func (p *TerminalHook) Levels() []uint32 {
	return []uint32{xlog.LevelFatal, xlog.LevelError, xlog.LevelWarn}
}

func (p *TerminalHook) Fire(outString string) error {
	fmt.Println(outString)
	return nil
}
