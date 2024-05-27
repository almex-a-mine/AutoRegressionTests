package usecases

import (
	"tex_money/config"
	"tex_money/domain"
)

type (
	currentStatus struct {
		moneySetting    *domain.MoneySetting
		moneyList       *domain.MoneyList
		overFlowBoxType *config.OverFlowBoxType
	}
	CurrentStatusRepository interface {
		MakeCurrentStatus(moneySetting *domain.MoneySetting, moneyList *domain.MoneyList) *[domain.CASH_TYPE_SHITEI]int
		ErrorCheckCurrentStatusTbl(currentStatusTbl [domain.CASH_TYPE_SHITEI]int) int
	}
)

func NewCurrentStatus(overFlowBoxType *config.OverFlowBoxType) CurrentStatusRepository {
	return &currentStatus{
		overFlowBoxType: overFlowBoxType,
	}
}

func (c *currentStatus) MakeCurrentStatus(moneySetting *domain.MoneySetting, moneyList *domain.MoneyList) *[domain.CASH_TYPE_SHITEI]int {
	return c.makeCurrentStatusTbl(moneySetting, moneyList)
}

func (c *currentStatus) makeCurrentStatusTbl(moneySetting *domain.MoneySetting, moneyList *domain.MoneyList) *[domain.CASH_TYPE_SHITEI]int {
	currentStatusTbl := &[domain.CASH_TYPE_SHITEI]int{}
	// あふれ処理優先の為、不足判定結果を元にあふれ判定は上書きする。
	/// 不足注意(1)
	c.CheckShortage(currentStatusTbl, moneySetting.ChangeShortageCount.RegisterDataTbl[0], moneyList, 1)
	/// あふれ注意(3)
	c.CheckExcess(currentStatusTbl, moneySetting.ExcessChangeCount.ExRegisterDataTbl[0], moneyList, 3)
	/// 不足エラー(2)
	c.CheckShortage(currentStatusTbl, moneySetting.ChangeShortageCount.RegisterDataTbl[1], moneyList, 2)
	/// あふれエラー(4)
	c.CheckExcess(currentStatusTbl, moneySetting.ExcessChangeCount.ExRegisterDataTbl[1], moneyList, 4)
	return currentStatusTbl
}

func (c *currentStatus) CheckShortage(currentTbl *[domain.CASH_TYPE_SHITEI]int, errorTbl domain.RegisterDataTbl, moneyList *domain.MoneyList, setStatus int) {
	type moneyAndShortage struct {
		requiredAmount int
		currentAmount  int
	}

	moneyMappings := []moneyAndShortage{
		{errorTbl.M10000Count, moneyList.M10000},
		{errorTbl.M5000Count, moneyList.M5000},
		{errorTbl.M2000Count, moneyList.M2000},
		{errorTbl.M1000Count, moneyList.M1000},
		{errorTbl.M500Count, moneyList.M500 + moneyList.S500},
		{errorTbl.M100Count, moneyList.M100 + moneyList.S100},
		{errorTbl.M50Count, moneyList.M50 + moneyList.S50},
		{errorTbl.M10Count, moneyList.M10 + moneyList.S10},
		{errorTbl.M5Count, moneyList.M5 + moneyList.S5},
		{errorTbl.M1Count, moneyList.M1 + moneyList.S1},
	}

	for i, mapping := range moneyMappings {
		if !c.shortage(mapping.requiredAmount, mapping.currentAmount) {
			currentTbl[i] = setStatus
		}
	}
}

func (c *currentStatus) CheckExcess(currentTbl *[domain.CASH_TYPE_SHITEI]int, errorTbl domain.ExRegisterDataTbl, moneyList *domain.MoneyList, setStatus int) {
	type moneyAndCount struct {
		money      int
		alertCount int
	}
	type billList struct {
		b10000 int
		b5000  int
		b2000  int
		b1000  int
	}
	type coinList struct {
		c500 int
		c100 int
		c50  int
		c10  int
		c5   int
		c1   int
	}

	bill := billList{
		b10000: moneyList.M10000,
		b5000:  moneyList.M5000,
		b2000:  moneyList.M2000,
		b1000:  moneyList.M1000,
	}
	// 回収庫有りの場合、あふれの金額をベースにする
	if c.overFlowBoxType.BillOverFlowBoxType {
		bill = billList{
			b10000: moneyList.A10000,
			b5000:  moneyList.A5000,
			b2000:  moneyList.A2000,
			b1000:  moneyList.A1000,
		}
	}

	// 基本、硬貨でメイン以外が還流するタイプは無い想定の為
	// メイン硬貨金庫の数値で比較
	coin := coinList{
		c500: moneyList.M500,
		c100: moneyList.M100,
		c50:  moneyList.M50,
		c10:  moneyList.M10,
		c5:   moneyList.M5,
		c1:   moneyList.M1,
	}
	// 回収庫有りの場合、あふれの金額をベースにする
	if c.overFlowBoxType.CoinOverFlowBoxType {
		coin = coinList{
			c500: moneyList.A500,
			c100: moneyList.A100,
			c50:  moneyList.A50,
			c10:  moneyList.A10,
			c5:   moneyList.A5,
			c1:   moneyList.A1,
		}
	}

	moneyMappings := []moneyAndCount{
		{bill.b10000, errorTbl.M10000Count},
		{bill.b5000, errorTbl.M5000Count},
		{bill.b2000, errorTbl.M2000Count},
		{bill.b1000, errorTbl.M1000Count},
		{coin.c500, errorTbl.M500Count},
		{coin.c100, errorTbl.M100Count},
		{coin.c50, errorTbl.M50Count},
		{coin.c10, errorTbl.M10Count},
		{coin.c5, errorTbl.M5Count},
		{coin.c1, errorTbl.M1Count},
	}

	for i, mapping := range moneyMappings {
		if !c.excess(mapping.alertCount, mapping.money) {
			currentTbl[i] = setStatus
		}
	}

	// 紙幣回収庫有りの場合、合計数比較を実施
	if c.overFlowBoxType.BillOverFlowBoxType {
		billSum := bill.b10000 + bill.b5000 + bill.b2000 + bill.b1000
		if !c.excess(errorTbl.BillOverBox, billSum) {
			for i := 0; i <= 3; i++ {
				currentTbl[i] = setStatus
			}
		}
	}

	// 硬貨回収庫有りの場合、合計数比較を実施
	if c.overFlowBoxType.CoinOverFlowBoxType {
		coinSum := coin.c500 + coin.c100 + coin.c50 + coin.c10 + coin.c5 + coin.c1
		if !c.excess(errorTbl.CoinOverBox, coinSum) {
			for i := 4; i <= 9; i++ {
				currentTbl[i] = setStatus
			}
		}
	}
}

func (c *currentStatus) ErrorCheckCurrentStatusTbl(currentStatusTbl [domain.CASH_TYPE_SHITEI]int) int {
	// 一致する特定のパターンの確認

	// 紙幣・硬貨がすべてオーバーエラーだった場合の処理
	// 0:正常 1:不足警告 2:不足エラー 3:オーバー警告 4:オーバーエラー
	if c.overFlowBoxType.BillOverFlowBoxType {
		if c.equals(currentStatusTbl[0:4], []int{4, 4, 4, 4}) {
			return ERROR_MANY_ALL_BILL
		}
	}
	if c.overFlowBoxType.CoinOverFlowBoxType {
		if c.equals(currentStatusTbl[4:], []int{4, 4, 4, 4, 4, 4}) {
			return ERROR_MANY_ALL_COIN
		}
	}

	// 紙幣、硬貨があふれワーニングの場合
	// 0:正常 1:不足警告 2:不足エラー 3:オーバー警告 4:オーバーエラー
	if c.overFlowBoxType.BillOverFlowBoxType {
		if c.equals(currentStatusTbl[0:4], []int{3, 3, 3, 3}) {
			return WARNING_MANY_ALL_BILL
		}
	}
	if c.overFlowBoxType.CoinOverFlowBoxType {
		if c.equals(currentStatusTbl[4:], []int{3, 3, 3, 3, 3, 3}) {
			return WARNING_MANY_ALL_COIN
		}
	}

	errorMappings := map[int][]int{
		1: {WARNING_NOTHING_TENTHOUSAND, WARNING_NOTHING_FIVETHOUSAND, WARNING_NOTHING_TWOTHOUSAND, WARNING_NOTHING_THOUSAND, WARNING_NOTHING_FIVEHUNDRED, WARNING_NOTHING_HUNDRED, WARNING_NOTHING_FIFTY, WARNING_NOTHING_TEN, WARNING_NOTHING_FIVE, WARNING_NOTHING_ONE},
		2: {ERROR_NOTHING_TENTHOUSAND, ERROR_NOTHING_FIVETHOUSAND, ERROR_NOTHING_TWOTHOUSAND, ERROR_NOTHING_THOUSAND, ERROR_NOTHING_FIVEHUNDRED, ERROR_NOTHING_HUNDRED, ERROR_NOTHING_FIFTY, ERROR_NOTHING_TEN, ERROR_NOTHING_FIVE, ERROR_NOTHING_ONE},
		3: {WARNING_MANY_TENTHOUSAND, WARNING_MANY_FIVETHOUSAND, WARNING_MANY_TWOTHOUSAND, WARNING_MANY_THOUSAND, WARNING_MANY_FIVEHUNDRED, WARNING_MANY_HUNDRED, WARNING_MANY_FIFTY, WARNING_MANY_TEN, WARNING_MANY_FIVE, WARNING_MANY_ONE},
		4: {ERROR_MANY_TENTHOUSAND, ERROR_MANY_FIVETHOUSAND, ERROR_MANY_TWOTHOUSAND, ERROR_MANY_THOUSAND, ERROR_MANY_FIVEHUNDRED, ERROR_MANY_HUNDRED, ERROR_MANY_FIFTY, ERROR_MANY_TEN, ERROR_MANY_FIVE, ERROR_MANY_ONE},
	}

	// あふれエラー、不足エラー、あふれワーニング、不足ワーニングの順に精査
	for _, priority := range []int{4, 2, 3, 1} {
		for i, v := range currentStatusTbl {
			if v == priority {
				return errorMappings[priority][i]
			}
		}
	}
	return 0 // Default no error code
}

func (c *currentStatus) equals(a, b []int) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// shortage alertCountと比較してmoney(有高)が下回るかどうかのチェック(不足チェック用)
func (c *currentStatus) shortage(alertCount int, money int) bool {
	if alertCount == 0 { // 0の場合チェックしない
		return true
	}
	if money < alertCount { // 有高が設定値以下であり、設定値が0ではない場合
		return false
	}
	return true
}

// excess alertCountと比較してmoney(有高)が上回るかどうかのチェック(あふれチェック用)
func (c *currentStatus) excess(alertCount int, money int) bool {
	if alertCount == 0 { // 0の場合チェックしない
		return true
	}
	if money > alertCount { // 有高が設定値以下であり、設定値が0ではない場合
		return false
	}
	return true
}
