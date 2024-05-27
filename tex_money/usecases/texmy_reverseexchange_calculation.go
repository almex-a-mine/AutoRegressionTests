package usecases

import (
	"errors"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/pkg/calculation"
)

type ReverseExchangeCalculationRepository interface {
	BaseExchange(texCon *domain.TexContext, exchangeType int) (int, [domain.EXTRA_CASH_TYPE_SHITEI]int, [domain.EXTRA_CASH_TYPE_SHITEI]int, error)     // 釣銭準備金に合わせた時の差分
	SpecifyExchange(texCon *domain.TexContext, amount int) (int, [domain.EXTRA_CASH_TYPE_SHITEI]int, error)                                            // 指定金額の逆両替算出
	SpecifyExchangeWithLowerDenominationLimit(texCon *domain.TexContext, amount int) (int, [domain.EXTRA_CASH_TYPE_SHITEI]int, error)                  // 指定金額(払出下位金種制限)の両替算出
	SalesMoneyExchange(texCon *domain.TexContext, exchangeType int, salesMoney int, overflowBox bool) (int, [domain.EXTRA_CASH_TYPE_SHITEI]int, error) // 売上金の逆両替
}

type reverseExchangeCalculationManager struct {
	logger       handler.LoggerRepository
	safe         SafeInfoManager
	config       config.Configuration
	texmyHandler TexMoneyHandlerRepository
}

func NewReverseExchangeCalculationManager(logger handler.LoggerRepository, config config.Configuration, safe SafeInfoManager, texmyHandler TexMoneyHandlerRepository) ReverseExchangeCalculationRepository {
	return &reverseExchangeCalculationManager{
		logger:       logger,
		safe:         safe,
		config:       config,
		texmyHandler: texmyHandler,
	}
}

// BaseExchange 釣銭準備金を基準として、現在のカセットとの差分を返却する
func (c *reverseExchangeCalculationManager) BaseExchange(texCon *domain.TexContext, exchangeType int) (int, [domain.EXTRA_CASH_TYPE_SHITEI]int, [domain.EXTRA_CASH_TYPE_SHITEI]int, error) {

	// 釣銭準備金
	data := c.texmyHandler.GetMoneySetting() //金銭設定情報取得
	// int配列に変換
	setUp := c.changeChangeReserveCountTo26IntTbl(data)

	// 現在の有高取得
	ok, safeZero := c.safe.GetSortInfo(texCon, 0)
	if !ok {
		return 0, [domain.EXTRA_CASH_TYPE_SHITEI]int{}, [domain.EXTRA_CASH_TYPE_SHITEI]int{}, errors.New("釣銭可能枚数の取得に失敗")
	}

	diffTotal, diff, reverseExchange := c.baseExchange(exchangeType, setUp, safeZero.ExCountTbl)

	return diffTotal, diff, reverseExchange, nil

}

// SpecifyExchange 指定金額に対する可能な逆両替値を算出する
func (c *reverseExchangeCalculationManager) SpecifyExchange(texCon *domain.TexContext, amount int) (int, [domain.EXTRA_CASH_TYPE_SHITEI]int, error) {
	// 現在の釣銭可能枚数を取得
	ok, safeZero := c.safe.GetSortInfo(texCon, 0)
	if !ok {
		return 0, [domain.EXTRA_CASH_TYPE_SHITEI]int{}, errors.New("釣銭可能枚数の取得に失敗")
	}
	exchange := c.specifyExchange(amount, safeZero.ExCountTbl)

	// 当初は要求金額そのまま返却していた。
	// ただし、要求が800円、払出可能が700円の場合に、Amountが800円、内訳が700円分という事象が発生する事が判明
	// amountを内訳の合計にする事で、上位で800円と700円の比較を行いエラーのハンドリングを実施する、又は
	// 800－700円で、100円足りない等を呼び出し元で判定できるように、
	// 要求金額に対して不足があった場合も、Resultはエラーにせずに
	// 可能な範囲の金額をのせて、情報の返却をする方向で対応する事で決定(中岡さん相談済) @2023/10/10 宮田
	exchangeAmount := calculation.NewCassette(exchange).GetTotalAmount()

	return exchangeAmount, exchange, nil

}

// SpecifyExchangeWithLowerDenominationLimit 指定金額に対する可能な両替値（払出下位金種制限有）を算出する
func (c *reverseExchangeCalculationManager) SpecifyExchangeWithLowerDenominationLimit(texCon *domain.TexContext, amount int) (int, [domain.EXTRA_CASH_TYPE_SHITEI]int, error) {
	// 現在の釣銭可能枚数を取得
	ok, safeZero := c.safe.GetSortInfo(texCon, 0)
	if !ok {
		return 0, [domain.EXTRA_CASH_TYPE_SHITEI]int{}, errors.New("釣銭可能枚数の取得に失敗")
	}
	exchange := c.specifyExchangeWithLowerDenominationLimit(amount, safeZero.ExCountTbl)
	exchangeAmount := calculation.NewCassette(exchange).GetTotalAmount()

	return exchangeAmount, exchange, nil

}

// SalesMoneyExchange 売上金と指定した金額、オーバーフローの考慮有無から、逆両替値を返却する
func (c *reverseExchangeCalculationManager) SalesMoneyExchange(texCon *domain.TexContext, exchangeType int, salesMoney int, overflowCashBox bool) (int, [domain.EXTRA_CASH_TYPE_SHITEI]int, error) {

	salesAmount := salesMoney
	if exchangeType == 2 {
		// 売上金回収済み金額を引く
		salesAmount -= c.safe.GetSafeInfo(texCon).SalesCompleteAmount
	}

	// 現在の釣銭可能枚数を取得
	ok, safeZero := c.safe.GetSortInfo(texCon, 0)
	if !ok {
		return 0, [domain.EXTRA_CASH_TYPE_SHITEI]int{}, errors.New("釣銭可能枚数の取得に失敗")
	}

	exchange := c.salesMoneyExchange(salesAmount, overflowCashBox, safeZero.ExCountTbl)
	exchangeTotalAmount := calculation.NewCassette(exchange).GetTotalAmount()

	return exchangeTotalAmount, exchange, nil

}

func (c *reverseExchangeCalculationManager) salesMoneyExchange(salesAmount int, overflowBox bool, safeZero [domain.EXTRA_CASH_TYPE_SHITEI]int) [domain.EXTRA_CASH_TYPE_SHITEI]int {

	if overflowBox {
		// オーバーフロー優先両替
		return calculation.NewCassette(safeZero).OverflowPriorityExchange(salesAmount)
	}
	// オーバーフローBOXを考慮しない
	// 通常の逆両替
	return calculation.NewCassette(safeZero).Exchange(salesAmount, 0)
}

func (c *reverseExchangeCalculationManager) specifyExchange(amount int, safeOne [domain.EXTRA_CASH_TYPE_SHITEI]int) [domain.EXTRA_CASH_TYPE_SHITEI]int {
	return calculation.NewCassette(safeOne).Exchange(amount, 0)

}

func (c *reverseExchangeCalculationManager) specifyExchangeWithLowerDenominationLimit(amount int, safeOne [domain.EXTRA_CASH_TYPE_SHITEI]int) [domain.EXTRA_CASH_TYPE_SHITEI]int {
	countTbl, result := calculation.NewCassette(safeOne).GetOutCountTbl(amount)
	var exchange [26]int
	copy(exchange[:16], countTbl[:])
	if !result {
		c.logger.Error("specifyExchangeWithLowerDenominationLimit 有高不足:2金種以内で払出不可")
		return exchange
	}
	return exchange
}

func (c *reverseExchangeCalculationManager) baseExchange(exchangeType int, setup [domain.EXTRA_CASH_TYPE_SHITEI]int, safeOne [domain.EXTRA_CASH_TYPE_SHITEI]int) (int, [domain.EXTRA_CASH_TYPE_SHITEI]int, [domain.EXTRA_CASH_TYPE_SHITEI]int) {
	switch exchangeType {
	case 10: // 紙幣
		diff := calculation.NewCassette(c.getBill(setup)).Subtract(c.getBill(safeOne)) // メイン紙幣釣銭準備金 - 現在のメイン紙幣有高枚数
		diffTotal := calculation.NewCassette(diff).GetTotalAmount()                    // 差分合計
		NewSafe := c.changeBill(safeOne, setup)                                        // 釣銭有高に釣銭準備金硬貨カセットセットした時の各金種の枚数
		reverseExchange := calculation.NewCassette(NewSafe).Exchange(diffTotal, 0)     // カセット交換後の各金種から逆両替を計算
		return diffTotal, diff, reverseExchange
	case 11: // 硬貨メイン
		diff := calculation.NewCassette(c.getMainCassette(setup)).Subtract(c.getMainCassette(safeOne)) // メイン硬貨カセット釣銭準備金 - 現在のメイン硬貨カセット有高枚数
		diffTotal := calculation.NewCassette(diff).GetTotalAmount()                                    // 差分合計
		NewSafe := c.changeMainCassette(safeOne, setup)                                                // 釣銭有高に釣銭準備金硬貨カセットセットした時の各金種の枚数
		reverseExchange := calculation.NewCassette(NewSafe).Exchange(diffTotal, 0)                     // カセット交換後の各金種から逆両替を計算
		return diffTotal, diff, reverseExchange
	case 12: // 硬貨サブ
		diff := calculation.NewCassette(c.getSubCassette(setup)).Subtract(c.getSubCassette(safeOne)) // サブ硬貨カセット釣銭準備金 - 現在のサブ硬貨カセット有高能枚数
		diffTotal := calculation.NewCassette(diff).GetTotalAmount()                                  // 差分合計
		NewSafe := c.changeSubCassette(safeOne, setup)                                               // 釣銭有高に釣銭準備金硬貨カセットセットした時の各金種の枚数
		reverseExchange := calculation.NewCassette(NewSafe).Exchange(diffTotal, 0)                   // カセット交換後の各金種から逆両替を計算
		return diffTotal, diff, reverseExchange
	case 13: // 硬貨メイン＆サブ
		diff := calculation.NewCassette(c.getCoinCassette(setup)).Subtract(c.getCoinCassette(safeOne)) // メイン＆サブ硬貨カセット釣銭準備金 - 現在のメイン＆サブ硬貨カセット有高枚数
		diffTotal := calculation.NewCassette(diff).GetTotalAmount()                                    // 差分合計
		NewSafe := c.changeCoinCassette(safeOne, setup)                                                // 釣銭有高に釣銭準備金硬貨カセットセットした時の各金種の枚数
		reverseExchange := calculation.NewCassette(NewSafe).Exchange(diffTotal, 0)                     // カセット交換後の各金種から逆両替を計算
		return diffTotal, diff, reverseExchange
	case 16: // 紙幣＆硬貨
		s := calculation.NewCassette(safeOne).ExCountTblToTenCountTbl()                                               // 現在有高を10金種に変換
		tempSafeOne := [domain.EXTRA_CASH_TYPE_SHITEI]int{s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7], s[8], s[9]} // 26金種領域に対して、10金種で現在有高情報をセットする
		diff := calculation.NewCassette(setup).Subtract(tempSafeOne)                                                  // 紙幣＆硬貨釣銭準備金 - 現在のメイン紙幣＆硬貨カセット有高枚数
		diffTotal := calculation.NewCassette(diff).GetTotalAmount()                                                   // 差分合計
		reverseExchange := calculation.NewCassette(tempSafeOne).Exchange(diffTotal, 0)                                // 現在有高から逆両替を計算
		return diffTotal, diff, reverseExchange
	default:

	}
	return 0, [domain.EXTRA_CASH_TYPE_SHITEI]int{}, [domain.EXTRA_CASH_TYPE_SHITEI]int{}
}

// getBill 紙幣の枚数のみを抽出
func (c *reverseExchangeCalculationManager) getBill(e [domain.EXTRA_CASH_TYPE_SHITEI]int) [domain.EXTRA_CASH_TYPE_SHITEI]int {
	return [domain.EXTRA_CASH_TYPE_SHITEI]int{e[0], e[1], e[2], e[3], 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
}

// getMainCassette 硬貨メインのみを抽出
func (c *reverseExchangeCalculationManager) getMainCassette(e [domain.EXTRA_CASH_TYPE_SHITEI]int) [domain.EXTRA_CASH_TYPE_SHITEI]int {
	return [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, e[4], e[5], e[6], e[7], e[8], e[9], 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
}

// getSubCassette 硬貨サブのみを抽出
func (c *reverseExchangeCalculationManager) getSubCassette(e [domain.EXTRA_CASH_TYPE_SHITEI]int) [domain.EXTRA_CASH_TYPE_SHITEI]int {
	return [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, e[10], e[11], e[12], e[13], e[14], e[15], 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
}

// getCoinCassette 硬貨のメイン＆サブを抽出
func (c *reverseExchangeCalculationManager) getCoinCassette(e [domain.EXTRA_CASH_TYPE_SHITEI]int) [domain.EXTRA_CASH_TYPE_SHITEI]int {
	return [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, e[4], e[5], e[6], e[7], e[8], e[9], e[10], e[11], e[12], e[13], e[14], e[15], 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
}

func (c *reverseExchangeCalculationManager) changeBill(e [domain.EXTRA_CASH_TYPE_SHITEI]int, set [domain.EXTRA_CASH_TYPE_SHITEI]int) [domain.EXTRA_CASH_TYPE_SHITEI]int {
	e[0] = set[0]
	e[1] = set[1]
	e[2] = set[2]
	e[3] = set[3]
	return e
}

func (c *reverseExchangeCalculationManager) changeMainCassette(e [domain.EXTRA_CASH_TYPE_SHITEI]int, set [domain.EXTRA_CASH_TYPE_SHITEI]int) [domain.EXTRA_CASH_TYPE_SHITEI]int {
	e[4] = set[4]
	e[5] = set[5]
	e[6] = set[6]
	e[7] = set[7]
	e[8] = set[8]
	e[9] = set[9]

	return e
}

func (c *reverseExchangeCalculationManager) changeSubCassette(e [domain.EXTRA_CASH_TYPE_SHITEI]int, set [domain.EXTRA_CASH_TYPE_SHITEI]int) [domain.EXTRA_CASH_TYPE_SHITEI]int {
	e[10] = set[10]
	e[11] = set[11]
	e[12] = set[12]
	e[13] = set[13]
	e[14] = set[14]
	e[15] = set[15]
	return e
}

func (c *reverseExchangeCalculationManager) changeCoinCassette(e [domain.EXTRA_CASH_TYPE_SHITEI]int, set [domain.EXTRA_CASH_TYPE_SHITEI]int) [domain.EXTRA_CASH_TYPE_SHITEI]int {
	e[4] = set[4]
	e[5] = set[5]
	e[6] = set[6]
	e[7] = set[7]
	e[8] = set[8]
	e[9] = set[9]
	e[10] = set[10]
	e[11] = set[11]
	e[12] = set[12]
	e[13] = set[13]
	e[14] = set[14]
	e[15] = set[15]
	return e
}

func (c *reverseExchangeCalculationManager) changeChangeReserveCountTo26IntTbl(d *domain.MoneySetting) [domain.EXTRA_CASH_TYPE_SHITEI]int {
	var r [domain.EXTRA_CASH_TYPE_SHITEI]int

	r[0] = d.ChangeReserveCount.M10000Count
	r[1] = d.ChangeReserveCount.M5000Count
	r[2] = d.ChangeReserveCount.M2000Count
	r[3] = d.ChangeReserveCount.M1000Count
	r[4] = d.ChangeReserveCount.M500Count
	r[5] = d.ChangeReserveCount.M100Count
	r[6] = d.ChangeReserveCount.M50Count
	r[7] = d.ChangeReserveCount.M10Count
	r[8] = d.ChangeReserveCount.M5Count
	r[9] = d.ChangeReserveCount.M1Count
	r[10] = d.ChangeReserveCount.S500Count
	r[11] = d.ChangeReserveCount.S100Count
	r[12] = d.ChangeReserveCount.S50Count
	r[13] = d.ChangeReserveCount.S10Count
	r[14] = d.ChangeReserveCount.S5Count
	r[15] = d.ChangeReserveCount.S1Count

	return r

}
