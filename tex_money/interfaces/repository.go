package interfaces

import "tex_money/domain"

type Repository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
}
