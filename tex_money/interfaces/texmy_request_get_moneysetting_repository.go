package interfaces

import "tex_money/domain"

type GetMoneySettingRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
}
