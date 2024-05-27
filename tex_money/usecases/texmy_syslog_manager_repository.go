package usecases

type SyslogManager interface {
	NoticeSystemLog(logType int, deviceNo string, logData string)
}
