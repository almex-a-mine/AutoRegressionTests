package report

import (
	"fmt"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/pkg/calculation"
)

type summaryReport struct {
	agg        domain.AggregateSafeInfo
	safeInfo   domain.SafeInfo
	logger     handler.LoggerRepository
	numInfoTbl []int
}

type SummaryReportRepository interface {
	GetReportSummary(texCon *domain.TexContext, amount []int, count []int, totalAmount int, totalCount int) []int // 日計レシート（FIT-B NEXTクリニック向け）
}

func NewSummaryReport(agg domain.AggregateSafeInfo, safeInfo domain.SafeInfo, logger handler.LoggerRepository) SummaryReportRepository {
	return &summaryReport{
		agg:        agg,
		safeInfo:   safeInfo,
		logger:     logger,
		numInfoTbl: make([]int, 0, 28),
	}
}

// 日計レシート（FIT-B NEXTクリニック向け）
func (c *summaryReport) GetReportSummary(texCon *domain.TexContext, amount []int, count []int, totalAmount int, totalCount int) []int {
	// 売上金欄 =============================================================================
	// 現金　クレジット　電子マネー　QRコード　Jデビット+その他
	for i := range amount {
		//売上金額 //売上回数
		c.numInfoTbl = append(c.numInfoTbl, amount[i], count[i])
	}
	//売上金額合計//売上回数合計
	c.numInfoTbl = append(c.numInfoTbl, totalAmount, totalCount)

	// 入出金欄 =============================================================================
	var totalInitial, totalIn, totalOut, totalSupply, totalBeforeCollect int
	//準備金
	for i, v := range c.safeInfo.SortInfoTbl[2].CountTbl {
		c.numInfoTbl = append(c.numInfoTbl, v)
		amount := v * domain.Cash[i]
		totalInitial += amount //合計金額を算出
	}
	// 釣銭準備金（合計）
	c.numInfoTbl = append(c.numInfoTbl, totalInitial)

	// 入金
	for i, v := range c.safeInfo.SortInfoTbl[3].CountTbl {
		c.numInfoTbl = append(c.numInfoTbl, v)
		amount := v * domain.Cash[i]
		totalIn += amount //合計金額を算出
	}
	// 取引入金（合計）
	c.numInfoTbl = append(c.numInfoTbl, totalIn)

	// 払出
	for i, v := range c.safeInfo.SortInfoTbl[4].CountTbl {
		c.numInfoTbl = append(c.numInfoTbl, v)
		amount := v * domain.Cash[i]
		totalOut += amount //合計金額を算出
	}
	// 取引払出（合計）
	c.numInfoTbl = append(c.numInfoTbl, totalOut)

	// 補充
	// 10金種配列に変換
	beforeReplenishCountTbl := calculation.NewCassette(c.agg.BeforeReplenishCountTbl).ExCountTblToTenCountTbl()
	for i, v := range beforeReplenishCountTbl {
		c.numInfoTbl = append(c.numInfoTbl, v)
		amount := v * domain.Cash[i]
		totalSupply += amount //合計金額を算出
	}
	// 補充（合計）
	c.numInfoTbl = append(c.numInfoTbl, totalSupply)

	// 回収
	// 10金種配列に変換
	beforeCollectCountTbl := calculation.NewCassette(c.agg.BeforeCollectCountTbl).ExCountTblToTenCountTbl()
	for i, v := range beforeCollectCountTbl {
		c.numInfoTbl = append(c.numInfoTbl, v)
		amount := v * domain.Cash[i]
		totalBeforeCollect += amount //合計金額を算出
	}
	// 回収（合計）
	c.numInfoTbl = append(c.numInfoTbl, totalBeforeCollect)

	// 締め処理時欄 =============================================================================
	var totalBefore, totalClosingOut, totalAfter int
	// 処理前有高
	// 10金種配列に変換
	beforeCountTbl := calculation.NewCassette(c.agg.BeforeAmountCountTbl).ExCountTblToTenCountTbl()

	for i, v := range beforeCountTbl {
		c.numInfoTbl = append(c.numInfoTbl, v)
		amount := v * domain.Cash[i]
		totalBefore += amount //合計金額を算出
	}
	// 処理前有高（合計）
	c.numInfoTbl = append(c.numInfoTbl, totalBefore)

	// 出金（締め処理中の出金情報:回収金種配列）
	// 10金種配列に変換
	outCountTbl := calculation.NewCassette(c.agg.CollectCountTbl).ExCountTblToTenCountTbl()

	for i, v := range outCountTbl {
		c.numInfoTbl = append(c.numInfoTbl, v)
		amount := v * domain.Cash[i]
		totalClosingOut += amount //合計金額を算出
	}
	// 出金（合計）
	c.numInfoTbl = append(c.numInfoTbl, totalClosingOut)

	// 締め後有高
	// 10金種配列に変換
	afterCountTbl := calculation.NewCassette(c.agg.AfterAmountCountTbl).ExCountTblToTenCountTbl()

	for i, v := range afterCountTbl {
		c.numInfoTbl = append(c.numInfoTbl, v)
		amount := v * domain.Cash[i]
		totalAfter += amount //合計金額を算出
	}
	// 締め後有高（合計）
	c.numInfoTbl = append(c.numInfoTbl, totalAfter)

	// 総回収（回収+出金）
	totalCollectExCountTbl := calculation.NewCassette(c.agg.BeforeCollectCountTbl).Add(c.agg.CollectCountTbl)
	// 10金種配列に変換
	totalCollectCountTbl := calculation.NewCassette(totalCollectExCountTbl).ExCountTblToTenCountTbl()

	var totalCollect int
	for i, v := range totalCollectCountTbl {
		c.numInfoTbl = append(c.numInfoTbl, v)
		amount := v * domain.Cash[i]
		totalCollect += amount //合計金額を算出
	}
	// 総回収（合計）
	c.numInfoTbl = append(c.numInfoTbl, totalCollect)
	// 内回収済み（合計）
	c.numInfoTbl = append(c.numInfoTbl, totalBeforeCollect)

	// ログ出力
	c.outputLogNumInfoTbl(texCon)
	return c.numInfoTbl
}

// 日計レシートログ出力
func (c *summaryReport) outputLogNumInfoTbl(texCon *domain.TexContext) {
	l := fmt.Sprintf("【%v】日計レシート作成データ ---", texCon.GetUniqueKey())
	l += "  （売上金）"
	l += fmt.Sprintf("  %v : %+v", "現金", c.numInfoTbl[0:2])
	l += fmt.Sprintf("  %v : %+v", "クレジット", c.numInfoTbl[2:4])
	l += fmt.Sprintf("  %v : %+v", "電子マネー", c.numInfoTbl[4:6])
	l += fmt.Sprintf("  %v : %+v", "ＱＲコード", c.numInfoTbl[6:8])
	l += fmt.Sprintf("  %v : %+v", "その他", c.numInfoTbl[8:10])
	l += fmt.Sprintf("  %v : %+v", "合計", c.numInfoTbl[10:12])
	l += "  （入出金）"
	l += fmt.Sprintf("  %v : %+v 合計 : %d", "準備", c.numInfoTbl[12:22], c.numInfoTbl[22])
	l += fmt.Sprintf("  %v : %+v 合計 : %d", "入金", c.numInfoTbl[23:33], c.numInfoTbl[33])
	l += fmt.Sprintf("  %v : %+v 合計 : %d", "払出", c.numInfoTbl[34:44], c.numInfoTbl[44])
	l += fmt.Sprintf("  %v : %+v 合計 : %d", "補充", c.numInfoTbl[45:55], c.numInfoTbl[55])
	l += fmt.Sprintf("  %v : %+v 合計 : %d", "回収", c.numInfoTbl[56:66], c.numInfoTbl[66])
	l += "  （締め情報）"
	l += fmt.Sprintf("  %v : %+v 合計 : %d", "処理前", c.numInfoTbl[67:77], c.numInfoTbl[77])
	l += fmt.Sprintf("  %v : %+v 合計 : %d", "出金", c.numInfoTbl[78:88], c.numInfoTbl[88])
	l += fmt.Sprintf("  %v : %+v 合計 : %d", "締め後", c.numInfoTbl[89:99], c.numInfoTbl[99])
	l += fmt.Sprintf("  %v : %+v 合計 : %d", "総回収", c.numInfoTbl[100:110], c.numInfoTbl[110])
	l += fmt.Sprintf("  %v : %+v", "内回収済", c.numInfoTbl[111])
	c.logger.Debug("%v", l)
}
