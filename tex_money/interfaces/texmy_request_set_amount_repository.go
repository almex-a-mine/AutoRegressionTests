package interfaces

import "tex_money/domain"

// 現在枚数変更要求
type SetAmountRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
	SenSorSendFinish(texCon *domain.TexContext, reqType int)
	ModeAnalogCollect(texCon *domain.TexContext, reqInfo domain.RequestSetAmount)
	SendResult(texCon *domain.TexContext, res domain.ResultCashctlSetAmount)
	ConnetctCoincasseteControl(texCon *domain.TexContext, cashTbl [26]int) //硬貨カセット操作要求との土管 pram1:変更後有高格納
}
