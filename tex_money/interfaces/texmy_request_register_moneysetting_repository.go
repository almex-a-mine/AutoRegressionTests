package interfaces

import "tex_money/domain"

type RegisterMoneySettingRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
}
