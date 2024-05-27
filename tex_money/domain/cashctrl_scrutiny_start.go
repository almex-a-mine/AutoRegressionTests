package domain

type RequestScrutinyStart struct {
	RequestInfo  RequestInfo `json:"requestInfo"`
	TargetDevice int         `json:"targetDevice"`
}

type ResultScrutinyStart struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
	CashControlId string      `json:"cashControlId"`
}

func NewRequestScrutinyStart(info RequestInfo, targetDevice int) *RequestScrutinyStart {
	return &RequestScrutinyStart{
		RequestInfo:  info,
		TargetDevice: targetDevice,
	}
}
