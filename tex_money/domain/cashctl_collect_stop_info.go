package domain

// 現金入出金機制御:回収停止要求
type RequestCollectStop struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	CashControlId string      `json:"cashControlId"`
}

type ResultCollectStop struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Result      bool        `json:"result"`                //処理結果
	ErrorCode   string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string      `json:"errorDetail,omitempty"` //エラー詳細
}

func NewRequestCollectStop(info RequestInfo, cashControlId string) RequestCollectStop {
	return RequestCollectStop{
		RequestInfo:   info,
		CashControlId: cashControlId,
	}
}
