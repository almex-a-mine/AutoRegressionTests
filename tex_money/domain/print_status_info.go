package domain

// 印刷制御:印刷ステータス要求
type RequestPrintStatus struct {
	RequestInfo RequestInfo `json:"requestInfo"`
}

// 印刷ステータス
type PrintStatus struct {
	RequestInfo    RequestInfo      `json:"requestInfo"`
	Result         bool             `json:"result"`                //処理結果
	ErrorCode      string           `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail    string           `json:"errorDetail,omitempty"` //エラー詳細
	PrintId        string           `json:"printId"`               //印刷管理ID
	PrintType      int              `json:"printType"`             //印刷タイプ
	DeviceNo       string           `json:"deviceNo"`              //印刷デバイス番号
	MonitorId      string           `json:"monitorId"`             //モニタ管理ID
	PrinterName    string           `json:"printerName"`           //出力プリンタ名
	PrinterType    int              `json:"printerType"`           //出力プリンタ種別
	OutletNo       int              `json:"outletNo"`              //出力プリンタ取出口番号
	DocumentName   string           `json:"documentName"`          //ドキュメント名
	CountPlan      int              `json:"countPlan"`             //出力予定枚数
	CountEnd       int              `json:"countEnd"`              //印刷完了枚数
	StatusPrint    int              `json:"statusPrint"`           //印刷ステータス
	StatusResult   bool             `json:"statusResult"`          //印刷結果
	PrinterInfoTbl []PrinterInfoTbl `json:"printerInfoTbl"`        //印刷プリンタ詳細情報
}

// 印刷プリンタ詳細情報
type PrinterInfoTbl struct {
	DeviceNo     string `json:"deviceNo"`     //印刷デバイス番号
	MonitorId    string `json:"monitorId"`    //モニタ管理ID
	PrinterName  string `json:"printerName"`  //出力プリンタ名
	PrinterType  int    `json:"printerType"`  //出力プリンタ種別
	OutletNo     int    `json:"outletNo"`     //出力プリンタ取出口番号
	DocumentName string `json:"documentName"` //ドキュメント名
	CountPlan    int    `json:"countPlan"`    //出力予定枚数
	CountEnd     int    `json:"countEnd"`     //印刷完了枚数
	StatusPrint  int    `json:"statusPrint"`  //印刷ステータス
}
