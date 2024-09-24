package common

import (
	"time"

	"github.com/google/uuid"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

//const LogAbsPath string = "/data/log" //日志绝对路径

////////////////////////////////////////////////////////////////////////////////////////////////////
//etcd 相关

// 各服务名称 注意此处名称需与bench.json里的ServiceName保持一致
const ServiceNameGateway string = "gateway"              //网关服
const ServiceNameLogin string = "login"                  //登录服
const ServiceNameWorld string = "world"                  //世界服
const ServiceNameMatchRoomList string = "match_roomlist" //匹配-房间列表方式
const ServiceNameRoom string = "room"                    //房间服
const ServiceNameBattleGateway string = "battle_gateway" //战斗网关
const ServiceNameBattleVerify string = "battle_verify"   //战斗校验
const ServiceNameBattleServer string = "battle_server"   //战斗计算
const ServiceNameGM string = "gm"                        //GM服务
const ServiceNameServerList string = "server_list"       //区服列表

////////////////////////////////////////////////////////////////////////////////////////////////////
//room 相关

// RoomUserMaxCnt 房间最大人数限制
const RoomUserMaxCnt = 4

// GetRoomListMaxCnt 获取房间列表,最大数量
const GetRoomListMaxCnt = 100

////////////////////////////////////////////////////////////////////////////////////////////////////
//心跳 相关

// TCPHeartBeatTimeOutSec 心跳时间 秒 // todo 可配置...
var TCPHeartBeatTimeOutSec int64 = 30 //原值30
// KCPHeartBeatTimeOutSec 心跳时间 秒 //todo 可配置...
var KCPHeartBeatTimeOutSec int64 = 30 //原值10

////////////////////////////////////////////////////////////////////////////////////////////////////
//服务配置参数-限流-阈值

// BenchJsonOverLoadRoomCountMaxDefault 房间数量 最大 默认值
const BenchJsonOverLoadRoomCountMaxDefault = 1000

// BenchJsonOverLoadUserNumberMaxDefault 用户数量 最大 默认值
const BenchJsonOverLoadUserNumberMaxDefault = 100000000 //todo [@] 测试时开启 最大值, 默认10000

////////////////////////////////////////////////////////////////////////////////////////////////////
//帧相关

// FrameWindowsSize 帧 窗口大小
const FrameWindowsSize uint32 = FrameOneSecondFrameCount * 5

// FrameBattleAverageSecond 一场战斗平均秒数
const FrameBattleAverageSecond = 600

// FrameTimeOutMillisecond 帧,超时,毫秒
const FrameTimeOutMillisecond = 50

// FrameOneSecondFrameCount 一秒 帧 数量
const FrameOneSecondFrameCount = 1000 / FrameTimeOutMillisecond

// FrameBattleAverageFrameCount 一场战斗 平均 帧数
const FrameBattleAverageFrameCount = FrameBattleAverageSecond * FrameOneSecondFrameCount

// FrameVerifyHashFrameInterval 帧 验证 hash 逻辑帧间隔 , 1000 为 1000帧. 0:不需要hash.
const FrameVerifyHashFrameInterval = 1000

// FrameVerifyHashMapInitCapacity 验证 hash  初始化 entityMap 容量(单场战斗平均帧数/多少帧间隔发送一次hash值)
const FrameVerifyHashMapInitCapacity = FrameBattleAverageFrameCount / FrameVerifyHashFrameInterval

// SnapShootInitCapacity 快照 初始化 容量 todo [*]需要客户端提供预估大小...
const SnapShootInitCapacity = 1024 * 1024

const FillFrameNumber = 10 // 填充帧数量 todo [*]需要客户端提供预估大小...
////////////////////////////////////////////////////////////////////////////////////////////////////
//杂

// LoginTokenDuration 登录token持续时间
var LoginTokenDuration = time.Second * 60

// ServiceInfoTimeOutSec 信息 超时时间 秒
const ServiceInfoTimeOutSec = 10

//限流

// RateLimitIntervalSec 限流 间隔 秒 // todo 可配置...
const RateLimitIntervalSec = 60

// RateLimitCount 限流 数量// todo 可配置...
const RateLimitCount = 120

// NameLengthMax 名字 长度 最大值
const NameLengthMax = 24

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
// zid: zone id
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
