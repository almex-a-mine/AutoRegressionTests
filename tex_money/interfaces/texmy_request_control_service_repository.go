package interfaces

import "tex_money/domain"

// 実行制御要求
type controlServiceRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
	NoticeControlService(pNoticeInfo *domain.StatusService)
}
