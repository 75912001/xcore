package log

const (
	TraceIDKey = "TraceID" // 日志traceId key
	UserIDKey  = "UID"     // 日志用户ID key
)
const (
	logChannelCapacity = 100000            // 日志通道最大容量
	logTimeFormat      = "15:04:05.000000" // 日志时间格式 时:分:秒.微秒
	callerInfoFormat   = "line:%d %s"      // 堆栈信息格式 例如 line:69 server/xxx/xx/x/log.TestExample
	fileFormat         = "%s/%s-%d-%s"     // 文件全路径名格式 例如 ${filePath}/${prefix}-20200101-normal.log
	bufferCapacity     = 300               // buffer 容量
	calldepth1         = 1
	calldepth2         = calldepth1 + 1
	calldepth3         = calldepth2 + 1
)

const (
	normalLogFileBaseName = "normal.log" // 正常日志文件名
	errorLogFileBaseName  = "error.log"  // 错误日志文件名
)
