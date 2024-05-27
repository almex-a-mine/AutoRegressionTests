package interfaces

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/pkg/calculation"
	"tex_money/usecases"
	"time"
)

type sendRecv struct {
	mqtt                     handler.MqttRepository
	logger                   handler.LoggerRepository
	errorMng                 usecases.ErrorManager
	texMoneyHandler          usecases.TexMoneyHandlerRepository
	waitManager              usecases.IWait
	safeInfoMng              usecases.SafeInfoManager
	texMoneyNoticeManager    usecases.TexMoneyNoticeManagerRepository
	moneyInit                MoneyInitRepository
	moneyExchange            MoneyExchangeRepository
	moneyAddReplenish        MoneyAddReplenishRepository
	moneyCollect             MoneyCollectRepository
	setAmount                SetAmountRepository
	statusCash               StatusCashRepository
	payCash                  PayCashRepository
	outCash                  OutCashRepository
	amountCash               AmountCashRepository
	printReport              PrintReportRepository
	salesInfo                SalesInfoRepository
	clearCashInfo            RequestClearCashInfoRepository
	onlyAmountChangeSequence int //有高変更用のシーケンス一時保管
	statusOfReq              int //request_collectかrequest_printか判断用 1:request_collect 2:request_print
	cashID                   string
	// 初回起動時に、tex_moneyが上位方の要求なくset_Amountを下位サービスに出すために利用している制御変数
	initializeCashCtrl bool // true : イニシャルsetAMount実施中 false : 完了
	// initializeCashStrlRecv bool // true 完了
	recvControlMapMutex    sync.Mutex
	recvControlMap         map[ControlID]*RecvControl
	waitNoticeAmountStatus bool
}

// 現金入出金機制御送信管理
func NewSendRecv(mqtt handler.MqttRepository,
	logger handler.LoggerRepository,
	errorMng usecases.ErrorManager,
	texMoneyHandler usecases.TexMoneyHandlerRepository,
	waitManager usecases.IWait,
	safeInfoMng usecases.SafeInfoManager,
	texMoneyNoticeManager usecases.TexMoneyNoticeManagerRepository) SendRecvRepository {
	return &sendRecv{
		mqtt:                  mqtt,
		logger:                logger,
		errorMng:              errorMng,
		texMoneyHandler:       texMoneyHandler,
		waitManager:           waitManager,
		safeInfoMng:           safeInfoMng,
		recvControlMap:        make(map[ControlID]*RecvControl),
		texMoneyNoticeManager: texMoneyNoticeManager,
	}
}

const TOPIC_NUMBER = 18 //Topic数

var topic [TOPIC_NUMBER]string
var topicName = [TOPIC_NUMBER]string{
	"result_in_start",
	"result_in_end",
	"result_out_start",
	"result_out_stop",
	"result_collect_start",
	"result_collect_stop",
	"result_in_status",
	"result_out_status",
	"result_collect_status",
	"result_amount_status",
	"result_status",
	"result_set_amount",
	"result_scrutiny_start",
	"notice_in_status",
	"notice_out_status",
	"notice_collect_status",
	"notice_amount_status",
	"notice_status"}

// 開始処理
func (c *sendRecv) Start() {
	var recvFunc = [TOPIC_NUMBER]func(string){
		c.RecvResultInStart,
		c.RecvResultInEnd,
		c.RecvResultOutStart,
		c.RecvResultOutStop,
		c.RecvResultCollectStart,
		c.RecvResultCollectStop,
		c.RecvResultInStatus,
		c.RecvResultOutStatus,
		c.RecvResultCollectStatus,
		c.RecvResultAmountStatus,
		c.RecvResultStatus,
		c.RecvResultCashctlSetAmount,
		c.RecvResultScrutinyStart,
		c.RecvNoticeInStatus,
		c.RecvNoticeOutStatus,
		c.RecvNoticeCollectStatus,
		c.RecvNoticeAmountStatus,
		c.RecvNoticeStatus}
	for i := 0; i < TOPIC_NUMBER; i++ {
		topic[i] = fmt.Sprintf("%v/%v", domain.TOPIC_CASHCTL_BASE, topicName[i])
		c.mqtt.Subscribe(topic[i], recvFunc[i])
	}
}

// 初回起動時に、tex_moneyが上位方の要求なくset_Amountを下位サービスに出すために利用している制御変数のon/off
func (c *sendRecv) InitializeCashCtrlFlagOn(change bool) {
	c.logger.Debug("変更前 InitializeCashCtrlFlagOn c.initializeCashCtrl=%t", c.initializeCashCtrl)
	if change {
		c.initializeCashCtrl = true
	} else {
		c.initializeCashCtrl = false
	}
	c.logger.Debug("変更後 InitializeCashCtrlFlagOn c.initializeCashCtrl=%t", c.initializeCashCtrl)
}

// 停止処理
func (c *sendRecv) Stop() {
	for i := 0; i < len(topic); i++ {
		c.mqtt.Unsubscribe(topic[i])
	}
}

// サービス制御要求検出
func (c *sendRecv) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

// registerRecvControlMap MAPにチャネルを登録
// MAP操作に対するロック処理漏れが無いように、登録だけの定義を実装
func (c *sendRecv) registerRecvControlMap(controlID ControlID, topic string, channel chan interface{}) {
	// 送信待機ロジック(同トピックに同時に送信しないようにする為)
	for {
		// 保持しているものがなければfor文を抜ける
		if _, ok := c.recvControlMap[controlID]; !ok {
			break
		} else {
			// 存在する場合，同Topicがなければ、for文を抜ける
			var b = false
			for _, v := range c.recvControlMap {
				if v.resultTopic == topic {
					b = true
					break
				}
			}
			if !b {
				break
			}
		}

	}

	c.recvControlMapMutex.Lock()
	defer c.recvControlMapMutex.Unlock()

	c.recvControlMap[controlID] = &RecvControl{resultTopic: topic, channel: channel}

}

// sendErrorReleaseRecvControlMap エラー発生時の処理
// 処理漏れしないように、エラー発生時に実行する
func (c *sendRecv) sendErrorReleaseRecvControlMap(controlID ControlID) {
	// 今のロジックだと100%タイムアウトするので、既にMAPが存在しなければ正常終了したものを見なす。
	if _, ok := c.recvControlMap[controlID]; ok {
		c.recvControlMap[controlID].channel <- errors.New(string(controlID) + "連携時にエラーが発生") // 上位に失敗を通知

		c.recvControlMapMutex.Lock()
		defer c.recvControlMapMutex.Unlock()
		delete(c.recvControlMap, controlID) // map情報削除
	}
}

func MakeControlId(processId, requestId string) ControlID {
	return ControlID(processId + "_" + requestId)
}

func (c *sendRecv) releaseRecvControlMap(id ControlID) {
	if _, ok := c.recvControlMap[id]; ok {
		c.recvControlMapMutex.Lock()
		defer c.recvControlMapMutex.Unlock()
		delete(c.recvControlMap, id) // map情報削除
	}

}

// 受信時エラー共通処理
func (c *sendRecv) RecvError(resultTopic string) {
	// 同Topicを全て解除する
	for key, value := range c.recvControlMap {
		if value.resultTopic == resultTopic {
			value.channel <- errors.New(resultTopic + "連携時にエラーが発生") // 上位に失敗を通知
			c.recvControlMapMutex.Lock()
			delete(c.recvControlMap, key) // map情報削除
			c.recvControlMapMutex.Unlock()
		}
	}
}

// 現金入出金機制御イニシャル動作
func (c *sendRecv) InitializeCashctl(wg *sync.WaitGroup) {
	defer wg.Done() // メソッドの最後で必ず実行されるようにする

	texCon := domain.NewTexContext(domain.RegisterTexContext{})
	c.texMoneyHandler.SetSequence(texCon, domain.INITIAL_CASHCTL)

	reqInfo := domain.RequestAmountStatus{
		RequestInfo: c.texMoneyHandler.NewRequestInfo(texCon),
	}

	//有高ステータス要求を2.1へ投げる
	for {
		c.SendRequestAmountStatus(texCon, c.texMoneyHandler.GetSequence(texCon), &reqInfo)
		if !c.initializeCashCtrl {
			break
		}
	}

	c.logger.Debug("InitializeCashctl()を抜けた c.initializeCashCtrl=%t", c.initializeCashCtrl)

	//入出金機ステータス取得要求を2.1へ投げる
	statusReqInfo := domain.RequestStatus{
		RequestInfo: c.texMoneyHandler.NewRequestInfo(texCon),
	}
	c.SendRequestStatus(texCon, &statusReqInfo)
}

// 初期補充要のアドレスを//現金入出金機制御送受信に渡す
func (c *sendRecv) SetAddressMoneyIni(moneyInit MoneyInitRepository) {
	c.moneyInit = moneyInit
}
func (c *sendRecv) SetAddressMoneyExchange(moneyExchange MoneyExchangeRepository) {
	c.moneyExchange = moneyExchange
}
func (c *sendRecv) SetAddressMoneyAddReplenish(moneyAddReplenish MoneyAddReplenishRepository) {
	c.moneyAddReplenish = moneyAddReplenish
}
func (c *sendRecv) SetAddressMoneyCollect(moneyCollect MoneyCollectRepository) {
	c.moneyCollect = moneyCollect
}
func (c *sendRecv) SetAddressSetAmount(setAmount SetAmountRepository) {
	c.setAmount = setAmount
}
func (c *sendRecv) SetAddressStatusCash(statusCash StatusCashRepository) {
	c.statusCash = statusCash
}
func (c *sendRecv) SetAddressPayCash(payCash PayCashRepository) {
	c.payCash = payCash
}
func (c *sendRecv) SetAddressOutCash(outCash OutCashRepository) {
	c.outCash = outCash
}
func (c *sendRecv) SetAddressAmountCash(amountCash AmountCashRepository) {
	c.amountCash = amountCash
}
func (c *sendRecv) SetAddressPrintReport(printReport PrintReportRepository) {
	c.printReport = printReport
}
func (c *sendRecv) SetAddressSalesInfo(salesInfo SalesInfoRepository) {
	c.salesInfo = salesInfo
}
func (c *sendRecv) SetAddressClearCashInfo(clearCashInfo RequestClearCashInfoRepository) {
	c.clearCashInfo = clearCashInfo
}

var noticeMutex = &sync.Mutex{}

func (c *sendRecv) unmarshalError(texCon *domain.TexContext, topic string, errorType int, err error) {
	c.logger.Error("【%v】- Topic=%s err=%v", texCon.GetUniqueKey(), topic, err)
	errorCode, errorDetail := c.errorMng.GetErrorInfo(errorType)
	c.texMoneyHandler.SetErrorFromRequest(texCon, false, errorCode, errorDetail)
}

// 入金ステータス通知
func (c *sendRecv) RecvNoticeInStatus(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "notice_in_status"})

	c.logger.Trace("【%v】START:通知受信 [現金入出金制御] notice_in_status", texCon.GetUniqueKey())
	noticeMutex.Lock()

	defer func() {
		c.logger.Trace("【%v】END:通知受信 [現金入出金制御] notice_in_status", texCon.GetUniqueKey())
		noticeMutex.Unlock()
	}()

	var stuInfo domain.InStatus
	err := json.Unmarshal([]byte(message), &stuInfo)
	if err != nil {
		c.unmarshalError(texCon, fmt.Sprintf("%v/%v", domain.TOPIC_CASHCTL_BASE, "notice_in_status"), usecases.ERROR_NOTICEINSTATUS_UNMARSHAL, err)
		return
	}

	c.logger.Debug("【%v】- cashControlId=%v", texCon.GetUniqueKey(), stuInfo.CashControlId)

	if stuInfo.Amount != 0 && stuInfo.ErrorCode == "" && c.statusCheck(stuInfo.CoinStatusCode, stuInfo.BillStatusCode) {
		c.SetWaitNoticeAmountStatus(true, texCon)
	}

	ok := c.texMoneyHandler.SensorCashctlNoticeInStatus(texCon, stuInfo) //入金データの登録
	c.logger.Debug("【%v】- ok=%v,BillStatusCode %v,CoinStatusCode %v", texCon.GetUniqueKey(), ok, stuInfo.BillStatusCode, stuInfo.CoinStatusCode)
	if !ok {
		return
	}

	switch c.texMoneyHandler.GetSequence(texCon) {
	case domain.INITIAL_ADDING_CONFIRM:
		if stuInfo.BillStatusCode == domain.IN_PAYMENT_COMPLETED && stuInfo.CoinStatusCode == domain.IN_PAYMENT_COMPLETED { //入金確定時かチェック
			c.moneyInit.SenSorSendFinish(texCon, domain.FINISH_IN_END)
		}

	case domain.REVERSE_EXCHANGEING_CONFIRM_INDATA,
		domain.NUMBER_OF_WITHDRAW_DESIGNATED_CONFIRM:
		if (stuInfo.BillStatusCode == domain.IN_PAYMENT_COMPLETED || stuInfo.BillStatusCode == domain.IN_PAYMENT_ERROR) &&
			(stuInfo.CoinStatusCode == domain.IN_PAYMENT_COMPLETED || stuInfo.CoinStatusCode == domain.IN_PAYMENT_ERROR) { //入金完了 or 入金異常時
			c.moneyExchange.SenSorIndataSendFinish(texCon, domain.FINISH_IN_END)
		}

	case domain.TRANSACTION_DEPOSIT_START: //取引入金開始
		c.payCash.SenSorSendFinish(texCon, domain.FINISH_IN_START)
	}

	/*硬貨ステータスと紙幣ステータスが完了以外はすべてここで通知を上げる
		完了時は関連サービスへの通信完了後完了を通知する
		notice_indata　開始
		notice_indata　途中　←ここまでは下記の処理で実装
	    2.4と1.1に通信後
		notice_indata　完了　←この通知は 2.4と1.1の　　　通信系の完了を検知できる場所に実装

		0円入金確定の場合はAmount=0でnotice_amountが2.1から来ない為、ここで通知を上げてデータを送らない
	*/
	if stuInfo.CoinStatusCode == domain.IN_PAYMENT_COMPLETED &&
		stuInfo.BillStatusCode == domain.IN_PAYMENT_COMPLETED &&
		stuInfo.Amount != 0 {
		return
	}
	c.texMoneyHandler.SetTexmyNoticeIndata(texCon, true)
}

// 補充レシートCashId
func (c *sendRecv) CashId(texCon *domain.TexContext, cashId string) {
	c.cashID = cashId
	c.logger.Trace("【%v】sendRecv CashI=%v", texCon.GetUniqueKey(), cashId)
}

// 出金ステータス通知
func (c *sendRecv) RecvNoticeOutStatus(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "notice_out_status"})
	noticeMutex.Lock()
	c.logger.Trace("【%v】START:通知受信 [現金入出金制御] notice_out_status", texCon.GetUniqueKey())
	defer func() {
		c.logger.Trace("【%v】END:通知受信 [現金入出金制御] notice_out_status", texCon.GetUniqueKey())
		noticeMutex.Unlock()
	}()
	var reqInfo domain.OutStatus
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.unmarshalError(texCon, fmt.Sprintf("%v/%v", domain.TOPIC_CASHCTL_BASE, "notice_out_status"), usecases.ERROR_NOTICEOUTSTATUS_UNMARSHAL, err)
		return
	}

	c.logger.Debug("【%v】- cashControlId=%v", texCon.GetUniqueKey(), reqInfo.CashControlId)

	if reqInfo.Amount != 0 && reqInfo.ErrorCode == "" && c.statusCheck(reqInfo.CoinStatusCode, reqInfo.BillStatusCode) {
		c.SetWaitNoticeAmountStatus(true, texCon)
	}

	/* 現在枚数有高変更要求中に出金ステータス通知が来た場合には
	無視する。
	①FIT-Aでは出金ステータス通知が無い。
	 notice_amountから内部的に出金情報を生成して金庫情報を更新する

	②Fit-Bでは出金ステータス通知が送信される
	 Fit-Bで発生した出金ステータス通知が、後半で補充出金のカウントを上げるが
	 出金完了後にくるnotice_amountでも再度カウントされてしまい、出金枚数が倍になるので破棄する処理を追加した。

	現在有高変更要求では、
	入金・出金処理→現在の有高を更新する
	処理が行われているので
	出金通知で、有高更新フラグをONにしなくても、後でONになる。
	*/
	if c.texMoneyHandler.GetSequence(texCon) == domain.SET_AMOUNT {
		c.logger.Debug("【%v】- 破棄", texCon.GetUniqueKey())
		if reqInfo.CoinStatusCode == domain.OUT_PAYMENT_COMPLETED &&
			reqInfo.BillStatusCode == domain.OUT_PAYMENT_COMPLETED {
			//出金データの登録だけ行ってその後のリクエストの送信は処理完了後に自動で行われる
			c.texMoneyNoticeManager.UpdateStatusBillBox(texCon, reqInfo)
		}
		return
	}

	c.texMoneyHandler.SensorCashctlNoticeOutStatus(texCon, reqInfo) //入出金管理の出金ステータスに2.1の出金ステータスを格納していく

	c.logger.Debug("【%v】- CoinStatusCode=%v,BillStatusCode=%v", texCon.GetUniqueKey(), reqInfo.CoinStatusCode, reqInfo.BillStatusCode)
	if reqInfo.CoinStatusCode == domain.OUT_PAYMENT_COMPLETED && reqInfo.BillStatusCode == domain.OUT_PAYMENT_COMPLETED {
		switch c.texMoneyHandler.GetSequence(texCon) {
		case domain.REVERSE_EXCHANGEING_CONFIRM_OUTDATA: //逆両替確定出金データ時
			c.moneyExchange.SenSorOutdataSendFinish(texCon, domain.FINISH_OUT_START)

		case domain.MONEY_ADD_REPLENISH_CANCEL: //追加補充払出
			c.moneyAddReplenish.SenSorSendFinish(texCon, domain.FINISH_OUT_START)

		case domain.SALESMONEY_START: //売上金回収開始
			c.moneyCollect.SenSorSendFinish(texCon, domain.FINISH_OUT_START)

		case domain.MIDDLE_START_OUT_START,
			domain.MIDDLE_START_OUT_STOP:
			c.moneyCollect.SenSorSendFinish(texCon, domain.FINISH_OUT_START)

		case domain.ALLCOLLECT_START_OUT_START, //全回収開始,
			domain.ALLCOLLECT_START_OUT_STOP: //全回収取消
			c.moneyCollect.SenSorSendFinish(texCon, domain.FINISH_OUT_START)

		case domain.TRANSACTION_OUT_REFUND_PAYMENT_OUT_START: //取引出金開始
			c.outCash.SetOutCashRefund(texCon, reqInfo)
		case domain.TRANSACTION_DEPOSIT_CANCEL:
			//出金額が0円の場合2.1からamountが来ない為、金庫情報遷移記録をスキップして通知のみを上げる
			if reqInfo.Amount == 0 {
				c.payCash.SenSorSendFinish(texCon, domain.FINISH_REPORT_SAFEINFO)
			}
		}
	}

	//No.345の変更で通知はすべて有高に伴う処理が終わった後に変更
	//出金完了以外は通知する
	if reqInfo.CoinStatusCode == domain.OUT_PAYMENT_COMPLETED &&
		reqInfo.BillStatusCode == domain.OUT_PAYMENT_COMPLETED {
		return
	}

	c.texMoneyHandler.SetTexmyNoticeOutdata(texCon, true)

	if c.texMoneyHandler.GetFlagCollect(texCon) == domain.REQUEST_HAVE {
		c.texMoneyHandler.SetTexmyNoticeCollectdata(texCon, true)

		if reqInfo.CoinStatusCode == domain.OUT_PAYMENT_ERROR &&
			reqInfo.BillStatusCode == domain.OUT_PAYMENT_ERROR {
			c.texMoneyHandler.SetFlagCollect(texCon, false) //回収完了時にフラグを初期化
		}
	}

}

// 回収ステータス通知
func (c *sendRecv) RecvNoticeCollectStatus(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "notice_collect_status"})

	c.logger.Trace("【%v】START:通知受信 [現金入出金制御] notice_collect_status", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:通知受信 [現金入出金制御] notice_collect_status", texCon.GetUniqueKey())
	var reqInfo domain.CollectStatus
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.unmarshalError(texCon, fmt.Sprintf("%v/%v", domain.TOPIC_CASHCTL_BASE, "notice_collect_status"), usecases.ERROR_NOTICECOLLECTSTATUS_UNMARSHAL, err)
		return
	}

	c.logger.Debug("【%v】- cashControlId=%v", texCon.GetUniqueKey(), reqInfo.CashControlId)

	if reqInfo.Amount != 0 && reqInfo.ErrorCode == "" && c.statusCheck(reqInfo.CoinStatusCode, reqInfo.BillStatusCode) {
		c.SetWaitNoticeAmountStatus(true, texCon)
	}

	//入出金管理:回収ステータスを投げる
	c.texMoneyHandler.SensorCashctlNoticeCollectStatus(texCon, reqInfo)
	switch c.texMoneyHandler.GetSequence(texCon) {
	case domain.MIDDLE_START_COLLECT_START, domain.MIDDLE_START_COLLECT_STOP: //途中回収開始、途中回収取消
		c.moneyCollect.SenSorSendFinish(texCon, domain.FINISH_COLLECT_START)
	}
}

func (c *sendRecv) statusCheck(coin int, bill int) bool {
	if !c.statusCodeCheck(coin) {
		return false
	}
	if !c.statusCodeCheck(bill) {
		return false
	}
	return true
}
func (c *sendRecv) statusCodeCheck(i int) bool {
	switch i {
	case 103:
		return true
	case 104:
		return true
	case 203:
		return true
	case 204:
		return true
	}
	return false
}

// 有高ステータス通知
func (c *sendRecv) RecvNoticeAmountStatus(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "notice_amount_status"})

	c.logger.Trace("【%v】START:通知受信 [現金入出金制御] notice_amount_status", texCon.GetUniqueKey())
	noticeMutex.Lock()
	defer func() {
		c.logger.Trace("【%v】END:通知受信 [現金入出金制御] notice_amount_status", texCon.GetUniqueKey())
		noticeMutex.Unlock()
	}()

	//TOPICのメッセージを構造体に格納
	var amStatus domain.AmountStatus
	err := json.Unmarshal([]byte(message), &amStatus)
	if err != nil {
		c.unmarshalError(texCon, fmt.Sprintf("%v/%v", domain.TOPIC_CASHCTL_BASE, "notice_amount_status"), usecases.ERROR_NOTICEAMOUNTSTATUS_UNMARSHAL, err)
		return
	}

	c.logger.Trace("【%v】- notice_amount_status=%+v", texCon.GetUniqueKey(), amStatus)

	//有高ステータス通知が来たらDepositTblの値を更新する
	//金庫情報の入金可能枚数の更新
	_, depositNumberSortInfo := c.safeInfoMng.GetSortInfo(texCon, domain.DEPOSIT_NUMBER)
	var amount int
	for i, d := range amStatus.DepositTbl {
		amount += d * domain.Cash[i]
	}
	c.logger.Debug("【%v】- amount=%+v", texCon.GetUniqueKey(), amount)

	depositNumberSortInfo = domain.SortInfoTbl{
		SortType: domain.DEPOSIT_NUMBER,
		Amount:   amount,
		CountTbl: amStatus.DepositTbl}
	c.safeInfoMng.UpdateSortInfo(texCon, depositNumberSortInfo)

	c.SetWaitNoticeAmountStatus(false, texCon)

	//有高枚数変更中かつ硬貨・紙幣結果通知コードが504or509なら　true
	sec := c.texMoneyHandler.GetSequence(texCon)
	c.logger.Trace("【%v】- sec=%+v", texCon.GetUniqueKey(), sec, amStatus.CoinStatusCode, amStatus.BillStatusCode)
	if sec == domain.SET_AMOUNT && c.noticeAmountStatusCheck(amStatus.CoinStatusCode, amStatus.BillStatusCode) {
		c.texMoneyHandler.SetSequence(texCon, c.onlyAmountChangeSequence)                               //有高枚数変更要求前のシーケンスを再設定
		c.texMoneyHandler.RecvSetAmountNoticeAmountStatus(texCon, amStatus, c.onlyAmountChangeSequence) //【Fit-A】有高枚数変更要求の有高ステータスから補充と出金の枚数を通知
	} else if sec != domain.SET_AMOUNT { // 現在枚数変更要求では完了時のみ有高データを1.3の各処理用にデータを格納してく
		c.texMoneyHandler.SensorCashctlNoticeAmountStatus(texCon, amStatus)
	}

	c.logger.Trace("【%v】- 有高完了処理 CoinStatusCode=%v , BillStatusCode=%v", texCon.GetUniqueKey(), amStatus.CoinStatusCode, amStatus.BillStatusCode)
	// 有高完了処理
	switch {
	case c.noticeAmountStatusCheck(amStatus.CoinStatusCode, amStatus.BillStatusCode):
		switch c.texMoneyHandler.GetSequence(texCon) {
		case domain.MANUAL_REPLENISHMENT_COLLECTION: //手動補充
			c.setAmount.SenSorSendFinish(texCon, domain.FINISH_IN_END)

		case domain.REJECTBOXCOLLECT_START: //リジェクトボックス回収完了
			c.setAmount.SenSorSendFinish(texCon, domain.FINISH_IN_END)

		case domain.UNRETURNEDCOLLECT_START: //非還流庫回収完了
			c.setAmount.SenSorSendFinish(texCon, domain.FINISH_IN_END)
		case domain.SCRUTINY_START: //精査モードの時は完了通知のみ上位に送る
			c.texMoneyHandler.SetTexmyNoticeAmountData(texCon)
		}
	case amStatus.CoinStatusCode == domain.AMO_RECEIPT_DATA_NOTIFICATION || amStatus.BillStatusCode == domain.AMO_RECEIPT_DATA_NOTIFICATION: //有高データ通知
		switch c.texMoneyHandler.GetSequence(texCon) {
		// case domain.REVERSE_EXCHANGEING_CONFIRM_INDATA: //入金ステータスで入金完了時に処理するためこのタイミングでは不要
		// 	c.moneyExchange.SenSorIndataSendFinish(texCon, domain.FINISH_IN_END)

		case domain.NUMBER_OF_WITHDRAW_DESIGNATED_CONFIRM:
			c.moneyExchange.SenSorIndataSendFinish(texCon, domain.FINISH_IN_END)

		case domain.EXCHANGEING_CANCEL: //両替取消
			c.moneyExchange.SenSorIndataSendFinish(texCon, domain.FINISH_IN_END)

		case domain.MONEY_ADD_REPLENISH_CONFIRM: //追加補充確定
			c.moneyAddReplenish.SenSorSendFinish(texCon, domain.FINISH_IN_END)

		case domain.MONEY_ADD_REPLENISH_CANCEL: //追加補充取消
			c.moneyAddReplenish.SenSorSendFinish(texCon, domain.FINISH_IN_END)

		case domain.REJECTBOXCOLLECT_START: //リジェクトボックス回収開始
			c.setAmount.SenSorSendFinish(texCon, domain.FINISH_IN_END)

		case domain.UNRETURNEDCOLLECT_START: //非還流庫回収開始
			c.setAmount.SenSorSendFinish(texCon, domain.FINISH_IN_END)

		case domain.TRANSACTION_DEPOSIT_CONFIRM,
			domain.TRANSACTION_DEPOSIT_CANCEL,
			domain.TRANSACTION_DEPOSIT_END_BILL,
			domain.TRANSACTION_DEPOSIT_END_COIN: //取引入金確定、//取引入金取消
			c.payCash.SenSorSendFinish(texCon, domain.FINISH_IN_END)

		case domain.TRANSACTION_OUT_START,
			domain.TRANSACTION_OUT_REFUND_PAYMENT_OUT_START:
			c.outCash.SenSorSendFinish(texCon, domain.FINISH_OUT_START)
		}
	}

	//有高ステータス通知送信
	/*硬貨ステータスと紙幣ステータスが完了以外はすべてここで通知を上げる
		完了時は関連サービスへの通信完了後完了を通知する
		notice_amount　開始
		notice_amount　途中　←ここまでは下記の処理で実装
	    supply_○○系通信
		notice_amount　完了　←この通知はsupply_○○系通信系の完了を検知できる場所に実装
	*/
	/*if sec == domain.SET_AMOUNT &&  ←移動
		c.noticeAmountStatusCheck(amStatus.CoinStatusCode, amStatus.BillStatusCode) {
		c.texMoneyHandler.SetTexmyNoticeAmountData(texCon)
	} else*/
	c.logger.Trace("【%v】- 有高ステータス通知送信=%+v", texCon.GetUniqueKey(), sec, amStatus.CoinStatusCode, amStatus.BillStatusCode)
	/*if sec != domain.SET_AMOUNT {
		c.texMoneyHandler.SetTexmyNoticeAmountData(texCon)
	}*/
}

// noticeAmountStatusCheck 両方が504or509ならtrue
func (c *sendRecv) noticeAmountStatusCheck(coin int, bill int) bool {
	if !c.noticeAmountStatusCodeCheck(coin) {
		return false
	}
	if !c.noticeAmountStatusCodeCheck(bill) {
		return false
	}
	return true
}

func (c *sendRecv) noticeAmountStatusCodeCheck(i int) bool {
	if i == domain.AMO_PAYMENT_COMPLETED || i == domain.AMO_PAYMENT_ERROR {
		return true
	}
	return false
}

// 現金入出金機ステータス通知
var noticeStatusMutex sync.Mutex

func (c *sendRecv) RecvNoticeStatus(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "notice_status"})
	noticeStatusMutex.Lock()
	c.logger.Trace("【%v】START:通知受信 [現金入出金制御] notice_status", texCon.GetUniqueKey())
	defer func() {
		c.logger.Trace("【%v】END:通知受信 [現金入出金制御] notice_status", texCon.GetUniqueKey())
		noticeStatusMutex.Unlock()
	}()

	var reqInfo domain.NoticeStatus
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.unmarshalError(texCon, fmt.Sprintf("%v/%v", domain.TOPIC_CASHCTL_BASE, "notice_status"), usecases.ERROR_NOTICESTATUS_UNMARSHAL, err)
		return
	}

	c.logger.Trace("【%v】- notice_status=%+v", texCon.GetUniqueKey(), reqInfo)

	//通知内容の分析
	c.texMoneyHandler.SensorCashctlNoticeStatus(texCon, reqInfo)
}

// 送信:入金開始要求
func (c *sendRecv) SendRequestInStart(texCon *domain.TexContext, resInfo *domain.RequestInStart) {
	topic := domain.TOPIC_CASHCTL_BASE + "/request_in_start"
	err := c.sendRequest(texCon, resInfo, topic, &resInfo.RequestInfo)
	if err != nil {
		errorCode, errorDetail := c.errorMng.GetErrorInfo(usecases.ERROR_SENDREQUESTINSTART_UNMARSHAL)
		c.texMoneyHandler.SetErrorFromRequest(texCon, false, errorCode, errorDetail)
	}
}

// 送信:入金終了要求
func (c *sendRecv) SendRequestInEnd(texCon *domain.TexContext, resInfo *domain.RequestInEnd) {
	topic := domain.TOPIC_CASHCTL_BASE + "/request_in_end"
	err := c.sendRequest(texCon, resInfo, topic, &resInfo.RequestInfo)
	if err != nil {
		errorCode, errorDetail := c.errorMng.GetErrorInfo(usecases.ERROR_SENDREQUESTINEND_UNMARSHAL)
		c.texMoneyHandler.SetErrorFromRequest(texCon, false, errorCode, errorDetail)
	}
}

// 送信:出金開始要求
func (c *sendRecv) SendRequestOutStart(texCon *domain.TexContext, resInfo *domain.RequestOutStart) {
	// 期待値としてnotice_amount_statusが存在する場合には論理値更新前に待機しておく
	// 両替等の連続処理時に論理値を先に更新してしまい、論理値更新後に更新前の有高ステータスを取得してしまうパターンが存在していた為
	c.LoopWaitNoticeAmountStatusFalse(texCon)

	// 出金時の論理枚数更新
	c.updateLogicalSortType(texCon, resInfo)

	topic := domain.TOPIC_CASHCTL_BASE + "/request_out_start"
	err := c.sendRequest(texCon, resInfo, topic, &resInfo.RequestInfo)
	if err != nil {
		errorCode, errorDetail := c.errorMng.GetErrorInfo(usecases.ERROR_SENDREQUESTOUTSTART_UNMARSHAL)
		c.texMoneyHandler.SetErrorFromRequest(texCon, false, errorCode, errorDetail)
	}
}

// 送信:出金停止要求
func (c *sendRecv) SendRequestCollectStart(texCon *domain.TexContext, resInfo *domain.RequestOutStop) {
	topic := domain.TOPIC_CASHCTL_BASE + "/request_out_stop"
	err := c.sendRequest(texCon, resInfo, topic, &resInfo.RequestInfo)
	if err != nil {
		errorCode, errorDetail := c.errorMng.GetErrorInfo(usecases.ERROR_SENDREQUESTCOLLECTSTART_UNMARSHAL)
		c.texMoneyHandler.SetErrorFromRequest(texCon, false, errorCode, errorDetail)
	}
}

// 送信:回収開始要求
func (c *sendRecv) SendRequestOutStop(texCon *domain.TexContext, resInfo *domain.RequestCollectStart) {
	topic := domain.TOPIC_CASHCTL_BASE + "/request_collect_start"
	err := c.sendRequest(texCon, resInfo, topic, &resInfo.RequestInfo)
	if err != nil {
		errorCode, errorDetail := c.errorMng.GetErrorInfo(usecases.ERROR_SENDREQUESTOUTSTOP_UNMARSHAL)
		c.texMoneyHandler.SetErrorFromRequest(texCon, false, errorCode, errorDetail)
	}
}

// 送信:回収停止要求
func (c *sendRecv) SendRequestCollectStop(texCon *domain.TexContext, resInfo *domain.RequestCollectStop) {
	topic := domain.TOPIC_CASHCTL_BASE + "/request_collect_stop"
	err := c.sendRequest(texCon, resInfo, topic, &resInfo.RequestInfo)
	if err != nil {
		errorCode, errorDetail := c.errorMng.GetErrorInfo(usecases.ERROR_SENDREQUESTCOLLECTSTOP_UNMARSHAL)
		c.texMoneyHandler.SetErrorFromRequest(texCon, false, errorCode, errorDetail)
	}
}

// 送信:入金ステータス取得要求
func (c *sendRecv) SendRequestInStatus(texCon *domain.TexContext, resInfo *domain.RequestRequestInStatus) {
	topic := domain.TOPIC_CASHCTL_BASE + "/request_in_status"
	err := c.sendRequest(texCon, resInfo, topic, &resInfo.RequestInfo)
	if err != nil {
		errorCode, errorDetail := c.errorMng.GetErrorInfo(usecases.ERROR_RECVRESULTCOLLECTSTOP_UNMARSHAL)
		c.texMoneyHandler.SetErrorFromRequest(texCon, false, errorCode, errorDetail)
	}
}

// 送信:出金ステータス取得要求
func (c *sendRecv) SendRequestOutStatus(texCon *domain.TexContext, resInfo *domain.RequestOutStatus) {
	topic := domain.TOPIC_CASHCTL_BASE + "/request_out_status"
	_ = c.sendRequest(texCon, resInfo, topic, &resInfo.RequestInfo)
}

// 送信:回収ステータス取得要求
func (c *sendRecv) SendRequestCollectStatus(texCon *domain.TexContext, resInfo *domain.RequestCollectStatus) {
	topic := domain.TOPIC_CASHCTL_BASE + "/request_collect_status"
	_ = c.sendRequest(texCon, resInfo, topic, &resInfo.RequestInfo)
}

// 送信:有高ステータス取得要求
func (c *sendRecv) SendRequestAmountStatus(texCon *domain.TexContext, statusOfReq int, resInfo *domain.RequestAmountStatus) {
	c.statusOfReq = statusOfReq //締め処理か補充かの判断用
	topic := domain.TOPIC_CASHCTL_BASE + "/request_amount_status"
	_ = c.sendRequest(texCon, resInfo, topic, &resInfo.RequestInfo)
}

// 送信:入出金機ステータス取得要求
func (c *sendRecv) SendRequestStatus(texCon *domain.TexContext, reqInfo *domain.RequestStatus) {
	topic := domain.TOPIC_CASHCTL_BASE + "/request_status"
	_ = c.sendRequest(texCon, reqInfo, topic, &reqInfo.RequestInfo)
}

// 送信:有高枚数変更要求
func (c *sendRecv) SendRequestCashctlSetAmount(texCon *domain.TexContext, resInfo *domain.RequestCashctlSetAmount, requestType int) {
	c.logger.Trace("【%v】START:有高枚数変更要求 SendRequestCashctlSetAmount", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:有高枚数変更要求 SendRequestCashctlSetAmount", texCon.GetUniqueKey())
	c.safeInfoMng.UpdateAllLogicalCashAvailable(texCon, domain.SortInfoTbl{
		Amount:     resInfo.Amount,
		CountTbl:   resInfo.CountTbl,
		ExCountTbl: resInfo.ExCountTbl,
	})
	c.onlyAmountChangeSequence = requestType
	c.texMoneyHandler.SetSequence(texCon, domain.SET_AMOUNT)

	// 要求内容（論理有高）とデバイス有高が一致している場合は2.1にリクエストを送信しない
	// ここの判断条件修正
	if mismatch, _ := c.safeInfoMng.GetAvailableBalance(texCon); !mismatch {
		c.logger.Debug("【%v】差分なし:下位レイヤーへの有高枚数変更要求送信スキップ", texCon.GetUniqueKey())
		c.texMoneyHandler.SetSequence(texCon, domain.MANUAL_REPLENISHMENT_COLLECTION)
		c.texMoneyHandler.RecvSetAmountNoticeAmountStatus(texCon, domain.AmountStatus{
			Amount:         resInfo.Amount,
			CountTbl:       resInfo.CountTbl,
			ExCountTbl:     resInfo.ExCountTbl,
			CoinStatusCode: 504,
			BillStatusCode: 504,
		}, requestType)
		go c.setAmount.SenSorSendFinish(texCon, domain.SKIP_SEND_SET_AMOUNT)
		// go c.setAmount.SenSorSendFinish(texCon, domain.FINISH_IN_END)
		return
	}

	topic := domain.TOPIC_CASHCTL_BASE + "/request_set_amount"
	_ = c.sendRequest(texCon, resInfo, topic, &resInfo.RequestInfo)
}

// 送信:精査モード開始要求
func (c *sendRecv) SendRequestScrutinyStart(texCon *domain.TexContext, resChan chan interface{}, reqInfo *domain.RequestScrutinyStart) {
	reqTopic := domain.TOPIC_CASHCTL_BASE + "/request_scrutiny_start"
	resTopic := domain.TOPIC_CASHCTL_BASE + "/result_scrutiny_start"
	_ = c.sendRequestPattern2(texCon, resChan, reqInfo, reqTopic, resTopic, &reqInfo.RequestInfo)
}

// リクエスト送信
// reqInfo:送信するリクエストのjson情報、topic:接頭語付きtopic名称、requestInfo:送信するリクエストのrequestInfo
func (c *sendRecv) sendRequest(texCon *domain.TexContext, reqInfo interface{}, topic string, requestInfo *domain.RequestInfo) error {
	c.logger.Trace("【%v】要求送信", texCon.GetUniqueKey())

	// 両替以外の要求では両替フラグを初期化する
	if !domain.IsExchangeSequences(c.texMoneyHandler.GetSequence(texCon)) {
		c.texMoneyHandler.SetFlagExchange(texCon, false)
	}

	// JSON形式に変換
	send, err := json.Marshal(reqInfo)
	if err != nil {
		c.logger.Error("【%v】- json.Marshal:%v", texCon.GetUniqueKey(), err)
		return err
	}

	// 待機情報作成
	waitInfo := c.waitManager.MakeWaitInfo(texCon, requestInfo.ProcessID, requestInfo.RequestID, reqInfo)
	// 不要になった待機情報は最後に必ず削除
	defer c.waitManager.DelWaitInfo(texCon, requestInfo.ProcessID, requestInfo.RequestID)

	// notice_in_status,notice_out_status,notice_collect_statusが来ていて、
	// notice_amount_statusがまだ来ていない場合には待機する。
	c.LoopWaitNoticeAmountStatusFalse(texCon)

	// 送信
	c.mqtt.Publish(topic, string(send))

	// 応答データ待機
	result := c.waitManager.WaitResultInfo(texCon, waitInfo, domain.WAIT_TIME)
	if !result {
		// タイムアウト
		c.logger.Debug("【%v】- wait.WaitResultInfo: timeout TOPIC=%v", texCon.GetUniqueKey(), topic)
		return nil
	}

	// 待機データ取得
	_, ok := c.waitManager.GetWaitInfo(texCon, requestInfo.ProcessID, requestInfo.RequestID)
	if !ok {
		c.logger.Debug("【%v】- wait.GetWaitInfo ok=%v", texCon.GetUniqueKey(), ok)
		return nil
	}

	return nil
}

// リクエスト送信パターン2
// TODO: 将来的にSendRequestとどちらか一方に統一する
func (c *sendRecv) sendRequestPattern2(texCon *domain.TexContext, resChan chan interface{}, reqInfo interface{}, reqTopic string, resTopic string, requestInfo *domain.RequestInfo) error {
	c.logger.Trace("【%v】要求送信2", texCon.GetUniqueKey())

	// トピックの返信先登録
	controlID := MakeControlId(requestInfo.ProcessID, requestInfo.RequestID)
	/// MAPに返却チャネル情報を登録
	c.registerRecvControlMap(controlID, resTopic, resChan)

	send, err := json.Marshal(reqInfo)
	if err != nil {
		c.sendErrorReleaseRecvControlMap(controlID)
		c.logger.Error("【%v】- json.Marshal:%v", texCon.GetUniqueKey(), err)
		return err
	}

	// 待機情報作成
	waitInfo := c.waitManager.MakeWaitInfo(texCon, requestInfo.ProcessID, requestInfo.RequestID, reqInfo)
	// 不要になった待機情報は最後に削除
	defer c.waitManager.DelWaitInfo(texCon, requestInfo.ProcessID, requestInfo.RequestID)

	// 送信
	c.mqtt.Publish(reqTopic, string(send))

	//〇応答データ待機
	result := c.waitManager.WaitResultInfo(texCon, waitInfo, domain.WAIT_TIME)
	if !result {
		if _, ok := c.recvControlMap[controlID]; ok {
			c.logger.Debug("【%v】- wait.WaitResultInfo: timeout TOPIC=%v", texCon.GetUniqueKey(), resTopic)
			c.sendErrorReleaseRecvControlMap(controlID)
		}
		return nil
	}

	//〇待機データ取得
	_, ok := c.waitManager.GetWaitInfo(texCon, requestInfo.ProcessID, requestInfo.RequestID)
	if !ok {
		c.logger.Debug("【%v】- wait.GetWaitInfo error", texCon.GetUniqueKey())
		c.sendErrorReleaseRecvControlMap(controlID)
		return nil
	}

	return nil
}

// 応答:入金開始要求
func (c *sendRecv) RecvResultInStart(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "result_in_start"})

	c.logger.Trace("【%v】START:応答受信[現金入出金制御] result_in_start", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:応答受信[現金入出金制御] result_in_start", texCon.GetUniqueKey())
	var reqInfo domain.ResultInStart
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.unmarshalError(texCon, fmt.Sprintf("%v/%v", domain.TOPIC_CASHCTL_BASE, "result_in_start"), usecases.ERROR_RECVRESULTINSTART_UNMARSHAL, err)
		return
	}

	//〇待機情報で検出データをセット
	ok := c.waitManager.SetWaitInfo(texCon, reqInfo.RequestInfo.ProcessID, reqInfo.RequestInfo.RequestID, reqInfo)
	if !ok {
		c.logger.Debug("【%v】待機情報無\n", texCon.GetUniqueKey())
		return
	}

	var ret bool
	switch c.texMoneyHandler.GetSequence(texCon) {
	case domain.INITIAL_ADDING_START: //初期補充開始
		c.moneyInit.SendResult(texCon, reqInfo.CashControlId) //初期補充要求応答

	case domain.REVERSE_EXCHANGEING_CONFIRM_INDATA: //逆両替確定入金データ時
		c.moneyExchange.SendResult(texCon, reqInfo) //両替要求応答

	case domain.MONEY_ADD_REPLENISH_START: //追加補充開始
		c.moneyAddReplenish.SendResult(texCon, reqInfo) //追加補充要求応答

	case domain.TRANSACTION_DEPOSIT_START: //取引入金開始
		ret = c.payCash.SendResult(texCon, reqInfo) //取引入金要求応答
	}
	if ret {
		c.texMoneyHandler.RecvCashctlALLRequest(texCon, reqInfo)
	}

}

// 応答:入金終了要求
func (c *sendRecv) RecvResultInEnd(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "result_in_end"})

	c.logger.Trace("【%v】START:応答受信[現金入出金制御] result_in_end", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:応答受信[現金入出金制御] result_in_end", texCon.GetUniqueKey())
	var reqInfo domain.ResultInEnd
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.unmarshalError(texCon, fmt.Sprintf("%v/%v", domain.TOPIC_CASHCTL_BASE, "result_in_end"), usecases.ERROR_RECVRESULTINEND_UNMARSHAL, err)
		return
	}
	c.logger.Trace("【%v】- RequestID %v", texCon.GetUniqueKey(), reqInfo.RequestInfo.RequestID)

	//〇待機情報で検出データをセット
	ok := c.waitManager.SetWaitInfo(texCon, reqInfo.RequestInfo.ProcessID, reqInfo.RequestInfo.RequestID, reqInfo)
	if !ok {
		c.logger.Debug("【%v】待機情報無\n", texCon.GetUniqueKey())
		return
	}

	switch c.texMoneyHandler.GetSequence(texCon) {
	case domain.INITIAL_ADDING_CONFIRM: //初期補充確定
		c.moneyInit.SendResult(texCon, reqInfo.CashControlId) //初期補充要求応答

	case domain.REVERSE_EXCHANGEING_CONFIRM_INDATA, //逆両替確定入金データ時
		domain.REVERSE_EXCHANGEING_CONFIRM_OUTDATA,   //逆両替確定出金データ時
		domain.FIVEONECASHTYPE_EXCHANGE_CONFIRM,      //1and5系金種両替確定
		domain.ONECASHTYPE_EXCHANGE_CONFIRM,          //1金種両替確定
		domain.EXCHANGEING_CANCEL,                    //両替取消
		domain.MONEY_EXCHANGE_CONFIRM,                ///両替 動作モード 確定
		domain.NUMBER_OF_WITHDRAW_DESIGNATED_CONFIRM: //出金枚数指定両替確定
		c.moneyExchange.SendResult(texCon, reqInfo) //両替要求応答

	case domain.MONEY_ADD_REPLENISH_CANCEL: ///追加補充要求 取消
		c.moneyAddReplenish.SendResultForInEnd(texCon, reqInfo) //追加補充要求

	case domain.TRANSACTION_DEPOSIT_CANCEL: //取引入金取消
		c.payCash.SendResultResultInEnd(texCon, reqInfo) //取引入金要求応答

	case domain.TRANSACTION_DEPOSIT_CONFIRM: //取引入金確定
		c.payCash.SendResultResultInEnd(texCon, reqInfo) //取引入金要求応答

	case domain.TRANSACTION_DEPOSIT_END_BILL: //取引入金終了紙幣
		c.payCash.SendResultResultInEnd(texCon, reqInfo) //取引入金要求応答

	case domain.TRANSACTION_DEPOSIT_END_COIN: //取引入金終了硬貨
		c.payCash.SendResultResultInEnd(texCon, reqInfo) //取引入金要求応答

	case domain.MONEY_ADD_REPLENISH_CONFIRM:
		go c.moneyAddReplenish.SendResultForInEndForDB(texCon, reqInfo)
	}
	c.texMoneyHandler.RecvCashctlALLRequest(texCon, reqInfo)
}

func (c *sendRecv) updateLogicalSortType(texCon *domain.TexContext, resInfo *domain.RequestOutStart) {
	switch resInfo.OutMode {
	case 0: // 金額指定
		// 釣銭可能枚数を取得
		_, change := c.safeInfoMng.GetSortInfo(texCon, domain.CHANGE_AVAILABLE)
		// 釣銭可能枚数から、出金金額をベースに出金予定枚数26金種を生成
		exCountTbl := calculation.NewCassette(change.ExCountTbl).Exchange(resInfo.Amount, 0)
		// 26金種から10金種を生成
		countTbl := calculation.NewCassette(exCountTbl).ExCountTblToTenCountTbl()
		// 論理有高を更新
		c.safeInfoMng.UpdateOutLogicalCashAvailable(texCon, domain.SortInfoTbl{
			Amount:     resInfo.Amount,
			CountTbl:   countTbl,
			ExCountTbl: exCountTbl,
		})

	case 1: // 枚数指定出金
		// 26金種配列を生成
		exCountTbl := [26]int{}
		// 16金種から26金種を生成
		copy(exCountTbl[:], resInfo.CountTbl[:])
		ex := calculation.NewCassette(exCountTbl)
		// 10金種を取得
		countTbl := ex.ExCountTblToTenCountTbl()
		// 合計金額を取得
		amount := ex.GetTotalAmount()
		// 論理有高を更新
		c.safeInfoMng.UpdateOutLogicalCashAvailable(texCon, domain.SortInfoTbl{
			Amount:     amount,
			CountTbl:   countTbl,
			ExCountTbl: exCountTbl,
		})

	}

}

// 応答:出金開始要求
func (c *sendRecv) RecvResultOutStart(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "result_out_start"})

	c.logger.Trace("【%v】START:応答受信[現金入出金制御] result_out_start", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:応答受信[現金入出金制御] result_out_start", texCon.GetUniqueKey())

	var reqInfo domain.ResultOutStart
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.unmarshalError(texCon, fmt.Sprintf("%v/%v", domain.TOPIC_CASHCTL_BASE, "result_out_start"), usecases.ERROR_RECVRESULTOUTSTART_UNMARSHAL, err)
		return
	}

	//〇待機情報で検出データをセット
	ok := c.waitManager.SetWaitInfo(texCon, reqInfo.RequestInfo.ProcessID, reqInfo.RequestInfo.RequestID, reqInfo)
	if !ok {
		c.logger.Debug("【%v】待機情報無\n", texCon.GetUniqueKey())
		return
	}

	c.texMoneyHandler.RecvCashctlALLRequest(texCon, reqInfo)

	switch c.texMoneyHandler.GetSequence(texCon) {
	case domain.MIDDLE_START_OUT_START: //途中回収開始出金開始
		c.moneyCollect.SendResult(texCon, reqInfo) //回収要求応答

	case domain.SALESMONEY_START: //売上金回収開始
		c.moneyCollect.SendResult(texCon, reqInfo) //回収要求応答

	case domain.ALLCOLLECT_START_OUT_START: //全回収開始出金開始
		c.moneyCollect.SendResult(texCon, reqInfo) //回収要求応答

	case domain.TRANSACTION_OUT_START: //取引出金開始
		c.outCash.SendResult(texCon, reqInfo) //取消出金応答

	case domain.TRANSACTION_OUT_REFUND_PAYMENT_OUT_START: //取引出金 返金残払出開始
		c.outCash.SendResult(texCon, reqInfo) //取消出金応答
		/*	case domain.MIDDLE_AND_SALES_COLLECT:
			c.moneyCollect.SendResult(texCon, reqInfo) //回収要求応答*/
	}
}

// 応答:出金停止要求
func (c *sendRecv) RecvResultCollectStart(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "result_collect_start"})

	c.logger.Trace("【%v】START:応答受信[現金入出金制御] result_collect_start", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:応答受信[現金入出金制御] result_collect_start", texCon.GetUniqueKey())

	var reqInfo domain.ResultOutStop
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.unmarshalError(texCon, fmt.Sprintf("%v/%v", domain.TOPIC_CASHCTL_BASE, "result_collect_start"), usecases.ERROR_RECVRESULTCOLLECTSTART_UNMARSHAL, err)
		return
	}

	//〇待機情報で検出データをセット
	ok := c.waitManager.SetWaitInfo(texCon, reqInfo.RequestInfo.ProcessID, reqInfo.RequestInfo.RequestID, reqInfo)
	if !ok {
		c.logger.Debug("【%v】待機情報無\n", texCon.GetUniqueKey())
		return
	}

	switch c.texMoneyHandler.GetSequence(texCon) {
	case domain.TRANSACTION_OUT_CANCEL: //取引出金取消
		c.outCash.SendResultForOutStop(texCon, reqInfo) //取消出金応答
		c.outCash.SenSorSendFinish(texCon, domain.FINISH_OUT_END)

	case domain.MIDDLE_START_OUT_STOP: //途中回収開始出金停止
		c.moneyCollect.SendResultCollectStart(texCon, reqInfo) //回収要求応答

	case domain.SALESMONEY_START: //売上金回収開始
		c.moneyCollect.SendResultCollectStart(texCon, reqInfo) //回収要求応答

	case domain.ALLCOLLECT_START_OUT_STOP: //全回収開始出金停止
		c.moneyCollect.SendResultCollectStart(texCon, reqInfo) //回収要求応答
	}
	c.texMoneyHandler.RecvCashctlALLRequest(texCon, reqInfo)
}

// 応答:回収開始要求
func (c *sendRecv) RecvResultOutStop(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "result_out_stop"})

	c.logger.Trace("【%v】START:応答受信[現金入出金制御] result_out_stop", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:応答受信[現金入出金制御] result_out_stop", texCon.GetUniqueKey())

	var reqInfo domain.ResultCollectStart
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.unmarshalError(texCon, fmt.Sprintf("%v/%v", domain.TOPIC_CASHCTL_BASE, "result_out_stop"), usecases.ERROR_RECVRESULTOUTSTOP_UNMARSHAL, err)
		return
	}

	//〇待機情報で検出データをセット
	ok := c.waitManager.SetWaitInfo(texCon, reqInfo.RequestInfo.ProcessID, reqInfo.RequestInfo.RequestID, reqInfo)
	if !ok {
		c.logger.Debug("【%v】待機情報無\n", texCon.GetUniqueKey())
		return
	}

	c.texMoneyHandler.RecvCashctlALLRequest(texCon, reqInfo)
	switch c.texMoneyHandler.GetSequence(texCon) {
	case domain.MIDDLE_START_COLLECT_START: //途中回収開始回収開始
		c.moneyCollect.SendResultOutStop(texCon, reqInfo) //回収要求応答

	case domain.ALLCOLLECT_START_COLLECT_START: //全回収開始回収開始
		c.moneyCollect.SendResultOutStop(texCon, reqInfo) //回収要求応答
	}
}

// 応答:回収停止要求
func (c *sendRecv) RecvResultCollectStop(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "result_collect_stop"})
	c.logger.Trace("【%v】START:応答受信[現金入出金制御] result_collect_stop", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:応答受信[現金入出金制御] result_collect_stop", texCon.GetUniqueKey())

	var reqInfo domain.ResultCollectStop
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.unmarshalError(texCon, fmt.Sprintf("%v/%v", domain.TOPIC_CASHCTL_BASE, "result_collect_stop"), usecases.ERROR_RECVRESULTCOLLECTSTOP_UNMARSHAL, err)
		return
	}

	//〇待機情報で検出データをセット
	ok := c.waitManager.SetWaitInfo(texCon, reqInfo.RequestInfo.ProcessID, reqInfo.RequestInfo.RequestID, reqInfo)
	if !ok {
		c.logger.Debug("【%v】待機情報無\n", texCon.GetUniqueKey())
		return
	}

	c.texMoneyHandler.RecvCashctlALLRequest(texCon, reqInfo)
	switch c.texMoneyHandler.GetSequence(texCon) {
	case domain.MIDDLE_START_COLLECT_STOP: //途中回収開始回収停止
		c.moneyCollect.SendResultCollectStop(texCon, reqInfo) //回収要求応答
	case domain.ALLCOLLECT_START_COLLECT_STOP: //全回収開始回収停止
		c.moneyCollect.SendResultCollectStop(texCon, reqInfo) //回収要求応答
	}
}

// 応答:入金ステータス取得要求
func (c *sendRecv) RecvResultInStatus(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "result_in_status"})
	c.logger.Trace("【%v】START:応答受信[現金入出金制御] result_in_status", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:応答受信[現金入出金制御] result_in_status", texCon.GetUniqueKey())

	var reqInfo domain.InStatus
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.logger.Error("【%v】- json.Unmarshal:%v", texCon.GetUniqueKey(), err)
		return
	}

	c.logger.Debug("【%v】- reqInfo=%+v", texCon.GetUniqueKey(), reqInfo)
}

// 応答:出金ステータス取得要求
func (c *sendRecv) RecvResultOutStatus(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "result_out_status"})
	c.logger.Trace("【%v】START:応答受信[現金入出金制御] result_out_status", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:応答受信[現金入出金制御] result_out_status", texCon.GetUniqueKey())

	var reqInfo domain.OutStatus
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.logger.Error("【%v】- json.Unmarshal:%v", texCon.GetUniqueKey(), err)
		return
	}

	c.logger.Debug("【%v】- reqInfo=%+v", texCon.GetUniqueKey(), reqInfo)
}

// 応答:回収ステータス取得要求
func (c *sendRecv) RecvResultCollectStatus(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "result_collect_status"})
	c.logger.Trace("【%v】START:応答受信[現金入出金制御] result_collect_status", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:応答受信[現金入出金制御] result_collect_status", texCon.GetUniqueKey())

	var reqInfo domain.CollectStatus
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.logger.Error("【%v】- json.Unmarshal:%v", texCon.GetUniqueKey(), err)
		return
	}

	c.logger.Debug("【%v】- reqInfo=%+v", texCon.GetUniqueKey(), reqInfo)
}

// 応答:有高ステータス取得要求
func (c *sendRecv) RecvResultAmountStatus(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "result_amount_status"})
	c.logger.Trace("【%v】START:応答受信[現金入出金制御] result_amount_status", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:応答受信[現金入出金制御] result_amount_status", texCon.GetUniqueKey())

	var resInfo domain.ResultAmountStatus
	err := json.Unmarshal([]byte(message), &resInfo)
	if err != nil {
		c.logger.Error("【%v】- json.Unmarshal:%v", texCon.GetUniqueKey(), err)
		return
	}
	c.logger.Debug("【%v】- c.statusOfReq == %v", texCon.GetUniqueKey(), c.statusOfReq)

	//〇待機情報で検出データをセット
	ok := c.waitManager.SetWaitInfo(texCon, resInfo.RequestInfo.ProcessID, resInfo.RequestInfo.RequestID, resInfo)
	if !ok {
		c.logger.Debug("【%v】待機情報無\n", texCon.GetUniqueKey())
		return
	}

	if c.statusOfReq == 4 {
		c.outCash.CheckStatusMode(texCon)
	}

	// 成功時には、trueの場合falseに変更する
	if c.initializeCashCtrl {
		c.InitializeCashCtrlFlagOn(false)
	}

	c.safeInfoMng.UpdateDeviceCashAvailable(texCon, domain.SortInfoTbl{
		Amount:     resInfo.Amount,
		CountTbl:   resInfo.CountTbl,
		ExCountTbl: resInfo.ExCountTbl,
	})

	// 釣銭不一致判定を行う為、notice_amountのロジックに送信する
	c.logger.Debug("START 釣銭不一致判定を行う為、notice_amountのロジックに送信する")
	c.texMoneyHandler.SensorCashctlNoticeAmountStatus(texCon, domain.AmountStatus{
		CoinStatusCode: resInfo.CoinStatusCode,
		BillStatusCode: resInfo.BillStatusCode,
		Amount:         resInfo.Amount,
		CountTbl:       resInfo.CountTbl,
		ExCountTbl:     resInfo.ExCountTbl,
		ErrorCode:      resInfo.ErrorCode,
		ErrorDetail:    resInfo.ErrorDetail,
	})
	c.logger.Debug("END 釣銭不一致判定を行う為、notice_amountのロジックに送信する")

	//入金可能枚数の更新
	_, depositNumberSortInfo := c.safeInfoMng.GetSortInfo(texCon, domain.DEPOSIT_NUMBER)
	var amount int
	for i, d := range resInfo.DepositTbl {
		amount += d * domain.Cash[i]
	}
	depositNumberSortInfo = domain.SortInfoTbl{
		SortType: domain.DEPOSIT_NUMBER,
		Amount:   amount,
		CountTbl: resInfo.DepositTbl}
	c.safeInfoMng.UpdateSortInfo(texCon, depositNumberSortInfo)

	// 起動時に受信したresult情報でのみ動かしたいロジックの調整用フラグ
	c.texMoneyHandler.InitialDiscrepanctStartOne()
}

// 応答:入出金機ステータス取得要求
func (c *sendRecv) RecvResultStatus(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "result_status"})
	c.logger.Trace("【%v】START:応答受信[現金入出金制御] result_status", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:応答受信[現金入出金制御] result_status", texCon.GetUniqueKey())

	var resInfo domain.ResultStatus
	err := json.Unmarshal([]byte(message), &resInfo)
	if err != nil {
		c.logger.Error("【%v】- json.Unmarshal:%v", texCon.GetUniqueKey(), err)
		return
	}

	//〇待機情報で検出データをセット
	ok := c.waitManager.SetWaitInfo(texCon, resInfo.RequestInfo.ProcessID, resInfo.RequestInfo.RequestID, resInfo)
	if !ok {
		c.logger.Debug("【%v】待機情報無\n", texCon.GetUniqueKey())
		return
	}

	// 応答値保存&通知送信
	c.texMoneyHandler.SensorCashctlNoticeStatus(texCon, domain.NoticeStatus{
		CoinStatusCode: resInfo.CoinStatusCode,
		BillStatusCode: resInfo.BillStatusCode,
		CoinNoticeStatusTbl: domain.NoticeStatusTbl{
			StatusError:       resInfo.Result,
			ErrorCode:         resInfo.ErrorCode,
			ErrorDetail:       resInfo.ErrorDetail,
			StatusCover:       resInfo.CoinStatusTbl.StatusCover,
			StatusUnitSet:     resInfo.CoinStatusTbl.StatusUnitSet,
			StatusInCassette:  resInfo.CoinStatusTbl.StatusInCassette,
			StatusOutCassette: resInfo.CoinStatusTbl.StatusOutCassette,
			StatusInsert:      resInfo.CoinStatusTbl.StatusInsert,
			StatusExit:        resInfo.CoinStatusTbl.StatusExit,
			StatusRjbox:       resInfo.CoinStatusTbl.StatusRjbox,
		},
		NoticeCoinResidueInfoTbl: resInfo.CoinResidueInfoTbl,
		BillNoticeStatusTbl: domain.NoticeStatusTbl{
			StatusCover:       resInfo.BillStatusTbl.StatusCover,
			StatusUnitSet:     resInfo.BillStatusTbl.StatusUnitSet,
			StatusInCassette:  resInfo.BillStatusTbl.StatusInCassette,
			StatusOutCassette: resInfo.BillStatusTbl.StatusOutCassette,
			StatusInsert:      resInfo.BillStatusTbl.StatusInsert,
			StatusExit:        resInfo.BillStatusTbl.StatusExit,
			StatusRjbox:       resInfo.BillStatusTbl.StatusRjbox,
		},
		BillNoticeResidueInfoTbl: resInfo.BillResidueInfoTbl,
		DeviceStatusInfoTbl:      resInfo.DeviceStatusInfoTbl,
		WarningInfoTbl:           resInfo.WarningInfoTbl,
	})
}

// 応答:有高枚数変更要求
func (c *sendRecv) RecvResultCashctlSetAmount(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "result_set_amount"})
	c.logger.Trace("【%v】START:応答受信[現金入出金制御] result_set_amount", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:応答受信[現金入出金制御] result_set_amount", texCon.GetUniqueKey())

	var reqInfo domain.ResultCashctlSetAmount
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.logger.Error("【%v】- json.Unmarshal:%v", texCon.GetUniqueKey(), err)
		return
	}
	//〇待機情報で検出データをセット
	ok := c.waitManager.SetWaitInfo(texCon, reqInfo.RequestInfo.ProcessID, reqInfo.RequestInfo.RequestID, reqInfo)
	if !ok {
		c.logger.Debug("【%v】待機情報無", texCon.GetUniqueKey())
		return
	}

	c.setAmount.SendResult(texCon, reqInfo)

	c.texMoneyHandler.RecvCashctlALLRequest(texCon, reqInfo)
}

// 応答:精査モード開始要求
func (c *sendRecv) RecvResultScrutinyStart(message string) {
	topic := domain.TOPIC_CASHCTL_BASE + "/result_scrutiny_start"
	texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "result_scrutiny_start"})
	c.logger.Trace("【%v】START:応答受信[現金入出金制御] result_scrutiny_start", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:応答受信[現金入出金制御] result_scrutiny_start", texCon.GetUniqueKey())

	var resInfo domain.ResultScrutinyStart
	err := json.Unmarshal([]byte(message), &resInfo)
	if err != nil {
		c.logger.Error("【%v】- json.Unmarshal:%v", texCon.GetUniqueKey(), err)
		c.RecvError(topic)
		return
	}

	//〇待機情報で検出データをセット
	ok := c.waitManager.SetWaitInfo(texCon, resInfo.RequestInfo.ProcessID, resInfo.RequestInfo.RequestID, resInfo)
	if !ok {
		c.logger.Debug("【%v】待機情報無", texCon.GetUniqueKey())
		return
	}

	controlID := MakeControlId(resInfo.RequestInfo.ProcessID, resInfo.RequestInfo.RequestID)
	if _, ok := c.recvControlMap[controlID]; ok {
		// 返信元チャネルのチャネルへ情報を返却
		c.recvControlMap[controlID].channel <- resInfo
		// 不要になったMAPを削除
		c.releaseRecvControlMap(controlID)
		return
	}
}

/////////////////////////////////////////
// 論理値を整える為の対応
// notice_in,out,collect後に,notice_amount_statusを待たないと
// 論理値とデバイス有高の整合性がずれる為、ずれないように内部で待つように修正
// 他のサービスが,notice_in,out,collectで判断して次の動作を要求してきた場合
// 期待値のnotice_amount_status受信より早く次の動作に移ってしまい
// 本来発生しない釣銭不一致が発生してしまう為
/////////////////////////////////////////

// SetWaitNoticeAmountStatus notice_amount_status受信を期待値として設定する
func (c *sendRecv) SetWaitNoticeAmountStatus(b bool, texCon *domain.TexContext) {

	a := "有高ステータス通知(notice_amount_status)通知監視有"
	if !b {
		a = "有高ステータス通知_監視無"
	}

	c.logger.Debug("【%v】 %v", texCon.GetUniqueKey(), a)
	c.waitNoticeAmountStatus = b
}

// LoopWaitNoticeAmountStatusFalse notice_amount_status受信期待中の場合には、送信をここで待機する
func (c *sendRecv) LoopWaitNoticeAmountStatusFalse(texCon *domain.TexContext) {

	var i int
	for {

		if !c.waitNoticeAmountStatus {
			break
		}
		c.logger.Debug("【%v】有高ステータス通知待", texCon.GetUniqueKey())
		time.Sleep(100 * time.Millisecond)
		if i == 20 {
			c.SetWaitNoticeAmountStatus(false, texCon)
			c.logger.Debug("【%v】有高ステータス通知_待機時間超過の為破棄", texCon.GetUniqueKey())
			break
		}
		i++
	}
}

/////////////////////////////////////////
