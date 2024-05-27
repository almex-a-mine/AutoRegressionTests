package domain

// 入出金レポート印刷ステータス情報
type StatusReport struct {
	SlipPrintId  string `json:"slipPrintId"`            //レポート印刷制御管理番号
	StatusPrint  int    `json:"statusPrint"`            //印刷状態
	CountPlan    int    `json:"countPlan"`              //出力予定枚数
	CountEnd     int    `json:"countEnd"`               //印刷完了枚数
	StatusResult *bool  `json:"statusResult,omitempty"` //印刷結果
}
