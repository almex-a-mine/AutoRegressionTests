package infrastructure

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"stub/domain"
	"stub/domain/handler"
	"sync"
	"time"
)

type logger struct {
	FileName    string
	FilePath    string
	File        *os.File
	MaxLength   int
	MaxRotation int
	Mtx         sync.Mutex
	Size        int64
	// 動作確認モード実装で追加
	VerifyMode int
	VerifyName string
	VerifyPath string
	//ログ設定値
	LogStopInfo  bool
	LogStopTrace bool
	LogStopMqtt  bool
	LogStopDebug bool
	LogStopMutex bool
	LogStopWarn  bool
	LogStopError bool
	LogStopFatal bool
}

var gLogger *logger

func NewLogger(maxLength int,
	maxRotation int,
	logStopInfo bool,
	logStopTrace bool,
	logStopMqtt bool,
	logStopDebug bool,
	logStopMutex bool,
	logStopWarn bool,
	logStopError bool,
	logStopFatal bool) handler.LoggerRepository {
	gLogger = &logger{}

	// ログファイル名作成
	almexpath := os.Getenv("ALMEXPATH")
	if len(almexpath) < 1 {
		almexpath = "."
	}
	_ = os.MkdirAll(almexpath+"/log", 0666)
	gLogger.FilePath = fmt.Sprintf("%v/log", almexpath)
	gLogger.FileName = fmt.Sprintf("%v/%s.log", gLogger.FilePath, domain.SrvName)

	gLogger.VerifyPath = fmt.Sprintf("%v/verify", almexpath)
	gLogger.VerifyName = fmt.Sprintf("%v/%s.log", gLogger.VerifyPath, domain.SrvName)

	gLogger.MaxLength = maxLength
	gLogger.MaxRotation = maxRotation

	//ログの設定値セット
	gLogger.LogStopInfo = logStopInfo
	gLogger.LogStopTrace = logStopTrace
	gLogger.LogStopMqtt = logStopMqtt
	gLogger.LogStopDebug = logStopDebug
	gLogger.LogStopMutex = logStopMutex
	gLogger.LogStopWarn = logStopWarn
	gLogger.LogStopError = logStopError
	gLogger.LogStopFatal = logStopFatal

	return gLogger
}

func (l *logger) Open() bool {
	_, fileName := l.GetFileInfo()
	var err error
	l.File, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("local log Open Error %s", err)
		return false
	}
	log.SetOutput(io.MultiWriter(l.File, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	return true
}

func (l *logger) Close() bool {
	err := l.File.Close()
	if err != nil {
		log.Printf("local log Close Error %s", err)
		return false
	}
	return true
}

func (l *logger) SetMaxLength(length int) {
	l.MaxLength = length
}

func (l *logger) SetMaxRotation(rotation int) {
	l.MaxRotation = rotation
	// 最大ローテーションを取得できた時点で、最大数より大きいログを削除する
	l.Mtx.Lock()
	defer l.Mtx.Unlock()
	l.DeleteLog()
}

func (l *logger) SetSize() {
	fileinfo, err := l.File.Stat()
	if nil != err {
		log.Printf("Get File Size Error:%s  l:%v  fileinfo:%v", err.Error(), l, fileinfo)
		return
	}
	l.Size = fileinfo.Size()
}

// ログファイルのローテーション
func (l *logger) CheckRotation() {
	// config読込前のエラーは必ず出力
	if -1 == l.MaxLength && -1 == l.MaxRotation {
		return
	}

	// 書き込み中のファイルが最大サイズ以上かチェック
	if l.Size < int64(l.MaxLength) {
		return
	}

	filePath, fileName := l.GetFileInfo()

	var fileCount int
	// ログローテーション数の指定が0(無制限)
	if 0 == l.MaxRotation {

		// ディレクトリ内のファイルリストを取得
		fileinfos, _ := os.ReadDir(filePath)

		// 現在のローテーションログファイルがいくつ存在するか数える
		for _, fileinfo := range fileinfos {

			// パスを結合
			checkname := filePath + "/" + fileinfo.Name()

			// 出力ファイル名(syslog_システム名.log)と同じならcontinue
			if checkname == fileName {
				continue
			}

			// 出力ファイル名を含んでいる場合はインクリメント
			if 0 == strings.Index(checkname, fileName) {
				fileCount++
			}
		}
	} else { // ログローテーション数が0以外(マイナスはありえない、.ini読込時にdefault値に書換えるため
		fileCount = l.MaxRotation - 1
	}

	// ログファイルのリネーム 4->5, 3->4, 2->3 と逆順にrenameする
	for i := fileCount; i > 0; i-- {
		// rename後のファイルが存在する場合は先に削除しておく
		if _, err := os.Stat(fileName + strconv.FormatInt(int64(i), 10)); nil == err {
			_ = os.Remove(fileName + strconv.FormatInt(int64(i), 10))
		}
		_ = os.Rename(fileName+strconv.FormatInt(int64(i-1), 10), fileName+strconv.FormatInt(int64(i), 10))
	}

	// 現行のログファイル名に0を付与 .log -> .log0
	_ = os.Rename(fileName, fileName+"0")
}

// ログファイル削除 ログローテーション最大数以上のログを削除する
func (l *logger) DeleteLog() {

	filePath, fileName := l.GetFileInfo()

	// ディレクトリ内のファイルリストを取得
	fileinfos, _ := os.ReadDir(filePath)

	// ログローテーション数以上のファイルを削除
	for _, fileinfo := range fileinfos {

		// pathを結合してFullPathを作成
		checkname := filePath + "/" + fileinfo.Name()

		// 出力ファイル名(syslog_システム名.log)と同じならcontinue
		if checkname == fileName {
			continue
		}

		// 出力ファイル名を含んでいるか
		if 0 == strings.Index(checkname, fileName) {
			flg := false
			// ローテーションファイル最大数の範囲内のファイル名かチェック
			for i := 0; i < l.MaxRotation; i++ {
				if (checkname) == (fileName + strconv.FormatInt(int64(i), 10)) {
					flg = true
					break
				}
			}

			if !flg {
				_ = os.Remove(checkname) // ファイル削除
				log.Printf("file deleted.:%s", checkname)
			}
		}
	}
}

// ローカルログ書き込み
func (l *logger) WriteLocalLog(stop bool, lv string, format *string, v ...interface{}) {
	if stop {
		return
	}

	l.Mtx.Lock()
	defer l.Mtx.Unlock()

	// ログを開く前にローテーション
	l.CheckRotation()

	// ログローテーション数以上のログを削除
	if 0 < l.MaxRotation {
		l.DeleteLog()
	}

	l.Open()        // ログを開く
	defer l.Close() // ログを閉じる deferはLIFOなので Close->Unlock の順に実行される
	log.SetFlags(0)
	log.SetPrefix(lv)

	t := getLogTime()
	text := t + fmt.Sprintf(*format, v...)

	log.Print(text)
	l.SetSize() // ログサイズを更新する

}

func getLogTime() string {
	t := time.Now()
	return fmt.Sprintf("%s.%06v ", t.Format("2006/01/02 15:04:05"), t.Nanosecond()/1e3)
}

func (l *logger) SetSystemOperation(mode int) {
	l.VerifyMode = mode
}

func (l *logger) GetSystemOperation() int {
	return l.VerifyMode
}

// filePathとfileNameを返却
func (l *logger) GetFileInfo() (string, string) {
	if 1 == l.VerifyMode {
		return l.VerifyPath, l.VerifyName
	}
	return l.FilePath, l.FileName
}

// 以降LVごとのログ出力method
func (l *logger) Debug(format string, v ...interface{}) {
	l.WriteLocalLog(l.LogStopDebug, "[DEBUG] ", &format, v...)
}

func (l *logger) Mutex(format string, v ...interface{}) {
	l.WriteLocalLog(l.LogStopMutex, "[MUTEX] ", &format, v...)
}

func (l *logger) Trace(format string, v ...interface{}) {
	l.WriteLocalLog(l.LogStopTrace, "[TRACE] ", &format, v...)
}

func (l *logger) Info(format string, v ...interface{}) {
	l.WriteLocalLog(l.LogStopInfo, "[INFO ] ", &format, v...)
}

func (l *logger) Warn(format string, v ...interface{}) {
	l.WriteLocalLog(l.LogStopWarn, "[WARN ] ", &format, v...)
}

func (l *logger) Error(format string, v ...interface{}) {
	l.WriteLocalLog(l.LogStopError, "[ERROR] ", &format, v...)
}

func (l *logger) Fatal(format string, v ...interface{}) {
	l.WriteLocalLog(l.LogStopFatal, "[FATAL] ", &format, v...)
}
