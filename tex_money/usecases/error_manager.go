package usecases

//import (
//	"almex_astcash_cru50/domain"
//	"encoding/json"
//	"fmt"
//	"strconv"
//	"time"
//)

const (
	ERROR_NO_SUPPORT                        = 900 //サポートしていない要求
	ERROR_NOTHING                           = 998 //エラーなし
	ERROR_INSIDE                            = 999 //内部エラー
	ERROR_NOTICEINSTATUS_UNMARSHAL          = 1   //入金ステータス通知:失敗
	ERROR_NOTICEOUTSTATUS_UNMARSHAL         = 2   //出金ステータス通知:失敗
	ERROR_NOTICECOLLECTSTATUS_UNMARSHAL     = 3   //回収ステータス通知:失敗
	ERROR_NOTICEAMOUNTSTATUS_UNMARSHAL      = 4   //有高ステータス通知:失敗
	ERROR_NOTICESTATUS_UNMARSHAL            = 5   //現金入出金機ステータス通知:失敗
	ERROR_SENDREQUESTINSTART_UNMARSHAL      = 6   //入金開始要求送信失敗
	ERROR_RECVRESULTINSTART_UNMARSHAL       = 7   //入金開始要求受信失敗
	ERROR_SENDREQUESTINEND_UNMARSHAL        = 8   //入金終了要求送信失敗
	ERROR_RECVRESULTINEND_UNMARSHAL         = 9   //入金終了要求応答失敗
	ERROR_SENDREQUESTOUTSTART_UNMARSHAL     = 10  //出金開始要求送信失敗
	ERROR_RECVRESULTOUTSTART_UNMARSHAL      = 11  //出金開始要求応答失敗
	ERROR_SENDREQUESTCOLLECTSTART_UNMARSHAL = 12  //出金停止要求送信失敗
	ERROR_RECVRESULTCOLLECTSTART_UNMARSHAL  = 13  //出金停止要求応答失敗
	ERROR_SENDREQUESTOUTSTOP_UNMARSHAL      = 14  //回収開始要求送信失敗
	ERROR_RECVRESULTOUTSTOP_UNMARSHAL       = 15  //回収開始要求応答失敗
	ERROR_SENDREQUESTCOLLECTSTOP_UNMARSHAL  = 16  //回収停止要求送信失敗
	ERROR_RECVRESULTCOLLECTSTOP_UNMARSHAL   = 17  //回収停止要求応答失敗
	ERROR_NOTHING_TENTHOUSAND               = 18  //不足エラー枚数:10000
	ERROR_NOTHING_FIVETHOUSAND              = 19  //不足エラー枚数:5000
	ERROR_NOTHING_TWOTHOUSAND               = 20  //不足エラー枚数:2000
	ERROR_NOTHING_THOUSAND                  = 21  //不足エラー枚数:1000
	ERROR_NOTHING_FIVEHUNDRED               = 22  //不足エラー枚数:500
	ERROR_NOTHING_HUNDRED                   = 23  //不足エラー枚数:100
	ERROR_NOTHING_FIFTY                     = 24  //不足エラー枚数:50
	ERROR_NOTHING_TEN                       = 25  //不足エラー枚数:10
	ERROR_NOTHING_FIVE                      = 26  //不足エラー枚数:5
	ERROR_NOTHING_ONE                       = 27  //不足エラー枚数:1
	ERROR_MANY_TENTHOUSAND                  = 28  //あふれエラー:10000
	ERROR_MANY_FIVETHOUSAND                 = 29  //あふれエラー:5000
	ERROR_MANY_TWOTHOUSAND                  = 30  //あふれエラー:2000
	ERROR_MANY_THOUSAND                     = 31  //あふれエラー:1000
	ERROR_MANY_FIVEHUNDRED                  = 32  //あふれエラー:500
	ERROR_MANY_HUNDRED                      = 33  //あふれエラー:100
	ERROR_MANY_FIFTY                        = 34  //あふれエラー:50
	ERROR_MANY_TEN                          = 35  //あふれエラー:10
	ERROR_MANY_FIVE                         = 36  //あふれエラー:5
	ERROR_MANY_ONE                          = 37  //あふれエラー:1
	ERROR_MANY_ALL_COIN                     = 38  //あふれ金庫硬貨枚数が制限値以上
	ERROR_PROCESS_ID_MISMATCH               = 39  //プロセスID不一致
	ERROR_PCID_MISMATCH                     = 40  //PCID不一致
	ERROR_REQUEST_ID_MISMATCH               = 41  //リクエストID不一致
	ERROR_COMMUNICATION_FAIL                = 42  //通信失敗
	ERROR_COIN_DOOR_OPENING                 = 43  //硬貨トビラ開
	ERROR_COINS_REMAINING                   = 44  //硬貨残留あり
	ERROR_BILL_DOOR_OPEN                    = 45  //紙幣トビラ開
	ERROR_BILL_REMAINING                    = 46  //紙幣残留あり
	WARNING_NOTHING_TENTHOUSAND             = 47  //不足注意枚数:10000
	WARNING_NOTHING_FIVETHOUSAND            = 48  //不足注意枚数:5000
	WARNING_NOTHING_TWOTHOUSAND             = 49  //不足注意枚数:2000
	WARNING_NOTHING_THOUSAND                = 50  //不足注意枚数:1000
	WARNING_NOTHING_FIVEHUNDRED             = 51  //不足注意枚数:500
	WARNING_NOTHING_HUNDRED                 = 52  //不足注意枚数:100
	WARNING_NOTHING_FIFTY                   = 53  //不足注意枚数:50
	WARNING_NOTHING_TEN                     = 54  //不足注意枚数:10
	WARNING_NOTHING_FIVE                    = 55  //不足注意枚数:5
	WARNING_NOTHING_ONE                     = 56  //不足注意枚数:1
	WARNING_MANY_TENTHOUSAND                = 57  //あふれ注意:10000
	WARNING_MANY_FIVETHOUSAND               = 58  //あふれ注意:5000
	WARNING_MANY_TWOTHOUSAND                = 59  //あふれ注意:2000
	WARNING_MANY_THOUSAND                   = 60  //あふれ注意:1000
	WARNING_MANY_FIVEHUNDRED                = 61  //あふれ注意:500
	WARNING_MANY_HUNDRED                    = 62  //あふれ注意:100
	WARNING_MANY_FIFTY                      = 63  //あふれ注意:50
	WARNING_MANY_TEN                        = 64  //あふれ注意:10
	WARNING_MANY_FIVE                       = 65  //あふれ注意:5
	WARNING_MANY_ONE                        = 66  //あふれ注意:1
	WARNING_MANY_ALL_COIN                   = 67  //あふれ金庫硬貨枚数が制限値以上
	ERROR_MANY_ALL_BILL                     = 68  //あふれ金庫紙幣枚数が制限値以上
	WARNING_MANY_ALL_BILL                   = 69  //あふれ金庫紙幣枚数が制限値以上
	ERROR_NOTHING_MONEY                     = 71  //有高不足エラー
	ERROR_REPORT_ID_MISMATCH                = 72  //レポートID不一致

	ERROR_CASH_DISCREPANCY = 100 // 釣銭不一致
	// ～ 120 迄、釣銭不一致で連番として使う可能性有りの為、開けておくこと

)

type errorManager struct {
}

type errorInfo struct {
	errorTypeCode int
	errorCode     string
	convert       string
	errorDetail   string
}

var mErrorInfoTbl []errorInfo

// システムログ管理
func NewErrorManager() ErrorManager {
	createErrorInfoTbl()
	return &errorManager{}
}

// エラーログ情報初期化
func createErrorInfoTbl() {
	mErrorInfoTbl = []errorInfo{
		//共通エラー
		{ERROR_NO_SUPPORT, "TEXMY900", "", "サポートしていない要求"},
		{ERROR_NOTHING, "", "", ""},
		{ERROR_INSIDE, "TEXMY999", "", "内部エラー"},
		//入出金機エラー一覧
		{ERROR_NOTICEINSTATUS_UNMARSHAL, "TXMYE001", "", "入金ステータス通知受信失敗"},
		{ERROR_NOTICEOUTSTATUS_UNMARSHAL, "TXMYE002", "", "出金ステータス通知受信失敗"},
		{ERROR_NOTICECOLLECTSTATUS_UNMARSHAL, "TXMYE003", "", "回収ステータス通知受信失敗"},
		{ERROR_NOTICEAMOUNTSTATUS_UNMARSHAL, "TXMYE004", "", "有高ステータス通知受信失敗"},
		{ERROR_NOTICESTATUS_UNMARSHAL, "TXMYE005", "", "現金入出金機ステータス通知受失敗"},
		{ERROR_SENDREQUESTINSTART_UNMARSHAL, "TXMYE006", "", "入金開始要求送信失敗"},
		{ERROR_RECVRESULTINSTART_UNMARSHAL, "TXMYE007", "", "入金開始要求受信失敗"},
		{ERROR_SENDREQUESTINEND_UNMARSHAL, "TXMYE008", "", "入金終了要求送信失敗"},
		{ERROR_RECVRESULTINEND_UNMARSHAL, "TXMYE009", "", "入金終了要求応答失敗"},
		{ERROR_SENDREQUESTOUTSTART_UNMARSHAL, "TXMYE010", "", "出金開始要求送信失敗"},
		{ERROR_RECVRESULTOUTSTART_UNMARSHAL, "TXMYE011", "", "出金開始要求応答失敗"},
		{ERROR_SENDREQUESTCOLLECTSTART_UNMARSHAL, "TXMYE012", "", "出金停止要求送信失敗"},
		{ERROR_RECVRESULTCOLLECTSTART_UNMARSHAL, "TXMYE013", "", "出金停止要求応答失敗"},
		{ERROR_SENDREQUESTOUTSTOP_UNMARSHAL, "TXMYE014", "", "回収開始要求送信失敗"},
		{ERROR_RECVRESULTOUTSTOP_UNMARSHAL, "TXMYE015", "", "回収開始要求応答失敗"},
		{ERROR_SENDREQUESTCOLLECTSTOP_UNMARSHAL, "TXMYE016", "", "回収停止要求送信失敗"},
		{ERROR_RECVRESULTCOLLECTSTOP_UNMARSHAL, "TXMYE017", "", "回収停止要求応答失敗"},
		{ERROR_NOTHING_TENTHOUSAND, "TXMYE018", "", "不足エラー枚数:10000"},
		{ERROR_NOTHING_FIVETHOUSAND, "TXMYE019", "", "不足エラー枚数:5000"},
		{ERROR_NOTHING_TWOTHOUSAND, "TXMYE020", "", "不足エラー枚数:2000"},
		{ERROR_NOTHING_THOUSAND, "TXMYE021", "", "不足エラー枚数:1000"},
		{ERROR_NOTHING_FIVEHUNDRED, "TXMYE022", "", "不足エラー枚数:500"},
		{ERROR_NOTHING_HUNDRED, "TXMYE023", "", "不足エラー枚数:100"},
		{ERROR_NOTHING_FIFTY, "TXMYE024", "", "不足エラー枚数:50"},
		{ERROR_NOTHING_TEN, "TXMYE025", "", "不足エラー枚数:10"},
		{ERROR_NOTHING_FIVE, "TXMYE026", "", "不足エラー枚数:5"},
		{ERROR_NOTHING_ONE, "TXMYE027", "", "不足エラー枚数:1"},
		{ERROR_MANY_TENTHOUSAND, "TXMYE028", "", "あふれエラー:10000"},
		{ERROR_MANY_FIVETHOUSAND, "TXMYE029", "", "あふれエラー:5000"},
		{ERROR_MANY_TWOTHOUSAND, "TXMYE030", "", "あふれエラー:2000"},
		{ERROR_MANY_THOUSAND, "TXMYE031", "", "あふれエラー:1000"},
		{ERROR_MANY_FIVEHUNDRED, "TXMYE032", "", "あふれエラー:500"},
		{ERROR_MANY_HUNDRED, "TXMYE033", "", "あふれエラー:100"},
		{ERROR_MANY_FIFTY, "TXMYE034", "", "あふれエラー:50"},
		{ERROR_MANY_TEN, "TXMYE035", "", "あふれエラー:10"},
		{ERROR_MANY_FIVE, "TXMYE036", "", "あふれエラー:5"},
		{ERROR_MANY_ONE, "TXMYE037", "", "あふれエラー:1"},
		{ERROR_MANY_ALL_COIN, "TXMYE038", "", "あふれ金庫硬貨枚数が制限値以上"},
		{ERROR_PROCESS_ID_MISMATCH, "TXMYE039", "", "プロセスID不一致"},
		{ERROR_PCID_MISMATCH, "TXMYE040", "", "PCID不一致"},
		{ERROR_REQUEST_ID_MISMATCH, "TXMYE041", "", "リクエストID不一致"},
		{ERROR_COMMUNICATION_FAIL, "TXMYE042", "", "通信失敗"},
		{ERROR_COIN_DOOR_OPENING, "TXMYE043", "", "硬貨トビラ開"},
		{ERROR_COINS_REMAINING, "TXMYE044", "", "硬貨残留あり"},
		{ERROR_BILL_DOOR_OPEN, "TXMYE045", "", "紙幣トビラ開"},
		{ERROR_BILL_REMAINING, "TXMYE046", "", "紙幣残留あり"},
		{WARNING_NOTHING_TENTHOUSAND, "TXMYE047", "", "不足注意枚数:10000"},
		{WARNING_NOTHING_FIVETHOUSAND, "TXMYE048", "", "不足注意枚数:5000"},
		{WARNING_NOTHING_TWOTHOUSAND, "TXMYE049", "", "不足注意枚数:2000"},
		{WARNING_NOTHING_THOUSAND, "TXMYE050", "", "不足注意枚数:1000"},
		{WARNING_NOTHING_FIVEHUNDRED, "TXMYE051", "", "不足注意枚数:500"},
		{WARNING_NOTHING_HUNDRED, "TXMYE052", "", "不足注意枚数:100"},
		{WARNING_NOTHING_FIFTY, "TXMYE053", "", "不足注意枚数:50"},
		{WARNING_NOTHING_TEN, "TXMYE054", "", "不足注意枚数:10"},
		{WARNING_NOTHING_FIVE, "TXMYE055", "", "不足注意枚数:5"},
		{WARNING_NOTHING_ONE, "TXMYE056", "", "不足注意枚数:1"},
		{WARNING_MANY_TENTHOUSAND, "TXMYE057", "", "あふれ注意:10000"},
		{WARNING_MANY_FIVETHOUSAND, "TXMYE058", "", "あふれ注意:5000"},
		{WARNING_MANY_TWOTHOUSAND, "TXMYE059", "", "あふれ注意:2000"},
		{WARNING_MANY_THOUSAND, "TXMYE060", "", "あふれ注意:1000"},
		{WARNING_MANY_FIVEHUNDRED, "TXMYE061", "", "あふれ注意:500"},
		{WARNING_MANY_HUNDRED, "TXMYE062", "", "あふれ注意:100"},
		{WARNING_MANY_FIFTY, "TXMYE063", "", "あふれ注意:50"},
		{WARNING_MANY_TEN, "TXMYE064", "", "あふれ注意:10"},
		{WARNING_MANY_FIVE, "TXMYE065", "", "あふれ注意:5"},
		{WARNING_MANY_ONE, "TXMYE066", "", "あふれ注意:1"},
		{WARNING_MANY_ALL_COIN, "TXMYE067", "", "あふれ金庫硬貨枚数が制限値以上"},
		{ERROR_MANY_ALL_BILL, "TXMYE068", "", "あふれ金庫紙幣枚数が制限値以上"},
		{WARNING_MANY_ALL_BILL, "TXMYE069", "", "あふれ金庫紙幣枚数が制限値以上"},
		{ERROR_NOTHING_MONEY, "TXMYE071", "", "有高不足エラー"},
		{ERROR_REPORT_ID_MISMATCH, "TXMYE072", "", "レポートID不一致"},
		{ERROR_CASH_DISCREPANCY, "TXMYE100", "", "釣銭不一致"},
		// ～ 120 迄、釣銭不一致で連番として使う可能性有りの為、開けておくこと
	}
}

// エラー詳細情報取得
func (c *errorManager) GetErrorInfo(errorTypeCode int) (string, string) {
	var errCode string
	var errDetail string

	errCode = ""
	errDetail = ""
	for i := 0; i < len(mErrorInfoTbl); i++ {
		if errorTypeCode == mErrorInfoTbl[i].errorTypeCode {
			if len(mErrorInfoTbl[i].convert) == 0 {
				errCode = mErrorInfoTbl[i].errorCode
			} else {
				errCode = mErrorInfoTbl[i].convert
			}
			errDetail = mErrorInfoTbl[i].errorDetail
			break
		}
	}
	return errCode, errDetail
}
