package interfaces

import "tex_money/domain"

// 印刷制御
type PrintSendRecvRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
	//アドレス設定
	SetAddressMoneyIni(moneyInit MoneyInitRepository)                          //初期補充要のアドレスを印刷制御通信に渡す
	SetAddressMoneyExchange(moneyExchange MoneyExchangeRepository)             //両替要求のアドレスを印刷制御通信に渡す
	SetAddressMoneyAddReplenish(moneyAddReplenish MoneyAddReplenishRepository) //追加補充要求のアドレスを印刷制御通信に渡す
	SetAddressMoneyCollect(moneyCollect MoneyCollectRepository)                //回収要求のアドレスを印刷制御通信に渡す
	SetAddressSetAmount(setAmount SetAmountRepository)                         //現在枚数変更要求のアドレスを印刷制御通信に渡す
	SetAddressStatusCash(statusCash StatusCashRepository)                      //現金入出金機制御ステータス要求のアドレスを印刷制御通信に渡す
	SetAddressPayCash(payCash PayCashRepository)                               //取引入金要求のアドレスを印刷制御通信に渡す
	SetAddressOutCash(outCash OutCashRepository)                               //取引出金要求のアドレスを印刷制御通信に渡す
	SetAddressAmountCash(amountCash AmountCashRepository)                      //有高枚数要求のアドレスを現金入出金機制御通信に渡す
	SetAddressPrintReport(printReport PrintReportRepository)                   //入出金レポート印刷のアドレスを印刷制御通信に渡す
	SetAddressSalesInfo(salesInfo SalesInfoRepository)                         //売上情報要求のアドレスを印刷制御通信に渡す
	SetAddressClearCashInfo(clearCashInfo RequestClearCashInfoRepository)      //入出金データクリア要求のアドレスを現金入出金機制御通信に渡す
	//要求
	SendRequestSupply(texCon *domain.TexContext, resInfo *domain.RequestSupply) //送信:補充レシート要求
	SendRequestStatus(texCon *domain.TexContext, resInfo *domain.RequestPrintStatus)
}
