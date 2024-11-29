package gateway

import xtime "xcore/lib/time"

const UserLoginTimeOut int64 = xtime.OneDaySecond // 用户登录超时时间-秒

const UserHeartbeatInterval int64 = 20 // 用户心跳时间间隔-秒
