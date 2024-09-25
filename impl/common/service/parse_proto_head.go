package service

//
//import (
//	"dawn-server/impl/common/msg"
//	xrerror "dawn-server/impl/xr/lib/error"
//	xrlog "dawn-server/impl/xr/lib/log"
//	xrutil "dawn-server/impl/xr/lib/util"
//	"github.com/pkg/errors"
//)
//
//// OnCheckLengthFromUser 处理客户端的包.
//func OnCheckLengthFromUser(length uint32) error {
//	if length < msg.GCSProtoHeadLength { //长度不足一个包头的长度大小
//		return errors.WithMessagef(xrerror.Packet, "%+v packetLength:%v", xrutil.GetCodeLocation(1).String(), length)
//	}
//	if GMgr.Bench.Json.Base.PacketLengthMax < length {
//		return errors.WithMessagef(xrerror.Packet, "%+v PacketLengthMax:%v, packetLength:%v",
//			xrutil.GetCodeLocation(1).String(), GMgr.Bench.Json.Base.PacketLengthMax, length)
//	}
//	return nil
//}
//
//// OnCheckLengthKCPUser 处理客户端的包-KCP
//func OnCheckLengthKCPUser(data []byte, length int) int {
//	if uint32(length) < msg.GCSProtoHeadLength { //长度不足一个包头的长度大小
//		return 0
//	}
//	packetLength := int(msg.GetPacketLength(data))
//	if uint32(packetLength) < msg.GCSProtoHeadLength {
//		xrlog.GetInstance().Errorf("packetLength:%v", packetLength)
//		return -1
//	}
//	if GMgr.Bench.Json.Base.PacketLengthMax < uint32(packetLength) {
//		xrlog.GetInstance().Errorf("PacketLengthMax:%v, packetLength:%v, data:%v",
//			GMgr.Bench.Json.Base.PacketLengthMax, packetLength, data[:msg.GCSProtoHeadLength])
//		return -1
//	}
//
//	if length < packetLength {
//		return 0
//	}
//
//	return packetLength
//}
