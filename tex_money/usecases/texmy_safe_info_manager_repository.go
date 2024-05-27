package usecases

import "tex_money/domain"

type SafeInfoManager interface {
	UpdateSalesCompleteAmount(texCon *domain.TexContext, salesCompleteAmount int)
	CountUpSalesCompleteCount(texCon *domain.TexContext)
	GetCollectCount() int
	CountUpCollectCount(texCon *domain.TexContext)
	ClearSalesInfo(texCon *domain.TexContext)
	ClearCashInfo(texCon *domain.TexContext)
	UpdateSortInfo(texCon *domain.TexContext, sortInfoTbl domain.SortInfoTbl) (result bool)
	InitSafeInfo()
	GetSortInfo(texCon *domain.TexContext, sortType int) (result bool, sortInfoTbl domain.SortInfoTbl)
	GetSafeInfo(texCon *domain.TexContext) domain.SafeInfo
	GetSalesInfo() (int, int)
	UpdateBeforeReplenishmentBalance(texCon *domain.TexContext)
	GetBeforeReplenishmentBalance() domain.SortInfoTbl
	OutputLogSafeInfoExCountTbl(texCon *domain.TexContext)

	UpdateBalanceInfo(texCon *domain.TexContext)                                                                                           // 差引情報更新
	UpdateSortInfoCumulative(texCon *domain.TexContext, sortType int, amount int, countTbl [10]int, exCountTbl [26]int)                    // 分類情報更新（累計）
	UpdateSortInfoCumulativeNoUpdateLogicalCash(texCon *domain.TexContext, sortType int, amount int, countTbl [10]int, exCountTbl [26]int) // 分類情報更新（累計）And 有高更新無し

	// デバイス有高 論理有高調整
	GetDeviceCashAvailable(texCon *domain.TexContext) domain.SortInfoTbl
	UpdateDeviceCashAvailable(texCon *domain.TexContext, tbl domain.SortInfoTbl)
	GetLogicalCashAvailable(texCon *domain.TexContext) domain.SortInfoTbl
	UpdateInLogicalCashAvailable(texCon *domain.TexContext, tbl domain.SortInfoTbl)
	UpdateOutLogicalCashAvailable(texCon *domain.TexContext, tbl domain.SortInfoTbl)
	UpdateAllLogicalCashAvailable(texCon *domain.TexContext, tbl domain.SortInfoTbl)
	GetAvailableBalance(texCon *domain.TexContext) (bool, domain.SortInfoTbl)

	// 更新する金庫情報の分類情報種別日本語取得
	GetInfoSafeSortTypeName(i int) string
}
