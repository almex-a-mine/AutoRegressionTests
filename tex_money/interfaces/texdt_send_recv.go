package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases"
	"time"
)

type texdtSendRecv struct {
	mqtt              handler.MqttRepository
	logger            handler.LoggerRepository
	errorMng          usecases.ErrorManager
	texMoneyHandler   usecases.TexMoneyHandlerRepository
	waitManager       usecases.IWait
	moneyInit         MoneyInitRepository
	moneyExchange     MoneyExchangeRepository
	moneyAddReplenish MoneyAddReplenishRepository
	moneyCollect      MoneyCollectRepository
	setAmount         SetAmountRepository
	statusCash        StatusCashRepository
	payCash           PayCashRepository
	outCash           OutCashRepository
	amountCash        AmountCashRepository
	printReport       PrintReportRepository
	salesInfo         SalesInfoRepository
	clearCashInfo     RequestClearCashInfoRepository
	reqStatus         *domain.RequestGetTermInfoNow //稼働データ管理：現在状態記録要求
	initializeDbData  bool                          // true : イニシャルGetTermInfoNow実施中 false : 完了
}

// 稼働データ送受信管理
func NewTexdtSendRecv(mqtt handler.MqttRepository,
	logger handler.LoggerRepository,
	errorMng usecases.ErrorManager,
	texMoneyHandler usecases.TexMoneyHandlerRepository,
	waitManager usecases.IWait) TexdtSendRecvRepository {
	return &texdtSendRecv{
		mqtt:            mqtt,
		logger:          logger,
		errorMng:        errorMng,
		texMoneyHandler: texMoneyHandler,
		waitManager:     waitManager}
}

const TOPIC_TEXDT_NUMBER = 2 //Topic数

var topicTexdt [TOPIC_TEXDT_NUMBER]string
var topicNameTexdt = [TOPIC_TEXDT_NUMBER]string{
	"result_report_safeinfo",  //金庫情報遷移記録要求
	"result_get_terminfo_now"} //現在端末取得要求

// 開始処理
func (c *texdtSendRecv) Start() {
	var recvFunc = [TOPIC_TEXDT_NUMBER]func(string){
		c.RecvResultReportSafeInfo,
		c.RecvResultGetTermInfoNow}
	for i := 0; i < TOPIC_TEXDT_NUMBER; i++ {
		topicTexdt[i] = fmt.Sprintf("%v/%v", domain.TOPIC_HELPERDBDATA_BASE, topicNameTexdt[i])
		c.mqtt.Subscribe(topicTexdt[i], recvFunc[i])
	}
}

// 停止処理
func (c *texdtSendRecv) Stop() {
	for i := 0; i < len(topicTexdt); i++ {
		c.mqtt.Unsubscribe(topicTexdt[i])
	}
}

// サービス制御要求検出
func (c *texdtSendRecv) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

// 稼働データイニシャル動作
func (c *texdtSendRecv) InitialDbData() {
	texCon := domain.NewTexContext(domain.RegisterTexContext{})
	c.texMoneyHandler.SetSequence(texCon, domain.INITIAL)
	reqInfo := domain.NewRequestGetTermInfoNow(c.texMoneyHandler.NewRequestInfo(texCon))

	c.initializeDbData = true

	//サービス起動時は現在端末状態を取得できるまでループ
	for {
		c.SendRequestGetTermInfoNow(texCon, &reqInfo)

		if !c.initializeDbData {
			break
		}
	}
}

func (c *texdtSendRecv) GetInitializeDbData() {
	for {

		if !c.initializeDbData {
			return
		}
		c.logger.Debug("DBデータ取得待")
		time.Sleep(1 * time.Second)

	}
}

// 送信:金庫情報遷移記録要求
func (c *texdtSendRecv) SendRequestReportSafeInfo(texCon *domain.TexContext, resInfo *domain.RequestReportSafeInfo) {
	// 以下場合は、送信時にresultを返すロジックに変更
	switch c.texMoneyHandler.GetSequence(texCon) {
	case domain.SET_AMOUNT,
		domain.REJECTBOXCOLLECT_START,
		domain.MANUAL_REPLENISHMENT_COLLECTION,
		domain.UNRETURNEDCOLLECT_START:
		go func() {
			texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "request_report_safeinfo"})
			c.setAmount.SenSorSendFinish(texCon, domain.FINISH_REPORT_SAFEINFO)
		}()
	}
	topic := domain.TOPIC_HELPERDBDATA_BASE + "/request_report_safeinfo"
	_ = c.sendRequest(texCon, resInfo, topic, &resInfo.RequestInfo)
}

// 送信:現在端末取得要求
func (c *texdtSendRecv) SendRequestGetTermInfoNow(texCon *domain.TexContext, reqInfo *domain.RequestGetTermInfoNow) {
	c.reqStatus = reqInfo
	topic := domain.TOPIC_HELPERDBDATA_BASE + "/request_get_terminfo_now"
	_ = c.sendRequest(texCon, reqInfo, topic, &reqInfo.RequestInfo)
}

// リクエスト送信
// reqInfo:送信するリクエストのjson情報、topic:接頭語付きtopic名称、requestInfo:送信するリクエストのrequestInfo
func (c *texdtSendRecv) sendRequest(texCon *domain.TexContext, reqInfo interface{}, topic string, requestInfo *domain.RequestInfo) error {
	c.logger.Trace("【%v】要求送信", texCon.GetUniqueKey())

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

// 応答:金庫情報遷移記録要求
func (c *texdtSendRecv) RecvResultReportSafeInfo(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "result_report_safeinfo",
	})
	c.logger.Trace("【%v】START:応答受信[稼働データ管理] result_report_safeinfo", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:応答受信[稼働データ管理] result_report_safeinfo", texCon.GetUniqueKey())

	var resInfo domain.ResultReportSafeInfo
	err := json.Unmarshal([]byte(message), &resInfo)
	if err != nil {
		c.logger.Error("texdtSendRecv RecvResultReportSafeInfo json.Unmarshal:%v", err)
		return
	}
	c.logger.Trace("【%v】- RequestID %v", texCon.GetUniqueKey(), resInfo.RequestInfo.RequestID)

	//〇待機情報で検出データをセット
	ok := c.waitManager.SetWaitInfo(texCon, resInfo.RequestInfo.ProcessID, resInfo.RequestInfo.RequestID, resInfo)
	if !ok {
		c.logger.Debug("【%v】待機情報無\n", texCon.GetUniqueKey())
		return
	}

	switch c.texMoneyHandler.GetSequence(texCon) {
	case domain.INITIAL_ADDING_CONFIRM,
		domain.INITIAL_ADDING_UPDATE:
		c.moneyInit.SenSorSendFinish(texCon, domain.FINISH_REPORT_SAFEINFO)

	case domain.REVERSE_EXCHANGEING_CONFIRM_INDATA, //逆両替確定入金データ時
		domain.NUMBER_OF_WITHDRAW_DESIGNATED_CONFIRM:
		c.moneyExchange.SenSorIndataSendFinish(texCon, domain.FINISH_CHANGE_SUPPLY)

	case domain.REVERSE_EXCHANGEING_CONFIRM_OUTDATA: //逆両替確定出金データ時
		c.moneyExchange.SenSorOutdataSendFinish(texCon, domain.FINISH_REPORT_SAFEINFO)

	case domain.MONEY_ADD_REPLENISH_CONFIRM: //追加補充確定
		c.moneyAddReplenish.SenSorSendFinish(texCon, domain.FINISH_REPORT_SAFEINFO)

	case domain.MONEY_ADD_REPLENISH_CANCEL: //追加補充取消

	case domain.REJECTBOXCOLLECT_START, //リジェクトボックス回収開始
		domain.MANUAL_REPLENISHMENT_COLLECTION, //手動補充・回収
		domain.UNRETURNEDCOLLECT_START:         //非還流庫回収開始

	case domain.TRANSACTION_DEPOSIT_CONFIRM,
		domain.TRANSACTION_DEPOSIT_END_BILL,
		domain.TRANSACTION_DEPOSIT_END_COIN,
		domain.TRANSACTION_DEPOSIT_CANCEL:
		c.payCash.SenSorSendFinish(texCon, domain.FINISH_REPORT_SAFEINFO)

	case domain.TRANSACTION_OUT_START: //取引出金開始
		c.outCash.SenSorSendFinish(texCon, domain.FINISH_REPORT_SAFEINFO)

	case domain.CLEAR_CASHINFO: //入出金データクリア要求
		c.clearCashInfo.SendResult(texCon, &resInfo)

	case domain.SALESMONEY_START, // 売上金回収
		domain.MIDDLE_START_OUT_START: // 途中回収
		c.moneyCollect.SenSorSendFinish(texCon, domain.FINISH_REPORT_SAFEINFO)
	}
	c.texMoneyHandler.RecvRequestReportSafeInfo(texCon, resInfo)
}

// 応答:現在端末取得要求
func (c *texdtSendRecv) RecvResultGetTermInfoNow(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "result_get_terminfo_now",
	})
	c.logger.Trace("【%v】START:応答受信[稼働データ管理] result_get_terminfo_now", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:応答受信[稼働データ管理] result_get_terminfo_now", texCon.GetUniqueKey())

	var resInfo domain.ResultGetTermInfoNow
	err := json.Unmarshal([]byte(message), &resInfo)
	if err != nil {
		c.logger.Error("texdtSendRecv RecvResultReportSafeInfo json.Unmarshal:%v", err)
		return
	}
	c.logger.Trace("【%v】- RequestID %v", texCon.GetUniqueKey(), resInfo.RequestInfo.RequestID)

	//〇待機情報で検出データをセット
	ok := c.waitManager.SetWaitInfo(texCon, resInfo.RequestInfo.ProcessID, resInfo.RequestInfo.RequestID, resInfo)
	if !ok {
		c.logger.Debug("【%v】待機情報無\n", texCon.GetUniqueKey())
		return
	}

	switch c.texMoneyHandler.GetSequence(texCon) {
	case domain.INITIAL:
		// 初回起動時の送信topic制御用
		// 成功時には、trueの場合falseに変更する
		if c.initializeDbData {
			c.initializeDbData = false
		}
		c.texMoneyHandler.TexdtInfoSave(texCon, resInfo)
		c.texMoneyHandler.SetSequence(texCon, domain.NO_SEQUENCE)

		// TexdtInfoSaveの処理と重複しているためコメントアウト
		// // 金庫情報更新
		// for _, s := range resInfo.InfoSafeTblGetTermNow.SortInfoTbl {
		// 	c.safeInfoMng.UpdateSortInfo(texCon, s)
		// 	if s.SortType == domain.CASH_AVAILABLE {
		// 		// 論理有高更新
		// 		c.safeInfoMng.UpdateAllLogicalCashAvailable(texCon, s)
		// 	}
		// }

	case domain.SALES_INFO: //売上金情報要求
		c.salesInfo.SenSorSendFinish(texCon)
	case domain.AMOUNT_CASH: //有高枚数要求
		c.amountCash.SendResult(texCon, &resInfo)
	case domain.REVERSE_EXCHANGEING_CONFIRM_INDATA:
		c.moneyExchange.SenSorIndataSendFinish(texCon, domain.FINISH_REPORT_SAFEINFO)
	default:
	}
}

// 初期補充要のアドレスを//現金入出金機制御送受信に渡す
func (c *texdtSendRecv) SetAddressMoneyIni(moneyInit MoneyInitRepository) {
	c.moneyInit = moneyInit
}
func (c *texdtSendRecv) SetAddressMoneyExchange(moneyExchange MoneyExchangeRepository) {
	c.moneyExchange = moneyExchange
}
func (c *texdtSendRecv) SetAddressMoneyAddReplenish(moneyAddReplenish MoneyAddReplenishRepository) {
	c.moneyAddReplenish = moneyAddReplenish
}
func (c *texdtSendRecv) SetAddressMoneyCollect(moneyCollect MoneyCollectRepository) {
	c.moneyCollect = moneyCollect
}
func (c *texdtSendRecv) SetAddressSetAmount(setAmount SetAmountRepository) {
	c.setAmount = setAmount
}
func (c *texdtSendRecv) SetAddressStatusCash(statusCash StatusCashRepository) {
	c.statusCash = statusCash
}
func (c *texdtSendRecv) SetAddressPayCash(payCash PayCashRepository) {
	c.payCash = payCash
}
func (c *texdtSendRecv) SetAddressOutCash(outCash OutCashRepository) {
	c.outCash = outCash
}
func (c *texdtSendRecv) SetAddressAmountCash(amountCash AmountCashRepository) {
	c.amountCash = amountCash
}
func (c *texdtSendRecv) SetAddressPrintReport(printReport PrintReportRepository) {
	c.printReport = printReport
}
func (c *texdtSendRecv) SetAddressSalesInfo(salesInfo SalesInfoRepository) {
	c.salesInfo = salesInfo
}
func (c *texdtSendRecv) SetAddressClearCashInfo(clearCashInfo RequestClearCashInfoRepository) {
	c.clearCashInfo = clearCashInfo
}
