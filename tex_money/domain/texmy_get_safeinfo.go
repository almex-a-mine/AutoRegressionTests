package domain

type RequestGetSafeInfo struct {
	RequestInfo RequestInfo `json:"requestInfo"`
}

type ResultGetSafeInfo struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
	SalesComplete int         `json:"salesComplete"`         //売上金回収済
	SalesCount    int         `json:"salesCount"`            //売上金回収回数
	CollectCount  int         `json:"collectCount"`          //回収操作回数
	InfoSafe      InfoSafe    `json:"infoSafe"`              //金庫情報
}

// 金庫情報
type InfoSafe struct {
	CurrentStatusTbl [CASH_TYPE_SHITEI]int `json:"currentStatusTbl"` //通常金種別状況
	SortInfoTbl      [14]SortInfoTbl       `json:"sortInfotbl"`      //分類情報
}
