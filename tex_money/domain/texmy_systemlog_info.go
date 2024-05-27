package domain

// シスログ情報
type NoticeSystemLogInfo struct {
	GenerateDate int    `json:"generateDate,omitempty"`
	GenerateTime int    `json:"generateTime,omitempty"`
	ForceEncrypt bool   `json:"forceEncrypt"`
	LogLevel     string `json:"logLevel"`
	RequestId    string `json:"requestId"`
	LogSummary   string `json:"logSummary"`
	LogData      string `json:"logData"`
}
