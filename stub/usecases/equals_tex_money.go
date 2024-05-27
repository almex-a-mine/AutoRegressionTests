package usecases

import (
	"reflect"
	"stub/domain"
	"stub/domain/handler"
)

type equalsTexMoney struct {
	mqtt   handler.MqttRepository
	logger handler.LoggerRepository
}

func NewEqualsTexMoney(mqtt handler.MqttRepository,
	logger handler.LoggerRepository) EqualsTexMoneyRepository {
	return &equalsTexMoney{
		mqtt:   mqtt,
		logger: logger}
}

// isEqualStatusCash は、2つの StatusCash オブジェクトが等しいかどうかをチェックします。
func (e *equalsTexMoney) IsEqualStatusCash(a, b domain.StatusCash) bool {
	// スライス以外のフィールドを比較
	if a.CashControlId != b.CashControlId ||
		a.StatusReady != b.StatusReady ||
		a.StatusMode != b.StatusMode ||
		a.StatusLine != b.StatusLine ||
		a.StatusError != b.StatusError ||
		a.StatusCover != b.StatusCover ||
		a.StatusAction != b.StatusAction ||
		a.StatusInsert != b.StatusInsert ||
		a.StatusExit != b.StatusExit ||
		a.StatusRjbox != b.StatusRjbox ||
		!reflect.DeepEqual(a.BillStatusTbl, b.BillStatusTbl) ||
		!reflect.DeepEqual(a.CoinStatusTbl, b.CoinStatusTbl) {
		return false
	}

	// スライスフィールドを比較
	return reflect.DeepEqual(a.BillResidueInfoTbl, b.BillResidueInfoTbl) &&
		reflect.DeepEqual(a.CoinResidueInfoTbl, b.CoinResidueInfoTbl) &&
		reflect.DeepEqual(a.DeviceStatusInfoTbl, b.DeviceStatusInfoTbl) &&
		reflect.DeepEqual(a.WarningInfoTbl, b.WarningInfoTbl)
}
