package domain

// 実行制御情報
type RequestControlService struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	StatusService bool        `json:"statusService"`
	IdDevice      string      `json:"idDevice"`
	IdExtSys      string      `json:"idExtSys"`
}

type ResultControlService struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Result      bool        `json:"result"`                //処理結果
	ErrorCode   string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string      `json:"errorDetail,omitempty"` //エラー詳細
}
