package usecases

import (
	"reflect"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/usecases/report"
)

type printDataManager struct {
	logger             handler.LoggerRepository
	config             config.Configuration
	syslogMng          SyslogManager
	errorMng           ErrorManager
	summarySales       domain.SummarySales //精算機日計表
	aggregateMng       AggregateManager
	safeInfoMng        SafeInfoManager
	coinCassetteMng    CoinCassetteControlManager
	maintenanceModeMng MaintenanceModeManager
	texmyHandler       TexMoneyHandlerRepository
	// changeReserveCountTbl []int //釣銭初期枚数
}

// 印刷データ保存
func NewPrintDataManager(logger handler.LoggerRepository,
	config config.Configuration,
	syslogMng SyslogManager,
	errorMng ErrorManager,
	aggregateMng AggregateManager,
	safeInfoMng SafeInfoManager,
	coinCassetteMng CoinCassetteControlManager,
	maintenanceModeMng MaintenanceModeManager,
	texmyHandler TexMoneyHandlerRepository) PrintDataManager {
	return &printDataManager{
		logger:             logger,
		config:             config,
		syslogMng:          syslogMng,
		errorMng:           errorMng,
		summarySales:       domain.SummarySales{},
		aggregateMng:       aggregateMng,
		safeInfoMng:        safeInfoMng,
		coinCassetteMng:    coinCassetteMng,
		maintenanceModeMng: maintenanceModeMng,
		texmyHandler:       texmyHandler,
	}
}

// 精算機別日計表:情報設定
func (c *printDataManager) SetSummarySales(resInfo domain.ResultGetSalesinfo) {
	c.logger.Trace("START:printDataManager SetSummarySales resInfo.InfoSalesTbl=%+v)", resInfo.InfoSales.SalesTypeTbl)

	c.summarySales.Amount = make([]int, 5)
	c.summarySales.Count = make([]int, 5)
	c.summarySales.TotalAmount = 0
	c.summarySales.TotalCount = 0
	for i := 0; i < len(resInfo.InfoSales.SalesTypeTbl); i++ {
		switch resInfo.InfoSales.SalesTypeTbl[i].PaymentType {

		case 0: //現金
			c.setSummarySales(0, resInfo.InfoSales.SalesTypeTbl[i].Amount, resInfo.InfoSales.SalesTypeTbl[i].Count)
		case 1: //クレジット
			c.setSummarySales(1, resInfo.InfoSales.SalesTypeTbl[i].Amount, resInfo.InfoSales.SalesTypeTbl[i].Count)
		case 2: //Jデビット
			c.setSummarySales(4, resInfo.InfoSales.SalesTypeTbl[i].Amount, resInfo.InfoSales.SalesTypeTbl[i].Count)
		case 3: //QRコード決済
			c.setSummarySales(2, resInfo.InfoSales.SalesTypeTbl[i].Amount, resInfo.InfoSales.SalesTypeTbl[i].Count)
		case 4: //電子マネー
			c.setSummarySales(3, resInfo.InfoSales.SalesTypeTbl[i].Amount, resInfo.InfoSales.SalesTypeTbl[i].Count)
		case 5: //その他
		}

	}

	c.logger.Trace("END:printDataManager SetSummarySales c.summarySales=%+v", c.summarySales)

}

func (c *printDataManager) setSummarySales(setNo int, amount int, count int) {
	c.summarySales.Amount[setNo] += amount
	c.summarySales.Count[setNo] += count
	c.summarySales.TotalAmount += amount
	c.summarySales.TotalCount += count

}

// 精算機別日計表:情報取得
func (c *printDataManager) GetSummarySales() ([]int, []int, int, int) {
	c.logger.Trace("START:printDataManager GetSummarySales")
	c.logger.Trace("END:printDataManager GetSummarySales c.summarySales=%+v", c.summarySales)
	return c.summarySales.Amount, c.summarySales.Count, c.summarySales.TotalAmount, c.summarySales.TotalCount
}

// 精算機別日計表
func (c *printDataManager) ReportSummarySales(texCon *domain.TexContext) []int {
	c.logger.Trace("【%v】START:精算機別日計表データ作成", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:精算機別日計表データ作成", texCon.GetUniqueKey())

	// 汎用数値情報の生成
	amount, count, totalAmount, totalCount := c.GetSummarySales()
	c.logger.Debug("【%v】- amount=%v, count=%v, totalAmount=%v, totalCount=%v", texCon.GetUniqueKey(), amount, count, totalAmount, totalCount)

	return report.NewSummarySalesReport(c.logger).GetSummarySalesReport(texCon, amount, count, totalAmount, totalCount)
}

// キャッシュカウントレポート
func (c *printDataManager) ReportCashCount(texCon *domain.TexContext) []int {
	c.logger.Trace("【%v】START:キャッシュカウントレポートデータ作成", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:キャッシュカウントレポートデータ作成", texCon.GetUniqueKey())

	//レポート用処理後有高金種配列にデータを格納
	var countTbl [domain.EXTRA_CASH_TYPE_SHITEI]int
	_, sortInfoTbl := c.safeInfoMng.GetSortInfo(texCon, 0)
	_, agg := c.aggregateMng.GetAggregateSafeInfo(texCon, c.maintenanceModeMng.GetMode(texCon))
	for i, v := range sortInfoTbl.ExCountTbl {
		if i >= 4 && i <= 15 {
			//硬貨データは「締め前有高配列 + 補充入金配列 - 回収配列」を格納
			countTbl[i] = agg.BeforeAmountCountTbl[i] + agg.ReplenishCountTbl[i] - agg.CollectCountTbl[i]
			continue
		}
		//この時点の有高を格納
		countTbl[i] = v
	}
	c.aggregateMng.UpdateAggregateCountTbl(texCon, c.maintenanceModeMng.GetMode(texCon), domain.AFTER_AMOUNT_COUNT_TBL, countTbl)

	// レポート用データ再取得
	_, agg = c.aggregateMng.GetAggregateSafeInfo(texCon, domain.CLOSING_MODE)

	return report.NewCashCountReport(agg, c.logger).GetCashCountReport(texCon)
}

// 紙幣補充
func (c *printDataManager) ReportSupplyBill(texCon *domain.TexContext) []int {
	c.logger.Trace("【%v】START:紙幣逆両替レポート作成データ", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:紙幣逆両替レポート作成データ", texCon.GetUniqueKey())
	//入金後の有高データ保存
	_, sortInfoTbl := c.safeInfoMng.GetSortInfo(texCon, 0)
	c.aggregateMng.UpdateAggregateCountTbl(texCon, c.maintenanceModeMng.GetMode(texCon), domain.AFTER_AMOUNT_COUNT_TBL, sortInfoTbl.ExCountTbl)

	// レポート用データ取得
	_, agg := c.aggregateMng.GetAggregateSafeInfo(texCon, domain.REPLENISHMENT_MODE)

	return report.NewBillReplenishReport(agg, c.logger).GetBillReplenishReport(texCon)
}

// 硬貨ユニット交換
// （通常硬貨ユニット交換 / 予備硬貨ユニット交換 / 全硬貨ユニット交換/ 硬貨手動追加）
func (c *printDataManager) ReportChangeCoinUnit(texCon *domain.TexContext, reportNo int) []int {
	c.logger.Trace("【%v】START:硬貨ユニット交換データ作成", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:硬貨ユニット交換データ作成", texCon.GetUniqueKey())

	//レポート用処理後有高金種配列に現在有高枚数を格納
	_, sortInfoTbl := c.safeInfoMng.GetSortInfo(texCon, 0)
	c.aggregateMng.UpdateAggregateCountTbl(texCon, c.maintenanceModeMng.GetMode(texCon), domain.AFTER_AMOUNT_COUNT_TBL, sortInfoTbl.ExCountTbl)

	// レポート用データ取得
	_, agg := c.aggregateMng.GetAggregateSafeInfo(texCon, domain.REPLENISHMENT_MODE)

	// 釣銭準備金枚数を配列に変換して取得
	changeReserveCountTbl := c.getChangeReserveCountArray(texCon, c.texmyHandler.GetMoneySetting().ChangeReserveCount)

	return report.NewCoinUnitReport(agg, c.logger, changeReserveCountTbl).GetCoinUnitReport(texCon, reportNo)
}

// 釣銭準備金枚数を配列にセット
func (c *printDataManager) getChangeReserveCountArray(texCon *domain.TexContext, changeReserveCount domain.ChangeReserveCount) (countArray []int) {
	// ChangeReserveCount構造体のリフレクションを取得
	rtCst := reflect.TypeOf(changeReserveCount)
	// ChangeReserveCount構造体の値を取得
	rvCst := reflect.ValueOf(changeReserveCount)

	// 構造体の情報を配列にセット
	for i := 0; i < rtCst.NumField(); i++ {
		if i < 2 { //枚数情報のみ取得
			continue
		}
		// フィールド名に対応する値を取得
		v := rvCst.FieldByName(rtCst.Field(i).Name).Interface()
		countArray = append(countArray, v.(int))
	}

	c.logger.Debug("【%v】printDataManager getChangeReserveCountArray countArray=%v", texCon.GetUniqueKey(), countArray)
	return
}

// 現金売上金回収
func (c *printDataManager) CashSalesCollectReport(texCon *domain.TexContext) []int {
	c.logger.Trace("【%v】START:現金売上金回収レポートデータ作成", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:現金売上金回収レポートデータ作成", texCon.GetUniqueKey())

	// 汎用数値情報の生成
	amount, _, _, _ := c.GetSummarySales()

	//入金後の有高データ保存
	_, sortInfoTbl := c.safeInfoMng.GetSortInfo(texCon, 0)
	c.aggregateMng.UpdateAggregateCountTbl(texCon, c.maintenanceModeMng.GetMode(texCon), domain.AFTER_AMOUNT_COUNT_TBL, sortInfoTbl.ExCountTbl)

	// レポート用データ再取得
	_, agg := c.aggregateMng.GetAggregateSafeInfo(texCon, domain.CLOSING_MODE)

	return report.NewCashSalesCollectReport(agg, c.logger).GetCashSalesCollectReport(texCon, amount[0])
}

// 硬貨ユニット補充差分レポート(補充予定枚数を印字する)
func (c *printDataManager) ReportCoinUnitDiff(texCon *domain.TexContext) []int {
	c.logger.Trace("【%v】START:硬貨ユニット補充差分レポートデータ作成", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:硬貨ユニット補充差分レポートデータ作成", texCon.GetUniqueKey())

	// 操作後の硬貨カセット情報取得
	// cassetteType= 3（メインカセットとサブカセット） を指定する
	resultCassette := c.coinCassetteMng.Exchange(texCon, 3)

	return report.NewCoinUnitDiffReport(c.logger).GetCoinUnitDiffReport(texCon, resultCassette.DifferenceExCountTbl)
}

// 精算機日計レシート:情報設定
func (c *printDataManager) SetReportSummarySalesInfo(resInfo domain.ResultGetSalesinfo) {
	c.logger.Trace("START:精算機日計レシート:情報設定 resInfo.InfoSalesTbl=%+v", resInfo.InfoSales.SalesTypeTbl)

	c.summarySales.Amount = make([]int, 5)
	c.summarySales.Count = make([]int, 5)
	c.summarySales.TotalAmount = 0
	c.summarySales.TotalCount = 0
	for i := 0; i < len(resInfo.InfoSales.SalesTypeTbl); i++ {
		// 売上種別=0:チェックイン のみ集計対象
		if resInfo.InfoSales.SalesTypeTbl[i].SalesType != 0 {
			continue
		}

		switch resInfo.InfoSales.SalesTypeTbl[i].PaymentType {
		case 0: //現金
			c.setSummarySales(0, resInfo.InfoSales.SalesTypeTbl[i].Amount, resInfo.InfoSales.SalesTypeTbl[i].Count)
		case 1: //クレジット
			c.setSummarySales(1, resInfo.InfoSales.SalesTypeTbl[i].Amount, resInfo.InfoSales.SalesTypeTbl[i].Count)
		case 2: //Jデビット
			c.setSummarySales(4, resInfo.InfoSales.SalesTypeTbl[i].Amount, resInfo.InfoSales.SalesTypeTbl[i].Count)
		case 3: //QRコード決済
			c.setSummarySales(3, resInfo.InfoSales.SalesTypeTbl[i].Amount, resInfo.InfoSales.SalesTypeTbl[i].Count)
		case 4: //電子マネー
			c.setSummarySales(2, resInfo.InfoSales.SalesTypeTbl[i].Amount, resInfo.InfoSales.SalesTypeTbl[i].Count)
		case 5: //その他
			c.setSummarySales(4, resInfo.InfoSales.SalesTypeTbl[i].Amount, resInfo.InfoSales.SalesTypeTbl[i].Count)
		}

	}

	c.logger.Trace("END:精算機日計レシート:情報設定 c.summarySales=%+v", c.summarySales)

}

// 日計レシート（FIT-B NEXTクリニック向け）
func (c *printDataManager) ReportSummary(texCon *domain.TexContext) []int {
	c.logger.Trace("【%v】START:日計レシートデータ作成", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:日計レシートデータ作成", texCon.GetUniqueKey())

	//レポート用処理後有高金種配列に現在有高を格納
	_, sortInfoTbl := c.safeInfoMng.GetSortInfo(texCon, 0)
	c.aggregateMng.UpdateAggregateCountTbl(texCon, c.maintenanceModeMng.GetMode(texCon), domain.AFTER_AMOUNT_COUNT_TBL, sortInfoTbl.ExCountTbl)

	// 現在有高取得
	safeInfo := c.safeInfoMng.GetSafeInfo(texCon)
	// レポート用データ取得
	_, agg := c.aggregateMng.GetAggregateSafeInfo(texCon, domain.CLOSING_MODE)

	// 汎用数値情報の生成
	amount, count, totalAmount, totalCount := c.GetSummarySales()
	c.logger.Debug("【%v】- amount=%v, count=%v, totalAmount=%v, totalCount=%v", texCon.GetUniqueKey(), amount, count, totalAmount, totalCount)

	return report.NewSummaryReport(agg, safeInfo, c.logger).GetReportSummary(texCon, amount, count, totalAmount, totalCount)
}

// 補充レポート（追加補充/回収庫から回収/指定枚数回収/逆両替）
func (c *printDataManager) ReportSupply(texCon *domain.TexContext) []int {
	c.logger.Trace("【%v】START:補充レポート作成データ", texCon.GetUniqueKey())
	defer c.logger.Trace("【%v】END:補充レポート作成データ", texCon.GetUniqueKey())
	//入金後の有高データ保存
	_, sortInfoTbl := c.safeInfoMng.GetSortInfo(texCon, 0)
	c.aggregateMng.UpdateAggregateCountTbl(texCon, c.maintenanceModeMng.GetMode(texCon), domain.AFTER_AMOUNT_COUNT_TBL, sortInfoTbl.ExCountTbl)

	// レポート用データ取得
	_, agg := c.aggregateMng.GetAggregateSafeInfo(texCon, domain.REPLENISHMENT_MODE)

	return report.NewReplenishReport(agg, c.logger).GetReplenishReport(texCon)
}
