package domain

// 印刷制御:補充レシート要求
type RequestSupply struct {
	RequestInfo RequestInfo  `json:"requestInfo"`
	Typename    string       `json:"typeName"`   //補充レシート種別名
	TermNo      int          `json:"termNo"`     //精算端末番号
	SupplyDate  int          `json:"supplyDate"` //補充処理日付
	SupplyTime  int          `json:"supplyTime"` //補充処理時刻
	NumInfoTbl  []int        `json:"numInfoTbl"` //汎用数値情報
	StrInfoTbl  []string     `json:"strInfoTbl"` //汎用文字列情報
	InfoTrade   infoTradeTbl `json:"infoTrade"`  //取引情報
	InfoSales   infoSalesTbl `json:"infoSales"`  //売上情報
	InfoSafe    infoSafeTbl  `json:"infoSafe"`   //金庫情報
}

type ResultSupply struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Result      bool        `json:"result"`                //処理結果
	ErrorCode   string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string      `json:"errorDetail,omitempty"` //エラー詳細
	PrintId     string      `json:"printId"`               //印刷管理ID
}

// 取引情報
type infoTradeTbl struct {
	StatusTrade       bool          `json:"statusTrade"`       //取引状況
	TypeTrade         bool          `json:"typeTrade"`         //取引種別
	BillingAmount     int           `json:"billingAmount"`     //請求金額
	DepositAmount     int           `json:"depositAmount"`     //入金金額
	PaymentPlanAmount int           `json:"paymentPlanAmount"` //出金予定金額
	PaymentAmount     int           `json:"paymentAmount"`     //出金金額
	PayoutBalance     int           `json:"payoutBlance"`      //払出残額
	PaymentType       int           `json:"paymentType"`       //決済方法
	CashInfoTbl       []CashInfoTbl `json:"cashInfoTbl"`       //入出金情報
}

// 入出金情報
type CashInfoTbl struct {
	InfoType   int     `json:"infoType"`   //入出金種別
	Amount     int     `json:"amount"`     //金額
	CountTbl   [10]int `json:"countTbl"`   //通常金種別枚数
	ExCountTbl [26]int `json:"exCountTbl"` //拡張金種別枚数
}

// 売上情報
type infoSalesTbl struct {
	SalesAmount   int            `json:"salesAmount"`   //売上金額合計
	ExchangeTotal int            `json:"exchangeTotal"` //両替金額合計
	SalesTypeTbl  []salesTypeTbl `json:"salesTypeTbl"`  //売上種別情報
}

// 売上種別情報
type salesTypeTbl struct {
	SalesType   int `json:"salesType"`   //売上種別
	PaymentType int `json:"paymentType"` //決済方法
	Amount      int `json:"amount"`      //金額
	Count       int `json:"count"`       //回数
}

// 金庫情報
type infoSafeTbl struct {
	CurrentStatusTbl int           `json:"currentStatusTbl"` //通常金種別状況
	SortInfoTbl      []SortInfoTbl `json:"sortInfotbl"`      //分類情報
}

var ReportName = map[int]string{
	0:  "summary_sales",        //補充レシート種別:定義値 精算機別日計表
	1:  "report_cashcount",     //補充レシート種別:定義値 キャッシュカウントレポート
	3:  "supply_addchange",     //補充レシート種別:定義値 釣銭補充
	4:  "supply_allcollect",    //補充レシート種別:定義値 全回収
	5:  "supply_safe",          //補充レシート種別:定義値 金庫回収
	7:  "supply_exchange",      //補充レシート種別:定義値 両替
	8:  "supply_bill",          //補充レシート種別:定義値 紙幣補充
	9:  "supply_changecollect", //補充レシート種別:定義値 釣銭回収
	10: "report_safetotal",     //補充レシート種別:定義値 現在有高レポート
	11: "change_coin_unit1",    //補充レシート種別:定義値 通常硬貨ユニット交換
	12: "change_coin_unit2",    //補充レシート種別:定義値 予備硬貨ユニット交換
	13: "change_coin_unit_all", //補充レシート種別:定義値 全硬貨ユニット交換
	14: "supply_coin_manual",   //補充レシート種別:定義値 硬貨手動追加
	15: "cashsales_collect",    //補充レシート種別:定義値 現金売上金回収
	16: "report_coinunit",      //補充レシート種別:定義値 硬貨ユニット補充差分レポート
	17: "supply_exchangebill",  //補充レシート種別:定義値 紙幣逆両替レポート
	18: "report_summary",       //補充レシート種別:定義値 精算機日計レシート（FIT-B NEXTクリニック向け）
	19: "supply_refill",        //補充レシート種別:定義値 追加補充
	20: "retrieve_collectbin",  //補充レシート種別:定義値 回収庫から回収
	21: "retrieve_quantity",    //補充レシート種別:定義値 指定枚数回収
	22: "reverse_exchange",     //補充レシート種別:定義値 逆両替
}

func GetReportName(reportId int) string {
	if val, ok := ReportName[reportId]; ok {
		return val
	}

	return ""
}

func NewRequestSupply(info RequestInfo,
	typename string,
	termNo int,
	supplyDate int,
	supplyTime int,
	numInfoTbl []int,
) *RequestSupply {
	return &RequestSupply{
		RequestInfo: info,
		Typename:    typename,
		TermNo:      termNo,
		SupplyDate:  supplyDate,
		SupplyTime:  supplyTime,
		NumInfoTbl:  numInfoTbl,
	}
}
