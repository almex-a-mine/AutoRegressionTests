package interfaces

import "tex_money/domain"

type RequestChangeSystemOperationRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
}
