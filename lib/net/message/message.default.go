package message

import (
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	xcallback "xcore/lib/callback"
	xerror "xcore/lib/error"
	xruntime "xcore/lib/runtime"
	xswitch "xcore/lib/xswitch"
)

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

// Marshal 序列化
func (p *defaultMessage) Marshal(message proto.Message) (data []byte, err error) {
	data, err = proto.Marshal(message)
	if err != nil {
		return nil, errors.WithMessage(err, xruntime.Location())
	}
	return data, nil
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

// JsonUnmarshal 反序列化
//
//	message: 反序列化 得到的 消息
func (p *defaultMessage) JsonUnmarshal(data []byte) (message proto.Message, err error) {
	message = p.newProtoMessage()
	err = protojson.Unmarshal(data, message)
	if err != nil {
		return nil, errors.WithMessage(err, xruntime.Location())
	}
	return message, nil
}
