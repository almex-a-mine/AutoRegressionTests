package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases"
)

type ReverseExchangeCalculationRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
	recvRequest(message string)
}

type reverseExchangeCalculation struct {
	mqtt            handler.MqttRepository
	logger          handler.LoggerRepository
	syslogMng       usecases.SyslogManager
	config          config.Configuration
	errorManager    usecases.ErrorManager
	texMoney        usecases.TexMoneyHandlerRepository
	paymentSendRecv PaymentSendRecvRepository
	reverseExchange usecases.ReverseExchangeCalculationRepository
	reqTopic        string
	resTopic        string
}

func NewReverseExchangeCalculation(
	mqtt handler.MqttRepository,
	logger handler.LoggerRepository,
	config config.Configuration,
	syslogMng usecases.SyslogManager,
	errorManager usecases.ErrorManager,
	texMoney usecases.TexMoneyHandlerRepository,
	reverseExchange usecases.ReverseExchangeCalculationRepository,
	paymentSendRecvRepository PaymentSendRecvRepository,
) ReverseExchangeCalculationRepository {
	return &reverseExchangeCalculation{
		mqtt:            mqtt,
		logger:          logger,
		syslogMng:       syslogMng,
		config:          config,
		errorManager:    errorManager,
		paymentSendRecv: paymentSendRecvRepository,
		texMoney:        texMoney,
		reverseExchange: reverseExchange,
		reqTopic:        fmt.Sprintf("%s/%s", domain.TOPIC_TEXMONEY_BASE, "request_reverseexchange_calculation"),
		resTopic:        fmt.Sprintf("%s/%s", domain.TOPIC_TEXMONEY_BASE, "result_reverseexchange_calculation"),
	}
}

// 開始処理
func (c *reverseExchangeCalculation) Start() {
	c.mqtt.Subscribe(c.reqTopic, c.recvRequest)
}

// 停止処理
func (c *reverseExchangeCalculation) Stop() {
	c.mqtt.Unsubscribe(c.reqTopic)
}

// サービス制御要求検出
func (c *reverseExchangeCalculation) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

func (c *reverseExchangeCalculation) recvRequest(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_revers_exchangecalculation",
	})

	c.logger.Trace("【%v】START:要求受信 request_reverseexchange_calculation 逆両替算出要求", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:要求受信 request_reverseexchange_calculation 逆両替算出要求", texCon.GetUniqueKey())

	inReq := &domain.RequestReverseExchangeCalculation{}
	err := json.Unmarshal([]byte(message), inReq)
	if err != nil {
		c.handleError(false, usecases.SYSLOG_LOGTYPE_REQUEST_REVERSE_EXCHNGE_CALCULATION_FATAL, fmt.Sprintf("coinCassetteControl request json.Unmarshal:%v", err), inReq.RequestInfo)
		return
	}
	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), inReq.RequestInfo.RequestID)

	inResult := &domain.ResultReverseExchangeCalculation{
		RequestInfo: inReq.RequestInfo,
		Result:      true,
	}

	// ExchangeTypeのバリデーション対応も含めて ここで、指定して分岐処理する
	c.logger.Debug("【%v】reverseExchangeCalculation recvRequest ExchangeType=%v", texCon.GetUniqueKey(), inReq.ExchangeType)
	switch inReq.ExchangeType {
	case 1, 2: // 売上金系指定(外部連携有関連)
		// リクエスト生成
		outReq := &domain.RequestGetSalesinfo{
			RequestInfo: c.texMoney.NewRequestInfo(texCon),
		}
		// 受信チャネルを生成
		var resChan = make(chan interface{})
		// 外部リクエスト
		go c.paymentSendRecv.SendRequestGetSalesInfo(texCon, resChan, outReq)
		// 外部リクエスト受信
		salesMoney := <-resChan
		// エラーの場合もあるので、型チェックでOKなら次の処理に進む
		outResInfo, ok := salesMoney.(domain.ResultGetSalesinfo)
		if !ok {
			err, _ = salesMoney.(error)
			c.handleError(true, usecases.SYSLOG_LOGTYPE_REQUEST_REVERSE_EXCHNGE_CALCULATION_FATAL, err.Error(), inReq.RequestInfo)
			return
		}

		var salesCashAmount int
		// 連携されてこない項目があった場合を考慮して、rangeで取り出し
		for _, v := range outResInfo.InfoSales.SalesTypeTbl {
			if v.PaymentType == 0 {
				switch v.SalesType {
				case 0, 1: // 0:チェックイン、1:チェックアウト
					salesCashAmount += v.Amount
				case 2, 3: // 2:チェックイン取消 3:チェックアウト取消
					salesCashAmount += v.Amount
				}

			}
		}

		salesAmount, exchangeTbl, err := c.reverseExchange.SalesMoneyExchange(texCon, inReq.ExchangeType, salesCashAmount, inReq.OverflowCashbox)
		if err != nil {
			c.handleError(true, usecases.SYSLOG_LOGTYPE_REQUEST_REVERSE_EXCHNGE_CALCULATION_FATAL, err.Error(), inReq.RequestInfo)
			return
		}
		inResult.TargetAmount = salesAmount
		inResult.TargetExCountTbl = nil
		inResult.ExchangeExCountTbl = exchangeTbl
	case 10, 11, 12, 13, 16: // 釣銭基準準備金ベースの逆両替
		diffTotal, diffTbl, exchange, err := c.reverseExchange.BaseExchange(texCon, inReq.ExchangeType)
		if err != nil {
			c.handleError(true, usecases.SYSLOG_LOGTYPE_REQUEST_REVERSE_EXCHNGE_CALCULATION_FATAL, err.Error(), inReq.RequestInfo)
			return
		}
		inResult.TargetAmount = diffTotal
		inResult.TargetExCountTbl = &diffTbl
		inResult.ExchangeExCountTbl = exchange
	case 30: // 指定金額
		amount, exchangeTbl, err := c.reverseExchange.SpecifyExchange(texCon, inReq.Amount)
		if err != nil {
			c.handleError(true, usecases.SYSLOG_LOGTYPE_REQUEST_REVERSE_EXCHNGE_CALCULATION_FATAL, err.Error(), inReq.RequestInfo)
			return
		}
		inResult.TargetAmount = amount
		inResult.TargetExCountTbl = nil
		inResult.ExchangeExCountTbl = exchangeTbl
	case 31: // 指定金額(払出下位金種制限)
		amount, exchangeTbl, err := c.reverseExchange.SpecifyExchangeWithLowerDenominationLimit(texCon, inReq.Amount)
		if err != nil {
			c.handleError(true, usecases.SYSLOG_LOGTYPE_REQUEST_REVERSE_EXCHNGE_CALCULATION_FATAL, err.Error(), inReq.RequestInfo)
			return
		}
		inResult.TargetAmount = amount
		inResult.TargetExCountTbl = nil
		inResult.ExchangeExCountTbl = exchangeTbl
	default: // その他(エラーResult)
		c.handleError(true, usecases.SYSLOG_LOGTYPE_REQUEST_REVERSE_EXCHNGE_CALCULATION_FATAL, "指定された逆両替種別が対象外", inReq.RequestInfo)
	}

	outResult, err := json.Marshal(inResult)
	if err != nil {
		c.handleError(true, usecases.SYSLOG_LOGTYPE_REQUEST_REVERSE_EXCHNGE_CALCULATION_FATAL, "指定された逆両替種別が対象外", inReq.RequestInfo)
		return
	}

	c.mqtt.Publish(c.resTopic, string(outResult))
}

func (c *reverseExchangeCalculation) handleError(resultSend bool, code int, message string, req domain.RequestInfo) {
	c.syslogMng.NoticeSystemLog(code, "", "入出金管理")
	c.logger.Error(message)

	if resultSend {
		errorCode, errorDetail := c.errorManager.GetErrorInfo(usecases.ERROR_INSIDE)
		resultInfo := domain.ResultReverseExchangeCalculation{
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
