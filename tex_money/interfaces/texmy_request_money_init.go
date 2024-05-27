package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases"
	"time"
)

type moneyInit struct {
	mqtt                  handler.MqttRepository
	logger                handler.LoggerRepository
	syslogMng             usecases.SyslogManager
	sendRecv              SendRecvRepository
	texdtSendRecv         TexdtSendRecvRepository
	statusTxSendRecv      StatusTxSendRecvRepository
	texmyHandler          usecases.TexMoneyHandlerRepository
	safeInfoMng           usecases.SafeInfoManager
	praReqInfo            domain.RequestMoneyInit
	changeStatus          usecases.ChangeStatusRepository
	texMoneyNoticeManager usecases.TexMoneyNoticeManagerRepository
}

// 初期補充要求
func NewRequestMoneyInit(mqtt handler.MqttRepository,
	logger handler.LoggerRepository,
	syslogMng usecases.SyslogManager,
	sendRecv SendRecvRepository,
	texdtSendRecv TexdtSendRecvRepository,
	statusTxSendRecv StatusTxSendRecvRepository,
	texmyHandler usecases.TexMoneyHandlerRepository,
	safeInfoMng usecases.SafeInfoManager,
	changeStatus usecases.ChangeStatusRepository,
	texMoneyNoticeManager usecases.TexMoneyNoticeManagerRepository) MoneyInitRepository {
	return &moneyInit{
		mqtt:                  mqtt,
		logger:                logger,
		syslogMng:             syslogMng,
		sendRecv:              sendRecv,
		texdtSendRecv:         texdtSendRecv,
		statusTxSendRecv:      statusTxSendRecv,
		texmyHandler:          texmyHandler,
		safeInfoMng:           safeInfoMng,
		changeStatus:          changeStatus,
		texMoneyNoticeManager: texMoneyNoticeManager,
	}
}

// 開始処理
func (c *moneyInit) Start() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_money_init")
	c.mqtt.Subscribe(topic, c.recvRequest)
}

// 停止処理
func (c *moneyInit) Stop() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_money_inits")
	c.mqtt.Unsubscribe(topic)
}

// サービス制御要求検出
func (c *moneyInit) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

func (c *moneyInit) recvRequest(message string) {

	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_money_init",
	})

	c.logger.Trace("【%v】START:要求受信 request_money_init 初期補充要求", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:要求受信 request_money_init 初期補充要求", texCon.GetUniqueKey())

	var reqInfo domain.RequestMoneyInit
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.logger.Error("moneyInit recvRequest json.Unmarshal:%v", err)
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_INIT_FATAL, "", "入出金管理")
		return
	}
	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), reqInfo.RequestInfo.RequestID)

	c.praReqInfo = reqInfo //リクエスト情報のセット

	//動作モード判定
	c.logger.Debug("【%v】- 初期補充要求  動作モード=%v", texCon.GetUniqueKey(), reqInfo.StatusMode)
	switch reqInfo.StatusMode {
	case 1: //開始
		c.statusModeStart(texCon, &reqInfo)
	case 0: //取消
		c.statusMoneyInitCancel(texCon, reqInfo)
	case 2: //確定
		c.statusMoneyInitConfirm(texCon, reqInfo)
	case 3: //更新
		c.statusModeUpdate(texCon, &reqInfo)
	}
}

// 開始
func (c *moneyInit) statusModeStart(texCon *domain.TexContext, pReqInfo *domain.RequestMoneyInit) {
	c.texmyHandler.SetSequence(texCon, domain.INITIAL_ADDING_START)

	// 入金枚数クリア無しの場合
	// 前回確定分にプラスした入金枚数が下位から再度通知される。
	// 前回確定分は既に論理枚数に加算してしまっている為
	// 処理前にnotice_indataに保持している枚数分の入金情報を一旦論理有高からマイナスを行う
	if !pReqInfo.CountClear {
		inData := c.texMoneyNoticeManager.GetStatusInData(texCon)
		c.safeInfoMng.UpdateOutLogicalCashAvailable(texCon, domain.SortInfoTbl{
			SortType:   domain.CASH_AVAILABLE,
			Amount:     inData.Amount,
			CountTbl:   inData.CountTbl,
			ExCountTbl: inData.ExCountTbl,
		})
	}

	// 入金開始要求リクエスト情報セット
	resInfo := domain.NewRequestInStart(c.texmyHandler.NewRequestInfo(texCon),
		pReqInfo.ModeOperation,
		pReqInfo.CountClear,
		pReqInfo.TargetDevice)
	// 入金開始要求送信
	c.sendRecv.SendRequestInStart(texCon, resInfo)
}

// 取消
func (c *moneyInit) statusMoneyInitCancel(texCon *domain.TexContext, pReqInfo domain.RequestMoneyInit) {
	c.texmyHandler.SetSequence(texCon, domain.INITIAL_ADDING_CANCEL)

	// 1.3応答送信
	c.SendResult(texCon, pReqInfo.CashControlId)

	// 入金終了要求リクエスト情報セット
	resInfo := domain.NewRequestInEnd(c.texmyHandler.NewRequestInfo(texCon), pReqInfo.CashControlId, pReqInfo.TargetDevice, pReqInfo.StatusMode)
	// 入金終了要求送信
	c.sendRecv.SendRequestInEnd(texCon, resInfo)
}

// 確定
func (c *moneyInit) statusMoneyInitConfirm(texCon *domain.TexContext, pReqInfo domain.RequestMoneyInit) {
	c.texmyHandler.SetSequence(texCon, domain.INITIAL_ADDING_CONFIRM)

	// 1.3応答送信
	c.SendResult(texCon, pReqInfo.CashControlId)

	// 入金終了要求リクエスト情報セット
	resInfo := domain.NewRequestInEnd(c.texmyHandler.NewRequestInfo(texCon), pReqInfo.CashControlId, pReqInfo.TargetDevice, pReqInfo.StatusMode)
	// 入金終了要求送信
	c.sendRecv.SendRequestInEnd(texCon, resInfo)
}

// 更新
func (c *moneyInit) statusModeUpdate(texCon *domain.TexContext, reqInfo *domain.RequestMoneyInit) {
	c.texmyHandler.SetSequence(texCon, domain.INITIAL_ADDING_UPDATE)

	// 1.3応答送信
	c.SendResult(texCon, texCon.GetUniqueKey()) //cashControlIdは何でもよいためユニークな値を詰めておく

	// 現在枚数を初期補充枚数にコピーする
	// 現在枚数を取得
	_, initialSortInfo := c.safeInfoMng.GetSortInfo(texCon, domain.CASH_AVAILABLE)
	// 登録時のSortTypeを変更
	initialSortInfo.SortType = domain.INITIAL_REPLENISHMENT
	// 登録
	c.safeInfoMng.UpdateSortInfo(texCon, initialSortInfo)

	//稼働データ管理に金庫状態記録を投げる
	c.SenSorSendFinish(texCon, domain.FINISH_IN_END)
}

// 処理結果応答
func (c *moneyInit) SendResult(texCon *domain.TexContext, cashControlId string) bool {
	c.logger.Trace("【%v】START: moneyInit SendResult", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END: moneyInit SendResult", texCon.GetUniqueKey())

	//初期補充応答
	result, errorCode, errorDetail := c.texmyHandler.GetErrorFromRequest(texCon)

	res := domain.ResultMoneyInit{
		RequestInfo:   c.praReqInfo.RequestInfo,
		Result:        result,
		ErrorCode:     errorCode,
		ErrorDetail:   errorDetail,
		CashControlId: cashControlId,
	}

	payment, err := json.Marshal(res)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_INIT_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Marshal:%v", texCon.GetUniqueKey(), err)
		return false
	}

	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_INIT_SUCCESS, "", "入出金管理")
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_money_init")
	c.mqtt.Publish(topic, string(payment))
	return true
}

// 各要求送信完了検知
func (c *moneyInit) SenSorSendFinish(texCon *domain.TexContext, reqType int) {
	c.logger.Trace("【%v】START:moneyInit SenSorSendFinish Type=%v", texCon.GetUniqueKey(), reqType)
	switch reqType {
	case domain.FINISH_IN_END: //入金禁止要求完了
		//稼働データ管理に金庫状態記録を投げる
		reqInfo := c.texmyHandler.RequestReportSafeInfo(texCon)
		c.texdtSendRecv.SendRequestReportSafeInfo(texCon, &reqInfo)
	case domain.FINISH_REPORT_SAFEINFO: //金庫情報遷移記録完了
		// リクエスト情報セット
		reqInfo := c.changeStatus.RequestChangeSupply(texCon, 0)
		c.statusTxSendRecv.SendRequestChangeSupply(texCon, &reqInfo)
	case domain.FINISH_CHANGE_SUPPLY: //状態変更完了
		//完了通知
		c.texmyHandler.SetTexmyNoticeIndata(texCon, true)

		time.Sleep(2 * time.Second) // 2秒待つ 下位レイヤーから有高が上がってくるまでの時間が2秒ほどある為
		c.texmyHandler.SetTexmyNoticeAmountData(texCon)
	}
	c.logger.Trace("【%v】END:moneyInit SenSorSendFinish", texCon.GetUniqueKey())
}
