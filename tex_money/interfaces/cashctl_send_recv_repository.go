package interfaces

import (
	"sync"
	"tex_money/domain"
)

// 現金入出金き
type SendRecvRepository interface {
	Start()
	Stop()
	InitializeCashCtrlFlagOn(change bool) // 初回起動時に、tex_moneyが上位方の要求なくset_Amountを下位サービスに出すために利用している制御変数のon/off
	ControlService(reqInfo domain.RequestControlService)
	InitializeCashctl(wg *sync.WaitGroup)
	//アドレス設定
	SetAddressMoneyIni(moneyInit MoneyInitRepository)                          //初期補充要のアドレスを現金入出金機制御通信に渡す
	SetAddressMoneyExchange(moneyExchange MoneyExchangeRepository)             //両替要求のアドレスを現金入出金機制御通信に渡す
	SetAddressMoneyAddReplenish(moneyAddReplenish MoneyAddReplenishRepository) //追加補充要求のアドレスを現金入出金機制御通信に渡す
	SetAddressMoneyCollect(moneyCollect MoneyCollectRepository)                //回収要求のアドレスを現金入出金機制御通信に渡す
	SetAddressSetAmount(setAmount SetAmountRepository)                         //現在枚数変更要求のアドレスを現金入出金機制御通信に渡す
	SetAddressStatusCash(statusCash StatusCashRepository)                      //現金入出金機制御ステータス要求のアドレスを現金入出金機制御通信に渡す
	SetAddressPayCash(payCash PayCashRepository)                               //取引入金要求のアドレスを現金入出金機制御送受信に渡す
	SetAddressOutCash(outCash OutCashRepository)                               //取引出金要求のアドレスを現金入出金機制御通信に渡す
	SetAddressAmountCash(amountCash AmountCashRepository)                      //有高枚数要求のアドレスを現金入出金機制御通信に渡す
	SetAddressPrintReport(printReport PrintReportRepository)                   //入出金レポート印刷のアドレスを現金入出金機制御通信に渡す
	SetAddressSalesInfo(salesInfo SalesInfoRepository)                         //売上金情報要求のアドレスを現金入出金機制御通信に渡す
	SetAddressClearCashInfo(clearCashInfo RequestClearCashInfoRepository)      //入出金データクリア要求のアドレスを現金入出金機制御通信に渡す
	//要求
	SendRequestInStart(texCon *domain.TexContext, resInfo *domain.RequestInStart)                                       //送信:入金開始要求
	RecvResultInStart(message string)                                                                                   //応答:入金開始要求
	SendRequestInEnd(texCon *domain.TexContext, resInfo *domain.RequestInEnd)                                           //送信:入金終了要求
	RecvResultInEnd(message string)                                                                                     //応答:入金終了要求
	SendRequestOutStart(texCon *domain.TexContext, resInfo *domain.RequestOutStart)                                     //送信:出金開始要求
	RecvResultOutStart(message string)                                                                                  //応答:出金開始要求
	SendRequestCollectStart(texCon *domain.TexContext, resInfo *domain.RequestOutStop)                                  //送信:出金停止要求
	RecvResultCollectStart(message string)                                                                              //応答:出金停止要求
	SendRequestOutStop(texCon *domain.TexContext, resInfo *domain.RequestCollectStart)                                  //送信:回収開始要求
	RecvResultOutStop(message string)                                                                                   //応答:回収開始要求
	SendRequestCollectStop(texCon *domain.TexContext, resInfo *domain.RequestCollectStop)                               //送信:回収停止要求
	RecvResultCollectStop(message string)                                                                               //応答:回収停止要求
	SendRequestInStatus(texCon *domain.TexContext, resInfo *domain.RequestRequestInStatus)                              //送信:入金ステータス取得要求
	RecvResultInStatus(message string)                                                                                  //応答:入金ステータス取得要求
	SendRequestOutStatus(texCon *domain.TexContext, resInfo *domain.RequestOutStatus)                                   //送信:出金ステータス取得要求
	RecvResultOutStatus(message string)                                                                                 //応答:出金ステータス取得要求
	SendRequestCollectStatus(texCon *domain.TexContext, resInfo *domain.RequestCollectStatus)                           //送信:回収ステータス取得要求
	RecvResultCollectStatus(message string)                                                                             //応答:回収ステータス取得要求
	SendRequestAmountStatus(texCon *domain.TexContext, statusOfReq int, resInfo *domain.RequestAmountStatus)            //送信:有高ステータス取得要求
	RecvResultAmountStatus(message string)                                                                              //応答:有高ステータス取得要求
	SendRequestStatus(texCon *domain.TexContext, resInfo *domain.RequestStatus)                                         //送信:入出金機ステータス取得要求
	RecvResultStatus(message string)                                                                                    //応答:入出金機ステータス取得要求
	SendRequestCashctlSetAmount(texCon *domain.TexContext, resInfo *domain.RequestCashctlSetAmount, requestType int)    //送信:有高枚数変更要求
	RecvResultCashctlSetAmount(message string)                                                                          //応答:有高枚数変更要求
	SendRequestScrutinyStart(texCon *domain.TexContext, resChan chan interface{}, reqInfo *domain.RequestScrutinyStart) // 送信:精査モード開始要求
	CashId(texCon *domain.TexContext, cashId string)
}
