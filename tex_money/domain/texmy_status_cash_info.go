package domain

// 現金入出金機制御ステータス情報
type RequestStatusCash struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	CashControlId string      `json:"cashControlId"` //入出金機制御管理番号
}

type ResultStatusCash struct {
	RequestInfo         RequestInfo        `json:"requestInfo,omitempty"`
	Result              bool               `json:"result"`                //通信結果
	CashControlId       string             `json:"cashControlId"`         //入出金機制御管理番号
	StatusReady         bool               `json:"statusReady"`           //制御状態
	StatusMode          int                `json:"statusMode"`            //動作状態
	StatusLine          bool               `json:"statusLine"`            //通信状態
	StatusError         bool               `json:"statusError"`           //エラー状態
	ErrorCode           string             `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail         string             `json:"errorDetail,omitempty"` //エラー詳細
	StatusCover         bool               `json:"statusCover"`           //トビラ状態
	StatusAction        int                `json:"statusAction"`          //動作状態
	StatusInsert        bool               `json:"statusInsert"`          //入金口状態
	StatusExit          bool               `json:"statusExit"`            //出金口状態
	StatusRjbox         bool               `json:"statusRjbox"`           //リジェクトBOX
	BillStatusTbl       TexmyBillStatusTbl `json:"billStatusTbl"`         //紙幣ステータス情報
	CoinStatusTbl       CoinStatusTbl      `json:"coinStatusTbl"`         //硬貨ステータス情報
	BillResidueInfoTbl  []BillResidueInfo  `json:"billResidueInfoTbl"`    //紙幣残留情報
	CoinResidueInfoTbl  []CoinResidueInfo  `json:"coinResidueInfoTbl"`    //硬貨残留情報
	DeviceStatusInfoTbl []string           `json:"deviceStatusInfoTbl"`   //デバイス詳細情報
	WarningInfoTbl      []int              `json:"warningInfoTbl"`        //警告情報
}

type StatusCash struct {
	CashControlId       string             `json:"cashControlId"`         //入出金機制御管理番号
	StatusReady         bool               `json:"statusReady"`           //制御状態
	StatusMode          int                `json:"statusMode"`            //動作状態
	StatusLine          bool               `json:"statusLine"`            //通信状態
	StatusError         bool               `json:"statusError"`           //エラー状態
	ErrorCode           string             `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail         string             `json:"errorDetail,omitempty"` //エラー詳細
	StatusCover         bool               `json:"statusCover"`           //トビラ状態
	StatusAction        int                `json:"statusAction"`          //動作状態
	StatusInsert        bool               `json:"statusInsert"`          //入金口状態
	StatusExit          bool               `json:"statusExit"`            //出金口状態
	StatusRjbox         bool               `json:"statusRjbox"`           //リジェクトBOX
	BillStatusTbl       TexmyBillStatusTbl `json:"billStatusTbl"`         //紙幣ステータス情報
	CoinStatusTbl       CoinStatusTbl      `json:"coinStatusTbl"`         //硬貨ステータス情報
	BillResidueInfoTbl  []BillResidueInfo  `json:"billResidueInfoTbl"`    //紙幣残留情報
	CoinResidueInfoTbl  []CoinResidueInfo  `json:"coinResidueInfoTbl"`    //硬貨残留情報
	DeviceStatusInfoTbl []string           `json:"deviceStatusInfoTbl"`   //デバイス詳細情報
	WarningInfoTbl      []int              `json:"warningInfoTbl"`        //警告情報
}

// 紙幣ステータス情報
type TexmyBillStatusTbl struct {
	StatusUnitSet     bool `json:"statusUnitSet"`               //ユニットセット状態
	StatusInCassette  bool `json:"statusInCassette"`            //補充カセット状態
	StatusOutCassette bool `json:"statusOutCassette"`           //回収カセット状態
	StatusAmountCount int  `json:"statusAmountCount,omitempty"` //有高枚数状態
}

// 硬貨ステータス情報
type CoinStatusTbl struct {
	StatusUnitSet     bool `json:"statusUnitSet"`               //ユニットセット状態
	StatusInCassette  bool `json:"statusInCassette"`            //補充カセット状態
	StatusOutCassette bool `json:"statusOutCassette"`           //回収カセット状態
	StatusAmountCount int  `json:"statusAmountCount,omitempty"` //有高枚数状態
}

// 紙幣残留情報
type BillResidueInfo struct {
	Title  string `json:"title"`  //管理名称
	Status bool   `json:"status"` //状態
}

// 硬貨残留情報
type CoinResidueInfo struct {
	Title  string `json:"title"`  //管理名称
	Status bool   `json:"status"` //状態
}
