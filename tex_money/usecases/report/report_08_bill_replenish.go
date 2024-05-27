package report

import (
	"fmt"
	"tex_money/domain"
	"tex_money/domain/handler"
)

type billReplenishReport struct {
	agg        domain.AggregateSafeInfo
	logger     handler.LoggerRepository
	numInfoTbl []int
}

type BillReplenishReportRepository interface {
	GetBillReplenishReport(texCon *domain.TexContext) []int // 紙幣補充:情報取得
}

func NewBillReplenishReport(agg domain.AggregateSafeInfo, logger handler.LoggerRepository) BillReplenishReportRepository {
	return &billReplenishReport{
		agg:        agg,
		logger:     logger,
		numInfoTbl: make([]int, 0, 19),
	}
}

// 紙幣補充:情報取得
func (c *billReplenishReport) GetBillReplenishReport(texCon *domain.TexContext) []int {

	//現枚数
	c.setBillData(c.agg.BeforeAmountCountTbl)
	//入金
	c.setBillData(c.agg.ReplenishCountTbl)
	//出金
	c.setBillData(c.agg.CollectCountTbl)
	//入出金後
	c.setBillData(c.agg.AfterAmountCountTbl)

	//入金額計
	c.setTotalBillAmount(c.agg.ReplenishCountTbl)
	// 出金額計
	c.setTotalBillAmount(c.agg.CollectCountTbl)

	// 紙幣補充ログ出力
	c.outputLogNumInfoTbl(texCon)
	return c.numInfoTbl
}

// 26金種配列から紙幣情報をnumInfoTblにセットする
func (c *billReplenishReport) setBillData(aggTbl [26]int) {
	for i, v := range aggTbl {
		if i == 4 {
			break
		}
		c.numInfoTbl = append(c.numInfoTbl, v)
	}
}

// 26金種配列のから紙幣の合計金額をnumInfoTblにセットする
func (c *billReplenishReport) setTotalBillAmount(aggTbl [26]int) {
	var total int
	for i, v := range aggTbl {
		if i == 4 {
			break
		}
		total += v * domain.AllCashInMachineTwentySix[i] //合計を算出
	}
	c.numInfoTbl = append(c.numInfoTbl, total)
}

// 紙幣補充ログ出力
func (c *billReplenishReport) outputLogNumInfoTbl(texCon *domain.TexContext) {
	l := fmt.Sprintf("【%v】紙幣逆両替レポート作成データ ---", texCon.GetUniqueKey())
	l += fmt.Sprintf("  %v : %+v", "現枚数", c.numInfoTbl[0:4])
	l += fmt.Sprintf("  %v : %+v", "入金", c.numInfoTbl[4:8])
	l += fmt.Sprintf("  %v : %+v", "出金", c.numInfoTbl[8:12])
	l += fmt.Sprintf("  %v : %+v", "入出金後", c.numInfoTbl[12:16])
	l += fmt.Sprintf("  %v : %+v", "入金額計", c.numInfoTbl[16])
	l += fmt.Sprintf("  %v : %+v", "出金額計", c.numInfoTbl[17])
	c.logger.Debug("%v", l)
}
