package interfaces

import "tex_money/domain"

// 取引出金要求
type OutCashRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
	SenSorSendFinish(texCon *domain.TexContext, reqType int)
	SendResult(texCon *domain.TexContext, reqInfo domain.ResultOutStart) bool
	SendResultForOutStop(texCon *domain.TexContext, reqInfo domain.ResultOutStop) bool // 処理結果応答:出金停止要求
	SetOutCashRefund(texCon *domain.TexContext, outStatus domain.OutStatus)
	CheckStatusMode(texCon *domain.TexContext)
}
