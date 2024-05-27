package domain

// 有高枚数要求
type RequestAmountCash struct {
	RequestInfo RequestInfo `json:"requestInfo"`
}

type ResultAmountCash struct {
	RequestInfo RequestInfo                 `json:"requestInfo"`
	Result      bool                        `json:"result"`                //処理結果
	ErrorCode   string                      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string                      `json:"errorDetail,omitempty"` //エラー詳細
	Amount      int                         `json:"amount"`                //金額
	CountTbl    [CASH_TYPE_SHITEI]int       `json:"countTbl"`              //通常金種別枚数
	ExCountTbl  [EXTRA_CASH_TYPE_SHITEI]int `json:"exCountTbl"`            //拡張金種別枚数
}
