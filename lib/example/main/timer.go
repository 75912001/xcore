package main

import (
	"context"
	"fmt"
	"time"
	xtimer "xcore/lib/timer"
	xutil "xcore/lib/util"
)

func cbSecond(arg interface{}) {
	fmt.Println("cbSecond:", arg.(uint64))
}
func exampleTimer() {
	if false {
		return
	}
	var timer xtimer.ITimer
	timer = xtimer.NewMgr()
	busChannel := make(chan interface{}, 100)
	err := timer.Start(context.Background(),
		xtimer.NewOptions().
			WithOutgoingTimerOutChan(busChannel),
	)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 3)
	for i := 0; i < 10; i++ {
		defaultCallBack := xutil.NewDefaultCallBack(cbSecond, uint64(i))
		second := timer.AddSecond(defaultCallBack, int64(i))
		_ = second
		//switch i {
		//case 30, 70, 90:
		//	timer.DelSecond(second)
		//default:
		//}
	}
	for {
		select {
		case v := <-busChannel:
			switch t := v.(type) {
			case xtimer.Second:
				_ = t.CallBackFunc()
			}
		}
	}
	for {
		time.Sleep(time.Second * 1)
	}
	return
}
