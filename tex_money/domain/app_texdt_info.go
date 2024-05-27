package domain

// 稼働データ
// 履歴分類コード
// 1:通常 2:補充 3:業務 4:エラー 5:入金 6:出金 7:その他
const (
	NORMAL_HIS     = 1
	ADD            = 2
	BUSINESS_HIS   = 3
	ERROR_HIS      = 4
	DEPOSIT_HIS    = 5
	WITHDRAWAL_HIS = 6
	OTHER          = 7
)

//分類情報種別

//状態遷移コード
/*62:初期補充 63:追加補充
65:枚数払出 66:金額払出 67:逆両替
71:手動補充 */
const (
	STATUSCODE_INITIAL_ADD       = 62
	ADDITIONAL_REPLENISHMENT_STC = 63
	NUMBER_PAYOUT                = 65
	AMOUNTP_AYOUT                = 66
	REVERSE_EXCHANGE_STC         = 67
	MANUAL_REPLENISHMENT         = 71
)

//操作オペレータコード
/*リモート操作及びシステム動作の場合は"0000"
不明な場合は"????"*/
const (
	SISTEM_DOING = "0000"
	UNKNOUWN     = "????"
)

// 0:正常 1:エラー 2:ワーニング
const (
	NORMAL = 0
	ERROR  = 1
)

// エラー発生状態:0:金銭取扱なし 1:入金中 2:出金中
const (
	NO_MONEY_HANDLING = 0
	ERROR_DEPOSITING  = 1
	WITHDRAWING       = 2
)

// 取扱中:0:取扱中 1:取扱中止 2:業務終了 3:メンテ中 4:取引中
const (
	HANDLED           = 0
	DISCONTINUED      = 1
	END_OF_BUSINESS   = 2
	UNDER_MAINTENANCE = 3
	TRADING           = 4
)

// セキュリティ状態:0:停止 1:作動
const (
	SECURITY_STOP   = 0
	SECURITY_SDOINT = 1
)

// 扉状態：0:クローズ 1:オープン
const (
	CLOSE = 0
	OPEN  = 1
)

// キーSW状態:0:自己診断 1:保守 2:通常 3:取扱中止 4:補充 5:業務 6:エラー確認
const (
	SELF_DIAGNOSIS    = 0
	MAINTENANCE       = 1
	KEYS_NORMAL       = 2
	KEYS_DISCONTINUED = 3
	REPLENISHMENT     = 4
	BUSINESS          = 5
	ERROR_CHECK       = 6
)

// 従業員呼出状態:0:呼出なし 1:呼出中
const (
	NO_CALL = 0
	CALLING = 1
)

//決済方法モード
/*0:指定なし 1:現金のみ 2:クレジットのみ 3:Jデビットのみ
4:クレジット優先＋Jデビット 5:Jデビット優先＋クレジット
6:現金優先＋クレジット 7:現金優先＋Jデビット 8:現金優先＋クレジット＋Jデビット*/
const (
	NOT_SPECIFIED                = 0
	CASH_ONLY                    = 1
	CREDIT_ONLY                  = 2
	J_DEBIT_ONLY                 = 3
	CREDIT_PRIORITY_J_DEBIT      = 4
	J_DEBIT_PRIORITY_CREDIT      = 5
	CASH_PRIORITY_CREDIT         = 6
	CASH_PRIORITY_JDEBIT         = 7
	CASH_PRIORITY_CREDIT_J_DEBIT = 8
)

// 決済方法
// 0:現金 1:クレジット 2:Jデビット 3:QRコード決済 4:電子マネー 5:その他
const (
	PAYWAY_CASH      = 0
	CREDIT           = 1
	J_DEBIT          = 2
	QR_CODE_PAYMENT  = 3
	ELECTRONIC_MONEY = 4
	OTHERS           = 5
)

// 入出金種別
// 0:入金 1:出金予定 2:出金 3:エラー入金 4:エラー出金 5:エラー出金済
const (
	DEPOSIT              = 0
	SCHEDULED_WITHDRAWAL = 1
	WITHDRAWAL_NOMAL     = 2
	DEPOSITED_ERROR      = 3
	WITHDRAWAL_ERROR     = 4
	WITHDRAWAL_AREADY    = 5
)

// 入出金種別配列
var InfoType = [6]int{
	DEPOSIT,
	SCHEDULED_WITHDRAWAL,
	WITHDRAWAL_NOMAL,
	DEPOSITED_ERROR,
	WITHDRAWAL_ERROR,
	WITHDRAWAL_AREADY}

// 分類情報種別:0:現金有高 1:釣銭可能 2:初期補充 3:取引入金 4:取引出金 5:取引差引 6:補充入金 7:補充出金 8:補充差引 9:売上金回収 10:入金可能
const (
	CASH_AVAILABLE           = 0
	CHANGE_AVAILABLE         = 1
	INITIAL_REPLENISHMENT    = 2
	TRANSACTION_DEPOSIT      = 3
	TRANSACTION_WITHDRAWAL   = 4
	TRANSACTION_BALANCE      = 5
	REPLENISHMENT_DEPOSIT    = 6
	REPLENISHMENT_WITHDRAWAL = 7
	REPLENISHMENT_BALANCE    = 8
	SALES_MONEY_COLLECT      = 9
	DEPOSIT_NUMBER           = 10
)

var SortType = [11]int{CASH_AVAILABLE,
	CHANGE_AVAILABLE,
	INITIAL_REPLENISHMENT,
	TRANSACTION_DEPOSIT,
	TRANSACTION_WITHDRAWAL,
	TRANSACTION_BALANCE,
	REPLENISHMENT_DEPOSIT,
	REPLENISHMENT_WITHDRAWAL,
	REPLENISHMENT_BALANCE,
	SALES_MONEY_COLLECT,
	DEPOSIT_NUMBER}

// Tex_money_固有 分類情報種別
const (
	// 集計中入出金差引
	AGGREGATE_WITHDRAWAL = 90
	// デバイス有高
	DEVICE_AVAILABLE = 91
	// 有高差引
	AVAILABLE_BALANCE = 92
)
