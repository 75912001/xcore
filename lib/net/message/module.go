package message

import (
	"github.com/pkg/errors"
	xcallback "xcore/lib/callback"
	xerror "xcore/lib/error"
	xnetpacket "xcore/lib/net/packet"
	xruntime "xcore/lib/runtime"
	"xcore/lib/xswitch"
)
import "context"
import "google.golang.org/protobuf/proto"

// HandlerFunc 处理函数
type HandlerFunc func(ctx context.Context, header xnetpacket.IHeader, message proto.Message, obj interface{}) error

type IMessage interface {
	xcallback.ICallBack
	Unmarshal(data []byte) (message proto.Message, err error)
}

type defaultMessage struct {
	xcallback.ICallBack
	newProtoMessage   func() proto.Message // 创建新的 proto.Message
	stateSwitch       xswitch.ISwitch      // 状态开关-该消息是否启用
	passThroughSwitch xswitch.ISwitch      // 透传开关-该消息是否透传
}

func newDefaultMessage(opts *options) IMessage {
	return &defaultMessage{
		ICallBack:         opts.callback,
		newProtoMessage:   opts.newProtoMessage,
		stateSwitch:       opts.stateSwitch,
		passThroughSwitch: opts.passThroughSwitch,
	}
}

func (p *defaultMessage) Execute() error {
	if p.stateSwitch.IsDisabled() { // 消息是否禁用
		return xerror.MessageIDDisable
	}
	return p.ICallBack.Execute()
}

// Unmarshal 反序列化
//
//	message: 反序列化 得到的 消息
func (p *defaultMessage) Unmarshal(data []byte) (message proto.Message, err error) {
	message = p.newProtoMessage()
	err = proto.Unmarshal(data, message)
	if err != nil {
		return nil, errors.WithMessage(err, xruntime.Location())
	}
	return message, nil
}
