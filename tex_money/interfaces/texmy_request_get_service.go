package interfaces

import (
	"encoding/json"
	"fmt"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases"
)

type getService struct {
	mqtt                  handler.MqttRepository
	logger                handler.LoggerRepository
	texMoneyNoticeManager usecases.TexMoneyNoticeManagerRepository
}

// 実行状態取得要求
func NewRequestGetService(mqtt handler.MqttRepository,
	logger handler.LoggerRepository,
	texMoneyNoticeManager usecases.TexMoneyNoticeManagerRepository) getServiceRepository {
	return &getService{
		mqtt:                  mqtt,
		logger:                logger,
		texMoneyNoticeManager: texMoneyNoticeManager,
	}
}

// 開始処理
func (g *getService) Start() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_get_service")
	g.mqtt.Subscribe(topic, g.recvResuestGetService)
}

// 停止処理
func (g *getService) Stop() {
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "request_get_service")
	g.mqtt.Unsubscribe(topic)
}

// サービス制御(制御なし)
func (g *getService) ControlService(reqInfo domain.RequestControlService) {
}

// サービス実行状態取得要求検出
func (g *getService) recvResuestGetService(message string) {
	texCon := domain.NewTexContext(domain.RegisterTexContext{
		ReceivingTopicName: "request_get_service",
	})

	var reqInfo domain.RequestGetService

	err := json.Unmarshal([]byte(message), &reqInfo)
	if err != nil {
		g.logger.Error("getService recvResuestGetService json.Unmarshal:%v", err)
		return
	}

	statusService := g.texMoneyNoticeManager.GetStatusServiceData(texCon)
	res := domain.ResultGetService{
		RequestInfo:   reqInfo.RequestInfo,
		Result:        true,
		StatusService: statusService.StatusService,
	}

	payment, err := json.Marshal(res)
	if err != nil {
		g.logger.Error("【%v】- json.Marshal:%v", texCon.GetUniqueKey(), err)
		return
	}
	topic := fmt.Sprintf("%v/%v", domain.TOPIC_TEXMONEY_BASE, "result_get_service")
	g.mqtt.Publish(topic, string(payment))

}
