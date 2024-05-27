package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/pkg/calculation"
	"tex_money/usecases"
	"time"
)

type moneyCollect struct {
	mqtt                          handler.MqttRepository
	logger                        handler.LoggerRepository
	syslogMng                     usecases.SyslogManager
	errorMng                      usecases.ErrorManager
	sendRecv                      SendRecvRepository
	texdtSendRecv                 TexdtSendRecvRepository
	statusTxSendRecv              StatusTxSendRecvRepository
	texmyHandler                  usecases.TexMoneyHandlerRepository
	praReqInfo                    domain.RequestMoneyCollect
	reqAmountStatusInfo           domain.RequestAmountStatus
	countTbl                      [domain.EXTRA_CASH_TYPE_SHITEI]int
	safeInfoMng                   usecases.SafeInfoManager
	gatherSalesInOverflowExTbl    [domain.EXTRA_CASH_TYPE_SHITEI]int
	gatherSalesInOverflowCountTbl [domain.CASH_TYPE_SHITEI]int
	gatherSalesInOverflowAmount   int
	changeStatus                  usecases.ChangeStatusRepository
}

// 回収要求(途中回収要求／全回収要求／売上金回収要求)
func NewRequestMoneyCollect(mqtt handler.MqttRepository, logger handler.LoggerRepository, syslogMng usecases.SyslogManager, errorMng usecases.ErrorManager, sendRecv SendRecvRepository, texdtSendRecv TexdtSendRecvRepository, statusTxSendRecv StatusTxSendRecvRepository, texmyHandler usecases.TexMoneyHandlerRepository, safeInfoMng usecases.SafeInfoManager, changeStatus usecases.ChangeStatusRepository) MoneyCollectRepository {
	return &moneyCollect{
		mqtt:                mqtt,
		logger:              logger,
		syslogMng:           syslogMng,
		errorMng:            errorMng,
		sendRecv:            sendRecv,
		texdtSendRecv:       texdtSendRecv,
		statusTxSendRecv:    statusTxSendRecv,
		texmyHandler:        texmyHandler,
		reqAmountStatusInfo: domain.RequestAmountStatus{},
		countTbl:            [domain.EXTRA_CASH_TYPE_SHITEI]int{},
		safeInfoMng:         safeInfoMng,
		changeStatus:        changeStatus,
	}

}

// 開始処理
func (c *moneyCollect) Start() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_money_collect")
	c.mqtt.Subscribe(topic, c.recvRequest)
}

// 停止処理
func (c *moneyCollect) Stop() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_money_collect")
	c.mqtt.Unsubscribe(topic)
}

// サービス制御要求検出
func (c *moneyCollect) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

func (c *moneyCollect) recvRequest(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_money_collect",
	})

	c.logger.Trace("【%v】START:要求受信 request_money_collect 回収要求", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:要求受信 request_money_collect 回収要求", texCon.GetUniqueKey())
	err := json.Unmarshal([]byte(message), &c.praReqInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_COLLECT_FATAL, "", "入出金管理")
		c.logger.Error("moneyCollect recvRequest json.Unmarshal:%v", err)
		return
	}

	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), c.praReqInfo.RequestInfo.RequestID)

	// 売上金回収以外で、中身がcashTblの中身が0の場合の処理
	// 売上金回収の場合には、売上金額(salesAmount)もチェックする為
	if c.checkCollectPattern(c.praReqInfo.CollectMode, c.praReqInfo.SalesAmount, c.praReqInfo.CashTbl) {
		c.result(texCon)
		c.texmyHandler.SensorZeroNoticeCollect(texCon)
		return
	}
	c.texmyHandler.SetFlagCollect(texCon, true) //フラグのセット

	c.sendRecv.CashId(texCon, c.praReqInfo.CashControlId)

	c.CheckCollectMode(texCon)

}

// checkCollectPattern 回収パターン毎の値セットチェック
func (c *moneyCollect) checkCollectPattern(collectMode int, salesAmount int, cashTbl [16]int) bool {
	switch collectMode {

	case domain.MONEY_COLLECT_SALESMONEY:
		ok := c.CheckCountTblOfReqMoneyCollect(cashTbl)
		if ok && salesAmount == 0 {
			return true
		}
		return false
	case domain.MONEY_COLLECT_MIDDLE_AND_SALES:
		ok := c.CheckCountTblOfReqMoneyCollect(cashTbl)
		if ok && salesAmount == 0 {
			return true
		}
		return false

	default:
		return c.CheckCountTblOfReqMoneyCollect(cashTbl)
	}

}

// チェック配列に値が格納されているか
func (c *moneyCollect) CheckCountTblOfReqMoneyCollect(countTbl [16]int) bool {

	var checkCountTbl [16]int

	// 回収対象がマイナスの場合には、0に置き換える。
	for i, v := range countTbl {
		checkCountTbl[i] = v
		if v < 0 {
			checkCountTbl[i] = 0
		}
	}

	// 回収対象が全て0かどうかのチェック
	for _, value := range checkCountTbl {
		if value != 0 {
			return false
		}
	}
	return true
}

func (c *moneyCollect) result(texCon *domain.TexContext) {
	var reqInfo domain.ResultOutStart
	reqInfo.RequestInfo = c.praReqInfo.RequestInfo
	reqInfo.Result = true
	c.SendResult(texCon, reqInfo)

}

// 回収モード判定
func (c *moneyCollect) CheckCollectMode(texCon *domain.TexContext) {
	c.logger.Trace("【%v】START:moneyCollect CheckCollectMode", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:moneyCollect CheckCollectMode()", texCon.GetUniqueKey())

	switch c.praReqInfo.CollectMode {
	case domain.MONEY_COLLECT_MIDDLE: //途中回収
		c.StatusModeMiddle(texCon, &c.praReqInfo)

	case domain.MONEY_COLLECT_ALL, //全回収
		domain.MONEY_COLLECT_INREJECT: // 全回収(リジェクト庫含)
		c.StatusModeAll(texCon, &c.praReqInfo)

	case domain.MONEY_COLLECT_SALESMONEY: //売上金回収
		c.StatusModeSalesMoney(texCon, &c.praReqInfo)

	case domain.MONEY_COLLECT_MIDDLE_AND_SALES: // 途中回収(売上金回収を含)
		c.StatusModeMiddleAndSales(texCon, &c.praReqInfo)
	}
}

// 処理結果応答:出金開始要求
func (c *moneyCollect) SendResult(texCon *domain.TexContext, reqInfo domain.ResultOutStart) bool {
	c.logger.Trace("【%v】START:moneyCollect SendResult", texCon.GetUniqueKey())

	res := domain.ResultMoneyCollect{
		RequestInfo:   c.praReqInfo.RequestInfo,
		Result:        reqInfo.Result,
		ErrorCode:     reqInfo.ErrorCode,
		ErrorDetail:   reqInfo.ErrorDetail,
		CashControlId: reqInfo.CashControlId,
	}

	payment, err := json.Marshal(res)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_COLLECT_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Marshal:%v", texCon.GetUniqueKey(), err)
	} else {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_COLLECT_SUCCESS, "", "入出金管理")
		topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_money_collect")
		c.mqtt.Publish(topic, string(payment))
	}

	c.logger.Trace("【%v】END: moneyCollect SendResult", texCon.GetUniqueKey())
	return true
}

// 処理結果応答:出金停止要求
func (c *moneyCollect) SendResultCollectStart(texCon *domain.TexContext, reqInfo domain.ResultOutStop) bool {
	c.logger.Trace("【%v】START: moneyCollect SendResult", texCon.GetUniqueKey())

	res := domain.ResultMoneyCollect{
		RequestInfo:   c.praReqInfo.RequestInfo,
		Result:        reqInfo.Result,
		ErrorCode:     reqInfo.ErrorCode,
		ErrorDetail:   reqInfo.ErrorDetail,
		CashControlId: "",
	}

	payment, err := json.Marshal(res)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_COLLECT_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Marshal:%v", texCon.GetUniqueKey(), err)
	} else {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_COLLECT_SUCCESS, "", "入出金管理")
		topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_money_collect")
		c.mqtt.Publish(topic, string(payment))
	}
	c.logger.Trace("【%v】END: moneyCollect SendResult", texCon.GetUniqueKey())
	return true
}

// 処理結果応答:回収開始要求
func (c *moneyCollect) SendResultOutStop(texCon *domain.TexContext, reqInfo domain.ResultCollectStart) bool {
	c.logger.Trace("【%v】START: moneyCollect SendResultOutStop", texCon.GetUniqueKey())
	//回収応答
	res := domain.ResultMoneyCollect{
		RequestInfo:   c.praReqInfo.RequestInfo,
		Result:        reqInfo.Result,
		ErrorCode:     reqInfo.ErrorCode,
		ErrorDetail:   reqInfo.ErrorDetail,
		CashControlId: reqInfo.CashControlId,
	}

	c.logger.Debug("【%v】- resInfo=%v", texCon.GetUniqueKey(), res)
	payment, err := json.Marshal(res)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_COLLECT_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Marshal:%v", texCon.GetUniqueKey(), err)
	} else {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_COLLECT_SUCCESS, "", "入出金管理")
		topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_money_collect")
		c.mqtt.Publish(topic, string(payment))
	}
	c.logger.Trace("【%v】END: moneyCollect SendResultOutStop", texCon.GetUniqueKey())
	return true
}

// 処理結果応答:回収停止要求
func (c *moneyCollect) SendResultCollectStop(texCon *domain.TexContext, reqInfo domain.ResultCollectStop) bool {
	c.logger.Trace("【%v】START: moneyCollect SendResultCollectStop", texCon.GetUniqueKey())
	//回収応答
	res := domain.ResultMoneyCollect{
		RequestInfo:   c.praReqInfo.RequestInfo,
		Result:        reqInfo.Result,
		ErrorCode:     reqInfo.ErrorCode,
		ErrorDetail:   reqInfo.ErrorDetail,
		CashControlId: "",
	}

	payment, err := json.Marshal(res)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_COLLECT_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Marshal:%v", texCon.GetUniqueKey(), err)
	} else {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MONEY_COLLECT_SUCCESS, "", "入出金管理")
		topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_money_collect")
		c.mqtt.Publish(topic, string(payment))
	}
	c.logger.Trace("【%v】END: moneyCollect SendResultCollectStop", texCon.GetUniqueKey())
	return true
}

// 途中回収要求
func (c *moneyCollect) StatusModeMiddle(texCon *domain.TexContext, pReqInfo *domain.RequestMoneyCollect) {
	c.logger.Trace("【%v】START:moneyCollect StatusModeMiddle(pReqInfo=%+v)", texCon.GetUniqueKey(), pReqInfo)
	preqInfo, preqInfo2, preqInfo3, preqInfo4 := c.texmyHandler.Collect(texCon, pReqInfo)
	c.logger.Debug("【%v】- pReqInfo.OutType=%v, preqInfo=%+v, preqInfo2=%+v, preqInfo3=%+v, preqInfo4 =%+v", texCon.GetUniqueKey(), pReqInfo.OutType, preqInfo, preqInfo2, preqInfo3, preqInfo4)
	switch pReqInfo.OutType {
	case domain.WITHDRAW_TO_OUTLET: //出金要求
		c.logger.Debug("【%v】- pReqInfo.StatusMode=%v", texCon.GetUniqueKey(), pReqInfo.StatusMode)
		if pReqInfo.StatusMode == domain.START {
			c.texmyHandler.SetSequence(texCon, domain.MIDDLE_START_OUT_START)
			c.sendRecv.SendRequestOutStart(texCon, &preqInfo3) //出金要求
		} else {
			c.texmyHandler.SetSequence(texCon, domain.MIDDLE_START_OUT_STOP)
			c.sendRecv.SendRequestCollectStart(texCon, &preqInfo4) //出金停止要求
		}

	case domain.COLLECT_TO_COLLECTION_BOX: //回収要求
		if pReqInfo.StatusMode == domain.START {
			c.texmyHandler.SetSequence(texCon, domain.MIDDLE_START_COLLECT_START)
			c.sendRecv.SendRequestOutStop(texCon, &preqInfo) //回収要求
		} else {
			c.texmyHandler.SetSequence(texCon, domain.MIDDLE_START_COLLECT_STOP)
			c.sendRecv.SendRequestCollectStop(texCon, &preqInfo2) //回収停止要求
		}
	}
	c.logger.Trace("【%v】END:moneyCollect StatusModeMiddle", texCon.GetUniqueKey())
}

// 売上金回収
func (c *moneyCollect) StatusModeSalesMoney(texCon *domain.TexContext, pReqInfo *domain.RequestMoneyCollect) {
	c.logger.Trace("【%v】START:moneyCollect StatusModeSalesMoney", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:moneyCollect StatusModeSalesMoney", texCon.GetUniqueKey())

	_, _, preqInfo, preqInfo2 := c.texmyHandler.Collect(texCon, pReqInfo)
	c.texmyHandler.SetSequence(texCon, domain.SALESMONEY_START)

	// okの場合、回収対象無し（あふれからの回収分のみ）
	ok := c.CheckCountTblOfReqMoneyCollect(pReqInfo.CashTbl)

	switch pReqInfo.StatusMode {
	case domain.START:
		// あふれ金庫からの回収分情報をセット
		c.gatherSalesInOverflowCountTbl, c.gatherSalesInOverflowExTbl, c.gatherSalesInOverflowAmount = c.gatherSalesOverflow(texCon, pReqInfo.SalesAmount, pReqInfo.CashTbl)
		// メソッドからあふれ金庫回収分をセット
		if c.gatherSalesInOverflowAmount != 0 {

			c.texmyHandler.SetOverflowCollectSales(texCon, ok, c.gatherSalesInOverflowAmount, c.gatherSalesInOverflowCountTbl, c.gatherSalesInOverflowExTbl)
			// ok の場合、cashTbl=払出枚数が0であふれからの回収のみとなるので、
			// result、notice_collectを送信して終了
			if ok {
				c.result(texCon)
				c.texmyHandler.SensorOverflowOnlyNoticeCollect(texCon, c.gatherSalesInOverflowAmount, c.gatherSalesInOverflowCountTbl, c.gatherSalesInOverflowExTbl)
				c.clearGatherSalesInOverflow()
				return
			}

			// メソッドの中身を初期化
			c.clearGatherSalesInOverflow()

		}

		c.sendRecv.SendRequestOutStart(texCon, &preqInfo) //出金開始要求

	case domain.STOP:

		// 出金停止要求
		c.sendRecv.SendRequestCollectStart(texCon, &preqInfo2)
	}

}

func (c *moneyCollect) clearGatherSalesInOverflow() {
	// メソッドの中身を初期化
	c.gatherSalesInOverflowExTbl = [domain.EXTRA_CASH_TYPE_SHITEI]int{}
	c.gatherSalesInOverflowAmount = 0
	c.gatherSalesInOverflowCountTbl = [domain.CASH_TYPE_SHITEI]int{}
}

func (c *moneyCollect) gatherSalesOverflow(texCon *domain.TexContext, reqSalesAmount int, reqCashTbl [16]int) ([domain.CASH_TYPE_SHITEI]int, [domain.EXTRA_CASH_TYPE_SHITEI]int, int) {
	// あふれ金庫回収有無判定
	// 16桁→26桁
	var cashToExCashTbl [26]int
	copy(cashToExCashTbl[:], reqCashTbl[:])

	// 回収要求分のcashTblをの金額を算出
	cashTblTotalAmount := calculation.NewCassette(cashToExCashTbl).GetTotalAmount()

	// 差分抽出 (連携されてくる金額はあふれ金庫を含む、連携されてくる金種テーブルはあふれ金庫を含まない)
	overflowAmount := reqSalesAmount - cashTblTotalAmount

	// 初期化漏れを防ぐ為、開始時には毎度データを新規でセットする
	// 有高取得
	_, safe0 := c.safeInfoMng.GetSortInfo(texCon, 0)
	// 有高をベースに、オーバーフロー分の金額でオーバーフローのみの逆両替を実行
	ex := calculation.NewCassette(safe0.ExCountTbl).OverflowOnlyExchange(overflowAmount)

	// 有高情報生成 10桁
	countTbl := [domain.CASH_TYPE_SHITEI]int{ex[16], ex[17], ex[18], ex[19], ex[20], ex[21], ex[22], ex[23], ex[24], ex[25]}

	return countTbl, ex, overflowAmount
}

// 全回収
func (c *moneyCollect) StatusModeAll(texCon *domain.TexContext, pReqInfo *domain.RequestMoneyCollect) {
	c.logger.Trace("【%v】START: moneyCollect StatusModeAll", texCon.GetUniqueKey())
	preqInfo, preqInfo2, preqInfo3, preqInfo4 := c.texmyHandler.Collect(texCon, pReqInfo)
	switch pReqInfo.OutType {
	case domain.WITHDRAW_TO_OUTLET: //現金入出金制御 出金要求:出金口に出金
		if pReqInfo.StatusMode == domain.START {
			c.texmyHandler.SetSequence(texCon, domain.ALLCOLLECT_START_OUT_START)
			c.sendRecv.SendRequestOutStart(texCon, &preqInfo3) //出金要求

		} else if pReqInfo.StatusMode == domain.STOP {
			c.texmyHandler.SetSequence(texCon, domain.ALLCOLLECT_START_OUT_STOP)
			c.sendRecv.SendRequestCollectStart(texCon, &preqInfo4) //出金停止要求
		}
	case domain.COLLECT_TO_COLLECTION_BOX: //現金入出金制御 出金要求:回収庫に回収
		if pReqInfo.StatusMode == domain.START {
			c.texmyHandler.SetSequence(texCon, domain.ALLCOLLECT_START_COLLECT_START)
			c.sendRecv.SendRequestOutStop(texCon, &preqInfo) //回収要求
		} else if pReqInfo.StatusMode == domain.STOP {
			c.texmyHandler.SetSequence(texCon, domain.ALLCOLLECT_START_COLLECT_STOP)
			c.sendRecv.SendRequestCollectStop(texCon, &preqInfo2) //回収停止要求
		}
	}
	c.logger.Trace("【%v】END: moneyCollect StatusModeAll", texCon.GetUniqueKey())
}

// 途中回収and売上金回収含
func (c *moneyCollect) StatusModeMiddleAndSales(texCon *domain.TexContext, pReqInfo *domain.RequestMoneyCollect) {
	c.logger.Trace("【%v】START: moneyCollect StatusModeMiddleAndSales", texCon.GetUniqueKey())

	// セットされた値が0の場合、回収対象無しとして終了
	if pReqInfo.SalesAmount == 0 {
		c.result(texCon)
		c.texmyHandler.SensorZeroNoticeCollect(texCon)
		return
	}

	switch pReqInfo.StatusMode {

	case domain.START:
		// シーケンス番号セット(途中回収AND売上金回収)
		c.texmyHandler.SetSequence(texCon, domain.MIDDLE_START_OUT_START)
		// リクエスト情報生成
		req, _ := c.texmyHandler.MiddleAndSalesCollect(texCon, pReqInfo)
		c.sendRecv.SendRequestOutStart(texCon, &req) //出金開始要求
	case domain.STOP:
		// シーケンス番号セット(途中回収AND売上金回収)
		c.texmyHandler.SetSequence(texCon, domain.MIDDLE_START_OUT_STOP)
		// リクエスト情報生成
		_, req := c.texmyHandler.MiddleAndSalesCollect(texCon, pReqInfo)
		c.sendRecv.SendRequestCollectStart(texCon, &req) //出金停止要求

	}

}

// 各要求送信完了検知
func (c *moneyCollect) SenSorSendFinish(texCon *domain.TexContext, reqType int) {
	c.logger.Trace("【%v】START: moneyCollect SenSorSendFinish", texCon.GetUniqueKey())
	switch reqType {
	case domain.FINISH_COLLECT_START, domain.FINISH_OUT_START: //回収開始要求完了,出金開始要求完了
		c.logger.Debug("【%v】- FINISH_COLLECT_START,FINISH_OUT_STAR", texCon.GetUniqueKey())
		//確定：稼働データ管理に金庫状態記録を投げる
		reqInfo := c.texmyHandler.RequestReportSafeInfo(texCon)
		c.texdtSendRecv.SendRequestReportSafeInfo(texCon, &reqInfo)

	case domain.FINISH_REPORT_SAFEINFO: //金庫情報遷移記録完了
		c.logger.Debug("【%v】- FINISH_REPORT_SAFEINFO", texCon.GetUniqueKey())
		c.SenSorSendFinish(texCon, domain.FINISH_PRINT_CHANGE_SUPPLY)

	case domain.FINISH_PRINT_CHANGE_SUPPLY: //印刷要求完了
		c.logger.Debug("【%v】- FINISH_PRINT_CHANGE_SUPPLY", texCon.GetUniqueKey())
		reqInfo := c.changeStatus.RequestChangeSupply(texCon, 0)
		c.statusTxSendRecv.SendRequestChangeSupply(texCon, &reqInfo)
	case domain.FINISH_CHANGE_SUPPLY: //精算機状態管理要求完了
		c.logger.Debug("【%v】- FINISH_CHANGE_SUPPLY", texCon.GetUniqueKey())
		//完了通知
		c.texmyHandler.SetTexmyNoticeOutdata(texCon, true)
		c.texmyHandler.SetTexmyNoticeCollectdata(texCon, true)
		c.texmyHandler.SetFlagCollect(texCon, false) //回収完了時にフラグを初期化

		time.Sleep(2 * time.Second) // 2秒待つ 下位レイヤーから有高が上がってくるまでの時間が2秒ほどある為
		c.texmyHandler.SetTexmyNoticeAmountData(texCon)

	}
	c.logger.Trace("【%v】END: moneyCollect SenSorSendFinish", texCon.GetUniqueKey())
}
