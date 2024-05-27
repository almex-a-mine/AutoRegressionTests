package interfaces

import "tex_money/domain"

// 取引入金要求
type PayCashRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
	SenSorSendFinish(texCon *domain.TexContext, reqType int)
	SendResult(texCon *domain.TexContext, reqInfo domain.ResultInStart) bool
	SendResultResultInEnd(texCon *domain.TexContext, reqInfo domain.ResultInEnd) bool
}
