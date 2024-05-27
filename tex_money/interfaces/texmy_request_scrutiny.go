package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases"
)

type ScrutinyRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
}

type scrutiny struct {
	mqtt                handler.MqttRepository
	logger              handler.LoggerRepository
	syslogMng           usecases.SyslogManager
	config              config.Configuration
	errorManager        usecases.ErrorManager
	texMoneyHandler     usecases.TexMoneyHandlerRepository
	cashControlSendRecv SendRecvRepository
	texStatusSendRecv   StatusTxSendRecvRepository
	reqTopic            string
	resTopic            string
	cashControlId       string
}

func NewScrutiny(
	mqtt handler.MqttRepository,
	logger handler.LoggerRepository,
	syslogMng usecases.SyslogManager,
	config config.Configuration,
	errorManager usecases.ErrorManager,
	texMoneyHandler usecases.TexMoneyHandlerRepository,
	cashControlSendRecv SendRecvRepository,
	texStatusSendRecv StatusTxSendRecvRepository,
) ScrutinyRepository {
	return &scrutiny{
		mqtt:                mqtt,
		logger:              logger,
		syslogMng:           syslogMng,
		config:              config,
		errorManager:        errorManager,
		texMoneyHandler:     texMoneyHandler,
		cashControlSendRecv: cashControlSendRecv,
		texStatusSendRecv:   texStatusSendRecv,
		reqTopic:            fmt.Sprintf("%s/%s", domain.TOPIC_TEXMONEY_BASE, "request_scrutiny"),
		resTopic:            fmt.Sprintf("%s/%s", domain.TOPIC_TEXMONEY_BASE, "result_scrutiny"),
	}
}

// 開始処理
func (c *scrutiny) Start() {
	c.mqtt.Subscribe(c.reqTopic, c.recvRequest)
}

// 停止処理
func (c *scrutiny) Stop() {
	c.mqtt.Unsubscribe(c.reqTopic)
}

// サービス制御要求検出
func (c *scrutiny) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

func (c *scrutiny) recvRequest(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_scrutiny",
	})

	var err error
	req := &domain.RequestScrutiny{}
	c.logger.Trace("【%v】START:要求受信 request_scrutiny 精査モード要求", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:要求受信 request_scrutiny 精査モード要求", texCon.GetUniqueKey())
	err = json.Unmarshal([]byte(message), req)
	if err != nil {
		errorCode, errorDetail := c.errorManager.GetErrorInfo(usecases.ERROR_INSIDE)
		c.handleError(true, usecases.SYSLOG_LOGTYPE_REQUEST_SCRUTINY_FATAL, fmt.Sprintf("scrutiny request json.Unmarshal:%v", err), req.RequestInfo, errorCode, errorDetail)
		return
	}
	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), req.RequestInfo.RequestID)

	c.texMoneyHandler.SetSequence(texCon, domain.SCRUTINY_START)

	// 2.1への精査開始要求送信
	if err := c.sendRequestScrutinyStart(texCon, req.TargetDevice); err != nil {
		errorCode, errorDetail := c.errorManager.GetErrorInfo(usecases.ERROR_COMMUNICATION_FAIL)
		c.handleError(true, usecases.SYSLOG_LOGTYPE_REQUEST_SCRUTINY_FATAL, err.Error(), req.RequestInfo, errorCode, errorDetail)
	} else {
		// 結果応答送信
		c.sendResult(req)
	}
	// 1.1への状態変更要求(スタッフ操作記録)送信
	c.sendRequestChangeStaffOperation(texCon)
}

// 2.1への精査開始要求送信
func (c *scrutiny) sendRequestScrutinyStart(texCon *domain.TexContext, targetDevice int) error {
	// リクエスト生成
	req := domain.NewRequestScrutinyStart(c.texMoneyHandler.NewRequestInfo(texCon), targetDevice)

	// 受信チャネルを生成
	var resChan = make(chan interface{})
	// 外部リクエスト
	go c.cashControlSendRecv.SendRequestScrutinyStart(texCon, resChan, req)

	// 外部リクエスト受信
	scrutinyStart := <-resChan
	// エラーの場合もあるので、型チェックでOKなら次の処理に進む
	resInfo, ok := scrutinyStart.(domain.ResultScrutinyStart)
	if !ok { // 型チェックでエラーの場合
		err, _ := scrutinyStart.(error) // 型チェックエラー
		return err
	}

	c.cashControlId = resInfo.CashControlId
	return nil
}

const SCRUTINY_OPERATION_DETAIL = "精査モードの実行"

// 1.1への状態変更要求(スタッフ操作記録)送信
func (c *scrutiny) sendRequestChangeStaffOperation(texCon *domain.TexContext) {
	// リクエスト生成
	req := domain.NewRequestChangeStaffOperation(c.texMoneyHandler.NewRequestInfo(texCon), 73, SCRUTINY_OPERATION_DETAIL)

	// 外部リクエスト
	c.texStatusSendRecv.SendRequestChangeStaffOperation(texCon, req)
}

// 結果応答送信
func (c *scrutiny) sendResult(req *domain.RequestScrutiny) {
	// 応答情報作成
	resultInfo := domain.NewResultScrutiny(req.RequestInfo, true, "", "", c.cashControlId)

	res, err := json.Marshal(resultInfo)
	if err != nil {
		errorCode, errorDetail := c.errorManager.GetErrorInfo(usecases.ERROR_INSIDE)
		c.handleError(true, usecases.SYSLOG_LOGTYPE_REQUEST_SCRUTINY_FATAL, fmt.Sprintf("scrutiny sendResult json.Unmarshal:%v", err), req.RequestInfo, errorCode, errorDetail)
		return
	}
	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_SCRUTINY_SUCCESS, "", "入出金管理")

	// 送信
	c.mqtt.Publish(c.resTopic, string(res))
}

func (c *scrutiny) handleError(resultSend bool, code int, message string, req domain.RequestInfo, errorCode string, errorDetail string) {
	c.syslogMng.NoticeSystemLog(code, "", "入出金管理")
	c.logger.Error(message)

	// エラーを載せた応答送信
	if resultSend {
		resultInfo := domain.NewResultScrutiny(req, false, errorCode, errorDetail, "")

		res, err := json.Marshal(resultInfo)
		if err == nil {
			c.mqtt.Publish(c.resTopic, string(res))
		}
	}
}
