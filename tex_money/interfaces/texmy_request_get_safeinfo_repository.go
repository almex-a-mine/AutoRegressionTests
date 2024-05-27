package interfaces

import "tex_money/domain"

type GetSafeInfoRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
}
