package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases"
	"time"
)

type payCash struct {
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
	praReqInfo       domain.RequestPayCash
}

// 取引入金要求
func NewRequestPayCash(mqtt handler.MqttRepository, logger handler.LoggerRepository, config config.Configuration, syslogMng usecases.SyslogManager, errorMng usecases.ErrorManager, sendRecv SendRecvRepository, texdtSendRecv TexdtSendRecvRepository, statusTxSendRecv StatusTxSendRecvRepository, printSendRecv PrintSendRecvRepository, texmyHandler usecases.TexMoneyHandlerRepository) PayCashRepository {
	return &payCash{
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
func (c *payCash) Start() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_pay_cash")
	c.mqtt.Subscribe(topic, c.recvRequest)
}

// 停止処理
func (c *payCash) Stop() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_pay_cash")
	c.mqtt.Unsubscribe(topic)
}

// サービス制御要求検出
func (c *payCash) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

func (c *payCash) recvRequest(message string) {

	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_pay_cash",
	})

	c.logger.Trace("【%v】START:要求受信 request_pay_cash 取引入金要求", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:要求受信 request_pay_cash 取引入金要求", texCon.GetUniqueKey())

	var reqInfo domain.RequestPayCash
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_PAY_CASH_FATAL, "", "入出金管理")
		c.logger.Error("payCash recvRequest json.Unmarshal:%v", err)
		return
	}

	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), reqInfo.RequestInfo.RequestID)

	c.praReqInfo = reqInfo //リクエスト情報のセット
	//動作モード判定
	switch reqInfo.StatusMode {
	case domain.MONEY_PAYCASH_START: // 開始
		c.logger.Debug("【%v】payCash recvRequest mode=開始", texCon.GetUniqueKey())
		c.StatusModeStart(texCon, reqInfo)
	case domain.MONEY_PAYCASH_CANCEL: // 取消
		c.logger.Debug("【%v】payCash recvRequest mode=取消", texCon.GetUniqueKey())
		c.StatusModeCancel(texCon, reqInfo)
	case domain.MONEY_PAYCASH_CONFIRM: //確定
		c.logger.Debug("【%v】payCash recvRequest mode=確定", texCon.GetUniqueKey())
		c.StatusModeConfirm(texCon, reqInfo)
	case domain.MONEY_PAYCASH_END: // 終了
		c.logger.Debug("【%v】payCash recvRequest mode=終了", texCon.GetUniqueKey())
		c.StatusModeEND(texCon, reqInfo)
	}

}

// 処理結果応答：入金許可時
func (c *payCash) SendResult(texCon *domain.TexContext, reqInfo domain.ResultInStart) bool {
	c.logger.Trace("【%v】START: payCash SendResult", texCon.GetUniqueKey())
	//取引入金要求

	c.logger.Debug("【%v】- praReqInfo", texCon.GetUniqueKey(), c.praReqInfo)

	res := domain.ResultPayCash{
		RequestInfo:   c.praReqInfo.RequestInfo,
		Result:        reqInfo.Result,
		ErrorCode:     reqInfo.ErrorCode,
		ErrorDetail:   reqInfo.ErrorDetail,
		CashControlId: reqInfo.CashControlId,
	}

	c.logger.Debug("【%v】- resInfo=%+v", texCon.GetUniqueKey(), res)
	payment, err := json.Marshal(res)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_PAY_CASH_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Marshal:%v", texCon.GetUniqueKey(), err)
		return false
	}
	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_PAY_CASH_SUCCESS, "", "入出金管理")
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_pay_cash")
	c.mqtt.Publish(topic, string(payment))

	c.logger.Trace("【%v】END:payCash SendResult", texCon.GetUniqueKey())
	return true
}

// 処理結果応答：入金禁止時
func (c *payCash) SendResultResultInEnd(texCon *domain.TexContext, reqInfo domain.ResultInEnd) bool {
	c.logger.Trace("【%v】START: payCash SendResult", texCon.GetUniqueKey())
	//取引入金要求

	c.logger.Debug("【%v】- praReqInfo=%v", texCon.GetUniqueKey(), c.praReqInfo)
	res := domain.ResultPayCash{
		RequestInfo:   c.praReqInfo.RequestInfo,
		Result:        reqInfo.Result,
		ErrorCode:     reqInfo.ErrorCode,
		ErrorDetail:   reqInfo.ErrorDetail,
		CashControlId: reqInfo.CashControlId,
	}

	c.logger.Debug("【%v】- resInfo=%+v", texCon.GetUniqueKey(), res)
	payment, err := json.Marshal(res)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_PAY_CASH_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Marshal:%v", texCon.GetUniqueKey(), err)
		return true
	}
	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_PAY_CASH_SUCCESS, "", "入出金管理")
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_pay_cash")
	c.mqtt.Publish(topic, string(payment))

	c.logger.Trace("【%v】END: payCash SendResult", texCon.GetUniqueKey())
	return true
}

// 開始要求
func (c *payCash) StatusModeStart(texCon *domain.TexContext, pReqInfo domain.RequestPayCash) {
	c.logger.Trace("【%v】START: payCash StatusModeStart", texCon.GetUniqueKey())
	c.texmyHandler.SetSequence(texCon, domain.TRANSACTION_DEPOSIT_START)

	// 入金開始要求リクエスト情報セット
	resInfo := domain.NewRequestInStart(c.texmyHandler.NewRequestInfo(texCon),
		pReqInfo.ModeOperation,
		pReqInfo.CountClear,
		pReqInfo.TargetDevice)
	// 入金開始要求送信
	c.sendRecv.SendRequestInStart(texCon, resInfo)
	c.logger.Trace("【%v】END: payCash StatusModeStart", texCon.GetUniqueKey())
}

// 取消要求
func (c *payCash) StatusModeCancel(texCon *domain.TexContext, pReqInfo domain.RequestPayCash) {
	c.logger.Trace("【%v】START:payCash StatusModeCancel", texCon.GetUniqueKey())
	c.texmyHandler.SetSequence(texCon, domain.TRANSACTION_DEPOSIT_CANCEL)

	// 入金終了要求リクエスト情報セット
	resInfo := domain.NewRequestInEnd(c.texmyHandler.NewRequestInfo(texCon), pReqInfo.CashControlId, pReqInfo.TargetDevice, pReqInfo.StatusMode)
	// 入金終了要求送信
	c.sendRecv.SendRequestInEnd(texCon, resInfo)
	c.logger.Trace("【%v】END: payCash StatusModeCancel", texCon.GetUniqueKey())
}

// 確定要求
func (c *payCash) StatusModeConfirm(texCon *domain.TexContext, pReqInfo domain.RequestPayCash) {
	c.logger.Trace("【%v】START: payCash StatusModeConfirm", texCon.GetUniqueKey())
	c.texmyHandler.SetSequence(texCon, domain.TRANSACTION_DEPOSIT_CONFIRM)

	// 入金終了要求リクエスト情報セット
	resInfo := domain.NewRequestInEnd(c.texmyHandler.NewRequestInfo(texCon), pReqInfo.CashControlId, pReqInfo.TargetDevice, pReqInfo.StatusMode)
	// 入金終了要求送信
	c.sendRecv.SendRequestInEnd(texCon, resInfo)
	c.logger.Trace("【%v】END: payCash StatusModeConfirm", texCon.GetUniqueKey())
}

// 終了要求
func (c *payCash) StatusModeEND(texCon *domain.TexContext, pReqInfo domain.RequestPayCash) {
	c.logger.Trace("【%v】START: payCash StatusModeEND TargetDevice =%+v", texCon.GetUniqueKey(), pReqInfo.TargetDevice)

	switch pReqInfo.TargetDevice {
	case 1: // 紙幣のみ
		c.texmyHandler.SetSequence(texCon, domain.TRANSACTION_DEPOSIT_END_BILL)
	case 2: // 硬貨のみ
		c.texmyHandler.SetSequence(texCon, domain.TRANSACTION_DEPOSIT_END_COIN)
	case 0: // 紙幣＆硬貨
		c.texmyHandler.SetSequence(texCon, domain.TRANSACTION_DEPOSIT_END_COIN)
	}

	// 入金終了要求リクエスト情報セット
	resInfo := domain.NewRequestInEnd(c.texmyHandler.NewRequestInfo(texCon), pReqInfo.CashControlId, pReqInfo.TargetDevice, 1)
	// 入金終了要求送信
	c.sendRecv.SendRequestInEnd(texCon, resInfo)

	c.logger.Trace("【%v】END: payCash StatusModeEND", texCon.GetUniqueKey())
}

// 各要求送信完了検知
func (c *payCash) SenSorSendFinish(texCon *domain.TexContext, reqType int) {
	c.logger.Trace("【%v】START: payCash SenSorSendFinish,reqType=%+v", texCon.GetUniqueKey(), reqType)
	switch reqType {
	case domain.FINISH_IN_END: //入金禁止
		reqInfo := c.texmyHandler.RequestReportSafeInfo(texCon) //確定：稼働データ管理に金庫状態記録を投げる
		c.texdtSendRecv.SendRequestReportSafeInfo(texCon, &reqInfo)

	case domain.FINISH_REPORT_SAFEINFO: //金庫情報遷移記録完了
		switch c.texmyHandler.GetSequence(texCon) {
		case domain.TRANSACTION_DEPOSIT_CANCEL:
			c.texmyHandler.SetTexmyNoticeOutdata(texCon, true)
		default:
			//c.SenSorSendFinish(texCon, domain.FINISH_PRINT_CHANGE_SUPPLY) TODO:不要な処理
			//入金ステータス・有高ステータスを通知する
			c.texmyHandler.SetTexmyNoticeIndata(texCon, true)
		}
		time.Sleep(2 * time.Second)                     // 2秒待つ 下位レイヤーから有高が上がってくるまでの時間が2秒ほどある為
		c.texmyHandler.SetTexmyNoticeAmountData(texCon) //ここで502 502に対する有高通知が出るはず

	}
	c.logger.Trace("【%v】END: payCash SenSorSendFinish", texCon.GetUniqueKey())
}
