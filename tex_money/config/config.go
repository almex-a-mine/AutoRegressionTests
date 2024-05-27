package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"tex_money/domain"
	"tex_money/pkg/file"
	"tex_money/pkg/pc"

	"gopkg.in/ini.v1"
)

// システム情報
type SystemConf struct {
	MaxLength       int
	MaxRotation     int
	StartUpStatus   int // 起動時のサービス動作状況	// 0:動作として起動 1:停止として起動
	LogStopInfo     bool
	LogStopTrace    bool
	LogStopMqtt     bool
	LogStopDebug    bool
	LogStopMutex    bool
	LogStopWarn     bool
	LogStopError    bool
	LogStopFatal    bool
	LogStopSequence bool
}

// MQTT接続情報
type MqttConf struct {
	TCP  string
	Port int
}

// リクエスト情報
type ReqInfo struct {
	ProcessID string
	PcId      string
}

// プログラム情報
type ProInfo struct {
	SupplyReceipt         bool //補充レシート発行有無(true:発行、false：発行無)
	SalesCompleteCount    int  //SafeInfoManager 売上回収回数
	CollectCount          int  //SafeInfoManager 回収操作回数
	MaintenanceModeStatus int  //保守業務モードステータス
}

// OverFlowBoxType オーバーフローボックス(回収庫)情報
// True: 有り false: 無し
type OverFlowBoxType struct {
	CoinOverFlowBoxType bool // 硬貨オーバーフローボックス存在有無
	BillOverFlowBoxType bool // 紙幣オーバーフローボックス存在有無
}

// 不足エラー枚数
type ErrorNothing struct {
	Yen []int //枚数配列
}

// 不足注意枚数
type AlertNothing struct {
	Yen []int //枚数配列
}

// あふれエラー枚数
type ErrorMany struct {
	Yen      []int //枚数配列
	AllMoney int   //全硬貨
}

// あふれ注意枚数
type AlertMany struct {
	Yen      []int //枚数配列
	AllMoney int   //全硬貨
}

// リミット条件チェックパターン:リミット条件値と枚数の比較方法を変更する
type LimitCheck struct {
	Pattern int //リミット条件チェックパターン
}

// 枚数制限情報
type LimitInfo struct {
	ErrorNothing ErrorNothing //不足エラー枚数
	AlertNothing AlertNothing //不足注意枚数
	ErrorMany    ErrorMany    //あふれエラー枚数
	AlertMany    AlertMany    //あふれ注意枚数
	LimitCheck   LimitCheck   //リミット条件チェックパターン
}

// コンフィグ情報
type Configuration struct {
	MqttConf        MqttConf
	SystemConf      SystemConf
	ReqInfo         ReqInfo
	ProInfo         ProInfo
	MoneySetting    domain.MoneySetting  //金銭設定
	AggregateData   domain.AggregateData //レポート用金庫情報
	TermNo          int                  //精算端末番号
	OverFlowBoxType OverFlowBoxType
}

var Config Configuration

// iniの設定値を読む
func Initialize(moduleName string) Configuration {
	ipAddrTbl, _ := pc.GetLocalIpAddrInfo()
	// 設定ファイルフォルダ取得
	dirPath := file.GetCurrentDir()
	// env name
	env := file.GetEnv("ALMEXPATH")
	if len(env) != 0 {
		env = file.AdjustFileName(env)
		if file.DirExists(env + "ini") {
			dirPath = env + "ini"
		}
	}
	dirPath = file.AdjustFileName(dirPath)

	filename := fmt.Sprintf("%v%v.ini", dirPath, moduleName)
	if fileExists := file.FileExists(filename); fileExists {
		cfg, err := ini.Load(filename)
		if err == nil {
			// MQTT接続情報の取得
			Config.MqttConf, _ = getMqttInfo(cfg, "MQTT")
			// SYSTEM情報
			Config.SystemConf, _ = getSystemInfo(cfg, "SYSTEM")
			//プロセスID取得
			Config.ReqInfo.ProcessID = fmt.Sprintf("%08x", os.Getpid())
			//対象IPアドレス
			ipaddrtbl, _ := getIpAddrList(cfg, "PROGRAM", ipAddrTbl)
			Config.ReqInfo.PcId = ipaddrtbl[0]
			//プログラム各種設定
			Config.ProInfo, _ = getProgramInfo(cfg, "PROGRAM")
			// 集計データ
			Config.AggregateData = getAggregateData(cfg)
			// 金銭設定
			Config.MoneySetting = getMoneySetting(cfg)
			// 回収庫情報
			Config.OverFlowBoxType = getOverFlowBoxType(cfg, "ORVERFLOWBOXTYPE")
		}
	}
	// 精算端末番号取得（端末番号Id+1が端末番号となる）
	Config.TermNo = getTermId(dirPath) + 1

	return Config
}

// 対象IPアドレス取得
func getIpAddrList(cfg *ini.File, appName string, localIpAddrTbl []string) ([]string, bool) {
	ipaddr := cfg.Section(appName).Key("IpAddrTbl").String()
	ipaddrtbl := []string{}
	if len(ipaddr) != 0 {
		array := strings.Split(ipaddr, ",")
		for i := 0; i < len(array); i++ {
			if len(array[i]) != 0 {
				ipaddrtbl = append(ipaddrtbl, array[i])
			}
		}
	} else {
		// 未設定の場合はローカルIPすべて
		for i := 0; i < len(localIpAddrTbl); i++ {
			ipaddrtbl = append(ipaddrtbl, localIpAddrTbl[i])
		}
	}
	return ipaddrtbl, true
}

// MQTT接続情報
func getMqttInfo(cfg *ini.File, appName string) (MqttConf, bool) {
	conf := MqttConf{}

	tcp := cfg.Section(appName).Key("Server").String()
	if len(tcp) == 0 {
		tcp = "localhost"
	}
	port, err := cfg.Section(appName).Key("Port").Int()
	if err != nil {
		port = 1883
	}

	conf.TCP = tcp
	conf.Port = port
	return conf, true
}

// システム情報取得
func getSystemInfo(cfg *ini.File, appName string) (SystemConf, bool) {
	conf := SystemConf{}

	// localLog(アプリログ) 設定
	maxLength, err := cfg.Section(appName).Key("MaxLength").Int()
	if nil != err || 0 >= maxLength {
		maxLength = 4194304
	}
	maxRotation, err := cfg.Section(appName).Key("MaxRotation").Int()
	if nil != err || 0 > maxRotation {
		maxRotation = 16
	}

	// 起動時サービスステータスタイプ
	startUpStatus, err := cfg.Section(appName).Key("StartupStatus").Int()
	if err != nil {
		startUpStatus = 0
	}

	//ローカルログ設定
	logStopInfo, err := cfg.Section(appName).Key("LogStopInfo").Bool()
	if err != nil {
		logStopInfo = false
	}
	logStopTrace, err := cfg.Section(appName).Key("LogStopTrace").Bool()
	if err != nil {
		logStopTrace = false
	}
	logStopMqtt, err := cfg.Section(appName).Key("LogStopMqtt").Bool()
	if err != nil {
		logStopMqtt = false
	}
	logStopDebug, err := cfg.Section(appName).Key("LogStopDebug").Bool()
	if err != nil {
		logStopDebug = false
	}
	logStopMutex, err := cfg.Section(appName).Key("LogStopMutex").Bool()
	if err != nil {
		logStopMutex = false
	}
	logStopWarn, err := cfg.Section(appName).Key("LogStopWarn").Bool()
	if err != nil {
		logStopWarn = false
	}
	logStopError, err := cfg.Section(appName).Key("LogStopError").Bool()
	if err != nil {
		logStopError = false
	}
	logStopFatal, err := cfg.Section(appName).Key("LogStopFatal").Bool()
	if err != nil {
		logStopFatal = false
	}

	conf.StartUpStatus = startUpStatus
	conf.MaxLength = maxLength
	conf.MaxRotation = maxRotation
	conf.LogStopInfo = logStopInfo
	conf.LogStopTrace = logStopTrace
	conf.LogStopMqtt = logStopMqtt
	conf.LogStopDebug = logStopDebug
	conf.LogStopMutex = logStopMutex
	conf.LogStopWarn = logStopWarn
	conf.LogStopError = logStopError
	conf.LogStopFatal = logStopFatal
	return conf, true
}

// プログラム設定情報
func getProgramInfo(cfg *ini.File, appName string) (ProInfo, bool) {
	conf := ProInfo{}

	conf.SupplyReceipt, _ = cfg.Section(appName).Key("Suplly_Recipt").Bool()
	conf.SalesCompleteCount, _ = cfg.Section(appName).Key("SalesCompleteCount").Int()
	conf.CollectCount, _ = cfg.Section(appName).Key("CollectCount").Int()
	conf.MaintenanceModeStatus, _ = cfg.Section(appName).Key("MaintenanceModeStatus").Int()

	return conf, true
}

// tex_controller.iniから端末IDを取得
func getTermId(dirPath string) (termId int) {
	// 設定ファイルフォルダ取得
	filename := fmt.Sprintf("%v%v.ini", dirPath, "tex_controller")

	//ファイル読込
	if fileExists := file.FileExists(filename); !fileExists {
		return
	}
	cfg, err := ini.Load(filename)
	if err != nil {
		return
	}

	// 端末ID取得
	termId, _ = cfg.Section("TERMINFO").Key("TermId").Int()
	return
}

func getAggregateData(cfg *ini.File) domain.AggregateData {
	aggregateData := domain.AggregateData{}
	//デフォルト値をセット
	for _, v := range domain.MODE {
		aggregateData[v] = &domain.AggregateSafeInfo{}
	}

	for _, mode := range domain.MODE {
		// セクション名作成
		section := fmt.Sprintf("AggregateData_%d", mode)

		//処理前有高金種配列
		beforeAmountCountTbl := cfg.Section(section).Key("BeforeCountTbl").String()
		for i, v := range changeStringToArray(beforeAmountCountTbl) {
			aggregateData[mode].BeforeAmountCountTbl[i] = v
		}
		//処理前補充入金金種配列
		beforeReplenishCountTbl := cfg.Section(section).Key("BeforeReplenishCountTbl").String()
		for i, v := range changeStringToArray(beforeReplenishCountTbl) {
			aggregateData[mode].BeforeReplenishCountTbl[i] = v
		}
		//補充入金金種配列
		replenishCountTbl := cfg.Section(section).Key("RepelenishCountTbl").String()
		for i, v := range changeStringToArray(replenishCountTbl) {
			aggregateData[mode].ReplenishCountTbl[i] = v
		}
		//処理前回収金種配列
		beforeCollectCountTbl := cfg.Section(section).Key("BeforeCollectCountTbl").String()
		for i, v := range changeStringToArray(beforeCollectCountTbl) {
			aggregateData[mode].BeforeCollectCountTbl[i] = v
		}
		//回収金種配列
		collectCountTbl := cfg.Section(section).Key("CollectCountTbl").String()
		for i, v := range changeStringToArray(collectCountTbl) {
			aggregateData[mode].CollectCountTbl[i] = v
		}
		//処理後有高金種配列
		afterAmountCountTbl := cfg.Section(section).Key("AfterAmountCountTbl").String()
		for i, v := range changeStringToArray(afterAmountCountTbl) {
			aggregateData[mode].AfterAmountCountTbl[i] = v
		}
		//売上回収金種配列
		salesCollectCountTbl := cfg.Section(section).Key("SalesCollectCountTbl").String()
		for i, v := range changeStringToArray(salesCollectCountTbl) {
			aggregateData[mode].SalesCollectCountTbl[i] = v
		}
	}
	return aggregateData
}

// 文字列を配列にセット
func changeStringToArray(stgTbl string) (intTbl []int) {

	// 文字数チェック
	if len(stgTbl) < 2 {
		return
	}

	// 数値型の配列にセット
	slice := strings.Split(stgTbl[1:len(stgTbl)-1], ",") //[]と,を削除
	for _, v := range slice {
		n, _ := strconv.Atoi(v)
		intTbl = append(intTbl, n)
	}
	return
}

func getMoneySetting(cfg *ini.File) domain.MoneySetting {
	moneySetting := domain.MoneySetting{}

	// 釣銭準備金枚数
	changeReserveCount := cfg.Section("MoneySetting").Key("ChangeReserveCount").String()
	_ = json.Unmarshal([]byte(changeReserveCount), &moneySetting.ChangeReserveCount)

	// 不足枚数
	changeShortageCount := cfg.Section("MoneySetting").Key("ChangeShortageCount").String()
	_ = json.Unmarshal([]byte(changeShortageCount), &moneySetting.ChangeShortageCount)
	if changeShortageCount == "" { // 初回登録時に中身がからの場合、空データを生成
		moneySetting.ChangeShortageCount.RegisterDataTbl = [2]domain.RegisterDataTbl{
			{AlertLevel: 1},
			{AlertLevel: 2},
		}
	}

	// あふれ枚数
	excessChangeCount := cfg.Section("MoneySetting").Key("ExcessChangeCount").String()
	_ = json.Unmarshal([]byte(excessChangeCount), &moneySetting.ExcessChangeCount)

	if changeShortageCount == "" { // 初回登録時に中身がからの場合、空データを生成
		moneySetting.ExcessChangeCount.ExRegisterDataTbl = [2]domain.ExRegisterDataTbl{
			{AlertLevel: 1},
			{AlertLevel: 2},
		}
	}
	return moneySetting
}

func getOverFlowBoxType(cfg *ini.File, sectionName string) OverFlowBoxType {

	coin, err := cfg.Section(sectionName).Key("CoinOverFlowBox").Bool()
	if err != nil {
		coin = true
	}
	bill, err := cfg.Section(sectionName).Key("BillOverFlowBox").Bool()
	if err != nil {
		bill = false
	}

	return OverFlowBoxType{
		CoinOverFlowBoxType: coin,
		BillOverFlowBoxType: bill,
	}
}
