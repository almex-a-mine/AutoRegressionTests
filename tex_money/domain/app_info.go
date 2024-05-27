package domain

// システム情報
const (
	// AppName APP NAME
	AppName = "tex_money"
	// サービス情報

	SrvName     = "texmy"
	DspName     = "ALMEX TEX Money"
	Description = ""
)

// ローカルログ出力
const (
	LOG_STOP_DEBUG = false
	LOG_STOP_MUTEX = false
	LOG_STOP_TRACE = false
	LOG_STOP_INFO  = false
	LOG_STOP_WARN  = false
)

// リクエストID識別データ
const (
	RequestIdData = "TexMoney_"
)

// Topic
const (
	// TOPIC_TEXMONEY_BASE 入出金管理
	TOPIC_TEXMONEY_BASE = "/tex/unifunc/money"
	// TOPIC_CASHCTL_BASE 現金入出金機制御
	TOPIC_CASHCTL_BASE = "/tex/helper/cashctl"
	// TOPIC_HELPERPRINT_BASE 印刷制御
	TOPIC_HELPERPRINT_BASE = "/almex/helper/print"
	// TOPIC_HELPERDBDATA_BASE 稼働データ管理
	TOPIC_HELPERDBDATA_BASE = "/tex/helper/dbdata"
	// TOPIC_UNIFUNCSTATUS_BASE 精算機状態管理
	TOPIC_UNIFUNCSTATUS_BASE = "/tex/unifunc/status"
	// TOPIC_UNIFUNCPAYMENT_BASE 精算取引管理
	TOPIC_UNIFUNCPAYMENT_BASE = "/tex/unifunc/payment"
)

// 応答待機時間 1000:1秒 30000:30秒 60000:1分 600000:1時間
const WAIT_TIME = 5000         //5000:5秒
const WAIT_TIME_INITIAL = 3000 //3000:3秒 リトライ時間
const RETRY_COUNTER = 30       //リトライ回数
const WAIT_TIME_OUT = 1000     //1秒 出金口のtrue待ち

const (
	//初期補充
	//動作モード
	// MONEY_INIT_CANCEL 取消
	MONEY_INIT_CANCEL = 0 //取消
	// MONEY_INIT_START 開始
	MONEY_INIT_START = 1 //開始
	// MONEY_INIT_CONFIRM 確定
	MONEY_INIT_CONFIRM = 2 //確定

	//両替
	//動作モード
	MONEY_EXCHANGE_CANCEL  = 0 //取消
	MONEY_EXCHANGE_START   = 1 //開始
	MONEY_EXCHANGE_CONFIRM = 2 //確定
	//追加補充
	//動作モード
	MONEY_ADDREPLENISH_CANCEL  = 0 //取消
	MONEY_ADDREPLENISH_START   = 1 //開始
	MONEY_ADDREPLENISH_CONFIRM = 2 //確定
	//回収
	//回収モード
	MONEY_COLLECT_MIDDLE           = 0 //途中回収
	MONEY_COLLECT_ALL              = 1 //全回収
	MONEY_COLLECT_SALESMONEY       = 2 //売上金回収
	MONEY_COLLECT_INREJECT         = 3 //全回収（リジェクト庫含）
	MONEY_COLLECT_MIDDLE_AND_SALES = 4 //途中回収(回収分に売上金回収を含む)
	//取引入金
	//動作モード
	MONEY_PAYCASH_CANCEL  = 0 //取消
	MONEY_PAYCASH_START   = 1 //開始
	MONEY_PAYCASH_CONFIRM = 2 //確定
	MONEY_PAYCASH_END     = 3 //終了
	//取引出金
	//動作モード
	MONEY_OUTCASH_STOP                 = 0 //停止
	MONEY_OUTCASH_START                = 1 //開始
	MONEY_REFUND_BALANCE_PAYMENT_START = 2 //返金残払出開始
	//現在枚数変更要求
	MONEY_SETAMOUNT_ANALOG_COLLECT               = 0 //手動補充/回収
	MONEY_SETAMOUNT_REJECTBOX_COLLECT            = 1 //リジェクトボックス回収
	MONEY_SETAMOUNT_UNRETURNED_COLLECT           = 2 //非還流庫回収
	MONEY_SETAMOUNT_UNRETURNED_AND_SALES_COLLECT = 3 //非還流庫回収and売上金回収
)

//現金入出金機制御ステータス通知項目
/*0:待機中、1:リセット中、2:入金許可中、3:入金禁止中
4:出金中、5:回収中、6:精査中
7:クリア処理中、8:警告中、9:エラー中
10:入金中、11:受取待ち中、12:保守動作中、13:処理停止中*/
const (
	NORMAL_OPERATION_MODE = 1 //1:通常運用モード
	CLEANING_MODE         = 2 //2:クリーニングモード
	MAINTENANCE_MOEDE     = 3 //3:メンテナンスモード
	//動作状態
	WAITING_TEXMY           = 0
	RESETTING               = 1
	DEPOSIT_ALLOWED         = 2
	DEPOSIT_PROHIBITED      = 3
	WITHDRAWAL              = 4
	COLLECTING              = 5
	UNDER_REVIEW            = 6
	DURING_CLEAR_PROCESSING = 7
	DURING_WARNING          = 8
	DURING_ERROR            = 9
	DEPOSITING              = 10
	WAITING_FOR_RECEIPT     = 11
	MAINTENANCE_IN_PROGRESS = 12
	PROCESSING_STOPPED      = 13
)

const (
	CASH_TYPE_SHITEI       = 10 //指定金種枚数
	EXTRA_CASH_TYPE_SHITEI = 26 //拡張金種枚数
	CASH_TYPE_UI           = 16 //UI指定金種枚数
)

const (
	REVERSE_EXCHAGE               = 0 //両替パターン:逆両替
	ONE_CASHTYPE_EXCHAGE          = 1 //両替パターン:全て1系金種で両替
	ONE_FIVE_CASHTYPE_EXCHAGE     = 2 //両替パターン:1系，5系混在で両替
	NUMBER_OF_WITHDRAW_DESIGNATED = 3 //両替パターン:出金枚数指定
)

const (
	OUT_AMOUNT_SHITEI = 0 //出金種別:金額指定回収
	OUT_SITEI_MAISUU  = 1 //出金種別:枚数指定回収
)

const (
	COLLECT_SITEI_AMOUNT = 0 //回収種別:金額指定回収
	COLLECT_SITEI_MAISUU = 1 //回収種別:枚数指定回収
)

// 途中回収
const (
	STOP  = 0 //動作モード:停止
	START = 1 //動作モード:開始
)

// 取引・有高
const (
	TRADE    = 1 //取引
	AMOUNT   = 2 //有高
	OUTTRADE = 3 //出金
)

const (
	AMOUNT_SHITEI_OUT         = 0 //現金入出金制御 出金要求:金額指定出金
	CASH_NUMBER_OUT           = 1 //現金入出金制御 出金要求:枚数指定出金
	WITHDRAW_TO_OUTLET        = 0 //現金入出金制御 出金要求:出金口に出金
	COLLECT_TO_COLLECTION_BOX = 1 //現金入出金制御 出金要求:回収庫に回収
)

// 結果通知コード
// 101:入金開始，102:入金データ通知，103:入金完了，109:入金異常
// 201:出金開始，202:出金データ通知，203:出金停止，204:出金完了，209:出金異常
// 300:回収機能無し，301:回収開始，302:回収データ通知，303:回収停止，304:回収完了，309:回収異常
// 501:有高処理開始，502:有高データ通知，504:有高処理完了，509:有高異常
const (
	//入金
	IN_DEPOSIT_START             = 101 //入金開始
	IN_RECEIPT_DATA_NOTIFICATION = 102 //入金データ通知
	IN_PAYMENT_COMPLETED         = 103 //入金完了
	IN_PAYMENT_PROHIBIT          = 104 //入金終了
	IN_PAYMENT_ERROR             = 109 //入金異常
	//入金取消
	CANCEL_DEPOSIT_START             = 201 //入金取消開始
	CANCEL_RECEIPT_DATA_NOTIFICATION = 202 //入金取消データ通知，
	CANCEL_PAYMENT_STOP              = 203 //入金取消停止
	CANCEL_PAYMENT_COMPLETE          = 204 //入金取消完了
	CANCEL_PAYMENT_ERR               = 209 //入金取消異常
	//出金
	OUT_DEPOSIT_START             = 201 //出金開始
	OUT_RECEIPT_DATA_NOTIFICATION = 202 //出金データ通知
	OUT_PAYMENT_STOP              = 203 //出金停止
	OUT_PAYMENT_COMPLETED         = 204 //出金完了
	OUT_PAYMENT_ERROR             = 209 //出金異常
	//回収
	COL_NOTHING_FUNC              = 300 //回収機能無し
	COL_DEPOSIT_START             = 301 //回収開始
	COL_RECEIPT_DATA_NOTIFICATION = 302 //回収データ通知
	COL_PAYMENT_STOP              = 303 //回収停止
	COL_PAYMENT_COMPLETED         = 304 //回収完了
	COL_PAYMENT_ERROR             = 309 //回収異常
	//有高
	AMO_DEPOSIT_START             = 501 //有高処理開始
	AMO_RECEIPT_DATA_NOTIFICATION = 502 //有高データ通知
	AMO_PAYMENT_COMPLETED         = 504 //有高処理完了
	AMO_PAYMENT_ERROR             = 509 //有高異常
)

const (
	NOMAL           = 0 //入金データ通知 拡張動作状況 0:通常(取消以外)
	CANCEL_START    = 1 //入金データ通知 拡張動作状況 1:取消開始
	CANCEL_COMPLETE = 4 //入金データ通知 拡張動作状況 4:取消完了
	CANCEL_ERROR    = 9 //入金データ通知 拡張動作状況 9:取消異常
)

const (
	FINISH_IN_START            = 6  //各要求完了:入金開始要求
	FINISH_IN_END              = 1  //各要求完了:入金終了要求
	FINISH_REPORT_SAFEINFO     = 2  //各要求完了:金庫情報遷移記録要求
	FINISH_REPORT_SUPPLY       = 3  //各要求完了:補充完了記録要求
	FINISH_CHANGE_SUPPLY       = 4  //各要求完了:状態変更要求（補充完了）
	FINISH_PRINT_CHANGE_SUPPLY = 5  //各要求完了:補充レシート印刷要求
	FINISH_OUT_START           = 7  //各要求完了:出金開始要求
	FINISH_COLLECT_START       = 8  //各要求完了:回収開始要求
	FINISH_OUT_END             = 9  //各要求完了:出金停止要求
	FINISH_GET_TERMINFO        = 10 //各要求完了:現在端末取得要求
	FINISH_SET_AMOUNT          = 11 //各要求完了:現在枚数変更要求
	SKIP_SEND_SET_AMOUNT       = 12 //各要求完了:現在枚数変更要求（送信スキップ）
)

// 入金終了要求:statusMode
const (
	CANCEL_IN_END  = 0 //取消（入金禁止＆返却）
	START_IN_END   = 1 //開始（入金禁止のみ）
	CONFIRM_IN_END = 2 //確定（入金禁止＆収納）
)

// 初期補充要求
const (
	SETAMOUNT = 1 //現在有高枚数を投げる時
	MONEYINIT = 2 //初期補充のみ投げる時
)

// 金庫情報遷移記録要求
const (
	SORTTYPE = 10 //分類情報種別項目数
)

// 現在端末状況取得要求
const (
	SALESTYPE = 4 //売上種別
)

// 売上情報取得要求
const (
	RETURNVALUE  = 0 //戻り値取得
	SALESCOUNTER = 1 //売上金回収要求カウンター
)

// 要求有無フラグ
const (
	REQUEST_NOTHING = 0 //要求無し
	REQUEST_HAVE    = 1 //要求あり
)

// リミット条件チェックパターン
const (
	LIMIT_CONTAIN_EQUAL   = 0 //枚数≧あふれ設定値、枚数≦不足設定値。不足W・不足Eが共に0設定の場合、不足エラー・不足注意のチェックをしない(既存動作)
	LIMIT_NOCONTAIN_EQUAL = 1 //枚数＞あふれ設定値、枚数＜不足設定値。
)

const (
	NO_CLOSING_PROCESS  = 0 //締め処理中状態:締め処理中でない
	NOW_CLOSING_PROCESS = 1 //締め処理中状態:締め処理中
)

const (
	CLOSING_PROCESS_CASHID        = "ClosingProcess_001"  //締め処理キャッシュID
	NORMAL_COIN_REPLENISH_CASHID  = "ReplenishProcess_11" //補充モードキャッシュID：通常硬貨ユニット交換（青カセット）
	SUB_COIN_REPLENISH_CASHID     = "ReplenishProcess_12" //補充モードキャッシュID：予備硬貨ユニット交換（ピンクカセット）
	ALL_COIN_REPLENISH_CASHID     = "ReplenishProcess_13" //補充モードキャッシュID：全硬貨ユニット交換（青・ピンクカセット）
	SUPPLY_COINUNIT_MANUAL_CASHID = "ReplenishProcess_14" //補充モードキャッシュID：硬貨手動追加
	MANUAL_BILL_REPLENISH_CASHID  = "ReplenishProcess_8"  //補充モードキャッシュID：紙幣逆両替レポート（紙幣手動補充）
)

const (
	SET_AMOUNT_INDATA_ONE     = "SetAmountIndata_1"  //有高変更要求時のnotice_inのCashId
	SET_AMOUNT_OUTDATA_ONE    = "RetAmountOutdata_1" //有高変更要求時のnotice_outのCashId
	NOTICE_COLLECT_CASHID_ONE = "NoticeCollect_1"    //有高変更要求時のnotice_outのCashId
	OUT_CASH_ONE              = "OutCash_1"          //取引出金要求時のCashId（払出予定金額0円用）
	SET_AMOUNT_ONE            = "SetAmount_1"        //有高変更要求時のCashId
)

const (
	NO_REPLENISH_PROCESS  = 0 //補充モード状態:補充モード中でない
	NOW_REPLENISH_PROCESS = 1 //補充モード状態:補充モード中
)

const (
	REPLENISHMENT_MODE = 1   //保守業務モード:補充モード
	CLOSING_MODE       = 100 //保守業務モード:締めモード
)

const (
	EXCHANGE_STATUS_IN  = 1 //両替ステータス：入金
	EXCHANGE_STATUS_OUT = 2 //両替ステータス：出金
)

// 保守業務モード
var MODE = [2]int{REPLENISHMENT_MODE, CLOSING_MODE}

const (
	REPLENISHMENT_START = 1 //保守業務動作状況:補充開始
	REPLENISHMENT_END   = 2 //保守業務動作状況:補充終了
	CLOSING_START       = 3 //保守業務動作状況:締め開始
	CLOSING_END         = 4 //保守業務動作状況:締め終了
)

const (
	BEFORE_AMOUNT_COUNT_TBL    = 1 //処理前有高金種配列
	BEFORE_REPLENISH_COUNT_TBL = 2 //処理前補充入金金種配列
	REPLENISH_COUNT_TBL        = 3 //補充入金金種配列
	BEFORE_COLLECT_COUNT_TBL   = 4 //処理前回収金種配列
	COLLECT_COUNT_TBL          = 5 //回収金種配列
	AFTER_AMOUNT_COUNT_TBL     = 6 //処理後有高金種配列
	SALES_COLLECT_TBL          = 7 //売上金回収配列
)
