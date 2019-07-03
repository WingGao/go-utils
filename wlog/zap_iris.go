package wlog

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"go.uber.org/zap/zapcore"
)

var (
	zapToIrisLevelMap = map[zapcore.Level]golog.Level{
		zapcore.DebugLevel:  golog.DebugLevel,
		zapcore.InfoLevel:   golog.InfoLevel,
		zapcore.WarnLevel:   golog.WarnLevel,
		zapcore.ErrorLevel:  golog.ErrorLevel,
		zapcore.PanicLevel:  golog.FatalLevel,
		zapcore.DPanicLevel: golog.FatalLevel,
	}
)

func NewZapToIris(app *iris.Application) *ZapToIris {
	return &ZapToIris{app: app}
}

type ZapToIris struct {
	app *iris.Application
}

func (c *ZapToIris) Enabled(zapcore.Level) bool {
	return true
}

func (c *ZapToIris) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return ce.AddCore(ent, c)
}

func (c *ZapToIris) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	e := errors.Wrap(error(nil), 4)
	linef := errors.NewStackFrame(e.Callers()[0])
	// 转换level
	il := zapToIrisLevelMap[ent.Level]
	c.app.Logger().Log(il, fmt.Sprintf("%s:%d\n    ", linef.File, linef.LineNumber), ent.Message)
	return nil
}

func (c *ZapToIris) Sync() error {
	return nil
}

func (c *ZapToIris) With([]zapcore.Field) zapcore.Core {
	return c
}
