package domain

// 金庫情報(0:現金有高~9:売上金回収までの情報)
type SafeInfo struct {
	SalesCompleteAmount int             //売上金回収済
	SalesCompleteCount  int             //売上金回収回数
	CollectCount        int             //回収操作回数
	SortInfoTbl         [11]SortInfoTbl //分類情報
}

// 分類別金庫情報
type SortInfoTbl struct {
	SortType   int                         `json:"sortType"`   //分類情報種別
	Amount     int                         `json:"amount"`     //金額
	CountTbl   [CASH_TYPE_SHITEI]int       `json:"countTbl"`   //通常金種別枚数
	ExCountTbl [EXTRA_CASH_TYPE_SHITEI]int `json:"exCountTbl"` //拡張金種別枚数
}
