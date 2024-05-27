package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/pkg/calculation"
	"tex_money/usecases"
	"time"
)

type outCash struct {
	mqtt                   handler.MqttRepository
	logger                 handler.LoggerRepository
	config                 config.Configuration
	syslogMng              usecases.SyslogManager
	errorMng               usecases.ErrorManager
	sendRecv               SendRecvRepository
	texdtSendRecv          TexdtSendRecvRepository
	statusTxSendRecv       StatusTxSendRecvRepository
	printSendRecv          PrintSendRecvRepository
	texmyHandler           usecases.TexMoneyHandlerRepository
	praReqInfo             domain.RequestOutCash
	reqOutStart            domain.RequestOutStart //レイヤー2へのリクエスト情報保存
	outPlanAmount          int                    //レイヤーレイヤー2へのリクエスト時の金額
	outStatus              domain.OutStatus       //返金残払出完了時の出金ステータス通知情報保存
	reqPaymentPlanCountTbl [domain.CASH_TYPE_SHITEI]int
	changeStatus           usecases.ChangeStatusRepository
}

// 取引出金要求
func NewRequestOutCash(mqtt handler.MqttRepository, logger handler.LoggerRepository, config config.Configuration, syslogMng usecases.SyslogManager, errorMng usecases.ErrorManager, sendRecv SendRecvRepository, texdtSendRecv TexdtSendRecvRepository, statusTxSendRecv StatusTxSendRecvRepository, printSendRecv PrintSendRecvRepository, texmyHandler usecases.TexMoneyHandlerRepository, changeStatus usecases.ChangeStatusRepository) OutCashRepository {
	return &outCash{
		mqtt:                   mqtt,
		logger:                 logger,
		config:                 config,
		syslogMng:              syslogMng,
		errorMng:               errorMng,
		sendRecv:               sendRecv,
		texdtSendRecv:          texdtSendRecv,
		statusTxSendRecv:       statusTxSendRecv,
		printSendRecv:          printSendRecv,
		texmyHandler:           texmyHandler,
		reqPaymentPlanCountTbl: [domain.CASH_TYPE_SHITEI]int{},
		changeStatus:           changeStatus,
	}
}

// 開始処理
func (c *outCash) Start() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_out_cash")
	c.mqtt.Subscribe(topic, c.recvRequest)
}

// 停止処理
func (c *outCash) Stop() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_out_cash")
	c.mqtt.Unsubscribe(topic)
}

// サービス制御要求検出
func (c *outCash) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

func (c *outCash) recvRequest(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_out_cash",
	})
	c.logger.Trace("【%v】START:要求受信 request_out_cash 取引出金要求", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:要求受信 request_out_cash 取引出金要求", texCon.GetUniqueKey())
	var reqInfo domain.RequestOutCash
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_OUT_CASH_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Unmarshal:%v", texCon.GetUniqueKey(), err)
		return
	}

	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), reqInfo.RequestInfo.RequestID)

	c.praReqInfo = reqInfo
	c.logger.Debug("【%v】- praReqInfo=%v", texCon.GetUniqueKey(), c.praReqInfo)

	// 有高ステータスの再取得不要？と、思われる為お試しコメントアウト
	// amountStatusReqInfo := domain.NewRequestAmountStatus(c.texmyHandler.NewRequestInfo(texCon))
	// c.sendRecv.SendRequestAmountStatus(texCon, 4, amountStatusReqInfo) //送信:有高ステータス取得要求
	c.CheckStatusMode(texCon)

}

// 有高チェック
func (c *outCash) CheckStatusMode(texCon *domain.TexContext) {
	switch c.praReqInfo.StatusMode {
	case domain.MONEY_OUTCASH_START:
		c.StatusModeStart(texCon, c.praReqInfo) //開始
	case domain.MONEY_OUTCASH_STOP:
		c.StatusModeCancel(texCon, c.praReqInfo) //停止
	case domain.MONEY_REFUND_BALANCE_PAYMENT_START: //返金残払出開始
		c.StatusModeRefundBalancePaymentStart(texCon, c.praReqInfo) //返金残払出開始
	}

}

// 処理結果応答:出金開始要求
func (c *outCash) SendResult(texCon *domain.TexContext, reqInfo domain.ResultOutStart) bool {
	c.logger.Trace("【%v】START: outCash SendResult", texCon.GetUniqueKey())
	//取引出金応答

	c.logger.Debug("【%v】- praReqInfo=%v", texCon.GetUniqueKey(), c.praReqInfo)

	res := domain.ResultOutCash{
		RequestInfo:   c.praReqInfo.RequestInfo,
		Result:        reqInfo.Result,
		ErrorCode:     reqInfo.ErrorCode,
		ErrorDetail:   reqInfo.ErrorDetail,
		CashControlId: reqInfo.CashControlId,
		StatusMode:    c.praReqInfo.StatusMode,
	}

	if res.StatusMode == 1 || res.StatusMode == 2 {
		res.PaymentPlanCountTbl = c.reqPaymentPlanCountTbl
	}

	c.logger.Debug("【%v】- resInfo=%v", texCon.GetUniqueKey(), res)

	defer func() { c.reqPaymentPlanCountTbl = [10]int{} }() // 初期化

	payment, err := json.Marshal(res)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_OUT_CASH_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Marshal:%v", texCon.GetUniqueKey(), err)
	} else {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_OUT_CASH_SUCCESS, "", "入出金管理")
		topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_out_cash")
		c.mqtt.Publish(topic, string(payment))
	}

	c.logger.Trace("【%v】END: outCash SendResult", texCon.GetUniqueKey())
	return true
}

// 処理結果応答:出金停止要求
func (c *outCash) SendResultForOutStop(texCon *domain.TexContext, reqInfo domain.ResultOutStop) bool {
	c.logger.Trace("【%v】START: moneyInit SendResultForOutStop(reqInfo=%+v)", texCon.GetUniqueKey(), reqInfo)
	//取引出金応答

	c.logger.Debug("【%v】- praReqInfo=%v", texCon.GetUniqueKey(), c.praReqInfo)

	res := domain.ResultOutCash{
		RequestInfo:   c.praReqInfo.RequestInfo,
		Result:        reqInfo.Result,
		ErrorCode:     reqInfo.ErrorCode,
		ErrorDetail:   reqInfo.ErrorDetail,
		CashControlId: "", //リザルト情報として項目が無い為
		StatusMode:    c.praReqInfo.StatusMode,
	}

	if res.StatusMode == 1 || res.StatusMode == 2 {
		res.PaymentPlanCountTbl = c.reqPaymentPlanCountTbl
	}

	defer func() { c.reqPaymentPlanCountTbl = [10]int{} }() // 初期化
	payment, err := json.Marshal(res)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_OUT_CASH_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Marshal:%v", texCon.GetUniqueKey(), err)
	} else {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_OUT_CASH_SUCCESS, "", "入出金管理")
		topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_out_cash")
		c.mqtt.Publish(topic, string(payment))
	}
	c.logger.Trace("【%v】END: moneyInit SendResultForOutStop", texCon.GetUniqueKey())
	return true
}

// 停止要求
func (c *outCash) StatusModeCancel(texCon *domain.TexContext, pReqInfo domain.RequestOutCash) {
	c.texmyHandler.SetSequence(texCon, domain.TRANSACTION_OUT_CANCEL)

	// 出金停止要求リクエスト情報セット
	resInfo := domain.NewRequestOutStop(c.texmyHandler.NewRequestInfo(texCon), pReqInfo.CashControlId)
	// 出金停止要求送信
	c.sendRecv.SendRequestCollectStart(texCon, &resInfo)
}

// 開始要求
func (c *outCash) StatusModeStart(texCon *domain.TexContext, reqInfo domain.RequestOutCash) {
	c.logger.Trace("【%v】START: outCash StatusModeStart", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END: outCash StatusModeStart", texCon.GetUniqueKey())

	c.texmyHandler.SetSequence(texCon, domain.TRANSACTION_OUT_START)

	resInfo := c.texmyHandler.OutCashStart(texCon, reqInfo)
	c.reqPaymentPlanCountTbl = c.changeCountTbl16IntTo10Int(resInfo.CountTbl) // resultで返却する値を10桁にしてセット

	c.setOutPlanAmount(resInfo) // 払出予定金額をセット

	// 有り高不足で払出予定金額が出金要求金額未満となる場合は出金しないため、2.1への要求をスキップ
	if c.outPlanAmount < reqInfo.OutData {
		c.sendResultSkipOutCash(texCon, reqInfo, resInfo)
		return
	}

	c.sendRecv.SendRequestOutStart(texCon, &resInfo)
}

func (c *outCash) changeCountTbl16IntTo10Int(i [16]int) [10]int {
	var countTbl [domain.CASH_TYPE_SHITEI]int

	countTbl[0] = i[0]
	countTbl[1] = i[1]
	countTbl[2] = i[2]
	countTbl[3] = i[3]
	countTbl[4] = i[4] + i[10]
	countTbl[5] = i[5] + i[11]
	countTbl[6] = i[6] + i[12]
	countTbl[7] = i[7] + i[13]
	countTbl[8] = i[8] + i[14]
	countTbl[9] = i[9] + i[15]

	return countTbl
}

// 払出予定金額をセット
func (c *outCash) setOutPlanAmount(resInfo domain.RequestOutStart) {
	maxInt := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	// レイヤー2への出金要求時に枚数指定出金の場合はamountが0になるため
	// 金種別枚数から払出予定金額を算出する
	var outPlanCountTbl [26]int
	for i := 0; i < len(outPlanCountTbl); i++ {
		if i < 16 {
			outPlanCountTbl[i] = resInfo.CountTbl[i]
		}
	}
	amount := calculation.NewCassette(outPlanCountTbl).GetTotalAmount()

	c.outPlanAmount = maxInt(amount, resInfo.Amount) //出金種別に応じて、どちらかが0円となるため大きい方の値をセットする
}

func (c *outCash) sendResultSkipOutCash(texCon *domain.TexContext, reqInfo domain.RequestOutCash, resInfo domain.RequestOutStart) {
	c.logger.Trace("【%v】START: outCash sendResultSkipOutCash", texCon.GetUniqueKey())

	// 返金残の金種別枚数内訳を算出
	balance := reqInfo.OutData - c.outPlanAmount
	exCountTbl := [26]int{}
	balanceCount := calculation.NewCassette(exCountTbl).AmountToTenCountTbl(balance) //差額の枚数内訳を取得
	for i, v := range balanceCount {
		c.reqPaymentPlanCountTbl[i] += v // 出金予定枚数 に差額の枚数を加算
	}

	c.SendResult(texCon,
		domain.ResultOutStart{
			CashControlId: domain.OUT_CASH_ONE,
			Result:        true,
		})
	c.texmyHandler.SensorFailedNoticeOutData(texCon)
	c.logger.Trace("【%v】END: outCash sendResultSkipOutCash", texCon.GetUniqueKey())
}

// 返金残払出開始要求
func (c *outCash) StatusModeRefundBalancePaymentStart(texCon *domain.TexContext, reqInfo domain.RequestOutCash) {
	c.logger.Trace("【%v】START: outCash StatusModeRefundBalancePaymentStart reqInfo=%+v", texCon.GetUniqueKey(), reqInfo)

	c.texmyHandler.SetSequence(texCon, domain.TRANSACTION_OUT_REFUND_PAYMENT_OUT_START)
	resInfo := c.texmyHandler.OutCashStart(texCon, reqInfo)
	c.reqOutStart = resInfo //返金残情報を算出する為に出金開始時のリクエスト情報を保持しておく

	c.setOutPlanAmount(resInfo) // 払出予定金額をセット

	// 有り高不足で払出予定金額が出金要求金額未満となる場合は出金しないため、2.1への要求をスキップ
	if c.outPlanAmount < reqInfo.OutData {
		c.sendResultSkipOutCash(texCon, reqInfo, resInfo)
		return
	}

	c.sendRecv.SendRequestOutStart(texCon, &c.reqOutStart)

	c.logger.Trace("【%v】END: outCash StatusModeRefundBalancePaymentStart", texCon.GetUniqueKey())
}

// 返金残払出での出金額確定データが来た場合の値の格納
func (c *outCash) SetOutCashRefund(texCon *domain.TexContext, outStatus domain.OutStatus) {
	c.outStatus = outStatus
	c.logger.Debug("【%v】SetOutCashRefund status=%v", texCon.GetUniqueKey(), outStatus)
}

// 各要求送信完了検知
func (c *outCash) SenSorSendFinish(texCon *domain.TexContext, reqType int) {
	c.logger.Trace("【%v】START: outCash SenSorSendFinish,reqType=%v", texCon.GetUniqueKey(), reqType)
	switch c.texmyHandler.GetSequence(texCon) {
	case domain.TRANSACTION_OUT_START,
		domain.MONEY_OUTCASH_STOP:
		switch reqType {
		case domain.FINISH_OUT_START, domain.FINISH_OUT_END: //出金開始/出金停止
			//確定：稼働データ管理に金庫状態記録を投げる
			reqInfo := c.texmyHandler.RequestReportSafeInfo(texCon)
			c.texdtSendRecv.SendRequestReportSafeInfo(texCon, &reqInfo)
		case domain.FINISH_REPORT_SAFEINFO: //金庫情報遷移記録完了
			//完了通知
			c.texmyHandler.SetTexmyNoticeOutdata(texCon, true)

			time.Sleep(2 * time.Second) // 2秒待つ 下位レイヤーから有高が上がってくるまでの時間が2秒ほどある為
			c.texmyHandler.SetTexmyNoticeAmountData(texCon)
		}

	case domain.TRANSACTION_OUT_REFUND_PAYMENT_OUT_START:
		c.logger.Debug("【%v】- 取引出金 返金残払出開始 outPlanAmount=%v,Amount=%v", texCon.GetUniqueKey(), c.outPlanAmount, c.outStatus.Amount)
		switch reqType {
		case domain.FINISH_OUT_START:
			//リクエストで送った出金額
			switch {
			case c.outPlanAmount == c.outStatus.Amount: //返金残無　 8:返金残抜取　返金残払い出しが完了というステータス
				reqInfo := c.changeStatus.RequestChangeSupply(texCon, 1)
				c.statusTxSendRecv.SendRequestChangeSupply(texCon, &reqInfo)
			case c.outPlanAmount > c.outStatus.Amount: //返金残あり　7:返金残払い出しをしているが、残金額があるので、返金残払い出しを継続する必要があるというステータス
				c.logger.Debug("【%v】- 返金残あり", texCon.GetUniqueKey())
				// 返金残情報をセット
				c.changeStatus.SetRefund(texCon, c.reqOutStart, c.outStatus, c.outPlanAmount)
				reqInfo := c.changeStatus.RequestChangeSupply(texCon, 0)
				c.statusTxSendRecv.SendRequestChangeSupply(texCon, &reqInfo)
			case c.outPlanAmount < c.outStatus.Amount: //返金し過ぎ
				c.logger.Debug("【%v】- 過返金", texCon.GetUniqueKey())
			}
			//データクリア
			c.reqOutStart = domain.RequestOutStart{}
			c.outStatus = domain.OutStatus{}
		case domain.FINISH_REPORT_SAFEINFO:
			//完了通知
			c.texmyHandler.SetTexmyNoticeOutdata(texCon, true)

			time.Sleep(2 * time.Second) // 2秒待つ 下位レイヤーから有高が上がってくるまでの時間が2秒ほどある為
			c.texmyHandler.SetTexmyNoticeAmountData(texCon)
		}

	}
	c.logger.Trace("【%v】END: outCash SenSorSendFinish", texCon.GetUniqueKey())
}
