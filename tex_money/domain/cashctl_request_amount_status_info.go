package domain

// 現金入出金機制御:有高ステータス取得要求
type RequestAmountStatus struct {
	RequestInfo RequestInfo `json:"requestInfo"`
}

type ResultAmountStatus struct {
	RequestInfo    RequestInfo                 `json:"requestInfo"`
	Result         bool                        `json:"result"`                //処理結果
	CoinStatusCode int                         `json:"coinStatusCode"`        //硬貨結果通知コード
	BillStatusCode int                         `json:"billStatusCode"`        //紙幣結果通知コード
	Amount         int                         `json:"amount"`                //金額
	CountTbl       [CASH_TYPE_SHITEI]int       `json:"countTbl"`              //通常金種別枚数
	ExCountTbl     [EXTRA_CASH_TYPE_SHITEI]int `json:"exCountTbl"`            //拡張金種別枚数
	DepositTbl     [CASH_TYPE_SHITEI]int       `json:"depositTbl"`            //入金可能枚数
	ErrorCode      string                      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail    string                      `json:"errorDetail,omitempty"` //エラー詳細
}

func NewRequestAmountStatus(info RequestInfo) *RequestAmountStatus {
	return &RequestAmountStatus{
		RequestInfo: info,
	}
}
