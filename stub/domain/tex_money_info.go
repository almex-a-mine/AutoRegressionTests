package domain

// 初期補充情報
type RequestMoneyInit struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	CashControlId string      `json:"cashControlId,omitempty"` //入出金制御管理番号
	ModeOperation int         `json:"modeOperation"`           //運用モード
	CountClear    bool        `json:"countClear"`              //入金枚数クリア
	TargetDevice  int         `json:"targetDevice"`            //対象デバイス
	StatusMode    int         `json:"statusMode"`              //動作モード
	CashTbl       [15]int     `json:"cashTbl"`                 //初期枚数
}

type ResultMoneyInit struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	Result        bool        `json:"result"`                  //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`     //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"`   //エラー詳細
	CashControlId string      `json:"cashControlId,omitempty"` //入出金制御管理番号
}

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

// 追加補充情報
type RequestMoneyAddReplenish struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	CashControlId string      `json:"cashControlId"` //入出金制御管理番号
	ModeOperation int         `json:"modeOperation"` //運用モード
	CountClear    bool        `json:"countClear"`    //入金枚数クリア
	TargetDevice  int         `json:"targetDevice"`  //対象デバイス
	StatusMode    int         `json:"statusMode"`    //動作モード
}

type ResultMoneyAddReplenish struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
	CashControlId string      `json:"cashControlId"`         //入出金機制御管理番号
}

// 回収情報(途中回収要求／全回収要求／売上金回収要求)
type RequestMoneyCollect struct {
	RequestInfo   RequestInfo       `json:"requestInfo"`
	CashControlId string            `json:"cashControlId"` //入出金制御管理番号
	CollectMode   int               `json:"collectMode"`   //回収モード
	OutType       int               `json:"outType"`       //払出方向
	StatusMode    int               `json:"statusMode"`    //動作モード
	SalesAmount   int               `json:"salesAmount"`   //売上金額
	CashTbl       [CASH_TYPE_UI]int `json:"cashTbl"`       //回収枚数
}

type ResultMoneyCollect struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
	CashControlId string      `json:"cashControlId"`         //入出金機制御管理番号
}

// 現在枚数変更情報
type RequestSetAmount struct {
	RequestInfo   RequestInfo                 `json:"requestInfo"`
	CashControlId string                      `json:"cashControlId"`
	OperationMode int                         `json:"operationMode"` //操作モード
	CashTbl       [EXTRA_CASH_TYPE_SHITEI]int `json:"cashTbl"`       //指定枚数

}

type ResultSetAmount struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	CashControlId string      `json:"cashControlId"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
}

// 現金入出金機制御ステータス情報
type RequestStatusCash struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	CashControlId string      `json:"cashControlId"` //入出金機制御管理番号
}

type ResultStatusCash struct {
	RequestInfo         RequestInfo        `json:"requestInfo,omitempty"`
	Result              bool               `json:"result"`                //通信結果
	CashControlId       string             `json:"cashControlId"`         //入出金機制御管理番号
	StatusReady         bool               `json:"statusReady"`           //制御状態
	StatusMode          int                `json:"statusMode"`            //動作状態
	StatusLine          bool               `json:"statusLine"`            //通信状態
	StatusError         bool               `json:"statusError"`           //エラー状態
	ErrorCode           string             `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail         string             `json:"errorDetail,omitempty"` //エラー詳細
	StatusCover         bool               `json:"statusCover"`           //トビラ状態
	StatusAction        int                `json:"statusAction"`          //動作状態
	StatusInsert        bool               `json:"statusInsert"`          //入金口状態
	StatusExit          bool               `json:"statusExit"`            //出金口状態
	StatusRjbox         bool               `json:"statusRjbox"`           //リジェクトBOX
	BillStatusTbl       TexmyBillStatusTbl `json:"billStatusTbl"`         //紙幣ステータス情報
	CoinStatusTbl       CoinStatusTbl      `json:"coinStatusTbl"`         //硬貨ステータス情報
	BillResidueInfoTbl  []BillResidueInfo  `json:"billResidueInfoTbl"`    //紙幣残留情報
	CoinResidueInfoTbl  []CoinResidueInfo  `json:"coinResidueInfoTbl"`    //硬貨残留情報
	DeviceStatusInfoTbl []string           `json:"deviceStatusInfoTbl"`   //デバイス詳細情報
	WarningInfoTbl      []int              `json:"warningInfoTbl"`        //警告情報
}

// 紙幣ステータス情報
type TexmyBillStatusTbl struct {
	StatusUnitSet     bool `json:"statusUnitSet"`               //ユニットセット状態
	StatusInCassette  bool `json:"statusInCassette"`            //補充カセット状態
	StatusOutCassette bool `json:"statusOutCassette"`           //回収カセット状態
	StatusAmountCount int  `json:"statusAmountCount,omitempty"` //有高枚数状態
}

// 硬貨ステータス情報
type CoinStatusTbl struct {
	StatusUnitSet     bool `json:"statusUnitSet"`               //ユニットセット状態
	StatusInCassette  bool `json:"statusInCassette"`            //補充カセット状態
	StatusOutCassette bool `json:"statusOutCassette"`           //回収カセット状態
	StatusAmountCount int  `json:"statusAmountCount,omitempty"` //有高枚数状態
}

// 紙幣残留情報
type BillResidueInfo struct {
	Title  string `json:"title"`  //管理名称
	Status bool   `json:"status"` //状態
}

// 硬貨残留情報
type CoinResidueInfo struct {
	Title  string `json:"title"`  //管理名称
	Status bool   `json:"status"` //状態
}

// 取引入金情報
type RequestPayCash struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	CashControlId string      `json:"cashControlId"` //入出金機制御管理番号
	ModeOperation int         `json:"modeOperation"` //運用モード
	CountClear    bool        `json:"countClear"`    //入金枚数クリア
	TargetDevice  int         `json:"targetDevice"`  //対象デバイス
	StatusMode    int         `json:"statusMode"`    //動作モード
}

type ResultPayCash struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
	CashControlId string      `json:"cashControlId"`         //入出金機制御管理番号
}

// 取引出金
type RequestOutCash struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	CashControlId string      `json:"cashControlId"` //入出金制御管理番号
	StatusMode    int         `json:"statusMode"`    //動作モード
	OutData       int         `json:"outData"`       //出金金額
}

type ResultOutCash struct {
	RequestInfo         RequestInfo           `json:"requestInfo"`
	Result              bool                  `json:"result"`                //処理結果
	ErrorCode           string                `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail         string                `json:"errorDetail,omitempty"` //エラー詳細
	CashControlId       string                `json:"cashControlId"`         //入出金機制御管理番号
	PaymentPlanCountTbl [CASH_TYPE_SHITEI]int `json:"paymentPlanCountTbl"`   // 出金予定枚数
	StatusMode          int                   `json:"statusMode"`            //動作モード
}

// 有高枚数要求
type RequestAmountCash struct {
	RequestInfo RequestInfo `json:"requestInfo"`
}

type ResultAmountCash struct {
	RequestInfo RequestInfo                 `json:"requestInfo"`
	Result      bool                        `json:"result"`                //処理結果
	ErrorCode   string                      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string                      `json:"errorDetail,omitempty"` //エラー詳細
	Amount      int                         `json:"amount"`                //金額
	CountTbl    [CASH_TYPE_SHITEI]int       `json:"countTbl"`              //通常金種別枚数
	ExCountTbl  [EXTRA_CASH_TYPE_SHITEI]int `json:"exCountTbl"`            //拡張金種別枚数
}

// 入出金レポート印刷情報
type RequestPrintReport struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	FilePath    string      `json:"filePath"` //精算結果ファイルパス名
	ReportId    int         `json:"reportId"` //レポート管理番号
}

type ResultPrintReport struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Result      bool        `json:"result"`                //処理結果
	ErrorCode   string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string      `json:"errorDetail,omitempty"` //エラー詳細
	SlipPrintId string      `json:"slipPrintId"`           //レポート印刷制御管理番号
}

// 売上金情報要求
type RequestSalesInfo struct {
	RequestInfo RequestInfo `json:"requestInfo"`
}

type ResultSalesInfo struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
	SalesAmount   int         `json:"salesAmount"`           //売上金額
	SalesComplete int         `json:"salesComplete"`         //売上金回収済金額
	SalesCount    int         `json:"salesCount"`            //売上金回収回数
}

// 入出金データクリア要求
type RequestClearCashInfo struct {
	RequestInfo RequestInfo `json:"requestInfo"`
}

type ResultClearCashInfo struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Result      bool        `json:"result"`                //処理結果
	ErrorCode   string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string      `json:"errorDetail,omitempty"` //エラー詳細
}

//保守業務モード要求
type RequestMaintenanceMode struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Mode        int         `json:"mode"`   //保守業務モード
	Action      bool        `json:"action"` //動作要求
}

type ResultMaintenanceMode struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Result      bool        `json:"result"`                //処理結果
	ErrorCode   string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string      `json:"errorDetail,omitempty"` //エラー詳細
}

//逆両替算出要求
type RequestReverseExchangeCalculation struct {
	RequestInfo     RequestInfo `json:"requestInfo"`
	ExchangeType    int         `json:"exchangeType"`
	OverflowCashbox bool        `json:"overflowCashbox"`
	Amount          int         `json:"amount"`
}

type ResultReverseExchangeCalculation struct {
	RequestInfo        RequestInfo                  `json:"requestInfo"`
	Result             bool                         `json:"result"`                //処理結果
	ErrorCode          string                       `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail        string                       `json:"errorDetail,omitempty"` //エラー詳細
	TargetAmount       int                          `json:"targetAmount"`
	TargetExCountTbl   *[EXTRA_CASH_TYPE_SHITEI]int `json:"targetExCountTbl,omitempty"`
	ExchangeExCountTbl [EXTRA_CASH_TYPE_SHITEI]int  `json:"exchangeExCountTbl"`
}

//硬貨カセット操作要求
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

//金庫情報取得要求
type RequestGetSafeInfo struct {
	RequestInfo RequestInfo `json:"requestInfo"`
}

type ResultGetSafeInfo struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
	SalesComplete int         `json:"salesComplete"`         //売上金回収済
	SalesCount    int         `json:"salesCount"`            //売上金回収回数
	CollectCount  int         `json:"collectCount"`          //回収操作回数
	InfoSafe      InfoSafe    `json:"infoSafe"`              //金庫情報
}

// 金庫情報
type InfoSafe struct {
	CurrentStatusTbl [CASH_TYPE_SHITEI]int `json:"currentStatusTbl"` //通常金種別状況
	SortInfoTbl      [14]SortInfoTbl       `json:"sortInfotbl"`      //分類情報
}

// 金庫情報(0:現金有高~9:売上金回収までの情報)
type SafeInfo struct {
	SalesCompleteAmount int             //売上金回収済
	SalesCompleteCount  int             //売上金回収回数
	CollectCount        int             //回収操作回数
	SortInfoTbl         [11]SortInfoTbl //分類情報
}

// 分類別金庫情報
type SortInfoTbl struct {
	SortType   int                         `json:"sortType"`   //分類情報種別
	Amount     int                         `json:"amount"`     //金額
	CountTbl   [CASH_TYPE_SHITEI]int       `json:"countTbl"`   //通常金種別枚数
	ExCountTbl [EXTRA_CASH_TYPE_SHITEI]int `json:"exCountTbl"` //拡張金種別枚数
}

//金銭設定登録要求
type RequestRegisterMoneySetting struct {
	RequestInfo         RequestInfo          `json:"requestInfo"`
	ChangeReserveCount  *ChangeReserveCount  `json:"changeReserveCount"`  //釣銭準備金枚数
	ChangeShortageCount *ChangeShortageCount `json:"changeShortageCount"` //不足枚数
	ExcessChangeCount   *ExcessChangeCount   `json:"excessChangeCount"`   //あふれ枚数
}

type ResultRegisterMoneySetting struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Result      bool        `json:"result"`                //処理結果
	ErrorCode   string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string      `json:"errorDetail,omitempty"` //エラー詳細
}

//金庫情報設定
type MoneySetting struct {
	ChangeReserveCount  ChangeReserveCount  `json:"changeReserveCount"`  //釣銭準備金枚数
	ChangeShortageCount ChangeShortageCount `json:"changeShortageCount"` //不足枚数
	ExcessChangeCount   ExcessChangeCount   `json:"excessChangeCount"`   //あふれ枚数
}

type ChangeReserveCount struct {
	LastRegistDate string `json:"lastRegistDate"` //最終登録日付
	LastRegistTime string `json:"lastRegistTime"` //最終登録時刻
	M10000Count    int    `json:"m10000Count"`    //10000円枚数
	M5000Count     int    `json:"m5000Count"`     //5000円枚数
	M2000Count     int    `json:"m2000Count"`     //2000円枚数
	M1000Count     int    `json:"m1000Count"`     //1000円枚数
	M500Count      int    `json:"m500Count"`      //500円枚数
	M100Count      int    `json:"m100Count"`      //100円枚数
	M50Count       int    `json:"m50Count"`       //50円枚数
	M10Count       int    `json:"m10Count"`       //10円枚数
	M5Count        int    `json:"m5Count"`        //5C円枚数
	M1Count        int    `json:"m1Count"`        //1C円枚数
	S500Count      int    `json:"s500Count"`      //500円枚数(サブ)
	S100Count      int    `json:"s100Count"`      //100円枚数(サブ)
	S50Count       int    `json:"s50Count"`       //5C円枚数(サブ)
	S10Count       int    `json:"s10Count"`       //10円枚数(サブ)
	S5Count        int    `json:"s5Count"`        //5C円枚数(サブ)
	S1Count        int    `json:"s1Count"`        //1C円枚数(サブ)
}

type ChangeShortageCount struct {
	LastRegistDate  string             `json:"lastRegistDate"` //最終登録日付
	LastRegistTime  string             `json:"lastRegistTime"` //最終登録時刻
	RegisterDataTbl [2]RegisterDataTbl `json:"registerDataTbl"`
}

type RegisterDataTbl struct {
	AlertLevel  int `json:"alertLevel"`  //アラートレベル
	M10000Count int `json:"m10000Count"` //10000円枚数
	M5000Count  int `json:"m5000Count"`  //5000円枚数
	M2000Count  int `json:"m2000Count"`  //2000円枚数
	M1000Count  int `json:"m1000Count"`  //1000円枚数
	M500Count   int `json:"m500Count"`   //500円枚数
	M100Count   int `json:"m100Count"`   //100円枚数
	M50Count    int `json:"m50Count"`    //50円枚数
	M10Count    int `json:"m10Count"`    //10円枚数
	M5Count     int `json:"m5Count"`     //5C円枚数
	M1Count     int `json:"m1Count"`     //1C円枚数
	S500Count   int `json:"s500Count"`   //500円枚数(サブ)
	S100Count   int `json:"s100Count"`   //100円枚数(サブ)
	S50Count    int `json:"s50Count"`    //5C円枚数(サブ)
	S10Count    int `json:"s10Count"`    //10円枚数(サブ)
	S5Count     int `json:"s5Count"`     //5C円枚数(サブ)
	S1Count     int `json:"s1Count"`     //1C円枚数(サブ)
}

type ExcessChangeCount struct {
	LastRegistDate    string               `json:"lastRegistDate"` //最終登録日付
	LastRegistTime    string               `json:"lastRegistTime"` //最終登録時刻
	ExRegisterDataTbl [2]ExRegisterDataTbl `json:"registerDataTbl"`
}

type ExRegisterDataTbl struct {
	AlertLevel  int `json:"alertLevel"`  //アラートレベル
	M10000Count int `json:"m10000Count"` //10000円枚数
	M5000Count  int `json:"m5000Count"`  //5000円枚数
	M2000Count  int `json:"m2000Count"`  //2000円枚数
	M1000Count  int `json:"m1000Count"`  //1000円枚数
	M500Count   int `json:"m500Count"`   //500円枚数
	M100Count   int `json:"m100Count"`   //100円枚数
	M50Count    int `json:"m50Count"`    //50円枚数
	M10Count    int `json:"m10Count"`    //10円枚数
	M5Count     int `json:"m5Count"`     //5C円枚数
	M1Count     int `json:"m1Count"`     //1C円枚数
	S500Count   int `json:"s500Count"`   //500円枚数(サブ)
	S100Count   int `json:"s100Count"`   //100円枚数(サブ)
	S50Count    int `json:"s50Count"`    //5C円枚数(サブ)
	S10Count    int `json:"s10Count"`    //10円枚数(サブ)
	S5Count     int `json:"s5Count"`     //5C円枚数(サブ)
	S1Count     int `json:"s1Count"`     //1C円枚数(サブ)
	BillOverBox int `json:"billOverBox"` //全紙幣
	CoinOverBox int `json:"coinOverBox"` //全硬貨
}

//金銭設定取得要求
type RequestGetMoneySetting struct {
	RequestInfo RequestInfo `json:"requestInfo"`
}

type ResultGetMoneySetting struct {
	RequestInfo         RequestInfo         `json:"requestInfo"`
	Result              bool                `json:"result"`                //処理結果
	ErrorCode           string              `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail         string              `json:"errorDetail,omitempty"` //エラー詳細
	ChangeReserveCount  ChangeReserveCount  `json:"changeReserveCount"`    //釣銭準備金枚数
	ChangeShortageCount ChangeShortageCount `json:"changeShortageCount"`   //不足枚数
	ExcessChangeCount   ExcessChangeCount   `json:"excessChangeCount"`     //あふれ枚数
}

//精査モード要求
type RequestScrutiny struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	CashControlId string      `json:"cashControlId"`
	TargetDevice  int         `json:"targetDevice"`
}

type ResultScrutiny struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
	CashControlId string      `json:"cashControlId"`
}

// 入金データ通知情報
type StatusIndata struct {
	CashControlId  string                      `json:"cashControlId"`          //入出金制御管理番号
	StatusAction   bool                        `json:"statusAction"`           //動作状況
	StatusResult   *bool                       `json:"statusResult,omitempty"` //入金結果
	Amount         int                         `json:"amount"`                 //金額
	CountTbl       [CASH_TYPE_SHITEI]int       `json:"countTbl"`               //通常金種別枚数
	ExCountTbl     [EXTRA_CASH_TYPE_SHITEI]int `json:"exCountTbl"`             //拡張金種別枚数
	ErrorCode      string                      `json:"errorCode,omitempty"`    //エラーコード
	ErrorDetail    string                      `json:"errorDetail,omitempty"`  //エラー詳細
	StatusActionEx int                         `json:"statusActionEx"`         //拡張動作状況
}

// 出金データ通知情報
type StatusOutdata struct {
	CashControlId string                      `json:"cashControlId"`          //入出金制御管理番号
	StatusAction  bool                        `json:"statusAction"`           //動作状況
	StatusResult  *bool                       `json:"statusResult,omitempty"` //出金結果
	Amount        int                         `json:"amount"`                 //金額
	CountTbl      [CASH_TYPE_SHITEI]int       `json:"countTbl"`               //通常金種別枚数
	ExCountTbl    [EXTRA_CASH_TYPE_SHITEI]int `json:"exCountTbl"`             //拡張金種別枚数
	ErrorCode     string                      `json:"errorCode,omitempty"`    //エラーコード
	ErrorDetail   string                      `json:"errorDetail,omitempty"`  //エラー詳細
}

// 回収データ情報
type StatusCollectData struct {
	CashControlId string                      `json:"cashControlId"`          //入出金制御管理番号
	StatusAction  bool                        `json:"statusAction"`           //動作状況
	StatusResult  *bool                       `json:"statusResult,omitempty"` //回収結果
	Amount        int                         `json:"amount"`                 //金額
	CountTbl      [CASH_TYPE_SHITEI]int       `json:"countTbl"`               //通常金種別枚数
	ExCountTbl    [EXTRA_CASH_TYPE_SHITEI]int `json:"exCountTbl"`             //拡張金種別枚数
	ErrorCode     string                      `json:"errorCode,omitempty"`    //エラーコード
	ErrorDetail   string                      `json:"errorDetail,omitempty"`  //エラー詳細
}

// 有高データ情報
type StatusAmount struct {
	Amount      int                         `json:"cashType"`              //金額
	CountTbl    [CASH_TYPE_SHITEI]int       `json:"countTbl"`              //通常金種別枚数
	ExCountTbl  [EXTRA_CASH_TYPE_SHITEI]int `json:"exCountTbl"`            //拡張金種別枚数
	ErrorCode   string                      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string                      `json:"errorDetail,omitempty"` //エラー詳細
}

//現金入出金機制御ステータス通知
type StatusCash struct {
	CashControlId       string             `json:"cashControlId"`         //入出金機制御管理番号
	StatusReady         bool               `json:"statusReady"`           //制御状態
	StatusMode          int                `json:"statusMode"`            //動作状態
	StatusLine          bool               `json:"statusLine"`            //通信状態
	StatusError         bool               `json:"statusError"`           //エラー状態
	ErrorCode           string             `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail         string             `json:"errorDetail,omitempty"` //エラー詳細
	StatusCover         bool               `json:"statusCover"`           //トビラ状態
	StatusAction        int                `json:"statusAction"`          //動作状態
	StatusInsert        bool               `json:"statusInsert"`          //入金口状態
	StatusExit          bool               `json:"statusExit"`            //出金口状態
	StatusRjbox         bool               `json:"statusRjbox"`           //リジェクトBOX
	BillStatusTbl       TexmyBillStatusTbl `json:"billStatusTbl"`         //紙幣ステータス情報
	CoinStatusTbl       CoinStatusTbl      `json:"coinStatusTbl"`         //硬貨ステータス情報
	BillResidueInfoTbl  []BillResidueInfo  `json:"billResidueInfoTbl"`    //紙幣残留情報
	CoinResidueInfoTbl  []CoinResidueInfo  `json:"coinResidueInfoTbl"`    //硬貨残留情報
	DeviceStatusInfoTbl []string           `json:"deviceStatusInfoTbl"`   //デバイス詳細情報
	WarningInfoTbl      []int              `json:"warningInfoTbl"`        //警告情報
}

// 入出金レポート印刷ステータス情報
type StatusReport struct {
	SlipPrintId  string `json:"slipPrintId"`            //レポート印刷制御管理番号
	StatusPrint  int    `json:"statusPrint"`            //印刷状態
	CountPlan    int    `json:"countPlan"`              //出力予定枚数
	CountEnd     int    `json:"countEnd"`               //印刷完了枚数
	StatusResult *bool  `json:"statusResult,omitempty"` //印刷結果
}
