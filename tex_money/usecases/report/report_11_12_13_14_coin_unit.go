package report

import (
	"fmt"
	"tex_money/domain"
	"tex_money/domain/handler"
)

type coinUnitReport struct {
	agg                   domain.AggregateSafeInfo
	logger                handler.LoggerRepository
	numInfoTbl            []int
	changeReserveCountTbl []int //釣銭初期枚数
}

type CoinUnitReportRepository interface {
	GetCoinUnitReport(texCon *domain.TexContext, reportNo int) []int // カセット交換レポート印刷用の配列を作成
}

func NewCoinUnitReport(agg domain.AggregateSafeInfo, logger handler.LoggerRepository, changeReserveCountTbl []int) CoinUnitReportRepository {
	return &coinUnitReport{
		agg:                   agg,
		logger:                logger,
		numInfoTbl:            make([]int, 0, 54),
		changeReserveCountTbl: changeReserveCountTbl,
	}
}

// カセット交換レポート印刷用の配列を作成
// （通常硬貨ユニット交換 / 予備硬貨ユニット交換 / 全硬貨ユニット交換/ 硬貨手動追加）
// exCountTbl
// [0]~[12]：締め前（紙幣/青カセット/ピンクカセット） [13]~[25]：入金（紙幣/青カセット/ピンクカセット）
// [26]~[38]：出金（紙幣/青カセット/ピンクカセット）[39]~[51]：入出金後（紙幣/青カセット/ピンクカセット）
// [52]：入金金庫合計 [53]：出金額合計
func (c *coinUnitReport) GetCoinUnitReport(texCon *domain.TexContext, reportNo int) []int {
	c.logger.Trace("【%v】START:printDataManager GetCoinUnitReportData", texCon.GetUniqueKey())

	// // レポート用金庫情報取得
	// r, agg := c.aggregateMng.GetAggregateSafeInfo(texCon, domain.REPLENISHMENT_MODE)
	// if !r {
	// 	return
	// }
	// // 釣銭準備金枚数を配列に変換して格納
	// c.changeReserveCountTbl = c.getChangeReserveCountArray(texCon, changeReserveCount)

	//現枚数
	for i, b := range c.agg.BeforeAmountCountTbl {
		if i > 15 {
			break
		}
		//予備カセットは1系のみ
		if i == 10 || i == 12 || i == 14 {
			continue
		}
		c.numInfoTbl = append(c.numInfoTbl, b)
	}

	//入金枚数
	for i, r := range c.agg.ReplenishCountTbl {
		if i > 15 {
			break
		}
		//予備カセットは1系のみ
		if i == 10 || i == 12 || i == 14 {
			continue
		}
		c.numInfoTbl = append(c.numInfoTbl, r)
	}

	//出金
	for i, col := range c.agg.CollectCountTbl {
		if i > 15 {
			break
		}
		//予備カセットは1系のみ
		if i == 10 || i == 12 || i == 14 {
			continue
		}
		c.numInfoTbl = append(c.numInfoTbl, col)
	}

	//入出金後
	for i, a := range c.agg.AfterAmountCountTbl {
		if i > 15 {
			break
		}
		//予備カセットは1系のみ
		if i == 10 || i == 12 || i == 14 {
			continue
		}
		c.numInfoTbl = append(c.numInfoTbl, a)
	}

	// レポート毎にデータを加工
	switch reportNo {
	case domain.CHANGE_COINTUNIT1: //通常硬貨ユニット交換(青カセット交換)
		c.setChangeNormalCoinUnitReport()
	case domain.CHANGE_COINTUNIT2: //予備硬貨ユニット交換(ピンクカセット交換)
		c.setChangeSubCoinUnitReport()
	case domain.CHANGE_COINTUNIT_ALL: //全硬貨ユニット交換（青ピンクカセット交換）
		c.setChangeAllCoinUnitReport()
	case domain.SUPPLY_COIN_MANUAL: //硬貨手動追加
	}

	//入金額計
	totalInAmount := 0
	for i, count := range c.numInfoTbl[13:26] {
		totalInAmount += count * domain.AllCashInMachine[i]
	}
	c.numInfoTbl = append(c.numInfoTbl, totalInAmount)
	//出金額計
	totalOutAmount := 0
	for i, count := range c.numInfoTbl[26:39] {
		totalOutAmount += count * domain.AllCashInMachine[i]
	}
	c.numInfoTbl = append(c.numInfoTbl, totalOutAmount)

	// ログ出力
	c.outputLogNumInfoTbl(texCon)
	return c.numInfoTbl
}

// 通常硬貨ユニット交換レポートのデータをセット
func (c *coinUnitReport) setChangeNormalCoinUnitReport() {
	n, j := 0, 0
	for i := 0; i < len(c.numInfoTbl); i++ {
		// 青カセット（入金）に釣銭初期枚数をセット
		if i >= 17 && i <= 22 {
			c.numInfoTbl[i] = c.changeReserveCountTbl[n+4] //青カセットの釣銭初期枚数[4]~[9]を格納
			n++
		}
		// 青カセット（出金）に現枚数をセット
		if i >= 30 && i <= 35 {
			// 有高＋出金金額の場合、有高が釣銭可能枚数より多い場合に想定値と異なる計算となる為
			// 有高＋入金-入出金後から計算する仕様に変更する
			// c.numInfoTbl[i] += c.numInfoTbl[4+j]
			c.numInfoTbl[i] = c.numInfoTbl[4+j] + c.numInfoTbl[17+j] - c.numInfoTbl[43+j]
			j++
		}
	}
	// 青カセット以外の入金情報をクリア
	c.numInfoTbl[13] = 0 // 紙幣
	c.numInfoTbl[14] = 0 // 紙幣
	c.numInfoTbl[15] = 0 // 紙幣
	c.numInfoTbl[16] = 0 // 紙幣
	c.numInfoTbl[23] = 0 // ピンク
	c.numInfoTbl[24] = 0 // ピンク
	c.numInfoTbl[25] = 0 // ピンク
}

// 予備硬貨ユニット交換レポートのデータをセット
func (c *coinUnitReport) setChangeSubCoinUnitReport() {
	n, j := 0, 0
	for i := 0; i < len(c.numInfoTbl); i++ {
		// ピンクカセット（入金）に釣銭初期枚数をセット
		if i >= 23 && i <= 25 {
			c.numInfoTbl[i] = c.changeReserveCountTbl[n+11] //ピンクセットの釣銭初期枚数[11]、[13]、[15]を格納
			n += 2
		}
		// ピンクセット（出金）に現枚数をセット
		if i >= 36 && i <= 38 {
			// 有高＋出金金額の場合、有高が釣銭可能枚数より多い場合に想定値と異なる計算となる為
			// 有高＋入金-入出金後から計算する仕様に変更する
			// c.numInfoTbl[i] += c.numInfoTbl[10+j]
			c.numInfoTbl[i] = c.numInfoTbl[10+j] + c.numInfoTbl[23+j] - c.numInfoTbl[49+j]
			j++
		}
	}
	// ピンクカセット以外の入金情報をクリア
	c.numInfoTbl[13] = 0 // 紙幣
	c.numInfoTbl[14] = 0 // 紙幣
	c.numInfoTbl[15] = 0 // 紙幣
	c.numInfoTbl[16] = 0 // 紙幣
	c.numInfoTbl[17] = 0 // 青
	c.numInfoTbl[18] = 0 // 青
	c.numInfoTbl[19] = 0 // 青
	c.numInfoTbl[20] = 0 // 青
	c.numInfoTbl[21] = 0 // 青
	c.numInfoTbl[22] = 0 // 青
}

// 全硬貨ユニット交換レポートのデータをセット
func (c *coinUnitReport) setChangeAllCoinUnitReport() {
	n, m, j := 0, 0, 0
	for i := 0; i < len(c.numInfoTbl); i++ {
		// 青カセット（入金）に釣銭初期枚数をセット
		if i >= 17 && i <= 22 {
			c.numInfoTbl[i] = c.changeReserveCountTbl[n+4] //青カセットの釣銭初期枚数[4]~[9]を格納
			n++
		}
		// ピンクカセット（入金）に釣銭初期枚数をセット
		if i >= 23 && i <= 25 {
			c.numInfoTbl[i] = c.changeReserveCountTbl[m+11] //ピンクセットの釣銭初期枚数[11]、[13]、[15]を格納
			m += 2
		}
		// 青カセット（出金）とピンクカセット（出金）に現枚数をセット
		if i >= 30 && i <= 38 {
			// 有高＋出金金額の場合、有高が釣銭可能枚数より多い場合に想定値と異なる計算となる為
			// 有高＋入金-入出金後から計算する仕様に変更する
			// c.numInfoTbl[i] += c.numInfoTbl[4+j]
			c.numInfoTbl[i] = c.numInfoTbl[4+j] + c.numInfoTbl[17+j] - c.numInfoTbl[43+j]
			j++
		}
	}

	// 硬貨カセット以外の入金情報をクリア
	c.numInfoTbl[13] = 0 // 紙幣
	c.numInfoTbl[14] = 0 // 紙幣
	c.numInfoTbl[15] = 0 // 紙幣
	c.numInfoTbl[16] = 0 // 紙幣
}

// カセット交換レポートログ出力
func (c *coinUnitReport) outputLogNumInfoTbl(texCon *domain.TexContext) {
	l := fmt.Sprintf("【%v】カセット交換レポート作成データ ---", texCon.GetUniqueKey())
	l += fmt.Sprintf("  %v : %+v", "現枚数", c.numInfoTbl[0:13])
	l += fmt.Sprintf("  %v : %+v", "入金", c.numInfoTbl[13:26])
	l += fmt.Sprintf("  %v : %+v", "出金", c.numInfoTbl[26:39])
	l += fmt.Sprintf("  %v : %+v", "入出金後", c.numInfoTbl[39:52])
	l += fmt.Sprintf("  %v : %+v", "入金額計", c.numInfoTbl[52])
	l += fmt.Sprintf("  %v : %+v", "出金額計", c.numInfoTbl[53])
	c.logger.Debug("%v", l)
}
