package runtime

import (
	"fmt"
	"runtime"
	"xcore/lib/constants"
)

// 代码位置
type codeLocation struct {
	fileName string //文件名
	funcName string //函数名
	line     int    //行数
}

// String 信息
func (p *codeLocation) String() string {
	return fmt.Sprintf("file:%v line:%v func:%v", p.fileName, p.line, p.funcName)
}

// Location 获取代码位置
func Location() string {
	location := &codeLocation{
		fileName: constants.Unknown,
		funcName: constants.Unknown,
	}
	pc, fileName, line, ok := runtime.Caller(1)
	if ok {
		location.fileName = fileName
		location.line = line
		location.funcName = runtime.FuncForPC(pc).Name()
	}
	return location.String()
}
