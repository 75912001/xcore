package gateway

import (
	xnettcp "github.com/75912001/xcore/lib/net/common"
)

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
