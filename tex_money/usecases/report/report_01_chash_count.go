package report

import (
	"fmt"
	"tex_money/domain"
	"tex_money/domain/handler"
)

type cashCountReport struct {
	agg        domain.AggregateSafeInfo
	logger     handler.LoggerRepository
	numInfoTbl []int
}

type CashCountReportRepository interface {
	GetCashCountReport(texCon *domain.TexContext) (numInfoTbl []int) // キャッシュカウントレポート情報取得
}

func NewCashCountReport(agg domain.AggregateSafeInfo, logger handler.LoggerRepository) CashCountReportRepository {
	return &cashCountReport{
		agg:        agg,
		logger:     logger,
		numInfoTbl: make([]int, 0, 70),
	}
}

// キャッシュカウントレポート情報取得
// numInfoTbl
// [0]~[12]：締め前（紙幣/青カセット/ピンクカセット） [13]~[25]：入金（紙幣/青カセット/ピンクカセット）
// [26]~[38]：出金（紙幣/青カセット/ピンクカセット）[39]~[51]：締め後（紙幣/青カセット/ピンクカセット）
// [52]：精算機内合計 [53]~[58]：収納枚数（入金金庫）[59]~[64]：金額（入金金庫）
// [65]：入金金庫合計 [66]：出金額合計 [67]：回収計 [68]：釣銭準備金 [69]：売上金
func (c *cashCountReport) GetCashCountReport(texCon *domain.TexContext) []int {

	//締め前有高
	c.setAggTbl(c.agg.BeforeAmountCountTbl)
	//入金枚数
	c.setAggTbl(c.agg.ReplenishCountTbl)
	//出金
	c.setAggTbl(c.agg.CollectCountTbl)
	//締め後
	c.setAggTbl(c.agg.AfterAmountCountTbl)

	//精算機内合計
	totalAmount := getTotalAmount(c.agg.AfterAmountCountTbl)
	c.numInfoTbl = append(c.numInfoTbl, totalAmount)

	//収納枚数
	for i := 20; i < 26; i++ {
		c.numInfoTbl = append(c.numInfoTbl, c.agg.BeforeAmountCountTbl[i]) //処理前有高金種配列のあふれ枚数をセット
	}

	//金額（入金金庫 金種別出金金額）
	var totalPaymentSafe int
	for i, v := range c.numInfoTbl[53:59] {
		amount := v * domain.Safe[i]
		c.numInfoTbl = append(c.numInfoTbl, amount) // 入金金庫出金金額をセット
		totalPaymentSafe += amount                  //入金金庫合計を算出
	}

	//入金金庫合計
	c.numInfoTbl = append(c.numInfoTbl, totalPaymentSafe)

	//出金額計 = 紙幣/青カセット/ピンクカセットの出金合計 + 入金金庫出金合計
	totalOutAmount := getTotalAmount(c.agg.CollectCountTbl)
	totalOutAmount += totalPaymentSafe
	c.numInfoTbl = append(c.numInfoTbl, totalOutAmount)

	//回収計 出金合計金額と同値
	c.numInfoTbl = append(c.numInfoTbl, totalOutAmount)

	//釣銭準備金 締め後の有高合計金額
	c.numInfoTbl = append(c.numInfoTbl, totalAmount)

	//売上金　売上金回収枚数*金額
	var totalSales int
	for i, s := range c.agg.SalesCollectCountTbl {
		totalSales += s * domain.AllCashInMachineTwentySix[i]

	}
	c.numInfoTbl = append(c.numInfoTbl, totalSales)

	c.outputLogCashCountNumInfoTbl(texCon)
	return c.numInfoTbl
}

// 26金種配列から普通金庫と予備金庫の情報をnumInfoTblにセットする
func (c *cashCountReport) setAggTbl(aggTbl [26]int) {
	for i, v := range aggTbl {
		if i > 15 {
			break
		}
		//予備カセットは1系のみ
		if i == 10 || i == 12 || i == 14 {
			continue
		}
		c.numInfoTbl = append(c.numInfoTbl, v)
	}
}

// 26金種配列のから普通金庫と予備金庫の合計金額を算出する
func getTotalAmount(aggTbl [26]int) (totalAmount int) {
	for i, v := range aggTbl {
		if i > 15 {
			break
		}
		//予備カセットは1系のみ
		if i == 10 || i == 12 || i == 14 {
			continue
		}
		totalAmount += v * domain.AllCashInMachineTwentySix[i] //合計を算出
	}
	return
}

// キャッシュカウントレポートログ出力
func (c *cashCountReport) outputLogCashCountNumInfoTbl(texCon *domain.TexContext) {
	l := fmt.Sprintf("【%v】キャッシュカウントレポート作成データ ---", texCon.GetUniqueKey())
	l += fmt.Sprintf("  %v : %+v", "締め前", c.numInfoTbl[0:13])
	l += fmt.Sprintf("  %v : %+v", "入金", c.numInfoTbl[13:26])
	l += fmt.Sprintf("  %v : %+v", "出金", c.numInfoTbl[26:39])
	l += fmt.Sprintf("  %v : %+v", "締め後", c.numInfoTbl[39:52])
	l += fmt.Sprintf("  %v : %+v", "精算機内合計", c.numInfoTbl[52])
	l += fmt.Sprintf("  %v : %+v", "収納枚数", c.numInfoTbl[53:59])
	l += fmt.Sprintf("  %v : %+v", "金額", c.numInfoTbl[59:65])
	l += fmt.Sprintf("  %v : %+v", "入金金庫合計", c.numInfoTbl[65])
	l += fmt.Sprintf("  %v : %+v", "出金額合計", c.numInfoTbl[66])
	l += fmt.Sprintf("  %v : %+v", "回収計", c.numInfoTbl[67])
	l += fmt.Sprintf("  %v : %+v", "釣銭準備金", c.numInfoTbl[68])
	l += fmt.Sprintf("  %v : %+v", "売上金", c.numInfoTbl[69])
	c.logger.Debug("%v", l)
}
