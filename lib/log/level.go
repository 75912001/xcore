package log

// 日志等级
const (
	LevelOff   int = 0 //关闭
	LevelFatal int = 1 //致命
	LevelError int = 2 //错误
	LevelWarn  int = 3 //警告
	LevelInfo  int = 4 //信息
	LevelDebug int = 5 //调试
	LevelTrace int = 6 //跟踪
	LevelOn    int = 7 //全部打开
)

// 等级描述
var levelDesc = []string{
	LevelOff:   "LevelOff", //关闭
	LevelFatal: "FAT",      //致命
	LevelError: "ERR",      //错误
	LevelWarn:  "WAR",      //警告
	LevelInfo:  "INF",      //信息
	LevelDebug: "DEB",      //调试
	LevelTrace: "TRA",      //跟踪
	LevelOn:    "LevelOn",  //全部打开
}
