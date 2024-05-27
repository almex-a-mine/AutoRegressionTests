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

type setAmount struct {
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
	reqInfo             domain.RequestSetAmount
	resInfo             domain.ResultSetAmount
	printDataManager    usecases.PrintDataManager
	reqAmountStatusInfo domain.RequestAmountStatus
	countTbl            [domain.EXTRA_CASH_TYPE_SHITEI]int
	coincassetFlag      bool //コインカセット要求経由　true:経由している　false：経由していない
	changeStatus        usecases.ChangeStatusRepository
	noticeManager       usecases.TexMoneyNoticeManagerRepository
	safeInfoMng         usecases.SafeInfoManager
	skipFlag            bool
}

// 現在枚数変更要求
func NewRequestSetAmount(mqtt handler.MqttRepository,
	logger handler.LoggerRepository,
	config config.Configuration,
	syslogMng usecases.SyslogManager,
	errorMng usecases.ErrorManager,
	sendRecv SendRecvRepository,
	texdtSendRecv TexdtSendRecvRepository,
	statusTxSendRecv StatusTxSendRecvRepository,
	printSendRecv PrintSendRecvRepository,
	texmyHandler usecases.TexMoneyHandlerRepository,
	printDataManager usecases.PrintDataManager,
	changeStatus usecases.ChangeStatusRepository,
	noticeManager usecases.TexMoneyNoticeManagerRepository,
	safeInfoMng usecases.SafeInfoManager) SetAmountRepository {
	return &setAmount{
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
		reqInfo:             domain.RequestSetAmount{},
		resInfo:             domain.ResultSetAmount{},
		printDataManager:    printDataManager,
		reqAmountStatusInfo: domain.RequestAmountStatus{},
		countTbl:            [domain.EXTRA_CASH_TYPE_SHITEI]int{},
		changeStatus:        changeStatus,
		noticeManager:       noticeManager,
		safeInfoMng:         safeInfoMng,
		skipFlag:            false,
	}
}

// 開始処理
func (c *setAmount) Start() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_set_amount")
	c.mqtt.Subscribe(topic, c.recvRequest)
}

// 停止処理
func (c *setAmount) Stop() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_set_amount")
	c.mqtt.Unsubscribe(topic)
}

// サービス制御要求検出
func (c *setAmount) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

func (c *setAmount) recvRequest(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_set_amount",
	})

	c.logger.Trace("【%v】START:要求受信 request_set_amount 現在枚数変更要求", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:要求受信 request_set_amount 現在枚数変更要求", texCon.GetUniqueKey())

	if c.coincassetFlag {
		// 1秒待ってから再開
		c.logger.Debug("setAmount recvRequest Cassette Control Wait")
		time.Sleep(time.Second * 1)
		// コインカセット操作フラグをOFFに変更
		c.coincassetFlag = false
	}

	var reqInfo domain.RequestSetAmount

	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_SET_AMOUNT_FATAL, "", "入出金管理")
		c.logger.Error("setAmount recvRequest json.Unmarshal:%v", err)
		return
	}

	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), reqInfo.RequestInfo.RequestID)

	c.reqInfo = reqInfo
	c.resInfo.RequestInfo = reqInfo.RequestInfo
	c.resInfo.Result = true
	c.resInfo.CashControlId = reqInfo.CashControlId
	c.texmyHandler.SetSequence(texCon, domain.SET_AMOUNT)

	c.CheckOperationMode(texCon)

}

// 操作モードチェック
func (c *setAmount) CheckOperationMode(texCon *domain.TexContext) {
	//操作モードチェック
	switch c.reqInfo.OperationMode {

	case domain.MONEY_SETAMOUNT_ANALOG_COLLECT: // 手動補充/回収
		c.ModeAnalogCollect(texCon, c.reqInfo)

	case domain.MONEY_SETAMOUNT_REJECTBOX_COLLECT: // リジェクトボックス回収
		c.ModeRejectBoxCollect(texCon, c.reqInfo)

	case domain.MONEY_SETAMOUNT_UNRETURNED_COLLECT: // 回収庫回収
		c.ModeSetamountUnreturnedCollect(texCon, c.reqInfo)

	case domain.MONEY_SETAMOUNT_UNRETURNED_AND_SALES_COLLECT: // 非還流庫回収and売上金回収
		c.ModeSetAmountUnreturnedAndSalesCollect(texCon, c.reqInfo)
	}

}

// 結果応答
func (c *setAmount) SendResult(texCon *domain.TexContext, res domain.ResultCashctlSetAmount) {
	c.logger.Trace("【%v】START:setAmount Result", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:setAmount Result", texCon.GetUniqueKey())
	if c.coincassetFlag {
		return
	}

	c.resInfo.CashControlId = res.CashControlId

	amount, err := json.Marshal(c.resInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_SET_AMOUNT_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Marshal:%v", texCon.GetUniqueKey(), err)
		return
	}
	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_SET_AMOUNT_SUCCESS, "", "入出金管理")

	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_set_amount")
	c.mqtt.Publish(topic, string(amount))

}

// 手動補充回収
func (c *setAmount) ModeAnalogCollect(texCon *domain.TexContext, reqInfo domain.RequestSetAmount) {
	c.logger.Trace("【%v】START:現在枚数変更要求 手動補充回収", texCon.GetUniqueKey())

	c.texmyHandler.SetSequence(texCon, domain.MANUAL_REPLENISHMENT_COLLECTION)
	preqInfo := c.texmyHandler.SetAmountRequestCreate(texCon, &reqInfo)
	c.sendRecv.SendRequestCashctlSetAmount(texCon, &preqInfo, domain.MANUAL_REPLENISHMENT_COLLECTION)
	c.logger.Trace("【%v】END:現在枚数変更要求 手動補充回収", texCon.GetUniqueKey())
}

// リジェクトボックス回収
func (c *setAmount) ModeRejectBoxCollect(texCon *domain.TexContext, reqInfo domain.RequestSetAmount) {
	c.logger.Trace("【%v】START:現在枚数変更要求 リジェクトボックス回収", texCon.GetUniqueKey())

	c.texmyHandler.SetSequence(texCon, domain.REJECTBOXCOLLECT_START)
	preqInfo := c.texmyHandler.SetAmountRequestCreate(texCon, &reqInfo)
	c.sendRecv.SendRequestCashctlSetAmount(texCon, &preqInfo, domain.REJECTBOXCOLLECT_START)
	c.logger.Trace("【%v】END:現在枚数変更要求 リジェクトボックス回収", texCon.GetUniqueKey())
}

// 非還流庫回収
func (c *setAmount) ModeSetamountUnreturnedCollect(texCon *domain.TexContext, reqInfo domain.RequestSetAmount) {
	c.logger.Trace("【%v】START:現在枚数変更要求 非還流庫回収", texCon.GetUniqueKey())

	c.texmyHandler.SetSequence(texCon, domain.UNRETURNEDCOLLECT_START)
	preqInfo := c.texmyHandler.SetAmountRequestCreate(texCon, &reqInfo)
	c.sendRecv.SendRequestCashctlSetAmount(texCon, &preqInfo, domain.UNRETURNEDCOLLECT_START)
	c.logger.Trace("【%v】END:現在枚数変更要求 	", texCon.GetUniqueKey())
}

// 非還流庫回収and売上金回収含
func (c *setAmount) ModeSetAmountUnreturnedAndSalesCollect(texCon *domain.TexContext, reqInfo domain.RequestSetAmount) {
	c.logger.Trace("【%v】START:現在枚数変更要求 非還流庫回収nad売上金回収含", texCon.GetUniqueKey())
	c.texmyHandler.SetSequence(texCon, domain.UNRETURNEDCOLLECT_START)
	preqInfo := c.texmyHandler.UnreturnedAndSalesCollect(texCon, &reqInfo)
	c.sendRecv.SendRequestCashctlSetAmount(texCon, &preqInfo, domain.UNRETURNEDCOLLECT_START)

	c.logger.Trace("【%v】END:現在枚数変更要求 非還流庫回収nad売上金回収含", texCon.GetUniqueKey())

}

// 硬貨カセット操作要求との土管 pram1:変更後有高格納
func (c *setAmount) ConnetctCoincasseteControl(texCon *domain.TexContext, cashTbl [26]int) {
	c.coincassetFlag = true //カセット経由の場合はset_amountのresultを返さないようにする為のフラグセット

	var reqInfo domain.RequestSetAmount
	reqInfo.CashTbl = cashTbl
	c.ModeAnalogCollect(texCon, reqInfo) //変更後有高で有高の手動補充に入れる
}

// 各要求送信完了検知
func (c *setAmount) SenSorSendFinish(texCon *domain.TexContext, reqType int) {
	c.logger.Trace("【%v】START:各要求送信完了検知 SenSorSendFinish", texCon.GetUniqueKey())
	switch reqType {
	case domain.FINISH_IN_END: //有高枚数変更要求完了
		c.logger.Debug("【%v】- 有高枚数変更要求完了", texCon.GetUniqueKey())
		reqInfo := c.texmyHandler.RequestReportSafeInfo(texCon) //確定：稼働データ管理に金庫状態記録を投げる
		c.texdtSendRecv.SendRequestReportSafeInfo(texCon, &reqInfo)

	case domain.FINISH_REPORT_SAFEINFO: //金庫情報遷移記録完了
		c.logger.Debug("【%v】- 金庫情報遷移記録完了", texCon.GetUniqueKey())
		go c.SenSorSendFinish(texCon, domain.FINISH_PRINT_CHANGE_SUPPLY)

	case domain.FINISH_PRINT_CHANGE_SUPPLY: //印刷要求完了
		c.logger.Debug("【%v】- 印刷要求完了", texCon.GetUniqueKey())
		reqInfo := c.changeStatus.RequestChangeSupply(texCon, 0)
		c.statusTxSendRecv.SendRequestChangeSupply(texCon, &reqInfo)

		if c.coincassetFlag { //カセット経由でのsetamountが完了したらfalseに戻す
			c.coincassetFlag = false
		}

	case domain.FINISH_CHANGE_SUPPLY: //精算機状態管理要求完了
		c.logger.Debug("【%v】- 精算機状態管理要求完了 skipFlag =%t", texCon.GetUniqueKey(), c.skipFlag)
		if c.skipFlag { //現在枚数変更要求送信スキップフラグが立っている場合のみここで通知を送る
			c.texmyHandler.SetTexmyNoticeAmountData(texCon) //通知送信
			c.skipFlag = false
		}

		time.Sleep(2 * time.Second) // 2秒待つ 下位レイヤーから有高が上がってくるまでの時間が2秒ほどある為
		c.texmyHandler.SetTexmyNoticeAmountData(texCon)

	case domain.SKIP_SEND_SET_AMOUNT: //現在枚数変更要求送信スキップ
		c.logger.Debug("【%v】- 現在枚数変更要求送信スキップ", texCon.GetUniqueKey())
		// 応答送信
		c.SendResult(texCon, domain.ResultCashctlSetAmount{
			RequestInfo:   c.reqInfo.RequestInfo,
			CashControlId: domain.SET_AMOUNT_ONE,
			Result:        true})
		// 有高データ通知送信
		cashAvailable := c.safeInfoMng.GetLogicalCashAvailable(texCon)
		statusAmount := domain.StatusAmount{
			Amount:     cashAvailable.Amount,
			CountTbl:   cashAvailable.CountTbl,
			ExCountTbl: cashAvailable.ExCountTbl,
		}
		c.noticeManager.UpdateStatusAmountData(texCon, statusAmount)
		c.SenSorSendFinish(texCon, domain.FINISH_IN_END) //有高枚数変更要求完了
		c.skipFlag = true

	default:
		c.logger.Warn("【%v】- 対象シーケンス無し", texCon.GetUniqueKey())

	}
	c.logger.Trace("【%v】END:各要求送信完了検知 SenSorSendFinish", texCon.GetUniqueKey())
}
