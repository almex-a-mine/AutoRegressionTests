package usecases

type ErrorManager interface {
	GetErrorInfo(errorTypeCode int) (string, string)
}
