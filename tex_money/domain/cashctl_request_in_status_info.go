package domain

// 現金入出金機制御:入金ステータス取得要求
type RequestRequestInStatus struct {
	RequestInfo RequestInfo `json:"requestInfo"`
}

type ResultRequestInStatus struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Result      bool        `json:"result"`                //処理結果
	ErrorCode   string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string      `json:"errorDetail,omitempty"` //エラー詳細
}
