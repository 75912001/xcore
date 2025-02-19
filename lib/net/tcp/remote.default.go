package tcp

import (
	"context"
	"encoding/binary"
	"github.com/pkg/errors"
	"io"
	"net"
	"runtime/debug"
	"strings"
	"time"
	xerror "xcore/lib/error"
	xlog "xcore/lib/log"
	"xcore/lib/packet"
	xpool "xcore/lib/pool"
	xruntime "xcore/lib/runtime"
	xutil "xcore/lib/util"
)

// Remote 远端
type Remote struct {
	Conn       *net.TCPConn     // 连接
	sendChan   chan interface{} // 发送管道
	cancelFunc context.CancelFunc
	Object     interface{} // 保存 应用层数据
}

func NewRemote(Conn *net.TCPConn, sendChan chan interface{}) *Remote {
	return &Remote{
		Conn:     Conn,
		sendChan: sendChan,
	}
}

// GetIP 获取IP地址
func (p *Remote) GetIP() string {
	slice := strings.Split(p.Conn.RemoteAddr().String(), ":")
	if len(slice) < 1 {
		return ""
	}
	return slice[0]
}

func (p *Remote) Start(tcpOptions *ConnOptions, event IEvent, handler IHandler) {
	//var err error
	//if err = p.Conn.SetKeepAlive(true); err != nil {
	//	xlog.PrintfErr("SetKeepAlive err:%v", err)
	//}
	//if err = p.Conn.SetKeepAlivePeriod(time.Second * 600); err != nil {
	//	xlog.PrintfErr("SetKeepAlivePeriod err:%v", err)
	//}
	//if err := p.Conn.SetNoDelay(true); err != nil {
	//	xlog.PrintfErr("SetNoDelay err:%v", err)
	//}
	if tcpOptions.ReadBuffer != nil {
		if err := p.Conn.SetReadBuffer(*tcpOptions.ReadBuffer); err != nil {
			xlog.PrintfErr("WithReadBuffer err:%v", err)
		}
	}
	if tcpOptions.WriteBuffer != nil {
		if err := p.Conn.SetWriteBuffer(*tcpOptions.WriteBuffer); err != nil {
			xlog.PrintfErr("WithWriteBuffer err:%v", err)
		}
	}
	ctx := context.Background()
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)
	p.cancelFunc = cancelFunc

	go p.onSend(ctxWithCancel)
	go p.onRecv(event, handler)
}

// IsConnect 是否连接
func (p *Remote) IsConnect() bool {
	return nil != p.Conn
}

// Send 发送数据
//
//	[NOTE]必须在处理 EventChan 事件中调用
//	参数:
//		packet: 未序列化的包. [NOTE]该数据会被引用,使用层不可写
func (p *Remote) Send(packet packet.IPacket) error {
	if !p.IsConnect() {
		return errors.WithMessage(xerror.Link, xruntime.Location())
	}
	p.sendChan <- packet
	return nil
}

func (p *Remote) Stop() {
	if p.IsConnect() {
		err := p.Conn.Close()
		if err != nil {
			xlog.PrintfErr("connect close err:%v", err)
		}
	}
	if p.cancelFunc != nil {
		p.cancelFunc()
	}
	p.cancelFunc = nil
	p.Conn = nil
}

// 写超时
//
//	只有超过50%时才更新写截止日期
//	参数:
//		lastTime:上次时间 (可能会更新)
//		thisTime:这次时间
//		writeTimeOutDuration:写超时时长
func (p *Remote) updateWriteDeadline(lastTime *time.Time, thisTime time.Time, writeTimeOutDuration time.Duration) error {
	if (writeTimeOutDuration >> 1) < thisTime.Sub(*lastTime) {
		if err := p.Conn.SetWriteDeadline(thisTime.Add(writeTimeOutDuration)); err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
		*lastTime = thisTime
	}
	return nil
}

// 重新整理-待发送数据
func rearrangeSendData(data []byte, cnt int, resetCnt int) []byte {
	if len(data) == cnt {
		if resetCnt <= cap(data) { // 占用空间过大,重新创建新的数据单元
			data = []byte{}
		} else {
			data = data[0:0]
		}
	} else {
		data = data[cnt:]
	}
	return data
}

// 将数据放入data中
func (p *Remote) push2Data(packet packet.IPacket, data []byte) ([]byte, error) {
	packetData, err := packet.Marshal()
	if err != nil {
		xlog.PrintfErr("packet marshal %v", packet)
		return nil, errors.WithMessage(err, xruntime.Location())
	}
	if len(data) == 0 { //当 data len == 0 时候, 直接发送 v.data 数据...
		data = packetData
	} else {
		data = append(data, packetData...)
	}
	return data, nil
}

// 处理发送
func (p *Remote) onSend(ctx context.Context) {
	defer func() {
		// 当 Conn 关闭, 该函数会引发 panic
		if err := recover(); err != nil {
			xlog.PrintErr(xerror.GoroutinePanic, p, err, debug.Stack())
		}
		xlog.PrintInfo(xerror.GoroutineDone, p)
	}()
	// 上次时间
	var lastTime time.Time
	// 这次时间
	var thisTime time.Time
	// 超时
	writeTimeOutDuration := time.Millisecond * 100
	var data []byte // 待发送数据
	var err error
	var writeCnt int
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-p.sendChan:
			data, err = p.push2Data(t.(packet.IPacket), data)
			if err != nil {
				xlog.PrintfErr("push2Data err:%v", err)
				continue
			}
			for {
				thisTime = time.Now()
				// 超时, 防止, 客户端不 read 数据, 导致主循环无法写入.到p.sendChan中, 阻塞主循环逻辑
				if err := p.updateWriteDeadline(&lastTime, thisTime, writeTimeOutDuration); err != nil {
					xlog.PrintfErr("updateWriteDeadline remote:%p err:%v", p, err)
				}
				writeCnt, err = p.Conn.Write(data)
				if 0 < writeCnt {
					data = rearrangeSendData(data, writeCnt, 10240)
					if len(data) == 0 {
						break
					} else {
						xlog.PrintfErr("Conn.Write remote:%p writeCnt:%v remaining:%v", p, writeCnt, len(data))
					}
				}
				for 0 < len(p.sendChan) { // 尽量取出待发送数据
					t := <-p.sendChan
					data, err = p.push2Data(t.(packet.IPacket), data)
					if err != nil {
						xlog.PrintfErr("push2Data err:%v", err)
						continue
					}
				}
				if nil != err {
					if xutil.IsNetErrorTimeout(err) { // 网络超时
						xlog.PrintfErr("Conn.Write timeOut. remote:%p writeCnt:%v remaining:%v err:%v",
							p, writeCnt, len(data), err)
						continue
					}
					xlog.PrintfErr("Conn.Write remote:%p writeCnt:%v remaining:%v err:%v", p, writeCnt, len(data), err)
					break
				}
			}
		}
	}
}

// 处理接收
func (p *Remote) onRecv(event IEvent, handler IHandler) {
	defer func() { // 断开链接
		// 当 Conn 关闭, 该函数会引发 panic
		if err := recover(); err != nil {
			xlog.PrintErr(xerror.GoroutinePanic, p, err, debug.Stack())
		}
		err := event.Disconnect(handler, p)
		if err != nil {
			xlog.PrintfErr("disconnect err:%v", err)
		}
		xlog.PrintInfo(xerror.GoroutineDone, p)
	}()
	// 消息总长度
	msgLengthBuf := make([]byte, packet.HeaderLengthFieldSize)
	for {
		if _, err := io.ReadFull(p.Conn, msgLengthBuf); err != nil {
			if !xutil.IsNetErrClosing(err) {
				xlog.PrintfInfo("remote:%p err:%v", p, err)
			}
			return
		}
		packetLength := binary.LittleEndian.Uint32(msgLengthBuf)
		if err := handler.OnCheckPacketLength(packetLength); err != nil {
			xlog.PrintfErr("remote:%p OnCheckPacketLength err:%v", p, err)
			return
		}
		buf := xpool.MakeByteSlice(int(packetLength))
		copy(buf, msgLengthBuf)
		if _, err := io.ReadFull(p.Conn, buf[packet.HeaderLengthFieldSize:]); err != nil {
			xlog.PrintfErr("remote:%p err:%v", p, err)
			_ = xpool.ReleaseByteSlice(buf)
			return
		}
		if err := handler.OnCheckPacketLimit(p); err != nil {
			xlog.PrintfErr("remote:%p buf:%v err:%v", p, buf, err)
			_ = xpool.ReleaseByteSlice(buf)
			continue
		}
		packet, err := handler.OnUnmarshalPacket(p, buf)
		_ = xpool.ReleaseByteSlice(buf)
		if err != nil {
			xlog.PrintfErr("remote:%p buf:%v err:%v", p, buf, err)
			continue
		}
		buf = nil
		_ = event.Packet(handler, p, packet)
	}
}
