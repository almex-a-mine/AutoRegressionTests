package interfaces

import "tex_money/domain"

// 精算機状態管理送信受信管理
type StatusTxSendRecvRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
	//アドレス設定
	SetAddressMoneyIni(moneyInit MoneyInitRepository)                          //初期補充要のアドレスを//現金入出金機制御送受信に渡す
	SetAddressMoneyExchange(moneyExchange MoneyExchangeRepository)             //両替要求のアドレスを//現金入出金機制御送受信に渡す
	SetAddressMoneyAddReplenish(moneyAddReplenish MoneyAddReplenishRepository) //追加補充要求のアドレスを//現金入出金機制御送受信に渡す
	SetAddressMoneyCollect(moneyCollect MoneyCollectRepository)                //回収要求のアドレスを//現金入出金機制御送受信に渡す
	SetAddressSetAmount(setAmount SetAmountRepository)                         //現在枚数変更要求のアドレスを//現金入出金機制御送受信に渡す
	SetAddressStatusCash(statusCash StatusCashRepository)                      //現金入出金機制御ステータス要求のアドレスを//現金入出金機制御送受信に渡す
	SetAddressPayCash(payCash PayCashRepository)                               //取引入金要求のアドレスを現金入出金機制御送受信に渡す
	SetAddressOutCash(outCash OutCashRepository)                               //取引出金要求のアドレスを//現金入出金機制御送受信に渡す
	SetAddressAmountCash(amountCash AmountCashRepository)                      //有高枚数要求のアドレスを現金入出金機制御通信に渡す
	SetAddressPrintReport(printReport PrintReportRepository)                   //入出金レポート印刷のアドレスを//現金入出金機制御送受信に渡す
	SetAddressSalesInfo(salesInfo SalesInfoRepository)                         //売上金情報取得のアドレスを//現金入出金機制御送受信に渡す
	SetAddressClearCashInfo(clearCashInfo RequestClearCashInfoRepository)      //入出金データクリア要求のアドレスを現金入出金機制御通信に渡す
	//要求
	SendRequestStatus(texCon *domain.TexContext, reqInfo *domain.RequestStatusStatusTx)
	SendRequestChangeSupply(texCon *domain.TexContext, resInfo *domain.RequestChangeSupply)
	SendRequestChangePayment(texCon *domain.TexContext, resInfo *domain.RequestChangePayment)
	SendRequestChangeStaffOperation(texCon *domain.TexContext, reqInfo *domain.RequestChangeStaffOperation) //状態変更要求(スタッフ操作記録)
}
