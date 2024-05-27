package domain

const SET_AMOUNT_CASHTYPE_NUMBER = 26 //設定金種数

// 現在枚数変更情報
type RequestSetAmount struct {
	RequestInfo   RequestInfo                 `json:"requestInfo"`
	CashControlId string                      `json:"cashControlId"`
	OperationMode int                         `json:"operationMode"` //操作モード
	CashTbl       [EXTRA_CASH_TYPE_SHITEI]int `json:"cashTbl"`       //指定枚数

}

type ResultSetAmount struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	CashControlId string      `json:"cashControlId"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
}
