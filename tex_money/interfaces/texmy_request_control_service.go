package interfaces

import (
	"encoding/json"
	"fmt"
	"sync"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases"
)

type controlService struct {
	mqtt                       handler.MqttRepository
	logger                     handler.LoggerRepository
	config                     config.Configuration
	syslogMng                  usecases.SyslogManager
	errorMng                   usecases.ErrorManager
	svcContTbl                 []Repository
	texMoneyNoticeManager      usecases.TexMoneyNoticeManagerRepository
	sendRecv                   SendRecvRepository                     //現金入出金機制御通信
	texdtSendRecv              TexdtSendRecvRepository                //稼働データ通信
	statusTxSendRecv           StatusTxSendRecvRepository             //精算機状態管理通信
	printSendRecv              PrintSendRecvRepository                //印刷制御通信
	moneyInit                  MoneyInitRepository                    //初期補充要求監視
	moneyExchange              MoneyExchangeRepository                //両替要求
	moneyAdd                   MoneyAddReplenishRepository            //追加補充要求
	moneyCol                   MoneyCollectRepository                 //回収要求(途中回収要求／全回収要求／売上金回収要求)
	setAmount                  SetAmountRepository                    //現在枚数変更要求
	statusCashRepo             StatusCashRepository                   //現金入出金機制御ステータス要求
	payCash                    PayCashRepository                      //取引入金要求
	outCash                    OutCashRepository                      //取引出金要求
	amountCash                 AmountCashRepository                   //有高枚数要求
	printReport                PrintReportRepository                  //入出金レポート印刷要求
	salesInfo                  SalesInfoRepository                    //売上金情報要求
	clearCashInfo              RequestClearCashInfoRepository         //入出金データクリア要求
	changeSystemOpe            RequestChangeSystemOperationRepository //システム動作モード変更要求
	maintenanceMode            MaintenanceModeRepository              //保守業務モード要求
	getSafeInfo                GetSafeInfoRepository                  //金庫情報取得要求
	registerMoneySetting       RegisterMoneySettingRepository         //金銭設定登録要求
	getMoneySetting            GetMoneySettingRepository              //金銭設定取得要求
	coinCassetteControl        CoinCassetteControlRepository          //硬貨カセット操作要求
	paymentSendRecv            PaymentSendRecvRepository              //取引管理
	reverseExchangeCalculation ReverseExchangeCalculationRepository   //逆両替算出
	scrutiny                   ScrutinyRepository                     //精査モード要求
	texmyHandler               usecases.TexMoneyHandlerRepository
}

// 実行制御要求
func NewRequestControlService(mqtt handler.MqttRepository,
	logger handler.LoggerRepository,
	config config.Configuration,
	syslogMng usecases.SyslogManager,
	errorMng usecases.ErrorManager,
	svcContTbl []Repository,
	texMoneyNoticeManager usecases.TexMoneyNoticeManagerRepository,
	sendRecv SendRecvRepository,
	texdtSendRecv TexdtSendRecvRepository,
	statusTxSendRecv StatusTxSendRecvRepository,
	printSendRecv PrintSendRecvRepository,
	moneyInit MoneyInitRepository,
	moneyExchange MoneyExchangeRepository,
	moneyAdd MoneyAddReplenishRepository,
	moneyCol MoneyCollectRepository,
	setAmount SetAmountRepository,
	statusCashRepo StatusCashRepository,
	payCash PayCashRepository,
	outCash OutCashRepository,
	amountCash AmountCashRepository,
	printReport PrintReportRepository,
	salesInfo SalesInfoRepository,
	clearCashInfo RequestClearCashInfoRepository,
	changeSystemOpe RequestChangeSystemOperationRepository,
	maintenanceMode MaintenanceModeRepository,
	registerMoneySetting RegisterMoneySettingRepository,
	getMoneySetting GetMoneySettingRepository,
	getSafeInfo GetSafeInfoRepository,
	coinCassetteControl CoinCassetteControlRepository,
	paymentSendRecv PaymentSendRecvRepository,
	reverseExchangeCalculation ReverseExchangeCalculationRepository,
	scrutiny ScrutinyRepository,
	texmyHandler usecases.TexMoneyHandlerRepository) controlServiceRepository {
	return &controlService{
		mqtt:                       mqtt,
		logger:                     logger,
		config:                     config,
		syslogMng:                  syslogMng,
		errorMng:                   errorMng,
		svcContTbl:                 svcContTbl,
		texMoneyNoticeManager:      texMoneyNoticeManager,
		sendRecv:                   sendRecv,
		texdtSendRecv:              texdtSendRecv,
		statusTxSendRecv:           statusTxSendRecv,
		printSendRecv:              printSendRecv,
		moneyInit:                  moneyInit,
		moneyExchange:              moneyExchange,
		moneyAdd:                   moneyAdd,
		moneyCol:                   moneyCol,
		setAmount:                  setAmount,
		statusCashRepo:             statusCashRepo,
		payCash:                    payCash,
		outCash:                    outCash,
		amountCash:                 amountCash,
		printReport:                printReport,
		salesInfo:                  salesInfo,
		clearCashInfo:              clearCashInfo,
		changeSystemOpe:            changeSystemOpe,
		maintenanceMode:            maintenanceMode,
		getSafeInfo:                getSafeInfo,
		registerMoneySetting:       registerMoneySetting,
		getMoneySetting:            getMoneySetting,
		coinCassetteControl:        coinCassetteControl,
		paymentSendRecv:            paymentSendRecv,
		reverseExchangeCalculation: reverseExchangeCalculation,
		scrutiny:                   scrutiny,
		texmyHandler:               texmyHandler}
}

// 開始処理
func (c *controlService) Start() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_control_service")
	c.mqtt.Subscribe(topic, c.recvResuestControlService)
}

// 停止処理
func (c *controlService) Stop() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_control_service")
	c.mqtt.Unsubscribe(topic)
}

// サービス制御要求検出
func (c *controlService) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

// サービス実行状態制御要求検出
func (c *controlService) recvResuestControlService(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_control_service",
	})
	c.logger.Trace("【%v】START:要求受信 controlService recvResuestControlService", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:要求受信 controlService recvResuestControlService", texCon.GetUniqueKey())
	var reqInfo domain.RequestControlService

	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.logger.Error("controlService recvResuestControlService json.Unmarshal:%v", err)
		return
	}
	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), reqInfo.RequestInfo.RequestID)

	res := domain.ResultControlService{
		RequestInfo: reqInfo.RequestInfo,
		Result:      true,
	}
	c.resultControlService(texCon, &res)

	c.logger.Debug("【%v】controlService recvResuestControlService START:各インターフェースのサービス制御呼び出し", texCon.GetUniqueKey())
	//各インターフェースのサービス制御呼び出す
	c.sendRecv.ControlService(reqInfo)                   //現金入出金機制御送信管理
	c.texdtSendRecv.ControlService(reqInfo)              //稼働データ送受信管理
	c.statusTxSendRecv.ControlService(reqInfo)           //精算機状態管理通信
	c.printSendRecv.ControlService(reqInfo)              //印刷制御通信
	c.moneyInit.ControlService(reqInfo)                  //初期補充要求監視
	c.moneyExchange.ControlService(reqInfo)              //両替要求
	c.moneyAdd.ControlService(reqInfo)                   //追加補充要求
	c.moneyCol.ControlService(reqInfo)                   //回収要求(途中回収要求／全回収要求／売上金回収要求)
	c.setAmount.ControlService(reqInfo)                  //現在枚数変更要求
	c.statusCashRepo.ControlService(reqInfo)             //現金入出金機制御ステータス要求
	c.payCash.ControlService(reqInfo)                    //取引入金要求
	c.outCash.ControlService(reqInfo)                    //取引出金要求
	c.amountCash.ControlService(reqInfo)                 //有高枚数要求
	c.printReport.ControlService(reqInfo)                //入出金レポート印刷要求
	c.salesInfo.ControlService(reqInfo)                  //売上金情報要求
	c.clearCashInfo.ControlService(reqInfo)              //入出金データクリア要求
	c.changeSystemOpe.ControlService(reqInfo)            //システム動作モード変更要求
	c.maintenanceMode.ControlService(reqInfo)            //保守業務モード要求
	c.getSafeInfo.ControlService(reqInfo)                //金庫情報取得要求
	c.registerMoneySetting.ControlService(reqInfo)       //金銭設定登録要求
	c.getMoneySetting.ControlService(reqInfo)            //金銭設定取得要求
	c.coinCassetteControl.ControlService(reqInfo)        //硬貨カセット操作要求
	c.paymentSendRecv.ControlService(reqInfo)            //取引管理
	c.reverseExchangeCalculation.ControlService(reqInfo) //逆両替算出
	c.scrutiny.ControlService(reqInfo)                   //精査モード要求
	c.logger.Debug("【%v】controlService recvResuestControlService END:各インターフェースのサービス制御呼び出し", texCon.GetUniqueKey())

	//再起動時、layer2へ有高要求を投げて有高ステータスを更新する
	if reqInfo.StatusService {
		// 釣銭不一致監視ON
		c.texmyHandler.InitialDiscrepantOn(true)
		c.sendRecv.InitializeCashCtrlFlagOn(true)

		var wg sync.WaitGroup
		wg.Add(1)
		go c.sendRecv.InitializeCashctl(&wg)
		wg.Wait()
	} else {
		// 釣銭不一致監視OFF
		c.texmyHandler.InitialDiscrepantOn(false)
	}

	// サービス実行状態更新
	ok := c.texMoneyNoticeManager.UpdateStatusServiceData(texCon,
		domain.StatusService{
			StatusService: reqInfo.StatusService,
		})
	// サービス実行状況を通知
	if ok {
		noticeInfo := c.texMoneyNoticeManager.GetStatusServiceData(texCon)
		c.NoticeControlService(&noticeInfo)
	}
}

// 処理結果応答
func (c *controlService) resultControlService(texCon *domain.TexContext, pResInfo *domain.ResultControlService) {
	c.logger.Trace("【%v】START:controlService resultControlService", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:controlService resultControlService", texCon.GetUniqueKey())

	payment, err := json.Marshal(pResInfo)
	if err != nil {
		c.logger.Error("【%v】- json.Marshal:%v", texCon.GetUniqueKey(), err)
		return
	}
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_control_service")
	c.mqtt.Publish(topic, string(payment))

}

// 実行状況通知
func (c *controlService) NoticeControlService(pNoticeInfo *domain.StatusService) {
	c.logger.Trace("START:controlService noticeControlService")
	defer c.logger.Trace("END:controlService noticeControlService")
	payment, err := json.Marshal(pNoticeInfo)
	if err != nil {
		c.logger.Error("controlService NoticeControlService json.Marshal:%v", err)
		return
	}
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "notice_status_service")
	c.mqtt.Publish(topic, string(payment))

}
