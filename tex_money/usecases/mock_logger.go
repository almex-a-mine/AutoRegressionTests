// LoggerRepositoryがusecasesに定義されているため、一旦usecasesフォルダに配置
// TODO: フォルダ構成を見直す際にmocksフォルダを作成してこのファイルを移動する
package usecases

import (
	"fmt"
	"tex_money/domain/handler"
)

type MockLogger struct {
}

func NewMockLogger() handler.LoggerRepository {
	return &MockLogger{}
}

func (l *MockLogger) SetMaxLength(int)                    {}
func (l *MockLogger) SetMaxRotation(int)                  {}
func (l *MockLogger) SetSystemOperation(int)              {}
func (l *MockLogger) GetSystemOperation() int             { return 1 }
func (l *MockLogger) Debug(s string, v ...interface{})    { fmt.Printf(s+"\n", v...) }
func (l *MockLogger) Mutex(s string, v ...interface{})    { fmt.Printf(s+"\n", v...) }
func (l *MockLogger) Trace(s string, v ...interface{})    { fmt.Printf(s+"\n", v...) }
func (l *MockLogger) Info(s string, v ...interface{})     { fmt.Printf(s+"\n", v...) }
func (l *MockLogger) Warn(s string, v ...interface{})     { fmt.Printf(s+"\n", v...) }
func (l *MockLogger) Error(s string, v ...interface{})    { fmt.Printf(s+"\n", v...) }
func (l *MockLogger) Fatal(s string, v ...interface{})    { fmt.Printf(s+"\n", v...) }
func (l *MockLogger) Sequence(s string, v ...interface{}) { fmt.Printf(s+"\n", v...) }
