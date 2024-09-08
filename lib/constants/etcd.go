package constants

const EtcdTtlSecondDefault int64 = 33 // etcd 默认 TTL 时间 秒

const EtcdWatchMsgTypeService string = "service"           // etcd watch 消息类型-服务
const EtcdWatchMsgTypeCommand string = "command"           // etcd watch 消息类型-命令
const EtcdWatchMsgTypeGM string = "gm"                     // etcd watch 消息类型-GM
const EtcdWatchMsgTypeServiceBench string = "serviceBench" // etcd watch 消息类型-服务配置
