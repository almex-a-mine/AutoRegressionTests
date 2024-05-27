package domain

// 現金入出金機制御:出金ステータス取得要求
type RequestOutStatus struct {
	RequestInfo RequestInfo `json:"requestInfo"`
}

type ResultOutStatus struct {
	RequestInfo      RequestInfo     `json:"requestInfo"`
	Result           bool            `json:"result"`                //処理結果
	ErrorCode        string          `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail      string          `json:"errorDetail,omitempty"` //エラー詳細
	CoinStatusAction bool            `json:"coinStatusAction"`      //硬貨動作状況
	CoinStatusCode   int             `json:"coinStatusCode"`        //硬貨結果通知コード
	BillStatusAction bool            `json:"billStatusAction"`      //紙幣動作状況
	BillStatusCode   int             `json:"billStatusCode"`        //紙幣結果通知コード
	OutCountKin      int             `json:"outCountKin"`           //出金金額
	CashTbl          []OutStatusCash `json:"cashTbl"`               //出金枚数
}

type OutStatusCash struct {
	CashType  string `json:"cashType"`  //金種
	CashCount int    `json:"cashCount"` //枚数
}
