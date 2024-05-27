package interfaces

import "tex_money/domain"

// 稼働データ管理
type TexdtSendRecvRepository interface {
	Start()
	Stop()
	ControlService(reqInfo domain.RequestControlService)
	InitialDbData()
	SetAddressMoneyIni(moneyInit MoneyInitRepository)                                           //初期補充要のアドレスを稼働データ管理通信に渡す
	SetAddressMoneyExchange(moneyExchange MoneyExchangeRepository)                              //両替要求のアドレスを稼働データ管理通信に渡す
	SetAddressMoneyAddReplenish(moneyAddReplenish MoneyAddReplenishRepository)                  //追加補充要求のアドレスを稼働データ管理通信に渡す
	SetAddressMoneyCollect(moneyCollect MoneyCollectRepository)                                 //回収要求のアドレスを稼働データ管理通信に渡す
	SetAddressSetAmount(setAmount SetAmountRepository)                                          //現在枚数変更要求のアドレスを稼働データ管理通信に渡す
	SetAddressStatusCash(statusCash StatusCashRepository)                                       //現金入出金機制御ステータス要求のアドレスを稼働データ管理通信に渡す
	SetAddressPayCash(payCash PayCashRepository)                                                //取引入金要求のアドレスを稼働データ管理通信に渡す
	SetAddressOutCash(outCash OutCashRepository)                                                //取引出金要求のアドレスを稼働データ管理通信に渡す
	SetAddressAmountCash(amountCash AmountCashRepository)                                       //有高枚数要求のアドレスを現金入出金機制御通信に渡す
	SetAddressPrintReport(printReport PrintReportRepository)                                    //入出金レポート印刷のアドレスを稼働データ管理通信に渡す
	SetAddressSalesInfo(salesInfo SalesInfoRepository)                                          //売上金情報要求のアドレスを稼働データ管理通信
	SetAddressClearCashInfo(clearCashInfo RequestClearCashInfoRepository)                       //入出金データクリア要求のアドレスを現金入出金機制御通信に渡す
	SendRequestReportSafeInfo(texCon *domain.TexContext, resInfo *domain.RequestReportSafeInfo) //送信:金庫情報遷移記録要求
	RecvResultReportSafeInfo(message string)                                                    //応答:金庫情報遷移記録要求
	SendRequestGetTermInfoNow(texCon *domain.TexContext, reqInfo *domain.RequestGetTermInfoNow) //送信:現在端末取得要求
	RecvResultGetTermInfoNow(message string)                                                    //応答:現在端末取得要求
	GetInitializeDbData()                                                                       // 初回起動時DBデータ取得確認
}
