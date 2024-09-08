package error

// [系统错误码] lib/system 级别的错误码 [0-0xffff]
var (
	// Success 成功
	Success = NewError(0x0000).WithName("Success").WithDesc("success")
	// Link 链接
	Link = NewError(0xf001).WithName("Link").WithDesc("link error")
	//// System 系统
	//System = NewError(0xf002, "System", "system error")
	// Param 参数
	Param = NewError(0xf003).WithName("Param").WithDesc("parameter error")
	//// Packet 数据包
	//Packet = NewError(0xf004, "Packet", "packet error")
	//// Timeout 超时
	//Timeout = NewError(0xf005, "Timeout", "time out")
	//// ChannelFull 通道 满
	//ChannelFull = NewError(0xf006, "ChannelFull", "channel full")
	//// ChannelEmpty 通道 空
	//ChannelEmpty = NewError(0xf007, "ChannelEmpty", "channel empty")
	//// OutOfRange 超出范围
	//OutOfRange = NewError(0xf008, "OutOfRange", "out of range")
	//// InvalidValue 无效数值
	//InvalidValue = NewError(0xf009, "InvalidValue", "invalid value")
	//// Conflict 冲突
	//Conflict = NewError(0xf00a, "Conflict", "conflict")
	//// TypeMismatch 类型不匹配
	//TypeMismatch = NewError(0xf00b, "TypeMismatch", "type mismatch")
	//// InvalidPointer 无效指针
	//InvalidPointer = NewError(0xf00c, "InvalidPointer", "invalid pointer")
	// LogLevel 日志等级
	LogLevel = NewError(0xf00d).WithName("LogLevel").WithDesc("log level error")
	//// NonExistent 不存在
	//NonExistent = NewError(0xf00e, "NonExistent", "non-existent")
	//// Exists 存在
	//Exists = NewError(0xf00f, "Exists", "exists")
	//// Marshal 序列化
	//Marshal = NewError(0xf010, "Marshal", "marshal")
	//// Unmarshal 反序列化
	//Unmarshal = NewError(0xf011, "Unmarshal", "unmarshal")
	//// Insert 插入
	//Insert = NewError(0xf012, "Insert", "insert error")
	//// Find 查找
	//Find = NewError(0xf013, "Find", "find error")
	//// Update 更新
	//Update = NewError(0xf014, "Update", "update error")
	//// Delete 删除
	//Delete = NewError(0xf015, "Delete", "delete error")
	//// Duplicate 重复
	//Duplicate = NewError(0xf016, "Duplicate", "duplicate error")
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
	//// PacketQuantityLimit 包数量限制
	//PacketQuantityLimit = NewError(0xf024, "PacketQuantityLimit", "packet quantity limit")
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
	//// Unavailable 不可用
	//Unavailable = NewError(0xf02d, "Unavailable", "unavailable")
	// NotImplemented 未实现
	NotImplemented = NewError(0xf02e).WithName("NotImplemented").WithDesc("not implemented")
	//// PacketHeaderLength 数据包头长度
	//PacketHeaderLength = NewError(0xf02f, "PacketHeaderLength", "packet header length error")
	//// ChannelClosed 通道 已关闭
	//ChannelClosed = NewError(0xf030, "ChannelClosed", "channel closed")
	//// Unregistered 未注册
	//Unregistered = NewError(0xf031, "Unregistered", "unregistered")
	// Unknown 未知
	Unknown = NewError(0xffff).WithName("Unknown").WithDesc("unknown error")
	// 0xffff
)
