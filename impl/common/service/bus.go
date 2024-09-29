package service

import (
	"time"
	xerror "xcore/lib/error"
	"xcore/lib/net/tcp"
)

// HandleEvent todo [重要] issue 在处理 event 时候, 向 eventChan 中插入 事件，注意超出eventChan的上限会阻塞.
func OnHandlerBus() error {
	//在消费eventChan时可能会往eventChan中写入事件，所以关闭服务时不能close eventChan（造成写入阻塞），通过定时检查eventChan大小来关闭
	for {
		select {
		case <-GBusChannelCheckChan:
			xrlog.GetInstance().Warn("receive GBusChannelCheckChan")
			if 0 == len(eventChan) && IsServerStopping() {
				xrlog.GetInstance().Warn("server is stopping, stop consume GEventChan with length 0")
				return
			} else {
				xrlog.GetInstance().Warnf("server is stopping, waiting for consume GEventChan with length:%d", len(eventChan))
			}
		case v := <-eventChan:
			//TODO [*] 应拿尽拿...
			GMgr.TimeMgr.Update()
			var err error
			switch t := v.(type) {
			//tcp
			case *tcp.EventConnect:
				err = t.IHandler.OnConnect(t.IRemote)
			case *tcp.EventPacket:
				err = t.IHandler.OnPacket(t.IPacket)
			case *tcp.EventDisconnect:
				err = t.IHandler.OnDisconnect(t.IRemote)
				if t.IRemote.IsConnect() {
					t.IRemote.Stop()
				}
				//timer
			case *timer.Second:
				if t.IsValid() {
					t.Function(t.Arg)
				}
			case *timer.Millisecond:
				if t.IsValid() {
					t.Function(t.Arg)
				}
				//kcp server
			case *xrkcp.EventConnect:
				err = t.Remote.Server.GetOnEvent().OnConn(t.Remote)
			case *xrkcp.EventDisconnect:
				err = t.Remote.Server.GetOnEvent().OnDisconnect(t.Remote)
			case *xrkcp.Packet:
				if !t.Remote.IsConn() {
					continue
				}
				err = t.Remote.Server.GetOnEvent().OnPacket(t)
			case *xretcd.KV:
				err = xretcd.GetInstance().Handler(t.Key, t.Value)
			case *mq_nats.Packet:
				err = onNatsFunc(t)
			default:
				if onDefaultFunc == nil {
					xrlog.GetInstance().Fatalf("non-existent event:%v %v", v, t)
				} else {
					err = onDefaultFunc(v)
				}
			}

			if err != nil { // todo 将日志放在每个 case 中,实例化输出对应的数据...
				xrlog.PrintErr(v, err)
			}

			if util.IsDebug() {
				dt := time.Now().Sub(GMgr.TimeMgr.Time).Milliseconds()
				if dt > 50 {
					xrlog.GetInstance().Warnf("cost time50: %v Millisecond with event type:%T", dt, v)
				} else if dt > 20 {
					xrlog.GetInstance().Warnf("cost time20: %v Millisecond with event type:%T", dt, v)
				} else if dt > 10 {
					xrlog.GetInstance().Warnf("cost time10: %v Millisecond with event type:%T", dt, v)
				}
			}
		}
	}
	return xerror.NotImplemented
}
