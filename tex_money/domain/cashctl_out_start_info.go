package domain

// 現金入出金機制御:出金開始要求
type RequestOutStart struct {
	RequestInfo        RequestInfo `json:"requestInfo"`
	StatusOutRejectBox bool        `json:"statusOutRejectBox"` //出金種別
	OutMode            int         `json:"outMode"`            //出金種別
	OutStatusCashInfoTbl
}

type ResultOutStart struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
	CashControlId string      `json:"cashControlId"`         //入出金機制御管理番号
}

type OutStatusCashInfoTbl struct {
	Amount   int               `json:"amount"`   //金種
	CountTbl [CASH_TYPE_UI]int `json:"countTbl"` //金種別枚数
}

func NewRequestOutStart(info RequestInfo, statusOutRejectBox bool, outMode int, amount int, countTbl [CASH_TYPE_UI]int) RequestOutStart {
	return RequestOutStart{
		RequestInfo:        info,
		StatusOutRejectBox: statusOutRejectBox,
		OutMode:            outMode,
		OutStatusCashInfoTbl: OutStatusCashInfoTbl{
			Amount:   amount,
			CountTbl: countTbl,
		},
	}
}
