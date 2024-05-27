package usecases

import (
	"fmt"
	"strconv"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/pkg/calculation"
)

type safeInfoManager struct {
	logger                     handler.LoggerRepository
	cfg                        config.Configuration
	syslogMng                  SyslogManager
	errorMng                   ErrorManager
	SafeInfo                   domain.SafeInfo
	iniService                 IniServiceRepository
	depositNumber              domain.SortInfoTbl
	BeforeReplenishmentBalance domain.SortInfoTbl
	deviceCashAvailable        domain.SortInfoTbl
	logicalCashAvailable       domain.SortInfoTbl
}

// 金庫情報保存
func NewSafeInfoManager(logger handler.LoggerRepository, cfg config.Configuration, syslogMng SyslogManager, errorMng ErrorManager, iniService IniServiceRepository) SafeInfoManager {
	c := &safeInfoManager{
		logger:     logger,
		cfg:        cfg,
		syslogMng:  syslogMng,
		errorMng:   errorMng,
		iniService: iniService,
		SafeInfo:   domain.SafeInfo{},
	}
	// 金庫情報初期化
	c.InitSafeInfo()
	return c
}

//============================================//
//分類情報種別:0:現金有高 1:釣銭可能 2:初期補充 3:取引入金 4:取引出金 5:取引差引 6:補充入金 7:補充出金 8:補充差引 9:売上金回収
// CASH_AVAILABLE           = 0
// CHANGE_AVAILABLE         = 1
// INITIAL_REPLENISHMENT    = 2
// TRANSACTION_DEPOSIT      = 3
// TRANSACTION_WITHDRAWAL   = 4
// TRANSACTION_BALANCE      = 5
// REPLENISHMENT_DEPOSIT    = 6
// REPLENISHMENT_WITHDRAWAL = 7
// REPLENISHMENT_BALANCE    = 8
// SALES_MONEY_COLLECT      = 9
//============================================//

func (c *safeInfoManager) GetInfoSafeSortTypeName(i int) string {
	switch i {
	case domain.CASH_AVAILABLE:
		return "0:現金有高"
	case domain.CHANGE_AVAILABLE:
		return "1:釣銭枚数"
	case domain.INITIAL_REPLENISHMENT:
		return "2:初期補充"
	case domain.TRANSACTION_DEPOSIT:
		return "3:取引入金"
	case domain.TRANSACTION_WITHDRAWAL:
		return "4:取引出金"
	case domain.TRANSACTION_BALANCE:
		return "5:取引差引"
	case domain.REPLENISHMENT_DEPOSIT:
		return "6:補充入金"
	case domain.REPLENISHMENT_WITHDRAWAL:
		return "7:補充出金"
	case domain.REPLENISHMENT_BALANCE:
		return "8:補充差引"
	case domain.SALES_MONEY_COLLECT:
		return "9:売上回収"
	case domain.DEPOSIT_NUMBER:
		return "10:入金可能"
	}
	return strconv.Itoa(i)
}

//============================================//
//金庫情報操作
//============================================//

// InitSafeInfo 金庫情報初期化
func (c *safeInfoManager) InitSafeInfo() {
	// 分類情報種別にデフォルト値をセット
	c.SafeInfo = domain.SafeInfo{}
	for i := range c.SafeInfo.SortInfoTbl {
		c.SafeInfo.SortInfoTbl[i].SortType = domain.SortType[i]
	}
	// 売上金回収回数にiniの値をセット
	c.SafeInfo.SalesCompleteCount = c.cfg.ProInfo.SalesCompleteCount
	// 回収操作回数にiniの値をセット
	c.SafeInfo.CollectCount = c.cfg.ProInfo.CollectCount
	c.logger.Info("金庫情報 起動時初期化 %+v", c.SafeInfo)
}

// 金庫情報取得
func (c *safeInfoManager) GetSafeInfo(texCon *domain.TexContext) domain.SafeInfo {
	c.OutputLogSafeInfoExCountTbl(texCon)
	return c.SafeInfo
}

// 金庫情報から入出金データをクリア
// 3:取引入金 4:取引出金 5:取引差引 6:補充入金 7:補充出金 8:補充差引 9:売上金回収  のみデータクリア
func (c *safeInfoManager) ClearCashInfo(texCon *domain.TexContext) {
	c.logger.Trace("【%v】START:入出金・売上金回収済みクリア", texCon.GetUniqueKey())
	// 入出金データのみクリアする
	for i := 0; i < len(c.SafeInfo.SortInfoTbl); i++ {
		if i < 3 || i == 10 { //0~2と10はクリアしない
			continue
		}
		c.SafeInfo.SortInfoTbl[i] = domain.SortInfoTbl{SortType: domain.SortType[i]}
	}
	// 回収操作回数をクリア
	c.clearCollectCount(texCon)
	c.OutputLogSafeInfoExCountTbl(texCon)
	c.logger.Trace("【%v】END:入出金・売上金回収済みクリア クリア後金庫情報=%v", texCon.GetUniqueKey(), c.SafeInfo)
}

// 金庫情報ログ出力
func (c *safeInfoManager) OutputLogSafeInfoExCountTbl(texCon *domain.TexContext) {
	l := fmt.Sprintf("【%v】\n", texCon.GetUniqueKey())
	l += "金庫情報出力 10000,5000,2000,1000,500,100,50,10,5,1,500予備,100予備,50予備,10予備,5予備,1予備,10000あふれ,5000あふれ,2000あふれ,1000あふれ,500あふれ,100あふれ,50あふれ,10あふれ,5あふれ,1あふれ\n"
	l += fmt.Sprintf("売上回収金額 %+v円\n", c.SafeInfo.SalesCompleteAmount)
	l += fmt.Sprintf("売上回収回数 %+v回\n", c.SafeInfo.SalesCompleteCount)
	l += fmt.Sprintf("回収操作回数 %+v回\n", c.SafeInfo.CollectCount)
	l += fmt.Sprintf("現金有高[00] %+v\n", c.addLog(c.SafeInfo.SortInfoTbl[0].ExCountTbl))
	l += fmt.Sprintf("釣銭可能[01] %+v\n", c.addLog(c.SafeInfo.SortInfoTbl[1].ExCountTbl))
	l += fmt.Sprintf("初期補充[02] %+v\n", c.addLog(c.SafeInfo.SortInfoTbl[2].ExCountTbl))
	l += fmt.Sprintf("取引入金[03] %+v\n", c.addLog(c.SafeInfo.SortInfoTbl[3].ExCountTbl))
	l += fmt.Sprintf("取引出金[04] %+v\n", c.addLog(c.SafeInfo.SortInfoTbl[4].ExCountTbl))
	l += fmt.Sprintf("取引差引[05] %+v\n", c.addLog(c.SafeInfo.SortInfoTbl[5].ExCountTbl))
	l += fmt.Sprintf("補充入金[06] %+v\n", c.addLog(c.SafeInfo.SortInfoTbl[6].ExCountTbl))
	l += fmt.Sprintf("補充出金[07] %+v\n", c.addLog(c.SafeInfo.SortInfoTbl[7].ExCountTbl))
	l += fmt.Sprintf("補充差引[08] %+v\n", c.addLog(c.SafeInfo.SortInfoTbl[8].ExCountTbl))
	l += fmt.Sprintf("売上回収[09] %+v\n", c.addLog(c.SafeInfo.SortInfoTbl[9].ExCountTbl))
	l += fmt.Sprintf("入金可能[10] %+v\n", c.addLog(c.SafeInfo.SortInfoTbl[10].ExCountTbl))
	l += fmt.Sprintf("デバ有高[12] %+v\n", c.addLog(c.deviceCashAvailable.ExCountTbl))
	l += fmt.Sprintf("論理有高[13] %+v", c.addLog(c.logicalCashAvailable.ExCountTbl))
	c.logger.Debug("%v", l)

}

func (c *safeInfoManager) addLog(i [26]int) []string {
	var s []string
	for _, v := range i {
		s = append(s, fmt.Sprintf("%3v", v))
	}
	return s
}

//============================================//
//売上金回収情報操作
//============================================//

// 売上回収情報取得
func (c *safeInfoManager) GetSalesInfo() (int, int) {
	return c.SafeInfo.SalesCompleteAmount, c.SafeInfo.SalesCompleteCount
}

// 売上金回収済金額更新
func (c *safeInfoManager) UpdateSalesCompleteAmount(texCon *domain.TexContext, amount int) {
	c.logger.Trace("【%v】START:売上金回収済金額更新 amount=%d", texCon.GetUniqueKey(), amount)
	//c.SafeInfo.SalesCompleteAmount += amount // No9068時に += に直しているが、足し算完了後の売上金回収countTblがベースなので、足すと倍々に増えると思われる
	c.SafeInfo.SalesCompleteAmount = amount //No2943
	c.logger.Trace("【%v】END:売上金回収済金額更新 更新後金庫情報=%+v", texCon.GetUniqueKey(), c.SafeInfo)
}

// 売上金回収回数をカウントアップ
func (c *safeInfoManager) CountUpSalesCompleteCount(texCon *domain.TexContext) {
	c.logger.Trace("【%v】売上金回収回数カウントアップ", texCon.GetUniqueKey())
	// カウントアップ
	c.SafeInfo.SalesCompleteCount++

	// 売上金回収回数をiniファイルに書き込む
	c.iniService.UpdateIni(texCon, "PROGRAM", "SalesCompleteCount", strconv.Itoa(c.SafeInfo.SalesCompleteCount))

}

// 売上金回収情報クリア
func (c *safeInfoManager) ClearSalesInfo(texCon *domain.TexContext) {
	c.logger.Trace("【%v】売上金回収情報クリア", texCon.GetUniqueKey())
	// データクリア
	c.SafeInfo.SalesCompleteAmount = 0
	c.SafeInfo.SalesCompleteCount = 0

	// 売上金回収回数をiniファイルに書き込む
	c.iniService.UpdateIni(texCon, "PROGRAM", "SalesCompleteCount", strconv.Itoa(c.SafeInfo.SalesCompleteCount))

}

//============================================//
// 回収操作回数（売上金回収もカウントに含まれる）
//============================================//

// 回収操作回数取得
func (c *safeInfoManager) GetCollectCount() int {
	return c.SafeInfo.CollectCount
}

// 回収操作回数をカウントアップ
func (c *safeInfoManager) CountUpCollectCount(texCon *domain.TexContext) {
	c.logger.Trace("【%v】回収操作回数カウントアップ", texCon.GetUniqueKey())
	// カウントアップ
	c.SafeInfo.CollectCount++

	// 回収操作回数をiniファイルに書き込む
	c.iniService.UpdateIni(texCon, "PROGRAM", "CollectCount", strconv.Itoa(c.SafeInfo.CollectCount))
}

// 回収操作回数クリア
func (c *safeInfoManager) clearCollectCount(texCon *domain.TexContext) {
	c.logger.Trace("【%v】回収操作回数クリア", texCon.GetUniqueKey())
	// データクリア
	c.SafeInfo.CollectCount = 0

	// 回収操作回数をiniファイルに書き込む
	c.iniService.UpdateIni(texCon, "PROGRAM", "CollectCount", strconv.Itoa(c.SafeInfo.CollectCount))
}

//============================================//
//分類別金庫情報操作
//============================================//

// GetSortInfo 分類別金庫情報取得
// sortType :0:現金有高 1:釣銭可能 2:初期補充 3:取引入金 4:取引出金 5:取引差引 6:補充入金 7:補充出金 8:補充差引 9:売上金回収 10:入金可能枚数
func (c *safeInfoManager) GetSortInfo(texCon *domain.TexContext, sortType int) (result bool, sortInfoTbl domain.SortInfoTbl) {
	defer c.logger.Trace("【%v】取得 分類別金庫情報 [%v][%+v]", texCon.GetUniqueKey(), c.GetInfoSafeSortTypeName(sortType), c.SafeInfo.SortInfoTbl[sortType])
	// 指定した分類の金庫情報を取得する
	for _, v := range c.SafeInfo.SortInfoTbl {
		if v.SortType == sortType {
			sortInfoTbl = v
			result = true
		}
	}

	if !result {
		c.logger.Trace("【%v】safeInfoManager GetSortInfo 取得失敗", texCon.GetUniqueKey())
	}

	return
}

// 分類別金庫情報更新
func (c *safeInfoManager) UpdateSortInfo(texCon *domain.TexContext, sortInfoTbl domain.SortInfoTbl) (result bool) {
	c.logger.Trace("【%v】更新 分類別金庫情報 [%v]  更新情報=%v", texCon.GetUniqueKey(), c.GetInfoSafeSortTypeName(sortInfoTbl.SortType), sortInfoTbl)

	// 指定した分類の金庫情報のみ更新する
	for i, v := range c.SafeInfo.SortInfoTbl {
		if v.SortType == sortInfoTbl.SortType {
			c.SafeInfo.SortInfoTbl[i] = sortInfoTbl
			result = true
			// 有高更新の場合には、釣り銭可能枚数も一緒に更新する
			if v.SortType == 0 {
				c.SafeInfo.SortInfoTbl[domain.CHANGE_AVAILABLE] = c.updateChangeAvailable(sortInfoTbl)
			}
			break
		}
	}
	if !result {
		c.logger.Error("【%v】- 金庫分類情報種別不一致", texCon.GetUniqueKey())
	}
	c.OutputLogSafeInfoExCountTbl(texCon)
	return
}

// 差引情報更新
func (c *safeInfoManager) UpdateBalanceInfo(texCon *domain.TexContext) {
	// 金庫情報ログ出力
	c.OutputLogSafeInfoExCountTbl(texCon)

	update := func(balance *domain.SortInfoTbl, deposit domain.SortInfoTbl, withdrawal domain.SortInfoTbl) {
		// 「入金 - 出金」で差引情報を更新
		balance.Amount = deposit.Amount - withdrawal.Amount
		for i := range balance.CountTbl {
			balance.CountTbl[i] = deposit.CountTbl[i] - withdrawal.CountTbl[i]
		}
		for i := range balance.ExCountTbl {
			balance.ExCountTbl[i] = deposit.ExCountTbl[i] - withdrawal.ExCountTbl[i]
		}
	}

	// 取引差引更新
	update(&c.SafeInfo.SortInfoTbl[domain.TRANSACTION_BALANCE], c.SafeInfo.SortInfoTbl[domain.TRANSACTION_DEPOSIT], c.SafeInfo.SortInfoTbl[domain.TRANSACTION_WITHDRAWAL])
	// 補充差引更新
	update(&c.SafeInfo.SortInfoTbl[domain.REPLENISHMENT_BALANCE], c.SafeInfo.SortInfoTbl[domain.REPLENISHMENT_DEPOSIT], c.SafeInfo.SortInfoTbl[domain.REPLENISHMENT_WITHDRAWAL])
	c.logger.Debug("【%v】取引・補充差引更新 after 取引差引=%v, 補充差引=%v", texCon.GetUniqueKey(), c.SafeInfo.SortInfoTbl[domain.TRANSACTION_BALANCE], c.SafeInfo.SortInfoTbl[domain.REPLENISHMENT_BALANCE])
}

// UpdateSortInfoCumulative 分類情報更新（累計）
// 現在のデータに加算する
// 論理有高の更新も行う
func (c *safeInfoManager) UpdateSortInfoCumulative(texCon *domain.TexContext, sortType int, amount int, countTbl [10]int, exCountTbl [26]int) {

	c.updateSortInfoCumulative(texCon, sortType, amount, countTbl, exCountTbl)

	// 入金の場合の論理有高更新
	if sortType == 3 || sortType == 6 {
		// 論理有高更新
		c.UpdateInLogicalCashAvailable(texCon, domain.SortInfoTbl{
			SortType:   domain.CASH_AVAILABLE,
			Amount:     amount,
			CountTbl:   countTbl,
			ExCountTbl: exCountTbl,
		})
	}

}

func (c *safeInfoManager) updateSortInfoCumulative(texCon *domain.TexContext, sortType int, amount int, countTbl [10]int, exCountTbl [26]int) {
	// 指定キーが存在するかチェック
	ok := c.checkKey(sortType)
	if !ok {
		c.logger.Error("【%v】更新 金庫分類情報 金庫分類情報種別不一致", texCon.GetUniqueKey())
		return
	}

	updateData := &c.SafeInfo.SortInfoTbl[sortType]
	c.logger.Trace("【%v】START:更新 金庫分類情報(累計)[%v] 更新前 =%+v", texCon.GetUniqueKey(), c.GetInfoSafeSortTypeName(sortType), c.SafeInfo.SortInfoTbl[sortType])

	updateData.Amount += amount

	for i, v := range countTbl {
		updateData.CountTbl[i] += v
	}

	for i, v := range exCountTbl {
		updateData.ExCountTbl[i] += v
	}

	c.logger.Trace("【%v】END:更新 金庫分類情報(累計)[%v] 更新後 =%+v", texCon.GetUniqueKey(), c.GetInfoSafeSortTypeName(sortType), c.SafeInfo.SortInfoTbl[sortType])
}

// UpdateSortInfoCumulativeNoUpdateLogicalCash 分類情報更新（累計）
// 現在のデータに加算する
// 論理有高の更新は実施しない
func (c *safeInfoManager) UpdateSortInfoCumulativeNoUpdateLogicalCash(texCon *domain.TexContext, sortType int, amount int, countTbl [10]int, exCountTbl [26]int) {
	c.updateSortInfoCumulative(texCon, sortType, amount, countTbl, exCountTbl)

}

// 指定キーが存在するかチェック
func (c *safeInfoManager) checkKey(key int) bool {
	for i := range c.SafeInfo.SortInfoTbl {
		if i == key {
			return true
		}
	}
	return false
}

func (c *safeInfoManager) updateChangeAvailable(s domain.SortInfoTbl) domain.SortInfoTbl {
	exC := s.ExCountTbl // 26配列を取得
	cT := s.CountTbl    // 10配列を取得
	oT := exC[16:]      // あふれ金庫分の情報を取得

	NewExCTbl := [domain.EXTRA_CASH_TYPE_SHITEI]int{exC[0], exC[1], exC[2], exC[3], exC[4], exC[5], exC[6], exC[7], exC[8], exC[9], exC[10], exC[11], exC[12], exC[13], exC[14], exC[15]} // 釣り銭可26能配列

	NewCTbl := [domain.CASH_TYPE_SHITEI]int{} // 有高-あふれ=釣り銭可能金種
	for i := 0; i < 10; i++ {
		NewCTbl[i] = c.subtract(cT[i], oT[i])
	}

	amount := calculation.NewCassette(NewExCTbl).GetTotalAmount()

	return domain.SortInfoTbl{
		Amount:     amount,
		SortType:   domain.CHANGE_AVAILABLE,
		CountTbl:   NewCTbl,
		ExCountTbl: NewExCTbl,
	}
}

func (c *safeInfoManager) subtract(i int, j int) int {
	x := i - j
	if x < 0 {
		return 0
	}
	return x
}

//============================================//
//処理前補充情報操作
//============================================//

// 処理前補充差引取得
func (c *safeInfoManager) GetBeforeReplenishmentBalance() domain.SortInfoTbl {
	return c.BeforeReplenishmentBalance
}

// 処理前補充差引更新
func (c *safeInfoManager) UpdateBeforeReplenishmentBalance(texCon *domain.TexContext) {
	c.logger.Trace("【%v】START:更新 処理前補充差引", texCon.GetUniqueKey())

	// データ初期化
	c.BeforeReplenishmentBalance = domain.SortInfoTbl{}

	// 金庫情報から現在の補充差引のデータをコピー
	c.BeforeReplenishmentBalance = c.SafeInfo.SortInfoTbl[domain.REPLENISHMENT_BALANCE]

	c.logger.Trace("【%v】END:更新 処理前補充差引 %v", texCon.GetUniqueKey(), c.BeforeReplenishmentBalance)
}

// ============================================//
// デバイス有高
// ============================================//

// GetDeviceCashAvailable デバイス有高を取得
func (c *safeInfoManager) GetDeviceCashAvailable(texCon *domain.TexContext) domain.SortInfoTbl {
	c.logger.Trace("【%v】取得 デバイス有高", texCon.GetUniqueKey())
	c.OutputLogSafeInfoExCountTbl(texCon)
	return c.deviceCashAvailable
}

// UpdateDeviceCashAvailable 受信した有高通知を登録する
func (c *safeInfoManager) UpdateDeviceCashAvailable(texCon *domain.TexContext, tbl domain.SortInfoTbl) {
	c.logger.Trace("【%v】更新 デバイス有高 %v", texCon.GetUniqueKey(), tbl)
	tbl.SortType = domain.DEVICE_AVAILABLE
	c.deviceCashAvailable = tbl

	c.OutputLogSafeInfoExCountTbl(texCon)

}

// ============================================//
// 論理有高
// ============================================//

// GetLogicalCashAvailable 論理有高を取得
func (c *safeInfoManager) GetLogicalCashAvailable(texCon *domain.TexContext) domain.SortInfoTbl {
	c.logger.Trace("【%v】取得 論理有高", texCon.GetUniqueKey())
	c.OutputLogSafeInfoExCountTbl(texCon)
	return c.logicalCashAvailable
}

// UpdateInLogicalCashAvailable 受信した入金情報を論理有高に保存する
func (c *safeInfoManager) UpdateInLogicalCashAvailable(texCon *domain.TexContext, tbl domain.SortInfoTbl) {
	c.logger.Trace("【%v】更新 論理有高入金 %v", texCon.GetUniqueKey(), tbl)

	c.logicalCashAvailable.Amount += tbl.Amount
	c.logicalCashAvailable.SortType = domain.CASH_AVAILABLE

	for i := range c.logicalCashAvailable.CountTbl {
		c.logicalCashAvailable.CountTbl[i] += tbl.CountTbl[i]
	}
	for i := range c.logicalCashAvailable.ExCountTbl {
		c.logicalCashAvailable.ExCountTbl[i] += tbl.ExCountTbl[i]
	}

	if LogicalChange {
		cashAvailable := domain.SortInfoTbl{
			SortType:   domain.CASH_AVAILABLE,
			Amount:     c.logicalCashAvailable.Amount,
			CountTbl:   c.logicalCashAvailable.CountTbl,
			ExCountTbl: c.logicalCashAvailable.ExCountTbl,
		}
		c.UpdateSortInfo(texCon, cashAvailable)
	}

	c.OutputLogSafeInfoExCountTbl(texCon)
}

// UpdateOutLogicalCashAvailable 受信した出金要求から、論理有高を計算する
func (c *safeInfoManager) UpdateOutLogicalCashAvailable(texCon *domain.TexContext, tbl domain.SortInfoTbl) {
	c.logger.Trace("【%v】更新 論理有高出金 %v", texCon.GetUniqueKey(), tbl)
	c.logicalCashAvailable.Amount -= tbl.Amount
	c.logicalCashAvailable.SortType = domain.CASH_AVAILABLE

	for i := range c.logicalCashAvailable.CountTbl {
		c.logicalCashAvailable.CountTbl[i] -= tbl.CountTbl[i]
	}
	for i := range c.logicalCashAvailable.ExCountTbl {
		c.logicalCashAvailable.ExCountTbl[i] -= tbl.ExCountTbl[i]
	}

	if LogicalChange {
		cashAvailable := domain.SortInfoTbl{
			SortType:   domain.CASH_AVAILABLE,
			Amount:     c.logicalCashAvailable.Amount,
			CountTbl:   c.logicalCashAvailable.CountTbl,
			ExCountTbl: c.logicalCashAvailable.ExCountTbl,
		}
		c.UpdateSortInfo(texCon, cashAvailable)
	}

	c.OutputLogSafeInfoExCountTbl(texCon)
}

// UpdateAllLogicalCashAvailable 論理有高を一括更新する
func (c *safeInfoManager) UpdateAllLogicalCashAvailable(texCon *domain.TexContext, tbl domain.SortInfoTbl) {
	c.logger.Trace("【%v】更新 論理有高一括 %v", texCon.GetUniqueKey(), tbl)
	c.logicalCashAvailable = tbl

	if LogicalChange {
		cashAvailable := domain.SortInfoTbl{
			SortType:   domain.CASH_AVAILABLE,
			Amount:     c.logicalCashAvailable.Amount,
			CountTbl:   c.logicalCashAvailable.CountTbl,
			ExCountTbl: c.logicalCashAvailable.ExCountTbl,
		}
		c.UpdateSortInfo(texCon, cashAvailable)
	}

	c.OutputLogSafeInfoExCountTbl(texCon)
}

// ============================================//
// デバイス有高‐論理有高
// ============================================//

// LogicalChange 0番有高に論理有高をセットするかどうか
// true:セットする false:セットしない
// 論理有高計算ロジックの印として一旦残しておく
// fit-bリリース次期からは、trueでの動作を対象として開発している為
// false時の処理は既に保証されない。
var LogicalChange = true

func (c *safeInfoManager) GetAvailableBalance(texCon *domain.TexContext) (bool, domain.SortInfoTbl) {
	c.logger.Trace("【%v】取得 デバイス有高‐論理有高差分", texCon.GetUniqueKey())

	defer c.OutputLogSafeInfoExCountTbl(texCon)

	var mismatch bool

	// 合計金額比較を実施
	dev := calculation.NewCassette(c.deviceCashAvailable.ExCountTbl).GetTotalAmount()
	logi := calculation.NewCassette(c.logicalCashAvailable.ExCountTbl).GetTotalAmount()

	// 合計値と同じ場合にはデバイスを正として更新
	// デバイスによっては金種を制御できないので,合計値で先にチェックする
	if dev == logi {
		mismatch = false

		// 差分無しとして,論理有高を更新
		c.UpdateAllLogicalCashAvailable(texCon, c.deviceCashAvailable)

		// 差分無しとしてresult生成
		result := domain.SortInfoTbl{
			SortType:   domain.AVAILABLE_BALANCE,
			Amount:     0,
			CountTbl:   [10]int{},
			ExCountTbl: [26]int{},
		}
		c.logger.Debug("【%v】- デバイス有高を上書 mismatch=%t result=%t", texCon.GetUniqueKey(), mismatch, result)
		return mismatch, result
	}

	// 合計が異なる場合
	mismatch = true
	// デバイス有高‐論理枚数
	availableBalanceExCountTbl := calculation.NewCassette(c.deviceCashAvailable.ExCountTbl).Subtract(c.logicalCashAvailable.ExCountTbl)

	availableBalance := calculation.NewCassette(availableBalanceExCountTbl)

	// 不一致枚数の合計金額
	availableBalanceAmount := availableBalance.GetTotalAmount()
	// 不一致枚数の 10金種配列
	availableBalanceCountTbl := availableBalance.ExCountTblToTenCountTbl()

	result := domain.SortInfoTbl{
		SortType:   domain.AVAILABLE_BALANCE,
		Amount:     availableBalanceAmount,
		CountTbl:   availableBalanceCountTbl,
		ExCountTbl: availableBalanceExCountTbl,
	}
	c.logger.Debug("【%v】- mismatch =%t result =%v", texCon.GetUniqueKey(), mismatch, result)

	return mismatch, result

}
