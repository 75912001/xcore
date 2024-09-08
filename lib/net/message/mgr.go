package message

import (
	"github.com/pkg/errors"
	xerror "xcore/lib/error"
	xlog "xcore/lib/log"
	xruntime "xcore/lib/runtime"
)

type messageMap map[uint32]*Message

// Mgr 管理器
type Mgr struct {
	messageMap messageMap
}

// Register 注册消息
// 重复会 panic
func (p *Mgr) Register(messageID uint32, messageSlice ...*Message) {
	if p.messageMap == nil {
		p.messageMap = make(messageMap)
	}
	if pb := p.Find(messageID); pb != nil {
		xlog.PrintErr(xerror.MessageIDExistent, "%v messageID:%#x %v", xruntime.Location(), messageID, messageID)
		panic(errors.WithMessagef(xerror.MessageIDExistent, "%v messageID:%#x %v",
			xruntime.Location(), messageID, messageID))
	}
	message := merge(messageSlice...)
	if err := configure(message); err != nil {
		xlog.PrintErr(xerror.MessageIDExistent, "%v messageID:%#x %v", xruntime.Location(), messageID, messageID)
		panic(errors.WithMessagef(xerror.MessageIDExistent, "%v messageID:%#x %v",
			xruntime.Location(), messageID, messageID))
	}
	p.messageMap[messageID] = message
}

func (p *Mgr) Find(messageID uint32) *Message {
	return p.messageMap[messageID]
}

// Replace 替换/覆盖(Override)
// 配置错误会 panic
func (p *Mgr) Replace(messageID uint32, messageEntity *Message) {
	if err := configure(messageEntity); err != nil {
		xlog.PrintErr(xerror.MessageIDExistent, "%v messageID:%#x %v", xruntime.Location(), messageID, messageID)
		panic(errors.WithMessagef(xerror.MessageIDExistent, "%v messageID:%#x %v",
			xruntime.Location(), messageID, messageID))
	}
	p.messageMap[messageID] = messageEntity
}
