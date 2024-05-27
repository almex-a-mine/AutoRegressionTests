package domain

// 現金入出金機制御:入金開始要求情報
type RequestInStart struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	ModeOperation int         `json:"modeOperation"` //運用モード
	CountClear    bool        `json:"countClear"`    //入金枚数クリア
	TargetDevice  int         `json:"targetDevice"`  //対象デバイス
}

type ResultInStart struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
	CashControlId string      `json:"cashControlId"`         //入出金機制御管理番号
}

func NewRequestInStart(info RequestInfo, modeOperation int, countClear bool, targetDevice int) *RequestInStart {
	return &RequestInStart{
		RequestInfo:   info,
		ModeOperation: modeOperation,
		CountClear:    countClear,
		TargetDevice:  targetDevice,
	}
}
