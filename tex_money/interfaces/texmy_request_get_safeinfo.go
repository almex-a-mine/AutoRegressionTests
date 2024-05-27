package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases"
	"time"
)

type getSafeInfo struct {
	mqtt         handler.MqttRepository
	logger       handler.LoggerRepository
	syslogMng    usecases.SyslogManager
	errorMng     usecases.ErrorManager
	safeInfoMng  usecases.SafeInfoManager
	texmyHandler usecases.TexMoneyHandlerRepository
}

func NewRequestGetSageInfo(mqtt handler.MqttRepository, logger handler.LoggerRepository, syslogMng usecases.SyslogManager, errorMng usecases.ErrorManager, safeInfoMng usecases.SafeInfoManager, texmyHandler usecases.TexMoneyHandlerRepository) GetSafeInfoRepository {
	return &getSafeInfo{
		mqtt:         mqtt,
		logger:       logger,
		syslogMng:    syslogMng,
		errorMng:     errorMng,
		safeInfoMng:  safeInfoMng,
		texmyHandler: texmyHandler,
	}
}

// 開始処理
func (c *getSafeInfo) Start() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_get_safeinfo")
	c.mqtt.Subscribe(topic, c.recvRequest)
}

// 停止処理
func (c *getSafeInfo) Stop() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_get_safeinfo")
	c.mqtt.Unsubscribe(topic)
}

// サービス制御要求検出
func (c *getSafeInfo) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

func (c *getSafeInfo) recvRequest(message string) {

	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_get_safeinfo",
	})

	c.logger.Trace("【%v】START:要求受信 request_get_safeinfo 金庫情報取得要求", texCon.GetUniqueKey())
	var reqInfo domain.RequestGetSafeInfo
	var resInfo domain.ResultGetSafeInfo
	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_GET_SAFEINFO_FATAL, "", "入出金管理")
		c.logger.Error("getSafeInfo recvRequest json.Unmarshal:%v", err)
		return
	}
	c.logger.Debug("【%v】- RequestID %v", texCon.GetUniqueKey(), reqInfo.RequestInfo.RequestID)

	// 出金完了時のシーケンスの為、500ms遅延させる
	// 出金完了時notice_amount後にDB更新するが
	// notice_outで、UIが次のシーケンスへ移行し
	// 1-2が、取得した情報をもとにDB更新を実行してしまう。
	// 1-3としては1-2更新前に出金情報をDBに保存したい。
	time.Sleep(500 * time.Millisecond)

	resInfo.RequestInfo = reqInfo.RequestInfo
	resInfo.Result = true
	// 売上金回収情報をセット
	resInfo.SalesComplete, resInfo.SalesCount = c.safeInfoMng.GetSalesInfo()
	// 回収操作回数をセット
	resInfo.CollectCount = c.safeInfoMng.GetCollectCount()
	sortInfoTbl := c.safeInfoMng.GetSafeInfo(texCon).SortInfoTbl
	for i := 0; i < 11; i++ {
		resInfo.InfoSafe.SortInfoTbl[i] = sortInfoTbl[i]
	}
	// 集計中入出金差引をセット
	amount, countTbl, exCountTbl := c.calculateAggregateBalance(texCon)
	resInfo.InfoSafe.SortInfoTbl[11] = domain.SortInfoTbl{
		SortType:   domain.AGGREGATE_WITHDRAWAL,
		Amount:     amount,
		CountTbl:   countTbl,
		ExCountTbl: exCountTbl,
	}

	// デバイス有高 91をセット
	resInfo.InfoSafe.SortInfoTbl[12] = c.safeInfoMng.GetDeviceCashAvailable(texCon)

	// デバイス有高‐論理有高 92をセット
	_, resInfo.InfoSafe.SortInfoTbl[13] = c.safeInfoMng.GetAvailableBalance(texCon)

	// 通常金種別状況をセット
	resInfo.InfoSafe.CurrentStatusTbl = c.texmyHandler.MakeCurrentStatusTbl(texCon)

	payment, err := json.Marshal(resInfo)
	if err != nil {
		c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_GET_SAFEINFO_FATAL, "", "入出金管理")
		c.logger.Error("【%v】- json.Unmarshal:%v", texCon.GetUniqueKey(), err)
		return
	}
	c.syslogMng.NoticeSystemLog(usecases.SYSLOG_LOGTYPE_REQUEST_GET_SAFEINFO_SUCCESS, "", "入出金管理")
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_get_safeinfo")
	c.mqtt.Publish(topic, string(payment))

	c.logger.Trace("【%v】END:要求受信 request_get_safeinfo 金庫情報取得要求", texCon.GetUniqueKey())
}

// 集計中入出金差引情報算出
func (c *getSafeInfo) calculateAggregateBalance(texCon *domain.TexContext) (amount int, countTbl [domain.CASH_TYPE_SHITEI]int, exCountTbl [domain.EXTRA_CASH_TYPE_SHITEI]int) {
	c.logger.Trace("【%v】START:集計中入出金差引情報算出", texCon.GetUniqueKey())

	beforeReplenishmentBalance := c.safeInfoMng.GetBeforeReplenishmentBalance()                              //処理前補充差引情報を取得
	afterReplenishmentBalance := c.safeInfoMng.GetSafeInfo(texCon).SortInfoTbl[domain.REPLENISHMENT_BALANCE] // 処理後(現在)補充差引情報を取得

	// 処理前補充差引 - 処理後補充差引 を算出
	// 金額
	amount = beforeReplenishmentBalance.Amount - afterReplenishmentBalance.Amount
	// 通常金種別枚数
	for i, v := range beforeReplenishmentBalance.CountTbl {
		countTbl[i] = v - afterReplenishmentBalance.CountTbl[i]
	}
	// 拡張金種別枚数
	for i, v := range beforeReplenishmentBalance.ExCountTbl {
		exCountTbl[i] = v - afterReplenishmentBalance.ExCountTbl[i]
	}

	c.logger.Trace("【%v】END:集計中入出金差引情報算出", texCon.GetUniqueKey())
	return
}
