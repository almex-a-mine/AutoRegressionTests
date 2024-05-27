package usecases

import (
	"fmt"
	"strconv"
	"sync"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/pkg/calculation"
	"time"
)

type texMoneyHandler struct {
	logger                           handler.LoggerRepository
	config                           config.Configuration
	errorMng                         ErrorManager
	safeInfoManager                  SafeInfoManager
	aggregateManager                 AggregateManager
	texMoneyNoticeManager            TexMoneyNoticeManagerRepository
	maintenanceModeMng               MaintenanceModeManager
	reqCollectFlag                   int    //回収要求フラグ：1の時回収要求を受信した0の時してない
	reqExchangeFlag                  int    //両替要求フラグ：1の時回収要求を受信した0の時してない
	reqExchangeCashControlId         string //両替用CashControlId保持
	reqIdValueCounter                int    //リクエストIDの末尾数字
	sequence                         int    //動作状態
	exchangeTargetDevice             int
	exchangePattern                  int
	beforeExchangeSafe               [10]int //両替：入金前の現在有高保存用
	callbackNoticeIndata             func(texCon *domain.TexContext, noticeInfo *domain.StatusIndata)
	callbackNoticeOutdata            func(texCon *domain.TexContext, noticeInfo *domain.StatusOutdata)
	callbackNoticeCollectdata        func(texCon *domain.TexContext, noticeInfo *domain.StatusCollectData)
	callbackNoticeAmountData         func(texCon *domain.TexContext, noticeInfo *domain.StatusAmount)
	callbackNoticeStatusdata         func(texCon *domain.TexContext, noticeInfo *domain.StatusCash)
	callbackNoticeReportStatusdata   func(texCon *domain.TexContext, noticeInfo *domain.StatusReport)
	callbackNoticeExchangeStatusdata func(texCon *domain.TexContext, noticeInfo *domain.StatusExchange)
	tempAmountData                   TempAmountData         //入出金管理：有高枚数変更用：補充・回収チェック用
	refundInfo                       domain.CashInfoTblInfo //返金残情報
	recvResultOutStartCashId         string                 //出金開始要求:入出金機制御管理番号
	tempCollectAndSales              tempCollectAndSales    // 回収ロジックで回収と売上金回収を同時に実施する際の金銭管理情報
	noticeAmountUpdateNg             bool                   // false:更新OK true:更新NG
	cashDiscrepanctFlg               bool                   // UIへの通知時に釣銭不一致チェックをするかどうか false:しない true:する
	initialDiscrepanctOn             bool                   // 起動時に監視したいタイミングまで進んだらtrueへ変更する
	initialDiscrepanctStartOne       bool                   // 初回起動時だけ、メンテナンスモードを無視して釣銭不一致判定させるためのフラグ
}

// 有高枚数変更要求実施前のデータ一時保存
type TempAmountData struct {
	Amount     int                                //金額
	CountTbl   [domain.CASH_TYPE_SHITEI]int       //通常金種別枚数
	ExCountTbl [domain.EXTRA_CASH_TYPE_SHITEI]int //拡張金種別枚数
}

// 回収ロジックで回収と売上金回収を同時に実施する際の金銭管理情報(途中回収,現在枚数変更要求で利用)
// 不要なロジックでは初期化しておく
type tempCollectAndSales struct {
	// SpecialSequence 既存動作に売上金回収カウントロジックを差し込む為の判定に利用する。
	specialSequence int
	// 回収金額合計
	collectTotalAmount int
	// 回収対象26金種
	collectTotalExCollectTbl [domain.EXTRA_CASH_TYPE_SHITEI]int
	// 売上金額
	salesAmount int
	// 売上金回収対象26金種
	salesExCountTbl [domain.EXTRA_CASH_TYPE_SHITEI]int
	// 回収金額
	collectAmount int
	// 回収金回収対象26金種
	collectExCountTbl [domain.EXTRA_CASH_TYPE_SHITEI]int
}

// 入出金管理
func NewTexmyHandler(
	logger handler.LoggerRepository,
	config config.Configuration,
	errorMng ErrorManager,
	safeInfoManager SafeInfoManager,
	aggregateManager AggregateManager,
	texMoneyNoticeManager TexMoneyNoticeManagerRepository,
	maintenanceModeMng MaintenanceModeManager,
) TexMoneyHandlerRepository {
	texMoneyHandler := texMoneyHandler{}
	texMoneyHandler.logger = logger
	texMoneyHandler.config = config
	texMoneyHandler.errorMng = errorMng
	texMoneyHandler.safeInfoManager = safeInfoManager
	texMoneyHandler.aggregateManager = aggregateManager
	texMoneyHandler.maintenanceModeMng = maintenanceModeMng
	texMoneyHandler.texMoneyNoticeManager = texMoneyNoticeManager
	texMoneyHandler.tempCollectAndSales = tempCollectAndSales{}
	return &texMoneyHandler
}

// 制御開始
func (c *texMoneyHandler) Start() {
	texCon := domain.NewTexContext(domain.RegisterTexContext{})
	c.logger.Info("【%v】texMoneyHandler Start", texCon.GetUniqueKey())
	statusCash := c.texMoneyNoticeManager.GetStatusCashData(texCon)
	statusCash.StatusMode = domain.NORMAL_OPERATION_MODE
	statusCash.StatusAction = domain.WAITING_TEXMY
	statusCash.StatusReady = true
	statusCash.StatusLine = true
	statusCash.StatusError = true
	statusCash.BillResidueInfoTbl = make([]domain.BillResidueInfo, 0)
	statusCash.CoinResidueInfoTbl = make([]domain.CoinResidueInfo, 0)
	statusCash.DeviceStatusInfoTbl = make([]string, 6)
	statusCash.WarningInfoTbl = make([]int, 10)

	c.SetTexmyNoticeStatus(texCon, statusCash)
}

// 制御終了
func (c *texMoneyHandler) Stop() {
	texCon := domain.NewTexContext(domain.RegisterTexContext{})
	c.logger.Info("【%v】texMoneyHandler Stop", texCon.GetUniqueKey())
	statusCash := c.texMoneyNoticeManager.GetStatusCashData(texCon)
	statusCash.StatusAction = domain.PROCESSING_STOPPED
	statusCash.StatusLine = false

	c.SetTexmyNoticeStatus(texCon, statusCash)
}

//////////////////////////////////////////////
// フラグ設定
//////////////////////////////////////////////

// 設定:回収要求時のフラグ
func (c *texMoneyHandler) SetFlagCollect(texCon *domain.TexContext, on bool) {
	if on {
		c.reqCollectFlag = domain.REQUEST_HAVE
	} else {
		c.reqCollectFlag = domain.REQUEST_NOTHING
	}
	c.logger.Debug("【%v】SetFlagCollect 回収要求フラグ=%v", texCon.GetUniqueKey(), on)
}

// 取得：回収要求時のフラグ
func (c *texMoneyHandler) GetFlagCollect(texCon *domain.TexContext) int {
	c.logger.Debug("【%v】GetFlagCollect 回収要求フラグ=%v", texCon.GetUniqueKey(), c.reqCollectFlag)
	return c.reqCollectFlag
}

// 設定：両替要求時のフラグ
func (c *texMoneyHandler) SetFlagExchange(texCon *domain.TexContext, on bool) {
	beforeFlag := c.reqExchangeFlag
	if on {
		c.reqExchangeFlag = domain.REQUEST_HAVE
	} else {
		c.reqExchangeFlag = domain.REQUEST_NOTHING
		c.reqExchangeCashControlId = ""
	}
	c.logger.Debug("【%v】SetFlagExchange 両替要求フラグ=%v before=%v after=%v", texCon.GetUniqueKey(), on, beforeFlag, c.reqExchangeFlag)
}

// 取得：両替要求時のフラグ
func (c *texMoneyHandler) GetFlagExchange(texCon *domain.TexContext) int {
	c.logger.Debug("【%v】GetFlagExchange 両替要求フラグ=%v", texCon.GetUniqueKey(), c.reqExchangeFlag)
	return c.reqExchangeFlag
}

//////////////////////////////////////////////
// 要求共通のデータ管理
//////////////////////////////////////////////

// 両替時のCashControlId保存
func (c *texMoneyHandler) SetExchangeCashControlId(texCon *domain.TexContext, id string) {
	c.reqExchangeCashControlId = id
}

// リクエストID採番
func (c *texMoneyHandler) RequestIdCalculation(texCon *domain.TexContext) (requestID string) {

	c.reqIdValueCounter++
	requestID = domain.RequestIdData + strconv.Itoa(c.reqIdValueCounter)
	c.logger.Debug("【%v】リクエストID採番=%s", texCon.GetUniqueKey(), requestID)
	return
}

// 処理状態のセット
func (c *texMoneyHandler) SetSequence(texCon *domain.TexContext, status int) {
	c.sequence = status
	details := domain.GetSquenceDetails(c.sequence)
	c.logger.Debug("【%v】SetSequence=%v,[%v]", texCon.GetUniqueKey(), status, details)
}

// 処理状態の取得
func (c *texMoneyHandler) GetSequence(texCon *domain.TexContext) (retSequence int) {
	details := domain.GetSquenceDetails(c.sequence)
	c.logger.Debug("【%v】GetSequence=%v,[%v]", texCon.GetUniqueKey(), c.sequence, details)
	return c.sequence
}

// RequestInfo生成
func (c *texMoneyHandler) NewRequestInfo(texCon *domain.TexContext) domain.RequestInfo {
	return domain.RequestInfo{
		ProcessID: c.config.ReqInfo.ProcessID,
		PcId:      c.config.ReqInfo.PcId,
		RequestID: c.RequestIdCalculation(texCon),
	}
}

//////////////////////////////////////////////
// 要求情報の設定
//////////////////////////////////////////////

// SetAmountRequestCreate CashCtrl向けSetAmountのRequest(operationMode:0~2)生成ロジック
func (c *texMoneyHandler) SetAmountRequestCreate(texCon *domain.TexContext, reqInfo *domain.RequestSetAmount) domain.RequestCashctlSetAmount {
	c.logger.Trace("【%v】START:texMoneyHandler SetAmountRequestCreate", texCon.GetUniqueKey())

	c.TempAmountData(texCon) //現在の有高枚数を一時セーブする

	reqAmountInfo := calculation.NewCassette(reqInfo.CashTbl)
	// 合計金額
	amount := reqAmountInfo.GetTotalAmount()
	countTbl := reqAmountInfo.ExCountTblToTenCountTbl()

	res := domain.RequestCashctlSetAmount{
		RequestInfo:   c.NewRequestInfo(texCon),
		OperationMode: reqInfo.OperationMode,
		Amount:        amount,
		CountTbl:      countTbl,
		ExCountTbl:    reqInfo.CashTbl,
	}
	c.tempCollectAndSales = tempCollectAndSales{} // 保持情報初期化
	c.logger.Trace("【%v】END:texMoneyHandler RejectBoxCollect resInfo=%v", texCon.GetUniqueKey(), res)
	return res
}

func (c *texMoneyHandler) UnreturnedAndSalesCollect(texCon *domain.TexContext, reqInfo *domain.RequestSetAmount) (resInfo domain.RequestCashctlSetAmount) {
	c.logger.Trace("【%v】START:texMoneyHandler UnreturnedAndSalesCollect", texCon.GetUniqueKey())

	defer func() {
		c.logger.Trace("【%v】END:texMoneyHandler UnreturnedAndSalesCollect resInfo=%+v", texCon.GetUniqueKey(), resInfo)
	}()

	// 売上金カウントしなければならないものを計算し、内部で保持しておく

	/// 現在有高情報を取得
	_, cashAvailableSortInfo := c.safeInfoManager.GetSortInfo(texCon, domain.CASH_AVAILABLE)
	/// 現在有高
	cashAvailableExCountTbl := calculation.NewCassette(cashAvailableSortInfo.ExCountTbl)
	/// 回収対象との差分を抽出
	collectTotalExCollectTbl := cashAvailableExCountTbl.Subtract(reqInfo.CashTbl)
	collectTotalAmount := calculation.NewCassette(collectTotalExCollectTbl).GetTotalAmount()

	/// 売上金金額の確認
	_, salesSortInfo := c.safeInfoManager.GetSortInfo(texCon, domain.TRANSACTION_BALANCE)
	salesAmount := salesSortInfo.Amount

	/// 非還流庫金額が売上金より大きい場合、売上金額を非還流庫回収金額に調整
	if salesAmount > collectTotalAmount {
		salesAmount = collectTotalAmount
	}

	/// オーバーフローのみを対象とした逆両替26金種を取得
	/// この26金種が回収時に売上金回収分としてカウントされなければならない枚数
	salesCollectExCountTbl := cashAvailableExCountTbl.OverflowOnlyExchange(salesAmount)

	/// 非還流庫から回収を先行して実施するが、非還流庫からの回収で対応できない枚数がある場合を想定し
	/// 再度、有高を計算した上で、回収対象の売上金として情報を保存する。
	salesAmountFix := calculation.NewCassette(salesCollectExCountTbl).GetTotalAmount()

	// FIT-B_No213対応
	// 回収庫回収できる金額の合計値が今回回収したい売上金合計と異なる場合
	// 売上金は要求金額とし、回収した売上金種は内部でいい感じにセットしておく。
	// 尚、この対応は売上金の内訳がどこにも利用されていない事を前提として実装する。
	if salesAmountFix != salesAmount {
		salesAmountFix = salesAmount
		fakeMoney := calculation.NewCassette([26]int{999, 999, 0, 999, 999, 999, 999, 999, 999, 999}).Exchange(salesAmount, 0)
		salesCollectExCountTbl = fakeMoney
	}

	// 売上金カウント情報を保存
	c.tempCollectAndSales = tempCollectAndSales{
		specialSequence:          domain.UNRETURNED_AND_SALES_COLLECT, // シーケンス登録
		collectTotalAmount:       collectTotalAmount,                  // 回収対象金額
		collectTotalExCollectTbl: collectTotalExCollectTbl,            // 回収対象配列
		salesAmount:              salesAmountFix,                      // 売上金
		salesExCountTbl:          salesCollectExCountTbl,              // 回収売上金対象配列
	}

	c.TempAmountData(texCon) //現在の有高枚数を一時セーブする
	resInfo.RequestInfo = c.NewRequestInfo(texCon)
	resInfo.OperationMode = domain.MONEY_SETAMOUNT_UNRETURNED_COLLECT

	reqAmountInfo := calculation.NewCassette(reqInfo.CashTbl)
	// 合計金額
	resInfo.Amount = reqAmountInfo.GetTotalAmount()
	resInfo.CountTbl = reqAmountInfo.ExCountTblToTenCountTbl()
	resInfo.ExCountTbl = reqInfo.CashTbl

	// ロジックは非還流庫からの回収を利用する
	return resInfo
}

// 回収
func (c *texMoneyHandler) Collect(texCon *domain.TexContext, reqInfo *domain.RequestMoneyCollect) (
	resInfo domain.RequestCollectStart,
	resInfo2 domain.RequestCollectStop,
	resInfo3 domain.RequestOutStart,
	resInfo4 domain.RequestOutStop) {
	c.logger.Trace("【%v】START:texMoneyHandler Collect reqInfo.OutType=%v", texCon.GetUniqueKey(), reqInfo.OutType)
	switch reqInfo.OutType {
	case domain.WITHDRAW_TO_OUTLET: //出金庫に出金
		var amount int
		var statusRjbox bool
		if reqInfo.StatusMode == domain.START {
			if reqInfo.CollectMode == 3 || reqInfo.CollectMode == 4 { //3:全回収（リジェクト庫含）, 4=途中回収(売上金含)
				statusRjbox = true
			}
			for i, valueTexHelperCash := range domain.TexHelperCash {
				amount += reqInfo.CashTbl[i] * valueTexHelperCash
			}
			resInfo3 = domain.NewRequestOutStart(c.NewRequestInfo(texCon), statusRjbox, domain.OUT_SITEI_MAISUU, amount, reqInfo.CashTbl)
			c.tempCollectAndSales = tempCollectAndSales{} // 保持情報初期化
		} else if reqInfo.StatusMode == domain.STOP {
			resInfo4 = domain.NewRequestOutStop(c.NewRequestInfo(texCon), reqInfo.CashControlId)
		}
	case domain.COLLECT_TO_COLLECTION_BOX: //回収庫に回収
		if reqInfo.StatusMode == domain.START {
			resInfo = domain.NewRequestCollectStart(c.NewRequestInfo(texCon), domain.COLLECT_SITEI_MAISUU, 0, reqInfo.CashTbl)
			c.tempCollectAndSales = tempCollectAndSales{} // 保持情報初期化
		} else if reqInfo.StatusMode == domain.STOP {
			resInfo2 = domain.NewRequestCollectStop(c.NewRequestInfo(texCon), reqInfo.CashControlId)
		}
	}
	c.logger.Trace("【%v】END:texMoneyHandler Collect", texCon.GetUniqueKey())
	return
}

// 途中回収（回収分に売上金回収を含む）
func (c *texMoneyHandler) MiddleAndSalesCollect(texCon *domain.TexContext, reqInfo *domain.RequestMoneyCollect) (domain.RequestOutStart, domain.RequestOutStop) {
	c.logger.Trace("【%v】START:texMoneyHandler MiddleAndSalesCollect", texCon.GetUniqueKey())

	reqStart := domain.RequestOutStart{}
	reqStop := domain.RequestOutStop{}

	switch reqInfo.StatusMode {

	case domain.START:
		//////////////////////////////////////////////
		//回収金額の計算
		//////////////////////////////////////////////
		// 回収予定金額
		collectTotalAmount := reqInfo.SalesAmount

		// 売上金取得
		_, salesTotal := c.safeInfoManager.GetSortInfo(texCon, domain.TRANSACTION_BALANCE)
		// 回収済売上金取得
		_, salesCollectTotal := c.safeInfoManager.GetSortInfo(texCon, domain.SALES_MONEY_COLLECT)

		// ★未回収売上金
		noCollectSalesAmount := salesTotal.Amount - salesCollectTotal.Amount

		// 回収額の調整(未回収売上金より指定金額が低い可能性がある)
		if noCollectSalesAmount > collectTotalAmount {
			noCollectSalesAmount = collectTotalAmount
		}
		// ★途中回収分金額(指定金額-未回収売上金)
		collectAmount := collectTotalAmount - noCollectSalesAmount

		//////////////////////////////////////////////
		//回収枚数の計算
		//////////////////////////////////////////////
		// 釣銭可能枚数を取得
		_, changeAvailable := c.safeInfoManager.GetSortInfo(texCon, domain.CHANGE_AVAILABLE)
		// 払出可能金種の配列を取得
		exchange := calculation.NewCassette(changeAvailable.ExCountTbl).Exchange(collectTotalAmount, 0)

		// 払出可能金種配列の再計算準備
		exchangeExCountTbl := calculation.NewCassette(exchange)
		// ★払出合計金額16金種
		countTblSixteen := exchangeExCountTbl.ExCountTblToSixteenCountTbl()

		// ★払出可能金種配列(売上金金種配列分）
		salesCollectExCountTbl := exchangeExCountTbl.Exchange(noCollectSalesAmount, 0)

		// ★払出可能金種配列(途中回収金種配列分）
		collectExCountTbl := exchangeExCountTbl.Subtract(salesCollectExCountTbl)

		// FIT-B_No213対応
		// 途中回収できる売上金額の合計値が今回回収したい売上金合計と異なる場合
		// 売上金は要求金額とし、途中回収した売上金種は内部でいい感じにセットしておく。
		// 尚、この対応は売上金の内訳がどこにも利用されていない事を前提として実装する。
		if noCollectSalesAmount != calculation.NewCassette(salesCollectExCountTbl).GetTotalAmount() {
			fakeMoney := calculation.NewCassette([26]int{999, 999, 0, 999, 999, 999, 999, 999, 999, 999}).Exchange(noCollectSalesAmount, 0)
			salesCollectExCountTbl = fakeMoney
		}

		// 回収完了時に更新する情報の保存
		c.tempCollectAndSales = tempCollectAndSales{
			specialSequence:          domain.MIDDLE_AND_SALES_COLLECT, // 途中回収and売上金回収含
			collectTotalAmount:       collectTotalAmount,              // 回収金額合計
			collectTotalExCollectTbl: exchange,                        // 払出対象26金種
			salesAmount:              noCollectSalesAmount,            // 売上金回収対象金額
			salesExCountTbl:          salesCollectExCountTbl,          // 売上金回収対象金額26配列
			collectAmount:            collectAmount,                   // 途中回収対象金額
			collectExCountTbl:        collectExCountTbl,               // 途中回収対象金額26配列
		}

		// 送信情報の生成
		reqStart = domain.NewRequestOutStart(
			c.NewRequestInfo(texCon),
			false,
			domain.OUT_SITEI_MAISUU,
			collectTotalAmount,
			countTblSixteen,
		)

	case domain.STOP:

		statusCash := c.texMoneyNoticeManager.GetStatusCashData(texCon)

		// 送信情報の生成
		reqStop = domain.NewRequestOutStop(c.NewRequestInfo(texCon), statusCash.CashControlId)

	}

	c.logger.Trace("【%v】STOP:texMoneyHandler MiddleAndSalesCollect", texCon.GetUniqueKey())

	return reqStart, reqStop
}

// 取引出金要求
// 開始
func (c *texMoneyHandler) OutCashStart(texCon *domain.TexContext, reqInfo domain.RequestOutCash) (resInfo domain.RequestOutStart) {
	defer c.logger.Trace("【%v】texMoneyHandler OutCashStart resInfo=%+v", texCon.GetUniqueKey(), resInfo)

	// 現在の釣銭可能枚数を取得
	ok, safeZero := c.safeInfoManager.GetSortInfo(texCon, 0)
	if !ok {
		c.logger.Error("【%v】- 釣銭可能枚数の取得に失敗", texCon.GetUniqueKey())
		return
	}

	countTbl, result := calculation.NewCassette(safeZero.ExCountTbl).GetOutCountTbl(reqInfo.OutData)
	if !result {
		c.logger.Error("【%v】- 有高不足:2金種以内で払出不可", texCon.GetUniqueKey())
		return
	}

	return domain.NewRequestOutStart(c.NewRequestInfo(texCon), false, domain.CASH_NUMBER_OUT, 0, countTbl)
}

// 有高不足で取引出金要求が失敗した場合にnotice_outdataを出すための処理
func (c *texMoneyHandler) SensorFailedNoticeOutData(texCon *domain.TexContext) {
	c.logger.Trace("【%v】START:取引出金失敗 出金ステータス通知 判定", texCon.GetUniqueKey())
	errorCode, errorDetail := c.errorMng.GetErrorInfo(71)

	res := domain.StatusOutdata{
		CashControlId: domain.OUT_CASH_ONE,
		ErrorCode:     errorCode,
		ErrorDetail:   errorDetail,
		StatusResult:  &domain.False,
		StatusAction:  false,
	}

	c.texMoneyNoticeManager.UpdateStatusOutData(texCon, res)

	c.SetTexmyNoticeOutdata(texCon, true) //前回値と差分がなくても必ず通知を出すためtrueをセット
	c.logger.Trace("【%v】END:取引出金失敗 出金ステータス通知 判定", texCon.GetUniqueKey())
}

// 入出金管理：入金ステータス通知 送信判定
func (c *texMoneyHandler) SetTexmyNoticeIndata(texCon *domain.TexContext, ok bool) {
	c.logger.Trace("【%v】START:入金ステータス通知 送信判定 送信=%v", texCon.GetUniqueKey(), ok)
	if ok {
		if c.callbackNoticeIndata != nil {
			statusIndata := c.texMoneyNoticeManager.GetStatusInData(texCon)
			c.callbackNoticeIndata(texCon, &statusIndata)
		}
	}
	c.logger.Trace("【%v】END:入金ステータス通知 送信判定", texCon.GetUniqueKey())
}

// 入出金管理：出金ステータス通知 送信判定
func (c *texMoneyHandler) SetTexmyNoticeOutdata(texCon *domain.TexContext, ok bool) {
	c.logger.Trace("【%v】START:出金ステータス通知 送信判定 送信=%v", texCon.GetUniqueKey(), ok)
	if ok {
		statusOutdata := c.texMoneyNoticeManager.GetStatusOutData(texCon)
		if c.callbackNoticeOutdata != nil {
			c.callbackNoticeOutdata(texCon, &statusOutdata)
		}
	}
	c.logger.Trace("【%v】END:出金ステータス通知 送信判定", texCon.GetUniqueKey())
}

// 紙幣補充画面からのnotice_exchange
func (c *texMoneyHandler) InStatusExchange(texCon *domain.TexContext) {
	c.logger.Trace("【%v】START:InStatusExchange", texCon.GetUniqueKey())
	statusExchange := c.texMoneyNoticeManager.GetStatusExchangeData(texCon)
	statusExchange.StatusAction = false
	c.texMoneyNoticeManager.UpdateStatusExchangeData(texCon, statusExchange)

	if c.callbackNoticeExchangeStatusdata != nil {
		c.callbackNoticeExchangeStatusdata(texCon, &statusExchange)
	}
}

// 入出金管理：両替ステータス通知
func (c *texMoneyHandler) SetTexmyNoticeExchangedata(texCon *domain.TexContext, ok bool) {
	c.logger.Trace("【%v】START:両替ステータス通知 送信判定 送信=%v", texCon.GetUniqueKey(), ok)
	if ok {
		statusExchange := c.texMoneyNoticeManager.GetStatusExchangeData(texCon)
		if c.callbackNoticeExchangeStatusdata != nil {
			c.callbackNoticeExchangeStatusdata(texCon, &statusExchange)
		}
	}
	c.logger.Trace("【%v】END:両替ステータス通知 送信判定", texCon.GetUniqueKey())
}

// 0円での回収要求が来た場合のnotice_collectを出すための処理
func (c *texMoneyHandler) SensorZeroNoticeCollect(texCon *domain.TexContext) {
	c.logger.Trace("【%v】START:回収ステータス通知(0円) 送信判定", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:回収ステータス通知(0円) 送信判定", texCon.GetUniqueKey())
	statusCollect := domain.StatusCollectData{
		CashControlId: domain.OUT_CASH_ONE,
		StatusAction:  false,
		StatusResult:  &domain.True,
	}

	c.texMoneyNoticeManager.UpdateStatusCollectData(texCon, statusCollect)

	c.SetTexmyNoticeCollectdata(texCon, true) //0円での回収要求が来た場合は前回値と差分がなくても必ず通知を出すためtrueをセット
	c.SetFlagCollect(texCon, false)           //回収完了時にフラグを初期化
}

func (c *texMoneyHandler) SensorOverflowOnlyNoticeCollect(texCon *domain.TexContext, amount int, cashTbl [10]int, exCashTbl [26]int) {
	c.logger.Trace("【%v】START:回収ステータス通知(回収庫) 送信判定", texCon.GetUniqueKey())
	if c.callbackNoticeCollectdata != nil {
		statusCollect := c.texMoneyNoticeManager.GetStatusCollectData(texCon)
		c.callbackNoticeCollectdata(texCon,
			&domain.StatusCollectData{
				CashControlId: statusCollect.CashControlId,
				StatusAction:  false,
				StatusResult:  &domain.True,
				Amount:        amount,
				CountTbl:      cashTbl,
				ExCountTbl:    exCashTbl,
				ErrorCode:     "",
				ErrorDetail:   "",
			})
	}
	c.logger.Trace("【%v】END:回収ステータス通知(回収庫) 送信判定", texCon.GetUniqueKey())
}

// 回収ステータス通知 送信判定
func (c *texMoneyHandler) SetTexmyNoticeCollectdata(texCon *domain.TexContext, ok bool) {
	c.logger.Trace("【%v】START:回収ステータス通知 送信判定 送信=%v", texCon.GetUniqueKey(), ok)
	if ok {
		if c.callbackNoticeCollectdata != nil {
			statusCollect := c.texMoneyNoticeManager.GetStatusCollectData(texCon)
			c.callbackNoticeCollectdata(texCon, &statusCollect)
		}
	}
	c.logger.Trace("【%v】END:回収ステータス通知 送信判定", texCon.GetUniqueKey())
}

// 入出金管理：有高ステータス通知 送信判定
func (c *texMoneyHandler) SetTexmyNoticeAmountData(texCon *domain.TexContext) {
	c.logger.Trace("【%v】START:有高ステータス通知 送信判定", texCon.GetUniqueKey())
	statusAmount := c.texMoneyNoticeManager.GetStatusAmountData(texCon)
	if c.callbackNoticeAmountData != nil {
		c.callbackNoticeAmountData(texCon, &statusAmount)
	}
	c.logger.Trace("【%v】END:有高ステータス通知 送信判定", texCon.GetUniqueKey())
}

// 入出金管理：入出金レポート印刷ステータス通知
func (c *texMoneyHandler) SetTexmyNoticeReportStatusdata(texCon *domain.TexContext, statusReport domain.StatusReport) {
	c.logger.Trace("【%v】START:入出金レポート印刷ステータス通知 送信判定", texCon.GetUniqueKey())

	ok := c.texMoneyNoticeManager.UpdateStatusReportData(texCon, statusReport)
	if ok {
		if c.callbackNoticeReportStatusdata != nil {
			c.callbackNoticeReportStatusdata(texCon, &statusReport)
		}
	}
	c.logger.Trace("【%v】END:入出金レポート印刷ステータス通知 送信判定", texCon.GetUniqueKey())
}

// 入出金管理：要求応答で起きたエラー状態
func (c *texMoneyHandler) SetErrorFromRequest(texCon *domain.TexContext, statusError bool, errorCode string, errorDetail string) {

	statusCash := c.texMoneyNoticeManager.GetStatusCashData(texCon)
	statusCash.StatusError = statusError
	statusCash.ErrorCode = errorCode
	statusCash.ErrorDetail = errorDetail

	c.SetTexmyNoticeStatus(texCon, statusCash)
	if errorCode != "" {
		c.logger.Debug("【%v】 SetError errorCode=%v errorDetail=%v", texCon.GetUniqueKey(), errorCode, errorDetail)
	}

}

func (c *texMoneyHandler) GetErrorFromRequest(texCon *domain.TexContext) (statusError bool, errorCode string, errorDetail string) {

	statusCash := c.texMoneyNoticeManager.GetStatusCashData(texCon)
	if statusCash.ErrorCode != "" {
		c.logger.Debug("【%v】texMoneyHandler GetErrorFromRequest StatusError=%v,ErrorCode=%v,ErrorDetail=%v", texCon.GetUniqueKey(), statusCash.StatusError, statusCash.ErrorCode, statusCash.ErrorDetail)
	}

	return statusCash.StatusError, statusCash.ErrorCode, statusCash.ErrorDetail
}

// 入出金管理：現金入出金機制御ステータス通知
func (c *texMoneyHandler) SetTexmyNoticeStatus(texCon *domain.TexContext, statusCash domain.StatusCash) {
	c.logger.Trace("【%v】START:現金入出金制御ステータス通知 送信判定", texCon.GetUniqueKey())
	//有高のチェックを毎回行う
	if statusCash.StatusError {
		statusCash = c.CheckAmountLimit(texCon, statusCash)
	}

	// 前回データが釣銭不一致かどうかのチェック
	beforeStatusCash := c.texMoneyNoticeManager.GetStatusCashData(texCon)

	// 両替中の場合には、両替用のIDをセットする。
	// からの場合にはセットしない
	if c.reqExchangeFlag == 1 && c.reqExchangeCashControlId != "" {
		statusCash.CashControlId = c.reqExchangeCashControlId

	}

	ok := c.texMoneyNoticeManager.UpdateStatusCashData(texCon, statusCash)

	// 金銭不一致チェックフラグがONの場合チェック
	if (c.cashDiscrepanctFlg || beforeStatusCash.ErrorCode == "TXMYE100") && c.initialDiscrepanctOn {
		c.cashDiscrepancyError(texCon)
		c.cashDiscrepanctFlg = false
		// TODO:  本採用時コメントアウト解除
		if LogicalChange {
			ok = true
			statusCash = c.texMoneyNoticeManager.GetStatusCashData(texCon)
		}

	}

	if ok {
		if c.callbackNoticeStatusdata != nil {
			c.callbackNoticeStatusdata(texCon, &statusCash)
		}
	}
	c.logger.Trace("【%v】END:現金入出金制御ステータス通知 送信判定", texCon.GetUniqueKey())
}

var inStatusMutex sync.Mutex

// 現金入出金機制御
// 入金状況
func (c *texMoneyHandler) SensorCashctlNoticeInStatus(texCon *domain.TexContext, stuInfo domain.InStatus) bool {
	inStatusMutex.Lock() //通知が連続できた場合に処理中に値が書き換わってしまうことがあったためロック処理を追加
	c.logger.Trace("【%v】START:入金ステータス通知精査", texCon.GetUniqueKey())
	defer func() {
		c.logger.Trace("【%v】END:入金ステータス通知精査", texCon.GetUniqueKey())
		inStatusMutex.Unlock()
	}()

	if (stuInfo.CoinStatusCode == 104 && stuInfo.BillStatusCode != 104) || (stuInfo.CoinStatusCode != 104 && stuInfo.BillStatusCode == 104) || stuInfo.CoinStatusCode == 102 || stuInfo.BillStatusCode == 102 {
		c.logger.Debug("【%v】- デバイス有高更新 無", texCon.GetUniqueKey())

		c.setNoticeAmountUpdateNg(texCon, true)
	} else {
		c.logger.Debug("【%v】- デバイス有高更新 有", texCon.GetUniqueKey())
		c.setNoticeAmountUpdateNg(texCon, false)

	}

	// notice_in_status でどちらかが109入金異常の場合には、紙幣・硬貨どちらも入金異常とみなす。
	// ※上位への通知時、硬貨、紙幣の両方のステータスを1つのステータスとして表現する為
	if stuInfo.CoinStatusCode == domain.IN_PAYMENT_ERROR || stuInfo.BillStatusCode == domain.IN_PAYMENT_ERROR {
		stuInfo.CoinStatusCode = domain.IN_PAYMENT_ERROR
		stuInfo.BillStatusCode = domain.IN_PAYMENT_ERROR
	}
	// どちらかが、209の場合も同様とする。
	if stuInfo.CoinStatusCode == domain.CANCEL_PAYMENT_ERR || stuInfo.BillStatusCode == domain.CANCEL_PAYMENT_ERR {
		stuInfo.CoinStatusCode = domain.CANCEL_PAYMENT_ERR
		stuInfo.BillStatusCode = domain.CANCEL_PAYMENT_ERR
	}

	statusIndata := c.texMoneyNoticeManager.GetStatusInData(texCon)

	//スタータスの設定
	switch {
	case stuInfo.CoinStatusCode == domain.IN_DEPOSIT_START || stuInfo.BillStatusCode == domain.IN_DEPOSIT_START:
		statusIndata.StatusAction = true //入金中
		statusIndata.StatusResult = nil  //入金失敗で初期化
		statusIndata.StatusActionEx = domain.NOMAL
	case stuInfo.CoinStatusCode == domain.IN_PAYMENT_COMPLETED && stuInfo.BillStatusCode == domain.IN_PAYMENT_COMPLETED,
		stuInfo.CoinStatusCode == domain.IN_PAYMENT_PROHIBIT && stuInfo.BillStatusCode == domain.IN_PAYMENT_PROHIBIT,
		stuInfo.CoinStatusCode == domain.IN_PAYMENT_COMPLETED && stuInfo.BillStatusCode == domain.CANCEL_PAYMENT_COMPLETE, //修正内容
		stuInfo.CoinStatusCode == domain.CANCEL_PAYMENT_COMPLETE && stuInfo.BillStatusCode == domain.IN_PAYMENT_COMPLETED, //入金取消開始
		stuInfo.CoinStatusCode == domain.IN_RECEIPT_DATA_NOTIFICATION && stuInfo.BillStatusCode == domain.CANCEL_PAYMENT_COMPLETE,
		stuInfo.CoinStatusCode == domain.CANCEL_PAYMENT_COMPLETE && stuInfo.BillStatusCode == domain.IN_RECEIPT_DATA_NOTIFICATION:
		statusIndata.StatusAction = false        //入金完了
		statusIndata.StatusResult = &domain.True //入金成功
		statusIndata.StatusActionEx = domain.NOMAL
	case stuInfo.CoinStatusCode == domain.IN_PAYMENT_ERROR && stuInfo.BillStatusCode == domain.IN_PAYMENT_ERROR: //確定取消情報
		statusIndata.StatusAction = false         //入金完了
		statusIndata.StatusResult = &domain.False //入金失敗
		statusIndata.StatusActionEx = domain.NOMAL
	case stuInfo.CoinStatusCode == domain.CANCEL_PAYMENT_COMPLETE && stuInfo.BillStatusCode == domain.CANCEL_PAYMENT_COMPLETE: //入金取消完了
		statusIndata.StatusAction = false                    //入金完了
		statusIndata.StatusResult = &domain.True             //入金失敗
		statusIndata.StatusActionEx = domain.CANCEL_COMPLETE //取消完了
	case stuInfo.CoinStatusCode == domain.CANCEL_PAYMENT_ERR && stuInfo.BillStatusCode == domain.CANCEL_PAYMENT_ERR: //入金取消異常
		statusIndata.StatusAction = false                 //入金完了
		statusIndata.StatusResult = &domain.False         //入金失敗
		statusIndata.StatusActionEx = domain.CANCEL_ERROR //取消エラー
	case stuInfo.CoinStatusCode == domain.CANCEL_DEPOSIT_START || stuInfo.BillStatusCode == domain.CANCEL_DEPOSIT_START: //入金取消開始
		statusIndata.StatusAction = true //入金中
		statusIndata.StatusResult = nil
		statusIndata.StatusActionEx = domain.CANCEL_START //取消開始
	default:
		c.logger.Debug("【%v】- 該当ステータス無し", texCon.GetUniqueKey())
	}

	if (stuInfo.CoinStatusCode == domain.IN_DEPOSIT_START || stuInfo.BillStatusCode == domain.IN_DEPOSIT_START) || //入金データ通知
		(stuInfo.CoinStatusCode == domain.IN_RECEIPT_DATA_NOTIFICATION || stuInfo.BillStatusCode == domain.IN_RECEIPT_DATA_NOTIFICATION) || //入金データ通知
		(stuInfo.CoinStatusCode == domain.IN_PAYMENT_COMPLETED && stuInfo.BillStatusCode == domain.IN_PAYMENT_COMPLETED) || //入金完了
		(stuInfo.CoinStatusCode == domain.IN_PAYMENT_ERROR && stuInfo.BillStatusCode == domain.IN_PAYMENT_ERROR) || //入金異常
		(stuInfo.CoinStatusCode == domain.CANCEL_PAYMENT_COMPLETE || stuInfo.BillStatusCode == domain.CANCEL_PAYMENT_COMPLETE) { //入金異常
		statusIndata.CashControlId = stuInfo.CashControlId
		statusIndata.Amount = stuInfo.Amount
		statusIndata.CountTbl = stuInfo.CountTbl
		statusIndata.ExCountTbl = stuInfo.ExCountTbl
		statusIndata.ErrorCode = stuInfo.ErrorCode
		statusIndata.ErrorDetail = stuInfo.ErrorDetail
	}
	c.logger.Debug("【%v】- 更新データ生成 入金データ通知 = %+v", texCon.GetUniqueKey(), statusIndata)

	ok := c.texMoneyNoticeManager.UpdateStatusInData(texCon, statusIndata)

	//両替中かどうかのチェック //texMoneyNoticeManagerのデータ更新後に変更
	if c.reqExchangeFlag == domain.REQUEST_HAVE {
		c.SensorCashctlNoticeExchangeStatus(texCon, stuInfo)
		return true
	}

	// もし値が前回分と同一で変化がなくても、入金開始(101)であれば、送信する。
	// リセット無しでの入金継続対応かつ、既に入金済みが請求金額を満たしている場合に、通知があがらずにUIが止まる件での対応
	if !ok && stuInfo.CoinStatusCode == domain.IN_DEPOSIT_START && stuInfo.BillStatusCode == domain.IN_DEPOSIT_START {
		ok = true
	}

	//金庫情報更新
	if (stuInfo.CoinStatusCode == domain.IN_PAYMENT_COMPLETED && stuInfo.BillStatusCode == domain.IN_PAYMENT_COMPLETED && ok) ||
		(stuInfo.CoinStatusCode == domain.IN_PAYMENT_PROHIBIT && stuInfo.BillStatusCode == domain.IN_PAYMENT_PROHIBIT && ok) {
		c.CalculationInfoSafeTbl(texCon, domain.TRADE)
	}

	/* 金庫情報更新完了前に次の要求を受信していた為、通知タイミングを金庫情報更新後に変更
	→No.345の変更で通知はすべて入金に伴う処理が終わった後に変更
	c.SetTexmyNoticeIndata(texCon, ok) //入金ステータス通知*/

	//現金入出金機制御ステータス
	statusCash := c.texMoneyNoticeManager.GetStatusCashData(texCon)
	statusCash.CashControlId = stuInfo.CashControlId //入出金管理ID
	if stuInfo.ErrorCode == "" {
		statusCash.StatusError = true
		statusCash.ErrorCode = ""
		statusCash.ErrorDetail = ""
	} else {
		statusCash.StatusError = false
		statusCash.ErrorCode = stuInfo.ErrorCode
		statusCash.ErrorDetail = stuInfo.ErrorDetail
	}
	c.SetTexmyNoticeStatus(texCon, statusCash) //ステータス通知

	return true
}

// 出金状況データ格納
func (c *texMoneyHandler) SensorCashctlNoticeOutStatus(texCon *domain.TexContext, stuInfo domain.OutStatus) {
	c.logger.Trace("【%v】START:出金状況データ格納, CoinStatusCode %v,BillStatusCode %v,reqExchangeFlag %v", texCon.GetUniqueKey(), stuInfo.CoinStatusCode, stuInfo.BillStatusCode, c.reqExchangeFlag)
	defer c.logger.Trace("【%v】END:出金状況データ格納", texCon.GetUniqueKey())

	// 出金時デバイス有高更新有無判定
	if (stuInfo.CoinStatusCode == 204 && stuInfo.BillStatusCode != 204) || (stuInfo.CoinStatusCode != 204 && stuInfo.BillStatusCode == 204) || stuInfo.CoinStatusCode == 202 || stuInfo.BillStatusCode == 202 || (stuInfo.CoinStatusCode == 201 && stuInfo.BillStatusCode == 201) {
		c.logger.Debug("【%v】- デバイス有高更新_無", texCon.GetUniqueKey())

		c.setNoticeAmountUpdateNg(texCon, true)
	} else {
		c.logger.Debug("【%v】- デバイス有高更新_有", texCon.GetUniqueKey())
		c.setNoticeAmountUpdateNg(texCon, false)
	}

	// 値の更新からnotice_out送信までの動作を担保する為
	// 抜粋して、Mutex管理できる関数へ移行
	c.noticeMutexControl(texCon, stuInfo)

	//現金入出金機制御ステータス
	statusCash := c.texMoneyNoticeManager.GetStatusCashData(texCon)
	statusCash.CashControlId = stuInfo.CashControlId //入出金管理IDのみ更新
	if stuInfo.ErrorCode == "" {
		statusCash.StatusError = true
		statusCash.ErrorCode = ""
		statusCash.ErrorDetail = ""
	} else {
		statusCash.StatusError = false
		statusCash.ErrorCode = stuInfo.ErrorCode
		statusCash.ErrorDetail = stuInfo.ErrorDetail
	}
	c.SetTexmyNoticeStatus(texCon, statusCash) //ステータス通知

	// 回収データ通知を送信すると金庫情報更新前に次の要求を受信してしまうためここに移動
	if c.reqCollectFlag == domain.REQUEST_HAVE {
		c.SensorCashctlNoticeCollectStatus(texCon, stuInfo)
	}
}

// notice_out専用MUTEX
var noticeOutMutex sync.Mutex

func (c *texMoneyHandler) noticeMutexControl(texCon *domain.TexContext, stuInfo domain.OutStatus) {
	c.logger.Trace("【%v】START:texMoneyHandler noticeMutexControl", texCon.GetUniqueKey())
	c.logger.Debug("【%v】- CoinStatusCode=%v, BillStatusCode=%v", texCon.GetUniqueKey(), stuInfo.CoinStatusCode, stuInfo.BillStatusCode)
	defer c.logger.Trace("【%v】END:texMoneyHandler noticeMutexControl", texCon.GetUniqueKey())
	// 通知情報を順序制御するように修正
	noticeOutMutex.Lock()
	defer noticeOutMutex.Unlock()

	statsusOut := c.texMoneyNoticeManager.GetStatusOutData(texCon)
	// 出金データ通知生成
	//出金データを通知する
	if (stuInfo.CoinStatusCode == domain.OUT_DEPOSIT_START || stuInfo.BillStatusCode == domain.OUT_DEPOSIT_START) || //出金開始
		(stuInfo.CoinStatusCode == domain.OUT_RECEIPT_DATA_NOTIFICATION || stuInfo.BillStatusCode == domain.OUT_RECEIPT_DATA_NOTIFICATION) { //出金データ通知
		statsusOut.StatusAction = true
		statsusOut.StatusResult = nil
	} else if stuInfo.CoinStatusCode == domain.OUT_PAYMENT_COMPLETED && stuInfo.BillStatusCode == domain.OUT_PAYMENT_COMPLETED {
		statsusOut.StatusAction = false
		statsusOut.StatusResult = &domain.True
	} else if stuInfo.CoinStatusCode == domain.OUT_PAYMENT_ERROR && stuInfo.BillStatusCode == domain.OUT_PAYMENT_ERROR {
		statsusOut.StatusAction = false
		statsusOut.StatusResult = &domain.False
	}
	statsusOut.CashControlId = stuInfo.CashControlId
	statsusOut.Amount = stuInfo.Amount
	statsusOut.CountTbl = stuInfo.CountTbl
	statsusOut.ExCountTbl = stuInfo.ExCountTbl
	statsusOut.ErrorCode = stuInfo.ErrorCode
	statsusOut.ErrorDetail = stuInfo.ErrorDetail
	c.logger.Debug("【%v】- indataStatusOutdata=%+v", texCon.GetUniqueKey(), statsusOut)

	c.texMoneyNoticeManager.UpdateStatusOutData(texCon, statsusOut)

	//両替中かどうかのチェック //texMoneyNoticeManagerのデータ更新後に変更
	if c.reqExchangeFlag == domain.REQUEST_HAVE {
		c.SensorCashctlNoticeExchangeStatus(texCon, stuInfo)
		return
	}

	c.CalculationInfoSafeTbl(texCon, domain.OUTTRADE) //金庫情報更新

	/* 金庫情報更新完了前に次の要求を受信していた為、通知タイミングを金庫情報更新後に変更
	→No.345の変更で通知はすべて出金に伴う処理が終わった後に変更
	c.SetTexmyNoticeOutdata(texCon, ok)*/

}

// 回収ステータス
func (c *texMoneyHandler) SensorCashctlNoticeCollectStatus(texCon *domain.TexContext, x interface{}) {
	c.logger.Trace("【%v】START:回収ステータス状況処理 ", texCon.GetUniqueKey())
	statusCollect := c.texMoneyNoticeManager.GetStatusCollectData(texCon)
	statusCash := c.texMoneyNoticeManager.GetStatusCashData(texCon)
	switch receiveData := x.(type) {
	case domain.OutStatus: //出金データ通知を変換する
		//回収データを通知する

		c.logger.Debug("【%v】- CoinStatusCode=%v, BillStatusCode=%v", texCon.GetUniqueKey(), receiveData.CoinStatusCode, receiveData.BillStatusCode)
		switch {
		case receiveData.CoinStatusCode == domain.OUT_PAYMENT_COMPLETED && receiveData.BillStatusCode == domain.OUT_PAYMENT_COMPLETED:
			statusCollect.StatusAction = false
			statusCollect.StatusResult = &domain.True
		case receiveData.CoinStatusCode == domain.OUT_PAYMENT_ERROR || receiveData.BillStatusCode == domain.OUT_PAYMENT_ERROR:
			statusCollect.StatusAction = false
			statusCollect.StatusResult = &domain.False
		case receiveData.CoinStatusCode == domain.OUT_DEPOSIT_START || receiveData.BillStatusCode == domain.OUT_DEPOSIT_START:
			statusCollect.StatusAction = true
			statusCollect.StatusResult = nil
		}
		statusCollect.CashControlId = receiveData.CashControlId
		statusCollect.Amount = receiveData.Amount
		statusCollect.CountTbl = receiveData.CountTbl
		statusCollect.ExCountTbl = receiveData.ExCountTbl
		statusCollect.ErrorCode = receiveData.ErrorCode
		statusCollect.ErrorDetail = receiveData.ErrorDetail

		for {
			statusCash = c.texMoneyNoticeManager.GetStatusCashData(texCon)
			time.Sleep(domain.WAIT_TIME_OUT * time.Millisecond)

			if !statusCash.StatusExit {
				break
			}
			c.logger.Debug("【%v】- 出金口紙幣抜取待", texCon.GetUniqueKey())
		}

		c.texMoneyNoticeManager.UpdateStatusCollectData(texCon, statusCollect)

		/* 金庫情報更新完了前に次の要求を受信していた為、通知タイミングを金庫情報更新後に変更
		→No.345の対応ですべてのTOPICが送信完了後にTOPICを送信するよう変更
		c.SetTexmyNoticeCollectdata(texCon, ok)*/

		//現金入出金機制御ステータス
		statusCash.CashControlId = receiveData.CashControlId //入出金管理IDのみ更新
		if statusCollect.ErrorCode == "" {
			statusCash.StatusError = true
			statusCash.ErrorCode = ""
			statusCash.ErrorDetail = ""
		} else {
			statusCash.StatusError = false
			statusCash.ErrorCode = statusCollect.ErrorCode
			statusCash.ErrorDetail = statusCollect.ErrorDetail
		}
		c.SetTexmyNoticeStatus(texCon, statusCash) //ステータス通知にデータをセット
	case domain.CollectStatus: //回収データ通知を変換する
		if (receiveData.CoinStatusCode == domain.COL_RECEIPT_DATA_NOTIFICATION && receiveData.BillStatusCode == domain.COL_RECEIPT_DATA_NOTIFICATION) || //回収データ通知
			(receiveData.CoinStatusCode == domain.COL_PAYMENT_COMPLETED && receiveData.BillStatusCode == domain.COL_PAYMENT_COMPLETED) { //回収完了
			//回収データを通知する
			statusCollect.CashControlId = receiveData.CashControlId
			if receiveData.CoinStatusCode == domain.IN_DEPOSIT_START || receiveData.CoinStatusCode == domain.IN_RECEIPT_DATA_NOTIFICATION ||
				receiveData.BillStatusCode == domain.IN_DEPOSIT_START || receiveData.BillStatusCode == domain.IN_RECEIPT_DATA_NOTIFICATION {
				statusCollect.StatusAction = true
			} else if receiveData.CoinStatusCode == domain.IN_PAYMENT_COMPLETED || receiveData.BillStatusCode == domain.IN_PAYMENT_COMPLETED {
				statusCollect.StatusAction = false
				statusCollect.StatusResult = &domain.True
			}
			statusCollect.CashControlId = receiveData.CashControlId
			statusCollect.Amount = receiveData.Amount
			statusCollect.CountTbl = receiveData.CountTbl
			statusCollect.ExCountTbl = receiveData.ExCountTbl
			statusCollect.ErrorCode = receiveData.ErrorCode
			statusCollect.ErrorDetail = receiveData.ErrorDetail
			c.logger.Debug("【%v】- indataStatusCollectdata = %v", texCon.GetUniqueKey(), statusCollect)

			ok := c.texMoneyNoticeManager.UpdateStatusCollectData(texCon, statusCollect)

			c.CalculationInfoSafeTbl(texCon, domain.TRADE) //金庫情報更新

			// 金庫情報更新完了前に次の要求を受信していた為、通知タイミングを金庫情報更新後に変更
			c.SetTexmyNoticeCollectdata(texCon, ok)
			c.SetFlagCollect(texCon, false) //回収完了時にフラグを初期化

		}
		//現金入出金機制御ステータス
		statusCash.CashControlId = receiveData.CashControlId //入出金管理IDのみ更新
		if receiveData.ErrorCode == "" {
			statusCash.StatusError = true
			statusCash.ErrorCode = ""
			statusCash.ErrorDetail = ""
		} else {
			statusCash.StatusError = false
			statusCash.ErrorCode = receiveData.ErrorCode
			statusCash.ErrorDetail = receiveData.ErrorDetail
		}
		c.SetTexmyNoticeStatus(texCon, statusCash) //ステータス通知
	}

	c.logger.Trace("【%v】END:回収ステータス状況処理", texCon.GetUniqueKey())
}

// 有高状況
func (c *texMoneyHandler) SensorCashctlNoticeAmountStatus(texCon *domain.TexContext, stuInfo domain.AmountStatus) {
	statusAmount := c.texMoneyNoticeManager.GetStatusAmountData(texCon)
	c.logger.Trace("【%v】START:有高ステータス状況処理 前回有高ステータス=%+v,CoinStatusCode=%+v,BillStatusCode=%+v", texCon.GetUniqueKey(), statusAmount, stuInfo.CoinStatusCode, stuInfo.BillStatusCode)
	if (stuInfo.CoinStatusCode == domain.AMO_RECEIPT_DATA_NOTIFICATION && stuInfo.BillStatusCode == domain.AMO_RECEIPT_DATA_NOTIFICATION) || //有高データ通知  //TODO: この行の条件を削除して問題ないか確認が必要
		((stuInfo.CoinStatusCode == domain.AMO_PAYMENT_COMPLETED || stuInfo.CoinStatusCode == domain.AMO_PAYMENT_ERROR) &&
			(stuInfo.BillStatusCode == domain.AMO_PAYMENT_COMPLETED || stuInfo.BillStatusCode == domain.AMO_PAYMENT_ERROR)) { //有高処理完了or有高異常 //TODO: 呼出し元で有高完了時（または有高異常で完了時）のみに絞り込めていればここの条件は不要

		statusAmount.Amount = stuInfo.Amount
		statusAmount.CountTbl = stuInfo.CountTbl
		statusAmount.ExCountTbl = stuInfo.ExCountTbl
		statusAmount.ErrorCode = stuInfo.ErrorCode
		statusAmount.ErrorDetail = stuInfo.ErrorDetail
		c.texMoneyNoticeManager.UpdateStatusAmountData(texCon, statusAmount)
		c.CalculationInfoSafeTbl(texCon, domain.AMOUNT) //金庫情報更新

		/*金庫情報更新完了前に次の要求を受信していた為、通知タイミングを金庫情報更新後に変更
		→No.345の変更で通知はすべて有高に伴う処理が終わった後に変更
		c.SetTexmyNoticeAmountData(texCon)
		*/
	}
	statusCash := c.texMoneyNoticeManager.GetStatusCashData(texCon)
	if stuInfo.ErrorCode == "" { //1回目のnotice_status_cashは2.1に乗ってきたエラーを格納する
		statusCash.StatusError = true
		statusCash.ErrorCode = ""
		statusCash.ErrorDetail = ""
	} else {
		statusCash.StatusError = false
		statusCash.ErrorCode = stuInfo.ErrorCode
		statusCash.ErrorDetail = stuInfo.ErrorDetail
	}
	c.SetTexmyNoticeStatus(texCon, statusCash) //ステータス通知
	c.logger.Trace("【%v】END:有高ステータス状況処理", texCon.GetUniqueKey())
}

// SensorCashctlNoticeExchangeStatus 両替ステータス通知
func (c *texMoneyHandler) SensorCashctlNoticeExchangeStatus(texCon *domain.TexContext, x interface{}) {
	c.logger.Trace("【%v】START:SensorCashctlNoticeExchangeStatus", texCon.GetUniqueKey())
	statusCash := c.texMoneyNoticeManager.GetStatusCashData(texCon)
	statusExchange := c.texMoneyNoticeManager.GetStatusExchangeData(texCon)
	switch receiveData := x.(type) {
	case domain.InStatus: //入金データ通知を変換する
		statusExchange.CashControlId = receiveData.CashControlId
		statusExchange.InOutType = domain.EXCHANGE_STATUS_IN
		statusExchange.StatusAction = true
		statusExchange.StatusResult = nil

		if receiveData.CoinStatusCode == 103 && receiveData.BillStatusCode == 103 { //入金確定
			statusExchange.StatusAction = false
			statusExchange.StatusResult = &domain.True
		}
		statusExchange.Amount = receiveData.Amount
		statusExchange.CountTbl = receiveData.CountTbl
		statusExchange.ExCountTbl = receiveData.ExCountTbl

		result, nowAmountTbl := c.safeInfoManager.GetSortInfo(texCon, 0)
		if result {
			for i, e := range receiveData.ExCountTbl { //現在有高に両替時の入金額をプラスする
				nowAmountTbl.ExCountTbl[i] += e
			}
		}
		var predictionExchangeTbl [26]int
		switch c.exchangePattern {
		case 0: //逆両替
			if c.exchangeTargetDevice == 1 { // targetDeviceが紙幣のみの場合
				predictionExchangeTbl = calculation.NewCassette(nowAmountTbl.ExCountTbl).Exchange(receiveData.Amount, 3)
			} else { // targetDeviceが紙幣のみ以外の場合
				predictionExchangeTbl = calculation.NewCassette(nowAmountTbl.ExCountTbl).Exchange(receiveData.Amount, 0)
			}
		case 1, 2: // 両替
			predictionExchangeTbl, _ = c.exchange(texCon, c.exchangePattern) // 出金予定枚数を算出
		}

		c.logger.Debug("【%v】- 両替金種別予定枚数 %+v", texCon.GetUniqueKey(), predictionExchangeTbl)
		var exchangeCountTbl [10]int
		for i := range exchangeCountTbl {
			exchangeCountTbl[i] = predictionExchangeTbl[i]
		}

		statusExchange.ExchangeCountTbl = exchangeCountTbl

		if receiveData.ErrorCode != "" {
			statusExchange.ErrorCode = receiveData.ErrorCode
			statusExchange.ErrorDetail = receiveData.ErrorDetail
			c.SetFlagExchange(texCon, false)
		}

		ok := c.texMoneyNoticeManager.UpdateStatusExchangeData(texCon, statusExchange)

		//金庫情報更新
		if (receiveData.CoinStatusCode == domain.IN_PAYMENT_ERROR && receiveData.BillStatusCode == domain.IN_PAYMENT_ERROR) ||
			(receiveData.CoinStatusCode == domain.IN_PAYMENT_COMPLETED && receiveData.BillStatusCode == domain.IN_PAYMENT_COMPLETED) {
			c.CalculationInfoSafeTbl(texCon, domain.TRADE)
		}

		/*No.345
		金庫情報更新完了前に次の要求を受信していた為、両替完了の通知タイミングを金庫情報更新後に変更
		c.SetTexmyNoticeExchangedata(texCon, ok)*/
		if receiveData.CoinStatusCode != domain.IN_PAYMENT_COMPLETED &&
			receiveData.BillStatusCode != domain.IN_PAYMENT_COMPLETED {
			c.SetTexmyNoticeExchangedata(texCon, ok)
		}

	case domain.OutStatus: //出金データ通知を変換する
		// 上位的には両替として1つのTopicなのに、2-1からは入金と出金で異なるIDが採番される為
		// 両替に関しては、出金の停止を提供していない(上位が知る必要のないID)ので、入金IDを一貫して利用するように調整する。 2024/1/15 No2070
		statusExchange.CashControlId = c.reqExchangeCashControlId
		if c.reqExchangeCashControlId == "" {
			statusExchange.CashControlId = receiveData.CashControlId
		}
		statusExchange.InOutType = domain.EXCHANGE_STATUS_OUT
		switch {
		case receiveData.CoinStatusCode == domain.OUT_DEPOSIT_START || //出金開始時
			receiveData.BillStatusCode == domain.OUT_DEPOSIT_START:
			if c.GetSequence(texCon) == domain.EXCHANGEING_CANCEL { //両替取消時は完了通知が来ないため回避措置
				statusExchange.StatusAction = false
				statusExchange.StatusResult = &domain.True
			} else {
				statusExchange.StatusAction = true
				statusExchange.StatusResult = nil
			}
		case receiveData.CoinStatusCode == domain.OUT_PAYMENT_COMPLETED && //出金完了時
			receiveData.BillStatusCode == domain.OUT_PAYMENT_COMPLETED:
			statusExchange.StatusAction = false
			statusExchange.StatusResult = &domain.True
		}
		statusExchange.Amount = receiveData.Amount
		statusExchange.CountTbl = receiveData.CountTbl
		statusExchange.ExCountTbl = receiveData.ExCountTbl
		statusExchange.ExchangeCountTbl = [domain.CASH_TYPE_SHITEI]int{} //TODO:両替金種別出金予定枚数
		if receiveData.ErrorCode != "" {
			statusExchange.ErrorCode = receiveData.ErrorCode
			statusExchange.ErrorDetail = receiveData.ErrorDetail
			c.SetFlagExchange(texCon, false)
		}

		ok := c.texMoneyNoticeManager.UpdateStatusExchangeData(texCon, statusExchange)
		//現在のデータをセット
		c.CalculationInfoSafeTbl(texCon, domain.OUTTRADE) //金庫情報更新

		/*No.345
		金庫情報更新完了前に次の要求を受信していた為、両替完了の通知タイミングを金庫情報更新後に変更
		c.SetTexmyNoticeExchangedata(texCon, ok)*/
		if receiveData.CoinStatusCode != domain.OUT_PAYMENT_COMPLETED &&
			receiveData.BillStatusCode != domain.OUT_PAYMENT_COMPLETED {
			c.SetTexmyNoticeExchangedata(texCon, ok)
		}

		//現金入出金機制御ステータス
		statusCash.CashControlId = receiveData.CashControlId //入出金管理IDのみ更新
		if receiveData.ErrorCode == "" {
			statusCash.StatusError = true
			statusCash.ErrorCode = ""
			statusCash.ErrorDetail = ""
		} else {
			statusCash.StatusError = false
			statusCash.ErrorCode = receiveData.ErrorCode
			statusCash.ErrorDetail = receiveData.ErrorDetail
		}
		c.SetTexmyNoticeStatus(texCon, statusCash) //ステータス通知

		if receiveData.CoinStatusCode == domain.OUT_PAYMENT_COMPLETED &&
			receiveData.BillStatusCode == domain.OUT_PAYMENT_COMPLETED {
			c.SetFlagExchange(texCon, false)
		}

	}

	c.logger.Trace("【%v】END:SensorCashctlNoticeExchangeStatus", texCon.GetUniqueKey())
}

// Exchange 両替 入金情報から両替内訳を算出する
func (c *texMoneyHandler) exchange(texCon *domain.TexContext, exchangePattern int) (outCountTbl [26]int, result bool) {
	statusIndata := c.texMoneyNoticeManager.GetStatusInData(texCon)

	// 最上位の入金金種を取得
	idx := targetMoneyType(statusIndata.CountTbl)
	if idx == -1 {
		c.logger.Debug("【%v】- 入金0円", texCon.GetUniqueKey())
		return outCountTbl, false
	}

	//両替パターンから両替内訳を取得
	exchangeCountTbl, exists := domain.GetExchangeData(domain.Cash[idx], exchangePattern)
	if !exists {
		// 両替パターンにない金種の場合はそのまま返金する
		c.logger.Error("【%v】- 両替対象外金種", texCon.GetUniqueKey())
		copy(outCountTbl[:], statusIndata.CountTbl[:])
		return outCountTbl, false
	}
	c.logger.Debug("【%v】- exchangeCountTbl=%v", texCon.GetUniqueKey(), exchangeCountTbl)

	// 現在の釣銭可能枚数を取得
	_, safeZero := c.safeInfoManager.GetSortInfo(texCon, 0)
	c.logger.Debug("【%v】両替時 金庫情報", texCon.GetUniqueKey())
	c.safeInfoManager.OutputLogSafeInfoExCountTbl(texCon)

	// あふれ枚数を除いた現在有高を10金種配列に変換
	var safe [10]int
	for i, v := range safeZero.ExCountTbl {
		if i < 10 { //普通金庫
			safe[i] = v
		} else if i >= 10 && i < 16 { //予備金庫
			safe[i-6] += v
		}
	}
	// 入金前の有高を保存 //入金確定時点で現在有高を更新しているため確定前のデータを保持しておく
	if statusIndata.StatusAction {
		c.beforeExchangeSafe = safe
	}
	c.logger.Debug("【%v】両替 入金前有高 c.beforeExchangeSafe=%v", texCon.GetUniqueKey(), c.beforeExchangeSafe)

	//入金前の現在有高から両替可能か判定
	ok := true
	for i, v := range c.beforeExchangeSafe {
		if v < exchangeCountTbl[i] {
			ok = false
		}
	}
	if !ok {
		// 現在有高から両替不可の場合は入金した金種でそのまま返金する
		c.logger.Debug("【%v】- 釣銭可能枚数から両替不可", texCon.GetUniqueKey())
		copy(outCountTbl[:], statusIndata.CountTbl[:])
		return outCountTbl, false
	}

	// 返金枚数算出
	refund := statusIndata.CountTbl
	refund[idx]--

	// 返金額を含めた出金枚数内訳を算出
	copy(outCountTbl[:], exchangeCountTbl[:])
	for i := range refund {
		outCountTbl[i] += refund[i]
	}

	return outCountTbl, result
}

// targetMoneyType 金種別枚数から最上位金種を判定
func targetMoneyType(countTbl [10]int) int {
	for i, v := range countTbl {
		if v > 0 {
			return i
		}
	}
	return -1
}

func (c *texMoneyHandler) SetExchangeTargetDevice(texCon *domain.TexContext, target int) {
	if target > 2 {
		c.exchangeTargetDevice = 0
		return
	}

	c.exchangeTargetDevice = target
}

func (c *texMoneyHandler) SetExchangePattern(texCon *domain.TexContext, pattern int) {
	if pattern > 2 {
		c.exchangePattern = 0
		return
	}

	c.exchangePattern = pattern
}

// CheckAmountLimit リミット有高チェック
// 注意の場合には精算機の停止は伴わない為、状態管理側と調整する事
func (c *texMoneyHandler) CheckAmountLimit(texCon *domain.TexContext, statusCash domain.StatusCash) domain.StatusCash {
	// チェック配列を返す関数
	currentStatusTbl := c.MakeCurrentStatusTbl(texCon)
	// チェック配列から、エラーコードを抽出する関数
	resultCode := NewCurrentStatus(&c.config.OverFlowBoxType).ErrorCheckCurrentStatusTbl(currentStatusTbl)
	// チェック配列から、紙幣と硬貨の有高枚数状態をセットする関数
	// リミット条件チェックで紙幣と硬貨の両方でエラーが発生した際にLEDを制御するのにそれぞれの状態が必要なため
	statusCash = c.setStatusAmountCount(currentStatusTbl, statusCash)

	if resultCode != 0 {
		c.logger.Debug("【%v】texMoneyHandler CheckAmountLimit resultCode=%v", texCon.GetUniqueKey(), resultCode)
		statusCash.StatusError = false
		statusCash.ErrorCode, statusCash.ErrorDetail = c.errorMng.GetErrorInfo(resultCode)
		return statusCash
	}
	statusCash.StatusError = true
	statusCash.ErrorCode = ""
	statusCash.ErrorDetail = ""
	return statusCash
}

// MakeCurrentStatusTbl 有高を、あふれ枚数、不足枚数と比較して通常金種別状況を作成する。
// 優先順位は互換性の為、あふれエラー、不足エラー、あふれ注意、不足注意の順でセットする。
func (c *texMoneyHandler) MakeCurrentStatusTbl(texCon *domain.TexContext) [domain.CASH_TYPE_SHITEI]int {
	c.logger.Trace("【%v】START:あふれ・不足枚数判定", texCon.GetUniqueKey())
	moneySetting := c.config.MoneySetting // 各種設定
	statusAmount := c.texMoneyNoticeManager.GetStatusAmountData(texCon)
	moneyList := domain.NewMoneyList(statusAmount.ExCountTbl) // 有高
	checkTbl := NewCurrentStatus(&c.config.OverFlowBoxType).MakeCurrentStatus(&moneySetting, moneyList)

	c.logger.Trace("【%v】END:あふれ・不足枚数判定 checkTbl=%v", texCon.GetUniqueKey(), *checkTbl)

	return *checkTbl
}

// 紙幣と硬貨の有高枚数状態をセットする
func (c *texMoneyHandler) setStatusAmountCount(currentStatusTbl [domain.CASH_TYPE_SHITEI]int, statusCash domain.StatusCash) domain.StatusCash {

	// 紙幣の有高枚数状態をセットする処理
	// billStatus 0:正常 1:エラー 2:警告
	var billStatus int
	for _, status := range currentStatusTbl[0:4] {
		switch status {
		case 2, 4: //不足エラー //オーバーエラー
			billStatus = 1
		case 1, 3: //不足警告 //オーバー警告
			if 0 == billStatus {
				billStatus = 2
			}
		default:
		}
	}
	statusCash.BillStatusTbl.StatusAmountCount = billStatus

	// 硬貨の有高枚数状態をセットする処理
	// coinStatus 0:正常 1:エラー 2:警告
	var coinStatus int
	for _, status := range currentStatusTbl[4:] {
		switch status {
		case 2, 4: //不足エラー //オーバーエラー
			coinStatus = 1
		case 1, 3: //不足警告 //オーバー警告
			if 0 == coinStatus {
				coinStatus = 2
			}
		default:
		}
	}
	statusCash.CoinStatusTbl.StatusAmountCount = coinStatus
	return statusCash
}

// 現金入出金機状況
func (c *texMoneyHandler) SensorCashctlNoticeStatus(texCon *domain.TexContext, stuInfo domain.NoticeStatus) {
	c.logger.Trace("【%v】START:texMoneyHandler SensorCashctlNoticeStatus stuInfo=%+v", texCon.GetUniqueKey(), stuInfo)
	statusCash := c.texMoneyNoticeManager.GetStatusCashData(texCon)
	//ステータスエラー検知
	statusCash.StatusError = stuInfo.CoinNoticeStatusTbl.StatusError && stuInfo.BillNoticeStatusTbl.StatusError
	statusCash.ErrorCode = stuInfo.CoinNoticeStatusTbl.ErrorCode
	statusCash.ErrorDetail = stuInfo.CoinNoticeStatusTbl.ErrorDetail
	if !stuInfo.BillNoticeStatusTbl.StatusError {
		// 紙幣エラーが発生している場合はエラー情報を上書きする
		statusCash.ErrorCode = stuInfo.BillNoticeStatusTbl.ErrorCode
		statusCash.ErrorDetail = stuInfo.BillNoticeStatusTbl.ErrorDetail
	}

	//トビラ状態
	// 両方true：閉じる状態の場合、trueにするように修正
	// 下位サービスは、紙幣・硬貨ごとに存在するが1-3からの通知は1種類しか存在しない為
	// どちらかが開いていれば、開いてる判定とする
	statusCash.StatusCover = stuInfo.CoinNoticeStatusTbl.StatusCover && stuInfo.BillNoticeStatusTbl.StatusCover

	//入出金口
	statusCash.StatusInsert = stuInfo.CoinNoticeStatusTbl.StatusInsert || stuInfo.BillNoticeStatusTbl.StatusInsert

	//出金口
	statusCash.StatusExit = stuInfo.CoinNoticeStatusTbl.StatusExit || stuInfo.BillNoticeStatusTbl.StatusExit

	//リジェクトボックス
	statusCash.StatusRjbox = stuInfo.CoinNoticeStatusTbl.StatusRjbox || stuInfo.BillNoticeStatusTbl.StatusRjbox

	// 紙幣ステータス情報
	statusCash.BillStatusTbl.StatusUnitSet = stuInfo.BillNoticeStatusTbl.StatusUnitSet
	statusCash.BillStatusTbl.StatusInCassette = stuInfo.BillNoticeStatusTbl.StatusInCassette
	statusCash.BillStatusTbl.StatusOutCassette = stuInfo.BillNoticeStatusTbl.StatusOutCassette

	// 硬貨ステータス情報
	statusCash.CoinStatusTbl.StatusUnitSet = stuInfo.CoinNoticeStatusTbl.StatusUnitSet
	statusCash.CoinStatusTbl.StatusInCassette = stuInfo.CoinNoticeStatusTbl.StatusInCassette
	statusCash.CoinStatusTbl.StatusOutCassette = stuInfo.CoinNoticeStatusTbl.StatusOutCassette

	//デバイス詳細情報
	statusCash.DeviceStatusInfoTbl = make([]string, 6)

	// 要素数の最小値を計算
	minLen := len(statusCash.DeviceStatusInfoTbl)
	if len(stuInfo.DeviceStatusInfoTbl) < minLen {
		minLen = len(stuInfo.DeviceStatusInfoTbl)
	}
	// コピーを実行
	copy(statusCash.DeviceStatusInfoTbl[:minLen], stuInfo.DeviceStatusInfoTbl[:minLen])

	//警告情報
	statusCash.WarningInfoTbl = make([]int, 10)

	// 要素数の最小値を計算
	minLen = len(statusCash.WarningInfoTbl)
	if len(stuInfo.WarningInfoTbl) < minLen {
		minLen = len(stuInfo.WarningInfoTbl)
	}
	// 要素をコピー
	copy(statusCash.WarningInfoTbl[:minLen], stuInfo.WarningInfoTbl[:minLen])

	c.logger.Debug("【%v】- StatusCover %t", texCon.GetUniqueKey(), statusCash.StatusCover)
	c.SetTexmyNoticeStatus(texCon, statusCash)
	c.logger.Trace("【%v】END:texMoneyHandler SensorCashctlNoticeStatus", texCon.GetUniqueKey())
}

// 有高枚数変更要求実施前のデータ一時保存
func (c *texMoneyHandler) TempAmountData(texCon *domain.TexContext) {
	// statusAmount := c.texMoneyNoticeManager.GetStatusAmountData(texCon)

	logicalCash := c.safeInfoManager.GetLogicalCashAvailable(texCon)

	c.logger.Trace("【%v】START:texMoneyHandler TempAmountData savedataStatusAmountData=%+v", texCon.GetUniqueKey(), logicalCash)
	c.tempAmountData.Amount = logicalCash.Amount
	c.tempAmountData.CountTbl = logicalCash.CountTbl
	c.tempAmountData.ExCountTbl = logicalCash.ExCountTbl
	c.logger.Trace("【%v】END:texMoneyHandler TempAmountData, tempAmountData=%+v", texCon.GetUniqueKey(), c.tempAmountData)
}

func (c *texMoneyHandler) GetTempAmountData(texCon *domain.TexContext) (int, [10]int, [26]int) {
	return c.tempAmountData.Amount, c.tempAmountData.CountTbl, c.tempAmountData.ExCountTbl
}

// 現金入出金機制御の応答
func (c *texMoneyHandler) RecvCashctlALLRequest(texCon *domain.TexContext, x interface{}) {
	c.logger.Trace("【%v】START:texMoneyHandler RecvCashctlALLRequest", texCon.GetUniqueKey())
	c.logger.Debug("【%v】- recvInfo %+v", texCon.GetUniqueKey(), x)
	statusCash := c.texMoneyNoticeManager.GetStatusCashData(texCon)

	// エラー情報のセット
	setErrorCode := func(result bool, errorCode, errorDetail string) {
		statusCash.StatusError = result
		if !result {
			statusCash.ErrorCode = errorCode
			statusCash.ErrorDetail = errorDetail
		} else {
			statusCash.ErrorCode = ""
			statusCash.ErrorDetail = ""
		}
	}

	//型チェック
	switch resInfo := x.(type) {
	case domain.ResultOutStart: //出金開始要求
		statusCash.CashControlId = resInfo.CashControlId //キャッシュIDの登録
		setErrorCode(resInfo.Result, resInfo.ErrorCode, resInfo.ErrorDetail)
	case domain.ResultInStart: //入金開始要求
		statusCash.CashControlId = resInfo.CashControlId //キャッシュIDの登録
		setErrorCode(resInfo.Result, resInfo.ErrorCode, resInfo.ErrorDetail)
	case domain.ResultInEnd: //入金終了要求
		statusCash.CashControlId = resInfo.CashControlId //キャッシュIDの登録
		setErrorCode(resInfo.Result, resInfo.ErrorCode, resInfo.ErrorDetail)
	case domain.ResultCashctlSetAmount: //現金入出金機制御:有高枚数変更要求
		statusCash.CashControlId = resInfo.CashControlId //キャッシュIDの登録
		setErrorCode(resInfo.Result, resInfo.ErrorCode, resInfo.ErrorDetail)
	case domain.ResultCollectStart: //回収開始
		statusCash.CashControlId = resInfo.CashControlId //キャッシュIDの登録
		setErrorCode(resInfo.Result, resInfo.ErrorCode, resInfo.ErrorDetail)
	case domain.ResultCollectStop: //回収停止
		setErrorCode(resInfo.Result, resInfo.ErrorCode, resInfo.ErrorDetail)
	}

	c.SetTexmyNoticeStatus(texCon, statusCash) //ステータスの更新
	c.logger.Trace("【%v】END:texMoneyHandler RecvCashctlALLRequest", texCon.GetUniqueKey())
}

// 有高枚数変更要求の有高ステータスから補充と出金の枚数を通知
func (c *texMoneyHandler) RecvSetAmountNoticeAmountStatus(texCon *domain.TexContext, amStatus domain.AmountStatus, onlyAmountChangeSequence int) bool {
	c.logger.Trace("【%v】START:有高からの入金出金枚数算出 onlyAmountChangeSequence=%v", texCon.GetUniqueKey(), onlyAmountChangeSequence)
	c.logger.Debug("【%v】- tempAmountData=%v", texCon.GetUniqueKey(), c.tempAmountData)

	// 釣り銭不一致判定
	baseAmount := amStatus
	logcalAmount := c.safeInfoManager.GetLogicalCashAvailable(texCon)
	if logcalAmount.Amount != amStatus.Amount {
		baseAmount.Amount = logcalAmount.Amount
		baseAmount.CountTbl = logcalAmount.CountTbl
		baseAmount.ExCountTbl = logcalAmount.ExCountTbl

	}

	var stuInfo domain.InStatus      //入金ステータス
	var stOutStatus domain.OutStatus //出金ステータス
	//通常金種別枚数
	for i := 0; i < len(baseAmount.CountTbl); i++ {
		if baseAmount.CountTbl[i] < c.tempAmountData.CountTbl[i] { //過去の有高枚数の方が大きいとき：回収→出金ステータスへ
			stOutStatus.CountTbl[i] = c.tempAmountData.CountTbl[i] - baseAmount.CountTbl[i]
		} else if baseAmount.CountTbl[i] > c.tempAmountData.CountTbl[i] { //過去の有高枚数の方が小さいとき：補充→入金ステータスへ
			stuInfo.CountTbl[i] = baseAmount.CountTbl[i] - c.tempAmountData.CountTbl[i]
		}
	}
	//拡張金種別枚数
	for j := 0; j < len(baseAmount.ExCountTbl); j++ {
		if baseAmount.ExCountTbl[j] < c.tempAmountData.ExCountTbl[j] { //過去の有高枚数の方が大きいとき：回収→出金ステータスへ
			stOutStatus.ExCountTbl[j] = c.tempAmountData.ExCountTbl[j] - baseAmount.ExCountTbl[j]
		} else if baseAmount.ExCountTbl[j] > c.tempAmountData.ExCountTbl[j] { //過去の有高枚数の方が小さいとき：補充→入金ステータスへ
			stuInfo.ExCountTbl[j] = baseAmount.ExCountTbl[j] - c.tempAmountData.ExCountTbl[j]
		}
	}

	//Amountの設定
	for k := 0; k < len(stOutStatus.CountTbl); k++ {
		stuInfo.Amount = stuInfo.Amount + (stuInfo.CountTbl[k] * domain.Cash[k])
		stOutStatus.Amount = stOutStatus.Amount + (stOutStatus.CountTbl[k] * domain.Cash[k])
	}

	stuInfo.CashControlId = domain.SET_AMOUNT_INDATA_ONE
	stuInfo.BillStatusCode = domain.IN_PAYMENT_COMPLETED
	stuInfo.CoinStatusCode = domain.IN_PAYMENT_COMPLETED
	stuInfo.ErrorCode = baseAmount.ErrorCode
	stuInfo.ErrorDetail = baseAmount.ErrorDetail

	stOutStatus.CashControlId = domain.SET_AMOUNT_OUTDATA_ONE
	stOutStatus.BillStatusCode = domain.OUT_PAYMENT_COMPLETED
	stOutStatus.CoinStatusCode = domain.OUT_PAYMENT_COMPLETED
	stOutStatus.ErrorCode = baseAmount.ErrorCode
	stOutStatus.ErrorDetail = baseAmount.ErrorDetail

	//入金ステータス通知
	c.SensorCashctlNoticeInStatus(texCon, stuInfo)
	//出金ステータス通知
	c.SensorCashctlNoticeOutStatus(texCon, stOutStatus)

	statusOut := c.texMoneyNoticeManager.GetStatusOutData(texCon)
	statusOut.Amount = stOutStatus.Amount
	statusOut.CountTbl = stOutStatus.CountTbl
	statusOut.ExCountTbl = stOutStatus.ExCountTbl
	statusOut.ErrorCode = amStatus.ErrorCode
	statusOut.ErrorDetail = amStatus.ErrorDetail
	c.texMoneyNoticeManager.UpdateStatusOutData(texCon, statusOut)

	//有高ステータス通知
	c.SensorCashctlNoticeAmountStatus(texCon, amStatus)

	// 非還流庫回収and売上金回収の場合のみ
	if c.tempCollectAndSales.specialSequence == domain.UNRETURNED_AND_SALES_COLLECT {
		c.logger.Debug("【%v】- 売上金更新実行", texCon.GetUniqueKey())
		salesCountTbl := calculation.NewCassette(c.tempCollectAndSales.salesExCountTbl).ExCountTblToTenCountTbl()
		c.SetCollectSales(texCon, c.tempCollectAndSales.salesAmount, salesCountTbl, c.tempCollectAndSales.salesExCountTbl)
		c.tempCollectAndSales = tempCollectAndSales{}
	}

	c.logger.Trace("【%v】END:有高からの入金出金枚数算出", texCon.GetUniqueKey())
	return true
}

// 稼働データ管理
// イニシャル時の稼働データの情報を保管
func (c *texMoneyHandler) TexdtInfoSave(texCon *domain.TexContext, resultStatus domain.ResultGetTermInfoNow) {
	c.logger.Trace("【%v】START:初回起動時DB取得取得保存 金庫情報=%v", texCon.GetUniqueKey(), resultStatus)

	// 取得データがあればセット
	for _, getTemNow := range resultStatus.InfoSafeTblGetTermNow.SortInfoTbl {
		c.safeInfoManager.UpdateSortInfo(texCon, getTemNow)
		// 有高の更新の場合には、論理枚数にもセットする
		if getTemNow.SortType == domain.CASH_AVAILABLE {
			c.safeInfoManager.UpdateAllLogicalCashAvailable(texCon, getTemNow)
			// tempAmontDataを一旦ここで更新しておく
			c.tempAmountData.Amount = getTemNow.Amount
			c.tempAmountData.CountTbl = getTemNow.CountTbl
			c.tempAmountData.ExCountTbl = getTemNow.ExCountTbl

		}
	}
	c.safeInfoManager.OutputLogSafeInfoExCountTbl(texCon)

	c.logger.Trace("【%v】END:初回起動時DB取得取得保存", texCon.GetUniqueKey())
}

// 金庫情報更新
func (c *texMoneyHandler) CalculationInfoSafeTbl(texCon *domain.TexContext, amountOrTrade int) {
	c.logger.Trace("【%v】texMoneyHandler CalculationInfoSafeTbl amountOrTrade=%v", texCon.GetUniqueKey(), amountOrTrade)
	defer func() {
		c.safeInfoManager.UpdateBalanceInfo(texCon) //差引情報更新
		c.updateReportData(texCon)                  // レポート用データ更新

		c.logger.Trace("【%v】END:texMoneyHandler CalculationInfoSafeTbl", texCon.GetUniqueKey())
	}()

	if amountOrTrade == domain.AMOUNT { //有高通知の時
		c.SetAmount(texCon)        //現金有高
		c.SetReserveCharge(texCon) //釣銭準備金
		return
	} else if amountOrTrade == domain.TRADE { //取引関連の動きの時
		c.updateInData(texCon)
		return
	} else if amountOrTrade == domain.OUTTRADE { //出金関連の動きの時
		c.updateOutData(texCon)
		return
	}
}

// 金庫情報更新(取引関連)
func (c *texMoneyHandler) updateInData(texCon *domain.TexContext) {
	statusIndata := c.texMoneyNoticeManager.GetStatusInData(texCon)

	switch c.GetSequence(texCon) {
	case domain.INITIAL_ADDING_CONFIRM, domain.INITIAL_ADDING_START: //初期補充の時
		//初期補充
		c.safeInfoManager.UpdateSortInfoCumulative(texCon, domain.INITIAL_REPLENISHMENT, statusIndata.Amount, statusIndata.CountTbl, statusIndata.ExCountTbl)

		//補充入金
		c.safeInfoManager.UpdateSortInfoCumulative(texCon, domain.REPLENISHMENT_DEPOSIT, statusIndata.Amount, statusIndata.CountTbl, statusIndata.ExCountTbl)
	case domain.SET_AMOUNT, //有高枚数変更中
		domain.MONEY_ADD_REPLENISH_START,   //追加補充開始 // TODO: 追加補充開始の場合は金庫情報を更新しないはずなので未使用の場合は削除する
		domain.MONEY_ADD_REPLENISH_CONFIRM: //追加補充確定

		//補充入金
		c.safeInfoManager.UpdateSortInfoCumulative(texCon, domain.REPLENISHMENT_DEPOSIT, statusIndata.Amount, statusIndata.CountTbl, statusIndata.ExCountTbl)
	case domain.REVERSE_EXCHANGEING_CONFIRM_INDATA: //逆両替入金データ時
		statusExchange := c.texMoneyNoticeManager.GetStatusExchangeData(texCon)

		// 両替は取引データを更新、それ以外は補充データを更新
		if c.exchangePattern == 1 || c.exchangePattern == 2 { // 両替
			//  取引入金
			c.safeInfoManager.UpdateSortInfoCumulative(texCon, domain.TRANSACTION_DEPOSIT, statusExchange.Amount, statusExchange.CountTbl, statusExchange.ExCountTbl)
		} else {
			// 補充入金
			c.safeInfoManager.UpdateSortInfoCumulative(texCon, domain.REPLENISHMENT_DEPOSIT, statusExchange.Amount, statusExchange.CountTbl, statusExchange.ExCountTbl)
		}

	case domain.TRANSACTION_DEPOSIT_CONFIRM, domain.TRANSACTION_DEPOSIT_END_BILL, domain.TRANSACTION_DEPOSIT_END_COIN: // 取引入金確定,取引終了紙幣,取引終了硬貨
		c.safeInfoManager.UpdateSortInfoCumulative(texCon, domain.TRANSACTION_DEPOSIT, statusIndata.Amount, statusIndata.CountTbl, statusIndata.ExCountTbl)

	case domain.MANUAL_REPLENISHMENT_COLLECTION:
		c.safeInfoManager.UpdateSortInfoCumulativeNoUpdateLogicalCash(texCon, domain.REPLENISHMENT_DEPOSIT, statusIndata.Amount, statusIndata.CountTbl, statusIndata.ExCountTbl)

	default:
		c.logger.Debug("【%v】- 対象無し", texCon.GetUniqueKey())
	}
}

// 金庫情報更新(出金関連)
func (c *texMoneyHandler) updateOutData(texCon *domain.TexContext) {
	statusOut := c.texMoneyNoticeManager.GetStatusOutData(texCon)
	switch c.GetSequence(texCon) {
	case domain.INITIAL_ADDING_CONFIRM,
		domain.INITIAL_ADDING_START,
		domain.SET_AMOUNT,
		domain.REJECTBOXCOLLECT_START,
		domain.UNRETURNEDCOLLECT_START,
		domain.MIDDLE_START_OUT_START,
		domain.MIDDLE_START_OUT_STOP,
		domain.MIDDLE_START_COLLECT_START,
		domain.MIDDLE_START_COLLECT_STOP,
		domain.ALLCOLLECT_START_OUT_START,
		domain.ALLCOLLECT_START_OUT_STOP,
		domain.ALLCOLLECT_START_COLLECT_START,
		domain.ALLCOLLECT_START_COLLECT_STOP,
		domain.MANUAL_REPLENISHMENT_COLLECTION: // TODO 仮追加
		//補充出金
		c.SetReplenishmentWithdrawal(texCon, statusOut.Amount, statusOut.CountTbl, statusOut.ExCountTbl)
	case domain.REVERSE_EXCHANGEING_CONFIRM_OUTDATA: //逆両替出金データ時
		statusExchange := c.texMoneyNoticeManager.GetStatusExchangeData(texCon)

		// 両替は取引データを更新、それ以外は補充データを更新
		if c.exchangePattern == 1 || c.exchangePattern == 2 { // 両替
			//  取引出金
			c.safeInfoManager.UpdateSortInfoCumulative(texCon, domain.TRANSACTION_WITHDRAWAL, statusExchange.Amount, statusExchange.CountTbl, statusExchange.ExCountTbl)
		} else {
			// 補充出金
			c.SetReplenishmentWithdrawal(texCon, statusExchange.Amount, statusExchange.CountTbl, statusExchange.ExCountTbl)
		}

	case domain.SALESMONEY_START: //売上金回収開始の時
		//補充出金
		c.SetReplenishmentWithdrawal(texCon, statusOut.Amount, statusOut.CountTbl, statusOut.ExCountTbl)
		// 出金完了時のみ売上金回収情報を更新する
		if !statusOut.StatusAction {
			//売上金回収
			c.SetCollectSales(texCon, statusOut.Amount, statusOut.CountTbl, statusOut.ExCountTbl)
		}
	case domain.TRANSACTION_OUT_START, //取引出金確定
		domain.TRANSACTION_OUT_CANCEL: //取引出金停止

		//取引出金
		c.safeInfoManager.UpdateSortInfoCumulative(texCon, domain.TRANSACTION_WITHDRAWAL, statusOut.Amount, statusOut.CountTbl, statusOut.ExCountTbl)
	case domain.TRANSACTION_DEPOSIT_CANCEL: //取引入金取消
		//取引出金
		c.safeInfoManager.UpdateSortInfoCumulative(texCon, domain.TRANSACTION_WITHDRAWAL, statusOut.Amount, statusOut.CountTbl, statusOut.ExCountTbl)
		// 論理有高更新 取引入金取消の場合、下位からの通知でマイナスする必要がある為
		c.safeInfoManager.UpdateOutLogicalCashAvailable(texCon, domain.SortInfoTbl{
			SortType:   domain.CASH_AVAILABLE,
			Amount:     statusOut.Amount,
			CountTbl:   statusOut.CountTbl,
			ExCountTbl: statusOut.ExCountTbl,
		})

	default:
		c.logger.Debug("【%v】- 対象無し", texCon.GetUniqueKey())
	}
}

// レポート用データ更新
func (c *texMoneyHandler) updateReportData(texCon *domain.TexContext) {
	status := c.maintenanceModeMng.GetMaintenanceMode(texCon)
	mode := c.maintenanceModeMng.GetMode(texCon)
	if mode == 0 || status == domain.REPLENISHMENT_END || status == domain.CLOSING_END {
		return
	}

	_, aggSafeInfo := c.aggregateManager.GetAggregateSafeInfo(texCon, mode)
	// レポート用補充入金金種配列に「処理前補充入金 - 補充入金」を格納
	_, replenishmentDeposit := c.safeInfoManager.GetSortInfo(texCon, domain.REPLENISHMENT_DEPOSIT) // 現在の補充入金情報取得
	replenishCountTbl := c.aggregateManager.DiffTbl(texCon, aggSafeInfo.BeforeReplenishCountTbl, replenishmentDeposit.ExCountTbl)
	c.aggregateManager.UpdateAggregateCountTbl(texCon, mode, domain.REPLENISH_COUNT_TBL, replenishCountTbl)
	// レポート用回収金種配列に「処理前補充出金 - 補充出金」を格納
	_, replenishmentWithdrawal := c.safeInfoManager.GetSortInfo(texCon, domain.REPLENISHMENT_WITHDRAWAL) // 現在の補充入金情報取得
	collectCountTbl := c.aggregateManager.DiffTbl(texCon, aggSafeInfo.BeforeCollectCountTbl, replenishmentWithdrawal.ExCountTbl)
	c.aggregateManager.UpdateAggregateCountTbl(texCon, mode, domain.COLLECT_COUNT_TBL, collectCountTbl)
}

// 金庫情報分類情報種別
// 現金有高更新
func (c *texMoneyHandler) SetAmount(texCon *domain.TexContext) {
	c.logger.Trace("【%v】START:texMoneyHandler SetAmount", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:texMoneyHandler SetAmount", texCon.GetUniqueKey())
	// 有高通知取得
	statusAmount := c.texMoneyNoticeManager.GetStatusAmountData(texCon)
	// 現金有高取得
	_, availableSortInfo := c.safeInfoManager.GetSortInfo(texCon, domain.CASH_AVAILABLE)
	// 有高情報の更新
	availableSortInfo.Amount = statusAmount.Amount
	availableSortInfo.CountTbl = statusAmount.CountTbl
	availableSortInfo.ExCountTbl = statusAmount.ExCountTbl
	// 現金有高更新
	if !c.getNoticeAmountUpdateNg(texCon) && !LogicalChange {
		c.safeInfoManager.UpdateSortInfo(texCon, availableSortInfo)
	}
	// デバイス有高更新

	// 入金時に,確定前だと還流庫カウントされ、確定時に非還流庫カウントとなるロジックの影響で
	// 入金時に確定の場合に、更新するようにしている。
	if !c.getNoticeAmountUpdateNg(texCon) {
		c.safeInfoManager.UpdateDeviceCashAvailable(texCon, availableSortInfo)
	} else {
		// 次回に向けてfalseをセットしておく
		c.setNoticeAmountUpdateNg(texCon, false)
	}

	//論理枚数とデバイス枚数の比較を実施する。
	//整合性不一致の場合には、エラー処理
	// ここの関数を通るときが、デバイス有高確定の為
	// ここでフラグを立てておき、送信前に再チェックを行って釣銭不一致を発生させる。
	// 上位への通知は以下関数で実行され、フラグはfalseに戻る。
	// (c *texMoneyHandler) SetTexmyNoticeStatus(texCon *domain.TexContext, statusCash domain.StatusCash)
	ng, _ := c.safeInfoManager.GetAvailableBalance(texCon)
	status := c.maintenanceModeMng.GetMaintenanceMode(texCon)
	// 差分有&初回起動時（メンテナンスモード無視）&監視タイミングの場合 又は、差分有&保守操作中外&監視タイミングの場合
	if (ng && !c.initialDiscrepanctStartOne && c.initialDiscrepanctOn) || (ng && status != 1 && status != 3 && c.initialDiscrepanctOn) {
		c.logger.Error("【%v】- 不一致エラー検知", texCon.GetUniqueKey())
		c.cashDiscrepanctFlg = true
	}

}

func (c *texMoneyHandler) getNoticeAmountUpdateNg(texCon *domain.TexContext) bool {

	var s string
	s = "取得 有高更新更新 有"
	if c.noticeAmountUpdateNg {
		s = "取得 有高更新更新 無"
	}

	c.logger.Debug("【%v】%v", texCon.GetUniqueKey(), s)
	return c.noticeAmountUpdateNg
}

func (c *texMoneyHandler) setNoticeAmountUpdateNg(texCon *domain.TexContext, b bool) {
	c.logger.Debug("【%v】setNoticeAmountUpdateNg set=%v", texCon.GetUniqueKey(), b)
	c.noticeAmountUpdateNg = b
}

/*
釣銭不一致 shouldCheckChange bool：true(チェックする)　false(チェックしない)
起動時に監視したいタイミングまで進んだらtrueへ変更する
*/
func (c *texMoneyHandler) InitialDiscrepantOn(shouldChangeCheck bool) {
	c.logger.Debug("不一致チェック開始(以降の処理で有高通知を受信した場合不一致チェックが実行される) shouldChangeCheck=%t", shouldChangeCheck)
	if shouldChangeCheck {
		c.initialDiscrepanctOn = true
	} else {
		c.initialDiscrepanctOn = false
	}
}

// 初回起動時だけメンテナンスモードを無視して判定処理を実施する為の設定
func (c *texMoneyHandler) InitialDiscrepanctStartOne() {
	c.logger.Debug("初回起動時 不一致判定終了(メンテナンスモード無視での不一致チェック)")
	c.initialDiscrepanctOn = true
}

func (c *texMoneyHandler) cashDiscrepancyError(texCon *domain.TexContext) {
	c.logger.Debug("【%v】texMoneyHandler cashDiscrepancyError", texCon.GetUniqueKey())
	status := c.maintenanceModeMng.GetMaintenanceMode(texCon)
	if (status == 1 || status == 3) && c.initialDiscrepanctStartOne {
		c.logger.Debug("【%v】- 補充・締め処理中は釣銭不一致対象外", texCon.GetUniqueKey())
		return
	}

	ng, balance := c.safeInfoManager.GetAvailableBalance(texCon)

	if ng {
		availableSortInfo := c.safeInfoManager.GetDeviceCashAvailable(texCon)
		logicalSortInfo := c.safeInfoManager.GetLogicalCashAvailable(texCon)

		var l string
		l += fmt.Sprintf("\n【%v】- 釣銭不一致発生\n", texCon.GetUniqueKey())
		l += fmt.Sprintf("【%v】- 有高不一致＿:%+v\n", texCon.GetUniqueKey(), balance)
		l += fmt.Sprintf("【%v】- 論理有高＿＿:%+v\n", texCon.GetUniqueKey(), logicalSortInfo)
		l += fmt.Sprintf("【%v】- デバイス有高:%+v\n", texCon.GetUniqueKey(), availableSortInfo)
		c.logger.Error("%v", l)

		// 釣銭不一致金種情報
		errorDetail := "不一致金種["

		eTbl := []string{
			" 1万円 ",
			" 5千円",
			" 2千円 ",
			" 千円 ",
			" 500円 ",
			" 100円 ",
			" 50円 ",
			" 10円 ",
			" 5円 ",
			" 1円 ",
		}

		for i, v := range balance.CountTbl {
			if v != 0 {
				errorDetail += eTbl[i]
			}
		}
		errorDetail += "]"

		// 今のステータス情報取得
		s := c.texMoneyNoticeManager.GetStatusCashData(texCon)

		// エラーコード変換
		e, _ := c.errorMng.GetErrorInfo(ERROR_CASH_DISCREPANCY)

		s.StatusError = false
		s.ErrorCode = e
		s.ErrorDetail = errorDetail

		c.logger.Error("【%v】- 送信予定不一致エラー %+v", texCon.GetUniqueKey(), s)

		// 送信予定ステータスの更新
		if LogicalChange {
			c.texMoneyNoticeManager.UpdateStatusCashData(texCon, s)
		}
	}

}

// 釣銭準備金更新
func (c *texMoneyHandler) SetReserveCharge(texCon *domain.TexContext) {

	_, changeAvailableSortInfo := c.safeInfoManager.GetSortInfo(texCon, domain.CHANGE_AVAILABLE)
	statusAmount := c.texMoneyNoticeManager.GetStatusAmountData(texCon)

	c.logger.Trace("【%v】START:texMoneyHandler SetReserveCharge 釣銭可能枚数=%+v", texCon.GetUniqueKey(), changeAvailableSortInfo)

	//釣銭可能枚数
	// 釣銭可能枚数金額算出
	amount, exCountTbl := calculation.NewCassette(statusAmount.ExCountTbl).GetChangeAvailable()
	countTbl := calculation.NewCassette(exCountTbl).ExCountTblToTenCountTbl()

	changeAvailableSortInfo.Amount = amount
	changeAvailableSortInfo.CountTbl = countTbl
	changeAvailableSortInfo.ExCountTbl = exCountTbl

	c.safeInfoManager.UpdateSortInfo(texCon, changeAvailableSortInfo)

	c.logger.Debug("【%v】- 釣銭可能枚数=%+v", texCon.GetUniqueKey(), changeAvailableSortInfo)
	c.logger.Trace("【%v】END:texMoneyHandler SetReserveCharge", texCon.GetUniqueKey())
}

// 補充出金更新
func (c *texMoneyHandler) SetReplenishmentWithdrawal(texCon *domain.TexContext, amount int, countTbl [domain.CASH_TYPE_SHITEI]int, exCountTbl [domain.EXTRA_CASH_TYPE_SHITEI]int) {
	c.logger.Trace("【%v】START:texMoneyHandler SetReplenishmentWithdrawal", texCon.GetUniqueKey())
	// TODO: 仮 確定前には足し算を行わない。 No9051
	// 出金中と出金完了の通知が発生した場合、2重で出金情報がたされてしまうパターンが存在する為
	statusOut := c.texMoneyNoticeManager.GetStatusOutData(texCon)
	if statusOut.StatusAction {
		c.logger.Trace("【%v】END:texMoneyHandler SetReplenishmentWithdrawal StatusAction Not Result 補充出金庫更新なし", texCon.GetUniqueKey())
		return
	}

	c.safeInfoManager.UpdateSortInfoCumulative(texCon, domain.REPLENISHMENT_WITHDRAWAL, amount, countTbl, exCountTbl)
	// 回収操作回数を更新
	c.safeInfoManager.CountUpCollectCount(texCon)

	if c.tempCollectAndSales.specialSequence == domain.MIDDLE_AND_SALES_COLLECT {
		c.logger.Debug("【%v】- 売上金更新実行", texCon.GetUniqueKey())
		salesCountTbl := calculation.NewCassette(c.tempCollectAndSales.salesExCountTbl).ExCountTblToTenCountTbl()
		// 売上金回収で使っているメソッドに横流しする
		c.SetCollectSales(texCon, c.tempCollectAndSales.salesAmount, salesCountTbl, c.tempCollectAndSales.salesExCountTbl)
		c.tempCollectAndSales = tempCollectAndSales{}
	}

	c.logger.Trace("【%v】END:texMoneyHandler SetReplenishmentWithdrawal", texCon.GetUniqueKey())
}

// 売上金回収更新
func (c *texMoneyHandler) SetCollectSales(texCon *domain.TexContext, amount int, countTbl [domain.CASH_TYPE_SHITEI]int, exCountTbl [domain.EXTRA_CASH_TYPE_SHITEI]int) {
	c.logger.Trace("【%v】START:texMoneyHandler SetCollectSales", texCon.GetUniqueKey())

	c.safeInfoManager.UpdateSortInfoCumulative(texCon, domain.SALES_MONEY_COLLECT, amount, countTbl, exCountTbl)

	// 売上金回収済み金額の更新、売上金回収回数の更新
	_, salesSortInfo := c.safeInfoManager.GetSortInfo(texCon, domain.SALES_MONEY_COLLECT)
	c.safeInfoManager.UpdateSalesCompleteAmount(texCon, salesSortInfo.Amount)
	c.safeInfoManager.CountUpSalesCompleteCount(texCon)

	// レポート用売上金回収配列データ更新
	c.aggregateManager.UpdateAggregateCountTbl(texCon, c.maintenanceModeMng.GetMode(texCon), domain.SALES_COLLECT_TBL, salesSortInfo.ExCountTbl)

	c.logger.Trace("【%v】END:texMoneyHandler SetCollectSales", texCon.GetUniqueKey())
}

// SetOverflowCollectSales あふれ金庫からの売上金回収分を保持情報に対して更新する。売上金回収は完了すると、一式「SetCollectSales」にて更新を行う為
// あふれ分を個別で更新してしまうと、二重にカウントされてしまう。
func (c *texMoneyHandler) SetOverflowCollectSales(texCon *domain.TexContext, update bool, amount int, countTbl [domain.CASH_TYPE_SHITEI]int, exCountTbl [domain.EXTRA_CASH_TYPE_SHITEI]int) {
	c.logger.Trace("【%v】START:texMoneyHandler SetOverflowCollectSales", texCon.GetUniqueKey())

	_, salesSortInfo := c.safeInfoManager.GetSortInfo(texCon, domain.SALES_MONEY_COLLECT)

	if update {

		c.logger.Debug("【%v】- before 売上金回収=%v", texCon.GetUniqueKey(), salesSortInfo)

		salesSortInfo.Amount += amount

		for i, v := range countTbl {
			salesSortInfo.CountTbl[i] += v
		}

		for i, v := range exCountTbl {
			salesSortInfo.ExCountTbl[i] += v
		}

		c.safeInfoManager.UpdateSortInfo(texCon, salesSortInfo)
		c.safeInfoManager.UpdateSalesCompleteAmount(texCon, salesSortInfo.Amount)
		c.safeInfoManager.CountUpSalesCompleteCount(texCon)

		// レポート用売上金回収配列データ更新
		c.aggregateManager.UpdateAggregateCountTbl(texCon, c.maintenanceModeMng.GetMode(texCon), domain.SALES_COLLECT_TBL, salesSortInfo.ExCountTbl)

	}

	c.logger.Debug("【%v】texMoneyHandler SetOverflowCollectSales after 売上金回収=%v", texCon.GetUniqueKey(), salesSortInfo)
}

// 金庫情報遷移記録要求情報作成
func (c *texMoneyHandler) RequestReportSafeInfo(texCon *domain.TexContext) (resInfo domain.RequestReportSafeInfo) {
	c.logger.Trace("【%v】START:texMoneyHandler RequestReportSafeInfo", texCon.GetUniqueKey())

	var historySortCode int
	switch c.GetSequence(texCon) {
	case domain.INITIAL_ADDING_CONFIRM, //初期補充確定
		domain.INITIAL_ADDING_UPDATE,           //初期補充現在枚数
		domain.MANUAL_REPLENISHMENT_COLLECTION: //手動補充／回収
		historySortCode = domain.ADD
	case domain.REVERSE_EXCHANGEING_CONFIRM_INDATA: //逆両替確定入金データ時
		if c.exchangePattern == 1 || c.exchangePattern == 2 { // 両替
			historySortCode = domain.DEPOSIT_HIS
		} else {
			historySortCode = domain.OTHER
		}
	case domain.REVERSE_EXCHANGEING_CONFIRM_OUTDATA: //逆両替確定出金データ時
		if c.exchangePattern == 1 || c.exchangePattern == 2 { // 両替
			historySortCode = domain.WITHDRAWAL_HIS
		} else {
			historySortCode = domain.OTHER
		}
	case domain.NUMBER_OF_WITHDRAW_DESIGNATED_CONFIRM,
		domain.REJECTBOXCOLLECT_START,        //リジェクトボックス回収開始
		domain.ALLCOLLECT_START_OUT_STOP,     //全回収開始出金停止
		domain.ALLCOLLECT_START_COLLECT_STOP: //全回収開始回収停止
		historySortCode = domain.OTHER
	case domain.MONEY_ADD_REPLENISH_CONFIRM, //追加補充確定
		domain.MONEY_ADD_REPLENISH_CANCEL: //追加補充取消
		historySortCode = domain.DEPOSIT_HIS
	case domain.MIDDLE_START_OUT_START, //途中回収開始
		domain.SALESMONEY_START,               //売上金回収開始
		domain.ALLCOLLECT_START_OUT_START,     //全回収開始出金開始
		domain.ALLCOLLECT_START_COLLECT_START: //全回収開始回収開始
		historySortCode = domain.WITHDRAWAL_HIS
	case domain.TRANSACTION_DEPOSIT_CONFIRM, //取引入金確定
		domain.TRANSACTION_DEPOSIT_CANCEL, //取引入金取消
		domain.TRANSACTION_OUT_START,      //取引出金開始
		domain.TRANSACTION_OUT_CANCEL:     //取引出金停止
		historySortCode = domain.NORMAL_HIS
	case domain.CLEAR_CASHINFO: //入金ステータス通知クリア
		historySortCode = domain.BUSINESS_HIS //業務
	}

	// 現在有高取得
	safeInfo := c.safeInfoManager.GetSafeInfo(texCon)

	resInfo = *domain.NewRequestReportSafeInfo(c.NewRequestInfo(texCon),
		historySortCode,
		domain.InfoSafeTbl{
			CurrentStatusTbl: c.MakeCurrentStatusTbl(texCon),
			SortInfoTbl:      safeInfo.SortInfoTbl[:],
		})

	c.logger.Trace("【%v】END:texMoneyHandler RequestReportSafeInfo", texCon.GetUniqueKey())
	return
}

// 金庫情報遷移記録要求応答あり
func (c *texMoneyHandler) RecvRequestReportSafeInfo(texCon *domain.TexContext, resInfo domain.ResultReportSafeInfo) {
	c.logger.Trace("【%v】START:texMoneyHandler RecvRequestReportSafeInfo", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:texMoneyHandler RecvRequestReportSafeInfo", texCon.GetUniqueKey())

	statusCash := c.texMoneyNoticeManager.GetStatusCashData(texCon)
	if resInfo.Result {
		statusCash.StatusError = true
		statusCash.ErrorCode = ""
		statusCash.ErrorDetail = ""
	} else {
		c.logger.Error(resInfo.ErrorCode, resInfo.ErrorDetail)
		//エラーコードのセット
		statusCash.StatusError = false
		statusCash.ErrorCode = resInfo.ErrorCode
		statusCash.ErrorDetail = resInfo.ErrorDetail
	}
	c.SetTexmyNoticeStatus(texCon, statusCash) //ステータスの更新
}

// 印刷制御
// 印刷制御の全ての要求の応答を受けるメソッド
func (c *texMoneyHandler) RecvPrintALLRequest(texCon *domain.TexContext, x interface{}) {
	c.logger.Trace("【%v】START: texMoneyHandler RecvPrintALLRequest", texCon.GetUniqueKey())
	//型チェック
	statusReport := c.texMoneyNoticeManager.GetStatusReportData(texCon)
	switch resInfo := x.(type) {
	// case domain.ResultSupply: //補充レシート印刷要求
	// 	statusReport.SlipPrintId = "" //レポート印刷制御管理番号
	// 	if !resInfo.Result {
	// 		statusReport.StatusResult = resInfo.Result //エラーコードのセット
	// 	}
	case domain.PrintStatus: //印刷ステータス通知
		statusReport.StatusPrint = resInfo.StatusPrint
		statusReport.CountPlan = resInfo.CountPlan
		statusReport.CountEnd = resInfo.CountEnd
		if resInfo.StatusPrint == 5 { // 5:印刷完了 の時のみ印刷結果をセット
			statusReport.StatusResult = &resInfo.StatusResult
		} else {
			statusReport.StatusResult = nil
		}
	}
	c.SetTexmyNoticeReportStatusdata(texCon, statusReport)
	c.logger.Trace("【%v】END:texMoneyHandler RecvPrintALLRequest", texCon.GetUniqueKey())
}

// コールバック登録
// 入金ステータスコールバック登録
func (c *texMoneyHandler) RegisterCallbackNoticeIndata(callbackFunc func(texCon *domain.TexContext, noticeInfo *domain.StatusIndata)) {
	c.callbackNoticeIndata = callbackFunc
}

// 出金ステータスコールバック登録
func (c *texMoneyHandler) RegisterCallbackNoticeOutdata(callbackFunc func(texCon *domain.TexContext, noticeInfo *domain.StatusOutdata)) {
	c.callbackNoticeOutdata = callbackFunc
}

// 回収ステータスコールバック登録
func (c *texMoneyHandler) RegisterCallbackNoticeCollectdata(callbackFunc func(texCon *domain.TexContext, noticeInfo *domain.StatusCollectData)) {
	c.callbackNoticeCollectdata = callbackFunc
}

// 有高ステータスコールバック登録
func (c *texMoneyHandler) RegisterCallbackNoticeAmountData(callbackFunc func(texCon *domain.TexContext, noticeInfo *domain.StatusAmount)) {
	c.callbackNoticeAmountData = callbackFunc
}

// 入出金機ステータススコールバック登録
func (c *texMoneyHandler) RegisterCallbackNoticeStatusdata(callbackFunc func(texCon *domain.TexContext, noticeInfo *domain.StatusCash)) {
	c.callbackNoticeStatusdata = callbackFunc
}

// 入出金レポート印刷ステータスコールバック登録
func (c *texMoneyHandler) RegisterCallbackNoticeReportStatusdata(callbackFunc func(texCon *domain.TexContext, noticeInfo *domain.StatusReport)) {
	c.callbackNoticeReportStatusdata = callbackFunc
}

// 両替ステータスコールバック登録
func (c *texMoneyHandler) RegisterCallbackNoticeExchangeStatusdata(callbackFunc func(texCon *domain.TexContext, noticeInfo *domain.StatusExchange)) {
	c.callbackNoticeExchangeStatusdata = callbackFunc
}

func (c *texMoneyHandler) SetMoneySetting(updateData *domain.MoneySetting) {
	c.config.MoneySetting = *updateData
}

func (c *texMoneyHandler) GetMoneySetting() *domain.MoneySetting {
	return &c.config.MoneySetting
}
