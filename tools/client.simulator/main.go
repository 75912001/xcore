package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	xlog "xcore/lib/log"
	xruntime "xcore/lib/runtime"
	"xcore/tools/client.simulator/codec/model"
)

func main() {
	var err error
	xruntime.SetRunMode(xruntime.RunModeDebug)
	// 启动日志
	GLog, err = xlog.NewMgr(xlog.NewOptions().
		WithLevelCallBack(logCallBackFunc, xlog.LevelFatal, xlog.LevelError, xlog.LevelWarn),
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = GLog.Stop()
		if err != nil {
			panic(err)
		}
	}()
	for {
		c := Client{}

		resFunc := func(res Res) {
			//if res.ID == 0 && !res.Result {
			//	return
			//}
			out, _ := json.MarshalIndent(res, "", "    ")
			fmt.Println("Res:\n", string(out))
			fmt.Print("Command: ")
		}

		c.Init(resFunc)

		err = c.Connect()
		if err != nil {
			fmt.Printf("connect fail:%+v", err)
			return
		}

		c.StartOnRec()

		fmt.Println("Input Command,Please.")
		fmt.Println("Example:{\"method\":\"/tevat.example.auth.Auth/Login\",\"msg\":{\"account_id\":\"1\",\"account_token\":\"1\"}}")
		fmt.Println("Example:login")

		for {
			fmt.Print("Command: ")
			_, err = fmt.Scan()
			command, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				fmt.Println("Scan fail, err:", err)
				err = nil
				continue
			}
			command = strings.TrimSpace(command)
			data := ApiData{}

			jsonDec := json.NewDecoder(strings.NewReader(command))
			jsonDec.UseNumber()
			err = jsonDec.Decode(&data)
			if err != nil {
				data.Method = command
				err = nil
			}

			if data.Msg != nil {
				msgBytes, err := c.MarshalMsg(data.Method, data.Msg)
				if err != nil {
					fmt.Println("MarshalMsg fail, err:", err)
					continue
				}
				req := &model.Request{
					Method: []byte(data.Method),
					Msg:    msgBytes,
				}
				id := c.reqPack.GetNewId(string(req.Method))
				req.ID = id
				c.Send(req)
			} else if data.Method == "LOAD" || data.Method == "L" || data.Method == "l" {
				fmt.Println("Start reload api data...")
				c.ReloadApi()
				fmt.Println("Reload api data finish.")
			} else if data.Method == "RESTART" {
				fmt.Println("Restart client...")
				break
			} else if data.Method == "EXIT" {
				fmt.Println("Goodbye.")
				return
			} else {
				apiData := c.GetApiDataByName(data.Method)
				if apiData.Commands == nil || len(apiData.Commands) == 0 {
					c.SendReq(data.Method)
				} else {
					for _, v := range apiData.Commands {
						fmt.Println(v)
						ok := c.SendReq(v)
						if !ok {
							break
						}
						res := c.WaitRes()
						if !res.Result {
							break
						}
						time.Sleep(time.Millisecond * 500)
					}
				}
			}
		}
	}
}
