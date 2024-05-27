package domain

type RequestCoincassetteControl struct {
	RequestInfo  RequestInfo                 `json:"requestInfo"`
	CoinCassette int                         `json:"coinCassette"`
	ControlMode  int                         `json:"controlMode"`
	AmountCount  [EXTRA_CASH_TYPE_SHITEI]int `json:"amountCount"`
}

type ResultCoincassetteControl struct {
	RequestInfo           RequestInfo                 `json:"requestInfo"`
	Result                bool                        `json:"result"`                //処理結果
	ErrorCode             string                      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail           string                      `json:"errorDetail,omitempty"` //エラー詳細
	DifferenceTotalAmount int                         `json:"differenceTotalAmount"`
	DifferenceExCountTbl  [EXTRA_CASH_TYPE_SHITEI]int `json:"differenceExCountTbl"`
	BeforeExCountTbl      [EXTRA_CASH_TYPE_SHITEI]int `json:"beforeExCountTbl"`
	AfterExCountTbl       [EXTRA_CASH_TYPE_SHITEI]int `json:"afterExCountTbl"`
	ExchangeExCountTbl    [EXTRA_CASH_TYPE_SHITEI]int `json:"exchangeExCountTbl"`
}

type CoinCassette struct {
	DifferenceTotalAmount int
	DifferenceExCountTbl  [EXTRA_CASH_TYPE_SHITEI]int
	BeforeExCountTbl      [EXTRA_CASH_TYPE_SHITEI]int
	AfterExCountTbl       [EXTRA_CASH_TYPE_SHITEI]int
	ExchangeExCountTbl    [EXTRA_CASH_TYPE_SHITEI]int
}

type Cassette struct {
	M10000 int // m:メイン
	M5000  int
	M2000  int
	M1000  int
	M500   int
	M100   int
	M50    int
	M10    int
	M5     int
	M1     int
	S500   int // s:サブ
	S100   int
	S50    int
	S10    int
	S5     int
	S1     int
	A10000 int // a:あふれ
	A5000  int
	A2000  int
	A1000  int
	A500   int
	A100   int
	A50    int
	A10    int
	A5     int
	A1     int
}

func NewCoinCassette(
	differenceTotalAmount int,
	differenceExCountTbl [EXTRA_CASH_TYPE_SHITEI]int,
	beforeExCountTbl [EXTRA_CASH_TYPE_SHITEI]int,
	afterExCountTbl [EXTRA_CASH_TYPE_SHITEI]int,
	exchangeExCountTbl [EXTRA_CASH_TYPE_SHITEI]int,
) *CoinCassette {
	return &CoinCassette{
		DifferenceTotalAmount: differenceTotalAmount,
		DifferenceExCountTbl:  differenceExCountTbl,
		BeforeExCountTbl:      beforeExCountTbl,
		AfterExCountTbl:       afterExCountTbl,
		ExchangeExCountTbl:    exchangeExCountTbl,
	}
}
