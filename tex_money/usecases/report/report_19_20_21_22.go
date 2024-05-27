package report

import (
	"fmt"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/pkg/calculation"
)

type replenishReport struct {
	agg        domain.AggregateSafeInfo
	logger     handler.LoggerRepository
	numInfoTbl []int
}

type ReplenishReportRepository interface {
	GetReplenishReport(texCon *domain.TexContext) []int
}

func NewReplenishReport(agg domain.AggregateSafeInfo, logger handler.LoggerRepository) ReplenishReportRepository {
	return &replenishReport{
		agg:        agg,
		logger:     logger,
		numInfoTbl: make([]int, 0, 43),
	}
}

// 補充レポート（追加補充/回収庫から回収/指定枚数回収/逆両替）
func (c *replenishReport) GetReplenishReport(texCon *domain.TexContext) []int {

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

	// ログ出力
	c.outputLogNumInfoTbl(texCon)
	return c.numInfoTbl
}

// 補充レポートログ出力
func (c *replenishReport) outputLogNumInfoTbl(texCon *domain.TexContext) {
	l := fmt.Sprintf("【%v】補充レポート作成データ ---", texCon.GetUniqueKey())
	l += fmt.Sprintf("  %v : %+v", "現枚数", c.numInfoTbl[0:10])
	l += fmt.Sprintf("  %v : %+v", "入金", c.numInfoTbl[10:20])
	l += fmt.Sprintf("  %v : %+v", "出金", c.numInfoTbl[20:30])
	l += fmt.Sprintf("  %v : %+v", "入出金後", c.numInfoTbl[30:40])
	l += fmt.Sprintf("  %v : %+v", "入金額計", c.numInfoTbl[40])
	l += fmt.Sprintf("  %v : %+v", "出金額計", c.numInfoTbl[41])
	c.logger.Debug("%v", l)
}
