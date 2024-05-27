package domain

type RequestMaintenanceMode struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Mode        int         `json:"mode"`   //保守業務モード
	Action      bool        `json:"action"` //動作要求
}

type ResultMaintenanceMode struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Result      bool        `json:"result"`                //処理結果
	ErrorCode   string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string      `json:"errorDetail,omitempty"` //エラー詳細
}

func NewResultMaintenanceMode(info RequestInfo, result bool, errorCode string, errorDetail string) *ResultMaintenanceMode {
	return &ResultMaintenanceMode{
		RequestInfo: info,
		Result:      result,
		ErrorCode:   errorCode,
		ErrorDetail: errorDetail,
	}
}
