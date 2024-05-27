package handler

type LoggerRepository interface {
	SetMaxLength(int)
	SetMaxRotation(int)
	SetSystemOperation(int)
	GetSystemOperation() int
	Debug(string, ...interface{})
	Mutex(string, ...interface{})
	Trace(string, ...interface{})
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
	Fatal(string, ...interface{})
}
