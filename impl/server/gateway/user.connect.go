package gateway

import xnetconnect "github.com/75912001/xcore/lib/net/common"

// UserConnect 用户-链接
type UserConnect struct {
	IRemote          xnetconnect.IRemote
	ExpireTimestamp  int64  // 到期时间戳-灰度
	heartbeatRandom  uint64 // 心跳随机数
	heartbeatTimeout int64  // 心跳超时时间
	Status           *xnetconnect.Status
}

func newUserConnect(remote xnetconnect.IRemote) *UserConnect {
	return &UserConnect{
		IRemote:          remote,
		heartbeatTimeout: gServer.TimeMgr.ShadowTimestamp() + UserHeartbeatIntervalMax,
		Status:           xnetconnect.NewStatus(),
	}
}
