package domain

// 有高データ情報
type StatusAmount struct {
	Amount      int                         `json:"cashType"`              //金額
	CountTbl    [CASH_TYPE_SHITEI]int       `json:"countTbl"`              //通常金種別枚数
	ExCountTbl  [EXTRA_CASH_TYPE_SHITEI]int `json:"exCountTbl"`            //拡張金種別枚数
	ErrorCode   string                      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string                      `json:"errorDetail,omitempty"` //エラー詳細
}
