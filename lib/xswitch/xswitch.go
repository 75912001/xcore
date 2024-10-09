// Package xswitch Package util 开关
// 两种状态：开启、关闭
package xswitch

// ISwitch interface
type ISwitch interface {
	Enable()
	Disable()
	IsEnabled() bool
	IsDisabled() bool
}
