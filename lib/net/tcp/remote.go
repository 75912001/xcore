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
	xconstants "xcore/lib/constants"
	xerror "xcore/lib/error"
	xlog "xcore/lib/log"
	xnetpacket "xcore/lib/net/packet"
	xpool "xcore/lib/pool"
	xruntime "xcore/lib/runtime"
	xutil "xcore/lib/util"
	"xcore/lib/xswitch"
)

type IRemote interface {
	ISend
	IsConnect() bool
	Stop()
	GetIP() string
	xswitch.ISwitch // 是否-可用. (己方-主动断开连接)
}

// defaultRemote 远端
type defaultRemote struct {
	Conn            *net.TCPConn     // 连接
	sendChan        chan interface{} // 发送管道
	cancelFunc      context.CancelFunc
	Object          interface{} // 保存 应用层数据
	xswitch.ISwitch             // 是否-己方-主动断开连接
}

func NewDefaultRemote(Conn *net.TCPConn, sendChan chan interface{}) *defaultRemote {
	defaultRemote := &defaultRemote{
		Conn:     Conn,
		sendChan: sendChan,
		ISwitch:  xswitch.NewDefaultSwitch(true),
	}
	return defaultRemote
}

// GetIP 获取IP地址
func (p *defaultRemote) GetIP() string {
	slice := strings.Split(p.Conn.RemoteAddr().String(), ":")
	if len(slice) < 1 {
		return ""
	}
	return slice[0]
}

func (p *defaultRemote) start(tcpOptions *connOptions, event IEvent, handler IHandler) {
	//if err = p.Conn.SetKeepAlive(true); err != nil {
	//	log.Printf("SetKeepAlive war:%v", err)
	//}
	//if err = p.Conn.SetKeepAlivePeriod(time.Second * 10); err != nil {
	//	log.Printf("SetKeepAlivePeriod war:%v", err)
	//}
	//if err := p.Conn.SetNoDelay(true); err != nil {
	//	xrlog.PrintfErr("SetNoDelay war:%v", err)
	//}
	if tcpOptions.readBuffer != nil {
		if err := p.Conn.SetReadBuffer(*tcpOptions.readBuffer); err != nil {
			xlog.PrintfErr("SetReadBuffer err:%v", err)
		}
	}
	if tcpOptions.writeBuffer != nil {
		if err := p.Conn.SetWriteBuffer(*tcpOptions.writeBuffer); err != nil {
			xlog.PrintfErr("SetWriteBuffer err:%v", err)
		}
	}
	ctx := context.Background()
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)
	p.cancelFunc = cancelFunc

	go p.onSend(ctxWithCancel)
	go p.onRecv(event, handler)
}

// IsConnect 是否连接
func (p *defaultRemote) IsConnect() bool {
	return nil != p.Conn
}

// Send 发送数据
//
//	[NOTE]必须在处理 EventChan 事件中调用
//	参数:
//		packet: 未序列化的包. [NOTE]该数据会被引用,使用层不可写
func (p *defaultRemote) Send(packet xnetpacket.IPacket) error {
	if !p.IsConnect() {
		return errors.WithMessage(xerror.Link, xruntime.Location())
	}
	p.sendChan <- packet
	return nil
}

func (p *defaultRemote) Stop() {
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
func (p *defaultRemote) updateWriteDeadline(lastTime *time.Time, thisTime time.Time, writeTimeOutDuration time.Duration) error {
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
func (p *defaultRemote) push2Data(packet xnetpacket.IPacket, data []byte) ([]byte, error) {
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
func (p *defaultRemote) onSend(ctx context.Context) {
	defer func() {
		// 当 Conn 关闭, 该函数会引发 panic
		if err := recover(); err != nil {
			xlog.PrintErr(xconstants.GoroutinePanic, p, err, debug.Stack())
		}
		xlog.PrintInfo(xconstants.GoroutineDone, p)
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
			data, err = p.push2Data(t.(xnetpacket.IPacket), data)
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
					data, err = p.push2Data(t.(xnetpacket.IPacket), data)
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

const MsgLengthFieldSize uint32 = 4 // 消息总长度字段 的 大小

// 处理接收
func (p *defaultRemote) onRecv(event IEvent, handler IHandler) {
	defer func() { // 断开链接
		// 当 Conn 关闭, 该函数会引发 panic
		if err := recover(); err != nil {
			xlog.PrintErr(xconstants.GoroutinePanic, p, err, debug.Stack())
		}
		err := event.Disconnect(p)
		if err != nil {
			xlog.PrintfErr("disconnect err:%v", err)
		}
		xlog.PrintInfo(xconstants.GoroutineDone, p)
	}()
	// 消息总长度
	msgLengthBuf := make([]byte, MsgLengthFieldSize)
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
		if _, err := io.ReadFull(p.Conn, buf[MsgLengthFieldSize:]); err != nil {
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
		_ = event.Packet(p, packet)
	}
}
