package service

import (
	"time"
)

type OnHandlerBus func(v interface{}) error

// HandlerBus todo [重要] issue 在处理 event 时候, 向 eventChan 中插入 事件，注意超出eventChan的上限会阻塞.
func (p *DefaultService) HandlerBus() {
	// 在消费eventChan时可能会往eventChan中写入事件，所以关闭服务时不能close eventChan（造成写入阻塞），通过定时检查eventChan大小来关闭
	for {
		time.Sleep(time.Second)
		//select {
		//case <-GBusChannelCheckChan:
		//	xrlog.GetInstance().Warn("receive GBusChannelCheckChan")
		//	if 0 == len(p.BusChannel) && IsServerStopping() {
		//		xrlog.GetInstance().Warn("server is stopping, stop consume GEventChan with length 0")
		//		return
		//	} else {
		//		xrlog.GetInstance().Warnf("server is stopping, waiting for consume GEventChan with length:%d", len(p.BusChannel))
		//	}
		//case v := <-p.BusChannel:
		//	//TODO [*] 应拿尽拿...
		//	GMgr.TimeMgr.Update()
		//	var err error
		//	// todo menglc 使用 映射关系来处理不同的事件类型,提高性能. 比如: 数组.
		//	switch t := v.(type) {
		//	//tcp
		//	case *tcp.EventConnect:
		//		err = t.Remote.Handler.OnConnect(t.Remote)
		//	case *tcp.DefaultPacket:
		//		if !t.Remote.IsConnect() {
		//			continue
		//		}
		//		err = t.Remote.Handler.OnPacket(t)
		//	case *tcp.EventDisconnect:
		//		if !t.Remote.IsConnect() {
		//			continue
		//		}
		//		err = t.Remote.Handler.OnDisconnect(t.Remote)
		//		//timer
		//	case *timer.Second:
		//		if t.IsValid() {
		//			t.Function(t.arg)
		//		}
		//	case *timer.Millisecond:
		//		if t.IsValid() {
		//			t.Function(t.arg)
		//		}
		//		//kcp server
		//	case *xrkcp.EventConnect:
		//		err = t.Remote.Server.GetOnEvent().OnConn(t.Remote)
		//	case *xrkcp.EventDisconnect:
		//		err = t.Remote.Server.GetOnEvent().OnDisconnect(t.Remote)
		//	case *xrkcp.Packet:
		//		if !t.Remote.IsConn() {
		//			continue
		//		}
		//		err = t.Remote.Server.GetOnEvent().OnPacket(t)
		//	case *xretcd.KV:
		//		err = xretcd.GetInstance().Handler(t.Key, t.Value)
		//	case *mq_nats.Packet:
		//		err = p.Opts.HandlerNats(t)
		//	default:
		//		if p.Opts.HandlerBus == nil {
		//			xrlog.GetInstance().Fatalf("non-existent event:%v %v", v, t)
		//		} else {
		//			err = p.Opts.HandlerBus(v)
		//		}
		//	}
		//
		//	if err != nil { // todo 将日志放在每个 case 中,实例化输出对应的数据...
		//		xrlog.PrintErr(v, err)
		//	}
		//
		//	if util.IsDebug() {
		//		dt := time.Now().Sub(GMgr.TimeMgr.Time).Milliseconds()
		//		if dt > 50 {
		//			xrlog.GetInstance().Warnf("cost time50: %v Millisecond with event type:%T", dt, v)
		//		} else if dt > 20 {
		//			xrlog.GetInstance().Warnf("cost time20: %v Millisecond with event type:%T", dt, v)
		//		} else if dt > 10 {
		//			xrlog.GetInstance().Warnf("cost time10: %v Millisecond with event type:%T", dt, v)
		//		}
		//	}
		//}
	}
}
