package interfaces

import "tex_money/domain"

// 回収要求(途中回収要求／全回収要求／売上金回収要求)
type MoneyCollectRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
	SenSorSendFinish(texCon *domain.TexContext, reqType int)
	SendResult(texCon *domain.TexContext, reqInfo domain.ResultOutStart) bool               //処理結果応答:出金開始要求
	SendResultCollectStart(texCon *domain.TexContext, reqInfo domain.ResultOutStop) bool    // 処理結果応答:出金停止要求
	SendResultOutStop(texCon *domain.TexContext, reqInfo domain.ResultCollectStart) bool    //処理結果応答:回収開始要求
	SendResultCollectStop(texCon *domain.TexContext, reqInfo domain.ResultCollectStop) bool //処理結果応答:回収停止要求
}
