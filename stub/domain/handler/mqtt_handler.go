package handler

type MqttRepository interface {
	Subscribe(string, func(string))
	Unsubscribe(string)
	Publish(string, string) bool
	InConnectionOpen() bool
}
