package common

import (
	"time"

	"github.com/google/uuid"
)

const ServerNameGateway string = "gateway" // 网关服
const ServerNameLogin string = "login"     // 登录服
const ServerNameLogic string = "logic"     // 逻辑服
const ServerNameRoom string = "room"       // 房间服
const ServerNameGM string = "gm"           // GM服务

var TCPHeartBeatTimeOutSec int64 = 30 // 心跳时间-秒 // todo menglc 可配置...
var KCPHeartBeatTimeOutSec int64 = 30 // 心跳时间-秒 // todo menglc 可配置...

var LoginTokenDuration = time.Second * 60 // LoginTokenDuration 登录 token 有效持续时间

const RateLimitIntervalSec = 60 // RateLimitIntervalSec 限流-间隔-秒 // todo 可配置...

const RateLimitCount = 120 // RateLimitCount 限流-数量// todo 可配置...

const NameLengthMax = 24 // NameLengthMax 名字 长度 最大值

//MongoDB 异步合并写

// BulkWriteMax DB合并写单个集合最大合批数量
const BulkWriteMax = 1000

// BatchExecMaxMilliSecond  DB合并写周期 单位毫秒
const BatchExecMaxMilliSecond = 30000

// //////////////////////////////////////////////////////////////////////////////////////////////////
// Mongodb
// var MongodbDatabaseNameDefault = ProjectName
var MongodbMaxPoolSizeDefault uint64 = 24 // 可以设置为 2*P + 30% 的数量
var MongodbMinPoolSizeDefault uint64 = 8  // 可以设置为 P 的数量
var MongodbTimeoutDurationDefault = time.Minute
var MongodbMaxConnIdleTimeDefault = time.Minute * 5
var MongodbMaxConnectingDefault uint64 = 4 // 设置为 P/2

////////////////////////////////////////////////////////////////////////////////////////////////////

//func init() {
//	if util.IsDebug() {
//		//HeartBeatTimeOutSec = 60 * 60 * 24
//		LoginTokenDuration = time.Hour
//		//etcd.TtlSecondDefault = 60 * 60 * 24
//		//etcd.KeepAliveSecond = 60 * 60
//	}
//}

// GenUID 生成UID
// zid: group id
// uidIdx: uid 在 redis 中的序号
func GenUID(zid uint32, uidIdx uint64) uint64 {
	return uint64(zid)*uint64(1000000000) + uidIdx
}

// GenUIDIdx 生成UID序号
func GenUIDIdx(uid uint64) uint64 {
	return uid % uint64(1000000000)
}

// GetZIDByUID 获取ZID 通过UID
func GetZIDByUID(uid uint64) uint32 {
	return uint32(uid / 1000000000)
}

// IsZoneSet 是否为同一组区域  比如zoneID为2的world和zoneID为20001的battle_gateway属于同一组区域
func IsZoneSet(zone1 uint32, zone2 uint32) bool {
	if zone1 == zone2 {
		return true
	}

	if zone1/10000 == zone2 || zone2/10000 == zone1 {
		return true
	}

	// 兼容内网测试服 id为10000+N的规则
	if zone1 == 10000+zone2 || zone2 == 10000+zone1 {
		return true
	}

	return false
}

// GetZoneSetPrefix 取得zoneID关联的所有zoneID前缀 用于清档删除etcd数据
//func GetZoneSetPrefix(zoneID uint32) []string {
//	var prefixSlice []string
//
//	var zoneIDShort uint32
//	if zoneID >= 10000 {
//		zoneIDShort = zoneID / 10000
//	} else {
//		zoneIDShort = zoneID
//	}
//
//	// world / login / gm
//	prefixSlice = append(prefixSlice, fmt.Sprintf("/%v/%v/%d/", ProjectName, EtcdWatchMsgTypeService, zoneIDShort))
//	prefixSlice = append(prefixSlice, fmt.Sprintf("/%v/%v/%d/", ProjectName, EtcdWatchMsgTypeCommand, zoneIDShort))
//
//	// battle_gateway / room / match_room_list / battle_verify
//	zoneIDStart := zoneIDShort*10000 + 1
//	zoneIDEnd := zoneIDShort*10000 + 9999
//
//	zoneIDRange := fmt.Sprintf("/%v/%v/%d-%d/", ProjectName, EtcdWatchMsgTypeService, zoneIDStart, zoneIDEnd) // 比如world为2 battle_gateway则为 20001 ~ 29999之间
//	prefixSlice = append(prefixSlice, zoneIDRange)
//	zoneIDRange = fmt.Sprintf("/%v/%v/%d-%d/", ProjectName, EtcdWatchMsgTypeCommand, zoneIDStart, zoneIDEnd) // 比如world为2 battle_gateway则为 20001 ~ 29999之间
//	prefixSlice = append(prefixSlice, zoneIDRange)
//
//	return prefixSlice
//}

// GenLogTraceID 生成日志traceID
func GenLogTraceID() string {
	genUUID, _ := uuid.NewRandom()
	return genUUID.String()
}

// UUIDGenerator 生成用户UUID的接口
type UUIDGenerator interface {
	GenUUID() (uint64, error)
}

const DBRetryCountMax = 10000

// GetMessageNameByID 根据proto定义的消息ID取得消息名称
//func GetMessageNameByID(cmd uint32) string {
//	cmdMapStr, ok := protobuf.CMDMap[cmd]
//	if !ok {
//		return fmt.Sprintf("%#x", cmd)
//	}
//	return cmdMapStr
//}
