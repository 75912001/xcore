package gateway

import xnettcp "xcore/lib/net/tcp"

type User struct {
	id           uint64          // 用户id
	login        bool            // 是否登录
	remote       xnettcp.IRemote // 连接
	LoginService *LoginService
	//LogicService *LogicService
	timeoutValid bool // 超时是否有效
}

func newUser(remote xnettcp.IRemote) *User {
	return &User{
		remote:       remote,
		timeoutValid: true,
	}
}

// 退出
func (p *User) exit() {
	p.timeoutValid = false
	p.remote.Stop()
}

// 释放用户
func (p *User) release() {
}
