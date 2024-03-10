package runtime

import "runtime"

// IsWindows win
func IsWindows() bool {
	return `windows` == runtime.GOOS
}

// IsLinux linux
func IsLinux() bool {
	return `linux` == runtime.GOOS
}
