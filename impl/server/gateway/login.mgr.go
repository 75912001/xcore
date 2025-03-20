package gateway

type LoginServiceMgr struct {
	LoginService map[uint32]*LoginService // key: loginId
}

func NewLoginServiceMgr() *LoginServiceMgr {
	return &LoginServiceMgr{
		LoginService: make(map[uint32]*LoginService),
	}
}

func (p *LoginServiceMgr) Add(loginId uint32, loginService *LoginService) {
	p.LoginService[loginId] = loginService
}

func (p *LoginServiceMgr) Get(loginId uint32) *LoginService {
	return p.LoginService[loginId]
}

func (p *LoginServiceMgr) Remove(loginId uint32) {
	delete(p.LoginService, loginId)
}

// GetLoginService 按照策略获取一个LoginService
func (p *LoginServiceMgr) GetLoginService() *LoginService {
	// todo menglc 暂时只返回第一个
	for _, loginService := range p.LoginService {
		return loginService
	}
	return nil
}
