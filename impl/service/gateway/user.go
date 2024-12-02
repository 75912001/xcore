package gateway

import xnettcp "xcore/lib/net/tcp"

// UserConnect 用户-链接
type UserConnect struct {
	IRemote          xnettcp.IRemote
	ExpireTimestamp  int64  // 到期时间戳-灰度
	heartbeatRandom  uint64 // 心跳随机数
	heartbeatTimeout int64  // 心跳超时时间
}

func newUserConnect(remote xnettcp.IRemote) *UserConnect {
	return &UserConnect{
		IRemote:          remote,
		heartbeatTimeout: gService.TimeMgr.ShadowTimestamp() + UserHeartbeatIntervalMax,
	}
}

type User struct {
	id           uint64       // 用户id
	login        bool         // 是否登录
	connect      *UserConnect // 连接
	LoginService *LoginService
	//LogicService *LogicService
	timeoutValid bool // 超时是否有效
}

func newUser(remote xnettcp.IRemote) *User {
	return &User{
		connect:      newUserConnect(remote),
		timeoutValid: true,
	}
}

// 退出
func (p *User) exit() {
	p.timeoutValid = false
	p.connect.IRemote.Stop()
}

// 释放用户
func (p *User) release() {
}
