// 服务main
// 参数
// [0]进程名. e.g.:gateway.exe
// [1]组ID-GroupID. e.g.:10001
// [2]服务名称. e.g.:gateway
// [3]服务ID. e.g.:10001

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"xcore/impl/common"
	commonservice "xcore/impl/common/service"
	"xcore/impl/service/gateway"
	xconstants "xcore/lib/constants"
	xerror "xcore/lib/error"
	xlog "xcore/lib/log"
	xruntime "xcore/lib/runtime"
)

func main() {
	var err error
	defaultService := commonservice.NewDefaultService()
	// 程序所在路径(如为link,则为link所在的路径)
	defaultService.ExecutablePath, err = xruntime.GetExecutablePath()
	if err != nil {
		xlog.PrintErr(err, xruntime.Location())
		return
	}
	args := os.Args
	argNum := len(args)
	const neededArgsNumber = 4
	if argNum != neededArgsNumber {
		xlog.PrintfErr("the number of parameters is incorrect, needed %v, but %v.", neededArgsNumber, argNum)
		return
	}

	{ // 解析启动参数
		groupID, err := strconv.ParseUint(args[1], 10, 32)
		if err != nil {
			xlog.PrintErr("groupID err:", err)
			return
		}
		defaultService.GroupID = uint32(groupID)
		defaultService.Name = args[2]
		serviceID, err := strconv.ParseUint(args[3], 10, 32)
		if err != nil {
			xlog.PrintErr("serviceID err", err)
			return
		}
		defaultService.ID = uint32(serviceID)
		xlog.PrintfInfo("groupID:%v name:%v, serviceID:%v",
			defaultService.GroupID, defaultService.Name, defaultService.ID)
	}
	var service commonservice.IService
	switch defaultService.Name {
	case common.ServiceNameGateway:
		service = gateway.NewService(defaultService)
	default:
		xlog.PrintErr(xerror.NotImplemented, "service name err", defaultService.Name)
		return
	}
	benchPath := path.Join(defaultService.ExecutablePath, fmt.Sprintf("%v.%v.%v.%v",
		defaultService.GroupID, defaultService.Name, defaultService.ID, xconstants.ServiceConfigFile))
	if err = defaultService.PreStart(context.Background(), commonservice.NewOptions().WithBenchPath(benchPath)); err != nil {
		xlog.PrintErr(err, xruntime.Location())
		return
	}
	if err = service.Start(); err != nil {
		xlog.PrintErr(err, xruntime.Location())
		return
	}

	// 退出服务
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
EXIT:
	for {
		select {
		case <-defaultService.QuitChan:
			defaultService.Log.Warn("service will shutdown in a few seconds")
			service.PreShutdown()
			_ = service.Stop()
			break EXIT // 退出循环
		case s := <-sigChan:
			defaultService.Log.Warnf("service got signal: %s, shutting down...", s)
			close(defaultService.QuitChan)
		}
	}
	return
}
