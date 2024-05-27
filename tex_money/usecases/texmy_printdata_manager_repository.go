package usecases

import "tex_money/domain"

// 印刷データ
type PrintDataManager interface {
	SetSummarySales(resInfo domain.ResultGetSalesinfo)                  //精算機別日計表情報設定
	GetSummarySales() ([]int, []int, int, int)                          //精算機別日計表:情報取得
	SetReportSummarySalesInfo(resInfo domain.ResultGetSalesinfo)        // 精算機日計レシート:情報設定
	ReportSummarySales(texCon *domain.TexContext) []int                 // 精算機日計表
	ReportCashCount(texCon *domain.TexContext) []int                    // キャッシュカウントレポート
	ReportSupplyBill(texCon *domain.TexContext) []int                   // 紙幣補充
	ReportChangeCoinUnit(texCon *domain.TexContext, reportNo int) []int // 硬貨ユニット交換
	CashSalesCollectReport(texCon *domain.TexContext) []int             // 現金売上金回収
	ReportCoinUnitDiff(texCon *domain.TexContext) []int                 // 硬貨ユニット補充差分レポート
	ReportSummary(texCon *domain.TexContext) []int                      // 精算機日計レシート（FIT-B NEXTクリニック向け）
	ReportSupply(texCon *domain.TexContext) []int                       // 補充レポート（追加補充/回収庫から回収/指定枚数回収/逆両替）
}
