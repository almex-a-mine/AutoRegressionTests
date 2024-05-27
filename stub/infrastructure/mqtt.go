package infrastructure

import (
	"encoding/json"
	"fmt"
	"strings"
	"stub/domain"
	"stub/domain/handler"
	"stub/usecases"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type pubInfo struct {
	topic   string
	message string
}

type mqttHandler struct {
	Client        mqtt.Client
	MsgCH         chan mqtt.Message
	PubCH         chan pubInfo
	logger        handler.LoggerRepository
	topicTbl      map[string]func(string)
	topicTblMutex sync.Mutex
}

// MQTT NewHandler
func NewMQTTHandler(logger handler.LoggerRepository, tcp string, port int, clientId string) handler.MqttRepository {
	mqtthanlder := &mqttHandler{}

	mqtthanlder.topicTbl = make(map[string]func(string), 20)
	mqtthanlder.MsgCH = make(chan mqtt.Message)
	mqtthanlder.PubCH = make(chan pubInfo)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%v:%v", tcp, port))
	opts.SetClientID(fmt.Sprintf("%v", clientId))
	opts.SetAutoReconnect(true)
	opts.SetResumeSubs(true)
	opts.SetCleanSession(true)
	opts.SetConnectRetry(true)
	opts.SetMaxReconnectInterval(1 * time.Second)
	opts.SetConnectRetryInterval(1 * time.Second)
	opts.OnConnect = mqtthanlder.cnctHandler
	opts.OnConnectionLost = mqtthanlder.cnctLostHandler
	client := mqtt.NewClient(opts)
	mqtthanlder.Client = client
	mqtthanlder.logger = logger
	mqtthanlder.topicTblMutex = sync.Mutex{}

	go mqtthanlder.watchRecvMQTT()
	go mqtthanlder.watchSendMQTT()
	return mqtthanlder
}

// MQTT Regist Subscribe

// Subscribe method is used to subscribe to a specified MQTT topic. If the specified topic is not already subscribed, it subscribes to that topic and registers the provided callback
func (m *mqttHandler) Subscribe(subtopic string, fnc func(string)) {
	m.topicTblMutex.Lock()
	defer m.topicTblMutex.Unlock()

	// 指定トピック登録済み？
	_, ok := m.topicTbl[subtopic]
	if ok {
		return
	}

	// 未登録なら指定トピックを登録
	m.topicTbl[subtopic] = fnc

	// 新規登録で、既にブローカとの接続が完了している場合は、subscribeを行う
	if m.InConnectionOpen() {
		if subscribeToken := m.Client.Subscribe(subtopic, 0, m.msgHandler); subscribeToken.Wait() && subscribeToken.Error() != nil {
			m.logger.Error("(MQTT) Subscribe client.Subscribe NG,%v,%v,%v\n", subtopic, subscribeToken.Error(), subscribeToken.Error())
		} else {
			m.logger.Info("(MQTT) Subscribe Success:%v", subtopic)
		}
	}

}

// Unsubscribe method is used to unsubscribe from a specified MQTT topic. If the specified topic is already subscribed, it unsubscribes from that topic. It takes the following parameters
func (m *mqttHandler) Unsubscribe(subtopic string) {

	m.topicTblMutex.Lock()
	defer m.topicTblMutex.Unlock()

	// 指定トピック登録済み？
	_, ok := m.topicTbl[subtopic]
	if !ok {
		return
	}

	if unsubscribeToken := m.Client.Unsubscribe(subtopic); unsubscribeToken.Wait() && unsubscribeToken.Error() != nil {
		m.logger.Error("(MQTT) Unsubscribe Error:%v, %v, %v", subtopic, unsubscribeToken.Error(), unsubscribeToken.Error())
	} else {
		m.logger.Info("(MQTT) Unsubscribe Success:%v", subtopic)
		delete(m.topicTbl, subtopic)
	}

}

// MQTT Connection Check
func (m *mqttHandler) InConnectionOpen() bool {
	if m.Client == nil {
		return false
	} else {
		return m.Client.IsConnectionOpen()
	}
}

// MQTT Publish
func (m *mqttHandler) Publish(topic string, message string) bool {
	if m.Client == nil {
		return false
	} else {
		pubInfo := pubInfo{
			topic:   topic,
			message: message,
		}
		m.PubCH <- pubInfo
	}
	return true
}

// connection handler
func (m *mqttHandler) cnctHandler(client mqtt.Client) {
	gLogger.Debug("(MQTT EVENT) Broker Connected.client  IN")
	for topic := range m.topicTbl {

		if subscribeToken := client.Subscribe(topic, 0, m.msgHandler); subscribeToken.Wait() && subscribeToken.Error() != nil {
			gLogger.Error("connectHandler client.Subscribe CALL ERROR,%v,%v,%v\n", topic, subscribeToken.Error(), subscribeToken.Error())
		} else {
			gLogger.Debug("connectHandler client.Subscribe OK,%v", topic)
		}
	}
}

// Connection Lost Handler
func (m *mqttHandler) cnctLostHandler(client mqtt.Client, err error) {
	gLogger.Warn("(MQTT EVENT) Broker Connect lost  IN")
	for t := range m.topicTbl {
		client.Unsubscribe(t)
		gLogger.Debug("Unsubscribe:%s", t)
	}
}

func (m *mqttHandler) msgHandler(client mqtt.Client, msg mqtt.Message) {

	m.MsgCH <- msg

}

// MQTT Connect
func (m *mqttHandler) Mqttconnect() bool {
	m.logger.Debug("(MQTT) Mqttconnect IN")
	var connected bool

	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		m.logger.Error("(MQTT) connected token: %v\n", token.Error())
		connected = false
	} else {
		connected = true
	}

	m.logger.Debug("(MQTT) Mqttconnect OUT")
	return connected
}

// Recv Watch Thread
func (m *mqttHandler) watchRecvMQTT() {
	// 接続完了まで
	for {
		ok := m.Mqttconnect()
		if ok {
			break
		}
	}
	// 接続検出したら通知メッセージを待つ
	for {
		// メッセージ検出待ち
		msg := <-m.MsgCH

		topic := msg.Topic()
		payload := string(msg.Payload())

		// 通知されたTOPICに対応したメソッドが登録済みか？
		fn, ok := m.topicTbl[topic]
		if ok {

			// 期待していない受信の精査(resultのみを対象とする)
			if strings.Contains(topic, "result") {
				// リクエストインフォ情報を取得
				var reqInfo domain.RequestInfoOnly
				err := json.Unmarshal([]byte(payload), &reqInfo)
				// エラーではない場合
				if err == nil {
					// notice以外の場合
					if reqInfo.RequestInfo.ProcessID != "" || reqInfo.RequestInfo.PcId != "" || reqInfo.RequestInfo.RequestID != "" {
						// 待機情報有無チェック
						ok := usecases.GetWaitInfoMqttRecv(reqInfo.RequestInfo.ProcessID, reqInfo.RequestInfo.RequestID)
						// 待機情報が無い場合には、受信待に戻る
						if !ok {
							m.logger.Debug("待機情報無[%v-%v]", reqInfo.RequestInfo.ProcessID, reqInfo.RequestInfo.RequestID)
							continue
						}
					}
				}
			}

			// printのnoticeは出力対象ではないものが多い為、noticeの処理側で必要な場合にログ出力する
			if !strings.Contains(topic, "money/request_get_service") && !strings.Contains(topic, "print/notice_status") {
				m.logger.Info("[Recv](MQTT)%v,%v", topic, payload)
			}

			go fn(payload)

		}
	}
}

// Publish Watch Thread
func (m *mqttHandler) watchSendMQTT() {

	for {
		// メッセージ検出待ち
		pub := <-m.PubCH
		topic := pub.topic
		message := pub.message

		if m.Client != nil {
			m.Client.Publish(topic, 0, false, message)

			// サービス監視のログは出力しないように修正
			if strings.Contains(topic, "money/result_get_service") {
				continue
			}
			m.logger.Info("[Send](MQTT)%v,%v\n", topic, message)

		}

	}
}
