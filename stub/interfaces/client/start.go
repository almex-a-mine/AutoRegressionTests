package interfaces

import (
	"encoding/json"
	"fmt"
	"os"
	"stub/config"
	"stub/domain"
	"stub/domain/handler"
	"stub/usecases"
)

// 起動シナリオのテスト
type start struct {
	mqtt           handler.MqttRepository
	logger         handler.LoggerRepository
	config         config.Configuration
	equalsTexMoney usecases.EqualsTexMoneyRepository
}

func NewStart(
	mqtt handler.MqttRepository,
	logger handler.LoggerRepository,
	config config.Configuration,
	equalsTexMoney usecases.EqualsTexMoneyRepository) Start {
	return &start{
		mqtt:           mqtt,
		logger:         logger,
		config:         config,
		equalsTexMoney: equalsTexMoney}
}

func (s *start) Senario1() {
	//tex/unifunc/money/notice_status_cashを受信する
	//s.mqtt.Subscribe("/tex/unifunc/money/notice_status_cash", s.RecvNoticeStatusCash)
	s.mqtt.Subscribe("/tex/helper/dbdata/request_get_terminfo_now", s.RecvGetTerminfoNow)
}

// 通知のチェック
func (s *start) RecvNoticeStatusCash(message string) {
	var statusCash domain.StatusCash

	err := json.Unmarshal([]byte(message), &statusCash)
	if err != nil {
		s.logger.Error("json unmarshal error: %v", err)
		return
	}
	status := 1
	var testStatusCash = domain.StatusCash{}
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

	// 比較 //TODO: isEqualStatusCash関数の実装
	if !s.equalsTexMoney.IsEqualStatusCash(testStatusCash, statusCash) {
		s.logger.Error("Senalio1_RecvNoticeStatus is not equal: %v", statusCash)
		os.Exit(1) // ステータスコード1でプログラムを終了
	}
}

// 稼働データ管理
func (s *start) RecvGetTerminfoNow(message string) {
	var req domain.RequestGetTermInfoNow

	err := json.Unmarshal([]byte(message), &req)
	if err != nil {
		s.logger.Error("json unmarshal error: %v", err)
		return
	}

	var test = domain.RequestGetTermInfoNow{}
	sinario := 1
	switch sinario {
	case 1:
		/*
			[INFO ] 2024/04/23 11:14:40.877889 [Send](MQTT)/tex/helper/dbdata/request_get_terminfo_now,
			{"requestInfo":{"processId":"00001358","pcId":"10.120.10.71","requestId":"TexMoney_1"}}
		*/
		test = domain.RequestGetTermInfoNow{
			RequestInfo: domain.RequestInfo{
				//ProcessID: config.Config.ReqInfo.ProcessID, //プロセスIDは実行されるたび値が変わるので、スコープに入れない
				PcId:      config.Config.ReqInfo.PcId,
				RequestID: "TexMoney_1",
			},
		}
		// 比較
		if ret := isEqual(test, req); !ret {
			s.logger.Error("RecvGetTerminfoNow is not equal: %t, %v", ret, sinario)
			os.Exit(1) // ステータスコード1でプログラムを終了
			return     //データ不整合がある場合は終了
		}

		s.SendResultGetTermInfoNow()
	}
}

func (s *start) SendResultGetTermInfoNow() {
	/*
		[INFO ] 2024/04/23 11:14:55.922897 [Recv](MQTT)/tex/helper/dbdata/result_get_terminfo_now,
		{"requestInfo":{"processId":"00001358","requestId":"TexMoney_1","pcId":"10.120.10.71"},
		"reportDate":20240423,"reportTime":111314414,
		"infoTerm":{"statusError":0,"termErrorCode":0,"termErrorState":0,"statusHandling":0,
		"statusSecurity":0,"statusDoor":0,"statusKeySw":2,"statusCall":0,"paymentMod":0},
		"infoTrade":{"statusTrade":false,"typeTrade":false,"billingAmount":0,"depositAmount":0,
		"paymentPlanAmount":0,"paymentAmount":0,"payoutBlance":0,"paymentType":0,
		"cashInfoTbl":[
		{"infoType":0,"amount":0,"countTbl":[0,0,0,0,0,0,0,0,0,0],"exCountTbl":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},
		{"infoType":1,"amount":0,"countTbl":[0,0,0,0,0,0,0,0,0,0],"exCountTbl":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},
		{"infoType":2,"amount":0,"countTbl":[0,0,0,0,0,0,0,0,0,0],"exCountTbl":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},
		{"infoType":3,"amount":0,"countTbl":[0,0,0,0,0,0,0,0,0,0],"exCountTbl":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},
		{"infoType":4,"amount":0,"countTbl":[0,0,0,0,0,0,0,0,0,0],"exCountTbl":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},
		{"infoType":5,"amount":0,"countTbl":[0,0,0,0,0,0,0,0,0,0],"exCountTbl":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]}]},
		"infoSales":{"salesAmount":6275900,"exchangeTotal":1000,"salesTypeTbl":[
		{"salesType":0,"paymentType":0,"amount":101700,"count":20},
		{"salesType":0,"paymentType":1,"amount":6174200,"count":9},
		{"salesType":0,"paymentType":2,"amount":0,"count":0},
		{"salesType":0,"paymentType":3,"amount":0,"count":0},
		{"salesType":0,"paymentType":4,"amount":0,"count":0},
		{"salesType":0,"paymentType":5,"amount":0,"count":0},
		{"salesType":1,"paymentType":0,"amount":0,"count":0},
		{"salesType":1,"paymentType":1,"amount":0,"count":0},
		{"salesType":1,"paymentType":2,"amount":0,"count":0},
		{"salesType":1,"paymentType":3,"amount":0,"count":0},
		{"salesType":1,"paymentType":4,"amount":0,"count":0},
		{"salesType":1,"paymentType":5,"amount":0,"count":0},
		{"salesType":2,"paymentType":0,"amount":0,"count":0},
		{"salesType":2,"paymentType":1,"amount":0,"count":0},
		{"salesType":2,"paymentType":2,"amount":0,"count":0},
		{"salesType":2,"paymentType":3,"amount":0,"count":0},
		{"salesType":2,"paymentType":4,"amount":0,"count":0},
		{"salesType":2,"paymentType":5,"amount":0,"count":0},
		{"salesType":3,"paymentType":0,"amount":0,"count":0},
		{"salesType":3,"paymentType":1,"amount":0,"count":0},
		{"salesType":3,"paymentType":2,"amount":0,"count":0},
		{"salesType":3,"paymentType":3,"amount":0,"count":0},
		{"salesType":3,"paymentType":4,"amount":0,"count":0},
		{"salesType":3,"paymentType":5,"amount":0,"count":0}]},
		"infoSafe":{"currentStatusTbl":[0,0,0,0,0,0,0,0,0],"sortInfotbl":[
			{"sortType":0,"amount":1000,"countTbl":[0,0,0,1,0,0,0,0,0,0],"exCountTbl":[0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},
			{"sortType":1,"amount":1000,"countTbl":[0,0,0,1,0,0,0,0,0,0],"exCountTbl":[0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},
			{"sortType":2,"amount":0,"countTbl":[0,0,0,0,0,0,0,0,0,0],"exCountTbl":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},
			{"sortType":3,"amount":135100,"countTbl":[11,0,0,18,1,66,0,0,0,0],"exCountTbl":[11,0,0,18,1,66,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},
			{"sortType":4,"amount":33400,"countTbl":[2,0,0,6,0,74,0,0,0,0],"exCountTbl":[2,0,0,6,0,74,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},
			{"sortType":5,"amount":101700,"countTbl":[9,0,0,12,1,-8,0,0,0,0],"exCountTbl":[9,0,0,12,1,-8,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},
			{"sortType":6,"amount":10600,"countTbl":[0,0,0,6,2,36,0,0,0,0],"exCountTbl":[0,0,0,6,2,36,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},
			{"sortType":7,"amount":111300,"countTbl":[9,0,0,17,3,28,0,0,0,0],"exCountTbl":[7,0,0,17,3,28,0,0,0,0,0,0,0,0,0,0,2,0,0,0,0,0,0,0,0,0]},
			{"sortType":8,"amount":-100700,"countTbl":[-9,0,0,-11,-1,8,0,0,0,0],"exCountTbl":[-7,0,0,-11,-1,8,0,0,0,0,0,0,0,0,0,0,-2,0,0,0,0,0,0,0,0,0]},{"sortType":9,"amount":0,"countTbl":[0,0,0,0,0,0,0,0,0,0],"exCountTbl":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},{"sortType":10,"amount":2031670,"countTbl":[100,100,100,249,110,170,160,170,160,170],"exCountTbl":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},{"sortType":90,"amount":100700,"countTbl":[9,0,0,11,1,-8,0,0,0,0],"exCountTbl":[7,0,0,11,1,-8,0,0,0,0,0,0,0,0,0,0,2,0,0,0,0,0,0,0,0,0]},{"sortType":91,"amount":1000,"countTbl":[0,0,0,1,0,0,0,0,0,0],"exCountTbl":[0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},{"sortType":92,"amount":0,"countTbl":[0,0,0,0,0,0,0,0,0,0],"exCountTbl":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]}]},"result":true}

	*/
	result := domain.ResultGetTermInfoNow{
		RequestInfo: domain.RequestInfo{
			ProcessID: "00001358",
			PcId:      "",
			RequestID: "TexMoney_1",
		},
		Result:      true,
		ErrorCode:   "",
		ErrorDetail: "",
		ReportDate:  20240423,
		ReportTime:  111440,
		InfoTermTbl: domain.InfoTermTbl{
			StatusError:    0,
			TermErrorCode:  0,
			TermErrorState: 0,
			StatusHandling: 0,
			StatusSecurity: 0,
			StatusDoor:     0,
			StatusKeySw:    0,
			StatusCall:     0,
			PaymentMode:    0,
		},
		InfoTradeTbl: domain.InfoTradeTbl{
			StatusTrade:       true,
			TypeTrade:         true,
			BillingAmount:     0,
			DepositAmount:     0,
			PaymentPlanAmount: 0,
			PaymentAmount:     0,
			PayoutBalance:     0,
			PaymentType:       0,
			CashInfoTbl: []domain.CashInfoTblGetTemNow{
				{
					InfoType:   0,
					Amount:     0,
					CountTbl:   [domain.CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					ExCountTbl: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				{
					InfoType:   1,
					Amount:     0,
					CountTbl:   [domain.CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					ExCountTbl: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				{
					InfoType:   2,
					Amount:     0,
					CountTbl:   [domain.CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					ExCountTbl: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},

				{
					InfoType:   3,
					Amount:     0,
					CountTbl:   [domain.CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					ExCountTbl: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				{
					InfoType:   4,
					Amount:     0,
					CountTbl:   [domain.CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					ExCountTbl: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				{

					InfoType:   5,
					Amount:     0,
					CountTbl:   [domain.CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					ExCountTbl: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
			},
		},
		InfoSalesTbl: domain.InfoSalesTbl{
			SalesAmount:   0,
			ExchangeTotal: 0,
			SalesTypeTbl: []domain.SalesTypeTbl{
				{
					SalesType:   0,
					PaymentType: 0,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   0,
					PaymentType: 1,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   0,
					PaymentType: 2,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   0,
					PaymentType: 3,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   0,
					PaymentType: 4,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   0,
					PaymentType: 5,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   1,
					PaymentType: 0,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   1,
					PaymentType: 1,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   1,
					PaymentType: 2,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   1,
					PaymentType: 3,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   1,
					PaymentType: 4,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   1,
					PaymentType: 5,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   2,
					PaymentType: 0,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   2,
					PaymentType: 1,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   2,
					PaymentType: 2,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   2,
					PaymentType: 3,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   2,
					PaymentType: 4,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   2,
					PaymentType: 5,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   3,
					PaymentType: 0,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   3,
					PaymentType: 1,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   3,
					PaymentType: 2,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   3,
					PaymentType: 3,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   3,
					PaymentType: 4,
					Amount:      0,
					Count:       0,
				},
				{
					SalesType:   3,
					PaymentType: 5,
					Amount:      0,
					Count:       0,
				},
			},
		},
		InfoSafeTblGetTermNow: domain.InfoSafeTblGetTermNow{
			CurrentStatusTbl: [domain.CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
			SortInfoTbl: []domain.SortInfoTbl{
				{
					SortType:   0,
					Amount:     0,
					CountTbl:   [domain.CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					ExCountTbl: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				{
					SortType:   1,
					Amount:     0,
					CountTbl:   [domain.CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					ExCountTbl: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				{
					SortType:   2,
					Amount:     0,
					CountTbl:   [domain.CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					ExCountTbl: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				{
					SortType:   3,
					Amount:     0,
					CountTbl:   [domain.CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					ExCountTbl: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				{
					SortType:   4,
					Amount:     0,
					CountTbl:   [domain.CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					ExCountTbl: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				{
					SortType:   5,
					Amount:     0,
					CountTbl:   [domain.CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					ExCountTbl: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				{
					SortType:   6,
					Amount:     0,
					CountTbl:   [domain.CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					ExCountTbl: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				{
					SortType:   7,
					Amount:     0,
					CountTbl:   [domain.CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					ExCountTbl: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				{
					SortType:   8,
					Amount:     0,
					CountTbl:   [domain.CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					ExCountTbl: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
			},
		},
	}

	res, _ := json.Marshal(result)
	topic := fmt.Sprintf("%v/%v", topic_base_db, "result_get_terminfo_now")
	s.mqtt.Publish(topic, string(res))
}
