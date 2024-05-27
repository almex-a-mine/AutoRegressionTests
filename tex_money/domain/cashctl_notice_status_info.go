package domain

// 現金入出金制御：入出金機ステータス通知
type NoticeStatus struct {
	CoinStatusCode           int              `json:"coinStatusCode"`      //硬貨結果通知コード
	BillStatusCode           int              `json:"billStatusCode"`      //紙幣結果通知コード
	CoinNoticeStatusTbl      NoticeStatusTbl  `json:"coinStatusTbl"`       //硬貨ステータス情報
	NoticeCoinResidueInfoTbl []ResidueInfoTbl `json:"coinResidueInfoTbl"`  //硬貨残留情報
	BillNoticeStatusTbl      NoticeStatusTbl  `json:"billStatusTbl"`       //紙幣ステータス情報
	BillNoticeResidueInfoTbl []ResidueInfoTbl `json:"billResidueInfoTbl"`  //紙幣残留情報
	DeviceStatusInfoTbl      []string         `json:"deviceStatusInfoTbl"` //デバイス詳細情報
	WarningInfoTbl           []int            `json:"warningInfoTbl"`      //警告情報
}

// ステータス情報
type NoticeStatusTbl struct {
	StatusError       bool   `json:"statusError"`           //処理結果
	ErrorCode         string `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail       string `json:"errorDetail,omitempty"` //エラー詳細
	StatusCover       bool   `json:"statusCover"`           //トビラ状態
	StatusUnitSet     bool   `json:"statusUnitSet"`         //ユニットセット状態
	StatusInCassette  bool   `json:"statusInCassette"`      //補充カセット状態
	StatusOutCassette bool   `json:"statusOutCassette"`     //回収カセット状態
	StatusInsert      bool   `json:"statusInsert"`          //入金口状態
	StatusExit        bool   `json:"statusExit"`            //出金口状態
	StatusRjbox       bool   `json:"statusRjbox"`           //リジェクトBOX状態
}
