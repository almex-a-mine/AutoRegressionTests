package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases"
)

type clearCashInfo struct {
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
	temReqInfo       domain.RequestClearCashInfo //一時リクエスト情報格納
	safeInfoMng      usecases.SafeInfoManager
}

// 入出金データクリア要求
func NewRequestClearCashInfo(mqtt handler.MqttRepository, logger handler.LoggerRepository, config config.Configuration, syslogMng usecases.SyslogManager, errorMng usecases.ErrorManager, sendRecv SendRecvRepository, texdtSendRecv TexdtSendRecvRepository, statusTxSendRecv StatusTxSendRecvRepository, printSendRecv PrintSendRecvRepository, texmyHandler usecases.TexMoneyHandlerRepository, safeInfoMng usecases.SafeInfoManager) RequestClearCashInfoRepository {
	return &clearCashInfo{
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
		safeInfoMng:      safeInfoMng,
	}
}

// 開始処理
func (c *clearCashInfo) Start() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_clear_cashinfo")
	c.mqtt.Subscribe(topic, c.recvRequest)
}

// 停止処理
func (c *clearCashInfo) Stop() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_clear_cashinfo")
	c.mqtt.Unsubscribe(topic)
}

// サービス制御要求検出
func (c *clearCashInfo) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

// リクエスト
func (c *clearCashInfo) recvRequest(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{ReceivingTopicName: "request_clear_cashinfo"})

	c.logger.Trace("【%v】START:要求受信 request_clear_cashinfo 入出金データクリア要求", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:要求受信 request_clear_cashinfo 入出金データクリア要求", texCon.GetUniqueKey())

	c.texmyHandler.SetSequence(texCon, domain.CLEAR_CASHINFO) //状態セット
	var reqInfo domain.RequestClearCashInfo
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_CLEAR_CASHINFO_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Unmarshal:%v", texCon.GetUniqueKey(), err)
		return
	}
	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), reqInfo.RequestInfo.RequestID)

	c.temReqInfo = reqInfo //リクエスト情報保存

	// 現在枚数を取得する
	_, initialSortInfo := c.safeInfoMng.GetSortInfo(texCon, domain.CASH_AVAILABLE) //現金有高取得
	// 現在枚数を初期補充にセットする。
	/// 登録時のSortTypeを変更
	initialSortInfo.SortType = domain.INITIAL_REPLENISHMENT
	/// 登録
	c.safeInfoMng.UpdateSortInfo(texCon, initialSortInfo)

	//safeInfoManagerの入出金データクリア
	c.safeInfoMng.ClearCashInfo(texCon)

	//safeInfoManagerの売上金情報クリア
	c.safeInfoMng.ClearSalesInfo(texCon)

	//稼働データへデータクリアを投げる
	//稼働データ管理に金庫状態記録を投げる
	reqSafeInfo := c.texmyHandler.RequestReportSafeInfo(texCon)
	c.texdtSendRecv.SendRequestReportSafeInfo(texCon, &reqSafeInfo)

}

// 結果応答
func (c *clearCashInfo) SendResult(texCon *domain.TexContext, presInfo *domain.ResultReportSafeInfo) {
	c.logger.Trace("【%v】START:clearCashInfo SendResult", texCon.GetUniqueKey())
	c.logger.Debug("【%v】- presInfo=%+v", texCon.GetUniqueKey(), presInfo)

	res := domain.ResultClearCashInfo{
		RequestInfo: c.temReqInfo.RequestInfo,
		Result:      presInfo.Result,
	}

	if !res.Result {
		res.ErrorCode = presInfo.ErrorCode
		res.ErrorDetail = presInfo.ErrorDetail
	}

	amount, err := json.Marshal(res)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_CLEAR_CASHINFO_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Unmarshal:%v", texCon.GetUniqueKey(), err)
	} else {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_CLEAR_CASHINFO_SUCCESS, "", "入出金管理")
		topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_clear_cashinfo")
		c.mqtt.Publish(topic, string(amount))
	}
	c.texmyHandler.SetErrorFromRequest(texCon, res.Result, res.ErrorCode, res.ErrorDetail)
	c.logger.Trace("【%v】END: clearCashInfo SendResult", texCon.GetUniqueKey())
}
