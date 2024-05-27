package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases"
)

type CoinCassetteControlRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
	recvRequest(message string)
}

type coinCassetteControl struct {
	mqtt                handler.MqttRepository
	logger              handler.LoggerRepository
	syslogMng           usecases.SyslogManager
	config              config.Configuration
	coinCassetteManager usecases.CoinCassetteControlManager
	setAmount           SetAmountRepository
	errorManager        usecases.ErrorManager
	reqTopic            string
	resTopic            string
	aggregateMng        usecases.AggregateManager
	texMoneyHandler     usecases.TexMoneyHandlerRepository
}

func NewCoinCassetteControl(
	mqtt handler.MqttRepository,
	logger handler.LoggerRepository,
	syslogMng usecases.SyslogManager,
	config config.Configuration,
	errorManager usecases.ErrorManager,
	coinCassetteManager usecases.CoinCassetteControlManager,
	setAmount SetAmountRepository,
	aggregateMng usecases.AggregateManager,
	texMoneyHandler usecases.TexMoneyHandlerRepository) CoinCassetteControlRepository {
	return &coinCassetteControl{
		mqtt:                mqtt,
		logger:              logger,
		syslogMng:           syslogMng,
		config:              config,
		coinCassetteManager: coinCassetteManager,
		setAmount:           setAmount,
		errorManager:        errorManager,
		reqTopic:            fmt.Sprintf("%s/%s", domain.TOPIC_TEXMONEY_BASE, "request_coincassette_control"),
		resTopic:            fmt.Sprintf("%s/%s", domain.TOPIC_TEXMONEY_BASE, "result_coincassette_control"),
		aggregateMng:        aggregateMng,
		texMoneyHandler:     texMoneyHandler}
}

// 開始処理
func (c *coinCassetteControl) Start() {
	c.mqtt.Subscribe(c.reqTopic, c.recvRequest)
}

// 停止処理
func (c *coinCassetteControl) Stop() {
	c.mqtt.Unsubscribe(c.reqTopic)
}

// サービス制御要求検出
func (c *coinCassetteControl) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

func (c *coinCassetteControl) recvRequest(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_coincassette_control",
	})
	var err error
	req := &domain.RequestCoincassetteControl{}
	c.logger.Trace("【%v】START:要求受信 request_coincassette_control 硬貨カセット操作要求", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:要求受信 request_coincassette_control 硬貨カセット操作要求", texCon.GetUniqueKey())
	err = json.Unmarshal([]byte(message), req)
	if err != nil {
		c.handleError(false, usecases.SYSLOG_LOGTYPE_REQUEST_COIN_CASSETTE_CONTROL_FATAL, fmt.Sprintf("coinCassetteControl request json.Unmarshal:%v", err), req.RequestInfo)
		return
	}
	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), req.RequestInfo.RequestID)

	resultCassette := &domain.CoinCassette{}

	switch req.ControlMode {
	case 1: // 回収
		resultCassette = c.coinCassetteManager.Collection(texCon, req.CoinCassette)
	case 2: // 交換
		resultCassette = c.coinCassetteManager.Exchange(texCon, req.CoinCassette)

	case 3: // 指定枚数補充
		resultCassette = c.coinCassetteManager.SpecificationReplenishment(texCon, req.CoinCassette, req.AmountCount)
	default:
		c.logger.Debug("【%v】coinCassetteControl-controlMode不正=%v", texCon.GetUniqueKey(), req.ControlMode)
		// エラー処理記述
		// return

	}

	go c.setAmount.ConnetctCoincasseteControl(texCon, resultCassette.AfterExCountTbl)

	resultInfo := domain.ResultCoincassetteControl{
		RequestInfo:           req.RequestInfo,
		Result:                true,
		DifferenceTotalAmount: resultCassette.DifferenceTotalAmount,
		DifferenceExCountTbl:  resultCassette.DifferenceExCountTbl,
		BeforeExCountTbl:      resultCassette.BeforeExCountTbl,
		AfterExCountTbl:       resultCassette.AfterExCountTbl,
		ExchangeExCountTbl:    resultCassette.ExchangeExCountTbl,
	}

	res, err := json.Marshal(resultInfo)
	if err != nil {
		c.handleError(true, usecases.SYSLOG_LOGTYPE_REQUEST_COIN_CASSETTE_CONTROL_FATAL, fmt.Sprintf("coinCassetteControl sendResult json.Unmarshal:%v", err), req.RequestInfo)
		return
	}
	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_COIN_CASSETTE_CONTROL_SUCCESS, "", "入出金管理")

	c.mqtt.Publish(c.resTopic, string(res))
}

func (c *coinCassetteControl) handleError(resultSend bool, code int, message string, req domain.RequestInfo) {
	c.syslogMng.NoticeSystemLog(code, "", "入出金管理")
	c.logger.Error(message)

	if resultSend {
		errorCode, errorDetail := c.errorManager.GetErrorInfo(usecases.ERROR_INSIDE)
		resultInfo := domain.ResultCoincassetteControl{
			RequestInfo: req,
			Result:      false,
			ErrorCode:   errorCode,
			ErrorDetail: errorDetail,
		}
		res, err := json.Marshal(resultInfo)
		if err == nil {
			c.mqtt.Publish(c.resTopic, string(res))
		}
	}
}
