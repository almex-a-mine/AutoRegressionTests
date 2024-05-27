package domain

type RequestRegisterMoneySetting struct {
	RequestInfo         RequestInfo          `json:"requestInfo"`
	ChangeReserveCount  *ChangeReserveCount  `json:"changeReserveCount"`  //釣銭準備金枚数
	ChangeShortageCount *ChangeShortageCount `json:"changeShortageCount"` //不足枚数
	ExcessChangeCount   *ExcessChangeCount   `json:"excessChangeCount"`   //あふれ枚数
}

type ResultRegisterMoneySetting struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Result      bool        `json:"result"`                //処理結果
	ErrorCode   string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string      `json:"errorDetail,omitempty"` //エラー詳細
}
