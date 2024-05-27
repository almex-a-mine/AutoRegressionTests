package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases"
)

type statusCash struct {
	mqtt                  handler.MqttRepository
	logger                handler.LoggerRepository
	syslogMng             usecases.SyslogManager
	texMoneyNoticeManager usecases.TexMoneyNoticeManagerRepository
}

// 現金入出金機制御ステータス要求
func NewRequestStatusCash(mqtt handler.MqttRepository,
	logger handler.LoggerRepository,
	syslogMng usecases.SyslogManager,
	texMoneyNoticeManager usecases.TexMoneyNoticeManagerRepository) StatusCashRepository {
	return &statusCash{
		mqtt:                  mqtt,
		logger:                logger,
		syslogMng:             syslogMng,
		texMoneyNoticeManager: texMoneyNoticeManager,
	}
}

// 開始処理
func (c *statusCash) Start() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_status_cash")
	c.mqtt.Subscribe(topic, c.recvRequest)
}

// 停止処理
func (c *statusCash) Stop() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_status_cash")
	c.mqtt.Unsubscribe(topic)
}

// サービス制御要求検出
func (c *statusCash) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

func (c *statusCash) recvRequest(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_status_cash",
	})
	c.logger.Trace("【%v】START:要求受信 request_status_cash 現金入出金制御ステータス要求", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:要求受信 request_status_cash 現金入出金制御ステータス要求", texCon.GetUniqueKey())
	var reqInfo domain.RequestStatusCash
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_STATUS_CASH_FATAL, "", "入出金管理")
		c.logger.Error("statusCash recvRequest json.Unmarshal:%v", err)
		return
	}

	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), reqInfo.RequestInfo.RequestID)

	// ステータス情報の取得
	statusCashData := c.texMoneyNoticeManager.GetStatusCashData(texCon)
	//値のセット
	resInfo := domain.ResultStatusCash{
		RequestInfo:         reqInfo.RequestInfo,
		Result:              true,
		CashControlId:       statusCashData.CashControlId,
		StatusReady:         statusCashData.StatusReady,
		StatusMode:          statusCashData.StatusMode,
		StatusLine:          statusCashData.StatusLine,
		StatusError:         statusCashData.StatusError,
		ErrorCode:           statusCashData.ErrorCode,
		ErrorDetail:         statusCashData.ErrorDetail,
		StatusCover:         statusCashData.StatusCover,
		StatusAction:        statusCashData.StatusAction,
		StatusInsert:        statusCashData.StatusInsert,
		StatusExit:          statusCashData.StatusExit,
		StatusRjbox:         statusCashData.StatusRjbox,
		BillStatusTbl:       statusCashData.BillStatusTbl,
		CoinStatusTbl:       statusCashData.CoinStatusTbl,
		BillResidueInfoTbl:  statusCashData.BillResidueInfoTbl,
		CoinResidueInfoTbl:  statusCashData.CoinResidueInfoTbl,
		DeviceStatusInfoTbl: statusCashData.DeviceStatusInfoTbl,
		WarningInfoTbl:      statusCashData.WarningInfoTbl,
	}

	payment, err := json.Marshal(resInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_STATUS_CASH_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Unmarshal:%v", texCon.GetUniqueKey(), err)
		return
	}
	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_STATUS_CASH_SUCCESS, "", "入出金管理")
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_status_cash")
	c.mqtt.Publish(topic, string(payment))

}
