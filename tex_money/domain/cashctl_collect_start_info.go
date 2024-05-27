package domain

// 現金入出金機制御:回収開始要求
type RequestCollectStart struct {
	RequestInfo RequestInfo       `json:"requestInfo"`
	CollectMode int               `json:"collectMode"` //回収種別
	CountTbl    [CASH_TYPE_UI]int `json:"countTbl"`    //枚数情報
	Amount      int               `json:"amount"`      //金額情報
}

type ResultCollectStart struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
	CashControlId string      `json:"cashControlId"`         //入出金機制御管理番号
}

func NewRequestCollectStart(info RequestInfo, collectMode int, amount int, countTbl [CASH_TYPE_UI]int) RequestCollectStart {
	return RequestCollectStart{
		RequestInfo: info,
		CollectMode: collectMode,
		Amount:      amount,
		CountTbl:    countTbl,
	}
}
