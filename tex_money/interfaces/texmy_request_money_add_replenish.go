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

type moneyAddReplenish struct {
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
	praReqInfo       domain.RequestMoneyAddReplenish
	changeStatus     usecases.ChangeStatusRepository
}

// 追加補充要求
func NewRequestMoneyAddReplenish(mqtt handler.MqttRepository, logger handler.LoggerRepository, config config.Configuration, syslogMng usecases.SyslogManager, errorMng usecases.ErrorManager, sendRecv SendRecvRepository, texdtSendRecv TexdtSendRecvRepository, statusTxSendRecv StatusTxSendRecvRepository, printSendRecv PrintSendRecvRepository, texmyHandler usecases.TexMoneyHandlerRepository, changeStatus usecases.ChangeStatusRepository) MoneyAddReplenishRepository {
	return &moneyAddReplenish{
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
		changeStatus:     changeStatus,
	}
}

// 開始処理
func (c *moneyAddReplenish) Start() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_money_add_replenish")
	c.mqtt.Subscribe(topic, c.recvRequest)
}

// 停止処理
func (c *moneyAddReplenish) Stop() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_money_add_replenish")
	c.mqtt.Unsubscribe(topic)
}

// サービス制御要求検出
func (c *moneyAddReplenish) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

func (c *moneyAddReplenish) recvRequest(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_money_add_replenish",
	})

	c.logger.Trace("【%v】START:要求受信 request_money_add_replenish 追加補充要求", texCon.GetUniqueKey())
	var reqInfo domain.RequestMoneyAddReplenish
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_ADD_REPLENISH_FATAL, "", "入出金管理")
		c.logger.Error("moneyAddReplenish recvRequest json.Unmarshal:%v", err)
		return
	}

	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), reqInfo.RequestInfo.RequestID)

	c.praReqInfo = reqInfo //リクエスト情報のセット
	//動作モード判定

	switch reqInfo.StatusMode {
	case domain.MONEY_ADDREPLENISH_START:
		c.logger.Debug("【%v】- 動作モード=開始", texCon.GetUniqueKey())
		c.StatusModeStart(texCon, reqInfo) //開始
	case domain.MONEY_ADDREPLENISH_CANCEL:
		c.logger.Debug("【%v】- 動作モード=取消", texCon.GetUniqueKey())
		c.StatusModeCancel(texCon, reqInfo) //取消
	case domain.MONEY_ADDREPLENISH_CONFIRM:
		c.logger.Debug("【%v】- 動作モード=確定", texCon.GetUniqueKey())
		c.StatusModeConfirm(texCon, reqInfo) //確定

	}

	c.logger.Trace("【%v】END:要求受信 request_money_add_replenish 追加補充要求", texCon.GetUniqueKey())

}

// 処理結果応答:入金開始要求
func (c *moneyAddReplenish) SendResult(texCon *domain.TexContext, reqInfo domain.ResultInStart) bool {
	c.logger.Trace("【%v】START:moneyAddReplenish SendResult", texCon.GetUniqueKey())

	c.logger.Debug("【%v】- c.praReqInfo=%v", texCon.GetUniqueKey(), c.praReqInfo)

	res := domain.ResultMoneyAddReplenish{
		RequestInfo:   c.praReqInfo.RequestInfo,
		Result:        reqInfo.Result,
		ErrorCode:     reqInfo.ErrorCode,
		ErrorDetail:   reqInfo.ErrorDetail,
		CashControlId: reqInfo.CashControlId,
	}

	payment, err := json.Marshal(res)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_ADD_REPLENISH_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Unmarshal:%v", texCon.GetUniqueKey(), err)
	} else {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_ADD_REPLENISH_SUCCESS, "", "入出金管理")
		topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_money_add_replenish")
		c.mqtt.Publish(topic, string(payment))
	}
	c.logger.Trace("【%v】END: moneyAddReplenish SendResult", texCon.GetUniqueKey())
	return true
}

// 処理結果応答:入金終了要求
func (c *moneyAddReplenish) SendResultForInEnd(texCon *domain.TexContext, reqInfo domain.ResultInEnd) bool {
	c.logger.Trace("【%v】START: moneyAddReplenish SendResult", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END: moneyAddReplenish SendResult", texCon.GetUniqueKey())

	c.logger.Debug("【%v】- c.praReqInfo=%v", texCon.GetUniqueKey(), c.praReqInfo)

	res := domain.ResultMoneyAddReplenish{
		RequestInfo:   c.praReqInfo.RequestInfo,
		Result:        reqInfo.Result,
		ErrorCode:     reqInfo.ErrorCode,
		ErrorDetail:   reqInfo.ErrorDetail,
		CashControlId: reqInfo.CashControlId,
	}

	payment, err := json.Marshal(res)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_ADD_REPLENISH_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Marshal:%v", texCon.GetUniqueKey(), err)
		return false
	}

	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_ADD_REPLENISH_SUCCESS, "", "入出金管理")
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_money_add_replenish")
	c.mqtt.Publish(topic, string(payment))

	return true
}

// 処理結果応答:入金終了要求　稼働データ管理で登録後検知
func (c *moneyAddReplenish) SendResultForInEndForDB(texCon *domain.TexContext, resinfo domain.ResultInEnd) bool {
	c.logger.Trace("【%v】START: moneyAddReplenish SendResultForInEndForDB", texCon.GetUniqueKey())

	res := domain.ResultMoneyAddReplenish{
		RequestInfo:   c.praReqInfo.RequestInfo,
		Result:        resinfo.Result,
		ErrorCode:     resinfo.ErrorCode,
		ErrorDetail:   resinfo.ErrorDetail,
		CashControlId: c.praReqInfo.CashControlId, //要求時に来た入金管理IDを返す
	}

	payment, err := json.Marshal(res)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_ADD_REPLENISH_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Marshal:%v", texCon.GetUniqueKey(), err)
	} else {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_ADD_REPLENISH_SUCCESS, "", "入出金管理")
		topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_money_add_replenish")
		c.mqtt.Publish(topic, string(payment))
	}
	c.logger.Trace("【%v】END: moneyAddReplenish SendResultForInEndForDB", texCon.GetUniqueKey())
	return true
}

// 入金要求
func (c *moneyAddReplenish) StatusModeStart(texCon *domain.TexContext, pReqInfo domain.RequestMoneyAddReplenish) {
	c.logger.Trace("【%v】START:moneyAddReplenish StatusModeStart", texCon.GetUniqueKey())
	c.texmyHandler.SetSequence(texCon, domain.MONEY_ADD_REPLENISH_START)

	// 入金開始要求リクエスト情報セット
	resInfo := domain.NewRequestInStart(c.texmyHandler.NewRequestInfo(texCon),
		pReqInfo.ModeOperation,
		pReqInfo.CountClear,
		pReqInfo.TargetDevice)
	// 入金開始要求送信
	c.sendRecv.SendRequestInStart(texCon, resInfo)
	c.logger.Trace("【%v】END:moneyAddReplenish StatusModeStart", texCon.GetUniqueKey())
}

// 取消要求
func (c *moneyAddReplenish) StatusModeCancel(texCon *domain.TexContext, pReqInfo domain.RequestMoneyAddReplenish) {
	c.logger.Trace("【%v】START:moneyAddReplenish StatusModeCancel", texCon.GetUniqueKey())
	c.texmyHandler.SetSequence(texCon, domain.MONEY_ADD_REPLENISH_CANCEL)

	// 入金終了要求リクエスト情報セット
	resInfo := domain.NewRequestInEnd(c.texmyHandler.NewRequestInfo(texCon), pReqInfo.CashControlId, pReqInfo.TargetDevice, pReqInfo.StatusMode)
	// 入金終了要求送信
	c.sendRecv.SendRequestInEnd(texCon, resInfo)
	c.logger.Trace("【%v】END:moneyAddReplenish StatusModeCancel", texCon.GetUniqueKey())
}

// 確定要求
func (c *moneyAddReplenish) StatusModeConfirm(texCon *domain.TexContext, pReqInfo domain.RequestMoneyAddReplenish) {
	c.logger.Trace("【%v】START:moneyAddReplenish StatusModeConfirm", texCon.GetUniqueKey())
	c.texmyHandler.SetSequence(texCon, domain.MONEY_ADD_REPLENISH_CONFIRM)

	// 入金終了要求リクエスト情報セット
	resInfo := domain.NewRequestInEnd(c.texmyHandler.NewRequestInfo(texCon), pReqInfo.CashControlId, pReqInfo.TargetDevice, pReqInfo.StatusMode)
	// 入金終了要求送信
	c.sendRecv.SendRequestInEnd(texCon, resInfo)

	c.logger.Trace("【%v】END:moneyAddReplenish StatusModeConfirm", texCon.GetUniqueKey())
}

// 各要求送信完了検知
func (c *moneyAddReplenish) SenSorSendFinish(texCon *domain.TexContext, reqType int) {
	c.logger.Trace("【%v】START:moneyAddReplenish SenSorSendFinish(reqType=%+v)", texCon.GetUniqueKey(), reqType)
	switch reqType {
	case domain.FINISH_IN_END: //入金禁止要求完了
		if c.texmyHandler.GetSequence(texCon) != domain.MONEY_ADD_REPLENISH_CANCEL { //追加補充取消
			reqInfo := c.texmyHandler.RequestReportSafeInfo(texCon) //確定：稼働データ管理に金庫状態記録を投げる
			c.texdtSendRecv.SendRequestReportSafeInfo(texCon, &reqInfo)
		}

	case domain.FINISH_OUT_START: //出金開始要求完了
		if c.texmyHandler.GetSequence(texCon) != domain.MONEY_ADD_REPLENISH_CANCEL { //追加補充取消
			reqInfo := c.texmyHandler.RequestReportSafeInfo(texCon) //確定：稼働データ管理に金庫状態記録を投げる
			c.texdtSendRecv.SendRequestReportSafeInfo(texCon, &reqInfo)
		}

	case domain.FINISH_REPORT_SAFEINFO: //金庫情報遷移記録完了
		c.SenSorSendFinish(texCon, domain.FINISH_PRINT_CHANGE_SUPPLY)

	case domain.FINISH_PRINT_CHANGE_SUPPLY: //印刷要求完了
		reqInfo := c.changeStatus.RequestChangeSupply(texCon, 0)
		c.statusTxSendRecv.SendRequestChangeSupply(texCon, &reqInfo)

	case domain.FINISH_CHANGE_SUPPLY:
		//完了通知
		c.texmyHandler.SetTexmyNoticeIndata(texCon, true)

		time.Sleep(2 * time.Second) // 2秒待つ 下位レイヤーから有高が上がってくるまでの時間が2秒ほどある為
		c.texmyHandler.SetTexmyNoticeAmountData(texCon)
	}
	c.logger.Trace("【%v】END:moneyAddReplenish SenSorSendFinish", texCon.GetUniqueKey())
}
