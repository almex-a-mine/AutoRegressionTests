package domain

// 取引入金情報
type RequestPayCash struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	CashControlId string      `json:"cashControlId"` //入出金機制御管理番号
	ModeOperation int         `json:"modeOperation"` //運用モード
	CountClear    bool        `json:"countClear"`    //入金枚数クリア
	TargetDevice  int         `json:"targetDevice"`  //対象デバイス
	StatusMode    int         `json:"statusMode"`    //動作モード
}

type ResultPayCash struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
	CashControlId string      `json:"cashControlId"`         //入出金機制御管理番号
}
