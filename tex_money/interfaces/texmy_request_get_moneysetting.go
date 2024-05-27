package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/pkg/lib"

	"tex_money/usecases"
)

type getMoneySetting struct {
	mqtt         handler.MqttRepository
	logger       handler.LoggerRepository
	syslogMng    usecases.SyslogManager
	errorMng     usecases.ErrorManager
	texmyHandler usecases.TexMoneyHandlerRepository
}

func NewRequestGetMoneySetting(
	mqtt handler.MqttRepository,
	logger handler.LoggerRepository,
	syslogMng usecases.SyslogManager,
	errorMng usecases.ErrorManager,
	texmyHandler usecases.TexMoneyHandlerRepository,
) RegisterMoneySettingRepository {
	return &getMoneySetting{
		mqtt:         mqtt,
		logger:       logger,
		syslogMng:    syslogMng,
		errorMng:     errorMng,
		texmyHandler: texmyHandler,
	}
}

// 開始処理
func (c *getMoneySetting) Start() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_get_moneysetting")
	c.mqtt.Subscribe(topic, c.recvRequest)
}

// 停止処理
func (c *getMoneySetting) Stop() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_get_moneysetting")
	c.mqtt.Unsubscribe(topic)
}

// サービス制御要求検出
func (c *getMoneySetting) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

func (c *getMoneySetting) recvRequest(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_get_moneysetting",
	})

	c.logger.Trace("【%v】START:要求受信 request_get_moneysetting 金銭設定取得要求", texCon.GetUniqueKey())
	var reqInfo domain.RequestGetMoneySetting
	var resInfo domain.ResultGetMoneySetting

	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_GET_MONEYSETTING_FATAL, "", "入出金管理")
		c.logger.Error("getMoneySetting recvRequest json.Unmarshal:%v", err)
		return
	}
	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), reqInfo.RequestInfo.RequestID)

	resInfo.RequestInfo = reqInfo.RequestInfo
	resInfo.Result = true

	data := c.texmyHandler.GetMoneySetting() //金銭設定情報取得

	// 最終登録日付，最終登録時刻が空の場合(初回設定時等)には，取得要求受信時の日時情報をセットする
	data, err = c.checkDateTime(texCon, data)
	if err != nil {
		resInfo.Result = false
		resInfo.ErrorCode, resInfo.ErrorDetail = c.errorMng.GetErrorInfo(usecases.ERROR_INSIDE)
	}
	resInfo.ChangeReserveCount = data.ChangeReserveCount
	resInfo.ChangeShortageCount = data.ChangeShortageCount
	resInfo.ExcessChangeCount = data.ExcessChangeCount

	payment, err := json.Marshal(resInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_GET_MONEYSETTING_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Unmarshal:%v", texCon.GetUniqueKey(), err)
		return
	}
	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_GET_MONEYSETTING_SUCCESS, "", "入出金管理")
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_get_moneysetting")
	c.mqtt.Publish(topic, string(payment))
	c.logger.Trace("【%v】END:要求受信 request_get_moneysetting 金銭設定取得要求", texCon.GetUniqueKey())
}

// 最終登録日付，最終登録時刻が空の場合(初回設定時等)には，取得要求受信時の日時情報をセットする
func (c *getMoneySetting) checkDateTime(texCon *domain.TexContext, data *domain.MoneySetting) (*domain.MoneySetting, error) {

	// 現在日時を取得
	date, time, err := lib.GeDateTime()
	if err != nil {
		c.logger.Error("【%v】getMoneySetting checkDateTime 現在日時取得失敗 err=%v", texCon.GetUniqueKey(), err)
		return data, err
	}

	// 日時が空の場合は現在日時をセット
	if len(data.ChangeReserveCount.LastRegistDate) == 0 && len(data.ChangeReserveCount.LastRegistTime) == 0 {
		data.ChangeReserveCount.LastRegistDate = date
		data.ChangeReserveCount.LastRegistTime = time
	}
	if len(data.ChangeShortageCount.LastRegistDate) == 0 && len(data.ChangeShortageCount.LastRegistTime) == 0 {
		data.ChangeShortageCount.LastRegistDate = date
		data.ChangeShortageCount.LastRegistTime = time
	}
	if len(data.ExcessChangeCount.LastRegistDate) == 0 && len(data.ExcessChangeCount.LastRegistTime) == 0 {
		data.ExcessChangeCount.LastRegistDate = date
		data.ExcessChangeCount.LastRegistTime = time
	}

	return data, nil
}
