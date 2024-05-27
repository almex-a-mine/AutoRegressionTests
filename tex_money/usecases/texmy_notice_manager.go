package usecases

import (
	"sync"
	"tex_money/domain"
	"tex_money/domain/handler"
)

type TexMoneyNoticeManager struct {
	statusInData          domain.StatusIndata
	statusOutData         domain.StatusOutdata
	statusOutDataBillBox  domain.StatusOutdataBillBox
	statusCollectData     domain.StatusCollectData
	statusAmountData      domain.StatusAmount
	statusCash            domain.StatusCash
	statusExchange        domain.StatusExchange
	statusReport          domain.StatusReport
	statusService         domain.StatusService
	statusSystemOperation domain.StatusSystemData
	logger                handler.LoggerRepository
	mtxStatus             sync.Mutex
	isDiff                bool // 回収データ通知の差分
	isDiffIndata          bool // 入金データ通知の差分
	isDiffOutdata         bool // 出金データ通知の差分
	isDiffAmountdata      bool // 有高データ通知の差分
}

type TexMoneyNoticeManagerRepository interface {
	GetStatusInData(texCon *domain.TexContext) domain.StatusIndata
	UpdateStatusInData(texCon *domain.TexContext, statusInData domain.StatusIndata) bool
	DiffCheckStatusInData(texCon *domain.TexContext) bool

	GetStatusOutData(texCon *domain.TexContext) domain.StatusOutdata
	UpdateStatusOutData(texCon *domain.TexContext, statusOutData domain.StatusOutdata) bool
	DiffCheckStatusOutdata(texCon *domain.TexContext) bool

	GetStatusOutDataBillBox(texCon *domain.TexContext) domain.StatusOutdataBillBox
	UpdateStatusBillBox(texCon *domain.TexContext, statusOutData domain.OutStatus)

	GetStatusCollectData(texCon *domain.TexContext) domain.StatusCollectData
	UpdateStatusCollectData(texCon *domain.TexContext, statusCollect domain.StatusCollectData) bool
	DiffCheckStatusCollectData(texCon *domain.TexContext) bool //回収データ通知 差分有り無しを外部インターフェース渡す

	GetStatusAmountData(texCon *domain.TexContext) domain.StatusAmount
	UpdateStatusAmountData(texCon *domain.TexContext, statusAmountData domain.StatusAmount) bool
	DiffCheckStatusAmountData(texCon *domain.TexContext) bool

	GetStatusCashData(texCon *domain.TexContext) domain.StatusCash
	UpdateStatusCashData(texCon *domain.TexContext, statusCash domain.StatusCash) bool

	GetStatusExchangeData(texCon *domain.TexContext) domain.StatusExchange
	UpdateStatusExchangeData(texCon *domain.TexContext, statusExchange domain.StatusExchange) bool

	GetStatusReportData(texCon *domain.TexContext) domain.StatusReport
	UpdateStatusReportData(texCon *domain.TexContext, statusReport domain.StatusReport) bool

	GetStatusServiceData(texCon *domain.TexContext) domain.StatusService
	UpdateStatusServiceData(texCon *domain.TexContext, statusService domain.StatusService) bool

	GetStatusSystemOperationData(texCon *domain.TexContext) domain.StatusSystemData
	UpdateStatusSystemOperationData(texCon *domain.TexContext, statusSystemData domain.StatusSystemData) bool
}

func NewTexMoneyNoticeManager(l handler.LoggerRepository) TexMoneyNoticeManagerRepository {
	return &TexMoneyNoticeManager{
		statusInData:          domain.StatusIndata{},
		statusOutData:         domain.StatusOutdata{},
		statusOutDataBillBox:  domain.StatusOutdataBillBox{},
		statusCollectData:     domain.StatusCollectData{},
		statusAmountData:      domain.StatusAmount{},
		statusCash:            domain.StatusCash{},
		statusExchange:        domain.StatusExchange{},
		statusReport:          domain.StatusReport{},
		statusService:         domain.StatusService{}, // 起動時実行状態FALSE
		statusSystemOperation: domain.StatusSystemData{},
		logger:                l,
		isDiff:                false,
		isDiffIndata:          false,
		isDiffOutdata:         false,
		isDiffAmountdata:      false}
}

// 用途
// Getで既存の情報を取得し、Updateで更新を行って管理を行う。
// 単一での更新は許容しない。

///////////////////
/// statusInData
///////////////////

func (c *TexMoneyNoticeManager) GetStatusInData(texCon *domain.TexContext) domain.StatusIndata {
	c.logger.Debug("【%v】取得 入金データ通知=%+v", texCon.GetUniqueKey(), c.statusInData)

	return c.statusInData
}

func (c *TexMoneyNoticeManager) UpdateStatusInData(texCon *domain.TexContext, statusIn domain.StatusIndata) bool {

	isDiffIndata := statusIn != c.statusInData
	c.statusInData = statusIn

	var s string
	s = "差分無"
	if isDiffIndata {
		s = "差分有"
	}

	c.logger.Debug("【%v】更新 入金データ通知 Diff=%v", texCon.GetUniqueKey(), s)

	return isDiffIndata
}

func (c *TexMoneyNoticeManager) DiffCheckStatusInData(texCon *domain.TexContext) bool {
	var s string
	s = "差分無"
	if c.isDiffIndata {
		s = "差分有"
	}
	c.logger.Debug("【%v】入金データ通知 差分= %t,%s", texCon.GetUniqueKey(), c.isDiffIndata, s)
	return c.isDiffIndata
}

///////////////////
/// statusOutData
///////////////////

func (c *TexMoneyNoticeManager) GetStatusOutData(texCon *domain.TexContext) domain.StatusOutdata {
	c.logger.Debug("【%v】取得 出金データ通知=%+v", texCon.GetUniqueKey(), c.statusOutData)
	return c.statusOutData
}

func (c *TexMoneyNoticeManager) UpdateStatusOutData(texCon *domain.TexContext, statusOut domain.StatusOutdata) bool {
	c.isDiffOutdata = statusOut != c.statusOutData

	c.statusOutData = statusOut
	var s string
	s = "差分無"
	if c.isDiffOutdata {
		s = "差分有"
	}
	c.logger.Debug("【%v】更新 出金データ通知 %v %v", texCon.GetUniqueKey(), statusOut, s)
	return c.isDiffOutdata
}

func (c *TexMoneyNoticeManager) DiffCheckStatusOutdata(texCon *domain.TexContext) bool {
	var s string
	s = "差分無"
	if c.isDiffOutdata {
		s = "差分有"
	}
	c.logger.Debug("【%v】出金データ通知 差分= %t,%s", texCon.GetUniqueKey(), c.isDiffOutdata, s)
	return c.isDiffOutdata
}

//No.346対応　無視しているnotice_out_dataをコントローラーまで送る為に作成
func (c *TexMoneyNoticeManager) GetStatusOutDataBillBox(texCon *domain.TexContext) domain.StatusOutdataBillBox {
	c.logger.Debug("【%v】取得 非還流庫回収時 出金データ通知情報=%+v", texCon.GetUniqueKey(), c.statusOutDataBillBox)
	return c.statusOutDataBillBox
}

func (c *TexMoneyNoticeManager) UpdateStatusBillBox(texCon *domain.TexContext, statusOutData domain.OutStatus) {
	c.statusOutDataBillBox.Amount = statusOutData.Amount
	c.statusOutDataBillBox.CountTbl = statusOutData.CountTbl
	c.statusOutDataBillBox.ExCountTbl = statusOutData.ExCountTbl
	c.logger.Debug("【%v】更新 非還流庫回収時 出金データ通知情報 %v %v", texCon.GetUniqueKey(), c.statusOutDataBillBox)
}

///////////////////
/// statusCollectData
///////////////////

func (c *TexMoneyNoticeManager) GetStatusCollectData(texCon *domain.TexContext) domain.StatusCollectData {
	c.logger.Debug("【%v】取得 回収データ通知=%+v", texCon.GetUniqueKey(), c.statusCollectData)
	return c.statusCollectData
}

func (c *TexMoneyNoticeManager) UpdateStatusCollectData(texCon *domain.TexContext, statusCollect domain.StatusCollectData) bool {
	c.isDiff = statusCollect != c.statusCollectData
	c.statusCollectData = statusCollect
	var s string
	s = "差分無"
	if c.isDiff {
		s = "差分有"
	}
	c.logger.Debug("【%v】更新 回収データ通知=%v %v", texCon.GetUniqueKey(), statusCollect, s)
	return c.isDiff
}

func (c *TexMoneyNoticeManager) DiffCheckStatusCollectData(texCon *domain.TexContext) bool {
	var s string
	s = "差分無"
	if c.isDiff {
		s = "差分有"
	}
	c.logger.Debug("【%v】回収データ通知 差分= %t,%s", texCon.GetUniqueKey(), c.isDiff, s)
	return c.isDiff
}

///////////////////
/// statusAmount
///////////////////

func (c *TexMoneyNoticeManager) GetStatusAmountData(texCon *domain.TexContext) domain.StatusAmount {
	c.logger.Debug("【%v】取得 有高データ通知=%+v", texCon.GetUniqueKey(), c.statusAmountData)
	return c.statusAmountData
}

func (c *TexMoneyNoticeManager) UpdateStatusAmountData(texCon *domain.TexContext, statusAmount domain.StatusAmount) bool {
	c.isDiffAmountdata = statusAmount != c.statusAmountData

	var s string
	s = "差分無"
	if c.isDiffAmountdata {
		s = "差分有"
	}

	c.statusAmountData = statusAmount
	c.logger.Debug("【%v】更新 有高データ通知=%+v %v", texCon.GetUniqueKey(), statusAmount, s)
	return c.isDiffAmountdata
}

func (c *TexMoneyNoticeManager) DiffCheckStatusAmountData(texCon *domain.TexContext) bool {
	var s string // c.isDiffAmountdata true:差分有り　false:差分無し
	s = "差分無"
	if c.isDiffAmountdata { //true:差分有り　false:差分無し
		s = "差分有"
	}
	c.logger.Debug("【%v】有高データ通知 差分= %t,%s", texCon.GetUniqueKey(), c.isDiffAmountdata, s)
	return c.isDiffAmountdata
}

///////////////////
/// statusCash
///////////////////

func (c *TexMoneyNoticeManager) GetStatusCashData(texCon *domain.TexContext) domain.StatusCash {
	c.logger.Debug("【%v】取得 現金入出金制御ステータス通知=%+v", texCon.GetUniqueKey(), c.statusCash)
	return c.statusCash
}

func (c *TexMoneyNoticeManager) UpdateStatusCashData(texCon *domain.TexContext, statusCash domain.StatusCash) bool {
	var isDiff bool

	defer func() {
		c.statusCash = statusCash
		var s string
		s = "差分無"
		if isDiff {
			s = "差分有"
		}
		c.logger.Debug("【%v】更新 現金入出金制御ステータス通知=%+v %v", texCon.GetUniqueKey(), statusCash, s)
	}()

	if isDiff = statusCash.CashControlId != c.statusCash.CashControlId; isDiff {
		return isDiff
	}

	if isDiff = statusCash.CashControlId != c.statusCash.CashControlId; isDiff {
		return isDiff
	}
	if isDiff = statusCash.StatusReady != c.statusCash.StatusReady; isDiff {
		return isDiff
	} //
	if isDiff = statusCash.StatusMode != c.statusCash.StatusMode; isDiff {
		return isDiff
	} //
	if isDiff = statusCash.StatusLine != c.statusCash.StatusLine; isDiff {
		return isDiff
	} //
	if isDiff = statusCash.StatusError != c.statusCash.StatusError; isDiff {
		return isDiff
	} //
	if isDiff = statusCash.ErrorCode != c.statusCash.ErrorCode; isDiff {
		return isDiff
	} //
	if isDiff = statusCash.ErrorDetail != c.statusCash.ErrorDetail; isDiff {
		return isDiff
	} //
	if isDiff = statusCash.StatusCover != c.statusCash.StatusCover; isDiff {
		return isDiff
	} //
	if isDiff = statusCash.StatusAction != c.statusCash.StatusAction; isDiff {
		return isDiff
	} //
	if isDiff = statusCash.StatusInsert != c.statusCash.StatusInsert; isDiff {
		return isDiff
	} //
	if isDiff = statusCash.StatusExit != c.statusCash.StatusExit; isDiff {
		return isDiff
	}
	if isDiff = statusCash.StatusRjbox != c.statusCash.StatusRjbox; isDiff {
		return isDiff
	} //
	if isDiff = statusCash.BillStatusTbl != c.statusCash.BillStatusTbl; isDiff {
		return isDiff
	} //
	if isDiff = statusCash.CoinStatusTbl != c.statusCash.CoinStatusTbl; isDiff {
		return isDiff
	} //

	// 紙幣残留上情報比較
	for i := 0; i < len(statusCash.BillResidueInfoTbl); i++ {
		if isDiff = statusCash.BillResidueInfoTbl[i].Title != c.statusCash.BillResidueInfoTbl[i].Title; isDiff {
			return isDiff
		} //
		if isDiff = statusCash.BillResidueInfoTbl[i].Status != c.statusCash.BillResidueInfoTbl[i].Status; isDiff {
			return isDiff
		} //
	}

	// 硬貨残留情報比較
	for j := 0; j < len(statusCash.CoinResidueInfoTbl); j++ {
		if isDiff = statusCash.CoinResidueInfoTbl[j].Title != c.statusCash.CoinResidueInfoTbl[j].Title; isDiff {
			return isDiff
		} //
		if isDiff = statusCash.CoinResidueInfoTbl[j].Status != c.statusCash.CoinResidueInfoTbl[j].Status; isDiff {
			return isDiff
		} //
	}

	// デバイス詳細情報比較
	for i := range statusCash.DeviceStatusInfoTbl {
		if isDiff = statusCash.DeviceStatusInfoTbl[i] != c.statusCash.DeviceStatusInfoTbl[i]; isDiff {
			return isDiff
		} //
	}

	// 警告情報比較
	for i := range statusCash.WarningInfoTbl {
		if isDiff = statusCash.WarningInfoTbl[i] != c.statusCash.WarningInfoTbl[i]; isDiff {
			return isDiff
		} //
	}

	return isDiff
}

///////////////////
/// statusExchange
///////////////////

func (c *TexMoneyNoticeManager) GetStatusExchangeData(texCon *domain.TexContext) domain.StatusExchange {
	c.logger.Debug("【%v】取得 両替ステータス通知=%+v", texCon.GetUniqueKey(), c.statusExchange)
	return c.statusExchange
}

func (c *TexMoneyNoticeManager) UpdateStatusExchangeData(texCon *domain.TexContext, statusExchange domain.StatusExchange) bool {

	isDiff := statusExchange != c.statusExchange

	c.statusExchange = statusExchange

	var s string
	s = "差分無"
	if isDiff {
		s = "差分有"
	}

	c.logger.Debug("【%v】更新 両替ステータス通知=%+v %v", texCon.GetUniqueKey(), c.statusExchange, s)
	return isDiff
}

///////////////////
/// statusReport
///////////////////

func (c *TexMoneyNoticeManager) GetStatusReportData(texCon *domain.TexContext) domain.StatusReport {
	c.logger.Debug("【%v】取得 入出金レポート印刷ステータス通知=%+v", texCon.GetUniqueKey(), c.statusReport)
	return c.statusReport
}

func (c *TexMoneyNoticeManager) UpdateStatusReportData(texCon *domain.TexContext, statusReport domain.StatusReport) bool {

	isDiff := statusReport != c.statusReport
	c.statusReport = statusReport

	var s string
	s = "差分無"
	if isDiff {
		s = "差分有"
	}

	c.logger.Debug("【%v】更新 入出金レポート印刷ステータス通知=%+v %v", texCon.GetUniqueKey(), c.statusReport, s)
	return isDiff
}

///////////////////
/// statusService
///////////////////

func (c *TexMoneyNoticeManager) GetStatusServiceData(texCon *domain.TexContext) domain.StatusService {
	return c.statusService
}

func (c *TexMoneyNoticeManager) UpdateStatusServiceData(texCon *domain.TexContext, statusService domain.StatusService) bool {
	isDiff := statusService != c.statusService
	c.statusService = statusService
	c.logger.Debug("【%v】更新 実行状態遷移通知=%+v Diff=%v", texCon.GetUniqueKey(), c.statusService, isDiff)
	return isDiff
}

//////////////////////////
/// statusSystemOperation
//////////////////////////

func (c *TexMoneyNoticeManager) GetStatusSystemOperationData(texCon *domain.TexContext) domain.StatusSystemData {
	c.logger.Debug("【%v】取得 システム動作モード遷移通知=%+v", texCon.GetUniqueKey(), c.statusSystemOperation)
	return c.statusSystemOperation
}

func (c *TexMoneyNoticeManager) UpdateStatusSystemOperationData(texCon *domain.TexContext, statusSystemData domain.StatusSystemData) bool {
	c.mtxStatus.Lock()
	defer c.mtxStatus.Unlock()

	isDiff := statusSystemData != c.statusSystemOperation
	c.statusSystemOperation = statusSystemData
	c.logger.Debug("【%v】更新 システム動作モード遷移通知=%+v Diff=%v", texCon.GetUniqueKey(), c.statusSystemOperation, isDiff)
	return isDiff
}
