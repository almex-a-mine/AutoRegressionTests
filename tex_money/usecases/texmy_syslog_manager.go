package usecases

import (
	"encoding/json"
	"fmt"
	"strconv"
	"tex_money/domain"
	"tex_money/domain/handler"
	"time"
)

const (
	SYSLOG_LOGTYPE_SERVICESTART_SUCCESS                        = 1  //サービス開始成功
	SYSLOG_LOGTYPE_SERVICESTART_FATAL                          = 2  //サービス開始失敗
	SYSLOG_LOGTYPE_SERVICESTOP_SUCCESS                         = 3  //サービス終了成功
	SYSLOG_LOGTYPE_SERVICESTOP_FATAL                           = 4  //サービス終了失敗
	SYSLOG_LOGTYPE_REQUEST_MONEY_INIT_SUCCESS                  = 5  //初期補充要求成功
	SYSLOG_LOGTYPE_REQUEST_MONEY_INIT_FATAL                    = 6  //初期補充要求失敗
	SYSLOG_LOGTYPE_REQUEST_MONEY_EXCHANGE_SUCCESS              = 7  //両替要求成功
	SYSLOG_LOGTYPE_REQUEST_MONEY_EXCHANGE_FATAL                = 8  //両替要求失敗
	SYSLOG_LOGTYPE_REQUEST_MONEY_ADD_REPLENISH_SUCCESS         = 9  //追加補充要求成功
	SYSLOG_LOGTYPE_REQUEST_MONEY_ADD_REPLENISH_FATAL           = 10 //追加補充要求失敗
	SYSLOG_LOGTYPE_REQUEST_MONEY_COLLECT_SUCCESS               = 11 //回収要求成功
	SYSLOG_LOGTYPE_REQUEST_MONEY_COLLECT_FATAL                 = 12 //回収要求失敗
	SYSLOG_LOGTYPE_REQUEST_SET_AMOUNT_SUCCESS                  = 13 //現在枚数変更要求成功
	SYSLOG_LOGTYPE_REQUEST_SET_AMOUNT_FATAL                    = 14 //現在枚数変更要求失敗
	SYSLOG_LOGTYPE_REQUEST_STATUS_CASH_SUCCESS                 = 15 //現金入出金機制御ステータス要求成功
	SYSLOG_LOGTYPE_REQUEST_STATUS_CASH_FATAL                   = 16 //現金入出金機制御ステータス要求失敗
	SYSLOG_LOGTYPE_REQUEST_PAY_CASH_SUCCESS                    = 17 //取引入金要求成功
	SYSLOG_LOGTYPE_REQUEST_PAY_CASH_FATAL                      = 18 //取引入金要求失敗
	SYSLOG_LOGTYPE_REQUEST_OUT_CASH_SUCCESS                    = 19 //取引出金要求成功
	SYSLOG_LOGTYPE_REQUEST_OUT_CASH_FATAL                      = 20 //取引出金要求失敗
	SYSLOG_LOGTYPE_REQUEST_AMOUNT_CASH_SUCCESS                 = 21 //有高枚数要求成功
	SYSLOG_LOGTYPE_REQUEST_AMOUNT_CASH_FATAL                   = 22 //有高枚数要求失敗
	SYSLOG_LOGTYPE_REQUEST_PRINT_REPORT_SUCCESS                = 23 //入出金レポート印刷要求成功
	SYSLOG_LOGTYPE_REQUEST_PRINT_REPORT_FATAL                  = 24 //入出金レポート印刷要求失敗
	SYSLOG_LOGTYPE_REQUEST_SALES_INFO_SUCCESS                  = 25 //売上金情報要求成功
	SYSLOG_LOGTYPE_REQUEST_SALES_INFO_FATAL                    = 26 //売上金情報要求失敗
	SYSLOG_LOGTYPE_REQUEST_CLEAR_CASHINFO_SUCCESS              = 27 //入出金データクリア要求成功
	SYSLOG_LOGTYPE_REQUEST_CLEAR_CASHINFO_FATAL                = 28 //入出金データクリア要求失敗
	SYSLOG_LOGTYPE_REQUEST_CHANGE_SYSTEM_OPERATION_SUCCESS     = 29 //システム動作モード変更要求成功
	SYSLOG_LOGTYPE_REQUEST_CHANGE_SYSTEM_OPERATION_FATAL       = 30 //システム動作モード変更要求失敗
	SYSLOG_LOGTYPE_NOTICE_INDATA_SUCCESS                       = 31 //入金データ通知成功(直し)
	SYSLOG_LOGTYPE_NOTICE_INDATA_FATAL                         = 32 //入金データ通知失敗
	SYSLOG_LOGTYPE_NOTICE_OUTDATA_SUCCESS                      = 33 //出金データ通知成功
	SYSLOG_LOGTYPE_NOTICE_OUTDATA_FATAL                        = 34 //出金データ通知失敗
	SYSLOG_LOGTYPE_NOTICE_COLLECTDATA_SUCCESS                  = 35 //回収データ通知成功
	SYSLOG_LOGTYPE_NOTICE_COLLECTDATA_FATAL                    = 36 //回収データ通知失敗
	SYSLOG_LOGTYPE_NOTICE_AMOUNT_SUCCESS                       = 37 //有高データ通知成功
	SYSLOG_LOGTYPE_NOTICE_AMOUNT_FATAL                         = 38 //有高データ通知失敗
	SYSLOG_LOGTYPE_NOTICE_STATUS_CASH_SUCCESS                  = 39 //現金入出金機制御ステータス通知成功
	SYSLOG_LOGTYPE_NOTICE_STATUS_CASH_FATAL                    = 40 //現金入出金機制御ステータス通知失敗
	SYSLOG_LOGTYPE_NOTICE_REPORT_STATUS_SUCCESS                = 41 //入出金レポート印刷ステータス通知成功
	SYSLOG_LOGTYPE_NOTICE_REPORT_STATUS_FATAL                  = 42 //入出金レポート印刷ステータス通知失敗
	SYSLOG_LOGTYPE_NOTICE_STATUS_SYSTEM_OPERATION_SUCCESS      = 43 //システム動作モード遷移通知成功
	SYSLOG_LOGTYPE_NOTICE_STATUS_SYSTEM_OPERATION_ERROR        = 44 //システム動作モード遷移通知失敗
	SYSLOG_LOGTYPE_NOTICE_STATUS_EXCHANGE_SUCCESS              = 45 //両替ステータス通知成功
	SYSLOG_LOGTYPE_NOTICE_STATUS_EXCHANGE_FATAL                = 46 //両替ステータス通知失敗
	SYSLOG_LOGTYPE_REQUEST_MAINTENANCE_MODE_SUCCESS            = 47 //保守業務モード要求成功
	SYSLOG_LOGTYPE_REQUEST_MAINTENANCE_MODE_FATAL              = 48 //保守業務モード要求失敗
	SYSLOG_LOGTYPE_REQUEST_COIN_CASSETTE_CONTROL_SUCCESS       = 49 //カセット交換要求成功
	SYSLOG_LOGTYPE_REQUEST_COIN_CASSETTE_CONTROL_FATAL         = 50 //カセット交換要求失敗
	SYSLOG_LOGTYPE_REQUEST_REVERSE_EXCHNGE_CALCULATION_SUCCESS = 51 //逆両替算出要求成功
	SYSLOG_LOGTYPE_REQUEST_REVERSE_EXCHNGE_CALCULATION_FATAL   = 52 //逆両替算出要求失敗
	SYSLOG_LOGTYPE_REQUEST_GET_SAFEINFO_SUCCESS                = 53 //金庫情報取得要求成功
	SYSLOG_LOGTYPE_REQUEST_GET_SAFEINFO_FATAL                  = 54 //金庫情報取得要求失敗
	SYSLOG_LOGTYPE_REQUEST_REGISTER_MONEYSETTING_SUCCESS       = 55 //金銭設定登録要求成功
	SYSLOG_LOGTYPE_REQUEST_REGISTER_MONEYSETTING_FATAL         = 56 //金銭設定登録要求失敗
	SYSLOG_LOGTYPE_REQUEST_GET_MONEYSETTING_SUCCESS            = 57 //金銭設定登録要求成功
	SYSLOG_LOGTYPE_REQUEST_GET_MONEYSETTING_FATAL              = 58 //金銭設定登録要求失敗
	SYSLOG_LOGTYPE_REQUEST_SCRUTINY_SUCCESS                    = 59 //精査モード要求成功
	SYSLOG_LOGTYPE_REQUEST_SCRUTINY_FATAL                      = 60 //精査モード要求失敗
)

type syslogManager struct {
	mqtt handler.MqttRepository
}

type syslogInfo struct {
	EventType int
	EventCode string
	EventName string
	LogLebel  string
}

var mSyslogInfoTbl []syslogInfo

// システムログ管理
func NewSysLogManager(mqtt handler.MqttRepository) SyslogManager {
	initialSyslogInfo()
	return &syslogManager{
		mqtt: mqtt,
	}
}

// システムログ情報初期化
func initialSyslogInfo() {
	mSyslogInfoTbl = []syslogInfo{
		{SYSLOG_LOGTYPE_SERVICESTART_SUCCESS, "TXMY0001", "サービス開始(入出金管理)", "DEBUG"},
		{SYSLOG_LOGTYPE_SERVICESTART_FATAL, "TXMY0002", "サービス開始失敗(入出金管理)", "FATAL"},
		{SYSLOG_LOGTYPE_SERVICESTOP_SUCCESS, "TXMY0003", "サービス停止(入出金管理)", "DEBUG"},
		{SYSLOG_LOGTYPE_SERVICESTOP_FATAL, "TXMY0004", "サービス停止失敗(入出金管理)", "FATAL"},
		{SYSLOG_LOGTYPE_REQUEST_MONEY_INIT_SUCCESS, "TXMY100", "入出金制御 初期補充要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_MONEY_INIT_FATAL, "TXMY101", "入出金制御 初期補充要求", "FATAL"},
		{SYSLOG_LOGTYPE_REQUEST_MONEY_EXCHANGE_SUCCESS, "TXMY102", "入出金制御 両替要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_MONEY_EXCHANGE_FATAL, "TXMY103", "入出金制御 両替要求", "FATAL"},
		{SYSLOG_LOGTYPE_REQUEST_MONEY_ADD_REPLENISH_SUCCESS, "TXMY104", "入出金制御 追加補充要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_MONEY_ADD_REPLENISH_FATAL, "TXMY105", "入出金制御 追加補充要求", "FATAL"},
		{SYSLOG_LOGTYPE_REQUEST_MONEY_COLLECT_SUCCESS, "TXMY106", "入出金制御 回収要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_MONEY_COLLECT_FATAL, "TXMY107", "入出金制御 回収要求", "FATAL"},
		{SYSLOG_LOGTYPE_REQUEST_SET_AMOUNT_SUCCESS, "TXMY108", "入出金制御 現在枚数変更要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_SET_AMOUNT_FATAL, "TXMY109", "入出金制御 現在枚数変更要求", "FATAL"},
		{SYSLOG_LOGTYPE_REQUEST_STATUS_CASH_SUCCESS, "TXMY110", "入出金制御 現金入出金機制御ステータス要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_STATUS_CASH_FATAL, "TXMY111", "入出金制御 現金入出金機制御ステータス要求", "FATAL"},
		{SYSLOG_LOGTYPE_REQUEST_PAY_CASH_SUCCESS, "TXMY112", "入出金制御 取引入金要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_PAY_CASH_FATAL, "TXMY113", "入出金制御 取引入金要求", "FATAL"},
		{SYSLOG_LOGTYPE_REQUEST_OUT_CASH_SUCCESS, "TXMY114", "入出金制御 取引出金要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_OUT_CASH_FATAL, "TXMY115", "入出金制御 取引出金要求", "FATAL"},
		{SYSLOG_LOGTYPE_REQUEST_AMOUNT_CASH_SUCCESS, "TXMY116", "入出金制御 有高枚数要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_AMOUNT_CASH_FATAL, "TXMY117", "入出金制御 有高枚数要求", "FATAL"},
		{SYSLOG_LOGTYPE_REQUEST_PRINT_REPORT_SUCCESS, "TXMY118", "入出金制御 入出金レポート印刷要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_PRINT_REPORT_FATAL, "TXMY119", "入出金制御 入出金レポート印刷要求", "FATAL"},
		{SYSLOG_LOGTYPE_REQUEST_SALES_INFO_SUCCESS, "TXMY120", "入出金制御 売上金情報要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_SALES_INFO_FATAL, "TXMY121", "入出金制御 売上金情報要求", "FATAL"},
		{SYSLOG_LOGTYPE_REQUEST_CLEAR_CASHINFO_SUCCESS, "TXMY122", "入出金制御 入出金データクリア要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_CLEAR_CASHINFO_FATAL, "TXMY123", "入出金制御 入出金データクリア要求", "FATAL"},
		{SYSLOG_LOGTYPE_REQUEST_CHANGE_SYSTEM_OPERATION_SUCCESS, "TXMY124", "入出金制御 システム動作モード変更要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_CHANGE_SYSTEM_OPERATION_FATAL, "TXMY125", "入出金制御 システム動作モード変更要求", "FATAL"},
		{SYSLOG_LOGTYPE_REQUEST_MAINTENANCE_MODE_SUCCESS, "TXMY126", "入出金制御 保守業務モード要求成功", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_MAINTENANCE_MODE_FATAL, "TXMY127", "入出金制御 保守業務モード要求成功", "FATAL"},
		{SYSLOG_LOGTYPE_REQUEST_COIN_CASSETTE_CONTROL_SUCCESS, "TXMY128", "硬貨カセット操作要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_COIN_CASSETTE_CONTROL_FATAL, "TXMY129", "硬貨カセット操作要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_GET_SAFEINFO_SUCCESS, "TXMY130", "入出金制御 金庫情報取得要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_GET_SAFEINFO_FATAL, "TXMY131", "入出金制御 金庫情報取得要求", "FATAL"},
		{SYSLOG_LOGTYPE_REQUEST_REGISTER_MONEYSETTING_SUCCESS, "TXMY132", "入出金制御 金銭設定登録要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_REGISTER_MONEYSETTING_FATAL, "TXMY133", "入出金制御 金銭設定登録要求", "FATAL"},
		{SYSLOG_LOGTYPE_REQUEST_GET_MONEYSETTING_SUCCESS, "TXMY134", "入出金制御 金銭設定取得要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_GET_MONEYSETTING_FATAL, "TXMY135", "入出金制御 金銭設定取得要求", "FATAL"},
		{SYSLOG_LOGTYPE_NOTICE_INDATA_SUCCESS, "TXMY301", "入出金制御 入金データ通知", "DEBUG"},
		{SYSLOG_LOGTYPE_NOTICE_INDATA_FATAL, "TXMY302", "入出金制御 入金データ通知", "FATAL"},
		{SYSLOG_LOGTYPE_NOTICE_OUTDATA_SUCCESS, "TXMY303", "入出金制御 出金データ通知", "DEBUG"},
		{SYSLOG_LOGTYPE_NOTICE_OUTDATA_FATAL, "TXMY304", "入出金制御 出金データ通知", "FATAL"},
		{SYSLOG_LOGTYPE_NOTICE_COLLECTDATA_SUCCESS, "TXMY305", "入出金制御 回収データ通知", "DEBUG"},
		{SYSLOG_LOGTYPE_NOTICE_COLLECTDATA_FATAL, "TXMY306", "入出金制御 回収データ通知", "FATAL"},
		{SYSLOG_LOGTYPE_NOTICE_AMOUNT_SUCCESS, "TXMY307", "入出金制御 有高データ通知", "DEBUG"},
		{SYSLOG_LOGTYPE_NOTICE_AMOUNT_FATAL, "TXMY308", "入出金制御 有高データ通知", "FATAL"},
		{SYSLOG_LOGTYPE_NOTICE_STATUS_CASH_SUCCESS, "TXMY309", "入出金制御 現金入出金機制御ステータス通知", "DEBUG"},
		{SYSLOG_LOGTYPE_NOTICE_STATUS_CASH_FATAL, "TXMY310", "入出金制御 現金入出金機制御ステータス通知", "FATAL"},
		{SYSLOG_LOGTYPE_NOTICE_REPORT_STATUS_SUCCESS, "TXMY311", "入出金制御 入出金レポート印刷ステータス通知", "DEBUG"},
		{SYSLOG_LOGTYPE_NOTICE_REPORT_STATUS_FATAL, "TXMY312", "入出金制御 入出金レポート印刷ステータス通知", "FATAL"},
		{SYSLOG_LOGTYPE_NOTICE_STATUS_SYSTEM_OPERATION_SUCCESS, "TXMY0310", "システム動作モード遷移通知成功", "DEBUG"},
		{SYSLOG_LOGTYPE_NOTICE_STATUS_SYSTEM_OPERATION_ERROR, "TXMY0311", "システム動作モード遷移通知失敗", "ERROR"},
		{SYSLOG_LOGTYPE_NOTICE_STATUS_EXCHANGE_SUCCESS, "TXMY0313", "両替ステータス通知成功", "DEBUG"},
		{SYSLOG_LOGTYPE_NOTICE_STATUS_EXCHANGE_FATAL, "TXMY0313", "両替ステータス通知失敗", "ERROR"},
		{SYSLOG_LOGTYPE_REQUEST_COIN_CASSETTE_CONTROL_SUCCESS, "TXMY0317", "硬貨カセット操作要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_COIN_CASSETTE_CONTROL_FATAL, "TXMY0318", "硬貨カセット操作要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_REVERSE_EXCHNGE_CALCULATION_SUCCESS, "TXMY0319", "逆両替算出要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_REVERSE_EXCHNGE_CALCULATION_FATAL, "TXMY0320", "逆両替算出要求", "FATAL"},
		{SYSLOG_LOGTYPE_REQUEST_SCRUTINY_SUCCESS, "TXMY0321", "精査モード要求", "DEBUG"},
		{SYSLOG_LOGTYPE_REQUEST_SCRUTINY_FATAL, "TXMY0322", "精査モード要求", "FATAL"},
	}
}

// システムログ情報取得
func (s *syslogManager) getlSyslogInfo(eventType int, deviceNo string) (string, string, string) {
	var logLebel string
	var eventCode string
	var eventName string

	logLebel = ""
	eventCode = ""
	eventName = ""
	for i := 0; i < len(mSyslogInfoTbl); i++ {
		if eventType == mSyslogInfoTbl[i].EventType {
			logLebel = mSyslogInfoTbl[i].LogLebel
			if len(deviceNo) != 0 {
				eventCode = fmt.Sprintf("%v-%v%v", mSyslogInfoTbl[i].EventCode, deviceNo[10:12], deviceNo[12:])
			} else {
				eventCode = mSyslogInfoTbl[i].EventCode
			}
			eventName = mSyslogInfoTbl[i].EventName
			break
		}
	}
	return eventCode, logLebel, eventName
}

// systtemログ発行
func (s *syslogManager) NoticeSystemLog(logType int, deviceNo string, logData string) {

	eventcode, loglebel, logSummary := s.getlSyslogInfo(logType, "")

	nowTime := fmt.Sprintf("%v", time.Now().Format("20060102150405.000"))

	dateTemp := fmt.Sprintf("%v", nowTime[0:8])
	TimeTemp := fmt.Sprintf("%v%v", nowTime[8:14], nowTime[15:])

	gDate, err := strconv.Atoi(dateTemp)
	if err != nil {
		gDate = 0
	}
	gTime, err := strconv.Atoi(TimeTemp)
	if err != nil {
		gTime = 0
	}

	noticeSystemLogInfo := domain.NoticeSystemLogInfo{
		GenerateDate: gDate,
		GenerateTime: gTime,
		ForceEncrypt: false,
		LogLevel:     loglebel,
		RequestId:    eventcode,
		LogSummary:   logSummary,
		LogData:      logData,
	}

	encodedJson, _ := json.Marshal(noticeSystemLogInfo)
	s.mqtt.Publish("/almex/function/syslog/request_log", string(encodedJson))
}
