package usecases

import (
	"fmt"
	"stub/domain"
	"stub/domain/handler"
	"sync"
	"time"
)

// 待機情報
type WaitInfo struct {
	ProcessId   string
	RequestId   string
	RequestData interface{}
	ResultData  interface{}
	WaitCh      chan bool
	waitTimer   *time.Timer
}

type waitManager struct {
	mutexWaitQues sync.Mutex
	waitInfoTbl   map[string]WaitInfo
	logger        handler.LoggerRepository
}

// 情報待機クラス
var pWaitManager *waitManager

// 情報待機クラス
func NewWaitManager(logger handler.LoggerRepository) IWait {
	if pWaitManager == nil {
		pWaitManager = &waitManager{
			waitInfoTbl: make(map[string]WaitInfo),
			logger:      logger}
	}
	return pWaitManager
}

// 待機情報を作成
//
//	paramater : プロセスID, リクエストID
//	return    : 待機情報
func (w *waitManager) MakeWaitInfo(texCon *domain.TexContext, processId, requestid string, requestdata interface{}) WaitInfo {
	//w.logger.Trace("【%v】START:waitManager MakeWaitInfo(processId=%s, requestid=%s, requestdata=%+v) WaitInfo ", texCon.GetUniqueKey(), processId, requestid, requestdata)
	w.mutexWaitQues.Lock()
	//w.logger.Debug("【%v】★★デバック用 waitManagerのlockの後", texCon.GetUniqueKey())
	defer w.mutexWaitQues.Unlock()

	//待機情報作成
	WaitInfo := WaitInfo{}
	WaitInfo.ProcessId = processId
	WaitInfo.RequestId = requestid
	WaitInfo.RequestData = requestdata
	WaitInfo.WaitCh = make(chan bool)

	// 作成した待機情報を管理テーブルへ追加
	id := w.makeWaitId(texCon, processId, requestid)
	w.waitInfoTbl[id] = WaitInfo
	//w.logger.Trace("【%v】waitManager MakeWaitInfo  WaitInfo%+v", texCon.GetUniqueKey(), WaitInfo)
	return WaitInfo
}

// 待機情報へセット
//
//	paramater : プロセスID, リクエストID, 更新データ
//	return    : 待機情報
func (w *waitManager) SetWaitInfo(texCon *domain.TexContext, processId, requestid string, data interface{}) bool {
	w.mutexWaitQues.Lock()
	defer w.mutexWaitQues.Unlock()

	id := w.makeWaitId(texCon, processId, requestid)

	//待機IDが管理テーブルに存在するか？
	if waitInfo, ok := w.waitInfoTbl[id]; ok {
		// 存在する場合、該当待機情報をデータをセットし、通知を行う
		waitData := w.waitInfoTbl[id]
		waitData.ResultData = data
		w.waitInfoTbl[id] = waitData
		waitInfo.WaitCh <- true

		return true
	}
	w.logger.Debug("【%v】待機情報重複  待機ID=%v", texCon.GetUniqueKey(), id)
	return false

}

// 待機処理
//
//	paramater : 待機情報, 待機時間
//	return    : 処理結果(true:OK false:タイムアウト)
func (w *waitManager) WaitResultInfo(texCon *domain.TexContext, wInfo WaitInfo, timerValue int) bool {
	//w.logger.Trace("【%v】START:waitManager WaitResultInfo(wInfo=%+v, timerValue =%v)", texCon.GetUniqueKey(), wInfo, timerValue)
	wInfo.waitTimer = time.NewTimer(time.Duration(timerValue) * time.Millisecond)
	wairRet := false

	if wInfo.waitTimer != nil {
		select {
		case <-wInfo.WaitCh:
			wairRet = true
		case <-wInfo.waitTimer.C:
			wairRet = false
		}

		if wInfo.waitTimer != nil {
			wInfo.waitTimer.Stop()
		}
		wInfo.waitTimer = nil
	}
	// w.logger.Trace("【%v】END:waitManager WaitResultInfo result=%t", texCon.GetUniqueKey(), wairRet)
	return wairRet
}

// 待機情報削除(不要になった待機情報を削除)
func (w *waitManager) DelWaitInfo(texCon *domain.TexContext, processId, requestid string) {

	w.mutexWaitQues.Lock()
	defer w.mutexWaitQues.Unlock()

	id := w.makeWaitId(texCon, processId, requestid)

	// 存在する場合、該当待機情報をデータを削除
	_, ok := w.waitInfoTbl[id]
	if ok {
		delete(w.waitInfoTbl, id)
	}

}

// 待機情報を取得
//
//	paramater : プロセスID, リクエストID
//	return    : 待機情報
func (w *waitManager) GetWaitInfo(texCon *domain.TexContext, processId, requestid string) (WaitInfo, bool) {

	w.mutexWaitQues.Lock()
	defer w.mutexWaitQues.Unlock()

	id := w.makeWaitId(texCon, processId, requestid)

	waitInfo, ok := w.waitInfoTbl[id]

	return waitInfo, ok
}

// 待機用IDの生成
//
//	paramater : プロセスID, リクエストID
//	return    : 待機用ID
func (w *waitManager) makeWaitId(texCon *domain.TexContext, processId, requestid string) string {
	ret := fmt.Sprintf("%v-%v", processId, requestid)
	return ret
}

func GetWaitInfoMqttRecv(processId, requestid string) bool {

	pWaitManager.mutexWaitQues.Lock()
	defer pWaitManager.mutexWaitQues.Unlock()

	id := fmt.Sprintf("%v-%v", processId, requestid)

	_, ok := pWaitManager.waitInfoTbl[id]

	return ok
}
