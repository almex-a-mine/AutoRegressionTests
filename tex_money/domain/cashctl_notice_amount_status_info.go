package domain

// 現金入出金制御:有高ステータス通知
type AmountStatus struct {
	CoinStatusCode int                         `json:"coinStatusCode"`        //硬貨結果通知コード
	BillStatusCode int                         `json:"billStatusCode"`        //紙幣結果通知コード
	Amount         int                         `json:"amount"`                //金額
	CountTbl       [CASH_TYPE_SHITEI]int       `json:"countTbl"`              //通常金種別枚数
	ExCountTbl     [EXTRA_CASH_TYPE_SHITEI]int `json:"exCountTbl"`            //拡張金種別枚数
	DepositTbl     [CASH_TYPE_SHITEI]int       `json:"depositTbl"`            //入金可能枚数
	ErrorCode      string                      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail    string                      `json:"errorDetail,omitempty"` //エラー詳細
}
