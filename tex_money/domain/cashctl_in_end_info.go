package domain

// 現金入出金機制御:入金終了要求情報
type RequestInEnd struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	CashControlId string      `json:"cashControlId"` //入出金機制御管理番号
	TargetDevice  int         `json:"targetDevice"`  //対象デバイス
	StatusMode    int         `json:"statusMode"`    //動作モード
}

type ResultInEnd struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
	CashControlId string      `json:"cashControlId"`         //入出金機制御管理番号
}

func NewRequestInEnd(info RequestInfo, cashControlId string, targetDevice int, statusMode int) *RequestInEnd {
	return &RequestInEnd{
		RequestInfo:   info,
		CashControlId: cashControlId,
		TargetDevice:  targetDevice,
		StatusMode:    statusMode,
	}
}
