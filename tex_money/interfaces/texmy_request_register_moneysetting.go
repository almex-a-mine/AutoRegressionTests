package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/pkg/lib"

	"tex_money/usecases"
)

type registerMoneySetting struct {
	mqtt                  handler.MqttRepository
	logger                handler.LoggerRepository
	config                config.Configuration
	syslogMng             usecases.SyslogManager
	errorMng              usecases.ErrorManager
	texmyHandler          usecases.TexMoneyHandlerRepository
	iniService            usecases.IniServiceRepository
	texMoneyNoticeManager usecases.TexMoneyNoticeManagerRepository
}

func NewRequestRegisterMoneySetting(
	mqtt handler.MqttRepository,
	logger handler.LoggerRepository,
	config config.Configuration,
	syslogMng usecases.SyslogManager,
	errorMng usecases.ErrorManager,
	texmyHandler usecases.TexMoneyHandlerRepository,
	iniService usecases.IniServiceRepository,
	texMoneyNoticeManager usecases.TexMoneyNoticeManagerRepository,
) RegisterMoneySettingRepository {
	return &registerMoneySetting{
		mqtt:                  mqtt,
		logger:                logger,
		config:                config,
		syslogMng:             syslogMng,
		errorMng:              errorMng,
		texmyHandler:          texmyHandler,
		iniService:            iniService,
		texMoneyNoticeManager: texMoneyNoticeManager,
	}
}

// 開始処理
func (c *registerMoneySetting) Start() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_register_moneysetting")
	c.mqtt.Subscribe(topic, c.recvRequest)
}

// 停止処理
func (c *registerMoneySetting) Stop() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_register_moneysetting")
	c.mqtt.Unsubscribe(topic)
}

// サービス制御要求検出
func (c *registerMoneySetting) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

func (c *registerMoneySetting) recvRequest(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_register_moneysetting",
	})

	c.logger.Trace("【%v】START:要求受信 request_register_moneysetting 金銭設定登録要求", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:要求受信 request_register_moneysetting 金銭設定登録要求", texCon.GetUniqueKey())
	reqInfo := &domain.RequestRegisterMoneySetting{}
	var resInfo domain.ResultRegisterMoneySetting

	err := json.Unmarshal([]byte(message), reqInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_REGISTER_MONEYSETTING_FATAL, "", "入出金管理")
		c.logger.Error("registerMoneySetting recvRequest json.Unmarshal:%v", err)
		return
	}
	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), reqInfo.RequestInfo.RequestID)

	resInfo.RequestInfo = reqInfo.RequestInfo
	resInfo.Result = true

	// 更新データ取得
	updateData, err := c.getDiffData(texCon, reqInfo)
	if err != nil {
		resInfo.Result = false
		resInfo.ErrorCode, resInfo.ErrorDetail = c.errorMng.GetErrorInfo(usecases.ERROR_INSIDE)
	}

	// texmyHandlerの内部変数更新
	c.texmyHandler.SetMoneySetting(updateData)

	// iniに書き込む
	err = c.updateMoneySettingIni(texCon, updateData)
	if err != nil {
		resInfo.Result = false
		resInfo.ErrorCode, resInfo.ErrorDetail = c.errorMng.GetErrorInfo(usecases.ERROR_INSIDE)
	}
	c.sendResult(texCon, &resInfo)
	statusCash := c.texMoneyNoticeManager.GetStatusCashData(texCon)
	statusCash = c.texmyHandler.CheckAmountLimit(texCon, statusCash) //リミット有高チェックを毎回行う

	c.texmyHandler.SetTexmyNoticeStatus(texCon, statusCash) //入出金機の状態が変化した場合に通知を発報する //設定の変更時に不足エラー等を発生/解消させるため

}

// 更新データを取得
func (c *registerMoneySetting) getDiffData(texCon *domain.TexContext, reqInfo *domain.RequestRegisterMoneySetting) (*domain.MoneySetting, error) {
	data := c.texmyHandler.GetMoneySetting() // 現在の金銭設定情報取得

	// 現在日時を取得
	date, time, err := lib.GeDateTime()
	if err != nil {
		c.logger.Error("【%v】registerMoneySetting getDiffData 現在日時取得失敗 err=%v", texCon.GetUniqueKey(), err)
		return data, err
	}

	// 更新がある箇所のみ更新データをセット
	if reqInfo.ChangeReserveCount != nil {
		data.ChangeReserveCount = *reqInfo.ChangeReserveCount
		data.ChangeReserveCount.LastRegistDate = date
		data.ChangeReserveCount.LastRegistTime = time
	}
	if reqInfo.ChangeShortageCount != nil {
		data.ChangeShortageCount = *reqInfo.ChangeShortageCount
		data.ChangeShortageCount.LastRegistDate = date
		data.ChangeShortageCount.LastRegistTime = time
	}
	if reqInfo.ExcessChangeCount != nil {
		data.ExcessChangeCount = *reqInfo.ExcessChangeCount
		data.ExcessChangeCount.LastRegistDate = date
		data.ExcessChangeCount.LastRegistTime = time
	}

	return data, nil
}

// iniのMoneySettingを更新
func (c *registerMoneySetting) updateMoneySettingIni(texCon *domain.TexContext, moneySetting *domain.MoneySetting) error {

	// 最外部のmapを初期化
	sectionKeyValue := make(map[string]map[string]string)

	if sectionKeyValue["MoneySetting"] == nil {
		sectionKeyValue["MoneySetting"] = make(map[string]string)
	}

	// json文字列に変換して書き込み
	// 釣銭準備金枚数
	bytes1, err := json.Marshal(moneySetting.ChangeReserveCount)
	if err != nil {
		c.logger.Error("【%v】registerMoneySetting updateMoneySettingIni err=%v", texCon.GetUniqueKey(), err)
		return err
	}
	sectionKeyValue["MoneySetting"]["ChangeReserveCount"] = string(bytes1)

	//不足枚数
	bytes2, err := json.Marshal(moneySetting.ChangeShortageCount)
	if err != nil {
		c.logger.Error("【%v】registerMoneySetting updateMoneySettingIni err=%v", texCon.GetUniqueKey(), err)
		return err
	}
	sectionKeyValue["MoneySetting"]["ChangeShortageCount"] = string(bytes2)

	//あふれ枚数
	bytes3, err := json.Marshal(moneySetting.ExcessChangeCount)
	if err != nil {
		c.logger.Error("【%v】registerMoneySetting updateMoneySettingIni err=%v", texCon.GetUniqueKey(), err)
		return err
	}
	sectionKeyValue["MoneySetting"]["ExcessChangeCount"] = string(bytes3)

	c.iniService.MultipleUpdateIni(texCon, sectionKeyValue)

	return nil
}

// 結果応答
func (c *registerMoneySetting) sendResult(texCon *domain.TexContext, resInfo *domain.ResultRegisterMoneySetting) {
	c.logger.Trace("【%v】START: registerMoneySetting SendResult", texCon.GetUniqueKey())
	payment, err := json.Marshal(resInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_REGISTER_MONEYSETTING_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Unmarshal:%v", texCon.GetUniqueKey(), err)
	} else {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_REGISTER_MONEYSETTING_SUCCESS, "", "入出金管理")
		topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_register_moneysetting")
		c.mqtt.Publish(topic, string(payment))
	}
	c.logger.Trace("【%v】END: registerMoneySetting SendResult", texCon.GetUniqueKey())
}
