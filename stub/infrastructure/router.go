package infrastructure

import (
	"stub/config"
	"stub/domain"
	"stub/domain/handler"
	interfaces "stub/interfaces/client"
	"stub/usecases"
	"sync"
	"time"
)

var wgrps sync.WaitGroup

func Router() {
	wgrps = sync.WaitGroup{}
	wgrps.Add(1)
	//設定取得
	config := config.Initialize(domain.SrvName)
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
		config.SystemConf.LogStopFatal)
	logger.Info("Program Start")
	logger.Info("config %+v", config)
	//MQTTハンドラ作成
	mqtt := NewMQTTHandler(logger, config.MqttConf.TCP, config.MqttConf.Port, domain.SrvName)

	// Start Contoller
	Controller(mqtt, logger, config)
}

// 停止処理
func RouterStop() {
	wgrps.Done()
}

func Controller(mqtt handler.MqttRepository, logger handler.LoggerRepository, config config.Configuration) {

	equalsTexMoney := usecases.NewEqualsTexMoney(mqtt, logger)

	//MQTTの接続を待機
	waitMqttConnected(mqtt)

	start := interfaces.NewStart(mqtt, logger, config, equalsTexMoney)
	start.Senario1()

	// 終了要求まで待機
	wgrps.Wait()

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
