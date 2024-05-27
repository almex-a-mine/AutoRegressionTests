package calculation

import (
	"tex_money/domain"
)

// 使い方
// 返り値 := NewCassette([26]int配列).Add([26]int配列)

type MoneyCalculation struct {
	m10000 int // m:メイン
	m5000  int
	m2000  int
	m1000  int
	m500   int
	m100   int
	m50    int
	m10    int
	m5     int
	m1     int
	s500   int // s:サブ
	s100   int
	s50    int
	s10    int
	s5     int
	s1     int
	a10000 int // a:あふれ
	a5000  int
	a2000  int
	a1000  int
	a500   int
	a100   int
	a50    int
	a10    int
	a5     int
	a1     int
}

type MoneyCalculationRepository interface {
	Add(plus [26]int) [26]int                                  // 足し算メソッドと引数
	Exchange(amount int, changeType int) [26]int               // 両替 メソッドから払出金種の算出 changeType: 0=逆両替，1=全て1系金種で両替，2=1系,5系混在で両替 3=逆両替(紙幣のみ対象)
	GetOutCountTbl(amount int) (countTbl [16]int, result bool) // 払出 払出金種の算出(下位金種制限有)
	Subtract(minus [26]int) [26]int                            // 引き算 メソッド-引数
	GetTotalAmount() int                                       // 合計金額 メソッドで設定した値を金額合計として取得
	GetPayableCurrencyCount() int                              // あふれ金種枚数を除く、釣銭可能な金種の合計金額
	GetChangeAvailable() (int, [26]int)                        // 釣銭可能な金種の合計金額と配列を返却
	OverflowPriorityExchange(amount int) [26]int               // オーバーフローBOX(あふれ優先) 逆両替
	OverflowOnlyExchange(amount int) [26]int                   // オーバーフローのみを対象とした逆両替

	ExCountTblToTenCountTbl() [10]int       // 26金種配列を10金種配列に変更する
	ExCountTblToSixteenCountTbl() [16]int   // 26金種配列を16機種配列に変更する(あふれ分カット)
	AmountToTenCountTbl(amount int) [10]int // 枚数内訳算出 金額を10金種配列に変更する
}

func NewCassette(moneyTbl [26]int) MoneyCalculationRepository {
	return toCassette(moneyTbl)
}

// Exchange 両替要求
// changeType: 0=逆両替，1=全て1系金種で両替，2=1系,5系混在で両替 3=逆両替(紙幣のみ対象)
func (c *MoneyCalculation) Exchange(amount int, changeType int) [26]int {
	var result *MoneyCalculation
	switch changeType {
	case 0:
		result = c.toExchange(amount)
	// case 1:
	// 	result = c.oneExchange(amount)
	// case 2:
	// 	result = c.oneAndFiveExchange(amount)
	case 3:
		result = c.toExchangeBill(amount)
	}
	return c.toIntTbl26(result)
}

func (c *MoneyCalculation) ExCountTblToTenCountTbl() [10]int {
	return c.toIntTbl10(c)
}

func (c *MoneyCalculation) ExCountTblToSixteenCountTbl() [16]int {
	return c.toIntTbl16(c)
}

// AmountToTenCountTbl 枚数内訳算出
func (c *MoneyCalculation) AmountToTenCountTbl(amount int) [10]int {
	return c.toIntTbl10(c.amountToCountTbl(amount))
}

func (c *MoneyCalculation) OverflowPriorityExchange(amount int) [26]int {
	result := c.overflowPriorityExchange(amount)
	return c.toIntTbl26(result)
}

func (c *MoneyCalculation) OverflowOnlyExchange(amount int) [26]int {
	result := c.overflowOnlyExchange(amount)
	return c.toIntTbl26(result)
}

// Subtract 引き算要求
// メソッド-引数の値を返却する
func (c *MoneyCalculation) Subtract(minus [26]int) [26]int {
	minusCassette := toCassette(minus)
	return c.toIntTbl26(c.subtract(minusCassette))
}

func (c *MoneyCalculation) Add(plus [26]int) [26]int {
	plusMoney := toCassette(plus)
	return c.toIntTbl26(c.add(plusMoney))
}

func (c *MoneyCalculation) GetTotalAmount() int {
	var total int
	total += c.m10000 * 10000
	total += c.m5000 * 5000
	total += c.m2000 * 2000
	total += c.m1000 * 1000
	total += c.m500 * 500
	total += c.m100 * 100
	total += c.m50 * 50
	total += c.m10 * 10
	total += c.m5 * 5
	total += c.m1 * 1
	total += c.s500 * 500
	total += c.s100 * 100
	total += c.s50 * 50
	total += c.s10 * 10
	total += c.s5 * 5
	total += c.s1 * 1
	total += c.a10000 * 10000
	total += c.a5000 * 5000
	total += c.a2000 * 2000
	total += c.a1000 * 1000
	total += c.a500 * 500
	total += c.a100 * 100
	total += c.a50 * 50
	total += c.a10 * 10
	total += c.a5 * 5
	total += c.a1 * 1
	return total
}
func (c *MoneyCalculation) GetPayableCurrencyCount() int {
	var total int
	total += c.m10000 * 10000
	total += c.m5000 * 5000
	total += c.m2000 * 2000
	total += c.m1000 * 1000
	total += c.m500 * 500
	total += c.m100 * 100
	total += c.m50 * 50
	total += c.m10 * 10
	total += c.m5 * 5
	total += c.m1 * 1
	total += c.s500 * 500
	total += c.s100 * 100
	total += c.s50 * 50
	total += c.s10 * 10
	total += c.s5 * 5
	total += c.s1 * 1
	return total
}

func (c *MoneyCalculation) GetChangeAvailable() (int, [26]int) {
	var total int
	total += c.m10000 * 10000
	total += c.m5000 * 5000
	total += c.m2000 * 2000
	total += c.m1000 * 1000
	total += c.m500 * 500
	total += c.m100 * 100
	total += c.m50 * 50
	total += c.m10 * 10
	total += c.m5 * 5
	total += c.m1 * 1
	total += c.s500 * 500
	total += c.s100 * 100
	total += c.s50 * 50
	total += c.s10 * 10
	total += c.s5 * 5
	total += c.s1 * 1

	return total, [26]int{c.m10000, c.m5000, c.m2000, c.m1000, c.m500, c.m100, c.m50, c.m10, c.m5, c.m1, c.s500, c.s100, c.s50, c.s10, c.s5, c.s1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
}

// subtract メソッドから引数の各金種の枚数を引き算して返却する
func (c *MoneyCalculation) subtract(m *MoneyCalculation) *MoneyCalculation {
	return &MoneyCalculation{
		m10000: c.m10000 - m.m10000,
		m5000:  c.m5000 - m.m5000,
		m2000:  c.m2000 - m.m2000,
		m1000:  c.m1000 - m.m1000,
		m500:   c.m500 - m.m500,
		m100:   c.m100 - m.m100,
		m50:    c.m50 - m.m50,
		m10:    c.m10 - m.m10,
		m5:     c.m5 - m.m5,
		m1:     c.m1 - m.m1,
		s500:   c.s500 - m.s500,
		s100:   c.s100 - m.s100,
		s50:    c.s50 - m.s50,
		s10:    c.s10 - m.s10,
		s5:     c.s5 - m.s5,
		s1:     c.s1 - m.s1,
		a10000: c.a10000 - m.a10000,
		a5000:  c.a5000 - m.a5000,
		a2000:  c.a2000 - m.a2000,
		a1000:  c.a1000 - m.a1000,
		a500:   c.a500 - m.a500,
		a100:   c.a100 - m.a100,
		a50:    c.a50 - m.a50,
		a10:    c.a10 - m.a10,
		a5:     c.a5 - m.a5,
		a1:     c.a1 - m.a1,
	}
}

// subtract メソッドから引数の各金種の枚数を足し算して返却する
func (c *MoneyCalculation) add(m *MoneyCalculation) *MoneyCalculation {
	return &MoneyCalculation{
		m10000: c.m10000 + m.m10000,
		m5000:  c.m5000 + m.m5000,
		m2000:  c.m2000 + m.m2000,
		m1000:  c.m1000 + m.m1000,
		m500:   c.m500 + m.m500,
		m100:   c.m100 + m.m100,
		m50:    c.m50 + m.m50,
		m10:    c.m10 + m.m10,
		m5:     c.m5 + m.m5,
		m1:     c.m1 + m.m1,
		s500:   c.s500 + m.s500,
		s100:   c.s100 + m.s100,
		s50:    c.s50 + m.s50,
		s10:    c.s10 + m.s10,
		s5:     c.s5 + m.s5,
		s1:     c.s1 + m.s1,
		a10000: c.a10000 + m.a10000,
		a5000:  c.a5000 + m.a5000,
		a2000:  c.a2000 + m.a2000,
		a1000:  c.a1000 + m.a1000,
		a500:   c.a500 + m.a500,
		a100:   c.a100 + m.a100,
		a50:    c.a50 + m.a50,
		a10:    c.a10 + m.a10,
		a5:     c.a5 + m.a5,
		a1:     c.a1 + m.a1,
	}
}

// toCassette 配列を構造体へ変換。配列での操作による誤記対策として置き換えて扱う
func toCassette(before [26]int) *MoneyCalculation {
	return &MoneyCalculation{
		m10000: before[0],
		m5000:  before[1],
		m2000:  before[2],
		m1000:  before[3],
		m500:   before[4],
		m100:   before[5],
		m50:    before[6],
		m10:    before[7],
		m5:     before[8],
		m1:     before[9],
		s500:   before[10],
		s100:   before[11],
		s50:    before[12],
		s10:    before[13],
		s5:     before[14],
		s1:     before[15],
		a10000: before[16],
		a5000:  before[17],
		a2000:  before[18],
		a1000:  before[19],
		a500:   before[20],
		a100:   before[21],
		a50:    before[22],
		a10:    before[23],
		a5:     before[24],
		a1:     before[25],
	}
}

// toExchange 指定された金額に対してカセットから払い出せる上位金種の枚数(逆両替)を実施する
// 但し、2000円を最優先として対象とする
func (c *MoneyCalculation) toExchange(amount int) *MoneyCalculation {
	result := &MoneyCalculation{}
	var minus bool
	balanceAmount := amount
	if amount < 0 {
		balanceAmount = amount * -1
		minus = true
	}

	// 先に2000円札を払い出す
	result.m2000 = minimumInt(balanceAmount/2000, c.m2000)
	balanceAmount -= result.m2000 * 2000

	result.m10000 = minimumInt(balanceAmount/10000, c.m10000)
	balanceAmount -= result.m10000 * 10000

	result.m5000 = minimumInt(balanceAmount/5000, c.m5000)
	balanceAmount -= result.m5000 * 5000

	result.m1000 = minimumInt(balanceAmount/1000, c.m1000)
	balanceAmount -= result.m1000 * 1000

	// 硬貨は高い順からメイン→サブで計算する

	result.m500 = minimumInt(balanceAmount/500, c.m500)
	balanceAmount -= result.m500 * 500

	result.s500 = minimumInt(balanceAmount/500, c.s500)
	balanceAmount -= result.s500 * 500

	result.m100 = minimumInt(balanceAmount/100, c.m100)
	balanceAmount -= result.m100 * 100

	result.s100 = minimumInt(balanceAmount/100, c.s100)
	balanceAmount -= result.s100 * 100

	result.m50 = minimumInt(balanceAmount/50, c.m50)
	balanceAmount -= result.m50 * 50

	result.s50 = minimumInt(balanceAmount/50, c.s50)
	balanceAmount -= result.s50 * 50

	result.m10 = minimumInt(balanceAmount/10, c.m10)
	balanceAmount -= result.m10 * 10
	result.s10 = minimumInt(balanceAmount/10, c.s10)
	balanceAmount -= result.s10 * 10

	result.m5 = minimumInt(balanceAmount/5, c.m5)
	balanceAmount -= result.m5 * 5

	result.s5 = minimumInt(balanceAmount/5, c.s5)
	balanceAmount -= result.s5 * 5

	result.m1 = minimumInt(balanceAmount, c.m1)
	balanceAmount -= result.m1 * 1

	result.s1 = minimumInt(balanceAmount, c.s1)

	// 取得金額がマイナスだった場合、最後にマイナス調整を実施
	// あふれ部分は不要の為置き換え未実施
	if minus {
		result.m10000 *= -1
		result.m5000 *= -1
		result.m2000 *= -1
		result.m1000 *= -1
		result.m500 *= -1
		result.m100 *= -1
		result.m50 *= -1
		result.m10 *= -1
		result.m5 *= -1
		result.m1 *= -1
		result.s500 *= -1
		result.s100 *= -1
		result.s50 *= -1
		result.s10 *= -1
		result.s5 *= -1
		result.s1 *= -1
		return result
	}

	return result
}

// // oneExchange 指定された金額に対して1系金種で払い出せる下位金種の枚数(両替)を実施する
// func (c *MoneyCalculation) oneExchange(amount int) *MoneyCalculation {
// 	result := &MoneyCalculation{}
// 	var minus bool
// 	balanceAmount := amount
// 	if amount < 0 {
// 		balanceAmount = amount * -1
// 		minus = true
// 	}

// 	result.m10000 = minimumInt(balanceAmount/10000, c.m10000)
// 	balanceAmount -= result.m10000 * 10000

// 	result.m1000 = minimumInt(balanceAmount/1000, c.m1000)
// 	balanceAmount -= result.m1000 * 1000

// 	// 硬貨は高い順からメイン→サブで計算する
// 	result.m100 = minimumInt(balanceAmount/100, c.m100)
// 	balanceAmount -= result.m100 * 100

// 	result.s100 = minimumInt(balanceAmount/100, c.s100)
// 	balanceAmount -= result.s100 * 100

// 	result.m10 = minimumInt(balanceAmount/10, c.m10)
// 	balanceAmount -= result.m10 * 10
// 	result.s10 = minimumInt(balanceAmount/10, c.s10)
// 	balanceAmount -= result.s10 * 10

// 	result.m1 = minimumInt(balanceAmount, c.m1)
// 	balanceAmount -= result.m1 * 1

// 	result.s1 = minimumInt(balanceAmount, c.s1)

// 	// 取得金額がマイナスだった場合、最後にマイナス調整を実施
// 	// あふれ部分は不要の為置き換え未実施
// 	if minus {
// 		result.m10000 *= -1
// 		result.m1000 *= -1
// 		result.m100 *= -1
// 		result.m10 *= -1
// 		result.m1 *= -1
// 		result.s100 *= -1
// 		result.s10 *= -1
// 		result.s1 *= -1
// 		return result
// 	}

// 	return result
// }

// // oneAndFiveExchange 指定された金額に対して1系&5系金種で払い出せる下位金種の枚数(両替)を実施する
// func (c *MoneyCalculation) oneAndFiveExchange(amount int) *MoneyCalculation {
// 	result := &MoneyCalculation{}
// 	var minus bool
// 	balanceAmount := amount
// 	if amount < 0 {
// 		balanceAmount = amount * -1
// 		minus = true
// 	}

// 	result.m10000 = minimumInt(balanceAmount/10000, c.m10000)
// 	balanceAmount -= result.m10000 * 10000

// 	result.m5000 = minimumInt(balanceAmount/5000, c.m5000)
// 	balanceAmount -= result.m5000 * 5000

// 	result.m1000 = minimumInt(balanceAmount/1000, c.m1000)
// 	balanceAmount -= result.m1000 * 1000

// 	// 硬貨は高い順からメイン→サブで計算する

// 	result.m500 = minimumInt(balanceAmount/500, c.m500)
// 	balanceAmount -= result.m500 * 500

// 	result.s500 = minimumInt(balanceAmount/500, c.m500)
// 	balanceAmount -= result.s500 * 500

// 	result.m100 = minimumInt(balanceAmount/100, c.m100)
// 	balanceAmount -= result.m100 * 100

// 	result.s100 = minimumInt(balanceAmount/100, c.s100)
// 	balanceAmount -= result.s100 * 100

// 	result.m50 = minimumInt(balanceAmount/50, c.m50)
// 	balanceAmount -= result.m50 * 50

// 	result.s50 = minimumInt(balanceAmount/50, c.s50)
// 	balanceAmount -= result.s50 * 50

// 	result.m10 = minimumInt(balanceAmount/10, c.m10)
// 	balanceAmount -= result.m10 * 10
// 	result.s10 = minimumInt(balanceAmount/10, c.s10)
// 	balanceAmount -= result.s10 * 10

// 	result.m5 = minimumInt(balanceAmount/5, c.m5)
// 	balanceAmount -= result.m5 * 5

// 	result.s5 = minimumInt(balanceAmount/5, c.s5)
// 	balanceAmount -= result.s5 * 5

// 	result.m1 = minimumInt(balanceAmount, c.m1)
// 	balanceAmount -= result.m1 * 1

// 	result.s1 = minimumInt(balanceAmount, c.s1)

// 	// 取得金額がマイナスだった場合、最後にマイナス調整を実施
// 	// あふれ部分は不要の為置き換え未実施
// 	if minus {
// 		result.m10000 *= -1
// 		result.m5000 *= -1
// 		result.m1000 *= -1
// 		result.m500 *= -1
// 		result.m100 *= -1
// 		result.m50 *= -1
// 		result.m10 *= -1
// 		result.m5 *= -1
// 		result.m1 *= -1
// 		result.s500 *= -1
// 		result.s100 *= -1
// 		result.s50 *= -1
// 		result.s10 *= -1
// 		result.s5 *= -1
// 		result.s1 *= -1
// 		return result
// 	}

// 	return result
// }

// min は2つの整数のうち小さい方を返す
func minimumInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 払出 払出金種の算出(下位金種制限有)
func (c *MoneyCalculation) GetOutCountTbl(amount int) (countTbl [16]int, result bool) {
	var ngIdx int

	// No.2158対応：2千円札をスキップして払い出す
	temp2000 := c.m2000
	c.m2000 = 0
	countTbl, result, ngIdx = c.getOutCountTbl(amount)

	// 紙幣の払出を有高から払い出せない場合は2千円札を入れて再計算
	if !result && ngIdx < 3 {
		c.m2000 = temp2000
		countTbl, result, _ = c.getOutCountTbl(amount)
	}

	return
}

func (c *MoneyCalculation) getOutCountTbl(amount int) (countTbl [16]int, result bool, ngIdx int) {

	//払出可能有高セット
	amountData := c.toIntTbl26(c)
	for i := 0; i < 16; i++ {
		if i >= 10 {
			amountData[i-6] += amountData[i] //予備金庫
		}
	}

	// 払出金種の選別
	for i := 0; i < 10; i++ {
		//払い出したい金種別枚数を算出
		wantOutCount := amount / domain.Cash[i]

		//払出不要の金種は処理をスキップ
		if wantOutCount <= 0 {
			continue
		}

		//有高から払出可能な枚数を払出枚数にセット
		outCount := minimumInt(wantOutCount, amountData[i])
		countTbl[i] += outCount
		amount -= outCount * domain.Cash[i]

		//有高不足のときは下位2金種まで払出可能か判定する
		if outCount != wantOutCount {
			balance := (wantOutCount - amountData[i]) * domain.Cash[i]

			// 下位金種制限値を設定
			limit := 2
			if domain.Cash[i] == 10000 || domain.Cash[i] == 5000 {
				// No.2158対応：1万円と5千円札の場合は下位3金種まで払出可能とする
				limit = 3
			}

			for ii := 1; ii <= limit; ii++ {
				idx := i + ii                // 判定する下位金種のインデックス
				if idx >= len(domain.Cash) { // 下位の2金種に1円が含まれる場合は1円まで判定
					break
				}

				//不足金額に対しての下位金種での払出予定枚数を算出
				additionalOutPlanCount := balance / domain.Cash[idx]
				//有高から払出可能な枚数を払出枚数にセット
				additionalOutCount := minimumInt(additionalOutPlanCount, amountData[idx])
				countTbl[idx] += additionalOutCount
				amountData[idx] -= additionalOutCount
				amount -= additionalOutCount * domain.Cash[idx]
				balance -= additionalOutCount * domain.Cash[idx]

				if 0 == balance {
					break
				}
			}

			//2金種以内で払出不可の場合は払出をしない（100円の場合1円では払出をしない）
			if balance > 0 {
				return countTbl, false, i
			}
		}
	}

	return countTbl, true, 0
}

// toExchangeBill 指定された金額に対して紙幣から払い出せる上位金種の枚数(逆両替)を実施する
// 但し、2000円を最優先として対象とする。
func (c *MoneyCalculation) toExchangeBill(amount int) *MoneyCalculation {
	var minus bool
	// 取得金額がマイナスだった場合、最後にマイナス調整を実施
	// あふれ部分は不要の為置き換え未実施
	minusCheck := func(minus bool, result *MoneyCalculation) *MoneyCalculation {
		if minus {
			result.m10000 *= -1
			result.m5000 *= -1
			result.m2000 *= -1
			result.m1000 *= -1
			result.m500 *= -1
			result.m100 *= -1
			result.m50 *= -1
			result.m10 *= -1
			result.m5 *= -1
			result.m1 *= -1
			result.s500 *= -1
			result.s100 *= -1
			result.s50 *= -1
			result.s10 *= -1
			result.s5 *= -1
			result.s1 *= -1
			return result
		}
		return result
	}
	balanceAmount := amount
	if amount < 0 {
		balanceAmount = amount * -1
		minus = true
	}

	var resultAmount int
	resultMoneyCalculation := &MoneyCalculation{}

	// 2000円を可能な限り利用して、逆両替
	for i := 0; i <= c.m2000; i++ {
		resultAmount, resultMoneyCalculation = c.priorityBill2000(balanceAmount, i)
		if resultAmount == 0 {
			break
		}
	}

	return minusCheck(minus, resultMoneyCalculation)
}

// priorityBill2000 2000円、10000円、5000円の順序で計算
// 2000円が多く1000円が少ない場合に、端数が最終的に払出できない可能性がある為
// 払出できないあまり金額がある場合には、2000円を1枚減らして再計算させる
func (c *MoneyCalculation) priorityBill2000(balanceAmount int, balanceCount int) (int, *MoneyCalculation) {
	result := &MoneyCalculation{}
	// 先に2000円札を払い出す
	result.m2000 = minimumInt(balanceAmount/2000, c.m2000-balanceCount)
	balanceAmount -= result.m2000 * 2000

	result.m10000 = minimumInt(balanceAmount/10000, c.m10000)
	balanceAmount -= result.m10000 * 10000

	result.m5000 = minimumInt(balanceAmount/5000, c.m5000)
	balanceAmount -= result.m5000 * 5000

	result.m1000 = minimumInt(balanceAmount/1000, c.m1000)
	balanceAmount -= result.m1000 * 1000
	return balanceAmount, result
}

// overflowPriorityExchange 指定された金額に対してオーバーフローを優先してカウントした後、カセットから払い出せる上位金種の枚数(逆両替)を実施する
// 但し、あふれ金庫も2000円を最優先として対象とする
func (c *MoneyCalculation) overflowPriorityExchange(amount int) *MoneyCalculation {
	result := &MoneyCalculation{}
	var minus bool
	balanceAmount := amount
	if amount < 0 {
		balanceAmount = amount * -1
		minus = true
	}

	// あふれ金庫優先
	result.a2000 = minimumInt(balanceAmount/2000, c.a2000)
	balanceAmount -= result.a2000 * 2000

	result.a10000 = minimumInt(balanceAmount/10000, c.a10000)
	balanceAmount -= result.a10000 * 10000

	result.a5000 = minimumInt(balanceAmount/5000, c.a5000)
	balanceAmount -= result.a5000 * 5000

	result.a1000 = minimumInt(balanceAmount/1000, c.a1000)
	balanceAmount -= result.a1000 * 1000

	result.a500 = minimumInt(balanceAmount/500, c.a500)
	balanceAmount -= result.a500 * 500

	result.a100 = minimumInt(balanceAmount/100, c.a100)
	balanceAmount -= result.a100 * 100

	result.a50 = minimumInt(balanceAmount/50, c.a50)
	balanceAmount -= result.a50 * 50

	result.a10 = minimumInt(balanceAmount/10, c.a10)
	balanceAmount -= result.a10 * 10

	result.a5 = minimumInt(balanceAmount/5, c.a5)
	balanceAmount -= result.a5 * 5

	result.a1 = minimumInt(balanceAmount, c.a1)
	balanceAmount -= result.a1 * 1

	// 先に2000円札を払い出す
	result.m2000 = minimumInt(balanceAmount/2000, c.m2000)
	balanceAmount -= result.m2000 * 2000

	result.m10000 = minimumInt(balanceAmount/10000, c.m10000)
	balanceAmount -= result.m10000 * 10000

	result.m5000 = minimumInt(balanceAmount/5000, c.m5000)
	balanceAmount -= result.m5000 * 5000

	result.m1000 = minimumInt(balanceAmount/1000, c.m1000)
	balanceAmount -= result.m1000 * 1000

	// 硬貨は高い順からメイン→サブで計算する

	result.m500 = minimumInt(balanceAmount/500, c.m500)
	balanceAmount -= result.m500 * 500

	result.s500 = minimumInt(balanceAmount/500, c.m500)
	balanceAmount -= result.s500 * 500

	result.m100 = minimumInt(balanceAmount/100, c.m100)
	balanceAmount -= result.m100 * 100

	result.s100 = minimumInt(balanceAmount/100, c.s100)
	balanceAmount -= result.s100 * 100

	result.m50 = minimumInt(balanceAmount/50, c.m50)
	balanceAmount -= result.m50 * 50

	result.s50 = minimumInt(balanceAmount/50, c.s50)
	balanceAmount -= result.s50 * 50

	result.m10 = minimumInt(balanceAmount/10, c.m10)
	balanceAmount -= result.m10 * 10
	result.s10 = minimumInt(balanceAmount/10, c.s10)
	balanceAmount -= result.s10 * 10

	result.m5 = minimumInt(balanceAmount/5, c.m5)
	balanceAmount -= result.m5 * 5

	result.s5 = minimumInt(balanceAmount/5, c.s5)
	balanceAmount -= result.s5 * 5

	result.m1 = minimumInt(balanceAmount, c.m1)
	balanceAmount -= result.m1 * 1

	result.s1 = minimumInt(balanceAmount, c.s1)

	// 取得金額がマイナスだった場合、最後にマイナス調整を実施
	// あふれ部分は不要の為置き換え未実施
	if minus {
		result.m10000 *= -1
		result.m5000 *= -1
		result.m2000 *= -1
		result.m1000 *= -1
		result.m500 *= -1
		result.m100 *= -1
		result.m50 *= -1
		result.m10 *= -1
		result.m5 *= -1
		result.m1 *= -1
		result.s500 *= -1
		result.s100 *= -1
		result.s50 *= -1
		result.s10 *= -1
		result.s5 *= -1
		result.s1 *= -1
		result.a10000 *= -1
		result.a5000 *= -1
		result.a2000 *= -1
		result.a1000 *= -1
		result.a500 *= -1
		result.a100 *= -1
		result.a50 *= -1
		result.a10 *= -1
		result.a5 *= -1
		result.a1 *= -1
		return result
	}

	return result
}

// overflowOnlyExchange 指定された金額に対してオーバーフローのみで、逆両替を実施する
// 但し、あふれ金庫も2000円を最優先として対象とする
func (c *MoneyCalculation) overflowOnlyExchange(amount int) *MoneyCalculation {
	result := &MoneyCalculation{}
	var minus bool
	balanceAmount := amount
	if amount < 0 {
		balanceAmount = amount * -1
		minus = true
	}

	// あふれ金庫優先
	result.a2000 = minimumInt(balanceAmount/2000, c.a2000)
	balanceAmount -= result.a2000 * 2000

	result.a10000 = minimumInt(balanceAmount/10000, c.a10000)
	balanceAmount -= result.a10000 * 10000

	result.a5000 = minimumInt(balanceAmount/5000, c.a5000)
	balanceAmount -= result.a5000 * 5000

	result.a1000 = minimumInt(balanceAmount/1000, c.a1000)
	balanceAmount -= result.a1000 * 1000

	result.a500 = minimumInt(balanceAmount/500, c.a500)
	balanceAmount -= result.a500 * 500

	result.a100 = minimumInt(balanceAmount/100, c.a100)
	balanceAmount -= result.a100 * 100

	result.a50 = minimumInt(balanceAmount/50, c.a50)
	balanceAmount -= result.a50 * 50

	result.a10 = minimumInt(balanceAmount/10, c.a10)
	balanceAmount -= result.a10 * 10

	result.a5 = minimumInt(balanceAmount/5, c.a5)
	balanceAmount -= result.a5 * 5

	result.a1 = minimumInt(balanceAmount, c.a1)
	// balanceAmount -= result.a1 * 1

	// 取得金額がマイナスだった場合、最後にマイナス調整を実施
	// あふれ部分は不要の為置き換え未実施
	if minus {
		result.a10000 *= -1
		result.a5000 *= -1
		result.a2000 *= -1
		result.a1000 *= -1
		result.a500 *= -1
		result.a100 *= -1
		result.a50 *= -1
		result.a10 *= -1
		result.a5 *= -1
		result.a1 *= -1
		return result
	}

	return result
}

// amountToCountTbl 指定された金額に対して金種毎の枚数内訳を算出する（現在有高枚数は考慮しない）
// 但し、2000円は対象とする
func (c *MoneyCalculation) amountToCountTbl(amount int) *MoneyCalculation {
	result := &MoneyCalculation{}
	balanceAmount := amount

	// 上位金種を優先して枚数内訳を算出
	// 紙幣
	result.m10000 = balanceAmount / 10000
	balanceAmount -= result.m10000 * 10000

	result.m5000 = balanceAmount / 5000
	balanceAmount -= result.m5000 * 5000

	result.m1000 = balanceAmount / 1000
	balanceAmount -= result.m1000 * 1000

	// 硬貨

	result.m500 = balanceAmount / 500
	balanceAmount -= result.m500 * 500

	result.m100 = balanceAmount / 100
	balanceAmount -= result.m100 * 100

	result.m50 = balanceAmount / 50
	balanceAmount -= result.m50 * 50

	result.m10 = balanceAmount / 10
	balanceAmount -= result.m10 * 10

	result.m5 = balanceAmount / 5
	balanceAmount -= result.m5 * 5

	result.m1 = balanceAmount

	return result
}

// toIntTbl10
// From [0]:10,000円 [1]:5,000円 [2]:2,000円 [3]:1,000円[4]:500円 [5]:100円 [6]:50円 [7]:10円 [8]:5円 [9]:1円[10]:500円予備 [11]:100円予備 [12]:50円予備 [13]:10円予備 [14]:5円予備 [15]:1円予備
//
//	[16]:10,000円あふれ [17]:5,000円あふれ [18]:2,000円あふれ [19]:1,000円あふれ[20]:500円あふれ [21]100円あふれ [22]:50円あふれ[23]:10円あふれ [24]:5円あふれ [25]:1円あふれ
//
// To   [0]:10,000円 [1]:5,000円 [2]:2,000円 [3]:1,000円[4]:500円 [5]:100円 [6]:50円 [7]:10円 [8]:5円 [9]:1円
func (c *MoneyCalculation) toIntTbl10(i *MoneyCalculation) [10]int {
	var result [10]int
	result[0] = i.m10000 + i.a10000
	result[1] = i.m5000 + i.a5000
	result[2] = i.m2000 + i.a2000
	result[3] = i.m1000 + i.a1000
	result[4] = i.m500 + i.s500 + i.a500
	result[5] = i.m100 + i.s100 + i.a100
	result[6] = i.m50 + i.s50 + i.a50
	result[7] = i.m10 + i.s10 + i.a10
	result[8] = i.m5 + i.s5 + i.a5
	result[9] = i.m1 + i.s1 + i.a1
	return result
}

// toExtraCashTypeShitei Cassetteをresult向けの数値配列に変換
func (c *MoneyCalculation) toIntTbl26(i *MoneyCalculation) [26]int {
	var result [26]int
	result[0] = i.m10000
	result[1] = i.m5000
	result[2] = i.m2000
	result[3] = i.m1000
	result[4] = i.m500
	result[5] = i.m100
	result[6] = i.m50
	result[7] = i.m10
	result[8] = i.m5
	result[9] = i.m1
	result[10] = i.s500
	result[11] = i.s100
	result[12] = i.s50
	result[13] = i.s10
	result[14] = i.s5
	result[15] = i.s1
	result[16] = i.a10000
	result[17] = i.a5000
	result[18] = i.a2000
	result[19] = i.a1000
	result[20] = i.a500
	result[21] = i.a100
	result[22] = i.a50
	result[23] = i.a10
	result[24] = i.a5
	result[25] = i.a1
	return result
}

func (c *MoneyCalculation) toIntTbl16(i *MoneyCalculation) [16]int {
	var result [16]int
	result[0] = i.m10000
	result[1] = i.m5000
	result[2] = i.m2000
	result[3] = i.m1000
	result[4] = i.m500
	result[5] = i.m100
	result[6] = i.m50
	result[7] = i.m10
	result[8] = i.m5
	result[9] = i.m1
	result[10] = i.s500
	result[11] = i.s100
	result[12] = i.s50
	result[13] = i.s10
	result[14] = i.s5
	result[15] = i.s1
	return result
}

func (c *MoneyCalculation) CountTblToAmount(countTbl [10]int) int {
	var amount int
	for i, c := range countTbl {
		amount += c * domain.Cash[i]
	}
	return amount
}
