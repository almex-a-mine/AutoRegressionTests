package interfaces

import "tex_money/domain"

// 現金入出金機制御ステータス要求
type StatusCashRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
}
