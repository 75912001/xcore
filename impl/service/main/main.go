// 服务main
// 参数
// [0]进程名. e.g.:gateway.exe
// [1]组ID-GroupID. e.g.:10001
// [2]服务名称. e.g.:gateway
// [3]服务ID. e.g.:10001

package main

import (
	"context"
	"os"
	common "xcore/impl/common"
	servicegateway "xcore/impl/service/gateway"
	xerror "xcore/lib/error"
	xlog "xcore/lib/log"
	xruntime "xcore/lib/runtime"
	xservice "xcore/lib/service"
)

func main() {
	var err error
	defaultService := xservice.NewService(os.Args)
	if defaultService == nil {
		panic("NewService failed")
	}
	var service xservice.IService
	switch defaultService.Name {
	case common.ServiceNameGateway:
		service = servicegateway.NewService(defaultService)
	default:
		xlog.PrintErr(xerror.NotImplemented, "service name err", defaultService.Name)
		return
	}
	if err = service.Start(context.Background()); err != nil {
		xlog.PrintErr(err, xruntime.Location())
		return
	}
	return
}
