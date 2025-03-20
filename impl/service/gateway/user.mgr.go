package gateway

import xnettcp "xcore/lib/net/tcp"

var gUserMgr *userMgr

func init() {
	gUserMgr = newUserMgr()
}

type userMgr struct {
	remoteMap map[xnettcp.IRemote]*User // key: remote
	idMap     map[uint64]*User
}

func newUserMgr() *userMgr {
	return &userMgr{
		remoteMap: make(map[xnettcp.IRemote]*User),
		idMap:     make(map[uint64]*User),
	}
}

// 增加一个用户,通过remote
func (p *userMgr) add(user *User, remote xnettcp.IRemote) {
	p.remoteMap[remote] = user
}

// 移除一个用户,通过remote
func (p *userMgr) remove(remote xnettcp.IRemote) {
	user, ok := p.remoteMap[remote]
	if !ok {
		return
	}
	user.release()
	delete(p.idMap, user.id)
	delete(p.remoteMap, remote)
}

// 移除一个用户的remote数据,设置用户的数据为非活跃(通过remote)
func (p *userMgr) removeRemoteAndSetInactive(remote xnettcp.IRemote) {
	user, ok := p.remoteMap[remote]
	if !ok {
		return
	}
	user.connect.Status.SetInactive(UserInactiveTimeout,
		xcontrol.NewCallBack(
			func(args ...interface{}) error {
				user := args[0].(*User)
				user.release()
				return nil
			},
			user,
		))
	user.connect.IRemote = nil
}

// 通过remote获取用户
func (p *userMgr) get(remote xnettcp.IRemote) *User {
	user, ok := p.remoteMap[remote]
	if !ok {
		return nil
	}
	return user
}
