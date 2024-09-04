package runtime

// RunMode 运行模式
type RunMode uint32

const (
	RunModeRelease RunMode = 0 // release 模式
	RunModeDebug   RunMode = 1 // debug 模式
)

// 程序运行模式
var programRunMode = RunModeRelease

func SetRunMode(mode RunMode) {
	programRunMode = mode
}

// IsDebug 是否为调试模式
func IsDebug() bool {
	return programRunMode == RunModeDebug
}

// IsRelease 是否为发行模式
func IsRelease() bool {
	return programRunMode == RunModeRelease
}
