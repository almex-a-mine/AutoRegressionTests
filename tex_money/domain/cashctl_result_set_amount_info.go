package domain

// 現金入出金機制御:有高枚数変更要求
type RequestCashctlSetAmount struct {
	RequestInfo   RequestInfo                 `json:"requestInfo"`
	OperationMode int                         `json:"operationMode"` //操作モード
	Amount        int                         `json:"amount"`        //金額
	CountTbl      [CASH_TYPE_SHITEI]int       `json:"countTbl"`      //通常金種別枚数
	ExCountTbl    [EXTRA_CASH_TYPE_SHITEI]int `json:"exCountTbl"`    //拡張金種別枚数
}

type ResultCashctlSetAmount struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	CashControlId string      `json:"cashControlId"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
}
