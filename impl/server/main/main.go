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

	common "github.com/75912001/xcore/impl/common"
	servergateway "github.com/75912001/xcore/impl/server/gateway"
	xerror "github.com/75912001/xcore/lib/error"
	xlog "github.com/75912001/xcore/lib/log"
	xruntime "github.com/75912001/xcore/lib/runtime"
	xserver "github.com/75912001/xcore/lib/server"
)

func main() {
	var err error
	defaultServer := xserver.NewServer(os.Args)
	if defaultServer == nil {
		panic("NewServer failed")
	}
	var server xserver.IServer
	switch defaultServer.Name {
	case common.ServerNameGateway:
		server = servergateway.NewService(defaultServer)
	default:
		xlog.PrintErr(xerror.NotImplemented, "server name err", defaultServer.Name)
		return
	}
	if err = server.Start(context.Background()); err != nil {
		xlog.PrintErr(err, xruntime.Location())
		return
	}
	return
}
