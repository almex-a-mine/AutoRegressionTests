package domain

type RequestReverseExchangeCalculation struct {
	RequestInfo     RequestInfo `json:"requestInfo"`
	ExchangeType    int         `json:"exchangeType"`
	OverflowCashbox bool        `json:"overflowCashbox"`
	Amount          int         `json:"amount"`
}

type ResultReverseExchangeCalculation struct {
	RequestInfo        RequestInfo                  `json:"requestInfo"`
	Result             bool                         `json:"result"`                //処理結果
	ErrorCode          string                       `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail        string                       `json:"errorDetail,omitempty"` //エラー詳細
	TargetAmount       int                          `json:"targetAmount"`
	TargetExCountTbl   *[EXTRA_CASH_TYPE_SHITEI]int `json:"targetExCountTbl,omitempty"`
	ExchangeExCountTbl [EXTRA_CASH_TYPE_SHITEI]int  `json:"exchangeExCountTbl"`
}
