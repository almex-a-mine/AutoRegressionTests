package interfaces

import "tex_money/domain"

// 追加補充要求
type MoneyAddReplenishRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
	SenSorSendFinish(texCon *domain.TexContext, reqType int)
	SendResult(texCon *domain.TexContext, reqInfo domain.ResultInStart) bool            // 処理結果応答:入金開始要求
	SendResultForInEnd(texCon *domain.TexContext, reqInfo domain.ResultInEnd) bool      // 処理結果応答:入金終了要求
	SendResultForInEndForDB(texCon *domain.TexContext, resinfo domain.ResultInEnd) bool // 処理結果応答:入金終了要求　in_end後検知
}
