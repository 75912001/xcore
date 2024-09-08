// Package util 开关
// 两种状态：开启、关闭
package util

// ISwitch interface
type ISwitch interface {
	Enable()
	Disable()
	IsEnabled() bool
	IsDisabled() bool
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type defaultSwitch struct {
	enabled bool // 是否启用
}

// NewDefaultSwitch creates a new defaultSwitch
func NewDefaultSwitch(enable bool) *defaultSwitch {
	s := new(defaultSwitch)
	if enable {
		s.Enable()
	}
	return s
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

// IsDisabled checks if the switch is disabled
func (s *defaultSwitch) IsDisabled() bool {
	return !s.enabled
}
