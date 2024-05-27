package domain

// システム動作モード変更要求情報
type RequestChangeSystemOperation struct {
	RequestInfo           RequestInfo `json:"requestInfo"`
	StatusSystemOperation int         `json:"statusSystemOperation"` //システム動作モード
	IdDevice              string      `json:"idDevice"`              //デバイス識別番号
	IdExtSys              string      `json:"idExtSys"`              //外部システム識別番号
}

// システム動作モード変更要求応答情報
type ResultChangeSystemOperation struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Result      bool        `json:"result"`
	ErrorCode   string      `json:"errorCode,omitempty"`
	ErrorDetail string      `json:"errorDetail,omitempty"`
}
