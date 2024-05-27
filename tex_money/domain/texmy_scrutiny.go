package domain

type RequestScrutiny struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	CashControlId string      `json:"cashControlId"`
	TargetDevice  int         `json:"targetDevice"`
}

type ResultScrutiny struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
	CashControlId string      `json:"cashControlId"`
}

func NewResultScrutiny(info RequestInfo, result bool, errorCode string, errorDetail string, cashControlId string) *ResultScrutiny {
	return &ResultScrutiny{
		RequestInfo:   info,
		Result:        result,
		ErrorCode:     errorCode,
		ErrorDetail:   errorDetail,
		CashControlId: cashControlId,
	}
}
