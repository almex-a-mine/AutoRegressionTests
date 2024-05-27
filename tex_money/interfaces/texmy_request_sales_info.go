package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases"
)

type salesInfo struct {
	mqtt             handler.MqttRepository
	logger           handler.LoggerRepository
	config           config.Configuration
	syslogMng        usecases.SyslogManager
	errorMng         usecases.ErrorManager
	sendRecv         SendRecvRepository
	texdtSendRecv    TexdtSendRecvRepository
	statusTxSendRecv StatusTxSendRecvRepository
	printSendRecv    PrintSendRecvRepository
	texmyHandler     usecases.TexMoneyHandlerRepository
	recvInfo         RecvInfo
	safeInfoManager  usecases.SafeInfoManager
}

type RecvInfo struct {
	processID string //プロセスID
	pcId      string //PCID
	requestID string //リエストID
	result    bool   //通信結果
}

// 売上金情報要求
func NewRequestSalesInfo(mqtt handler.MqttRepository, logger handler.LoggerRepository, config config.Configuration, syslogMng usecases.SyslogManager, errorMng usecases.ErrorManager, sendRecv SendRecvRepository, texdtSendRecv TexdtSendRecvRepository, statusTxSendRecv StatusTxSendRecvRepository, printSendRecv PrintSendRecvRepository, texmyHandler usecases.TexMoneyHandlerRepository, safeInfoManager usecases.SafeInfoManager) SalesInfoRepository {
	return &salesInfo{
		mqtt:             mqtt,
		logger:           logger,
		config:           config,
		syslogMng:        syslogMng,
		errorMng:         errorMng,
		sendRecv:         sendRecv,
		texdtSendRecv:    texdtSendRecv,
		statusTxSendRecv: statusTxSendRecv,
		printSendRecv:    printSendRecv,
		texmyHandler:     texmyHandler,
		safeInfoManager:  safeInfoManager,
	}
}

// 開始処理
func (c *salesInfo) Start() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_sales_info")
	c.mqtt.Subscribe(topic, c.recvRequest)
}

// 停止処理
func (c *salesInfo) Stop() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_sales_info")
	c.mqtt.Unsubscribe(topic)
}

// サービス制御要求検出
func (c *salesInfo) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

func (c *salesInfo) recvRequest(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_sales_info",
	})
	c.logger.Trace("【%v】START:要求受信 request_sales_info 売上金情報取得要求", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:要求受信 request_sales_info 売上金情報取得要求", texCon.GetUniqueKey())
	var reqInfo domain.RequestSalesInfo
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_SALES_INFO_FATAL, "", "入出金管理")
		c.logger.Error("salesInfo recvRequest json.Unmarshal:%v", err)
		return
	}
	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), reqInfo.RequestInfo.RequestID)

	c.recvInfo.processID = reqInfo.RequestInfo.ProcessID
	c.recvInfo.pcId = reqInfo.RequestInfo.PcId
	c.recvInfo.requestID = reqInfo.RequestInfo.RequestID
	c.recvInfo.result = true

	c.texmyHandler.SetSequence(texCon, domain.SALES_INFO)

	resInfo := domain.NewRequestGetTermInfoNow(c.texmyHandler.NewRequestInfo(texCon))
	c.texdtSendRecv.SendRequestGetTermInfoNow(texCon, &resInfo)

}

// 各要求送信完了検知
func (c *salesInfo) SenSorSendFinish(texCon *domain.TexContext) {
	c.logger.Trace("【%v】START:salesInfo SenSorSendFinish", texCon.GetUniqueKey())

	salesAmount, salesComplete, salesCompleteCount := c.recvSalesInfo(texCon) //レシーブ情報を取りに行く
	//レシーブ情報をセット
	resInfo := domain.ResultSalesInfo{
		RequestInfo: domain.RequestInfo{
			ProcessID: c.recvInfo.processID,
			PcId:      c.recvInfo.pcId,
			RequestID: c.recvInfo.requestID,
		},
		SalesAmount:   salesAmount,
		SalesCount:    salesCompleteCount,
		SalesComplete: salesComplete,
		Result:        c.recvInfo.result,
	}

	amount, err := json.Marshal(resInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_SALES_INFO_FATAL, "", "入出金管理")
		c.logger.Error("salesInfo SendResultSetAmount json.Marshal:%v", err)
		return
	}
	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_SALES_INFO_SUCCESS, "", "入出金管理")

	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_sales_info")
	c.mqtt.Publish(topic, string(amount))

	c.logger.Trace("【%v】END:salesInfo SenSorSendFinish", texCon.GetUniqueKey())
}

// 売上金情報レシーブ情報作成
func (c *salesInfo) recvSalesInfo(texCon *domain.TexContext) (int, int, int) {

	// 取引差引
	_, balance := c.safeInfoManager.GetSortInfo(texCon, domain.TRANSACTION_BALANCE)
	salesAmount := balance.Amount

	// 売上回収の金額
	_, salesCollect := c.safeInfoManager.GetSortInfo(texCon, domain.SALES_MONEY_COLLECT)
	salesComplete := salesCollect.Amount

	// 売上金回収回数
	_, saleDataCounter := c.safeInfoManager.GetSalesInfo()

	c.logger.Debug("【%v】texMoneyHandler recvSalesInfo , SalesAmount=%v, SalesComplete=%v, SaleDataCounter=%v", texCon.GetUniqueKey(), salesAmount, salesComplete, saleDataCounter)
	return salesAmount, salesComplete, saleDataCounter
}
