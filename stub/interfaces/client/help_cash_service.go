package interfaces

import (
	"stub/domain/handler"
)

type helpCashService struct {
	mqtt   handler.MqttRepository
	logger handler.LoggerRepository
}

func NewHelpCashService(mqtt handler.MqttRepository,
	logger handler.LoggerRepository) HelpCashService {
	return &helpCashService{
		mqtt:   mqtt,
		logger: logger}
}

var topicHelpCash = [19]string{
	"request_in_start",
	"request_in_end",
	"request_out_start",
	"request_out_stop",
	"request_collect_start",
	"request_collect_stop",
	"request_in_status",
	"request_out_status",
	"request_collect_status",
	"request_amount_status",
	"request_status",
	"request_set_amount",
	"request_warning_reset",
	"request_scrutiny_start",
	"notice_in_status",
	"notice_out_status",
	"notice_collect_status",
	"notice_amount_status",
	"notice_status"}

func (c *helpCashService) StartHelpCash() {
	c.mqtt.Subscribe("/tex/helper/cashctl/request_amount_status", c.RecvRequestAmountStatus)
}

func (c *helpCashService) RecvRequestAmountStatus(message string) {

}
