package domain

type RequestGetService struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	IdDevice    string      `json:"idDevice,omitempty"`
	IdExtSys    string      `json:"idExtSys,omitempty"`
}

type ResultGetService struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
	StatusService bool        `json:"statusService"`
	IdDevice      string      `json:"idDevice,omitempty"`
	IdExtSys      string      `json:"idExtSys,omitempty"`
}

type StatusService struct {
	StatusService bool   `json:"statusService"`
	IdDevice      string `json:"idDevice,omitempty"`
	IdExtSys      string `json:"idExtSys,omitempty"`
}
