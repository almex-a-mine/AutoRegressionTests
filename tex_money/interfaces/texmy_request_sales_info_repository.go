package interfaces

import "tex_money/domain"

// 売上金情報要求
type SalesInfoRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
	SenSorSendFinish(texCon *domain.TexContext)
}
