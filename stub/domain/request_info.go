package domain

type RequestInfo struct {
	ProcessID string `json:"processId"`
	PcId      string `json:"pcId"`
	RequestID string `json:"requestId"`
}

type RequestInfoOnly struct {
	RequestInfo RequestInfo `json:"requestInfo"`
}
