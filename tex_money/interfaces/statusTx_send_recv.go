package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases"
)

type statusTxSendRecv struct {
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
	reqStatus         *domain.RequestStatusStatusTx
}

// 精算機状態管理送信受信管理
func NewStatusTXSendRecv(mqtt handler.MqttRepository, logger handler.LoggerRepository, errorMng usecases.ErrorManager, texMoneyHandler usecases.TexMoneyHandlerRepository, waitManager usecases.IWait) StatusTxSendRecvRepository {
	return &statusTxSendRecv{
		mqtt:            mqtt,
		logger:          logger,
		errorMng:        errorMng,
		texMoneyHandler: texMoneyHandler,
		waitManager:     waitManager}
}

const TOPIC_STATUSTX_NUMBER = 4 //Topic数

var topicStatusTx [TOPIC_STATUSTX_NUMBER]string
var topicNameStatusTx = [TOPIC_STATUSTX_NUMBER]string{
	"result_change_supply",         //状態変更要求(補充完了)
	"result_change_payment",        //状態変更要求(精算完了)
	"result_status",                //精算機状態取得要求
	"result_change_staffoperation", //状態変更要求(スタッフ操作記録)
}

// 開始処理
func (c *statusTxSendRecv) Start() {
	var recvFunc = [TOPIC_STATUSTX_NUMBER]func(string){
		c.RecvResultChangeSupply,
		c.RecvResultChangePayment,
		c.RecvResultStatus,
		c.RecvResultChangeStaffOperation,
	}
	for i := 0; i < TOPIC_STATUSTX_NUMBER; i++ {
		topicStatusTx[i] = fmt.Sprintf("%v/%v", domain.TOPIC_UNIFUNCSTATUS_BASE, topicNameStatusTx[i])
		c.mqtt.Subscribe(topicStatusTx[i], recvFunc[i])
	}

}

// 停止処理
func (c *statusTxSendRecv) Stop() {
	for i := 0; i < len(topicStatusTx); i++ {
		c.mqtt.Unsubscribe(topicStatusTx[i])
	}
}

// サービス制御要求検出
func (c *statusTxSendRecv) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

// 初期補充要のアドレスを//現金入出金機制御送受信に渡す
func (c *statusTxSendRecv) SetAddressMoneyIni(moneyInit MoneyInitRepository) {
	c.moneyInit = moneyInit
}
func (c *statusTxSendRecv) SetAddressMoneyExchange(moneyExchange MoneyExchangeRepository) {
	c.moneyExchange = moneyExchange
}
func (c *statusTxSendRecv) SetAddressMoneyAddReplenish(moneyAddReplenish MoneyAddReplenishRepository) {
	c.moneyAddReplenish = moneyAddReplenish
}
func (c *statusTxSendRecv) SetAddressMoneyCollect(moneyCollect MoneyCollectRepository) {
	c.moneyCollect = moneyCollect
}
func (c *statusTxSendRecv) SetAddressSetAmount(setAmount SetAmountRepository) {
	c.setAmount = setAmount
}
func (c *statusTxSendRecv) SetAddressStatusCash(statusCash StatusCashRepository) {
	c.statusCash = statusCash
}
func (c *statusTxSendRecv) SetAddressPayCash(payCash PayCashRepository) {
	c.payCash = payCash
}
func (c *statusTxSendRecv) SetAddressOutCash(outCash OutCashRepository) {
	c.outCash = outCash
}
func (c *statusTxSendRecv) SetAddressAmountCash(amountCash AmountCashRepository) {
	c.amountCash = amountCash
}
func (c *statusTxSendRecv) SetAddressPrintReport(printReport PrintReportRepository) {
	c.printReport = printReport
}
func (c *statusTxSendRecv) SetAddressSalesInfo(salesInfo SalesInfoRepository) {
	c.salesInfo = salesInfo
}
func (c *statusTxSendRecv) SetAddressClearCashInfo(clearCashInfo RequestClearCashInfoRepository) {
	c.clearCashInfo = clearCashInfo
}

// 送信:状態変更要求（補充完了）
func (c *statusTxSendRecv) SendRequestChangeSupply(texCon *domain.TexContext, resInfo *domain.RequestChangeSupply) {
	topic := domain.TOPIC_UNIFUNCSTATUS_BASE + "/request_change_supply"
	_ = c.sendRequest(texCon, resInfo, topic, &resInfo.RequestInfo)
}

// 送信:状態変更要求（精算完了）
func (c *statusTxSendRecv) SendRequestChangePayment(texCon *domain.TexContext, resInfo *domain.RequestChangePayment) {
	topic := domain.TOPIC_UNIFUNCSTATUS_BASE + "/request_change_payment"
	_ = c.sendRequest(texCon, resInfo, topic, &resInfo.RequestInfo)
}

// 送信:精算機状態取得要求
func (c *statusTxSendRecv) SendRequestStatus(texCon *domain.TexContext, reqInfo *domain.RequestStatusStatusTx) {
	if reqInfo != nil {
		c.reqStatus = reqInfo
	}
	topic := domain.TOPIC_UNIFUNCSTATUS_BASE + "/request_status"
	_ = c.sendRequest(texCon, reqInfo, topic, &reqInfo.RequestInfo)
}

// 送信:状態変更要求(スタッフ操作記録)
func (c *statusTxSendRecv) SendRequestChangeStaffOperation(texCon *domain.TexContext, reqInfo *domain.RequestChangeStaffOperation) {
	topic := domain.TOPIC_UNIFUNCSTATUS_BASE + "/request_change_staffoperation"
	_ = c.sendRequest(texCon, reqInfo, topic, &reqInfo.RequestInfo)
}

// リクエスト送信
// reqInfo:送信するリクエストのjson情報、topic:接頭語付きtopic名称、requestInfo:送信するリクエストのrequestInfo
func (c *statusTxSendRecv) sendRequest(texCon *domain.TexContext, reqInfo interface{}, topic string, requestInfo *domain.RequestInfo) error {
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

// 応答:状態変更要求（補充完了）
func (c *statusTxSendRecv) RecvResultChangeSupply(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "result_change_supply",
	})
	c.logger.Trace("【%v】START:応答受信[精算機状態管理] result_change_supply", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:応答受信[精算機状態管理] result_change_supply", texCon.GetUniqueKey())

	var resInfo domain.ResultStatusStatusTx
	err := json.Unmarshal([]byte(message), &resInfo)
	if err != nil {
		c.logger.Error("statusTxSendRecv RecvResultStatus json.Unmarshal:%v", err)
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
	case domain.INITIAL_ADDING_CONFIRM:
		c.moneyInit.SenSorSendFinish(texCon, domain.FINISH_CHANGE_SUPPLY)

	case domain.REVERSE_EXCHANGEING_CONFIRM_OUTDATA: //逆両替確定出金データ時
		c.moneyExchange.SenSorOutdataSendFinish(texCon, domain.FINISH_CHANGE_SUPPLY)

	case domain.NUMBER_OF_WITHDRAW_DESIGNATED_CONFIRM: //出金枚数指定両替確定
		c.moneyExchange.SenSorIndataSendFinish(texCon, domain.FINISH_CHANGE_SUPPLY)

	case domain.MONEY_ADD_REPLENISH_CONFIRM: //追加補充確定
		c.moneyAddReplenish.SenSorSendFinish(texCon, domain.FINISH_CHANGE_SUPPLY)

	case domain.REJECTBOXCOLLECT_START, //リジェクトボックス回収確定
		domain.UNRETURNEDCOLLECT_START,
		domain.MANUAL_REPLENISHMENT_COLLECTION:
		c.setAmount.SenSorSendFinish(texCon, domain.FINISH_CHANGE_SUPPLY)

	case domain.MIDDLE_START_OUT_START,
		domain.ALLCOLLECT_START_OUT_START,
		domain.SALESMONEY_START:
		c.moneyCollect.SenSorSendFinish(texCon, domain.FINISH_CHANGE_SUPPLY)

	case domain.TRANSACTION_OUT_REFUND_PAYMENT_OUT_START:
		c.outCash.SenSorSendFinish(texCon, domain.FINISH_REPORT_SAFEINFO)
	}
}

// 応答:状態変更要求（補充完了）
func (c *statusTxSendRecv) RecvResultChangePayment(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "result_change_payment",
	})
	c.logger.Trace("【%v】START:応答受信[精算機状態管理] result_change_payment", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:応答受信[精算機状態管理] result_change_payment", texCon.GetUniqueKey())

	var resInfo domain.ResultChangePayment
	err := json.Unmarshal([]byte(message), &resInfo)
	if err != nil {
		c.logger.Error("statusTxSendRecv RecvResultChangePayment json.Unmarshal:%v", err)
		return
	}
	c.logger.Trace("【%v】- RequestID %v", texCon.GetUniqueKey(), resInfo.RequestInfo.RequestID)

	//〇待機情報で検出データをセット
	ok := c.waitManager.SetWaitInfo(texCon, resInfo.RequestInfo.ProcessID, resInfo.RequestInfo.RequestID, resInfo)
	if !ok {
		c.logger.Debug("【%v】待機情報無", texCon.GetUniqueKey())
		return
	}

	switch c.texMoneyHandler.GetSequence(texCon) {
	case domain.REVERSE_EXCHANGEING_CONFIRM_OUTDATA:
		c.moneyExchange.SenSorOutdataSendFinish(texCon, domain.FINISH_CHANGE_SUPPLY)
	}
}

// 応答:精算機状態取得要求
func (c *statusTxSendRecv) RecvResultStatus(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "result_status",
	})
	var resInfo domain.ResultStatusStatusTx
	err := json.Unmarshal([]byte(message), &resInfo)
	if err != nil {
		c.logger.Error("statusTxSendRecv RecvResultStatus json.Unmarshal:%v", err)
		return
	}
	c.logger.Trace("【%v】- RequestID %v", texCon.GetUniqueKey(), resInfo.RequestInfo.RequestID)

	//〇待機情報で検出データをセット
	ok := c.waitManager.SetWaitInfo(texCon, resInfo.RequestInfo.ProcessID, resInfo.RequestInfo.RequestID, resInfo)
	if !ok {
		c.logger.Debug("【%v】待機情報無\n", texCon.GetUniqueKey())
		return
	}

	c.logger.Debug("【%v】応答受信[精算機状態管理] result_status", texCon.GetUniqueKey())
}

// 応答:状態変更要求(スタッフ操作記録)
func (c *statusTxSendRecv) RecvResultChangeStaffOperation(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "result_change_staffoperation",
	})
	c.logger.Trace("【%v】START:応答受信[精算機状態管理] result_change_staffoperation", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:応答受信[精算機状態管理] result_change_staffoperation", texCon.GetUniqueKey())

	var resInfo domain.ResultChangeStaffOperation
	err := json.Unmarshal([]byte(message), &resInfo)
	if err != nil {
		c.logger.Error("statusTxSendRecv RecvResultChangeStaffOperation json.Unmarshal:%v", err)
		return
	}
	c.logger.Trace("【%v】- RequestID %v", texCon.GetUniqueKey(), resInfo.RequestInfo.RequestID)

	//〇待機情報で検出データをセット
	ok := c.waitManager.SetWaitInfo(texCon, resInfo.RequestInfo.ProcessID, resInfo.RequestInfo.RequestID, resInfo)
	if !ok {
		c.logger.Debug("【%v】待機情報無 \n", texCon.GetUniqueKey())
		return
	}
}
