package domain

// レポート用金庫情報
type AggregateSafeInfo struct {
	BeforeAmountCountTbl    [EXTRA_CASH_TYPE_SHITEI]int //処理前有高金種配列
	BeforeReplenishCountTbl [EXTRA_CASH_TYPE_SHITEI]int //処理前補充入金金種配列
	ReplenishCountTbl       [EXTRA_CASH_TYPE_SHITEI]int //補充入金金種配列
	BeforeCollectCountTbl   [EXTRA_CASH_TYPE_SHITEI]int //処理前回収金種配列
	CollectCountTbl         [EXTRA_CASH_TYPE_SHITEI]int //回収金種配列
	AfterAmountCountTbl     [EXTRA_CASH_TYPE_SHITEI]int //処理後有高金種配列
	SalesCollectCountTbl    [EXTRA_CASH_TYPE_SHITEI]int //売上回収金種配列
}

type AggregateData map[int]*AggregateSafeInfo
