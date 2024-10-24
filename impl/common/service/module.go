package service

import "context"

type IService interface {
	Start(ctx context.Context) (err error) // 启动服务
	PreStop() error                        // 服务关闭前的处理
	Stop() (err error)                     // 停止服务
	// 获取消息范围
	//GetMsgRange() (uint32, uint32)
}
