package infrastructure

import (
	"sync"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/interfaces"
	"tex_money/pkg/pc"
	"tex_money/usecases"
	"time"
)

var wgrps sync.WaitGroup
var waitManager usecases.IWait

// 開始処理
func Router() {
	wgrps = sync.WaitGroup{}
	wgrps.Add(1)
	//設定取得
	config := config.Initialize(domain.AppName)
	//ログハンドラ作成
	logger := NewLogger(config.SystemConf.MaxLength,
		config.SystemConf.MaxRotation,
		config.SystemConf.LogStopInfo,
		config.SystemConf.LogStopTrace,
		config.SystemConf.LogStopMqtt,
		config.SystemConf.LogStopDebug,
		config.SystemConf.LogStopMutex,
		config.SystemConf.LogStopWarn,
		config.SystemConf.LogStopError,
		config.SystemConf.LogStopFatal,
		config.SystemConf.LogStopSequence)
	logger.Info("Program Start")
	logger.Info("config %+v", config)

	//MQTTハンドラ作成
	mqtt := NewMQTTHandler(logger, config.MqttConf.TCP, config.MqttConf.Port, domain.AppName)
	sysMqtt := NewMQTTHandler(logger, config.MqttConf.TCP, config.MqttConf.Port, domain.AppName+"_sys")

	// Start Contoller
	Controller(mqtt, sysMqtt, logger, config)
}

// 停止処理
func RouterStop() {
	wgrps.Done()
}

func Controller(mqtt handler.MqttRepository, sysMqtt handler.MqttRepository, logger handler.LoggerRepository, config config.Configuration) {
	statusService := config.SystemConf.StartUpStatus == 0
	logger.Info("statusService = %t", statusService)

	funcTbl := make([]interfaces.Repository, 0)

	//通知管理
	noticeManager := usecases.NewTexMoneyNoticeManager(logger)

	//システムログ管理
	syslogManager := usecases.NewSysLogManager(mqtt)
	//エラー管理
	errorManager := usecases.NewErrorManager()
	//iniファイル管理
	iniService := usecases.NewIniService(logger)
	//金庫情報管理
	safeInfoManager := usecases.NewSafeInfoManager(logger, config, syslogManager, errorManager, iniService)
	//レポート用金庫情報管理
	aggregateManager := usecases.NewAggregateManager(logger, config, syslogManager, errorManager, safeInfoManager, iniService)

	//保守業務モード要求管理
	maintenanceModeManager := usecases.NewMaintenanceModeManager(logger, config, aggregateManager, safeInfoManager, iniService)

	//入出金管理ハンドラ実装
	texmyHandler := usecases.NewTexmyHandler(
		logger,
		config,
		errorManager,
		safeInfoManager,
		aggregateManager,
		noticeManager,
		maintenanceModeManager,
	)

	//実行状態取得要求監視
	requestStatSvc := interfaces.NewRequestGetService(sysMqtt, logger, noticeManager)
	//requestStatSvc.Start() // 108.8.0.9 起動シーケンス完了後に移動

	//システム動作モード変更要求
	requestChangeSystemOperation := interfaces.NewRequestChangeSystemOperation(sysMqtt, logger, syslogManager, noticeManager)
	//requestChangeSystemOperation.Start() // 108.8.0.9 起動シーケンス完了後に移動

	//硬貨カセット操作管理
	coinCassetteControlManager := usecases.NewCoinCassetteControlManager(logger, safeInfoManager, config, texmyHandler)

	//レシート管理
	printDataManager := usecases.NewPrintDataManager(logger, config, syslogManager, errorManager, aggregateManager, safeInfoManager, coinCassetteControlManager, maintenanceModeManager, texmyHandler)

	//逆両替要求管理
	reverseExchangeCalculationManager := usecases.NewReverseExchangeCalculationManager(logger, config, safeInfoManager, texmyHandler)

	// 状態変更管理
	changeStatus := usecases.NewChangeStatus(logger, safeInfoManager, noticeManager, texmyHandler)

	//MQTTの接続を待機
	waitMqttConnected(mqtt)

	//待機クラス生成
	waitManager = usecases.NewWaitManager(logger)

	//稼働データ通信
	texdtSendController := interfaces.NewTexdtSendRecv(mqtt, logger, errorManager, texmyHandler, waitManager)
	funcTbl = append(funcTbl, texdtSendController)
	if statusService {
		texdtSendController.Start()
	}

	//現金入出金機制御通信
	sendController := interfaces.NewSendRecv(mqtt, logger, errorManager, texmyHandler, waitManager, safeInfoManager, noticeManager)
	funcTbl = append(funcTbl, sendController)
	if statusService {
		sendController.InitializeCashCtrlFlagOn(true)
		sendController.Start()
	}

	//精算機状態管理通信
	statusTXController := interfaces.NewStatusTXSendRecv(mqtt, logger, errorManager, texmyHandler, waitManager)
	funcTbl = append(funcTbl, statusTXController)
	if statusService {
		statusTXController.Start()
	}

	//印刷制御通信
	printController := interfaces.NewPrintSendRecv(mqtt, logger, config, syslogManager, errorManager, texmyHandler, waitManager)
	funcTbl = append(funcTbl, printController)
	if statusService {
		printController.Start()
	}

	//入出金機制御
	//初期補充要求監視
	requestMoneyInitController := interfaces.NewRequestMoneyInit(mqtt, logger, syslogManager, sendController, texdtSendController, statusTXController, texmyHandler, safeInfoManager, changeStatus, noticeManager)
	sendController.SetAddressMoneyIni(requestMoneyInitController)      //初期補充要のアドレスを//現金入出金機制御送受信に渡す
	texdtSendController.SetAddressMoneyIni(requestMoneyInitController) //初期補充要のアドレスを//稼働データ管理送受信に渡す
	statusTXController.SetAddressMoneyIni(requestMoneyInitController)  //初期補充要のアドレスを//精算機状態管理送受信に渡す
	printController.SetAddressMoneyIni(requestMoneyInitController)     //初期補充要のアドレスを//印刷制御送受信に渡す
	funcTbl = append(funcTbl, requestMoneyInitController)
	if statusService {
		requestMoneyInitController.Start()
	}

	// 精算取引管理要求
	paymentSendRecv := interfaces.NewPaymentSendRecv(mqtt, logger, config, syslogManager, errorManager, waitManager)
	if statusService {
		paymentSendRecv.Start()
	}

	//両替要求
	requestMoneyExchangeController := interfaces.NewRequestMoneyExchange(mqtt, logger, config, syslogManager, errorManager, sendController, texdtSendController, statusTXController, printController, texmyHandler, printDataManager, aggregateManager, noticeManager, maintenanceModeManager, changeStatus, paymentSendRecv)
	sendController.SetAddressMoneyExchange(requestMoneyExchangeController)      //両替要求のアドレスを//現金入出金機制御送受信に渡す
	texdtSendController.SetAddressMoneyExchange(requestMoneyExchangeController) //両替要求のアドレスを//稼働データ管理送受信に渡す
	statusTXController.SetAddressMoneyExchange(requestMoneyExchangeController)  //両替要求のアドレスを//精算機状態管理送受信に渡す
	printController.SetAddressMoneyExchange(requestMoneyExchangeController)     //両替要求のアドレスを//印刷制御送受信に渡す
	funcTbl = append(funcTbl, requestMoneyExchangeController)
	if statusService {
		requestMoneyExchangeController.Start()
	}

	//追加補充要求
	requestMoneyAddReplenishController := interfaces.NewRequestMoneyAddReplenish(mqtt, logger, config, syslogManager, errorManager, sendController, texdtSendController, statusTXController, printController, texmyHandler, changeStatus)
	sendController.SetAddressMoneyAddReplenish(requestMoneyAddReplenishController)      //追加補充要求のアドレスを//現金入出金機制御送受信に渡す
	texdtSendController.SetAddressMoneyAddReplenish(requestMoneyAddReplenishController) //追加補充要求のアドレスを//稼働データ管理送受信に渡す
	statusTXController.SetAddressMoneyAddReplenish(requestMoneyAddReplenishController)  //追加補充要求のアドレスを//精算機状態管理送受信に渡す
	printController.SetAddressMoneyAddReplenish(requestMoneyAddReplenishController)     //追加補充要求のアドレスを//印刷制御送受信に渡す
	funcTbl = append(funcTbl, requestMoneyAddReplenishController)
	if statusService {
		requestMoneyAddReplenishController.Start()
	}

	//回収要求(途中回収要求／全回収要求／売上金回収要求)
	requestMoneyCollectController := interfaces.NewRequestMoneyCollect(mqtt, logger, syslogManager, errorManager, sendController, texdtSendController, statusTXController, texmyHandler, safeInfoManager, changeStatus)
	sendController.SetAddressMoneyCollect(requestMoneyCollectController)      //回収要求のアドレスを//現金入出金機制御送受信に渡す
	texdtSendController.SetAddressMoneyCollect(requestMoneyCollectController) //回収要求のアドレスを//稼働データ管理送受信に渡す
	statusTXController.SetAddressMoneyCollect(requestMoneyCollectController)  //回収要求のアドレスを//精算機状態管理送受信に渡す
	printController.SetAddressMoneyCollect(requestMoneyCollectController)     //回収要求のアドレスを//印刷制御送受信に渡す
	funcTbl = append(funcTbl, requestMoneyCollectController)
	if statusService {
		requestMoneyCollectController.Start()
	}

	//現在枚数変更要求
	requestSetAmountController := interfaces.NewRequestSetAmount(mqtt, logger, config, syslogManager, errorManager, sendController, texdtSendController, statusTXController, printController, texmyHandler, printDataManager, changeStatus, noticeManager, safeInfoManager)
	sendController.SetAddressSetAmount(requestSetAmountController)      //現在枚数変更要求のアドレスを//現金入出金機制御送受信に渡す
	texdtSendController.SetAddressSetAmount(requestSetAmountController) //現在枚数変更要求のアドレスを//稼働データ管理送受信に渡す
	statusTXController.SetAddressSetAmount(requestSetAmountController)  //現在枚数変更要求のアドレスを//精算機状態管理送受信に渡す
	printController.SetAddressSetAmount(requestSetAmountController)     //現在枚数変更要求のアドレスを//印刷制御送受信に渡す
	funcTbl = append(funcTbl, requestSetAmountController)
	if statusService {
		requestSetAmountController.Start()
	}

	//現金入出金機制御ステータス要求
	requestStatusCashController := interfaces.NewRequestStatusCash(mqtt, logger, syslogManager, noticeManager)
	sendController.SetAddressStatusCash(requestStatusCashController)      //現金入出金機制御ステータス要求のアドレスを//現金入出金機制御送受信に渡す
	texdtSendController.SetAddressStatusCash(requestStatusCashController) //現金入出金機制御ステータス要求のアドレスを//稼働データ管理送受信に渡す
	statusTXController.SetAddressStatusCash(requestStatusCashController)  //現金入出金機制御ステータス要求のアドレスを//精算機状態管理送受信に渡す
	printController.SetAddressStatusCash(requestStatusCashController)     //現金入出金機制御ステータス要求のアドレスを//印刷制御送受信に渡す
	funcTbl = append(funcTbl, requestStatusCashController)
	if statusService {
		requestStatusCashController.Start()
	}

	//取引入金要求
	requestPayCashController := interfaces.NewRequestPayCash(mqtt, logger, config, syslogManager, errorManager, sendController, texdtSendController, statusTXController, printController, texmyHandler)
	sendController.SetAddressPayCash(requestPayCashController)      //取引入金要求のアドレスを現金入出金機制御送受信に渡す
	texdtSendController.SetAddressPayCash(requestPayCashController) //取引入金要求のアドレスを稼働データ管理送受信に渡す
	statusTXController.SetAddressPayCash(requestPayCashController)  //取引入金要求のアドレスを精算機状態管理送受信に渡す
	printController.SetAddressPayCash(requestPayCashController)     //取引入金要求のアドレスを印刷制御送受信に渡す
	funcTbl = append(funcTbl, requestPayCashController)
	if statusService {
		requestPayCashController.Start()
	}

	//取引出金要求
	requestOutCashController := interfaces.NewRequestOutCash(mqtt, logger, config, syslogManager, errorManager, sendController, texdtSendController, statusTXController, printController, texmyHandler, changeStatus)
	sendController.SetAddressOutCash(requestOutCashController)      //取引出金要求のアドレスを//現金入出金機制御送受信に渡す
	texdtSendController.SetAddressOutCash(requestOutCashController) //取引出金要求のアドレスを//稼働データ管理送受信に渡す
	statusTXController.SetAddressOutCash(requestOutCashController)  //取引出金要求のアドレスを//精算機状態管理送受信に渡す
	printController.SetAddressOutCash(requestOutCashController)     //取引出金要求のアドレスを//印刷制御送受信に渡す
	funcTbl = append(funcTbl, requestOutCashController)
	if statusService {
		requestOutCashController.Start()
	}

	//有高枚数要求
	requestAmoutCashController := interfaces.NewRequestAmountCash(mqtt, logger, config, syslogManager, errorManager, sendController, texdtSendController, statusTXController, printController, texmyHandler)
	sendController.SetAddressAmountCash(requestAmoutCashController)      //有高枚数要求のアドレスを//現金入出金機制御送受信に渡す
	texdtSendController.SetAddressAmountCash(requestAmoutCashController) //有高枚数要求のアドレスを//稼働データ管理送受信に渡す
	statusTXController.SetAddressAmountCash(requestAmoutCashController)  //有高枚数要求のアドレスを//精算機状態管理送受信に渡す
	printController.SetAddressAmountCash(requestAmoutCashController)     //有高枚数要求のアドレスを//印刷制御送受信に渡す
	funcTbl = append(funcTbl, requestAmoutCashController)
	if statusService {
		requestAmoutCashController.Start()
	}

	//入出金レポート印刷要求
	requestPrintReportController := interfaces.NewRequestPrintReport(mqtt, logger, config, syslogManager, errorManager, printController, paymentSendRecv, texmyHandler, printDataManager, aggregateManager, safeInfoManager)
	sendController.SetAddressPrintReport(requestPrintReportController)      // 入出金レポート印刷のアドレスを//現金入出金機制御送受信に渡す
	texdtSendController.SetAddressPrintReport(requestPrintReportController) // 入出金レポート印刷のアドレスを//稼働データ管理送受信に渡す
	statusTXController.SetAddressPrintReport(requestPrintReportController)  // 入出金レポート印刷のアドレスを//精算機状態管理送受信に渡す
	printController.SetAddressPrintReport(requestPrintReportController)     // 入出金レポート印刷のアドレスを//印刷制御送受信に渡す
	funcTbl = append(funcTbl, requestPrintReportController)
	if statusService {
		requestPrintReportController.Start()
	}

	//売上金情報要求
	requestSalesInfoController := interfaces.NewRequestSalesInfo(mqtt, logger, config, syslogManager, errorManager, sendController, texdtSendController, statusTXController, printController, texmyHandler, safeInfoManager)
	sendController.SetAddressSalesInfo(requestSalesInfoController)      // 売上金情報のアドレスを//現金入出金機制御送受信に渡す
	texdtSendController.SetAddressSalesInfo(requestSalesInfoController) // 売上金情報のアドレスを//稼働データ管理送受信に渡す
	statusTXController.SetAddressSalesInfo(requestSalesInfoController)  // 売上金情報のアドレスを//精算機状態管理送受信に渡す
	printController.SetAddressSalesInfo(requestSalesInfoController)     // 売上金情報のアドレスを//印刷制御送受信に渡す
	funcTbl = append(funcTbl, requestSalesInfoController)
	if statusService {
		requestSalesInfoController.Start()
	}

	//入出金データクリア要求
	requestClearCashInfoController := interfaces.NewRequestClearCashInfo(mqtt, logger, config, syslogManager, errorManager, sendController, texdtSendController, statusTXController, printController, texmyHandler, safeInfoManager)
	sendController.SetAddressClearCashInfo(requestClearCashInfoController)      //入出金データクリア要求のアドレスを//現金入出金機制御送受信に渡す
	texdtSendController.SetAddressClearCashInfo(requestClearCashInfoController) //入出金データクリア要求のアドレスを//稼働データ管理送受信に渡す
	statusTXController.SetAddressClearCashInfo(requestClearCashInfoController)  //入出金データクリア要求のアドレスを//精算機状態管理送受信に渡す
	printController.SetAddressClearCashInfo(requestClearCashInfoController)     //入出金データクリア要求のアドレスを//印刷制御送受信に渡す
	funcTbl = append(funcTbl, requestClearCashInfoController)
	if statusService {
		requestClearCashInfoController.Start()
	}

	// 硬貨カセット操作要求
	coinCassetteControl := interfaces.NewCoinCassetteControl(mqtt, logger, syslogManager, config, errorManager, coinCassetteControlManager, requestSetAmountController, aggregateManager, texmyHandler)
	if statusService {
		coinCassetteControl.Start()
	}

	// 逆両替計算
	reverseExchangeCalculation := interfaces.NewReverseExchangeCalculation(mqtt, logger, config, syslogManager, errorManager, texmyHandler, reverseExchangeCalculationManager, paymentSendRecv)
	if statusService {
		reverseExchangeCalculation.Start()
	}

	//保守業務モード要求
	requestMaintenanceModeController := interfaces.NewRequestMaintenanceMode(mqtt, logger, syslogManager, errorManager, aggregateManager, safeInfoManager, maintenanceModeManager)
	funcTbl = append(funcTbl, requestMaintenanceModeController)
	if statusService {
		requestMaintenanceModeController.Start()
	}

	// 金庫情報取得要求
	requestGetSageInfoController := interfaces.NewRequestGetSageInfo(mqtt, logger, syslogManager, errorManager, safeInfoManager, texmyHandler)
	if statusService {
		requestGetSageInfoController.Start()
	}

	// 金銭設定登録要求
	requestRegisterMoneySettingController := interfaces.NewRequestRegisterMoneySetting(mqtt, logger, config, syslogManager, errorManager, texmyHandler, iniService, noticeManager)
	if statusService {
		requestRegisterMoneySettingController.Start()
	}

	// 金銭設定取得要求
	requestGetMoneySettingController := interfaces.NewRequestGetMoneySetting(mqtt, logger, syslogManager, errorManager, texmyHandler)
	if statusService {
		requestGetMoneySettingController.Start()
	}

	// 精査モード要求
	requestScrutinyController := interfaces.NewScrutiny(mqtt, logger, syslogManager, config, errorManager, texmyHandler, sendController, statusTXController)
	if statusService {
		requestScrutinyController.Start()
	}

	//制御開始
	controller := interfaces.NewController(mqtt, logger, syslogManager, texmyHandler)
	controller.Start()

	//実行制御要求
	requestContSvc := interfaces.NewRequestControlService(
		sysMqtt,
		logger,
		config,
		syslogManager,
		errorManager,
		funcTbl,
		noticeManager,
		sendController,
		texdtSendController,
		statusTXController,
		printController,
		requestMoneyInitController,
		requestMoneyExchangeController,
		requestMoneyAddReplenishController,
		requestMoneyCollectController,
		requestSetAmountController,
		requestStatusCashController,
		requestPayCashController,
		requestOutCashController,
		requestAmoutCashController,
		requestPrintReportController,
		requestSalesInfoController,
		requestClearCashInfoController,
		requestChangeSystemOperation,
		requestMaintenanceModeController,
		requestGetSageInfoController,
		requestRegisterMoneySettingController,
		requestGetMoneySettingController,
		coinCassetteControl,
		paymentSendRecv,
		reverseExchangeCalculation,
		requestScrutinyController,
		texmyHandler)

	// requestContSvc.Start() // 108.8.0.9 起動シーケンス完了後に移動

	//入出金管理
	texmyHandler.Start()

	// DBから現在端末情報を取得
	texdtSendController.InitialDbData()
	// DBから現在端末情報を取得するまで待つ
	texdtSendController.GetInitializeDbData()
	// 釣銭不一致監視ON
	texmyHandler.InitialDiscrepantOn(true)

	// 2.1に起動時の有高ステータスと入出金機ステータスを要求
	var wg sync.WaitGroup
	wg.Add(1) // WaitGroupに対して、待つべきゴルーチンの数を1増や
	go sendController.InitializeCashctl(&wg)
	wg.Wait() // すべてのゴルーチンがDone()を呼び出すまで待つ

	appVersion, ok := pc.GetAppVersion()
	if !ok {
		appVersion = "アプリ情報取得失敗"
	}

	logger.Debug("appVersion %s", appVersion)
	//system log
	syslogManager.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_SERVICESTART_SUCCESS, "", appVersion)

	// 準備完了としてステータス更新
	texCon := domain.NewTexContext(domain.RegisterTexContext{})
	noticeManager.UpdateStatusServiceData(texCon, domain.StatusService{StatusService: statusService})

	// 起動時最後に送信
	noticeInfo := noticeManager.GetStatusServiceData(texCon)
	requestContSvc.NoticeControlService(&noticeInfo)

	pc.GetAppVersion()

	requestStatSvc.Start()               // 108.8.0.9 起動シーケンス完了後に移動
	requestChangeSystemOperation.Start() // 108.8.0.9 起動シーケンス完了後に移動
	requestContSvc.Start()               // 108.8.0.9 起動シーケンス完了後に移動

	// 終了要求まで待機
	wgrps.Wait()

	//system log
	syslogManager.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_SERVICESTOP_SUCCESS, "", "入出金管理")
	time.Sleep(1000 * time.Millisecond)

	//通信
	sendController.Stop()      //現金入出金機制御通信
	texdtSendController.Stop() //稼働データ通信
	statusTXController.Stop()  //精算機状態管理通信
	printController.Stop()     //印刷制御通信
	//入出金機制御
	requestMoneyInitController.Stop()            //初期補充要求監視
	requestMoneyExchangeController.Stop()        //両替要求
	requestMoneyAddReplenishController.Stop()    //追加補充要求
	requestMoneyCollectController.Stop()         //回収要求(途中回収要求／全回収要求／売上金回収要求)
	requestSetAmountController.Stop()            //現在枚数変更要求
	requestStatusCashController.Stop()           //現金入出金機制御ステータス要求
	requestPayCashController.Stop()              //取引入金要求
	requestOutCashController.Stop()              //取引出金要求
	requestPrintReportController.Stop()          //入出金レポート印刷要求
	requestSalesInfoController.Stop()            //売上金情報要求
	requestMaintenanceModeController.Stop()      // 保守業務モード要求
	requestGetSageInfoController.Stop()          //金庫情報取得要求
	requestRegisterMoneySettingController.Stop() //金銭設定登録要求
	requestGetMoneySettingController.Stop()      //金銭設定取得要求
	controller.Stop()                            //制御開始
	requestStatSvc.Stop()                        //サービス状態取得要求監視
	requestContSvc.Stop()                        //サービス制御監視
	texmyHandler.Stop()                          //入出金管理
	coinCassetteControl.Stop()                   //硬貨カセット操作
	paymentSendRecv.Stop()                       //精算取引管理
	reverseExchangeCalculation.Stop()            //逆両替算出要求
	requestScrutinyController.Stop()             //精査モード要求
}

// MQTT接続まで待機
func waitMqttConnected(mqtt handler.MqttRepository) {
	for {
		if mqtt.InConnectionOpen() {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}
