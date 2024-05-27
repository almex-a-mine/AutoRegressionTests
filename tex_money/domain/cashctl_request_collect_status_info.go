package domain

// 現金入出金機制御:回収ステータス取得要求
type RequestCollectStatus struct {
	RequestInfo RequestInfo `json:"requestInfo"`
}

type ResultCollectStatus struct {
	RequestInfo      RequestInfo         `json:"requestInfo"`
	Result           bool                `json:"result"`                //処理結果
	ErrorCode        string              `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail      string              `json:"errorDetail,omitempty"` //エラー詳細
	CoinStatusAction bool                `json:"coinStatusAction"`      //硬貨動作状況
	CoinStatusCode   int                 `json:"coinStatusCode"`        //硬貨結果通知コード
	BillStatusAction bool                `json:"billStatusAction"`      //紙幣動作状況
	BillStatusCode   int                 `json:"billStatusCode"`        //紙幣結果通知コード
	CollectCountKin  int                 `json:"collectCountKin"`       //回収金額
	CashTbl          []CollectStatusCash `json:"cashTbl"`               //回収枚数
}

type CollectStatusCash struct {
	CashType  string `json:"cashType"`  //金種
	CashCount int    `json:"cashCount"` //枚数
}
