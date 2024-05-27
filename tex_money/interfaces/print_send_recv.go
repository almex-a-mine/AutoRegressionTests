package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases"
)

type printSendRecv struct {
	mqtt              handler.MqttRepository
	logger            handler.LoggerRepository
	config            config.Configuration
	syslogMng         usecases.SyslogManager
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
	savePrintId       string //resultで来た印刷管理IDの保存用変数
}

// 印刷制御通信
func NewPrintSendRecv(mqtt handler.MqttRepository, logger handler.LoggerRepository, config config.Configuration, syslogMng usecases.SyslogManager, errorMng usecases.ErrorManager, texMoneyHandler usecases.TexMoneyHandlerRepository, waitManager usecases.IWait) PrintSendRecvRepository {
	return &printSendRecv{
		mqtt:            mqtt,
		logger:          logger,
		config:          config,
		syslogMng:       syslogMng,
		errorMng:        errorMng,
		texMoneyHandler: texMoneyHandler,
		waitManager:     waitManager,
		savePrintId:     "",
	}
}

const TOPIC_PRINT_NUMBER = 3 //Topic数

var topicPrint [TOPIC_PRINT_NUMBER]string
var topicNamePrint = [TOPIC_PRINT_NUMBER]string{
	"result_status", //印刷ステータス要求応答
	"result_supply", //補充レシート印刷要求応答
	"notice_status"} //印刷ステータス通知

// 開始処理
func (c *printSendRecv) Start() {
	var recvFunc = [TOPIC_PRINT_NUMBER]func(string){
		c.RecvResultStatus,
		c.RecvRequestSupply,
		c.RecvNoticeStatus}
	for i := 0; i < TOPIC_PRINT_NUMBER; i++ {
		topicPrint[i] = fmt.Sprintf("%v/%v", domain.TOPIC_HELPERPRINT_BASE, topicNamePrint[i])
		c.mqtt.Subscribe(topicPrint[i], recvFunc[i])
	}

}

// 停止処理
func (c *printSendRecv) Stop() {
	for i := 0; i < len(topicPrint); i++ {
		c.mqtt.Unsubscribe(topicPrint[i])
	}
}

// サービス制御要求検出
func (c *printSendRecv) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

// 初期補充要のアドレスを//現金入出金機制御送受信に渡す
func (c *printSendRecv) SetAddressMoneyIni(moneyInit MoneyInitRepository) {
	c.moneyInit = moneyInit
}
func (c *printSendRecv) SetAddressMoneyExchange(moneyExchange MoneyExchangeRepository) {
	c.moneyExchange = moneyExchange
}
func (c *printSendRecv) SetAddressMoneyAddReplenish(moneyAddReplenish MoneyAddReplenishRepository) {
	c.moneyAddReplenish = moneyAddReplenish
}
func (c *printSendRecv) SetAddressMoneyCollect(moneyCollect MoneyCollectRepository) {
	c.moneyCollect = moneyCollect
}
func (c *printSendRecv) SetAddressSetAmount(setAmount SetAmountRepository) {
	c.setAmount = setAmount
}
func (c *printSendRecv) SetAddressStatusCash(statusCash StatusCashRepository) {
	c.statusCash = statusCash
}
func (c *printSendRecv) SetAddressPayCash(payCash PayCashRepository) {
	c.payCash = payCash
}
func (c *printSendRecv) SetAddressOutCash(outCash OutCashRepository) {
	c.outCash = outCash
}
func (c *printSendRecv) SetAddressAmountCash(amountCash AmountCashRepository) {
	c.amountCash = amountCash
}
func (c *printSendRecv) SetAddressPrintReport(printReport PrintReportRepository) {
	c.printReport = printReport
}
func (c *printSendRecv) SetAddressSalesInfo(salesInfo SalesInfoRepository) {
	c.salesInfo = salesInfo
}
func (c *printSendRecv) SetAddressClearCashInfo(clearCashInfo RequestClearCashInfoRepository) {
	c.clearCashInfo = clearCashInfo
}

// 送信:補充レシート要求
func (c *printSendRecv) SendRequestSupply(texCon *domain.TexContext, resInfo *domain.RequestSupply) {
	topic := domain.TOPIC_HELPERPRINT_BASE + "/request_supply"
	_ = c.sendRequest(texCon, resInfo, topic, &resInfo.RequestInfo)
}

// 送信:印刷ステータス要求
func (c *printSendRecv) SendRequestStatus(texCon *domain.TexContext, resInfo *domain.RequestPrintStatus) {
	topic := domain.TOPIC_HELPERPRINT_BASE + "/request_status"
	_ = c.sendRequest(texCon, resInfo, topic, &resInfo.RequestInfo)
}

// リクエスト送信
// reqInfo:送信するリクエストのjson情報、topic:接頭語付きtopic名称、requestInfo:送信するリクエストのrequestInfo
func (c *printSendRecv) sendRequest(texCon *domain.TexContext, reqInfo interface{}, topic string, requestInfo *domain.RequestInfo) error {
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

// 応答::補充レシート要求要求
func (c *printSendRecv) RecvRequestSupply(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "result_supply",
	})
	var resInfo domain.ResultSupply
	err := json.Unmarshal([]byte(message), &resInfo)
	if err != nil {
		c.logger.Error("printSendRecv RecvRequestSupply json.Unmarshal:%v", err)
		return
	}
	c.logger.Trace("【%v】- RequestID %v", texCon.GetUniqueKey(), resInfo.RequestInfo.RequestID)

	//〇待機情報で検出データをセット
	ok := c.waitManager.SetWaitInfo(texCon, resInfo.RequestInfo.ProcessID, resInfo.RequestInfo.RequestID, resInfo)
	if !ok {
		c.logger.Debug("【%v】待機情報無", texCon.GetUniqueKey())
		return
	}

	c.logger.Debug("【%v】応答受信[印刷制御] result_supply", texCon.GetUniqueKey())
	c.savePrintId = resInfo.PrintId
	// c.texMoneyHandler.RecvPrintALLRequest(texCon, resInfo)
	c.printReport.SendResult(texCon, resInfo)
}

// 応答:印刷ステータス取得要求
func (c *printSendRecv) RecvResultStatus(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "result_status",
	})
	var resInfo domain.PrintStatus
	err := json.Unmarshal([]byte(message), &resInfo)
	if err != nil {
		c.logger.Error("printSendRecv RecvResultStatus json.Unmarshal:%v", err)
		return
	}

	//〇待機情報で検出データをセット
	ok := c.waitManager.SetWaitInfo(texCon, resInfo.RequestInfo.ProcessID, resInfo.RequestInfo.RequestID, resInfo)
	if !ok {
		c.logger.Debug("【%v】待機情報無", texCon.GetUniqueKey())
		return
	}
	c.logger.Debug("【%v】応答受信[印刷制御] result_status", texCon.GetUniqueKey())

}

// 応答:印刷ステータス通知
func (c *printSendRecv) RecvNoticeStatus(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "notice_status",
	})

	var noticeInfo domain.PrintStatus
	err := json.Unmarshal([]byte(message), &noticeInfo)
	if err != nil {
		c.logger.Error("通知受信[印刷制御] 受信データ不正:%v", err)
		return
	}

	// 待機中のIDと通知のIDが異なる場合、待機中IDが存在しない場合
	if c.savePrintId != noticeInfo.PrintId || c.savePrintId == "" {
		c.logger.Debug("【%v】印刷管理ID不一致", texCon.GetUniqueKey())
		return
	}
	c.texMoneyHandler.RecvPrintALLRequest(texCon, noticeInfo)
}
