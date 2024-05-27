package interfaces

import "tex_money/domain"

// 実行状態取得要求
type getServiceRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
}
