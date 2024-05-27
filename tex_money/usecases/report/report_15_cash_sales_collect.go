package report

import (
	"fmt"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/pkg/calculation"
)

type cashSalesCollectReport struct {
	agg        domain.AggregateSafeInfo
	logger     handler.LoggerRepository
	numInfoTbl []int
}

type CashSalesCollectReportRepository interface {
	GetCashSalesCollectReport(texCon *domain.TexContext, salesAmount int) []int // 現金売上金回収
}

func NewCashSalesCollectReport(agg domain.AggregateSafeInfo, logger handler.LoggerRepository) CashSalesCollectReportRepository {
	return &cashSalesCollectReport{
		agg:        agg,
		logger:     logger,
		numInfoTbl: make([]int, 0, 28),
	}
}

// 現金売上金回収
func (c *cashSalesCollectReport) GetCashSalesCollectReport(texCon *domain.TexContext, salesAmount int) []int {
	// FIT-A ======================================================================
	//枚数
	for i, s := range c.agg.SalesCollectCountTbl {
		if i == 10 {
			break
		}
		c.numInfoTbl = append(c.numInfoTbl, s)
	}
	c.numInfoTbl = append(c.numInfoTbl, c.agg.SalesCollectCountTbl[11], c.agg.SalesCollectCountTbl[13], c.agg.SalesCollectCountTbl[15])

	//金額
	limit := len(domain.AllCashInMachine)
	for i, e := range c.numInfoTbl {
		if i > limit {
			break
		}
		c.numInfoTbl = append(c.numInfoTbl, e*domain.AllCashInMachine[i])
	}

	// レポート保持情報の売上金回収のあふれ金庫を除いた金額をセット
	salesCollect := calculation.NewCassette(c.agg.SalesCollectCountTbl)
	salesCollectAmount := salesCollect.GetPayableCurrencyCount()
	//回収計
	/// あふれ金庫回収分を抜いた額に変更する
	c.numInfoTbl = append(c.numInfoTbl, salesCollectAmount)

	//売上金
	c.numInfoTbl = append(c.numInfoTbl, salesAmount)

	// FIT-B ======================================================================
	//現枚数
	// 10金種配列に変換
	beforeCountTbl := calculation.NewCassette(c.agg.BeforeAmountCountTbl).ExCountTblToTenCountTbl()
	c.numInfoTbl = append(c.numInfoTbl, beforeCountTbl[:]...)

	//入金
	inCountTbl := calculation.NewCassette(c.agg.ReplenishCountTbl).ExCountTblToTenCountTbl()
	c.numInfoTbl = append(c.numInfoTbl, inCountTbl[:]...)

	//出金
	outCountTbl := calculation.NewCassette(c.agg.CollectCountTbl).ExCountTblToTenCountTbl()
	c.numInfoTbl = append(c.numInfoTbl, outCountTbl[:]...)

	//入出金後
	afterCountTbl := calculation.NewCassette(c.agg.AfterAmountCountTbl).ExCountTblToTenCountTbl()
	c.numInfoTbl = append(c.numInfoTbl, afterCountTbl[:]...)

	//入金額計
	totalIn := calculation.NewCassette(c.agg.ReplenishCountTbl).GetTotalAmount()
	c.numInfoTbl = append(c.numInfoTbl, totalIn)

	// 出金額計
	totalOut := calculation.NewCassette(c.agg.CollectCountTbl).GetTotalAmount()
	c.numInfoTbl = append(c.numInfoTbl, totalOut)

	//回収計 // 出金額計と同じ
	c.numInfoTbl = append(c.numInfoTbl, totalOut)

	//売上金
	c.numInfoTbl = append(c.numInfoTbl, salesAmount)

	c.outputLogNumInfoTbl(texCon)
	return c.numInfoTbl
}

// 現金売上金回収レポートログ出力
func (c *cashSalesCollectReport) outputLogNumInfoTbl(texCon *domain.TexContext) {
	l := fmt.Sprintf("【%v】現金売上金回収レポート作成データ ---", texCon.GetUniqueKey())
	l += "  (FIT-A)"
	l += fmt.Sprintf("  %v : %+v", "枚数", c.numInfoTbl[0:13])
	l += fmt.Sprintf("  %v : %+v", "金額", c.numInfoTbl[13:26])
	l += fmt.Sprintf("  %v : %+v", "回収計", c.numInfoTbl[26])
	l += fmt.Sprintf("  %v : %+v", "売上計", c.numInfoTbl[27])
	l += "  |  (FIT-B)"
	l += fmt.Sprintf("  %v : %+v", "現枚数", c.numInfoTbl[28:38])
	l += fmt.Sprintf("  %v : %+v", "入金", c.numInfoTbl[38:48])
	l += fmt.Sprintf("  %v : %+v", "出金", c.numInfoTbl[48:58])
	l += fmt.Sprintf("  %v : %+v", "入出金後", c.numInfoTbl[58:68])
	l += fmt.Sprintf("  %v : %+v", "入金額計", c.numInfoTbl[68])
	l += fmt.Sprintf("  %v : %+v", "出金額計", c.numInfoTbl[69])
	l += fmt.Sprintf("  %v : %+v", "回収計", c.numInfoTbl[70])
	l += fmt.Sprintf("  %v : %+v", "売上計", c.numInfoTbl[71])
	c.logger.Debug("%v", l)
}
