package interfaces

import (
	"encoding/json"
	"fmt"
	"strconv"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/pkg/lib"
	"tex_money/usecases"
)

type printReport struct {
	mqtt             handler.MqttRepository
	logger           handler.LoggerRepository
	config           config.Configuration
	syslogMng        usecases.SyslogManager
	errorMng         usecases.ErrorManager
	printSendRecv    PrintSendRecvRepository
	paymentSendRecv  PaymentSendRecvRepository
	texmyHandler     usecases.TexMoneyHandlerRepository
	praReqInfo       domain.RequestPrintReport
	printDataManager usecases.PrintDataManager
	aggregateMng     usecases.AggregateManager
	safeInfoMng      usecases.SafeInfoManager
	resTopic         string
}

// 入出金レポート印刷要求
func NewRequestPrintReport(mqtt handler.MqttRepository,
	logger handler.LoggerRepository,
	config config.Configuration,
	syslogMng usecases.SyslogManager,
	errorMng usecases.ErrorManager,
	printSendRecv PrintSendRecvRepository,
	paymentSendRecv PaymentSendRecvRepository,
	texmyHandler usecases.TexMoneyHandlerRepository,
	printDataManager usecases.PrintDataManager,
	aggregateMng usecases.AggregateManager,
	safeInfoMng usecases.SafeInfoManager,
) PrintReportRepository {
	return &printReport{
		mqtt:             mqtt,
		logger:           logger,
		config:           config,
		syslogMng:        syslogMng,
		errorMng:         errorMng,
		paymentSendRecv:  paymentSendRecv,
		printSendRecv:    printSendRecv,
		texmyHandler:     texmyHandler,
		printDataManager: printDataManager,
		aggregateMng:     aggregateMng,
		safeInfoMng:      safeInfoMng,
		resTopic:         fmt.Sprintf("%s/%s", domain.TOPIC_TEXMONEY_BASE, "result_print_report")}
}

// 開始処理
func (c *printReport) Start() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_print_report")
	c.mqtt.Subscribe(topic, c.recvRequest)
}

// 停止処理
func (c *printReport) Stop() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_print_report")
	c.mqtt.Unsubscribe(topic)
}

// サービス制御要求検出
func (c *printReport) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

func (c *printReport) recvRequest(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_print_report",
	})

	c.logger.Trace("【%v】START:要求受信 request_print_report 入出金レポート印刷要求", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:要求受信 request_print_report 入出金レポート印刷要求", texCon.GetUniqueKey())

	err := json.Unmarshal([]byte(message), &c.praReqInfo)
	if err != nil {
		errorCode, errorDetail := c.errorMng.GetErrorInfo(usecases.ERROR_INSIDE)
		c.handleError(true, usecases.SYSLOG_LOGTYPE_REQUEST_PRINT_REPORT_FATAL, fmt.Sprintf("printReport recvRequest json.Unmarshal:%v", err), c.praReqInfo.RequestInfo, errorCode, errorDetail)
		return
	}
	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), c.praReqInfo.RequestInfo.RequestID)

	if err, errorCode, errorDetail := c.reqPrintReport(texCon, c.praReqInfo.FilePath, c.praReqInfo.ReportId); err != nil {
		c.handleError(true, usecases.SYSLOG_LOGTYPE_REQUEST_PRINT_REPORT_FATAL, fmt.Sprintf("【%v】printReport recvRequest:%v", texCon.GetUniqueKey(), err), c.praReqInfo.RequestInfo, errorCode, errorDetail)
		return
	}
}

// 入出金レポート印刷要求:戻り値　印刷制御:補充レシート要求
func (c *printReport) reqPrintReport(texCon *domain.TexContext, filepath string, reportId int) (error, string, string) {
	c.logger.Trace("【%v】START:texMoneyHandler reqPrintReport(filepath=%s, reportId=%v)", texCon.GetUniqueKey(), filepath, reportId)
	defer c.logger.Trace("【%v】END:texMoneyHandler reqPrintReport", texCon.GetUniqueKey())

	// 動作開始/終了時の金庫・レポート情報
	c.logger.Debug("【%v】レシート発行時点保持データ", texCon.GetUniqueKey())
	c.safeInfoMng.OutputLogSafeInfoExCountTbl(texCon)
	c.aggregateMng.OutputLogAggregateExCountTbl()

	// データセット
	var numInfo []int
	switch reportId {
	case domain.SUMMARY_SALES: //精算機別日計表
		c.texmyHandler.SetSequence(texCon, domain.PRINT_SUMMARY_SALES)
		numInfo = c.summarySales(texCon)
		if numInfo == nil {
			errorCode, errorDetail := c.errorMng.GetErrorInfo(usecases.ERROR_COMMUNICATION_FAIL)
			return fmt.Errorf("精算機別日計表失敗"), errorCode, errorDetail
		}
	case domain.REPORT_CASHCOUNT: //キャッシュカウントレポート
		numInfo = c.printDataManager.ReportCashCount(texCon)

	case domain.SUPPLY_BILL: //紙幣補充
		numInfo = c.printDataManager.ReportSupplyBill(texCon)

	case domain.CHANGE_COINTUNIT1, //通常硬貨ユニット交換
		domain.CHANGE_COINTUNIT2,    //予備硬貨ユニット交換
		domain.CHANGE_COINTUNIT_ALL, //全硬貨ユニット交換
		domain.SUPPLY_COIN_MANUAL:   //硬貨手動追加
		numInfo = c.printDataManager.ReportChangeCoinUnit(texCon, reportId)
	case domain.CASHSALES_COLLECT: //現金売上金回収
		numInfo = c.cashSalesCollectReport(texCon)
		if numInfo == nil {
			errorCode, errorDetail := c.errorMng.GetErrorInfo(usecases.ERROR_COMMUNICATION_FAIL)
			return fmt.Errorf("現金売上金回収レポート失敗"), errorCode, errorDetail
		}
	case domain.REPORT_COINUNIT: //硬貨ユニット補充差分レポート
		numInfo = c.printDataManager.ReportCoinUnitDiff(texCon)

	case domain.REPORT_SUMMARY: //補充レシート種別:精算機日計レシート（FIT-B NEXTクリニック向け）
		c.texmyHandler.SetSequence(texCon, domain.PRINT_REPORT_SUMMARY)
		numInfo = c.reportSummary(texCon)
		if numInfo == nil {
			errorCode, errorDetail := c.errorMng.GetErrorInfo(usecases.ERROR_COMMUNICATION_FAIL)
			return fmt.Errorf("精算機日計レシート失敗"), errorCode, errorDetail
		}
	case 19, 20, 21, 22: // 補充レポート（追加補充/回収庫から回収/指定枚数回収/逆両替）
		numInfo = c.printDataManager.ReportSupply(texCon)
	default:
		errorCode, errorDetail := c.errorMng.GetErrorInfo(usecases.ERROR_REPORT_ID_MISMATCH)
		return fmt.Errorf("リクエスト情報不正:サポート外のreportId(%v)", reportId), errorCode, errorDetail
	}

	// 現在日時を取得
	date, time, err := lib.GeDateTime()
	if err != nil {
		errorCode, errorDetail := c.errorMng.GetErrorInfo(usecases.ERROR_INSIDE)
		return fmt.Errorf("入出金レポート印刷要求 現在日時取得失敗 err=%v", err), errorCode, errorDetail
	}
	dateInt, _ := strconv.Atoi(date)
	timeInt, _ := strconv.Atoi(time)

	//補充レシート印刷要求情報作成
	reqInfo := domain.NewRequestSupply(c.texmyHandler.NewRequestInfo(texCon),
		domain.GetReportName(reportId),
		c.config.TermNo,
		dateInt,
		timeInt,
		numInfo)
	// 送信
	c.printSendRecv.SendRequestSupply(texCon, reqInfo)

	return nil, "", ""
}

// 精算機別日計表
func (c *printReport) summarySales(texCon *domain.TexContext) []int {
	c.logger.Trace("【%v】START:精算機別日計表", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:精算機別日計表", texCon.GetUniqueKey())

	/// payment から情報取得
	ok := c.getSalesInfo(texCon)
	if !ok {
		c.logger.Warn("【%v】END:printReport SummarySales Error", texCon.GetUniqueKey())
		return nil
	}

	return c.printDataManager.ReportSummarySales(texCon)

}

func (c *printReport) getSalesInfo(texCon *domain.TexContext) bool {
	// リクエスト生成
	outReq := &domain.RequestGetSalesinfo{
		RequestInfo: c.texmyHandler.NewRequestInfo(texCon),
	}
	// 受信チャネルを生成
	var resChan = make(chan interface{})
	// 外部リクエスト
	go c.paymentSendRecv.SendRequestGetSalesInfo(texCon, resChan, outReq)
	// 外部リクエスト受信
	salesMoney := <-resChan
	// エラーの場合もあるので、型チェックでOKなら次の処理に進む
	outResInfo, ok := salesMoney.(domain.ResultGetSalesinfo)
	if !ok { // 型チェックでエラーの場合
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_REVERSE_EXCHNGE_CALCULATION_FATAL, "", "入出金管理")
		err, ok := salesMoney.(error) // 型チェックエラー
		if ok {
			c.logger.Error(err.Error()) // エラーでセットされたメッセージをログ出力
		}
		return false
	}

	// 受信値を格納
	switch c.texmyHandler.GetSequence(texCon) {
	case domain.PRINT_REPORT_SUMMARY: // 精算機日計レシート
		c.printDataManager.SetReportSummarySalesInfo(outResInfo)
	default:
		c.printDataManager.SetSummarySales(outResInfo)
	}

	return true
}

// 現金売上金回収
func (c *printReport) cashSalesCollectReport(texCon *domain.TexContext) []int {
	c.logger.Trace("【%v】START:現金売上金回収レポート", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END::現金売上金回収レポート", texCon.GetUniqueKey())

	/// payment から情報取得
	ok := c.getSalesInfo(texCon)
	if !ok {
		c.logger.Warn("【%v】END:printReport SummarySales Error", texCon.GetUniqueKey())
		return nil
	}

	return c.printDataManager.CashSalesCollectReport(texCon)
}

// 精算機日計レシート（FIT-B NEXTクリニック向け）
func (c *printReport) reportSummary(texCon *domain.TexContext) []int {
	c.logger.Trace("【%v】START:精算機日計レシートデータ作成", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:精算機日計レシートデータ作成", texCon.GetUniqueKey())
	/// payment から売上情報取得
	ok := c.getSalesInfo(texCon)
	if !ok {
		c.logger.Warn("【%v】- 精算機日計レシート売上情報取得 Error", texCon.GetUniqueKey())
		return nil
	}

	return c.printDataManager.ReportSummary(texCon)

}

// 処理結果応答
func (c *printReport) SendResult(texCon *domain.TexContext, qresInfo domain.ResultSupply) {
	c.logger.Trace("【%v】START:printReport SendResult c.praReqInfo.RequestInfo=%+v", texCon.GetUniqueKey(), c.praReqInfo.RequestInfo)
	defer c.logger.Trace("【%v】END:printReport SendResult", texCon.GetUniqueKey())

	if qresInfo.PrintId == "" || !qresInfo.Result {
		c.handleError(true, usecases.SYSLOG_LOGTYPE_REQUEST_PRINT_REPORT_FATAL, fmt.Sprintf("printReport SendResult Error:%s", "印刷要求応答情報エラー"), c.praReqInfo.RequestInfo, qresInfo.ErrorCode, qresInfo.ErrorDetail)
		return
	}

	resInfo := &domain.ResultPrintReport{
		RequestInfo: c.praReqInfo.RequestInfo,
		Result:      true,
		SlipPrintId: qresInfo.PrintId,
	}

	payment, err := json.Marshal(resInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_PRINT_REPORT_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Marshal:%v", texCon.GetUniqueKey(), err)
		return
	}

	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_PRINT_REPORT_SUCCESS, "", "入出金管理")
	c.mqtt.Publish(c.resTopic, string(payment))
}

func (c *printReport) handleError(resultSend bool, code int, message string, req domain.RequestInfo, errorCode string, errorDetail string) {
	c.syslogMng.NoticeSystemLog(code, "", "入出金管理")
	c.logger.Error(message)

	if resultSend {
		resultInfo := domain.ResultPrintReport{
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
