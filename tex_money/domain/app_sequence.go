package domain

const (
	NO_SEQUENCE = 0 //シーケンス無し
	INITIAL     = 1 //イニシャル動作：現金入出金制御

	INITIAL_ADDING_START   = 2 //初期補充開始
	INITIAL_ADDING_CANCEL  = 3 //初期補充取消
	INITIAL_ADDING_CONFIRM = 4 //初期補充確定
	INITIAL_ADDING_UPDATE  = 5 //初期補充更新

	EXCHANGEING_CANCEL                    = 6  //両替取消
	REVERSE_EXCHANGEING_CONFIRM_INDATA    = 8  //両替確定入金データ時
	REVERSE_EXCHANGEING_CONFIRM_OUTDATA   = 9  //両替確定出金データ時
	ONECASHTYPE_EXCHANGE_CONFIRM          = 13 //1金種両替確定
	FIVEONECASHTYPE_EXCHANGE_CONFIRM      = 14 //1and5系金種両替確定
	NUMBER_OF_WITHDRAW_DESIGNATED_CONFIRM = 15 //出金枚数指定両替確定

	MONEY_ADD_REPLENISH_START   = 10 //追加補充開始
	MONEY_ADD_REPLENISH_CANCEL  = 11 //追加補充取消
	MONEY_ADD_REPLENISH_CONFIRM = 12 //追加補充確定

	REJECTBOXCOLLECT_START          = 16 //リジェクトボックス回収開始
	UNRETURNEDCOLLECT_START         = 17 //非還流庫回収開始
	UNRETURNED_AND_SALES_COLLECT    = 49 // 非還流庫回収開始and売上金回収
	MIDDLE_START_OUT_START          = 18 //途中回収開始出金開始
	MIDDLE_START_OUT_STOP           = 19 //途中回収開始出金停止
	MIDDLE_START_COLLECT_START      = 20 //途中回収開始回収開始
	MIDDLE_START_COLLECT_STOP       = 21 //途中回収開始回収停止
	ALLCOLLECT_START_OUT_START      = 22 //全回収開始出金開始
	ALLCOLLECT_START_OUT_STOP       = 23 //全回収開始出金停止
	ALLCOLLECT_START_COLLECT_START  = 24 //全回収開始回収開始
	ALLCOLLECT_START_COLLECT_STOP   = 25 //全回収開始回収停止
	SALESMONEY_START                = 26 //売上金回収開始
	MANUAL_REPLENISHMENT_COLLECTION = 28 //手動補充・回収

	MIDDLE_AND_SALES_COLLECT = 48

	TRANSACTION_DEPOSIT_START                = 29 //取引入金開始
	TRANSACTION_DEPOSIT_CONFIRM              = 30 //取引入金確定
	TRANSACTION_DEPOSIT_CANCEL               = 31 //取引入金取消
	TRANSACTION_DEPOSIT_END_BILL             = 32 //取引入金終了紙幣
	TRANSACTION_DEPOSIT_END_COIN             = 33 //取引入金終了硬貨
	TRANSACTION_OUT_START                    = 34 //取引出金開始
	TRANSACTION_OUT_CONFIRM                  = 35 //取引出金確定
	TRANSACTION_OUT_CANCEL                   = 36 //取引出金取消
	TRANSACTION_OUT_REFUND_PAYMENT_OUT_START = 45 //取引出金 返金残払出開始

	SALES_INFO = 37 //売上金情報開始

	SET_AMOUNT  = 38 //有高枚数変更中
	AMOUNT_CASH = 39 //有高枚数要求

	CLEAR_CASHINFO = 40 //入金データ通知クリア

	//MIDLE_ADDING = 8 //途中追加の時
	//PAYING       = 9 //取引入金の時

	INITIAL_DATABASE   = 41 //イニシャル動作：稼働データ管理
	INITIAL_STATUSMANG = 42 //イニシャル動作：精算機状態管理
	INITIAL_PRINTER    = 43 //イニシャル動作：印刷
	INITIAL_CASHCTL    = 47 //イニシャル動作：現金入出金機制御

	PRINT_SUMMARY_SALES  = 44 //入出金レポート印刷要求：精算機別日計表
	PRINT_REPORT_SUMMARY = 50 //入出金レポート印刷要求：精算機日計レシート

	SCRUTINY_START = 46 //精査モード要求開始
)

func GetSquenceDetails(i int) string {
	list := map[int]string{
		NO_SEQUENCE:                              "シーケンス無し",
		INITIAL:                                  "イニシャル動作",
		INITIAL_ADDING_START:                     "初期補充開始",
		INITIAL_ADDING_CANCEL:                    "初期補充取消",
		INITIAL_ADDING_CONFIRM:                   "初期補充確定",
		INITIAL_ADDING_UPDATE:                    "初期補充更新",
		EXCHANGEING_CANCEL:                       "両替取消",
		REVERSE_EXCHANGEING_CONFIRM_INDATA:       "逆両替確定入金データ時",
		REVERSE_EXCHANGEING_CONFIRM_OUTDATA:      "逆両替確定出金データ時",
		ONECASHTYPE_EXCHANGE_CONFIRM:             "1金種両替確定",
		FIVEONECASHTYPE_EXCHANGE_CONFIRM:         "1and5系金種両替確定",
		NUMBER_OF_WITHDRAW_DESIGNATED_CONFIRM:    "出金枚数指定両替確定",
		MONEY_ADD_REPLENISH_START:                "追加補充開始",
		MONEY_ADD_REPLENISH_CANCEL:               "追加補充取消",
		MONEY_ADD_REPLENISH_CONFIRM:              "追加補充確定",
		REJECTBOXCOLLECT_START:                   "リジェクトボックス回収開始",
		UNRETURNEDCOLLECT_START:                  "非還流庫回収開始",
		UNRETURNED_AND_SALES_COLLECT:             "非還流庫回収開始and売上金回収",
		MIDDLE_START_OUT_START:                   "途中回収開始出金開始",
		MIDDLE_START_OUT_STOP:                    "途中回収開始出金停止",
		MIDDLE_START_COLLECT_START:               "途中回収開始回収開始",
		MIDDLE_START_COLLECT_STOP:                "途中回収開始回収停止",
		ALLCOLLECT_START_OUT_START:               "全回収開始出金開始",
		ALLCOLLECT_START_OUT_STOP:                "全回収開始出金停止",
		ALLCOLLECT_START_COLLECT_START:           "全回収開始回収開始",
		ALLCOLLECT_START_COLLECT_STOP:            "全回収開始回収停止",
		SALESMONEY_START:                         "売上金回収開始",
		MANUAL_REPLENISHMENT_COLLECTION:          "手動補充・回収",
		MIDDLE_AND_SALES_COLLECT:                 "未使用",
		TRANSACTION_DEPOSIT_START:                "取引入金開始",
		TRANSACTION_DEPOSIT_CONFIRM:              "取引入金確定",
		TRANSACTION_DEPOSIT_CANCEL:               "取引入金取消",
		TRANSACTION_DEPOSIT_END_BILL:             "取引入金終了紙幣",
		TRANSACTION_DEPOSIT_END_COIN:             "取引入金終了硬貨",
		TRANSACTION_OUT_START:                    "取引出金開始",
		TRANSACTION_OUT_CONFIRM:                  "取引出金確定",
		TRANSACTION_OUT_CANCEL:                   "取引出金取消",
		TRANSACTION_OUT_REFUND_PAYMENT_OUT_START: "取引出金 返金残払出開始",
		SALES_INFO:                               "売上金情報開始",
		SET_AMOUNT:                               "有高枚数変更中",
		AMOUNT_CASH:                              "有高枚数要求",
		CLEAR_CASHINFO:                           "入金データ通知クリア",
		INITIAL_DATABASE:                         "イニシャル動作：稼働データ管理",
		INITIAL_STATUSMANG:                       "イニシャル動作：精算機状態管理",
		INITIAL_PRINTER:                          "イニシャル動作：印刷",
		INITIAL_CASHCTL:                          "イニシャル動作：現金入出金機制御",
		PRINT_SUMMARY_SALES:                      "入出金レポート印刷要求：精算機別日計表",
		PRINT_REPORT_SUMMARY:                     "入出金レポート印刷要求：精算機日計レシート",
		SCRUTINY_START:                           "精査モード要求開始",
	}

	val, ok := list[i]
	if ok {
		return val
	}
	return "シーケンス詳細無" // 任意のデフォルト値

}

// 両替シーケンス
var exchangeSequences = []int{
	EXCHANGEING_CANCEL,
	REVERSE_EXCHANGEING_CONFIRM_INDATA,
	REVERSE_EXCHANGEING_CONFIRM_OUTDATA,
	NUMBER_OF_WITHDRAW_DESIGNATED_CONFIRM,
}

func IsExchangeSequences(i int) bool {
	for _, v := range exchangeSequences {
		if v == i {
			return true
		}
	}
	return false
}
