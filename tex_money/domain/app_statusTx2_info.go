package domain

//補充種別
const (
	START_OF_BUSINESS_DEPOSIT_WITHDRAWAL_CLEAR_SALES = 1  //1:業務開始（入出金＆売上クリア
	START_OF_BUSINESS_CLEAR_PAYMENT_WITHDRAWAL       = 2  //2:業務開始（入出金クリア）
	START_OF_BUSINESS_CLEAR_SALES_SUP                = 3  //3:業務開始（売上クリア）
	START_OF_BUSINESS_HOLD_DATA                      = 4  //4:業務開始（データ保持）
	END_OF_BUSINESS_SUP                              = 5  //5:業務終了
	CLOSING_PROCESS                                  = 6  //6:締め処理
	REFUND_BALANCE_PAYMENT                           = 7  //7:返金残払出
	REFUND_BALANCE_WITHDRAWAL                        = 8  //8:返金残抜取
	KEY_SW_REPLENISHMENT                             = 9  //9:キーSW（補充）
	INITIAL_REPLENISHMENT_SUP                        = 10 //10:初期補充
	ADDITIONAL_REPLENISHMENT                         = 11 //11:追加補充
	EXTRACTION_SUP                                   = 12 //12:抜取
	PAYOUT_NUMBER_SUP                                = 13 //13:枚数払出
	PAYOUT_AMOUNT_SUP                                = 14 //14:金額払出
	REVERSE_EXCHANGE                                 = 15 //15:逆両替
	OVERFLOW_SAFE_COLLECTION                         = 16 //16:あふれ金庫回収
	BILL_FRONT_BOX_COLLECTION                        = 17 //17:紙幣フロントBOX回収
	REPLENISHMENT_OF_RESERVE_SAFE                    = 18 //18:予備金庫補給
	MANUAL_REPLENISHMENT_SUP                         = 19 //19:手動補充
	COLLECTION_OF_BANKNOTE_REJECT_BOX                = 20 //20:紙幣リジェクトBOX回収
)

// 決済方法
const (
	PAYWAY_CASH_SUP      = 0 //現金
	CREDIT_SUP           = 1 //クレジット
	J_DEBIT_SUP          = 2 //Jデビット
	QR_CODE_PAYMENT_SUP  = 3 //QRコード決済
	ELECTRONIC_MONEY_SUP = 4 //電子マネー
	OTHERS_SUP           = 5 //その他
)
