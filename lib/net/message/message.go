package message

import (
	"context"
	xconstants "xcore/lib/constants"
	xerror "xcore/lib/error"
	"xcore/lib/net/packet"
	xruntime "xcore/lib/runtime"
	xutil "xcore/lib/util"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type Message struct {
	handler           Handler              // [required] 消息处理函数
	newProtoMessage   func() proto.Message // [required] 创建新的 proto.Message
	name              string               // [optional] [default]: "Unknown" 名称, 在设置 handler 时, 会自动设置
	stateSwitch       xutil.ISwitch        // [optional] 状态开关-该消息是否启用 [default]:true
	passThroughSwitch xutil.ISwitch        // [optional] 透传开关-该消息是否透传 [default]:false
}

// NewMessage 创建 Message
func NewMessage() *Message {
	return new(Message)
}

func (p *Message) WithHandler(handler Handler) *Message {
	p.handler = handler
	p.name = getFuncName(handler, '.')
	return p
}

func (p *Message) WithNewProtoMessage(newProtoMessage func() proto.Message) *Message {
	p.newProtoMessage = newProtoMessage
	return p
}

func (p *Message) WithName(name string) *Message {
	p.name = name
	return p
}

func (p *Message) WithStateSwitch(stateSwitch xutil.ISwitch) *Message {
	p.stateSwitch = stateSwitch
	return p
}

func (p *Message) WithPassThroughSwitch(passThroughSwitch xutil.ISwitch) *Message {
	p.passThroughSwitch = passThroughSwitch
	return p
}

// Unmarshal 反序列化
//
//	message: 反序列化 得到的 消息
func (p *Message) Unmarshal(data []byte) (message proto.Message, err error) {
	message = p.newProtoMessage()
	err = proto.Unmarshal(data, message)
	if err != nil {
		return nil, errors.WithMessage(err, xruntime.Location())
	}
	return message, nil
}

// Handler 处理
func (p *Message) Handler(ctx context.Context, header packet.IHeader, message proto.Message, obj interface{}) error {
	// 消息是否禁用
	if p.stateSwitch.IsDisabled() {
		return xerror.MessageIDDisable
	}
	return p.handler(ctx, header, message, obj)
}

func merge(opts ...*Message) *Message {
	so := NewMessage()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.handler != nil {
			so.handler = opt.handler
		}
		if opt.newProtoMessage != nil {
			so.newProtoMessage = opt.newProtoMessage
		}
		if 0 < len(opt.name) {
			so.name = opt.name
		}
		if opt.stateSwitch != nil {
			so.stateSwitch = opt.stateSwitch
		}
		if opt.passThroughSwitch != nil {
			so.passThroughSwitch = opt.passThroughSwitch
		}
	}
	return so
}

// 配置
func configure(opts *Message) error {
	if len(opts.name) == 0 {
		opts.name = xconstants.Unknown
	}
	if opts.stateSwitch == nil {
		opts.stateSwitch = xutil.NewDefaultSwitch(true)
	}
	if opts.passThroughSwitch == nil {
		opts.passThroughSwitch = xutil.NewDefaultSwitch(false)
	}
	if opts.passThroughSwitch.IsDisabled() { // 非 透传
		if opts.handler == nil { // 没有处理函数
			return errors.WithMessage(xerror.Param, xruntime.Location())
		}
		if opts.newProtoMessage == nil { // 没有创建消息函数
			return errors.WithMessage(xerror.Param, xruntime.Location())
		}
	}
	return nil
}
