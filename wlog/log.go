package wlog

import (
	"go.uber.org/zap"
)

var (
	_logger *zap.Logger
)

func init() {
	_logger, _ = zap.NewDevelopment()
}
func SetLogger(logger *zap.Logger) {
	_logger = logger
}

func L() *zap.Logger {
	return _logger
}
func S() *zap.SugaredLogger {
	return _logger.Sugar()
}
