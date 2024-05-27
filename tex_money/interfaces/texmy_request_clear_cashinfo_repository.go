package interfaces

import "tex_money/domain"

//入出金データクリア要求
type RequestClearCashInfoRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
	SendResult(texCon *domain.TexContext, presInfo *domain.ResultReportSafeInfo)
}
