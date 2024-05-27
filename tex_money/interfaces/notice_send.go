package interfaces

import (
	"encoding/json"
	"sync"
	"tex_money/domain"
	"tex_money/domain/handler"
)

type (
	noticeSend struct {
		mqtt      handler.MqttRepository
		logger    handler.LoggerRepository
		sendMutex sync.Mutex
	}
	NoticeSendRepository interface {
		Send(texCon *domain.TexContext, noticeInfo interface{}, topic string) error
	}
)

func NewNoticeSend(mqtt handler.MqttRepository, logger handler.LoggerRepository) NoticeSendRepository {
	return &noticeSend{
		mqtt:   mqtt,
		logger: logger,
	}
}

// Send 通知送信（応答待機なし）
// noticeInfo:送信する通知のjson情報、top..ic:接頭語付きtopic名称
func (c *noticeSend) Send(texCon *domain.TexContext, noticeInfo interface{}, topic string) error {
	c.sendMutex.Lock() // 同時に同じ送信が発生した場合に調査が難しいので、1件単位で処理を行う
	defer c.sendMutex.Unlock()
	// JSON形式に変換
	data, err := json.Marshal(noticeInfo)
	if err != nil {
		c.logger.Error("【%v】noticeSend Send json.Marshal:%v", texCon.GetUniqueKey(), err)

		return err
	}
	c.mqtt.Publish(topic, string(data))

	return nil
}
