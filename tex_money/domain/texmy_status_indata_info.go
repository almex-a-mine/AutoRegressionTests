package domain

// 入金データ通知情報
type StatusIndata struct {
	CashControlId  string                      `json:"cashControlId"`          //入出金制御管理番号
	StatusAction   bool                        `json:"statusAction"`           //動作状況
	StatusResult   *bool                       `json:"statusResult,omitempty"` //入金結果
	Amount         int                         `json:"amount"`                 //金額
	CountTbl       [CASH_TYPE_SHITEI]int       `json:"countTbl"`               //通常金種別枚数
	ExCountTbl     [EXTRA_CASH_TYPE_SHITEI]int `json:"exCountTbl"`             //拡張金種別枚数
	ErrorCode      string                      `json:"errorCode,omitempty"`    //エラーコード
	ErrorDetail    string                      `json:"errorDetail,omitempty"`  //エラー詳細
	StatusActionEx int                         `json:"statusActionEx"`         //拡張動作状況
}

// アドレス取得用
var True, False = true, false
