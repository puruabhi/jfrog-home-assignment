package types

//go:generate mockgen -destination=./mocks/mock_logger.go -source=logger.go -package=mocks . Logger

type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

type loggerStub struct{}

func NewLoggerStub() *loggerStub {
	return &loggerStub{}
}

func (l *loggerStub) Debugf(format string, args ...any) {}
func (l *loggerStub) Infof(format string, args ...any)  {}
func (l *loggerStub) Warnf(format string, args ...any)  {}
func (l *loggerStub) Errorf(format string, args ...any) {}
