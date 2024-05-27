package usecases

import (
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/pkg/calculation"
)

type (
	ChangeStatus struct {
		logger                handler.LoggerRepository
		safeInfoManager       SafeInfoManager
		texMoneyNoticeManager TexMoneyNoticeManagerRepository
		texMoneyHandler       TexMoneyHandlerRepository
		refundInfo            domain.CashInfoTblInfo //返金残情報
	}

	ChangeStatusRepository interface {
		SetRefund(texCon *domain.TexContext, reqOutStart domain.RequestOutStart, outStatus domain.OutStatus, outPlanAmount int)
		RequestChangeSupply(texCon *domain.TexContext, status int) (resInfo domain.RequestChangeSupply)
		RequestChangePayment(texCon *domain.TexContext, resultCode int, resultGetSalesinfo domain.ResultGetSalesinfo) (resInfo domain.RequestChangePayment)
	}
)

func NewChangeStatus(logger handler.LoggerRepository,
	safeInfoManager SafeInfoManager,
	texMoneyNoticeManager TexMoneyNoticeManagerRepository,
	texMoneyHandler TexMoneyHandlerRepository,
) ChangeStatusRepository {
	return &ChangeStatus{
		logger:                logger,
		safeInfoManager:       safeInfoManager,
		texMoneyNoticeManager: texMoneyNoticeManager,
		texMoneyHandler:       texMoneyHandler,
		refundInfo:            domain.CashInfoTblInfo{}}
}

// SetRefund 状態セット
func (c *ChangeStatus) SetRefund(texCon *domain.TexContext, reqOutStart domain.RequestOutStart, outStatus domain.OutStatus, outPlanAmount int) {
	c.logger.Trace("【%v】texMoneyHandler SetRefundCountTbl reqOutStart=%+v, outStatus=%+v", texCon.GetUniqueKey(), reqOutStart, outStatus)
	// 返金残金額
	// 出金予定額 - 実際の出金額
	amount := reqOutStart.Amount
	if amount == 0 {
		amount = outPlanAmount
	}

	c.refundInfo.Amount = amount - outStatus.Amount
	if c.refundInfo.Amount < 0 {
		c.refundInfo.Amount = 0
	}

	// TODO 以降のロジックで問題が発生する場合には、1-3独自で返金残をベースに逆両替算出を行って通知するように修正することも検討する事
	// ※何故か硬貨のメイン、サブどちらもマイナスであったりした場合等

	// 返金残の内訳（拡張金種別枚数）
	// 出金予定枚数 - 実際の出金枚数
	for i, v := range reqOutStart.CountTbl {
		c.refundInfo.ExCountTbl[i] = v - outStatus.ExCountTbl[i]

	}

	// 硬貨配列のみメイン、サブの排出が指定できずマイナス値となる可能性がある為、チェックする。
	// メインがマイナスの場合にはサブ、サブがマイナスの場合にはメインから出力するように変更する。
	for i := 0; i < 6; i++ {

		// メインがマイナスの場合
		if c.refundInfo.ExCountTbl[i+4] < 0 {
			// マイナス分の値を
			c.refundInfo.ExCountTbl[i+10] += c.refundInfo.ExCountTbl[i+4] * -1
			c.refundInfo.ExCountTbl[i+4] = 0
		}

		// サブがマイナスの場合
		if c.refundInfo.ExCountTbl[i+10] < 0 {
			c.refundInfo.ExCountTbl[i+4] += c.refundInfo.ExCountTbl[i+10] * -1
			c.refundInfo.ExCountTbl[i+10] = 0
		}
	}

	// ExTblから、CountTblを作成
	c.refundInfo.CountTbl = calculation.NewCassette(c.refundInfo.ExCountTbl).ExCountTblToTenCountTbl()

	c.logger.Debug("【%v】texMoneyHandler refundInfo=%+v", texCon.GetUniqueKey(), c.refundInfo)
}

// RequestChangeSupply 状態変更要求
func (c *ChangeStatus) RequestChangeSupply(texCon *domain.TexContext, status int) (resInfo domain.RequestChangeSupply) {
	var supplyType, payOutBalance int

	// 入出金データ取得
	statusIndata := c.texMoneyNoticeManager.GetStatusInData(texCon)
	statusOut := c.texMoneyNoticeManager.GetStatusOutData(texCon)
	statusOutDataBillBox := c.texMoneyNoticeManager.GetStatusOutDataBillBox(texCon)

	// 入出金情報作成
	cashInfoTbl := make([]domain.CashInfoTblInfo, 6)
	for i := range cashInfoTbl {
		cashInfoTbl[i].InfoType = domain.InfoType[i]
	}

	// 入金情報をセットする
	setInData := func(cashInfoTbl []domain.CashInfoTblInfo) []domain.CashInfoTblInfo {
		cashInfoTbl[0].Amount = statusIndata.Amount
		cashInfoTbl[0].CountTbl = statusIndata.CountTbl
		cashInfoTbl[0].ExCountTbl = statusIndata.ExCountTbl
		return cashInfoTbl
	}

	// 出金情報をセットする
	setOutData := func(cashInfoTbl []domain.CashInfoTblInfo) []domain.CashInfoTblInfo {
		cashInfoTbl[2].Amount = -1 * statusOut.Amount
		for i, c := range statusOut.CountTbl {
			cashInfoTbl[2].CountTbl[i] = -1 * c
		}
		for i, ex := range statusOut.ExCountTbl {
			cashInfoTbl[2].ExCountTbl[i] = -1 * ex
		}
		return cashInfoTbl
	}

	// 返金残情報をセットする
	setRefundData := func(cashInfoTbl []domain.CashInfoTblInfo) []domain.CashInfoTblInfo {
		// 出金予定
		cashInfoTbl[1].Amount = -1 * c.refundInfo.Amount
		for i, rc := range c.refundInfo.CountTbl {
			cashInfoTbl[1].CountTbl[i] = -1 * rc
		}
		for i, rEx := range c.refundInfo.ExCountTbl {
			cashInfoTbl[1].ExCountTbl[i] = -1 * rEx
		}
		// エラー出金済
		cashInfoTbl[5].Amount = -1 * statusOut.Amount
		for i, errC := range statusOut.CountTbl {
			cashInfoTbl[5].CountTbl[i] = -1 * errC
		}
		for i, errEx := range statusOut.ExCountTbl {
			cashInfoTbl[5].ExCountTbl[i] = -1 * errEx
		}
		return cashInfoTbl
	}

	// 出金予定情報をセットする
	setOutPlan := func(cashInfoTbl []domain.CashInfoTblInfo) []domain.CashInfoTblInfo {
		cashInfoTbl[1].Amount = -1 * statusOut.Amount
		for i, planC := range statusOut.CountTbl {
			cashInfoTbl[1].CountTbl[i] = -1 * planC
		}
		for i, planEx := range statusOut.ExCountTbl {
			cashInfoTbl[1].ExCountTbl[i] = -1 * planEx
		}
		return cashInfoTbl
	}

	//非還流金庫回収データを出金額と出金予定情報にセットする
	setBillBoxData := func(cashInfoTbl []domain.CashInfoTblInfo) []domain.CashInfoTblInfo {
		// 出金情報
		cashInfoTbl[2].Amount = -1 * statusOutDataBillBox.Amount
		for i, c := range statusOutDataBillBox.CountTbl {
			cashInfoTbl[2].CountTbl[i] = -1 * c
		}
		for i, ex := range statusOutDataBillBox.ExCountTbl {
			cashInfoTbl[2].ExCountTbl[i] = -1 * ex
		}
		// 出金予定
		cashInfoTbl[1].Amount = -1 * statusOutDataBillBox.Amount
		for i, rc := range statusOutDataBillBox.CountTbl {
			cashInfoTbl[1].CountTbl[i] = -1 * rc
		}
		for i, rEx := range statusOutDataBillBox.ExCountTbl {
			cashInfoTbl[1].ExCountTbl[i] = -1 * rEx
		}
		return cashInfoTbl
	}

	switch c.texMoneyHandler.GetSequence(texCon) {
	case domain.INITIAL_ADDING_CONFIRM: //初期補充確定
		supplyType = domain.INITIAL_REPLENISHMENT_SUP
		cashInfoTbl = setInData(cashInfoTbl)

	case domain.INITIAL_ADDING_UPDATE: //初期補充更新
		supplyType = domain.INITIAL_REPLENISHMENT_SUP

	case domain.NUMBER_OF_WITHDRAW_DESIGNATED_CONFIRM, //出金枚数指定両替確定時
		domain.REVERSE_EXCHANGEING_CONFIRM_OUTDATA: //逆両替確定出金データ時
		supplyType = domain.REVERSE_EXCHANGE
		cashInfoTbl = setInData(cashInfoTbl)
		cashInfoTbl = setOutPlan(cashInfoTbl)
		cashInfoTbl = setOutData(cashInfoTbl)

	case domain.MONEY_ADD_REPLENISH_CONFIRM: //追加補充
		supplyType = domain.ADDITIONAL_REPLENISHMENT
		cashInfoTbl = setInData(cashInfoTbl)

	case domain.REJECTBOXCOLLECT_START: //リジェクトボックス回収
		supplyType = domain.COLLECTION_OF_BANKNOTE_REJECT_BOX
		//入金出金してないのでデータ格納しない

	case domain.UNRETURNEDCOLLECT_START: //非還流庫回収開始
		supplyType = domain.BILL_FRONT_BOX_COLLECTION //17:紙幣フロントBOX回収
		cashInfoTbl = setBillBoxData(cashInfoTbl)

	case domain.MIDDLE_START_OUT_START, //途中回収出金開始
		domain.ALLCOLLECT_START_OUT_START,     //全回収開始出金開始
		domain.ALLCOLLECT_START_OUT_STOP,      //全回収開始出金停止
		domain.ALLCOLLECT_START_COLLECT_START, //全回収開始回収開始
		domain.ALLCOLLECT_START_COLLECT_STOP:  //全回収開始回収停止
		supplyType = domain.PAYOUT_NUMBER_SUP //枚数払出
		cashInfoTbl = setOutPlan(cashInfoTbl)
		cashInfoTbl = setOutData(cashInfoTbl)

	case domain.MANUAL_REPLENISHMENT_COLLECTION: //手動回収
		supplyType = domain.MANUAL_REPLENISHMENT_SUP //手動補充
		cashInfoTbl = setInData(cashInfoTbl)
		cashInfoTbl = setOutPlan(cashInfoTbl)
		cashInfoTbl = setOutData(cashInfoTbl)

	case domain.SALESMONEY_START: //売上金回収
		supplyType = domain.PAYOUT_AMOUNT_SUP //金額払出
		cashInfoTbl = setInData(cashInfoTbl)
		cashInfoTbl = setOutPlan(cashInfoTbl)
		cashInfoTbl = setOutData(cashInfoTbl)

	case domain.TRANSACTION_OUT_REFUND_PAYMENT_OUT_START: //取引出金 返金残払出開始
		if status == 0 { // 返金残有りパターン
			supplyType = domain.REFUND_BALANCE_PAYMENT
			payOutBalance = c.refundInfo.Amount - statusOut.Amount
			cashInfoTbl = setOutPlan(cashInfoTbl)
			cashInfoTbl = setOutData(cashInfoTbl)
			cashInfoTbl = setRefundData(cashInfoTbl)
		} else { // 返金残無しパターン
			supplyType = domain.REFUND_BALANCE_WITHDRAWAL
			cashInfoTbl = setOutPlan(cashInfoTbl)
			cashInfoTbl = setOutData(cashInfoTbl)
			cashInfoTbl = setOutPlan(cashInfoTbl)
		}
	}

	//金庫情報取得
	safeInfo := c.safeInfoManager.GetSafeInfo(texCon)

	return *domain.NewRequestChangeSupply(c.texMoneyHandler.NewRequestInfo(texCon),
		supplyType,
		domain.InfoTrade{
			BillingAmount:     cashInfoTbl[0].Amount + cashInfoTbl[2].Amount,
			DepositAmount:     cashInfoTbl[0].Amount,
			PaymentPlanAmount: cashInfoTbl[1].Amount,
			PaymentAmount:     cashInfoTbl[2].Amount,
			PayoutBalance:     payOutBalance,
			CashInfoTbl:       cashInfoTbl,
		},
		domain.InfoSafeInfo{
			CurrentStatusTbl: c.texMoneyHandler.MakeCurrentStatusTbl(texCon),
			SortInfoTbl:      safeInfo.SortInfoTbl[:],
		},
	)
}

// 状態変更要求（精算完了）
func (c *ChangeStatus) RequestChangePayment(texCon *domain.TexContext, resultCode int, resultGetSalesinfo domain.ResultGetSalesinfo) (resInfo domain.RequestChangePayment) {
	// 入出金データ取得
	statusIndata := c.texMoneyNoticeManager.GetStatusInData(texCon)
	statusOut := c.texMoneyNoticeManager.GetStatusOutData(texCon)

	// 入出金情報作成
	cashInfoTbl := make([]domain.CashInfoTblInfo, 6)
	for i := range cashInfoTbl {
		cashInfoTbl[i].InfoType = domain.InfoType[i]
	}

	// 入金情報をセットする
	cashInfoTbl[0].Amount = statusIndata.Amount
	cashInfoTbl[0].CountTbl = statusIndata.CountTbl
	cashInfoTbl[0].ExCountTbl = statusIndata.ExCountTbl
	// 出金情報をセットする
	cashInfoTbl[2].Amount = statusOut.Amount
	cashInfoTbl[2].CountTbl = statusOut.CountTbl
	cashInfoTbl[2].ExCountTbl = statusOut.ExCountTbl

	//金庫情報取得
	safeInfo := c.safeInfoManager.GetSafeInfo(texCon)

	return *domain.NewRequestChangePayment(c.texMoneyHandler.NewRequestInfo(texCon),
		9, //両替
		"EXCHANGE",
		resultCode,
		domain.InfoTrade{
			DepositAmount: statusIndata.Amount,
			PaymentAmount: statusOut.Amount,
			CashInfoTbl:   cashInfoTbl,
		},
		domain.InfoSales{
			SalesAmount:   resultGetSalesinfo.InfoSales.SalesAmount,
			ExchangeTotal: statusOut.Amount + resultGetSalesinfo.InfoSales.ExchangeTotal,
			SalesTypeTbl:  resultGetSalesinfo.InfoSales.SalesTypeTbl,
		},
		domain.InfoSafeInfo{
			CurrentStatusTbl: c.texMoneyHandler.MakeCurrentStatusTbl(texCon),
			SortInfoTbl:      safeInfo.SortInfoTbl[:],
		})
}
