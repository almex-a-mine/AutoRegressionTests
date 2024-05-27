package domain

// 現金入出金機制御:出金停止要求
type RequestOutStop struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	CashControlId string      `json:"cashControlId"`
}

type ResultOutStop struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Result      bool        `json:"result"`                //処理結果
	ErrorCode   string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string      `json:"errorDetail,omitempty"` //エラー詳細
}

func NewRequestOutStop(info RequestInfo, cashControlId string) RequestOutStop {
	return RequestOutStop{
		RequestInfo:   info,
		CashControlId: cashControlId,
	}
}
