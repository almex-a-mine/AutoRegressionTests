package report

import (
	"fmt"
	"tex_money/domain"
	"tex_money/domain/handler"
)

type coinUnitDiffReport struct {
	logger     handler.LoggerRepository
	numInfoTbl []int
}

type CoinUnitDiffReportRepository interface {
	GetCoinUnitDiffReport(texCon *domain.TexContext, differenceExCountTbl [26]int) []int
}

func NewCoinUnitDiffReport(logger handler.LoggerRepository) CoinUnitDiffReportRepository {
	return &coinUnitDiffReport{
		logger:     logger,
		numInfoTbl: make([]int, 0, 12),
	}
}

// 硬貨ユニット補充差分レポート(補充予定枚数を印字する)
func (c *coinUnitDiffReport) GetCoinUnitDiffReport(texCon *domain.TexContext, differenceExCountTbl [26]int) []int {
	//各金種の釣銭初期枚数と現在有高との差額をレシートに印字する

	//枚数
	// cassetteType := 3 //画面的にメインカセットとサブカセットが存在する為
	// resultCassette := c.coinCassetteMng.Exchange(texCon, cassetteType)
	for i, d := range differenceExCountTbl {
		if i > 15 {
			break
		}
		if (i >= 4 && i < 10) || i == 11 || i == 13 || i == 15 {
			c.numInfoTbl = append(c.numInfoTbl, d)
		}
	}

	//金額
	for i, e := range c.numInfoTbl {
		c.numInfoTbl = append(c.numInfoTbl, e*domain.AllCashInMachineNine[i])
	}

	// ログ出力
	c.outputLogNumInfoTbl(texCon)
	return c.numInfoTbl
}

// 硬貨ユニット補充差分レポートログ出力
func (c *coinUnitDiffReport) outputLogNumInfoTbl(texCon *domain.TexContext) {
	l := fmt.Sprintf("【%v】硬貨ユニット補充差分レポート作成データ ---", texCon.GetUniqueKey())
	l += fmt.Sprintf("  %v : %+v", "枚数", c.numInfoTbl[0:9])
	l += fmt.Sprintf("  %v : %+v", "金額", c.numInfoTbl[9:18])
	c.logger.Debug("%v", l)
}
