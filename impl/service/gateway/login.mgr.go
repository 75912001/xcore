package gateway

type LoginServiceMgr struct {
	LoginService map[uint32]*LoginService // key: loginId
}
