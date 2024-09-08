// Package util 开关
// 两种状态：开启、关闭
package util

// ISwitch interface
type ISwitch interface {
	Enable()
	Disable()
	IsEnabled() bool
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type defaultSwitch struct {
	enabled bool // [default:false] 是否启用
}

// NewDefaultSwitch creates a new defaultSwitch
func NewDefaultSwitch() *defaultSwitch {
	return new(defaultSwitch)
}

// Enable the switch
func (s *defaultSwitch) Enable() {
	s.enabled = true
}

// Disable the switch
func (s *defaultSwitch) Disable() {
	s.enabled = false
}

// IsEnabled checks if the switch is enabled
func (s *defaultSwitch) IsEnabled() bool {
	return s.enabled
}
