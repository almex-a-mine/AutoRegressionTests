package interfaces

import "tex_money/domain"

//有高枚数要求
type AmountCashRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
	SendResult(texCon *domain.TexContext, pResInfo *domain.ResultGetTermInfoNow)
}
