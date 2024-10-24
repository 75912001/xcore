package gateway

import (
	"github.com/pkg/errors"
	xcommon "xcore/impl/common"
	xconstants "xcore/lib/constants"
	xerror "xcore/lib/error"
	xetcd "xcore/lib/etcd"
	xlog "xcore/lib/log"
)

// EtcdKeyValue etcd 刷新 key value
func EtcdKeyValue(arg ...interface{}) error {
	key := arg[0].(string)
	valueJson := arg[1].(*xetcd.ValueJson)
	if valueJson == nil { // 失效
		xlog.PrintfInfo("etcd key:%v, 失效", key)
		return nil
	}
	xlog.PrintfInfo("etcd key:%v, value:%v", key, valueJson)
	msgType, groupID, serviceName, serviceID := xetcd.Parse(key)
	xlog.PrintfInfo("msgType:%v, groupID:%v, serviceName:%v, serviceID:%v", msgType, groupID, serviceName, serviceID)
	if groupID != gservice.GroupID {
		return nil
	}
	if serviceName == gservice.Name {
		return nil
	}
	switch msgType {
	case xconstants.EtcdWatchMsgTypeService:
		switch serviceName {
		case xcommon.ServiceNameLogin:
			loginService := gservice.LoginServiceMgr.Get(serviceID)
			if loginService != nil {
				xlog.PrintfInfo("login service Existent %v", serviceID)
				return errors.WithMessage(xerror.Existent, "login services Existent")
			}
			loginService = NewLoginService()
			//address := fmt.Sprintf("%v:%v", etcdValueJson.ServiceNetTCP.IP, etcdValueJson.ServiceNetTCP.Port)
			//
			//if err := room.Client.Connect(context.TODO(),
			//	xrtcp.NewClientOptions().
			//		SetAddress(address).
			//		SetEventChan(server.GMgr.BusChannel).
			//		SetSendChanCapacity(server.GMgr.Bench.Json.Base.SendChanCapacity*100). // bench里的SendChanCapacity是单个用户的消息容量 此处乘以100
			//		SetPacket(&msg.SSPacket{}).
			//		SetOnUnmarshalPacket(OnUnmarshalPacketFromRoom).
			//		SetOnPacket(OnPacketFromRoom).
			//		SetOnDisconnect(OnDisconnectFromRoom),
			//); err != nil {
			//	xrlog.GetInstance().Fatalf("%v zoneID:%v, serviceName:%v, serviceID:%v, ip:%v, port:%v connect service err:%v",
			//		xrconstant.Etcd, zoneIDU32, serviceName, serviceIDU32, etcdValueJson.ServiceNetTCP.IP,
			//		etcdValueJson.ServiceNetTCP.Port, err)
			//	return errors.WithMessage(err, "connect room service err")
			//}
			//xrlog.GetInstance().Debug(address, zoneIDU32, serviceName, serviceIDU32, &room.Client.Remote)
			//room.ID = serviceIDU32
			//room.Client.Remote.Object = room
			//room_mgr.GMgr.Add(room)
			//
			//err := room.Client.Remote.Send(
			//	&xrpb.UnserializedPacket{
			//		CTX: context.TODO(),
			//		Header: &msg.SSHeader{
			//			MessageID: room_proto.RoomBattleGatewayRegisterMsg_CMD,
			//		},
			//		Message: &room_proto.RoomBattleGatewayRegisterMsg{
			//			ID: server.GMgr.ServiceID,
			//		},
			//	},
			//)
			//if err != nil {
			//	xrlog.GetInstance().Warnf("%v zoneID:%v, serviceName:%v, serviceID:%v, ip:%v, port:%v connect service err:%v",
			//		xrconstant.Etcd, zoneIDU32, serviceName, serviceIDU32, etcdValueJson.ServiceNetTCP.IP,
			//		etcdValueJson.ServiceNetTCP.Port, err)
			//	return errors.WithMessage(err, "register room service err")
			//}
		default:
		}
	default:
	}
	return nil
}
