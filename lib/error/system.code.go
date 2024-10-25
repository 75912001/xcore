package error

var errorCode uint32 = 0x0000

// 获取当前错误码, 并且自增
func getErrorCode() uint32 {
	code := errorCode
	errorCode++
	return code
}

// [系统错误码] lib/system 级别的错误码
var (
	Success = NewError(0x0000).WithName("Success").WithDesc("success-成功")
	// 通用
	Fail           = NewError(0x0101).WithName("Fail").WithDesc("fail-失败")
	Valid          = NewError(0x0102).WithName("Valid").WithDesc("valid-有效")
	Invalid        = NewError(0x0103).WithName("Invalid").WithDesc("invalid-无效")
	Available      = NewError(0x0104).WithName("Available").WithDesc("available-可用")
	Unavailable    = NewError(0x0105).WithName("Unavailable").WithDesc("unavailable-不可用")
	Existent       = NewError(0x0106).WithName("Exists").WithDesc("exists-存在")
	NonExistent    = NewError(0x0107).WithName("NonExistent").WithDesc("nonexistent-不存在")
	Legal          = NewError(0x0108).WithName("Legal").WithDesc("legal-合法")
	Illegal        = NewError(0x0109).WithName("Illegal").WithDesc("illegal-非法")
	Permitted      = NewError(0x010a).WithName("Permitted").WithDesc("permitted-允许的")
	Prohibited     = NewError(0x010b).WithName("Prohibited").WithDesc("prohibited-禁止的")
	Expect         = NewError(0x010c).WithName("Expect").WithDesc("expect-期望")
	Unexpected     = NewError(0x010d).WithName("Unexpected").WithDesc("unexpected-不期望")
	Enable         = NewError(0x010e).WithName("Enable").WithDesc("enable-启用")
	Disable        = NewError(0x010f).WithName("Disable").WithDesc("disable-禁用")
	Normal         = NewError(0x0110).WithName("Normal").WithDesc("normal-正常")
	Abnormal       = NewError(0x0111).WithName("Abnormal").WithDesc("abnormal-异常")
	NotTimeout     = NewError(0x0112).WithName("NotTimeout").WithDesc("not-timeout-未超时")
	Timeout        = NewError(0x0113).WithName("Timeout").WithDesc("timeout-超时")
	NotOutOfRange  = NewError(0x0114).WithName("NotOutOfRange").WithDesc("not-out-of-range-未超出范围")
	OutOfRange     = NewError(0x0115).WithName("OutOfRange").WithDesc("out-of-range-超出范围")
	NotConflict    = NewError(0x0116).WithName("NotConflict").WithDesc("not-conflict-未冲突")
	Conflict       = NewError(0x0117).WithName("Conflict").WithDesc("conflict-冲突")
	Matched        = NewError(0x0118).WithName("Matched").WithDesc("matched-匹配")
	Mismatch       = NewError(0x0119).WithName("Mismatch").WithDesc("mismatch-不匹配")
	Implemented    = NewError(0x011a).WithName("Implemented").WithDesc("implemented-已实现")
	NotImplemented = NewError(0x011b).WithName("NotImplemented").WithDesc("not-implemented-未实现")
	Registered     = NewError(0x011c).WithName("Registered").WithDesc("registered-已注册")
	Unregistered   = NewError(0x011d).WithName("Unregistered").WithDesc("unregistered-未注册")
	Marshal        = NewError(0x011e).WithName("Marshal").WithDesc("marshal-序列化")
	Unmarshal      = NewError(0x011f).WithName("Unmarshal").WithDesc("unmarshal-反序列化")

	// 协程
	GoroutinePanic = NewError(0x0200).WithName("GoroutinePanic").WithDesc("goroutine-panic-协程-panic")
	GoroutineDone  = NewError(0x0201).WithName("GoroutineDone").WithDesc("goroutine-done-协程-结束")
	// 函数
	FunctionPanic = NewError(0x0300).WithName("FunctionPanic").WithDesc("function-panic-函数-panic")
	FunctionDone  = NewError(0x0301).WithName("FunctionDone").WithDesc("function-done-函数-结束")
	// channel
	ChannelFull      = NewError(0x0400).WithName("ChannelFull").WithDesc("channel-full-通道-满")
	ChannelEmpty     = newError(0x0401).WithName("ChannelEmpty").WithDesc("channel-empty-通道-空")
	ChannelNotClosed = NewError(0x0402).WithName("ChannelNotClosed").WithDesc("channel-not-closed-通道-未关闭")
	ChannelClosed    = NewError(0x0403).WithName("ChannelClosed").WithDesc("channel-closed-通道-已关闭")

	Retry     = "retry"     // 重试
	Parameter = "parameter" // 参数
	Etcd      = "etcd"      // Etcd
	Redis     = "redis"     // Redis
	Mongodb   = "mongodb"   // Mongodb
	Mysql     = "mysql"     // Mysql
	Kafka     = "kafka"     // Kafka
	Nats      = "nats"      // Nats
	Nil       = "nil"       // 空
	// Link 链接
	Link = NewError(0xf001).WithName("Link").WithDesc("link error")
	// System 系统
	System = NewError(0xf002).WithName("System").WithDesc("system error")
	// Param 参数
	Param = NewError(0xf003).WithName("Param").WithDesc("parameter error")
	//// Packet 数据包
	//Packet = NewError(0xf004, "Packet", "packet error")

	// LogLevel 日志等级
	LogLevel = NewError(0xf00d).WithName("LogLevel").WithDesc("log level error")

	//// Insert 插入
	//Insert = NewError(0xf012, "Insert", "insert error")
	//// Find 查找
	//Find = NewError(0xf013, "Find", "find error")
	//// Update 更新
	//Update = NewError(0xf014, "Update", "update error")
	//// Delete 删除
	//Delete = NewError(0xf015, "Delete", "delete error")
	// Duplicate 重复
	Duplicate = NewError(0xf016).WithName("Duplicate").WithDesc("duplicate error")
	// Config 配置
	Config = NewError(0xf017).WithName("Config").WithDesc("config error")

	//// InvalidOperation 无效操作
	//InvalidOperation = NewError(0xf018, "InvalidOperation", "invalid operation")
	//// IllConditioned 条件不足
	//IllConditioned = NewError(0xf019, "IllConditioned", "ill conditioned")
	//// PermissionDenied 没有权限
	//PermissionDenied = NewError(0xf01a, "PermissionDenied", "permission denied")
	//// BlockedAccount 冻结账号
	//BlockedAccount = NewError(0xf01b, "BlockedAccount", "blocked account")
	//// Send 发送
	//Send = NewError(0xf01c, "Send", "send")
	//// Configure 给配置
	//Configure = NewError(0xf01d, "Configure", "configure")
	//// Retry 重试
	//Retry = NewError(0xf01e, "Retry", "retry")
	//// MessageIDNonExistent 消息ID 不存在
	//MessageIDNonExistent = NewError(0xf01f, "MessageIDNonExistent", "message id non-existent")
	//// Redis 系统 Redis
	//Redis = NewError(0xf020, "Redis", "redis")
	//// Busy 繁忙
	//Busy = NewError(0xf021, "Busy", "busy")
	//// OutOfResources 资源不足
	//OutOfResources = NewError(0xf022, "OutOfResources", "out of resources")
	//// NATS NATS错误
	//NATS = NewError(0xf023, "NATS", "nats")
	// PacketQuantityLimit 包数量限制
	PacketQuantityLimit = NewError(0xf024).WithName("PacketQuantityLimit").WithDesc("packet quantity limit")
	//// OverloadWarning 过载-告警
	//OverloadWarning = NewError(0xf025, "OverloadWarning", "overload warning")
	//// OverloadError 过载-错误
	//OverloadError = NewError(0xf026, "OverloadError", "overload error")
	// MessageIDDisable 消息ID禁用
	MessageIDDisable = NewError(0xf027).WithName("MessageIDDisable").WithDesc("message id is disabled")
	// MessageIDExistent 消息ID 存在
	MessageIDExistent = NewError(0xf028).WithName("MessageIDExistent").WithDesc("message id existent")
	//// ModeMismatch 模式 不匹配
	//ModeMismatch = NewError(0xf029, "ModeMismatch", "mode mismatch")
	//// FormatMismatch 格式 不匹配
	//FormatMismatch = NewError(0xf02a, "FormatMismatch", "format mismatch")
	//// MISSING 找不到,丢失,未命中
	//MISSING = NewError(0xf02b, "MISSING", "missing")
	//// VersionMismatch 版本 不匹配
	//VersionMismatch = NewError(0xf02c, "VersionMismatch", "version mismatch")

	// PacketHeaderLength 数据包头长度
	PacketHeaderLength = NewError(0xf02f).WithName("PacketHeaderLength").WithDesc("packet header length error")

	// ChannelNil 通道 未初始化
	ChannelNil = NewError(0xf032).WithName("ChannelNil").WithDesc("channel is nil")
	// MessageIDNonExistent 消息ID 不存在
	MessageIDNonExistent = NewError(0xf033).WithName("MessageIDNonExistent").WithDesc("message id non-existent")

	// Unknown 未知
	Unknown = NewError(0xffff).WithName("Unknown").WithDesc("unknown error")
	// 0xffff
)
