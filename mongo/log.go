package mongo

import (
	"github.com/Sirupsen/logrus"
	"os"
	"fmt"
	"strings"
)

type MLogger struct {
	logrus.Logger
	OnlyDebugOp bool
}

func (l MLogger) Output(calldepth int, s string) error {
	if !l.OnlyDebugOp {
		l.Info(s)
	}
	return nil
}

func (l MLogger) Debug(calldepth int, s string) error {
	if !l.OnlyDebugOp {
		l.Logger.Debug(s)
	}
	return nil
}

func (l MLogger) DebugOp(format string, v ...interface{}) error {
	str := fmt.Sprintf(format, v...)
	if strings.Contains(str, "bson.DocElem{Name:\"ping\"") || strings.Contains(str, "bson.DocElem{Name:\"ismaster\"") {
		return nil
	}
	l.Logger.Debugf(format, v...)
	return nil
}

func GetLogger() *MLogger {
	var log = logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.DebugLevel
	return &MLogger{Logger: *log}
}
