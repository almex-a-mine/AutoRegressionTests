package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases"
)

type maintenanceMode struct {
	mqtt               handler.MqttRepository
	logger             handler.LoggerRepository
	syslogMng          usecases.SyslogManager
	errorMng           usecases.ErrorManager
	aggregateMng       usecases.AggregateManager
	safeInfoMng        usecases.SafeInfoManager
	maintenanceModeMng usecases.MaintenanceModeManager
	reqTopic           string
	resTopic           string
}

func NewRequestMaintenanceMode(
	mqtt handler.MqttRepository,
	logger handler.LoggerRepository,
	syslogMng usecases.SyslogManager,
	errorMng usecases.ErrorManager,
	aggregateMng usecases.AggregateManager,
	safeInfoMng usecases.SafeInfoManager,
	maintenanceModeMng usecases.MaintenanceModeManager) MaintenanceModeRepository {
	return &maintenanceMode{
		mqtt:               mqtt,
		logger:             logger,
		syslogMng:          syslogMng,
		errorMng:           errorMng,
		aggregateMng:       aggregateMng,
		safeInfoMng:        safeInfoMng,
		maintenanceModeMng: maintenanceModeMng,
		reqTopic:           fmt.Sprintf("%s/%s", domain.TOPIC_TEXMONEY_BASE, "request_maintenance_mode"),
		resTopic:           fmt.Sprintf("%s/%s", domain.TOPIC_TEXMONEY_BASE, "result_maintenance_mode"),
	}
}

// 開始処理
func (c *maintenanceMode) Start() {
	c.mqtt.Subscribe(c.reqTopic, c.recvRequest)
}

// 停止処理
func (c *maintenanceMode) Stop() {
	c.mqtt.Unsubscribe(c.reqTopic)
}

// サービス制御要求検出
func (c *maintenanceMode) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

// 保守業務モード要求 受信時処理
func (c *maintenanceMode) recvRequest(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_maintenance_mode",
	})

	c.logger.Trace("【%v】START:要求受信 request_maintenance_mode 保守業務モード要求", texCon.GetUniqueKey())
	var reqInfo domain.RequestMaintenanceMode

	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.handleError(true, usecases.SYSLOG_LOGTYPE_REQUEST_MAINTENANCE_MODE_FATAL, fmt.Sprintf("maintenanceMode request json.Unmarshal:%v", err), reqInfo.RequestInfo)
		return
	}
	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), reqInfo.RequestInfo.RequestID)

	// 受信値を保存
	c.maintenanceModeMng.SetReceiveData(texCon, reqInfo)

	var result bool
	if reqInfo.Action {
		//開始処理
		result = c.maintenanceModeMng.SetStatusStart(texCon)
	} else {
		//終了処理
		result = c.maintenanceModeMng.SetStatusEnd(texCon)
	}
	if !result {
		c.handleError(true, usecases.SYSLOG_LOGTYPE_REQUEST_MAINTENANCE_MODE_FATAL, "maintenanceMode 失敗", reqInfo.RequestInfo)
		return
	}

	// 応答メッセージ作成
	resInfo := domain.NewResultMaintenanceMode(reqInfo.RequestInfo, true, "", "")
	res, err := json.Marshal(resInfo)
	if err != nil {
		c.handleError(true, usecases.SYSLOG_LOGTYPE_REQUEST_MAINTENANCE_MODE_FATAL, fmt.Sprintf("maintenanceMode sendResult json.Unmarshal:%v", err), reqInfo.RequestInfo)
		return
	}
	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_MAINTENANCE_MODE_SUCCESS, "", "入出金管理")

	// 応答送信
	c.mqtt.Publish(c.resTopic, string(res))

	// 動作開始/終了時の金庫・レポート情報 ログ出力
	c.logger.Debug("保守モード操作時点保持データ")
	c.safeInfoMng.OutputLogSafeInfoExCountTbl(texCon)
	c.aggregateMng.OutputLogAggregateExCountTbl()
	c.logger.Trace("【%v】END:要求受信 request_maintenance_mode 保守業務モード要求", texCon.GetUniqueKey())
}

// エラー発報
func (c *maintenanceMode) handleError(resultSend bool, code int, message string, req domain.RequestInfo) {
	c.syslogMng.NoticeSystemLog(code, "", "入出金管理")
	c.logger.Error(message)

	if resultSend {
		errorCode, errorDetail := c.errorMng.GetErrorInfo(usecases.ERROR_INSIDE)
		resultInfo := domain.NewResultMaintenanceMode(req, false, errorCode, errorDetail)
		res, err := json.Marshal(resultInfo)
		if err == nil {
			c.mqtt.Publish(c.resTopic, string(res))
		}
	}
}
