package gateway

import xnettcp "github.com/75912001/xcore/lib/net/tcp"
import xnetconnect "github.com/75912001/xcore/lib/net/connect"

// UserConnect 用户-链接
type UserConnect struct {
	IRemote          xnettcp.IRemote
	ExpireTimestamp  int64  // 到期时间戳-灰度
	heartbeatRandom  uint64 // 心跳随机数
	heartbeatTimeout int64  // 心跳超时时间
	Status           *xnetconnect.Status
}

func newUserConnect(remote xnettcp.IRemote) *UserConnect {
	return &UserConnect{
		IRemote:          remote,
		heartbeatTimeout: gService.TimeMgr.ShadowTimestamp() + UserHeartbeatIntervalMax,
		Status:           xnetconnect.NewStatus(),
	}
}
