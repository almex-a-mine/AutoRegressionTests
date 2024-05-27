package domain

// 印刷制御通信
const (
	SUMMARY_SALES        = 0  //補充レシート種別:精算機別日計表
	REPORT_CASHCOUNT     = 1  //補充レシート種別:キャッシュカウントレポート
	SUPPLY_ADDCHANGE     = 3  //補充レシート種別:釣銭補充
	SUPPLY_ALLCOLLECT    = 4  //補充レシート種別:全回収
	SUPPLY_SAFE          = 5  //補充レシート種別:金庫回収
	SUPPLY_EXCHANGE      = 7  //補充レシート種別:両替
	SUPPLY_BILL          = 8  //補充レシート種別:紙幣補充
	SUPPLY_CHANGECOLLECT = 9  //補充レシート種別:釣銭回収
	REPORT_SAFETOTAL     = 10 //補充レシート種別:現在有高レポート
	CHANGE_COINTUNIT1    = 11 //補充レシート種別:通常硬貨ユニット交換(青カセット交換)
	CHANGE_COINTUNIT2    = 12 //補充レシート種別:予備硬貨ユニット交換(ピンクカセット交換)
	CHANGE_COINTUNIT_ALL = 13 //補充レシート種別:全硬貨ユニット交換（青ピンクカセット交換）
	SUPPLY_COIN_MANUAL   = 14 //補充レシート種別:硬貨手動追加
	CASHSALES_COLLECT    = 15 //補充レシート種別:現金売上金回収
	REPORT_COINUNIT      = 16 //補充レシート種別:硬貨ユニット補充差分レポート
	SUPPLY_EXCHANGEBILL  = 17 //補充レシート種別:紙幣逆両替レポート
	REPORT_SUMMARY       = 18 //補充レシート種別:精算機日計レシート（FIT-B NEXTクリニック向け）
)

// const (
// 	NAME_SUMMARY_SALES               = "summary_sales"               //補充レシート種別:定義値 精算機別日計表
// 	NAME_REPORT_CASHCOUNT            = "report_cashcount"            //補充レシート種別:定義値 キャッシュカウントレポート
// 	NAME_SUPPLY_ADDCHANGE            = "supply_addchange"            //補充レシート種別:定義値 釣銭補充
// 	NAME_SUPPLY_ALLCOLLECT           = "supply_allcollect"           //補充レシート種別:定義値 全回収
// 	NAME_SUPPLY_SAFE                 = "supply_safe"                 //補充レシート種別:定義値 金庫回収
// 	NAME_SUPPLY_EXCHANGE             = "supply_exchange"             //補充レシート種別:定義値 両替
// 	NAME_SUPPLY_BILL                 = "supply_bill"                 //補充レシート種別:定義値 紙幣補充
// 	NAME_SUPPLY_CHANGECOLLECT        = "supply_changecollect"        //補充レシート種別:定義値 釣銭回収
// 	NAME_REPORT_SAFETOTAL            = "report_safetotal"            //補充レシート種別:定義値 現在有高レポート
// 	NAME_CHANGE_COINTUNIT1           = "change_coin_unit1"           //補充レシート種別:定義値 通常硬貨ユニット交換
// 	NAME_CHANGE_COINTUNIT2           = "change_coin_unit2"           //補充レシート種別:定義値 予備硬貨ユニット交換
// 	NAME_CHANGE_COINTUNIT_ALL        = "change_coin_unit_all"        //補充レシート種別:定義値 全硬貨ユニット交換
// 	NAME_SUPPLY_COIN_MANUAL          = "supply_coin_manual"          //補充レシート種別:定義値 硬貨手動追加
// 	NAME_CASHSALES_COLLECT           = "cashsales_collect"           //補充レシート種別:定義値 現金売上金回収
// 	NAME_REPORT_COINUNIT             = "report_coinunit"             //補充レシート種別:定義値 硬貨ユニット補充差分レポート
// 	NAME_SUPPLY_EXCHANGEBILL         = "supply_exchangebill"         //補充レシート種別:定義値 紙幣逆両替レポート
// )

const (
	SUMMARY_SALES_START     = 1 //精算機別日計表：開始
	SUMMARY_SALES_PRINT_REQ = 2 //精算機別日計表：補充レシート要求
)

const (
	BEFORE_CROSSING_PROCESS  = 1 //カウントレポート：締め前有高
	AFTER_CROSSING_PROCESS   = 2 //カウントレポート：締め後有高
	INDATA_CROSSING_PROCESS  = 3 //カウントレポート：締め中入金枚数
	OUTDATA_CROSSING_PROCESS = 4 //カウントレポート：締め中出金枚数
	COMPLETE_CLOSING_PROCESS = 5 //カウントレポート：締め処理完了
)
const (
	NOW_AMOUNT            = 1 //レシート状態：締め前有高
	INDATA_EXCOUNTTBL     = 2 //レシート状態：入金
	OUTDATA_EXCOUNTTBL    = 3 //レシート状態：出金
	AFTER_AMOUNT          = 4 //レシート状態：入出金後
	TOTAL_INDATA          = 5 //レシート状態：入金額計
	TOTAL_OUTDATA         = 6 //レシート状態：出金額計
	COMPLET_REPLENISHMENT = 7 //補充完了
)

const (
	PRINT_DATA_CASH_COUNT   = 66 //キャッシュカウントレポート:拡張金種別枚数要素数
	OVERFLOW_DATA           = 6  //あふれ要素数
	PRINT_CASH_DATA         = 13 //印刷用データ13個
	PRINT_CASH_DATA_SIX     = 6  //印刷用データ6個
	PRINT_CASH_DATA_FOUR    = 4  //印刷用データ4個
	PRINT_CASH_DATA_SIXTEEN = 16 //印刷用データ8個
	PRINT_CASH_DATA_EXCOUNT = 26 //印刷用データ26個
)
