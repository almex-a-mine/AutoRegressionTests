package domain

// 両替情報
type RequestMoneyExchange struct {
	RequestInfo     RequestInfo `json:"requestInfo"`
	CashControlId   string      `json:"cashControlId"`   //入出金制御管理番号
	ModeOperation   int         `json:"modeOperation"`   //運用モード
	CountClear      bool        `json:"countClear"`      //入金枚数クリア
	TargetDevice    int         `json:"targetDevice"`    //対象デバイス
	StatusMode      int         `json:"statusMode"`      //動作モード
	ExchangePattern int         `json:"exchangePattern"` //両替パターン
	PaymentPlanTbl  []int       `json:"paymentPlanTbl"`  //出金予定枚数
}

type ResultMoneyExchange struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
	CashControlId string      `json:"cashControlId"`         //入出金機制御管理番号
}

func GetExchangeData(moneyType int, exchangePattern int) ([16]int, bool) {
	// 両替パターン=1: 全て1系金種で両替
	patternOneList := map[int][16]int{
		10000: {0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		5000:  {0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		2000:  {0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		1000:  {0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		500:   {0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		100:   {0, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0},
		50:    {0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0},
		10:    {0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0},
	}

	//  両替パターン=2: 1系,5系混在で両替
	patternTwoList := map[int][16]int{
		10000: {0, 1, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		5000:  {0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		2000:  {0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		1000:  {0, 0, 0, 0, 1, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		500:   {0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		100:   {0, 0, 0, 0, 0, 0, 1, 5, 0, 0, 0, 0, 0, 0, 0, 0},
		50:    {0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0},
		10:    {0, 0, 0, 0, 0, 0, 0, 0, 1, 5, 0, 0, 0, 0, 0, 0},
	}

	switch exchangePattern {
	case 1:
		val, ok := patternOneList[moneyType]
		return val, ok

	default:
		val, ok := patternTwoList[moneyType]
		return val, ok
	}
}
