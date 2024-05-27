package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases"
)

type amountCash struct {
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
	reqInfo          domain.RequestAmountCash
}

// 有高枚数要求
func NewRequestAmountCash(mqtt handler.MqttRepository, logger handler.LoggerRepository, config config.Configuration, syslogMng usecases.SyslogManager, errorMng usecases.ErrorManager, sendRecv SendRecvRepository, texdtSendRecv TexdtSendRecvRepository, statusTxSendRecv StatusTxSendRecvRepository, printSendRecv PrintSendRecvRepository, texmyHandler usecases.TexMoneyHandlerRepository) AmountCashRepository {
	return &amountCash{
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
	}
}

// 開始処理
func (c *amountCash) Start() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_amount_cash")
	c.mqtt.Subscribe(topic, c.recvRequest)
}

// 停止処理
func (c *amountCash) Stop() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_amount_cash")
	c.mqtt.Unsubscribe(topic)
}

// サービス制御要求検出
func (c *amountCash) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

func (c *amountCash) recvRequest(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_amount_cash",
	})
	c.logger.Trace("【%v】START:要求受信 request_amount_cash 有高枚数要求", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:要求受信 request_amount_cash 有高枚数要求", texCon.GetUniqueKey())
	var reqInfo domain.RequestAmountCash
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.logger.Error("amountCash recvRequest json.Unmarshal:%v", err)
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_AMOUNT_CASH_FATAL, "", "入出金管理")
		return
	}
	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), reqInfo.RequestInfo.RequestID)

	c.texmyHandler.SetSequence(texCon, domain.AMOUNT_CASH)
	c.reqInfo = reqInfo //リクエスト情報格納
	//稼働データ管理へ金庫情報を要求
	req := domain.NewRequestGetTermInfoNow(c.texmyHandler.NewRequestInfo(texCon))
	c.texdtSendRecv.SendRequestGetTermInfoNow(texCon, &req)

}

// 処理結果応答
func (c *amountCash) SendResult(texCon *domain.TexContext, pResInfo *domain.ResultGetTermInfoNow) {
	c.logger.Trace("【%v】START: amountCash SendResult,pResInfo=%v", texCon.GetUniqueKey(), pResInfo)
	//データ格納

	if len(pResInfo.InfoSafeTblGetTermNow.SortInfoTbl) == 0 {
		pResInfo.InfoSafeTblGetTermNow.SortInfoTbl = make([]domain.SortInfoTbl, 0)
	}
	tbl := pResInfo.InfoSafeTblGetTermNow.SortInfoTbl[0]

	result := domain.ResultAmountCash{
		RequestInfo: c.reqInfo.RequestInfo,
		Result:      pResInfo.Result,
		Amount:      tbl.Amount,
		CountTbl:    tbl.CountTbl,
		ExCountTbl:  tbl.ExCountTbl,
	}

	if !result.Result {
		result.ErrorCode = pResInfo.ErrorCode
		result.ErrorDetail = pResInfo.ErrorDetail
	}

	payment, err := json.Marshal(result)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_AMOUNT_CASH_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Unmarshal:%v", texCon.GetUniqueKey(), err)
	} else {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_AMOUNT_CASH_SUCCESS, "", "入出金管理")
		topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_amount_cash")
		c.mqtt.Publish(topic, string(payment))
	}
	c.texmyHandler.SetErrorFromRequest(texCon, result.Result, result.ErrorCode, result.ErrorDetail)
	c.logger.Trace("【%v】END: amountCash SendResult", texCon.GetUniqueKey())
}
