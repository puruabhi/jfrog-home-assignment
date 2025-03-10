package types

import "fmt"

//go:generate mockgen -destination=./mocks/mock_logger.go -source=logger.go -package=mocks . Logger

type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Logf(format string, args ...any)
}

type fmtLogger struct{}

func NewFmtLogger() *fmtLogger {
	return &fmtLogger{}
}

func (l *fmtLogger) Debugf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func (l *fmtLogger) Infof(format string, args ...any) {
	fmt.Printf(format, args...)
}

func (l *fmtLogger) Warnf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func (l *fmtLogger) Errorf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func (l *fmtLogger) Logf(format string, args ...any) {
	fmt.Printf(format, args...)
}

type loggerStub struct{}

func NewLoggerStub() *loggerStub {
	return &loggerStub{}
}

func (l *loggerStub) Debugf(format string, args ...any) {}
func (l *loggerStub) Infof(format string, args ...any)  {}
func (l *loggerStub) Warnf(format string, args ...any)  {}
func (l *loggerStub) Errorf(format string, args ...any) {}
func (l *loggerStub) Logf(format string, args ...any)   {}
