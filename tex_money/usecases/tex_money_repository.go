package usecases

import "tex_money/domain"

type TexMoneyHandlerRepository interface {
	Start()
	Stop()
	SetSequence(texCon *domain.TexContext, status int)
	GetSequence(texCon *domain.TexContext) (retSequence int)
	SetMoneySetting(updateData *domain.MoneySetting)
	GetMoneySetting() *domain.MoneySetting
	MakeCurrentStatusTbl(texCon *domain.TexContext) [domain.CASH_TYPE_SHITEI]int

	NewRequestInfo(texCon *domain.TexContext) domain.RequestInfo // RequestInfo生成

	SetFlagCollect(texCon *domain.TexContext, on bool)  // 設定:回収要求時のフラグ
	GetFlagCollect(texCon *domain.TexContext) int       // 取得：回収要求時のフラグ
	SetFlagExchange(texCon *domain.TexContext, on bool) // 設定：両替要求時のフラグ
	GetFlagExchange(texCon *domain.TexContext) int      // 取得：両替要求時のフラグ

	SetExchangeTargetDevice(texCon *domain.TexContext, target int) // 両替要求時のターゲットデバイス(紙幣のみの場合、紙幣だけで両替する等の判定に利用)
	SetExchangePattern(texCon *domain.TexContext, pattern int)
	SetExchangeCashControlId(texCon *domain.TexContext, id string) // 両替時に入金時のCashControlIdを登録する

	SetAmountRequestCreate(texCon *domain.TexContext, reqInfo *domain.RequestSetAmount) domain.RequestCashctlSetAmount // CashCtrl向けrequest_set_amount生成

	UnreturnedAndSalesCollect(texCon *domain.TexContext, reqInfo *domain.RequestSetAmount) (resInfo domain.RequestCashctlSetAmount)                                                                                   // 非還流庫回収and売上金回収
	Collect(texCon *domain.TexContext, reqInfo *domain.RequestMoneyCollect) (resInfo domain.RequestCollectStart, resInfo2 domain.RequestCollectStop, resInfo3 domain.RequestOutStart, resInfo4 domain.RequestOutStop) //回収
	MiddleAndSalesCollect(texCon *domain.TexContext, reqInfo *domain.RequestMoneyCollect) (domain.RequestOutStart, domain.RequestOutStop)                                                                             // 途中回収And売上金回収

	OutCashStart(texCon *domain.TexContext, reqInfo domain.RequestOutCash) (resInfo domain.RequestOutStart) //開始:取引出金要求
	SensorFailedNoticeOutData(texCon *domain.TexContext)                                                    //有高不足で取引出金要求が失敗した場合にnotice_outdataを出すための処理

	SetCollectSales(texCon *domain.TexContext, amount int, countTbl [domain.CASH_TYPE_SHITEI]int, exCountTbl [domain.EXTRA_CASH_TYPE_SHITEI]int)                      // 売上金回収情報
	SetOverflowCollectSales(texCon *domain.TexContext, update bool, amount int, countTbl [domain.CASH_TYPE_SHITEI]int, exCountTbl [domain.EXTRA_CASH_TYPE_SHITEI]int) // 売上金回収のあふれ情報
	SetErrorFromRequest(texCon *domain.TexContext, statusError bool, errorCode string, errorDetail string)                                                            //エラー状態のセット
	GetErrorFromRequest(texCon *domain.TexContext) (statusError bool, errorCode string, errorDetail string)                                                           //入出金管理：要求応答で起きたエラー状態取得
	SetTexmyNoticeStatus(texCon *domain.TexContext, statusCash domain.StatusCash)                                                                                     //現金入出金機ステータス通知検知

	//現金入出金機制御
	SensorCashctlNoticeInStatus(texCon *domain.TexContext, stuInfo domain.InStatus) bool        //入金状況
	SensorCashctlNoticeOutStatus(texCon *domain.TexContext, stuInfo domain.OutStatus)           //出金状況
	SensorCashctlNoticeCollectStatus(texCon *domain.TexContext, x interface{})                  //回収状況
	SensorCashctlNoticeAmountStatus(texCon *domain.TexContext, stuInfo domain.AmountStatus)     //有高状況
	SensorCashctlNoticeExchangeStatus(texCon *domain.TexContext, x interface{})                 //両替ステータス通知
	SensorCashctlNoticeStatus(texCon *domain.TexContext, stuInfo domain.NoticeStatus)           //現金入出金機状況
	CheckAmountLimit(texCon *domain.TexContext, statusCash domain.StatusCash) domain.StatusCash //リミット有高チェック

	RecvCashctlALLRequest(texCon *domain.TexContext, x interface{}) //全ての要求の応答が受けるメソッド

	RecvSetAmountNoticeAmountStatus(texCon *domain.TexContext, amStatus domain.AmountStatus, onlyAmountChangeSequence int) bool //有高枚数変更要求の有高ステータスが返ってきた

	//稼働データ管理
	TexdtInfoSave(texCon *domain.TexContext, resultStatus domain.ResultGetTermInfoNow)        //イニシャル時のデータ保存
	RequestReportSafeInfo(texCon *domain.TexContext) (resInfo domain.RequestReportSafeInfo)   //稼働データ管理:金庫情報遷移記録要求
	RecvRequestReportSafeInfo(texCon *domain.TexContext, resInfo domain.ResultReportSafeInfo) //稼働データ管理:金庫情報遷移記録要求Recv

	//印刷要求
	RecvPrintALLRequest(texCon *domain.TexContext, x interface{}) //全ての要求の応答を受けるメソッド

	//入出金管理：コールバック
	RegisterCallbackNoticeIndata(callbackFunc func(texCon *domain.TexContext, noticeInfo *domain.StatusIndata))           //入金ステータスコールバック登録
	RegisterCallbackNoticeOutdata(callbackFunc func(texCon *domain.TexContext, noticeInfo *domain.StatusOutdata))         //出金ステータスコールバック登録
	RegisterCallbackNoticeCollectdata(callbackFunc func(texCon *domain.TexContext, noticeInfo *domain.StatusCollectData)) //回収ステータスコールバック登録
	RegisterCallbackNoticeAmountData(callbackFunc func(texCon *domain.TexContext, noticeInfo *domain.StatusAmount))       //有高ステータスコールバック登録
	RegisterCallbackNoticeStatusdata(callbackFunc func(texCon *domain.TexContext, noticeInfo *domain.StatusCash))         //入出金機ステータススコールバック登録
	RegisterCallbackNoticeReportStatusdata(callbackFunc func(texCon *domain.TexContext, noticeInfo *domain.StatusReport))
	RegisterCallbackNoticeExchangeStatusdata(callbackFunc func(texCon *domain.TexContext, noticeInfo *domain.StatusExchange))

	SensorZeroNoticeCollect(texCon *domain.TexContext)                                                         //cashTblが0円での回収要求が来た場合のnotice_collectを出すための処理
	SensorOverflowOnlyNoticeCollect(texCon *domain.TexContext, amount int, cashTbl [10]int, exCashTbl [26]int) // あふれ回収だけある売上金回収でnotice_collectを出すため
	InStatusExchange(texCon *domain.TexContext)
	SetTexmyNoticeExchangedata(texCon *domain.TexContext, ok bool) // 入出金管理：両替ステータス通知

	InitialDiscrepantOn(shouldChangeCheck bool) // 釣銭不一致監視を開始するチェックフラグ(routerからのアクセス専用)
	InitialDiscrepanctStartOne()                // 起動時だけ、メンテナスモードを無視して釣銭不一致判定を実施するためのロジック

	SetTexmyNoticeIndata(texCon *domain.TexContext, ok bool)      //入金ステータス通知 送信判定
	SetTexmyNoticeOutdata(texCon *domain.TexContext, ok bool)     //出金ステータス通知 送信判定
	SetTexmyNoticeCollectdata(texCon *domain.TexContext, ok bool) //回収ステータス通知 送信判定
	SetTexmyNoticeAmountData(texCon *domain.TexContext)           //有高ステータス通知 送信判定
}
