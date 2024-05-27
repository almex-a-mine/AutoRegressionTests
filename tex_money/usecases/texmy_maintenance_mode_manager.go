package usecases

import (
	"strconv"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
)

type MaintenanceModeManager interface {
	GetMode(texCon *domain.TexContext) int
	GetMaintenanceMode(texCon *domain.TexContext) int
	SetReceiveData(texCon *domain.TexContext, req domain.RequestMaintenanceMode)
	SetStatusEnd(texCon *domain.TexContext) bool
	SetStatusStart(texCon *domain.TexContext) bool
}

type maintenanceModeManager struct {
	logger          handler.LoggerRepository
	config          config.Configuration
	aggregateMng    AggregateManager
	safeInfoMng     SafeInfoManager
	iniService      IniServiceRepository
	mode            int
	maintenanceMode int
}

func NewMaintenanceModeManager(logger handler.LoggerRepository,
	config config.Configuration,
	aggregateMng AggregateManager,
	safeInfoMng SafeInfoManager,
	iniService IniServiceRepository,
) MaintenanceModeManager {
	return &maintenanceModeManager{
		logger:          logger,
		config:          config,
		aggregateMng:    aggregateMng,
		safeInfoMng:     safeInfoMng,
		iniService:      iniService,
		maintenanceMode: config.ProInfo.MaintenanceModeStatus, //iniの値を入れる
	}
}

// 保守モード設定
func (c *maintenanceModeManager) SetReceiveData(texCon *domain.TexContext, req domain.RequestMaintenanceMode) {
	c.mode = req.Mode
	c.setRecvStatus(texCon, req.Action) // 受信値の保守業務モードステータスをセット
}

// 受信値から保守業務モードステータスをセット
// 1：補充開始 2：補充終了 3：締め開始 4：締め終了
func (c *maintenanceModeManager) setRecvStatus(texCon *domain.TexContext, action bool) {
	switch {
	case c.mode == domain.REPLENISHMENT_MODE && action: //補充開始
		c.updateMaintenanceMode(texCon, domain.REPLENISHMENT_START)
	case c.mode == domain.REPLENISHMENT_MODE && !action: //補充終了
		c.updateMaintenanceMode(texCon, domain.REPLENISHMENT_END)
	case c.mode == domain.CLOSING_MODE && action: //締め開始
		c.updateMaintenanceMode(texCon, domain.CLOSING_START)
	case c.mode == domain.CLOSING_MODE && !action: //締め終了
		c.updateMaintenanceMode(texCon, domain.CLOSING_END)
	case c.mode == 0 && action: //締め日計表
		c.updateMaintenanceMode(texCon, domain.CLOSING_START)
	default:
	}
}

// 保守モード取得
func (c *maintenanceModeManager) GetMode(texCon *domain.TexContext) int {
	c.logger.Trace("【%v】GetMode(mode %v)", texCon.GetUniqueKey(), c.mode)
	return c.mode
}

// 保守モード取得
func (c *maintenanceModeManager) updateMaintenanceMode(texCon *domain.TexContext, maintenanceMode int) {
	c.logger.Trace("【%v】updateMaintenanceMode(maintenanceMode %v)", texCon.GetUniqueKey(), maintenanceMode)
	c.maintenanceMode = maintenanceMode
	// iniへの書き込み
	c.iniService.UpdateIni(texCon, "PROGRAM", "maintenanceModeStatus", strconv.Itoa(c.maintenanceMode))
}

// 保守モード取得
func (c *maintenanceModeManager) GetMaintenanceMode(texCon *domain.TexContext) int {
	c.logger.Trace("【%v】GetMaintenanceMode(maintenanceMode %v)", texCon.GetUniqueKey(), c.maintenanceMode)
	return c.maintenanceMode
}

// 動作要求：開始 の処理
func (c *maintenanceModeManager) SetStatusStart(texCon *domain.TexContext) bool {

	// 保守業務モード=0の処理
	if c.mode == 0 {
		return c.modeZeroStart(texCon)
	}

	// 前回処理が正常に終了している場合はレポート用金庫情報をクリア
	switch {
	case c.mode == domain.REPLENISHMENT_MODE && c.maintenanceMode != domain.REPLENISHMENT_START: // 補充
		if !c.aggregateMng.ClearAggregateSafeInfo(texCon, c.mode) { //補充のレポート用金庫情報クリア
			return false
		}
	case c.mode == domain.CLOSING_MODE && c.maintenanceMode != domain.CLOSING_START: // 締め
		c.aggregateMng.ClearAggregateData(texCon) // レポート用金庫情報クリア
	default:
	}

	// 処理前の補充差引情報を保存
	c.safeInfoMng.UpdateBeforeReplenishmentBalance(texCon)

	// レポート用有高金種配列を更新
	if !c.aggregateMng.UpdateBeforeCountTbl(texCon, c.mode, domain.BEFORE_AMOUNT_COUNT_TBL) {
		return false
	}
	if !c.aggregateMng.UpdateBeforeCountTbl(texCon, c.mode, domain.BEFORE_REPLENISH_COUNT_TBL) {
		return false
	}
	if !c.aggregateMng.UpdateBeforeCountTbl(texCon, c.mode, domain.BEFORE_COLLECT_COUNT_TBL) {
		return false
	}

	return true
}

// 保守業務モード=0の処理
// 保守業務モード=0のときはその時点の金庫情報を帳票用データとして登録する
func (c *maintenanceModeManager) modeZeroStart(texCon *domain.TexContext) bool {

	c.aggregateMng.ClearAggregateData(texCon) // レポート用金庫情報クリア

	// 現在の金庫情報取得
	_, sortInfo := c.safeInfoMng.GetSortInfo(texCon, domain.CASH_AVAILABLE)

	// レポート用有高金種配列を現在有高枚数で更新
	// 補充_処理前
	if !c.aggregateMng.UpdateBeforeCountTbl(texCon, domain.REPLENISHMENT_MODE, domain.BEFORE_AMOUNT_COUNT_TBL) {
		return false
	}
	// 補充_処理後
	if !c.aggregateMng.UpdateAggregateCountTbl(texCon, domain.REPLENISHMENT_MODE, domain.AFTER_AMOUNT_COUNT_TBL, sortInfo.ExCountTbl) {
		return false
	}
	// 締め_処理前
	if !c.aggregateMng.UpdateBeforeCountTbl(texCon, domain.CLOSING_MODE, domain.BEFORE_AMOUNT_COUNT_TBL) {
		return false
	}
	// 締め_処理後
	if !c.aggregateMng.UpdateAggregateCountTbl(texCon, domain.CLOSING_MODE, domain.AFTER_AMOUNT_COUNT_TBL, sortInfo.ExCountTbl) {
		return false
	}

	// 処理前の補充差引情報を保存
	c.safeInfoMng.UpdateBeforeReplenishmentBalance(texCon)

	// 保守業務モードステータスを更新
	c.updateMaintenanceMode(texCon, domain.CLOSING_END)
	return true
}

// 動作要求：終了 の処理
func (c *maintenanceModeManager) SetStatusEnd(texCon *domain.TexContext) bool {

	// 金庫情報をクリア
	switch c.mode {
	case domain.REPLENISHMENT_MODE:
		if result := c.aggregateMng.ClearAggregateSafeInfo(texCon, c.mode); !result { //補充のレポート用金庫情報クリア
			return false
		}
	case domain.CLOSING_MODE:
		c.aggregateMng.ClearAggregateData(texCon) // レポート用金庫情報クリア
	default:
		break
	}

	return true
}
