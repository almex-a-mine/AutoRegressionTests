package domain

type RequestGetSalesinfo struct {
	RequestInfo RequestInfo `json:"requestInfo"`
}

type ResultGetSalesinfo struct {
	RequestInfo RequestInfo  `json:"requestInfo"`
	Result      bool         `json:"result"`                //処理結果
	ErrorCode   string       `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string       `json:"errorDetail,omitempty"` //エラー詳細
	InfoSales   InfoSalesTbl `json:"infoSales"`
}
