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

import "xcore/lib/log"

func main() {
	element, _ := log.NewMgr()
	element.Debug("this is debug log")
	element.Stop()
	return
}
