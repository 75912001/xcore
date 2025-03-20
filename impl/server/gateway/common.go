package gateway

import xtime "github.com/75912001/xcore/lib/time"

const UserInactiveTimeout int64 = xtime.OneDaySecond // 用户非活跃时间-秒
const UserLoginTimeout int64 = xtime.OneDaySecond    // 用户登录超时时间-秒

const UserHeartbeatInterval int64 = 20 // 用户心跳时间间隔-秒

const UserHeartbeatIntervalMax int64 = UserHeartbeatInterval * 2 // 用户心跳时间间隔-最大-秒
