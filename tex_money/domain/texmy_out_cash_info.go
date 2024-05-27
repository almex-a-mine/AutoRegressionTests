package domain

// 取引出金
type RequestOutCash struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	CashControlId string      `json:"cashControlId"` //入出金制御管理番号
	StatusMode    int         `json:"statusMode"`    //動作モード
	OutData       int         `json:"outData"`       //出金金額
}

type ResultOutCash struct {
	RequestInfo         RequestInfo           `json:"requestInfo"`
	Result              bool                  `json:"result"`                //処理結果
	ErrorCode           string                `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail         string                `json:"errorDetail,omitempty"` //エラー詳細
	CashControlId       string                `json:"cashControlId"`         //入出金機制御管理番号
	PaymentPlanCountTbl [CASH_TYPE_SHITEI]int `json:"paymentPlanCountTbl"`   // 出金予定枚数
	StatusMode          int                   `json:"statusMode"`            //動作モード
}
