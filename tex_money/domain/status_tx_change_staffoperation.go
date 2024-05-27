package domain

// 精算機状態管理：状態変更要求(スタッフ操作記録)
type RequestChangeStaffOperation struct {
	RequestInfo     RequestInfo `json:"requestInfo"`
	OperationType   int         `json:"operationType"`             //操作種別
	OperationDetail string      `json:"operationDetail,omitempty"` //操作内容
}

type ResultChangeStaffOperation struct {
	RequestInfo RequestInfo `json:"requestInfo"`
	Result      bool        `json:"result"`                //処理結果
	ErrorCode   string      `json:"errorCode,omitempty"`   //エラーコード
	ErrorDetail string      `json:"errorDetail,omitempty"` //エラー詳細
}

func NewRequestChangeStaffOperation(info RequestInfo, operationType int, operationDetail string) *RequestChangeStaffOperation {
	return &RequestChangeStaffOperation{
		RequestInfo:     info,
		OperationType:   operationType,
		OperationDetail: operationDetail,
	}
}
