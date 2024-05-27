package report

import (
	"fmt"
	"tex_money/domain"
	"tex_money/domain/handler"
)

type summarySalesReport struct {
	logger     handler.LoggerRepository
	numInfoTbl []int
}

type SummarySalesReportRepository interface {
	GetSummarySalesReport(texCon *domain.TexContext, amount []int, count []int, totalAmount int, totalCount int) []int // 精算機別日計表
}

func NewSummarySalesReport(logger handler.LoggerRepository) SummarySalesReportRepository {
	return &summarySalesReport{
		logger:     logger,
		numInfoTbl: make([]int, 0, 12),
	}
}

// 精算機別日計表
func (c *summarySalesReport) GetSummarySalesReport(texCon *domain.TexContext, amount []int, count []int, totalAmount int, totalCount int) []int {
	//売上金額　現金　クレジット　QRコード
	for i := range amount {
		if i > 2 {
			break
		}
		c.numInfoTbl = append(c.numInfoTbl, amount[i])
		c.numInfoTbl = append(c.numInfoTbl, count[i])

	}
	//売上金額合計
	c.numInfoTbl = append(c.numInfoTbl, totalAmount)
	//売上回数合計
	c.numInfoTbl = append(c.numInfoTbl, totalCount)

	//売上金額　電子マネー　J-Debit
	for i := range amount {
		if i < 3 {
			continue
		}
		c.numInfoTbl = append(c.numInfoTbl, amount[i])
		c.numInfoTbl = append(c.numInfoTbl, count[i])

	}

	c.outputLogNumInfoTbl(texCon)
	return c.numInfoTbl
}

// 精算機別日計表ログ出力
func (c *summarySalesReport) outputLogNumInfoTbl(texCon *domain.TexContext) {
	l := fmt.Sprintf("【%v】精算機別日計表作成データ ---", texCon.GetUniqueKey())
	l += fmt.Sprintf("  %v : %+v", "現金", c.numInfoTbl[0:2])
	l += fmt.Sprintf("  %v : %+v", "クレジット", c.numInfoTbl[2:4])
	l += fmt.Sprintf("  %v : %+v", "ＱＲコード", c.numInfoTbl[4:6])
	l += fmt.Sprintf("  %v : %+v", "合計", c.numInfoTbl[6:8])
	l += fmt.Sprintf("  %v : %+v", "電子マネー", c.numInfoTbl[8:10])
	l += fmt.Sprintf("  %v : %+v", "J-Debit", c.numInfoTbl[10:12])
	c.logger.Debug("%v", l)
}
