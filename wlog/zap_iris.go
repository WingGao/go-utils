package wlog

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"go.uber.org/zap/zapcore"
	"strings"
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
	// 0 /Users/ppd-03020144/Projs/go-web/src/github.com/WingGao/go-utils/wlog/zap_iris.go
	// 1 /Users/ppd-03020144/go/pkg/mod/go.uber.org/zap@v1.10.0/zapcore/entry.go
	// 2 /Users/ppd-03020144/go/pkg/mod/go.uber.org/zap@v1.10.0/sugar.go
	// 3 /Users/ppd-03020144/go/pkg/mod/go.uber.org/zap@v1.10.0/sugar.go
	//e.StackFrames()
	var linef errors.StackFrame
	for _, l := range e.Callers() {
		linef = errors.NewStackFrame(l)
		if !strings.HasPrefix(linef.Package, "go.uber.org") {
			break
		}
	}
	//linef := errors.NewStackFrame(e.Callers()[0])
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
