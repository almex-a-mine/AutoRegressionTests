package domain

// 稼働データ管理:現在端末状態取得要求
type RequestGetTermInfoNow struct {
	RequestInfo RequestInfo `json:"requestInfo"`
}

type ResultGetTermInfoNow struct {
	RequestInfo  RequestInfo  `json:"requestInfo"`
	Result       bool         `json:"result"`                //処理結果
	ErrorCode    string       `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail  string       `json:"errorDetail,omitempty"` //エラー詳細
	ReportDate   int          `json:"reportDate"`            //記録日付
	ReportTime   int          `json:"reportTime"`            //記録時刻
	InfoTermTbl  InfoTermTbl  `json:"infoTerm"`              //端末情報
	InfoTradeTbl InfoTradeTbl `json:"infoTrade"`             //取引情報

	InfoSalesTbl          InfoSalesTbl          `json:"infoSales"` //売上情報
	InfoSafeTblGetTermNow InfoSafeTblGetTermNow `json:"infoSafe"`  //金庫情報
}

// 端末情報
type InfoTermTbl struct {
	StatusError    int `json:"statusError"`    //エラー状態
	TermErrorCode  int `json:"termErrorCode"`  //エラーコード
	TermErrorState int `json:"termErrorState"` //エラー発生状態
	StatusHandling int `json:"statusHandling"` //取扱状態
	StatusSecurity int `json:"statusSecurity"` //セキュリティ状態
	StatusDoor     int `json:"statusDoor"`     //扉状態
	StatusKeySw    int `json:"statuskeySw"`    //キーSW状態
	StatusCall     int `json:"statusCall"`     //従業員呼出状態
	PaymentMode    int `json:"paymentMode"`    //決済方法モード
}

// 取引情報
type InfoTradeTbl struct {
	StatusTrade       bool                   `json:"statusTrade"`       //取引状況
	TypeTrade         bool                   `json:"typeTrade"`         //取引種別
	BillingAmount     int                    `json:"billingAmount"`     //請求金額
	DepositAmount     int                    `json:"depositAmount"`     //入金金額
	PaymentPlanAmount int                    `json:"paymentPlanAmount"` //出金予定金額
	PaymentAmount     int                    `json:"paymentAmount"`     //出金金額
	PayoutBalance     int                    `json:"payoutBlance"`      //払出残額
	PaymentType       int                    `json:"paymentType"`       //決済方法
	CashInfoTbl       []CashInfoTblGetTemNow `json:"cashInfoTbl"`       //入出金情報
}

// 入出金情報
type CashInfoTblGetTemNow struct {
	InfoType   int                         `json:"infoType"`   //入出金種別
	Amount     int                         `json:"amount"`     //金額
	CountTbl   [CASH_TYPE_SHITEI]int       `json:"countTbl"`   //通常金種別枚数
	ExCountTbl [EXTRA_CASH_TYPE_SHITEI]int `json:"exCountTbl"` //拡張金種別枚数
}

// 売上情報
type InfoSalesTbl struct {
	SalesAmount   int            `json:"salesAmount"`   //売上金額合計
	ExchangeTotal int            `json:"exchangeTotal"` //両替金額合計
	SalesTypeTbl  []SalesTypeTbl `json:"salesTypeTbl"`  //売上種別情報
}

// 売上種別情報
type SalesTypeTbl struct {
	SalesType   int `json:"salesType"`   //売上種別
	PaymentType int `json:"paymentType"` //決済方法
	Amount      int `json:"amount"`      //金額
	Count       int `json:"count"`       //回数
}

// 金庫情報
type InfoSafeTblGetTermNow struct {
	CurrentStatusTbl [CASH_TYPE_SHITEI]int `json:"currentStatusTbl"` //通常金種別状況
	SortInfoTbl      []SortInfoTbl         `json:"sortInfotbl"`      //分類情報
}

func NewRequestGetTermInfoNow(info RequestInfo) RequestGetTermInfoNow {
	return RequestGetTermInfoNow{
		RequestInfo: info,
	}
}
