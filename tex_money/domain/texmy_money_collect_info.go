package domain

// 回収情報(途中回収要求／全回収要求／売上金回収要求)
type RequestMoneyCollect struct {
	RequestInfo   RequestInfo       `json:"requestInfo"`
	CashControlId string            `json:"cashControlId"` //入出金制御管理番号
	CollectMode   int               `json:"collectMode"`   //回収モード
	OutType       int               `json:"outType"`       //払出方向
	StatusMode    int               `json:"statusMode"`    //動作モード
	SalesAmount   int               `json:"salesAmount"`   //売上金額
	CashTbl       [CASH_TYPE_UI]int `json:"cashTbl"`       //回収枚数
}

type ResultMoneyCollect struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
	CashControlId string      `json:"cashControlId"`         //入出金機制御管理番号
}
