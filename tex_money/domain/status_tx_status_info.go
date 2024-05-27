package domain

// 精算機状態管理：精算機状態取得要求
type RequestStatusStatusTx struct {
	RequestInfo RequestInfo `json:"requestInfo"`
}

type ResultStatusStatusTx struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Result      bool        `json:"result"`                //処理結果
	ErrorCode   string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string      `json:"errorDetail,omitempty"` //エラー詳細
	InfoTerm    InfoTerm    `json:"infoTerm"`              //端末情報
}

type InfoTerm struct {
	StatusError         int    `json:"statusError"`         //エラー状態
	TermErrorCode       int    `json:"termErrorCode"`       //エラーコード
	TermErrorState      int    `json:"termErrorState"`      //エラー発生状態
	StatusHandling      int    `json:"statusHandling"`      //取扱状態
	StatusSecurity      int    `json:"statusSecurity"`      //セキュリティ状態
	StatusDoor          int    `json:"statusDoor"`          //扉状態
	StatusKeySw         int    `json:"statuskeySw"`         //キーSW状態
	StatusCall          int    `json:"statusCall"`          //従業員呼出状態
	PaymentMode         int    `json:"paymentMode"`         //決済方法モード
	ErrorContent        string `json:"errorContent"`        //エラー内容
	ErrorDetailCode     string `json:"errorDetailCode"`     //エラーコード（詳細）
	ErrorDetailContent  string `json:"errorDetailContent"`  //エラー内容（詳細）
	RecoverGuideMessage string `json:"recoverGuideMessage"` //エラー復旧ガイドライン
	RecoverGuidePicture string `json:"recoverGuidePicture"` //エラー復旧案内画像
	RecoverGuideMovie   string `json:"recoverGuideMovie"`   //エラー復旧案内動画
}
