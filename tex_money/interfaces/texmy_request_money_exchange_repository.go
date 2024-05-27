package interfaces

import "tex_money/domain"

// 両替要求
type MoneyExchangeRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
	SenSorIndataSendFinish(texCon *domain.TexContext, reqType int)
	SenSorOutdataSendFinish(texCon *domain.TexContext, reqType int)
	SendResult(texCon *domain.TexContext, recvinfo interface{}) //処理結果応答
}
