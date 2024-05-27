package domain

// 入出金レポート印刷情報
type RequestPrintReport struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	FilePath    string      `json:"filePath"` //精算結果ファイルパス名
	ReportId    int         `json:"reportId"` //レポート管理番号
}

type ResultPrintReport struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Result      bool        `json:"result"`                //処理結果
	ErrorCode   string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string      `json:"errorDetail,omitempty"` //エラー詳細
	SlipPrintId string      `json:"slipPrintId"`           //レポート印刷制御管理番号
}
