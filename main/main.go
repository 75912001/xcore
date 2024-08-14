//package main
//
//import (
//	"fmt"
//	"os"
//	"path/filepath"
//	"reflect"
//)
//
//// TODO 注意:不支持 link/快捷方式
//func GetCurrentPath() (currentPath string, err error) {
//	currentPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
//	if err != nil {
//		return
//	}
//	return
//}
//
//func main() {
//	currentPath, err := os.Getwd()
//	if err != nil {
//		fmt.Println("Failed to get current directory:", err)
//		return
//	}
//
//	fmt.Println("Current directory:", currentPath)
//	fmt.Println(GetCurrentPath())
//	fmt.Println(os.Executable())
//	s, _ := os.Executable()
//	fmt.Println(filepath.Dir(s))
//}

package main

import (
	"context"
	"time"
	"xcore/lib/log"
	"xcore/lib/timer"
)

func main() {
	//l, _ := log.NewMgr()
	//l.Debug("this is debug log")
	log.PrintInfo("this is info log")
	log.PrintErr("this is error log")
	//l.Stop()
	busChannel := make(chan interface{}, 100)
	if true {
		// 测试 定时器
		timerMgr := timer.NewMgr()
		err := timerMgr.Start(context.Background(),
			timer.NewOptions().
				WithOutgoingTimerOutChan(busChannel))
		if err != nil {
			panic(err)
		}
		timerMgr.AddSecond(func(arg interface{}) {
			argValue := arg.(string)
			log.PrintInfo("second timer %s", argValue)
		}, "1", time.Now().Unix()+10)

	}

	for {
		// busChannel 取出数据
		select {
		case v := <-busChannel:
			switch t := v.(type) {
			case *timer.Second:
				t.Function(t.Arg)
			}
		}
	}
	return
}
