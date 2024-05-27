package usecases

import "tex_money/domain"

// 集計データ管理
type AggregateManager interface {
	UpdateAggregateCountTbl(texCon *domain.TexContext, mode int, tbl int, countTbl [domain.EXTRA_CASH_TYPE_SHITEI]int) (result bool)
	UpdateBeforeCountTbl(texCon *domain.TexContext, mode int, tbl int) (result bool)
	ClearAggregateData(texCon *domain.TexContext)
	ClearAggregateSafeInfo(texCon *domain.TexContext, mode int) (result bool)
	GetAggregateSafeInfo(texCon *domain.TexContext, mode int) (result bool, aggregateSafeInfo domain.AggregateSafeInfo)
	DiffTbl(texCon *domain.TexContext, beforeTbl [domain.EXTRA_CASH_TYPE_SHITEI]int, afterTbl [domain.EXTRA_CASH_TYPE_SHITEI]int) (countTbl [domain.EXTRA_CASH_TYPE_SHITEI]int)
	GetAggregateSafeInfoAll() domain.AggregateData
	OutputLogAggregateExCountTbl()
}
