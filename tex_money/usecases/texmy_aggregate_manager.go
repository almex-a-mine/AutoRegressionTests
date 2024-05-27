package usecases

import (
	"encoding/json"
	"fmt"
	"sync"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
)

type aggregateManager struct {
	logger        handler.LoggerRepository
	cfg           config.Configuration
	syslogMng     SyslogManager
	errorMng      ErrorManager
	safeInfoMng   SafeInfoManager
	iniService    IniServiceRepository
	mu            sync.Mutex
	aggregateData domain.AggregateData
}

// レポート用金庫情報
func NewAggregateManager(logger handler.LoggerRepository, cfg config.Configuration, syslogMng SyslogManager, errorMng ErrorManager, safeInfoMng SafeInfoManager, iniService IniServiceRepository) AggregateManager {
	c := &aggregateManager{
		logger:        logger,
		cfg:           cfg,
		syslogMng:     syslogMng,
		errorMng:      errorMng,
		safeInfoMng:   safeInfoMng,
		iniService:    iniService,
		aggregateData: domain.AggregateData{}}
	c.init() //初期化
	return c
}

// レポート用金庫情報初期化
func (c *aggregateManager) init() {
	// iniの値をセット
	for k, v := range c.cfg.AggregateData {
		c.aggregateData[k] = v
	}
	c.logger.Info("集計管理 初期化")
}

// レポート用金庫情報の金種配列を指定枚数で更新 mode:業務モード　tbl:格納したい配列　countTbl
func (c *aggregateManager) UpdateAggregateCountTbl(
	texCon *domain.TexContext,
	mode int,
	tbl int,
	countTbl [domain.EXTRA_CASH_TYPE_SHITEI]int) (result bool) {
	c.logger.Trace("【%v】START:更新 集計管理  mode=%d指定, tbl=%d, countTbl=%v", texCon.GetUniqueKey(), mode, tbl, countTbl)

	// 保守業務モードのデータが存在するかチェック
	if result = c.checkDataExist(texCon, mode); !result {
		return
	}

	// 指定した金種配列を更新
	c.mu.Lock()
	switch tbl {
	case domain.REPLENISH_COUNT_TBL:
		c.aggregateData[mode].ReplenishCountTbl = countTbl
	case domain.COLLECT_COUNT_TBL:
		c.aggregateData[mode].CollectCountTbl = countTbl
	case domain.AFTER_AMOUNT_COUNT_TBL:
		c.aggregateData[mode].AfterAmountCountTbl = countTbl
	case domain.SALES_COLLECT_TBL:
		c.aggregateData[mode].SalesCollectCountTbl = countTbl
	}
	c.mu.Unlock()

	// iniファイルに書き込み
	c.aggregateDataIni(texCon)

	c.logger.Trace("【%v】END:更新 集計管理 mode=%d指定 %+v", texCon.GetUniqueKey(), mode, *c.aggregateData[mode])
	return
}

// レポート用金庫情報の金種配列に現在の金庫情報をコピー
// mode:業務モード tbl:格納したい配列
func (c *aggregateManager) UpdateBeforeCountTbl(texCon *domain.TexContext, mode int, tbl int) (result bool) {
	c.logger.Trace("【%v】START:更新 集計管理 前回データコピー", texCon.GetUniqueKey())

	// 保守業務モードのデータが存在するかチェック
	if result = c.checkDataExist(texCon, mode); !result {
		return
	}

	// レポート用金種配列にコピー
	c.mu.Lock()
	var sortInfo domain.SortInfoTbl
	switch tbl {
	case domain.BEFORE_AMOUNT_COUNT_TBL:
		result, sortInfo = c.safeInfoMng.GetSortInfo(texCon, domain.CASH_AVAILABLE)
		c.aggregateData[mode].BeforeAmountCountTbl = sortInfo.ExCountTbl
	case domain.BEFORE_REPLENISH_COUNT_TBL:
		result, sortInfo = c.safeInfoMng.GetSortInfo(texCon, domain.REPLENISHMENT_DEPOSIT)
		c.aggregateData[mode].BeforeReplenishCountTbl = sortInfo.ExCountTbl
	case domain.BEFORE_COLLECT_COUNT_TBL:
		result, sortInfo = c.safeInfoMng.GetSortInfo(texCon, domain.REPLENISHMENT_WITHDRAWAL)
		c.aggregateData[mode].BeforeCollectCountTbl = sortInfo.ExCountTbl
	}
	c.mu.Unlock()

	// iniファイルに書き込み
	c.aggregateDataIni(texCon)

	c.logger.Trace("【%v】END:更新 集計管理 前回データコピー mode=%d %+v", texCon.GetUniqueKey(), mode, c.aggregateData[mode])
	return
}

// レポート用金庫情報クリア
func (c *aggregateManager) ClearAggregateData(texCon *domain.TexContext) {
	c.logger.Trace("【%v】START:クリア 集計管理", texCon.GetUniqueKey())

	//初期化
	c.mu.Lock()
	c.aggregateData = domain.AggregateData{}
	//デフォルト値をセット
	for _, v := range domain.MODE {
		c.aggregateData[v] = &domain.AggregateSafeInfo{}
	}
	c.mu.Unlock()

	// iniファイルに書き込み
	c.aggregateDataIni(texCon)

	c.logger.Trace("【%v】END:クリア 集計管理 %+v", texCon.GetUniqueKey(), c.aggregateData)
}

// 指定した保守業務モードのレポート用金庫情報クリア
func (c *aggregateManager) ClearAggregateSafeInfo(texCon *domain.TexContext, mode int) (result bool) {
	c.logger.Trace("【%v】START:クリア 集計管理 mode=%v指定", texCon.GetUniqueKey(), mode)

	// 保守業務モードのデータが存在するかチェック
	if result = c.checkDataExist(texCon, mode); !result {
		return
	}

	//クリア
	c.mu.Lock()
	c.aggregateData[mode] = &domain.AggregateSafeInfo{}
	c.mu.Unlock()

	// iniファイルに書き込み
	c.aggregateDataIni(texCon)

	c.logger.Trace("【%v】END:クリア 集計管理 mode=%v指定 %+v", texCon.GetUniqueKey(), mode, *c.aggregateData[mode])
	return
}

// レポート用金庫情報取得
func (c *aggregateManager) GetAggregateSafeInfo(texCon *domain.TexContext, mode int) (result bool, aggregateSafeInfo domain.AggregateSafeInfo) {
	c.logger.Trace("【%v】START:取得 集計管理 mode=%d指定", texCon.GetUniqueKey(), mode)
	defer c.logger.Trace("【%v】END:取得 集計管理 mode=%d指定 result=%t, %+v", texCon.GetUniqueKey(), mode, result, aggregateSafeInfo)

	// 保守業務モードのデータが存在するかチェック
	if result = c.checkDataExist(texCon, mode); !result {
		return
	}

	aggregateSafeInfo = *c.aggregateData[mode]
	return
}

// レポート用金庫情報取得
func (c *aggregateManager) GetAggregateSafeInfoAll() domain.AggregateData {
	return c.aggregateData

}

// 保守業務モードのデータが存在するかチェック
func (c *aggregateManager) checkDataExist(texCon *domain.TexContext, mode int) bool {
	if _, ok := c.aggregateData[mode]; !ok {
		switch mode {
		case domain.REPLENISHMENT_MODE, domain.CLOSING_MODE:
			// modeは正しいが、データが存在しない場合は要素を追加
			c.mu.Lock()
			c.aggregateData[mode] = &domain.AggregateSafeInfo{}
			c.mu.Unlock()
		default:
			c.logger.Error("【%v】aggregateManager checkDataExist 動作要求情報不正 mode=%d", texCon.GetUniqueKey(), mode)
			return false
		}
	}
	return true
}

// iniへの書き込み
func (c *aggregateManager) aggregateDataIni(texCon *domain.TexContext) {

	// 最外部のmapを初期化
	sectionKeyValue := make(map[string]map[string]string)

	for k, v := range c.aggregateData {
		// セクション名作成
		section := fmt.Sprintf("AggregateData_%d", k)
		// セクション毎に初期化
		if sectionKeyValue[section] == nil {
			sectionKeyValue[section] = make(map[string]string)
		}

		// json文字列に変換して書き込み TODO:For文へ変換可能
		// 有高金種配列
		bytes1, err := json.Marshal(v.BeforeAmountCountTbl)
		if err != nil {
			c.logger.Error("【%v】aggregateManager aggregateDataIni err=%v", texCon.GetUniqueKey(), err)
			return
		}
		sectionKeyValue[section]["BeforeAmountCountTbl"] = string(bytes1)

		//処理前補充入金金種配列
		bytes2, err := json.Marshal(v.BeforeReplenishCountTbl)
		if err != nil {
			c.logger.Error("【%v】aggregateManager aggregateDataIni err=%v", texCon.GetUniqueKey(), err)
			return
		}
		sectionKeyValue[section]["BeforeReplenishCountTbl"] = string(bytes2)

		//補充金種配列
		bytes3, err := json.Marshal(v.ReplenishCountTbl)
		if err != nil {
			c.logger.Error("【%v】aggregateManager aggregateDataIni err=%v", texCon.GetUniqueKey(), err)
			return
		}
		sectionKeyValue[section]["RepelenishCountTbl"] = string(bytes3)

		//処理前回収金種配列
		bytes4, err := json.Marshal(v.BeforeCollectCountTbl)
		if err != nil {
			c.logger.Error("【%v】aggregateManager aggregateDataIni err=%v", texCon.GetUniqueKey(), err)
			return
		}
		sectionKeyValue[section]["BeforeCollectCountTbl"] = string(bytes4)

		//回収金種配列
		bytes5, err := json.Marshal(v.CollectCountTbl)
		if err != nil {
			c.logger.Error("【%v】aggregateManager aggregateDataIni err=%v", texCon.GetUniqueKey(), err)
			return
		}
		sectionKeyValue[section]["CollectCountTbl"] = string(bytes5)

		//処理後有高金種配列
		bytes6, err := json.Marshal(v.AfterAmountCountTbl)
		if err != nil {
			c.logger.Error("【%v】aggregateManager aggregateDataIni err=%v", texCon.GetUniqueKey(), err)
			return
		}
		sectionKeyValue[section]["AfterAmountCountTbl"] = string(bytes6)

		//売上回収金種配列
		bytes7, err := json.Marshal(v.SalesCollectCountTbl)
		if err != nil {
			c.logger.Error("【%v】aggregateManager aggregateDataIni err=%v", texCon.GetUniqueKey(), err)
			return
		}
		sectionKeyValue[section]["SalesCollectCountTbl"] = string(bytes7)

	}
	c.iniService.MultipleUpdateIni(texCon, sectionKeyValue)
}

// 指定した配列の差分算出
func (c *aggregateManager) DiffTbl(texCon *domain.TexContext, beforeTbl [domain.EXTRA_CASH_TYPE_SHITEI]int, afterTbl [domain.EXTRA_CASH_TYPE_SHITEI]int) (countTbl [domain.EXTRA_CASH_TYPE_SHITEI]int) {
	c.logger.Trace("【%v】START:差分 集計管理 beforeTbl=%v afterTbl=%v", texCon.GetUniqueKey(), beforeTbl, afterTbl)

	// 差分枚数の算出
	for i := range beforeTbl {
		countTbl[i] = afterTbl[i] - beforeTbl[i]
	}

	c.logger.Trace("【%v】END:差分 集計管理 countTbl=%v", texCon.GetUniqueKey(), countTbl)
	return
}

func (c *aggregateManager) OutputLogAggregateExCountTbl() {
	c.logger.Debug("レポート用保持情報出力")
	c.logger.Debug("補充処理----")
	c.logger.Debug("[補充]処理前有高金種 %+v", c.aggregateData[1].BeforeAmountCountTbl)
	c.logger.Debug("[補充]処理後有高金種 %+v", c.aggregateData[1].AfterAmountCountTbl)
	c.logger.Debug("[補充]処理前補充入金 %+v", c.aggregateData[1].BeforeReplenishCountTbl)
	c.logger.Debug("[補充]処理前回収金種 %+v", c.aggregateData[1].BeforeCollectCountTbl)
	c.logger.Debug("[補充]処理中補充金種 %+v", c.aggregateData[1].ReplenishCountTbl)
	c.logger.Debug("[補充]処理中回収金種 %+v", c.aggregateData[1].CollectCountTbl)
	c.logger.Debug("[補充]売上金回収金種 %+v", c.aggregateData[1].SalesCollectCountTbl)
	c.logger.Debug("締め処理----")
	c.logger.Debug("[締め]処理前有高金種 %+v", c.aggregateData[100].BeforeAmountCountTbl)
	c.logger.Debug("[締め]処理後有高金種 %+v", c.aggregateData[100].AfterAmountCountTbl)
	c.logger.Debug("[締め]処理前補充入金 %+v", c.aggregateData[100].BeforeReplenishCountTbl)
	c.logger.Debug("[締め]処理前回収金種 %+v", c.aggregateData[100].BeforeCollectCountTbl)
	c.logger.Debug("[締め]処理中補充金種 %+v", c.aggregateData[100].ReplenishCountTbl)
	c.logger.Debug("[締め]処理中回収金種 %+v", c.aggregateData[100].CollectCountTbl)
	c.logger.Debug("[締め]売上金回収金種 %+v", c.aggregateData[100].SalesCollectCountTbl)
}
