package domain

// 入出金データクリア要求
type RequestClearCashInfo struct {
	RequestInfo RequestInfo `json:"requestInfo"`
}

type ResultClearCashInfo struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Result      bool        `json:"result"`                //処理結果
	ErrorCode   string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string      `json:"errorDetail,omitempty"` //エラー詳細
}
