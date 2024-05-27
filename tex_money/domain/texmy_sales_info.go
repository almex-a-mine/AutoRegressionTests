package domain

// 売上金情報要求
type RequestSalesInfo struct {
	RequestInfo RequestInfo `json:"requestInfo"`
}

type ResultSalesInfo struct {
	RequestInfo   RequestInfo `json:"requestInfo"`
	Result        bool        `json:"result"`                //処理結果
	ErrorCode     string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail   string      `json:"errorDetail,omitempty"` //エラー詳細
	SalesAmount   int         `json:"salesAmount"`           //売上金額
	SalesComplete int         `json:"salesComplete"`         //売上金回収済金額
	SalesCount    int         `json:"salesCount"`            //売上金回収回数
}
