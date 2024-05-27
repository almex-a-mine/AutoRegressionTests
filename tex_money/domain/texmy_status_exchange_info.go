package domain

// 両替ステータス情報
type StatusExchange struct {
	CashControlId    string                      `json:"cashControlId"`          //入出金制御管理番号
	InOutType        int                         `json:"inOutType"`              //入出金区分
	StatusAction     bool                        `json:"statusAction"`           //動作状況
	StatusResult     *bool                       `json:"statusResult,omitempty"` //出金結果
	Amount           int                         `json:"amount"`                 //金額
	CountTbl         [CASH_TYPE_SHITEI]int       `json:"countTbl"`               //通常金種別枚数
	ExCountTbl       [EXTRA_CASH_TYPE_SHITEI]int `json:"exCountTbl"`             //拡張金種別枚数
	ExchangeCountTbl [CASH_TYPE_SHITEI]int       `json:"exchangeCountTbl"`       //両替金種別予定枚数
	ErrorCode        string                      `json:"errorCode,omitempty"`    //エラーコード
	ErrorDetail      string                      `json:"errorDetail,omitempty"`  //エラー詳細
}
