package domain

// 現金入出金機制御:入出金機ステータス取得要求
type RequestStatus struct {
	RequestInfo RequestInfo `json:"requestInfo"`
}

type ResultStatus struct {
	RequestInfo         RequestInfo      `json:"requestInfo"`
	Result              bool             `json:"result"`                //処理結果
	ErrorCode           string           `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail         string           `json:"errorDetail,omitempty"` //エラー詳細
	CoinStatusCode      int              `json:"coinStatusCode"`        //硬貨結果通知コード
	BillStatusCode      int              `json:"billStatusCode"`        //紙幣結果通知コード
	CoinStatusTbl       StatusTbl        `json:"coinStatusTbl"`         //硬貨ステータス情報
	CoinResidueInfoTbl  []ResidueInfoTbl `json:"coinResidueInfoTbl"`    //硬貨残留情報
	BillStatusTbl       StatusTbl        `json:"billStatusTbl"`         //紙幣ステータス情報
	BillResidueInfoTbl  []ResidueInfoTbl `json:"billResidueInfoTbl"`    //紙幣残留情報
	DeviceStatusInfoTbl []string         `json:"deviceStatusInfoTbl"`   //デバイス詳細情報
	WarningInfoTbl      []int            `json:"warningInfoTbl"`        //警告情報
}

// ステータス情報
type StatusTbl struct {
	StatusCover       bool `json:"statusCover"`       //トビラ状態
	StatusUnitSet     bool `json:"statusUnitSet"`     //ユニットセット状態
	StatusInCassette  bool `json:"statusInCassette"`  //補充カセット状態
	StatusOutCassette bool `json:"statusOutCassette"` //回収カセット状態
	StatusInsert      bool `json:"statusInsert"`      //入金口状態
	StatusExit        bool `json:"statusExit"`        //出金口状態
	StatusRjbox       bool `json:"statusRjbox"`       //リジェクトBOX状態
}

// 残留情報
type ResidueInfoTbl struct {
	Title  string `json:"title"`  //管理名称
	Status bool   `json:"status"` //状態
}
