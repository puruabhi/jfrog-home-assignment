package logger

import "fmt"

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
