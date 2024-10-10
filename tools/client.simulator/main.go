package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path"
	"strconv"
	"strings"
	xconstants "xcore/lib/constants"
	xerror "xcore/lib/error"
	xlog "xcore/lib/log"
	xnetpacket "xcore/lib/net/packet"
	xnettcp "xcore/lib/net/tcp"
	xruntime "xcore/lib/runtime"
)

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
	apiDataJsonPath := path.Join(executablePath, "apiData.json")

	// todo menglc 加载所有proto的协议
	for {
		busChannel := make(chan interface{}, xconstants.BusChannelCapacityDefault)
		client := &defaultClient{}
		client.Client = xnettcp.NewClient(xnetpacket.NewDefaultPacket(xnetpacket.NewDefaultHeader()), client)
		err := client.Connect(ctx, xnettcp.NewClientOptions().
			WithAddress("127.0.0.1:30201").
			WithEventChan(busChannel).
			WithSendChanCapacity(1000))
		if err != nil {
			fmt.Println("connect fail:", err)
			panic(err)
		}
		for {
			fmt.Print("Command:")
			_, err = fmt.Scan()
			command, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				fmt.Println("Scan fail, err:", err)
				err = nil
				continue
			}
			command = strings.TrimSpace(command)
			// 创建一个 map 来存储 JSON 数据
			data := make(map[string]ApiData)
			apiData := ApiData{}
			parseCommand := func(command string) error {
				file, err := os.Open(apiDataJsonPath)
				if err != nil {
					fmt.Printf("Error opening file:%v %v", apiDataJsonPath, err)
					panic(err)
				}
				defer file.Close()
				// 创建一个新的解码器
				decoder := json.NewDecoder(file)
				// 解码 JSON 数据到 map 中
				err = decoder.Decode(&data)
				if err != nil {
					fmt.Println("Error decoding JSON:", err)
					panic(err)
				}
				if info, ok := data[command]; ok {
					fmt.Printf("apiData: %+v\n", info)
					apiData = info
				} else {
					fmt.Printf("\033[31m%s\033[0m\n", "apiData not found")
					return xerror.NonExistent
				}
				return nil
			}
			err = parseCommand(command)
			if err != nil {
				if errors.Is(err, xerror.NonExistent) {
					continue
				}
				panic(err)
			}
			// todo menglc 打包数据
			num, err := strconv.ParseUint(apiData.ID, 0, 32)
			if err != nil {
				fmt.Println("strconv.ParseUint fail, err:", err)
				return
			}
			messageID := uint32(num)
			fmt.Printf("%v %#x", messageID, messageID)
			// 发送消息给服务器
			//h := xnetpacket.NewDefaultHeader()
			//h.SetPacketLength(24)
			//packet := xnetpacket.NewDefaultPacket()
			//
			//if err = client.Send(ctx, data[command].Method, data[command].Msg); err != nil {
			//
			//}
		}
	}
}

type ApiData struct {
	ID       string
	Method   string
	Msg      map[string]interface{}
	Commands []string
}
