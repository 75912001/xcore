package common

import (
	xlog "github.com/75912001/xcore/lib/log"
	xpacket "github.com/75912001/xcore/lib/packet"
	xruntime "github.com/75912001/xcore/lib/runtime"
	"github.com/pkg/errors"
)

// 将数据 packet 放到 data 中
func PushPacket2Data(data []byte, packet xpacket.IPacket) ([]byte, error) {
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

// RearRangeData 重新整理-数据
// removeLen: 移除长度
// resetCnt: 重置长度, 大于该长度则重新创建新的数据单元
func RearRangeData(dataSlice []byte, removeLen int, resetCnt int) []byte {
	if len(dataSlice) == removeLen {
		if resetCnt <= cap(dataSlice) { // 占用空间过大,重新创建新的数据单元
			dataSlice = []byte{}
		} else {
			dataSlice = dataSlice[0:0]
		}
	} else {
		dataSlice = dataSlice[removeLen:]
	}
	return dataSlice
}
