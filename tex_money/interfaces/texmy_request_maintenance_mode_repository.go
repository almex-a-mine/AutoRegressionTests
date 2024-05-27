package interfaces

import "tex_money/domain"

type MaintenanceModeRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
}
