package service

type IService interface {
	Start() (err error) // 启动服务
	Stop() (err error)  // 停止服务

	PreShutdown() // 服务关闭前的处理
}
