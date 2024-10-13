package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	xconstants "xcore/lib/constants"
	xlog "xcore/lib/log"
	xnettcp "xcore/lib/net/tcp"
	xruntime "xcore/lib/runtime"
)

var apiDataJsonPath string
var client *defaultClient

func main() {
	var err error
	xruntime.SetRunMode(xruntime.RunModeDebug)
	// 启动日志
	glog, err = xlog.NewMgr(xlog.NewOptions().
		WithLevelCallBack(logCallBackFunc, xlog.LevelFatal, xlog.LevelError, xlog.LevelWarn),
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = glog.Stop()
		if err != nil {
			panic(err)
		}
	}()
	ctx := context.Background()
	// 程序所在路径(如为link,则为link所在的路径)
	var executablePath string
	if executablePath, err = xruntime.GetExecutablePath(); err != nil {
		panic(err)
	}
	apiDataJsonPath = path.Join(executablePath, "apiData.json")
	{
		busChannel := make(chan interface{}, xconstants.BusChannelCapacityDefault)
		go func() {
			_ = Handle(busChannel)
		}()
		client = &defaultClient{}
		client.Client = xnettcp.NewClient(client)
		err := client.Connect(ctx, xnettcp.NewClientOptions().
			WithAddress("127.0.0.1:30201").
			WithEventChan(busChannel).
			WithSendChanCapacity(1000))
		if err != nil {
			fmt.Println("connect fail:", err)
			panic(err)
		}
		for {
			_, err = fmt.Scan()
			command, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				fmt.Println("Scan fail, err:", err)
				err = nil
				continue
			}
			command = strings.TrimSpace(command)
			busChannel <- &EventCommand{Command: command}
		}
	}
}

type EventCommand struct {
	Command string
}

type ApiData struct {
	ID       string
	Method   string
	Msg      map[string]interface{}
	Commands []string
}
