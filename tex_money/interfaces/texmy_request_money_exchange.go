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

type moneyExchange struct {
	mqtt                handler.MqttRepository
	logger              handler.LoggerRepository
	config              config.Configuration
	syslogMng           usecases.SyslogManager
	errorMng            usecases.ErrorManager
	sendRecv            SendRecvRepository
	texdtSendRecv       TexdtSendRecvRepository
	statusTxSendRecv    StatusTxSendRecvRepository
	printSendRecv       PrintSendRecvRepository
	texmyHandler        usecases.TexMoneyHandlerRepository
	praReqInfo          domain.RequestMoneyExchange
	cashIdIndata        string //helperprint　のcashID
	reqAmountStatusInfo domain.RequestAmountStatus
	printdataMng        usecases.PrintDataManager
	aggregateMng        usecases.AggregateManager
	noticeMng           usecases.TexMoneyNoticeManagerRepository
	maintenanceModeMng  usecases.MaintenanceModeManager
	changeStatus        usecases.ChangeStatusRepository
	paymentSendRecv     PaymentSendRecvRepository
}

// 両替要求
func NewRequestMoneyExchange(mqtt handler.MqttRepository,
	logger handler.LoggerRepository,
	config config.Configuration,
	syslogMng usecases.SyslogManager,
	errorMng usecases.ErrorManager,
	sendRecv SendRecvRepository,
	texdtSendRecv TexdtSendRecvRepository,
	statusTxSendRecv StatusTxSendRecvRepository,
	printSendRecv PrintSendRecvRepository,
	texmyHandler usecases.TexMoneyHandlerRepository,
	printdataMng usecases.PrintDataManager,
	aggregateMng usecases.AggregateManager,
	noticeMng usecases.TexMoneyNoticeManagerRepository,
	maintenanceModeMng usecases.MaintenanceModeManager,
	changeStatus usecases.ChangeStatusRepository,
	paymentSendRecv PaymentSendRecvRepository) MoneyExchangeRepository {
	return &moneyExchange{
		mqtt:                mqtt,
		logger:              logger,
		config:              config,
		syslogMng:           syslogMng,
		errorMng:            errorMng,
		sendRecv:            sendRecv,
		texdtSendRecv:       texdtSendRecv,
		statusTxSendRecv:    statusTxSendRecv,
		printSendRecv:       printSendRecv,
		texmyHandler:        texmyHandler,
		reqAmountStatusInfo: domain.RequestAmountStatus{},
		printdataMng:        printdataMng,
		aggregateMng:        aggregateMng,
		noticeMng:           noticeMng,
		maintenanceModeMng:  maintenanceModeMng,
		changeStatus:        changeStatus,
		paymentSendRecv:     paymentSendRecv,
	}
}

// 開始処理
func (c *moneyExchange) Start() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_money_exchange")
	c.mqtt.Subscribe(topic, c.recvRequest)
}

// 停止処理
func (c *moneyExchange) Stop() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_money_exchange")
	c.mqtt.Unsubscribe(topic)
}

// サービス制御要求検出
func (c *moneyExchange) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

func (c *moneyExchange) recvRequest(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_money_exchange",
	})

	c.logger.Trace("【%v】START:要求受信 request_money_exchange 両替要求", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:要求受信 request_money_exchange 両替要求", texCon.GetUniqueKey())
	err := json.Unmarshal([]byte(message), &c.praReqInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_EXCHANGE_FATAL, "", "入出金管理")
		c.logger.Error("moneyExchange recvRequest json.Unmarshal:%v", err)
		return
	}
	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), c.praReqInfo.RequestInfo.RequestID)

	c.texmyHandler.SetFlagExchange(texCon, true)
	c.texmyHandler.SetExchangeTargetDevice(texCon, c.praReqInfo.TargetDevice)
	c.texmyHandler.SetExchangePattern(texCon, c.praReqInfo.ExchangePattern)
	//処理要求の振り分け
	c.logger.Debug("【%v】‐ StatusMode %v", texCon.GetUniqueKey(), c.praReqInfo.StatusMode)

	switch c.praReqInfo.StatusMode { //動作モード判定
	case domain.MONEY_EXCHANGE_CANCEL: //取消
		c.StatusModeCancel(texCon)

	case domain.MONEY_EXCHANGE_START: //開始
		c.StatusModeStart(texCon)

	case domain.MONEY_EXCHANGE_CONFIRM: //確定
		c.StatusModeConfirm(texCon)

	}

}

// 処理結果応答：両替要求
func (c *moneyExchange) SendResult(texCon *domain.TexContext, recvinfo interface{}) {
	c.logger.Trace("【%v】START: moneyExchange SendResult recvinfo=%+v", texCon.GetUniqueKey(), recvinfo)
	var resInfo domain.ResultMoneyExchange
	//型チェック
	switch receiveData := recvinfo.(type) {
	case domain.ResultInStart: //入金開始要求応答受信
		resInfo.RequestInfo = c.praReqInfo.RequestInfo
		resInfo.Result, resInfo.ErrorCode, resInfo.ErrorDetail = receiveData.Result, receiveData.ErrorCode, receiveData.ErrorDetail
		resInfo.CashControlId = receiveData.CashControlId
		c.cashIdIndata = receiveData.CashControlId //入金時のCashID
		c.texmyHandler.SetExchangeCashControlId(texCon, receiveData.CashControlId)

	case domain.ResultInEnd: //入金終了要求応答受信
		resInfo.RequestInfo = c.praReqInfo.RequestInfo
		resInfo.Result, resInfo.ErrorCode, resInfo.ErrorDetail = receiveData.Result, receiveData.ErrorCode, receiveData.ErrorDetail
		resInfo.CashControlId = receiveData.CashControlId
	case domain.RequestMoneyExchange: //待機中に要求を受け付けた時
		resInfo.RequestInfo = receiveData.RequestInfo
		resInfo.Result = true
	}
	c.logger.Debug("【%v】- resInfo =%+v", texCon.GetUniqueKey(), resInfo)
	payment, err := json.Marshal(resInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_EXCHANGE_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Unmarshal:%v", texCon.GetUniqueKey(), err)
	} else {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_EXCHANGE_SUCCESS, "", "入出金管理")
		topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_money_exchange")
		c.mqtt.Publish(topic, string(payment))
	}
}

// 開始
func (c *moneyExchange) StatusModeStart(texCon *domain.TexContext) {
	c.logger.Trace("【%v】START:moneyExchange StatusModeStart", texCon.GetUniqueKey())
	c.texmyHandler.SetSequence(texCon, domain.REVERSE_EXCHANGEING_CONFIRM_INDATA) //両替中をシーケンスに登録する

	// 入金開始要求リクエスト情報セット
	resInfo := domain.NewRequestInStart(c.texmyHandler.NewRequestInfo(texCon),
		c.praReqInfo.ModeOperation,
		c.praReqInfo.CountClear,
		c.praReqInfo.TargetDevice)
	// 入金開始要求送信
	c.sendRecv.SendRequestInStart(texCon, resInfo)
	c.logger.Trace("【%v】END:moneyExchange StatusModeStart", texCon.GetUniqueKey())
}

// 取消
func (c *moneyExchange) StatusModeCancel(texCon *domain.TexContext) {
	c.texmyHandler.SetSequence(texCon, domain.EXCHANGEING_CANCEL) //両替中をシーケンスに登録する

	// 入金終了要求リクエスト情報セット
	resInfo := domain.NewRequestInEnd(c.texmyHandler.NewRequestInfo(texCon), c.praReqInfo.CashControlId, c.praReqInfo.TargetDevice, c.praReqInfo.StatusMode)
	// 入金終了要求送信
	c.sendRecv.SendRequestInEnd(texCon, resInfo) // 2.1に入金終了要求送信
}

// 確定
func (c *moneyExchange) StatusModeConfirm(texCon *domain.TexContext) {
	c.logger.Trace("【%v】START:moneyExchange StatusModeConfirm", texCon.GetUniqueKey())

	c.logger.Debug("【%v】- c.praReqInfo.ExchangePattern=%v", texCon.GetUniqueKey(), c.praReqInfo.ExchangePattern)
	switch c.praReqInfo.ExchangePattern {
	case domain.REVERSE_EXCHAGE:
		c.texmyHandler.SetSequence(texCon, domain.REVERSE_EXCHANGEING_CONFIRM_INDATA)
		//処理前の有高を格納
		c.aggregateMng.UpdateBeforeCountTbl(texCon, c.maintenanceModeMng.GetMode(texCon), domain.BEFORE_AMOUNT_COUNT_TBL)

	case domain.ONE_CASHTYPE_EXCHAGE,
		domain.ONE_FIVE_CASHTYPE_EXCHAGE:
		c.texmyHandler.SetSequence(texCon, domain.REVERSE_EXCHANGEING_CONFIRM_INDATA)

	case domain.NUMBER_OF_WITHDRAW_DESIGNATED:
		c.texmyHandler.SetSequence(texCon, domain.NUMBER_OF_WITHDRAW_DESIGNATED_CONFIRM)
	}

	// 入金終了要求リクエスト情報セット
	resInfo := domain.NewRequestInEnd(c.texmyHandler.NewRequestInfo(texCon), c.praReqInfo.CashControlId, c.praReqInfo.TargetDevice, c.praReqInfo.StatusMode)
	// 入金終了要求送信
	c.sendRecv.SendRequestInEnd(texCon, resInfo) // 2.1に入金終了要求送信

	c.logger.Trace("【%v】END:moneyExchange StatusModeConfirm", texCon.GetUniqueKey())
}

// 両替入金時：要求送信完了検知
func (c *moneyExchange) SenSorIndataSendFinish(texCon *domain.TexContext, reqType int) {
	c.logger.Trace("【%v】START:moneyExchange SenSorIndataSendFinish reqType=%v Sequence=%v", texCon.GetUniqueKey(), reqType, c.texmyHandler.GetSequence(texCon))
	switch reqType {
	case domain.FINISH_IN_END:
		switch c.texmyHandler.GetSequence(texCon) {
		case domain.NUMBER_OF_WITHDRAW_DESIGNATED_CONFIRM,
			domain.REVERSE_EXCHANGEING_CONFIRM_INDATA,
			domain.ONECASHTYPE_EXCHANGE_CONFIRM,
			domain.FIVEONECASHTYPE_EXCHANGE_CONFIRM,
			domain.NUMBER_OF_WITHDRAW_DESIGNATED:

			reqInfo := c.texmyHandler.RequestReportSafeInfo(texCon) //金庫情報遷移記録に投げる
			c.texdtSendRecv.SendRequestReportSafeInfo(texCon, &reqInfo)
		}

	case domain.FINISH_CHANGE_SUPPLY:
		//逆両替入金時の入金ステータス・両替ステータスを通知する
		c.texmyHandler.SetTexmyNoticeIndata(texCon, true)

		c.texmyHandler.SetTexmyNoticeExchangedata(texCon, true)

		time.Sleep(1 * time.Second)                     // 1秒待つ 下位レイヤーから有高が上がってくるまでの時間が秒ほどある為
		c.texmyHandler.SetTexmyNoticeAmountData(texCon) //ここで502 502に対する有高通知が出るはず

		//ステータス通知
		in := c.noticeMng.GetStatusInData(texCon)
		e := c.noticeMng.GetStatusExchangeData(texCon)
		e.CashControlId = in.CashControlId
		c.noticeMng.UpdateStatusExchangeData(texCon, e)

		var countTbl [16]int
		c.logger.Debug("【%v】- ExchangePattern=%v", texCon.GetUniqueKey(), c.praReqInfo.ExchangePattern)
		switch c.praReqInfo.ExchangePattern {
		case domain.REVERSE_EXCHAGE, //逆両替
			domain.ONE_CASHTYPE_EXCHAGE,      //1系金種両替
			domain.ONE_FIVE_CASHTYPE_EXCHAGE: //1and5系金種両替
			c.texmyHandler.SetSequence(texCon, domain.REVERSE_EXCHANGEING_CONFIRM_OUTDATA)
			copy(countTbl[:], e.ExchangeCountTbl[:])
		case domain.NUMBER_OF_WITHDRAW_DESIGNATED: //出金枚数指定
			c.logger.Debug("【%v】- 紙幣補充0円", texCon.GetUniqueKey())
			if c.CheckRequestOutCashTbl(texCon, c.praReqInfo.PaymentPlanTbl) {
				//notice_exchangeだけ発射して完了させる
				c.texmyHandler.InStatusExchange(texCon)
				return
			}
			copy(countTbl[:], c.praReqInfo.PaymentPlanTbl)
		}

		// 出金開始要求リクエスト情報セット
		reqInfo := domain.NewRequestOutStart(c.texmyHandler.NewRequestInfo(texCon), false, domain.SPECIFIED_NUMBER_OF_WITHDRAWALS, 0, countTbl)
		// 出金開始要求送信
		c.sendRecv.SendRequestOutStart(texCon, &reqInfo)
	}
	c.logger.Trace("【%v】END:moneyExchange SenSorIndataSendFinish", texCon.GetUniqueKey())
}

func (c *moneyExchange) CheckRequestOutCashTbl(texCon *domain.TexContext, paymentPlanTbl []int) bool {
	c.logger.Trace("【%v】START:moneyExchange CheckRequestOutCashTbl paymentPlanTbl=%v", texCon.GetUniqueKey(), paymentPlanTbl)
	var result = true //結果
	for _, paymentPlan := range paymentPlanTbl {
		if paymentPlan > 0 {
			result = false
			break
		}
	}
	c.logger.Trace("【%v】END:moneyExchange CheckRequestOutCashTbl result=%v", texCon.GetUniqueKey(), result)
	return result
}

// 両替出金時：要求送信完了検知
func (c *moneyExchange) SenSorOutdataSendFinish(texCon *domain.TexContext, reqType int) {
	c.logger.Trace("【%v】START:moneyExchange SenSorOutdataSendFinish reqType=%v", texCon.GetUniqueKey(), reqType)
	switch reqType {
	case domain.FINISH_OUT_START: //出金開始要求
		reqInfo := c.texmyHandler.RequestReportSafeInfo(texCon)
		c.texdtSendRecv.SendRequestReportSafeInfo(texCon, &reqInfo) //金庫情報遷移記録に投げる

	case domain.FINISH_REPORT_SAFEINFO: //金庫情報遷移記録要求
		c.SenSorOutdataSendFinish(texCon, domain.FINISH_PRINT_CHANGE_SUPPLY)

	case domain.FINISH_PRINT_CHANGE_SUPPLY:
		switch c.praReqInfo.ExchangePattern {
		case 1, 2: //両替
			// リクエスト生成
			outReq := &domain.RequestGetSalesinfo{
				RequestInfo: c.texmyHandler.NewRequestInfo(texCon),
			}
			// 受信チャネルを生成
			var resChan = make(chan interface{})
			// 外部リクエスト
			go c.paymentSendRecv.SendRequestGetSalesInfo(texCon, resChan, outReq)
			// 外部リクエスト受信
			salesMoney := <-resChan
			// エラーの場合もあるので、型チェックでOKなら次の処理に進む
			outResInfo, ok := salesMoney.(domain.ResultGetSalesinfo)
			if !ok {
				c.logger.Error("【%v】両替出金時 型チェックNG %v", texCon.GetUniqueKey(), reqType)
				return
			}
			reqInfo := c.changeStatus.RequestChangePayment(texCon, 0, outResInfo)
			c.statusTxSendRecv.SendRequestChangePayment(texCon, &reqInfo)
		default:
			reqInfo := c.changeStatus.RequestChangeSupply(texCon, 0)
			c.statusTxSendRecv.SendRequestChangeSupply(texCon, &reqInfo)
		}
	case domain.FINISH_CHANGE_SUPPLY:
		//両替時の出金ステータス・両替ステータスを通知する
		c.texmyHandler.SetTexmyNoticeOutdata(texCon, true)
		c.texmyHandler.SetTexmyNoticeExchangedata(texCon, true)

		time.Sleep(2 * time.Second)                     // 2秒待つ 下位レイヤーから有高が上がってくるまでの時間が2秒ほどある為
		c.texmyHandler.SetTexmyNoticeAmountData(texCon) //ここで502 502に対する有高通知が出るはず
	}
	c.logger.Trace("【%v】END:moneyExchange SenSorOutdataSendFinish", texCon.GetUniqueKey())
}
