package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strconv"
	xerror "xcore/lib/error"
	xnettcp "xcore/lib/net/tcp"
)

func Handle(busChannel chan interface{}) error {
	for {
		select {
		case value := <-busChannel:
			var err error
			switch event := value.(type) {
			//tcp
			case *xnettcp.Connect:
				err = event.IHandler.OnConnect(event.IRemote)
			case *xnettcp.Packet:
				err = event.IHandler.OnPacket(event.IRemote, event.IPacket)
			case *xnettcp.Disconnect:
				err = event.IHandler.OnDisconnect(event.IRemote)
				if event.IRemote.IsConnect() {
					event.IRemote.Stop()
				}
			case *EventCommand:
				// 创建一个 map 来存储 JSON 数据
				data := make(map[string]ApiData)
				apiData := ApiData{}
				parseCommand := func(command string) error {
					file, err := os.Open(apiDataJsonPath)
					if err != nil {
						fmt.Printf("Error opening file:%v %v", apiDataJsonPath, err)
						panic(err)
					}
					defer func() {
						_ = file.Close()
					}()
					// 创建一个新的解码器
					decoder := json.NewDecoder(file)
					// 解码 JSON 数据到 map 中
					err = decoder.Decode(&data)
					if err != nil {
						fmt.Println("Error decoding JSON:", err)
						panic(err)
					}
					if info, ok := data[command]; ok {
						//fmt.Printf("apiData: %+v\n", info)
						apiData = info
					} else {
						fmt.Printf("\033[31m%s\033[0m\n", "api not found in apiData.json")
						return xerror.NotExist
					}
					return nil
				}
				err = parseCommand(event.Command)
				if err != nil {
					if errors.Is(err, xerror.NotExist) {
						continue
					}
					panic(err)
				}
				// todo menglc 打包数据
				num, err := strconv.ParseUint(apiData.ID, 0, 32)
				if err != nil {
					fmt.Println("strconv.ParseUint fail, err:", err)
					panic(err)
				}
				messageID := uint32(num)
				//fmt.Printf("%v %#x\n", messageID, messageID)

				// gateway 中, 查找消息
				message := GMessage.Find(messageID)
				if message == nil {
					fmt.Printf("\033[31m%s\033[0m\n", "message not found")
					continue
				} else {
					//fmt.Printf("message: %v\n", message)
				}
				// 将 apiData 的数据,构建成消息
				msgData, err := json.Marshal(apiData.Msg)
				if err != nil {
					fmt.Println("json.Marshal fail, err:", err)
					continue
				}
				protoMsg, err := message.JsonUnmarshal(msgData)
				if err != nil {
					fmt.Println("message.Unmarshal fail, err:", err)
					continue
				}
				//sendData, err := message.Marshal(protoMsg)
				if err != nil {
					fmt.Println("message.Marshal fail, err:", err)
					continue
				}
				fmt.Println()
				fmt.Printf("\033[34mmessageID: 0x%x\033[0m\n", messageID)
				fmt.Printf("\033[34mMessage: %v\033[0m\n", protoMsg)
				log.Infof("\n======send message======\n%v\nmessageID: 0x%x\nMessage: %v", event.Command, messageID, protoMsg)
				if err := xnettcp.Send(client.IRemote, protoMsg, messageID, 0, 0); err != nil {
					fmt.Println("client.Send fail, err:", err)
					log.Errorf("client.Send fail, err: %v", err)
				}
			default:
			}
			if err != nil {
			}

		}
	}
}
