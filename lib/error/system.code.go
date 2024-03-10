package error

// 错误信息
var errMap map[uint32]struct{}

var (
	// Success 成功
	Success = CreateObject(0x0000, "Success", "success")
	// Link 链接
	Link = CreateObject(0xf001, "Link", "link error")
	//// System 系统
	//System = CreateObject(0xf002, "System", "system error")
	//// Param 参数
	//Param = CreateObject(0xf003, "Param", "parameter error")
	//// Packet 数据包
	//Packet = CreateObject(0xf004, "Packet", "packet error")
	//// Timeout 超时
	//Timeout = CreateObject(0xf005, "Timeout", "time out")
	//// ChannelFull 通道 满
	//ChannelFull = CreateObject(0xf006, "ChannelFull", "channel full")
	//// ChannelEmpty 通道 空
	//ChannelEmpty = CreateObject(0xf007, "ChannelEmpty", "channel empty")
	//// OutOfRange 超出范围
	//OutOfRange = CreateObject(0xf008, "OutOfRange", "out of range")
	//// InvalidValue 无效数值
	//InvalidValue = CreateObject(0xf009, "InvalidValue", "invalid value")
	//// Conflict 冲突
	//Conflict = CreateObject(0xf00a, "Conflict", "conflict")
	//// TypeMismatch 类型不匹配
	//TypeMismatch = CreateObject(0xf00b, "TypeMismatch", "type mismatch")
	//// InvalidPointer 无效指针
	//InvalidPointer = CreateObject(0xf00c, "InvalidPointer", "invalid pointer")
	//// Level 等级
	//Level = CreateObject(0xf00d, "level", "level error")
	//// NonExistent 不存在
	//NonExistent = CreateObject(0xf00e, "NonExistent", "non-existent")
	//// Exists 存在
	//Exists = CreateObject(0xf00f, "Exists", "exists")
	//// Marshal 序列化
	//Marshal = CreateObject(0xf010, "Marshal", "marshal")
	//// Unmarshal 反序列化
	//Unmarshal = CreateObject(0xf011, "Unmarshal", "unmarshal")
	//// Insert 插入
	//Insert = CreateObject(0xf012, "Insert", "insert error")
	//// Find 查找
	//Find = CreateObject(0xf013, "Find", "find error")
	//// Update 更新
	//Update = CreateObject(0xf014, "Update", "update error")
	//// Delete 删除
	//Delete = CreateObject(0xf015, "Delete", "delete error")
	//// Duplicate 重复
	//Duplicate = CreateObject(0xf016, "Duplicate", "duplicate error")
	//// Config 配置
	//Config = CreateObject(0xf017, "Config", "config error")
	//// InvalidOperation 无效操作
	//InvalidOperation = CreateObject(0xf018, "InvalidOperation", "invalid operation")
	//// IllConditioned 条件不足
	//IllConditioned = CreateObject(0xf019, "IllConditioned", "ill conditioned")
	//// PermissionDenied 没有权限
	//PermissionDenied = CreateObject(0xf01a, "PermissionDenied", "permission denied")
	//// BlockedAccount 冻结账号
	//BlockedAccount = CreateObject(0xf01b, "BlockedAccount", "blocked account")
	//// Send 发送
	//Send = CreateObject(0xf01c, "Send", "send")
	//// Configure 给配置
	//Configure = CreateObject(0xf01d, "Configure", "configure")
	//// Retry 重试
	//Retry = CreateObject(0xf01e, "Retry", "retry")
	//// MessageIDNonExistent 消息ID 不存在
	//MessageIDNonExistent = CreateObject(0xf01f, "MessageIDNonExistent", "message id non-existent")
	//// Redis 系统 Redis
	//Redis = CreateObject(0xf020, "Redis", "redis")
	//// Busy 繁忙
	//Busy = CreateObject(0xf021, "Busy", "busy")
	//// OutOfResources 资源不足
	//OutOfResources = CreateObject(0xf022, "OutOfResources", "out of resources")
	//// NATS NATS错误
	//NATS = CreateObject(0xf023, "NATS", "nats")
	//// PacketQuantityLimit 包数量限制
	//PacketQuantityLimit = CreateObject(0xf024, "PacketQuantityLimit", "packet quantity limit")
	//// OverloadWarning 过载-告警
	//OverloadWarning = CreateObject(0xf025, "OverloadWarning", "overload warning")
	//// OverloadError 过载-错误
	//OverloadError = CreateObject(0xf026, "OverloadError", "overload error")
	//// MessageIDDisable 消息ID禁用
	//MessageIDDisable = CreateObject(0xf027, "MessageIDDisable", "message id is disabled")
	//// MessageIDExistent 消息ID 存在
	//MessageIDExistent = CreateObject(0xf028, "MessageIDExistent", "message id existent")
	//// ModeMismatch 模式 不匹配
	//ModeMismatch = CreateObject(0xf029, "ModeMismatch", "mode mismatch")
	//// FormatMismatch 格式 不匹配
	//FormatMismatch = CreateObject(0xf02a, "FormatMismatch", "format mismatch")
	//// MISSING 找不到,丢失,未命中
	//MISSING = CreateObject(0xf02b, "MISSING", "missing")
	//// VersionMismatch 版本 不匹配
	//VersionMismatch = CreateObject(0xf02c, "VersionMismatch", "version mismatch")
	//// Unavailable 不可用
	//Unavailable = CreateObject(0xf02d, "Unavailable", "unavailable")
	//// NotImplemented 未实现
	//NotImplemented = CreateObject(0xf02e, "NotImplemented", "not implemented")
	//// PacketHeaderLength 数据包头长度
	//PacketHeaderLength = CreateObject(0xf02f, "PacketHeaderLength", "packet header length error")
	//// ChannelClosed 通道 已关闭
	//ChannelClosed = CreateObject(0xf030, "ChannelClosed", "channel closed")
	//// Unregistered 未注册
	//Unregistered = CreateObject(0xf031, "Unregistered", "unregistered")
	// Unknown 未知
	Unknown = CreateObject(0xffff, "Unknown", "unknown")
	// 0xffff
)
