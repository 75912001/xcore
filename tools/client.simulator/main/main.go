package main

import (
	"bufio"
	"context"
	"fmt"
	xconstants "github.com/75912001/xcore/lib/constants"
	xlog "github.com/75912001/xcore/lib/log"
	xnettcp "github.com/75912001/xcore/lib/net/tcp"
	xruntime "github.com/75912001/xcore/lib/runtime"
	"os"
	"path"
	"strings"
)

func main() {
	var err error
	xruntime.SetRunMode(xruntime.RunModeDebug)
	// 启动日志
	log, err = xlog.NewMgr(xlog.NewOptions().
		WithLevelCallBack(logCallBackFunc, xlog.LevelFatal, xlog.LevelError, xlog.LevelWarn),
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = log.Stop()
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
	configJsonPath = path.Join(executablePath, "config.json")
	err = parseConfigJson(configJsonPath)
	if err != nil {
		panic(err)
	}
	{
		busChannel := make(chan interface{}, xconstants.BusChannelCapacityDefault)
		go func() {
			_ = Handle(busChannel)
		}()
		client = &defaultClient{}
		client.Client = xnettcp.NewClient(client)
		err := client.Connect(ctx, xnettcp.NewClientOptions().
			WithAddress(configJson.Addr).
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
			//for i := 0; i < 10; i++ {
			busChannel <- &EventCommand{Command: command}
			//}
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
