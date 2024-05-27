package domain

// 精算機状態管理：状態変更要求（補充完了）
type RequestChangeSupply struct {
	RequestInfo RequestInfo  `json:"requestInfo"`
	SupplyType  int          `json:"supplyType"` //補充種別
	InfoTrade   InfoTrade    `json:"infoTrade"`  //取引情報
	InfoSafe    InfoSafeInfo `json:"infoSafe"`   //金庫情報
}

type ResultChangeSupply struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Result      bool        `json:"result"`                //処理結果
	ErrorCode   string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string      `json:"errorDetail,omitempty"` //エラー詳細
}

type InfoTrade struct {
	StatusTrade       bool              `json:"statusTrade"`       //取引状況
	TypeTrade         bool              `json:"typeTrade"`         //取引種別
	BillingAmount     int               `json:"billingAmount"`     //請求金額
	DepositAmount     int               `json:"depositAmount"`     //入金金額
	PaymentPlanAmount int               `json:"paymentPlanAmount"` //出金予定金額
	PaymentAmount     int               `json:"paymentAmount"`     //出金金額
	PayoutBalance     int               `json:"payoutBlance"`      //払出残額
	PaymentType       int               `json:"paymentType"`       //決済方法
	CashInfoTbl       []CashInfoTblInfo `json:"cashInfoTbl"`       //入出金情報
}

type CashInfoTblInfo struct {
	InfoType   int                         `json:"infoType"`   //入出金種別
	Amount     int                         `json:"amount"`     //金額
	CountTbl   [CASH_TYPE_SHITEI]int       `json:"countTbl"`   //通常金種別枚数
	ExCountTbl [EXTRA_CASH_TYPE_SHITEI]int `json:"exCountTbl"` //拡張金種別枚数
}

type InfoSafeInfo struct {
	CurrentStatusTbl [CASH_TYPE_SHITEI]int `json:"currentStatusTbl"` //通常金種別状況
	SortInfoTbl      []SortInfoTbl         `json:"sortInfotbl"`      //分類情報
}

func NewRequestChangeSupply(info RequestInfo, supplyType int, infoTrade InfoTrade, infoSafe InfoSafeInfo) *RequestChangeSupply {
	return &RequestChangeSupply{
		RequestInfo: info,
		SupplyType:  supplyType,
		InfoTrade:   infoTrade,
		InfoSafe:    infoSafe,
	}
}
