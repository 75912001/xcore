package service

import (
	"time"
	xlog "xcore/lib/log"
	xnetevent "xcore/lib/net/event"
	xruntime "xcore/lib/runtime"
	xtimer "xcore/lib/timer"
)

// HandleEvent todo [重要] issue 在处理 event 时候, 向 eventChan 中插入 事件，注意超出eventChan的上限会阻塞.
func (p *DefaultService) Handle() error {
	//在消费eventChan时可能会往eventChan中写入事件，所以关闭服务时不能close eventChan（造成写入阻塞），通过定时检查eventChan大小来关闭
	for {
		select {
		//case <-GBusChannelCheckChan:
		//	xrlog.GetInstance().Warn("receive GBusChannelCheckChan")
		//	if 0 == len(eventChan) && IsServerStopping() {
		//		xrlog.GetInstance().Warn("server is stopping, stop consume GEventChan with length 0")
		//		return
		//	} else {
		//		xrlog.GetInstance().Warnf("server is stopping, waiting for consume GEventChan with length:%d", len(eventChan))
		//	}
		case value := <-p.BusChannel:
			//TODO [*] 应拿尽拿...
			p.TimeMgr.Update()
			var err error
			switch event := value.(type) {
			//tcp
			case *xnetevent.Connect:
				err = event.IHandler.OnConnect(event.IRemote)
			case *xnetevent.Packet:
				err = event.IHandler.OnPacket(event.IPacket)
			case *xnetevent.Disconnect:
				err = event.IHandler.OnDisconnect(event.IRemote)
				if event.IRemote.IsConnect() {
					event.IRemote.Stop()
				}
				//timer
			case *xtimer.EventTimerSecond:
				if event.IsDisabled() {
					continue
				}
				_ = event.Execute()
			case *xtimer.EventTimerMillisecond:
				if event.IsDisabled() {
					continue
				}
				_ = event.Execute()
				//kcp server
			//case *xrkcp.EventConnect:
			//	err = event.Remote.Server.GetOnEvent().OnConn(event.Remote)
			//case *xrkcp.EventDisconnect:
			//	err = event.Remote.Server.GetOnEvent().OnDisconnect(event.Remote)
			//case *xrkcp.Packet:
			//	if !event.Remote.IsConn() {
			//		continue
			//	}
			//	err = event.Remote.Server.GetOnEvent().OnPacket(event)
			//case *xretcd.KV:
			//	err = xretcd.GetInstance().Handler(event.Key, event.Value)
			//case *mq_nats.Packet:
			//	err = onNatsFunc(event)
			default:
				//	xrlog.GetInstance().Fatalf("non-existent event:%value %value", value, event)
				//	if onDefaultFunc == nil {
				//} else {
				//	err = onDefaultFunc(value)
				//}
			}
			if err != nil {
				p.Log.Errorf("Handle event:%v error:%value", value, err)
			}

			if xruntime.IsDebug() {
				dt := time.Now().Sub(p.TimeMgr.NowTime()).Milliseconds()
				if dt > 50 {
					xlog.PrintfErr("cost time50: %value Millisecond with event type:%T", dt, value)
				} else if dt > 20 {
					xlog.PrintfErr("cost time20: %value Millisecond with event type:%T", dt, value)
				} else if dt > 10 {
					xlog.PrintfErr("cost time10: %value Millisecond with event type:%T", dt, value)
				}
			}
		}
	}
	return nil
}
