package wlog

import (
	"go.uber.org/zap"
)

var (
	_logger  *zap.Logger
	_loggerS *SugaredLogger
)

func init() {
	g, _ := zap.NewDevelopment()
	SetLogger(g)
}
func SetLogger(logger *zap.Logger) {
	_logger = logger
	_loggerS = &SugaredLogger{SugaredLogger: *logger.Sugar()}
}

func L() *zap.Logger {
	return _logger
}
func S() *SugaredLogger {
	return _loggerS
}

type SugaredLogger struct {
	zap.SugaredLogger
}

// 打印错误，true=有错误
func (s *SugaredLogger) IfError(e interface{}) bool {
	if e != nil {
		s.Error(e)
		return true
	}
	return false
}
