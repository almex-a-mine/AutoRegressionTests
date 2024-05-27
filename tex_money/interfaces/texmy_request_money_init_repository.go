package interfaces

import "tex_money/domain"

// 初期補充要求
type MoneyInitRepository interface {
	Start()
	Stop()
	SendResult(texCon *domain.TexContext, cashControlId string) bool //処理結果応答
	ControlService(reqInfo domain.RequestControlService)
	SenSorSendFinish(texCon *domain.TexContext, reqType int)
}
