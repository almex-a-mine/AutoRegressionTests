package domain

// 稼働データ管理:金庫情報遷移記録要求
type RequestReportSafeInfo struct {
	RequestInfo     RequestInfo `json:"requestInfo"`
	GenerateDate    int         `json:"generateDate"`    //発生日付
	GenerateTime    int         `json:"generateTime"`    //発生時刻
	HistorySortCode int         `json:"historySortCode"` //履歴分類コード
	StatusDetail    string      `json:"statusDetail"`    //状態遷移内容
	InfoSafe        InfoSafeTbl `json:"infoSafe"`        //金庫情報
}

type ResultReportSafeInfo struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Result      bool        `json:"result"`                //処理結果
	ErrorCode   string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string      `json:"errorDetail,omitempty"` //エラー詳細
}

type InfoSafeTbl struct {
	CurrentStatusTbl [CASH_TYPE_SHITEI]int `json:"currentStatusTbl"` //通常金種別状況
	SortInfoTbl      []SortInfoTbl         `json:"sortInfotbl"`      //分類情報
}

func NewRequestReportSafeInfo(info RequestInfo, historySortCode int, infoSafe InfoSafeTbl) *RequestReportSafeInfo {
	return &RequestReportSafeInfo{
		RequestInfo:     info,
		HistorySortCode: historySortCode,
		InfoSafe:        infoSafe,
	}
}
