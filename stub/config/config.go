package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"stub/pkg/file"
	"stub/pkg/pc"

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

// コンフィグ情報
type Configuration struct {
	MqttConf   MqttConf
	SystemConf SystemConf
	ReqInfo    ReqInfo
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
		}
	}

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
