package interfaces

import "tex_money/domain"

// 入出金レポート印刷要求
type PrintReportRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
	SendResult(texCon *domain.TexContext, qresInfo domain.ResultSupply)
}
