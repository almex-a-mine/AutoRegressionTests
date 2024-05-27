package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases"
)

type changeSystemOperation struct {
	mqtt                  handler.MqttRepository
	logger                handler.LoggerRepository
	syslogMng             usecases.SyslogManager
	texMoneyNoticeManager usecases.TexMoneyNoticeManagerRepository
	noticeSend            NoticeSendRepository
}

// システム動作モード管理
func NewRequestChangeSystemOperation(mqtt handler.MqttRepository,
	logger handler.LoggerRepository,
	syslogMng usecases.SyslogManager,
	texMoneyNoticeManager usecases.TexMoneyNoticeManagerRepository) RequestChangeSystemOperationRepository {
	return &changeSystemOperation{
		mqtt:                  mqtt,
		logger:                logger,
		syslogMng:             syslogMng,
		texMoneyNoticeManager: texMoneyNoticeManager,
		noticeSend:            NewNoticeSend(mqtt, logger),
	}
}

// 開始処理
func (c *changeSystemOperation) Start() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_change_system_operation")
	c.mqtt.Subscribe(topic, c.recvRequestChangeSystemOperation)
}

// 停止処理
func (c *changeSystemOperation) Stop() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_change_system_operation")
	c.mqtt.Unsubscribe(topic)
}

// サービス制御要求検出
func (c *changeSystemOperation) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

// システム動作モード変更要求検出
func (c *changeSystemOperation) recvRequestChangeSystemOperation(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_change_systemoperation",
	})
	c.logger.Trace("【%v】START:要求受信 request_change_system_operation システム動作モード変更要求", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:要求受信 request_change_system_operation システム動作モード変更要求", texCon.GetUniqueKey())

	var reqInfo domain.RequestChangeSystemOperation

	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.logger.Error("changeSystemOperation recvRequestChangeSystemOperation json.Unmarshal:%v", err)
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_CHANGE_SYSTEM_OPERATION_FATAL, "", "入出金管理")
		return
	}
	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), reqInfo.RequestInfo.RequestID)

	res := domain.ResultChangeSystemOperation{
		RequestInfo: reqInfo.RequestInfo,
		Result:      true,
	}

	c.resultChangeSystemOperationService(texCon, &res)
	//loggerの設定変更
	c.logger.SetSystemOperation(reqInfo.StatusSystemOperation)
	//状態更新
	ok := c.texMoneyNoticeManager.UpdateStatusSystemOperationData(texCon,
		domain.StatusSystemData{
			StatusSystemOperation: reqInfo.StatusSystemOperation,
			IdDevice:              reqInfo.IdDevice,
			IdExtSys:              reqInfo.IdExtSys})
	//状況を通知
	if ok {
		statusSystemData := c.texMoneyNoticeManager.GetStatusSystemOperationData(texCon)

		topic := domain.TOPIC_TEXMONEY_BASE + "/notice_status_system_operation"
		err := c.noticeSend.Send(texCon, statusSystemData, topic)
		if err != nil {
			c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_NOTICE_STATUS_SYSTEM_OPERATION_ERROR, "", "入出金管理")
			return
		}
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_NOTICE_STATUS_SYSTEM_OPERATION_SUCCESS, "", "入出金管理")
	}

}

// 処理結果応答
func (c *changeSystemOperation) resultChangeSystemOperationService(texCon *domain.TexContext, pResInfo *domain.ResultChangeSystemOperation) {
	c.logger.Trace("【%v】START: changeSystemOperation resultChangeSystemOperationService", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END: changeSystemOperation resultChangeSystemOperationService", texCon.GetUniqueKey())
	payment, err := json.Marshal(pResInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_CHANGE_SYSTEM_OPERATION_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Unmarshal:%v", texCon.GetUniqueKey(), err)
		return
	}
	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_CHANGE_SYSTEM_OPERATION_SUCCESS, "", "入出金管理")
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_change_system_operation")
	c.mqtt.Publish(topic, string(payment))

}
