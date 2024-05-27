package usecases

import "stub/domain"

type IWait interface {
	MakeWaitInfo(texCon *domain.TexContext, processId, requestid string, requestdata interface{}) WaitInfo
	GetWaitInfo(texCon *domain.TexContext, processId, requestid string) (WaitInfo, bool)
	SetWaitInfo(texCon *domain.TexContext, processId, requestid string, result interface{}) bool
	DelWaitInfo(texCon *domain.TexContext, processId, requestid string)
	WaitResultInfo(texCon *domain.TexContext, wInfo WaitInfo, timerValue int) bool
}
