package domain

// 精算機状態管理：状態変更要求（精算完了）
type RequestChangePayment struct {
	RequestInfo        RequestInfo  `json:"requestInfo"`
	TradeType          int          `json:"tradeType"`                    // 取引種別
	CancelFlag         bool         `json:"cancelFlag,omitempty"`         // 取消処理フラグ
	RequestId          string       `json:"requestId"`                    // 精算問合せ番号
	DispRequestId      string       `json:"dispRequestId,omitempty"`      // 表示用精算問合せ番号
	ResultCode         int          `json:"resultCode"`                   // 精算結果コード
	PaymentDetail      string       `json:"paymentDetail,omitempty"`      // 精算内容
	ReceiptCount       int          `json:"receiptCount,omitempty"`       // 領収書枚数
	SpecificationCount int          `json:"specificationCount,omitempty"` // 明細書枚数
	OthersPrintCount1  int          `json:"othersPrintCount1,omitempty"`  // その他発行枚数1
	OthersPrintCount2  int          `json:"othersPrintCount2,omitempty"`  // その他発行枚数2
	FilePath           string       `json:"filePath,omitempty"`           // 精算結果ファイルパス名
	InfoTrade          InfoTrade    `json:"infoTrade,omitempty"`          // 取引情報
	InfoSales          InfoSales    `json:"infoSales,omitempty"`          // 売上情報
	InfoSafe           InfoSafeInfo `json:"infoSafe,omitempty"`           // 金庫情報
}

type ResultChangePayment struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Result      bool        `json:"result"`                //処理結果
	ErrorCode   string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string      `json:"errorDetail,omitempty"` //エラー詳細
}

type InfoSales struct {
	SalesAmount   int            `json:"salesAmount"`   // 売上金額合計
	ExchangeTotal int            `json:"exchangeTotal"` // 両替金額合計
	SalesTypeTbl  []SalesTypeTbl `json:"salesTypeTbl"`  // 売上種別情報
}

func NewRequestChangePayment(info RequestInfo, tradeType int, requestId string, resultCode int, infoTrade InfoTrade, infoSales InfoSales, infoSafe InfoSafeInfo) *RequestChangePayment {
	return &RequestChangePayment{
		RequestInfo: info,
		TradeType:   tradeType,
		RequestId:   requestId,
		ResultCode:  resultCode,
		InfoTrade:   infoTrade,
		InfoSales:   infoSales,
		InfoSafe:    infoSafe}
}
