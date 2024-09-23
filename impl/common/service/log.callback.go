package service

// todo menglc 该文件夹需要提交
import (
	"fmt"
	xruntime "xcore/lib/runtime"
)

func logCallBackFunc(level uint32, outString string) {
	if xruntime.IsDebug() {
		fmt.Println(level, outString)
	}
	return
}
