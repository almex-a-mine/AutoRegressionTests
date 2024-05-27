package usecases

import (
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
)

// TODO 残作業
// 金庫情報への更新登録対応
// 読み込んだ釣銭準備金情報を利用した動き(今は仮置き)
// 更新発生時のDBへの登録要求

type CoinCassetteControlManager interface {
	Collection(texCon *domain.TexContext, cassetteType int) *domain.CoinCassette
	Exchange(texCon *domain.TexContext, cassetteType int) *domain.CoinCassette
	SpecificationReplenishment(texCon *domain.TexContext, cassetteType int, amountCount [domain.EXTRA_CASH_TYPE_SHITEI]int) *domain.CoinCassette
}

type coinCassetteControlManager struct {
	logger       handler.LoggerRepository
	safe         SafeInfoManager
	config       config.Configuration
	texmyHandler TexMoneyHandlerRepository
}

func NewCoinCassetteControlManager(logger handler.LoggerRepository,
	safe SafeInfoManager,
	config config.Configuration,
	texmyHandler TexMoneyHandlerRepository,
) CoinCassetteControlManager {
	return &coinCassetteControlManager{
		logger:       logger,
		safe:         safe,
		config:       config,
		texmyHandler: texmyHandler,
	}
}

// Collection 硬貨カセット回収操作
// 指定カセットを0にする
func (c *coinCassetteControlManager) Collection(texCon *domain.TexContext, cassetteType int) *domain.CoinCassette {
	// 操作前の有高枚数
	_, safeTwo := c.safe.GetSortInfo(texCon, 0)
	beforeSafeTwo := c.toCassette(safeTwo.ExCountTbl)

	// 操作後の釣銭可能枚
	afterSafeTwo := &domain.Cassette{}
	switch cassetteType {
	case 1: // メインカセット回収
		afterSafeTwo = c.ClearMainCassette(beforeSafeTwo) // メインカセットの有高を0に変更
	case 2: // サブカセット回収
		afterSafeTwo = c.ClearSubCassette(beforeSafeTwo) // サブカセットの有高を0に変更
	case 3: // メイン＆サブカセット回収
		afterSafeTwo = c.ClearSubCassette(c.ClearMainCassette(beforeSafeTwo)) //メイン＆サブのカセットの有高を0に変更
	}

	differenceSafeTwo := c.subtract(afterSafeTwo, beforeSafeTwo) // 差分枚数の抽出
	differenceTotalAmount := c.toTalAmount(differenceSafeTwo)    // 差分の合計金額

	return domain.NewCoinCassette(
		differenceTotalAmount,
		c.toExtraCashTypeShitei(differenceSafeTwo),
		c.toExtraCashTypeShitei(beforeSafeTwo),
		c.toExtraCashTypeShitei(afterSafeTwo),
		[domain.EXTRA_CASH_TYPE_SHITEI]int{}, // 回収の場合、出力無し
	)
}

// Exchange 硬貨カセットを釣銭準備金がセットされたカセットと交換する
// 指定されたカセットは、釣銭準備金の枚数がセットされる
func (c *coinCassetteControlManager) Exchange(texCon *domain.TexContext, cassetteType int) *domain.CoinCassette {
	// 操作前の有高枚数
	_, safeTwo := c.safe.GetSortInfo(texCon, 0)
	beforeSafeTwo := c.toCassette(safeTwo.ExCountTbl)

	// 釣銭準備金
	data := c.texmyHandler.GetMoneySetting() //金銭設定情報取得
	// int配列に変換
	setUp := c.toCassette(c.changeChangeReserveCountTo26IntTbl(data))

	// 操作後の釣銭可能枚
	afterSafeTwo := &domain.Cassette{}
	switch cassetteType {
	case 1:
		afterSafeTwo = c.changeMainCassette(beforeSafeTwo, setUp)
	case 2:
		afterSafeTwo = c.changeSubCassette(beforeSafeTwo, setUp)
	case 3:
		afterSafeTwo = c.changeSubCassette(c.changeMainCassette(beforeSafeTwo, setUp), setUp)
	}

	differenceSafeTwo := c.subtract(afterSafeTwo, beforeSafeTwo)  // 差分枚数の抽出
	differenceTotalAmount := c.toTalAmount(differenceSafeTwo)     // 差分の合計金額
	exchange := c.toExchange(differenceTotalAmount, afterSafeTwo) // 交換後のカセットから、差額分の合計金額に対して払出枚数を算出する。

	return domain.NewCoinCassette(
		differenceTotalAmount,
		c.toExtraCashTypeShitei(differenceSafeTwo),
		c.toExtraCashTypeShitei(beforeSafeTwo),
		c.toExtraCashTypeShitei(afterSafeTwo),
		c.toExtraCashTypeShitei(exchange),
	)

}

// SpecificationReplenishment 指定されたカセットに対して、指定枚数を追加する
// 指定されたカセットを正として足し算を行う
func (c *coinCassetteControlManager) SpecificationReplenishment(texCon *domain.TexContext, cassetteType int, inAmountCount [domain.EXTRA_CASH_TYPE_SHITEI]int) *domain.CoinCassette {
	// 操作前の有高枚数
	_, safeTwo := c.safe.GetSortInfo(texCon, 0)
	beforeSafeTwo := c.toCassette(safeTwo.ExCountTbl)

	// 入金された金種枚数
	inAmount := c.toCassette(inAmountCount)

	// 操作後の釣銭可能枚数
	afterSafeTwo := &domain.Cassette{}
	switch cassetteType {
	case 1:
		afterSafeTwo = c.addMainCassette(beforeSafeTwo, inAmount)
	case 2:
		afterSafeTwo = c.addSubCassette(beforeSafeTwo, inAmount)
	case 3:
		afterSafeTwo = c.addSubCassette(c.addMainCassette(beforeSafeTwo, inAmount), inAmount)
	}

	differenceSafeTwo := c.subtract(afterSafeTwo, beforeSafeTwo)  // 差分枚数の抽出
	differenceTotalAmount := c.toTalAmount(differenceSafeTwo)     // 差分の合計金額
	exchange := c.toExchange(differenceTotalAmount, afterSafeTwo) // 交換後のカセットから、差額分の合計金額に対して払出枚数を算出する。

	return domain.NewCoinCassette(
		differenceTotalAmount,
		c.toExtraCashTypeShitei(differenceSafeTwo),
		c.toExtraCashTypeShitei(beforeSafeTwo),
		c.toExtraCashTypeShitei(afterSafeTwo),
		c.toExtraCashTypeShitei(exchange),
	)

}

// toCassette 配列を構造体へ変換。配列での操作による誤記載対策として置き換えて扱う
func (c *coinCassetteControlManager) toCassette(before [domain.EXTRA_CASH_TYPE_SHITEI]int) *domain.Cassette {
	return &domain.Cassette{
		M10000: before[0],
		M5000:  before[1],
		M2000:  before[2],
		M1000:  before[3],
		M500:   before[4],
		M100:   before[5],
		M50:    before[6],
		M10:    before[7],
		M5:     before[8],
		M1:     before[9],
		S500:   before[10],
		S100:   before[11],
		S50:    before[12],
		S10:    before[13],
		S5:     before[14],
		S1:     before[15],
		A10000: before[16],
		A5000:  before[17],
		A2000:  before[18],
		A1000:  before[19],
		A500:   before[20],
		A100:   before[21],
		A50:    before[22],
		A10:    before[23],
		A5:     before[24],
		A1:     before[25],
	}
}

// ClearMainCassette メインカセットの数値を0にする
func (c *coinCassetteControlManager) ClearMainCassette(before *domain.Cassette) *domain.Cassette {
	return &domain.Cassette{
		M10000: before.M10000,
		M5000:  before.M5000,
		M2000:  before.M2000,
		M1000:  before.M1000,
		M500:   0,
		M100:   0,
		M50:    0,
		M10:    0,
		M5:     0,
		M1:     0,
		S500:   before.S500,
		S100:   before.S100,
		S50:    before.S50,
		S10:    before.S10,
		S5:     before.S5,
		S1:     before.S1,
		A10000: before.A10000,
		A5000:  before.A5000,
		A2000:  before.A2000,
		A1000:  before.A1000,
		A500:   before.A500,
		A100:   before.A100,
		A50:    before.A50,
		A10:    before.A10,
		A5:     before.A5,
		A1:     before.A1,
	}
}

// ClearSubCassette サブカセットの数値を0にする
func (c *coinCassetteControlManager) ClearSubCassette(before *domain.Cassette) *domain.Cassette {
	return &domain.Cassette{
		M10000: before.M10000,
		M5000:  before.M5000,
		M2000:  before.M2000,
		M1000:  before.M1000,
		M500:   before.M500,
		M100:   before.M100,
		M50:    before.M50,
		M10:    before.M10,
		M5:     before.M5,
		M1:     before.M1,
		S500:   0,
		S100:   0,
		S50:    0,
		S10:    0,
		S5:     0,
		S1:     0,
		A10000: before.A10000,
		A5000:  before.A5000,
		A2000:  before.A2000,
		A1000:  before.A1000,
		A500:   before.A500,
		A100:   before.A100,
		A50:    before.A50,
		A10:    before.A10,
		A5:     before.A5,
		A1:     before.A1,
	}
}

// addMainCassette 操作前に追加分を合算する(メインカセット)
// 要求カセットとセット枚数が異なる場合には、要求カセットを正とする為にメインとサブそれぞれを実装
func (c *coinCassetteControlManager) addMainCassette(before *domain.Cassette, addAmountCount *domain.Cassette) *domain.Cassette {
	return &domain.Cassette{
		M10000: before.M10000,
		M5000:  before.M5000,
		M2000:  before.M2000,
		M1000:  before.M1000,
		M500:   before.M500 + addAmountCount.M500,
		M100:   before.M100 + addAmountCount.M100,
		M50:    before.M50 + addAmountCount.M50,
		M10:    before.M10 + addAmountCount.M10,
		M5:     before.M5 + addAmountCount.M5,
		M1:     before.M1 + addAmountCount.M1,
		S500:   before.S500,
		S100:   before.S100,
		S50:    before.S50,
		S10:    before.S10,
		S5:     before.S5,
		S1:     before.S1,
		A10000: before.A10000,
		A5000:  before.A5000,
		A2000:  before.A2000,
		A1000:  before.A1000,
		A500:   before.A500,
		A100:   before.A100,
		A50:    before.A50,
		A10:    before.A10,
		A5:     before.A5,
		A1:     before.A1,
	}
}

func (c *coinCassetteControlManager) addSubCassette(before *domain.Cassette, addAmountCount *domain.Cassette) *domain.Cassette {
	return &domain.Cassette{
		M10000: before.M10000,
		M5000:  before.M5000,
		M2000:  before.M2000,
		M1000:  before.M1000,
		M500:   before.M500,
		M100:   before.M100,
		M50:    before.M50,
		M10:    before.M10,
		M5:     before.M5,
		M1:     before.M1,
		S500:   before.S500 + addAmountCount.S500,
		S100:   before.S100 + addAmountCount.S100,
		S50:    before.S50 + addAmountCount.S50,
		S10:    before.S10 + addAmountCount.S10,
		S5:     before.S5 + addAmountCount.S5,
		S1:     before.S1 + addAmountCount.S1,
		A10000: before.A10000,
		A5000:  before.A5000,
		A2000:  before.A2000,
		A1000:  before.A1000,
		A500:   before.A500,
		A100:   before.A100,
		A50:    before.A50,
		A10:    before.A10,
		A5:     before.A5,
		A1:     before.A1,
	}
}

// subtract 操作後から操作前の各金種の枚数を引き算して返却する
func (c *coinCassetteControlManager) subtract(after *domain.Cassette, before *domain.Cassette) *domain.Cassette {
	return &domain.Cassette{
		M10000: after.M10000 - before.M10000,
		M5000:  after.M5000 - before.M5000,
		M2000:  after.M2000 - before.M2000,
		M1000:  after.M1000 - before.M1000,
		M500:   after.M500 - before.M500,
		M100:   after.M100 - before.M100,
		M50:    after.M50 - before.M50,
		M10:    after.M10 - before.M10,
		M5:     after.M5 - before.M5,
		M1:     after.M1 - before.M1,
		S500:   after.S500 - before.S500,
		S100:   after.S100 - before.S100,
		S50:    after.S50 - before.S50,
		S10:    after.S10 - before.S10,
		S5:     after.S5 - before.S5,
		S1:     after.S1 - before.S1,
		A10000: after.A10000 - before.A10000,
		A5000:  after.A5000 - before.A5000,
		A2000:  after.A2000 - before.A2000,
		A1000:  after.A1000 - before.A1000,
		A500:   after.A500 - before.A500,
		A100:   after.A100 - before.A100,
		A50:    after.A50 - before.A50,
		A10:    after.A10 - before.A10,
		A5:     after.A5 - before.A5,
		A1:     after.A1 - before.A1,
	}
}

func (c *coinCassetteControlManager) changeMainCassette(before *domain.Cassette, setStartCount *domain.Cassette) *domain.Cassette {
	// configから取得するようにする。

	return &domain.Cassette{
		M10000: before.M10000,
		M5000:  before.M5000,
		M2000:  before.M2000,
		M1000:  before.M1000,
		M500:   setStartCount.M500,
		M100:   setStartCount.M100,
		M50:    setStartCount.M50,
		M10:    setStartCount.M10,
		M5:     setStartCount.M5,
		M1:     setStartCount.M1,
		S500:   before.S500,
		S100:   before.S100,
		S50:    before.S50,
		S10:    before.S10,
		S5:     before.S5,
		S1:     before.S1,
		A10000: before.A10000,
		A5000:  before.A5000,
		A2000:  before.A2000,
		A1000:  before.A1000,
		A500:   before.A500,
		A100:   before.A100,
		A50:    before.A50,
		A10:    before.A10,
		A5:     before.A5,
		A1:     before.A1,
	}
}

func (c *coinCassetteControlManager) changeSubCassette(before *domain.Cassette, setStartCount *domain.Cassette) *domain.Cassette {
	// configから取得するようにする。

	return &domain.Cassette{
		M10000: before.M10000,
		M5000:  before.M5000,
		M2000:  before.M2000,
		M1000:  before.M1000,
		M500:   before.M500,
		M100:   before.M100,
		M50:    before.M50,
		M10:    before.M10,
		M5:     before.M5,
		M1:     before.M1,
		S500:   setStartCount.S500,
		S100:   setStartCount.S100,
		S50:    setStartCount.S50,
		S10:    setStartCount.S10,
		S5:     setStartCount.S5,
		S1:     setStartCount.S1,
		A10000: before.A10000,
		A5000:  before.A5000,
		A2000:  before.A2000,
		A1000:  before.A1000,
		A500:   before.A500,
		A100:   before.A100,
		A50:    before.A50,
		A10:    before.A10,
		A5:     before.A5,
		A1:     before.A1,
	}

}

// toExchange 指定された金額に対してカセットから払い出せる上位金種の枚数(逆両替)を実施する
// 但し、2000円を最優先として対象とする
func (c *coinCassetteControlManager) toExchange(amount int, cassette *domain.Cassette) *domain.Cassette {
	result := &domain.Cassette{}
	var minus bool
	balanceAmount := amount
	if amount < 0 {
		balanceAmount = amount * -1
		minus = true
	}

	// 先に2000円札を払い出す
	result.M2000 = minimum(balanceAmount/2000, cassette.M2000)
	balanceAmount -= result.M2000 * 2000

	result.M10000 = minimum(balanceAmount/10000, cassette.M10000)
	balanceAmount -= result.M10000 * 10000

	result.M5000 = minimum(balanceAmount/5000, cassette.M5000)
	balanceAmount -= result.M5000 * 5000

	result.M1000 = minimum(balanceAmount/1000, cassette.M1000)
	balanceAmount -= result.M1000 * 1000

	// 硬貨は高い順からメイン→サブで計算する

	result.M500 = minimum(balanceAmount/500, cassette.M500)
	balanceAmount -= result.M500 * 500

	result.S500 = minimum(balanceAmount/500, cassette.S500)
	balanceAmount -= result.S500 * 500

	result.M100 = minimum(balanceAmount/100, cassette.M100)
	balanceAmount -= result.M100 * 100

	result.S100 = minimum(balanceAmount/100, cassette.S100)
	balanceAmount -= result.S100 * 100

	result.M50 = minimum(balanceAmount/50, cassette.M50)
	balanceAmount -= result.M50 * 50

	result.S50 = minimum(balanceAmount/50, cassette.S50)
	balanceAmount -= result.S50 * 50

	result.M10 = minimum(balanceAmount/10, cassette.M10)
	balanceAmount -= result.M10 * 10
	result.S10 = minimum(balanceAmount/10, cassette.S10)
	balanceAmount -= result.S10 * 10

	result.M5 = minimum(balanceAmount/5, cassette.M5)
	balanceAmount -= result.M5 * 5

	result.S5 = minimum(balanceAmount/5, cassette.S5)
	balanceAmount -= result.S5 * 5

	result.M1 = minimum(balanceAmount, cassette.M1)
	balanceAmount -= result.M1 * 1

	result.S1 = minimum(balanceAmount, cassette.S1)

	// 取得金額がマイナスだった場合、最後にマイナス調整を実施
	// あふれ部分は不要の為置き換え未実施
	if minus {
		result.M10000 *= -1
		result.M5000 *= -1
		result.M2000 *= -1
		result.M1000 *= -1
		result.M500 *= -1
		result.M100 *= -1
		result.M50 *= -1
		result.M10 *= -1
		result.M5 *= -1
		result.M1 *= -1
		result.S500 *= -1
		result.S100 *= -1
		result.S50 *= -1
		result.S10 *= -1
		result.S5 *= -1
		result.S1 *= -1
		return result
	}

	return result
}

// minimum は2つの整数のうち小さい方を返す
func minimum(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// toTalAmount 各金種の金額合計
func (c *coinCassetteControlManager) toTalAmount(i *domain.Cassette) int {
	var total int
	total += i.M10000 * 10000
	total += i.M5000 * 5000
	total += i.M2000 * 2000
	total += i.M1000 * 1000
	total += i.M500 * 500
	total += i.M100 * 100
	total += i.M50 * 50
	total += i.M10 * 10
	total += i.M5 * 5
	total += i.M1 * 1
	total += i.S500 * 500
	total += i.S100 * 100
	total += i.S50 * 50
	total += i.S10 * 10
	total += i.S5 * 5
	total += i.S1 * 1
	total += i.A10000 * 10000
	total += i.A5000 * 5000
	total += i.A2000 * 2000
	total += i.A1000 * 1000
	total += i.A500 * 500
	total += i.A100 * 100
	total += i.A50 * 50
	total += i.A10 * 10
	total += i.A5 * 5
	total += i.A1 * 1
	return total
}

// toExtraCashTypeShitei Cassetteをresult向けの数値配列に変換
func (c *coinCassetteControlManager) toExtraCashTypeShitei(i *domain.Cassette) [domain.EXTRA_CASH_TYPE_SHITEI]int {
	var result [domain.EXTRA_CASH_TYPE_SHITEI]int
	result[0] = i.M10000
	result[1] = i.M5000
	result[2] = i.M2000
	result[3] = i.M1000
	result[4] = i.M500
	result[5] = i.M100
	result[6] = i.M50
	result[7] = i.M10
	result[8] = i.M5
	result[9] = i.M1
	result[10] = i.S500
	result[11] = i.S100
	result[12] = i.S50
	result[13] = i.S10
	result[14] = i.S5
	result[15] = i.S1
	result[16] = i.A10000
	result[17] = i.A5000
	result[18] = i.A2000
	result[19] = i.A1000
	result[20] = i.A500
	result[21] = i.A100
	result[22] = i.A50
	result[23] = i.A10
	result[24] = i.A5
	result[25] = i.A1
	return result

}

func (c *coinCassetteControlManager) changeChangeReserveCountTo26IntTbl(d *domain.MoneySetting) [domain.EXTRA_CASH_TYPE_SHITEI]int {
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
