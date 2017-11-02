package mongo

import (
	"github.com/Sirupsen/logrus"
	"os"
)

type MLogger struct {
	logrus.Logger
}

func (l MLogger) Output(calldepth int, s string) error {
	l.Debugf("%d => %s", calldepth, s)
	return nil
}

func GetLogger() *MLogger {
	var log = logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.DebugLevel
	return &MLogger{Logger: *log}
}
