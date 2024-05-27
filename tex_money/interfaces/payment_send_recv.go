package interfaces

import (
	"encoding/json"
	"errors"
	"sync"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases"
)

type (
	paymentSendRecv struct {
		mqtt                handler.MqttRepository
		logger              handler.LoggerRepository
		config              config.Configuration
		syslogMng           usecases.SyslogManager
		errorMng            usecases.ErrorManager
		waitManager         usecases.IWait
		recvControlMapMutex sync.Mutex
		recvControlMap      map[ControlID]*RecvControl
	}
	PaymentSendRecvRepository interface {
		Start()
		Stop()
		ControlService(reqInfo domain.RequestControlService)
		SendRequestGetSalesInfo(texCon *domain.TexContext, resChan chan interface{}, resInfo *domain.RequestGetSalesinfo) // 返却先情報をチャネルにセット、リクエスト情報をセット
	}

	RecvControl struct {
		resultTopic string           // 期待するresultのtopic名をセット
		channel     chan interface{} // 呼び出し元から取得したチャネルを登録
	}
	ControlID string
)

const ( // resultTopic
	resultSalesInfoTopic = domain.TOPIC_UNIFUNCPAYMENT_BASE + "/result_get_salesinfo"
)

func NewPaymentSendRecv(mqtt handler.MqttRepository,
	logger handler.LoggerRepository,
	config config.Configuration,
	syslogMng usecases.SyslogManager,
	errorMng usecases.ErrorManager,
	waitManager usecases.IWait,
) PaymentSendRecvRepository {

	return &paymentSendRecv{
		mqtt:           mqtt,
		logger:         logger,
		config:         config,
		syslogMng:      syslogMng,
		errorMng:       errorMng,
		waitManager:    waitManager,
		recvControlMap: make(map[ControlID]*RecvControl),
	}
}

// Start 開始
func (c *paymentSendRecv) Start() {
	c.mqtt.Subscribe(resultSalesInfoTopic, c.RecvRequestGetSalesInfo)
}

// Stop 終了
func (c *paymentSendRecv) Stop() {
	c.mqtt.Unsubscribe(resultSalesInfoTopic)
}

// ControlService サービス制御要求検出
func (c *paymentSendRecv) ControlService(reqInfo domain.RequestControlService) {
	if reqInfo.StatusService {
		c.Start()
	} else {
		c.Stop()
	}
}

// registerRecvControlMap MAPにチャネルを登録
// MAP操作に対するロック処理漏れが無いように、登録だけの定義を実装
func (c *paymentSendRecv) registerRecvControlMap(controlID ControlID, topic string, channel chan interface{}) {
	// 送信待機ロジック(同トピックに同時に送信しないようにする為)
	for {
		// 保持しているものがなければfor分を抜ける
		if _, ok := c.recvControlMap[controlID]; !ok {
			break
		} else {
			// 存在する場合，同Topicがなければ、for文を抜ける
			var b = false
			for _, v := range c.recvControlMap {
				if v.resultTopic == topic {
					b = true
					break
				}
			}
			if !b {
				break
			}
		}

	}

	c.recvControlMapMutex.Lock()
	defer c.recvControlMapMutex.Unlock()

	c.recvControlMap[controlID] = &RecvControl{resultTopic: topic, channel: channel}

}

// sendErrorReleaseRecvControlMap エラー発生時の処理
// 処理漏れしないように、エラー発生時に実行する
func (c *paymentSendRecv) sendErrorReleaseRecvControlMap(controlID ControlID) {
	// 今のロジックだと100%タイムアウトするので、既にMAPが存在しなければ正常終了したものを見なす。
	if _, ok := c.recvControlMap[controlID]; ok {
		c.recvControlMap[controlID].channel <- errors.New(string(controlID) + "連携時にエラーが発生") // 上位に失敗を通知

		c.recvControlMapMutex.Lock()
		defer c.recvControlMapMutex.Unlock()
		delete(c.recvControlMap, controlID) // map情報削除
	}
}

func (c *paymentSendRecv) makeControlID(processId, requestId string) ControlID {
	return ControlID(processId + "_" + requestId)
}

// 送信：request_get_salesinfo
func (c *paymentSendRecv) SendRequestGetSalesInfo(texCon *domain.TexContext, resChan chan interface{}, resInfo *domain.RequestGetSalesinfo) {
	var reqTopic = domain.TOPIC_UNIFUNCPAYMENT_BASE + "/request_get_salesinfo"
	var resTopic = resultSalesInfoTopic
	_ = c.sendRequest(texCon, resChan, resInfo, reqTopic, resTopic, &resInfo.RequestInfo)
}

// リクエスト送信
// resChan:返却チャネル情報、reqInfo:送信するリクエストのjson情報、reqTopic:接頭語付き送信topic名称、resTopic:接頭語付き応答topic名称、requestInfo:送信するリクエストのrequestInfo
func (c *paymentSendRecv) sendRequest(texCon *domain.TexContext, resChan chan interface{}, reqInfo interface{}, reqTopic string, resTopic string, requestInfo *domain.RequestInfo) error {
	c.logger.Trace("【%v】START:paymentSendRecv sendRequest", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:paymentSendRecv sendRequest", texCon.GetUniqueKey())

	// トピックの返信先登録
	controlID := MakeControlId(requestInfo.ProcessID, requestInfo.RequestID)
	/// MAPに返却チャネル情報を登録
	c.registerRecvControlMap(controlID, resTopic, resChan)

	send, err := json.Marshal(reqInfo)
	if err != nil {
		c.sendErrorReleaseRecvControlMap(controlID)
		c.logger.Error("【%v】- json.Marshal:%v", texCon.GetUniqueKey(), err)
		return err
	}

	// 待機情報作成
	waitInfo := c.waitManager.MakeWaitInfo(texCon, requestInfo.ProcessID, requestInfo.RequestID, reqInfo)
	// 不要になった待機情報は最後に削除
	defer c.waitManager.DelWaitInfo(texCon, requestInfo.ProcessID, requestInfo.RequestID)

	// 送信
	c.mqtt.Publish(reqTopic, string(send))

	//〇応答データ待機
	result := c.waitManager.WaitResultInfo(texCon, waitInfo, domain.WAIT_TIME)
	if !result {
		if _, ok := c.recvControlMap[controlID]; ok {
			c.logger.Debug("【%v】- wait.WaitResultInfo: timeout TOPIC=%v", texCon.GetUniqueKey(), resTopic)
			c.sendErrorReleaseRecvControlMap(controlID)
		}
		return nil
	}

	//〇待機データ取得
	_, ok := c.waitManager.GetWaitInfo(texCon, requestInfo.ProcessID, requestInfo.RequestID)
	if !ok {
		c.logger.Debug("【%v】- wait.GetWaitInfo error", texCon.GetUniqueKey())
		c.sendErrorReleaseRecvControlMap(controlID)
		return nil
	}

	return nil
}

// 受信：result_get_salesinfo
func (c *paymentSendRecv) RecvRequestGetSalesInfo(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "result_get_salesinfo",
	})
	var topic = resultSalesInfoTopic
	c.logger.Trace("【%v】START:paymentSendRecv SendRequestGetSalesInfo", texCon.GetUniqueKey())
	var outResInfo domain.ResultGetSalesinfo
	err := json.Unmarshal([]byte(message), &outResInfo)
	if err != nil {
		c.logger.Error("paymentSendRecv SendRequestGetSalesInfo json.Unmarshal:%v", err)
		c.RecvError(topic)
		return
	}
	c.logger.Trace("【%v】- RequestID %v", texCon.GetUniqueKey(), outResInfo.RequestInfo.RequestID)

	ok := c.waitManager.SetWaitInfo(texCon, outResInfo.RequestInfo.ProcessID, outResInfo.RequestInfo.RequestID, outResInfo)
	if !ok {
		c.logger.Debug("【%v】待機情報無\n", texCon.GetUniqueKey())
		return
	}

	controlID := c.makeControlID(outResInfo.RequestInfo.ProcessID, outResInfo.RequestInfo.RequestID)
	if _, ok := c.recvControlMap[controlID]; ok {
		// 返信元チャネルのチャネルへ情報を返却
		c.recvControlMap[controlID].channel <- outResInfo
		// 不要になったMAPを削除
		c.releaseRecvControlMap(controlID)
		c.logger.Trace("【%v】END:paymentSendRecv SendRequestGetSalesInfo", texCon.GetUniqueKey())
		return
	}

}

func (c *paymentSendRecv) releaseRecvControlMap(id ControlID) {
	if _, ok := c.recvControlMap[id]; ok {
		c.recvControlMapMutex.Lock()
		defer c.recvControlMapMutex.Unlock()
		delete(c.recvControlMap, id) // map情報削除
	}

}

// 受信時エラー共通処理
func (c *paymentSendRecv) RecvError(resultTopic string) {
	// 同Topicを全て解除する
	for key, value := range c.recvControlMap {
		if value.resultTopic == resultTopic {
			value.channel <- errors.New(resultTopic + "連携時にエラーが発生") // 上位に失敗を通知
			c.recvControlMapMutex.Lock()
			delete(c.recvControlMap, key) // map情報削除
			c.recvControlMapMutex.Unlock()
		}
	}

}
