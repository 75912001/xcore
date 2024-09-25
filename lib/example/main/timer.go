package main

import (
	"context"
	"fmt"
	"time"
	xutil "xcore/lib/callback"
	xconstants "xcore/lib/constants"
	xtimer "xcore/lib/timer"
)

func cbSecond(arg interface{}) error {
	fmt.Println("cbSecond:", arg.(uint64))
	return nil
}

func cbMillisecond(arg interface{}) error {
	fmt.Println("cbMillisecond:", arg.(uint64))
	return nil
}

type addSecondSignal struct {
}
type addMillisecondSignal struct {
}

func exampleTimer() {
	if true {
		return
	}
	var timer xtimer.ITimer
	timer = xtimer.NewMgr()
	busChannel := make(chan interface{}, xconstants.BusChannelCapacityDefault)
	err := timer.Start(context.Background(),
		xtimer.NewOptions().
			WithOutgoingTimerOutChan(busChannel),
	)
	if err != nil {
		panic(err)
	}

	busChannel <- addSecondSignal{}
	busChannel <- addMillisecondSignal{}
	for {
		select {
		case v := <-busChannel:
			switch t := v.(type) {
			case addSecondSignal:
				for i := 0; i < 10; i++ {
					defaultCallBack := xutil.NewCallBack(cbSecond, uint64(i))
					second := timer.AddSecond(defaultCallBack, time.Now().Unix()+int64(i))
					switch i {
					case 3, 7, 9:
						timer.DelSecond(second)
					default:
					}
				}
			case addMillisecondSignal:
				for i := 0; i < 10000; i += 1000 {
					defaultCallBack := xutil.NewCallBack(cbMillisecond, uint64(i))
					millisecond := timer.AddMillisecond(defaultCallBack, time.Now().UnixMilli()+int64(i))
					switch i {
					case 3000, 7000, 9000:
						timer.DelMillisecond(millisecond)
					default:
					}
				}
			case xtimer.EventTimerSecond:
				_ = t.Function()
			case xtimer.EventTimerMillisecond:
				_ = t.Function()
			}
		}
	}
	return
}
