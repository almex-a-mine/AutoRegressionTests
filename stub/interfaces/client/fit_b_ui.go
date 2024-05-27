package interfaces

import (
	"encoding/json"
	"fmt"
	"stub/domain"
	"stub/domain/handler"
)

type communicateWithClient struct {
	mqtt   handler.MqttRepository
	logger handler.LoggerRepository
}

func NewFitB_Ui(mqtt handler.MqttRepository,
	logger handler.LoggerRepository) CommunicateWithClient {
	return &communicateWithClient{
		mqtt:   mqtt,
		logger: logger}
}

const (
	TOPIC_NUMBER = 25 //Topic数
	TOPIC_BASE   = "/tex/unifunc/money/"
)

var topicUI [TOPIC_NUMBER]string
var topicName = [TOPIC_NUMBER]string{
	"request_money_init",
	"request_money_exchange",
	"request_money_add_replenish",
	"request_money_collect",
	"request_set_amount",
	"request_status_cash",
	"request_pay_cash",
	"request_out_cash",
	"request_amount_cash",
	"request_print_report",
	"request_sales_info",
	"request_clear_cashinfo",
	"request_maintenance_mode",
	"request_reverseexchange_calculation",
	"request_coincassette_control",
	"request_get_safeinfo",
	"request_register_moneysetting",
	"request_get_moneysetting",
	"request_scrutiny",
	"notice_indata",
	"notice_outdata",
	"notice_collectdata",
	"notice_amount",
	"notice_status_cash",
	"notice_report_status"}

// 開始処理
func (c *communicateWithClient) Snalio1() {
	c.mqtt.Subscribe("/tex/unifunc/money/notice_status_cash", c.RecvNoticeStatusCash)
	/*var recvFunc = [TOPIC_NUMBER]func(string){
		//"request_money_init",
		//"request_money_exchange",
		//"request_money_add_replenish",
		//"request_money_collect",
		//"request_set_amount",
		//"request_status_cash",
		//"request_pay_cash",
		//"request_out_cash",
		//"request_amount_cash",
		//"request_print_report",
		//"request_sales_info",
		//"request_clear_cashinfo",
		//"request_maintenance_mode",
		//"request_reverseexchange_calculation",
		//"request_coincassette_control",
		//"request_get_safeinfo",
		//"request_register_moneysetting",
		//"request_get_moneysetting",
		//"request_scrutiny",
		//"notice_indata",
		//"notice_outdata",
		//"notice_collectdata",
		//"notice_amount",
		c.RecvNoticeStatusCash, //"notice_status_cash",
		//"notice_report_status"
	}
	for i := 0; i < TOPIC_NUMBER; i++ {
		topic[i] = fmt.Sprintf("%v/%v", TOPIC_BASE, topicName[i])
		c.mqtt.Subscribe(topic[i], recvFunc[i])
	}*/
}

// 通知のチェック
func (c *communicateWithClient) RecvNoticeStatusCash(message string) {
	var statusCash domain.StatusCash

	err := json.Unmarshal([]byte(message), &statusCash)
	if err != nil {
		c.logger.Error("json unmarshal error: %v", err)
		return
	}
	status := 1
	var testStatusCash = domain.StatusCash{}
	fmt.Println("testStatusCash: ", testStatusCash)
	switch status {
	case 1:
		/*テストデータ
		[INFO ] 2024/04/23 11:14:40.874798 [Send](MQTT)/tex/unifunc/money/notice_status_cash,
		  {"cashControlId":"","statusReady":true,"statusMode":1,"statusLine":true,"statusError":true,
		  "statusCover":false,"statusAction":0,"statusInsert":false,"statusExit":false,"statusRjbox":false,
		  "billStatusTbl":{"statusUnitSet":false,"statusInCassette":false,"statusOutCassette":false},
		  "coinStatusTbl":{"statusUnitSet":false,"statusInCassette":false,"statusOutCassette":false},
		  "billResidueInfoTbl":[],"coinResidueInfoTbl":[],"deviceStatusInfoTbl":["","","","","",""],
		  "warningInfoTbl":[0,0,0,0,0,0,0,0,0,0]}*/
		testStatusCash = domain.StatusCash{
			CashControlId: "",
			StatusReady:   true,
			StatusMode:    1,
			StatusLine:    true,
			StatusError:   true,
			StatusCover:   false,
			StatusAction:  0,
			StatusInsert:  false,
			StatusExit:    false,
			StatusRjbox:   false,
			BillStatusTbl: domain.TexmyBillStatusTbl{
				StatusUnitSet:     false,
				StatusInCassette:  false,
				StatusOutCassette: false,
				StatusAmountCount: 0,
			},
			CoinStatusTbl: domain.CoinStatusTbl{
				StatusUnitSet:     false,
				StatusInCassette:  false,
				StatusOutCassette: false,
				StatusAmountCount: 0,
			},
			BillResidueInfoTbl:  []domain.BillResidueInfo{},
			CoinResidueInfoTbl:  []domain.CoinResidueInfo{},
			DeviceStatusInfoTbl: []string{"", "", "", "", "", ""},
			WarningInfoTbl:      []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		}

	}

	// 比較
	/*if !isEqualStatusCash(testStatusCash, statusCash) {
		c.logger.Error("Senalio1_RecvNoticeStatus is not equal: %v", statusCash)
	}*/
}
