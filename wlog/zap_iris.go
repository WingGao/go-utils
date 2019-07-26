package wlog

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"go.uber.org/zap"
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

	fullEncoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
)

func NewZapToIris(app *iris.Application) *ZapToIris {
	alogger := &ZapToIris{app: app}
	alogger.encoder = getFullEncoder()
	return alogger
}

type ZapToIris struct {
	app     *iris.Application
	encoder zapcore.Encoder
}

func (c *ZapToIris) Enabled(zapcore.Level) bool {
	return true
}

func (c *ZapToIris) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return ce.AddCore(ent, c)
}

func (c *ZapToIris) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	//e := errors.Wrap(error(nil), 3)
	//e := errors.New(error(nil))
	// 0 /Users/ppd-03020144/Projs/go-web/src/github.com/WingGao/go-utils/wlog/zap_iris.go
	// 1 /Users/ppd-03020144/go/pkg/mod/go.uber.org/zap@v1.10.0/zapcore/entry.go
	// 2 /Users/ppd-03020144/go/pkg/mod/go.uber.org/zap@v1.10.0/sugar.go
	// 3 /Users/ppd-03020144/go/pkg/mod/go.uber.org/zap@v1.10.0/sugar.go
	//e.StackFrames()
	//var linef errors.StackFrame
	//for _, l := range e.Callers() {
	//	linef = errors.NewStackFrame(l)
	//	//fmt.Println(linef)
	//	if !strings.HasPrefix(linef.Package, "go.uber.org") && !strings.Contains(linef.File, "mod/go.uber.org") {
	//		break
	//	}
	//}
	bf, err := c.encoder.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}
	strMsg := bf.String()
	//linef := errors.NewStackFrame(e.Callers()[0])
	// 转换level
	il := zapToIrisLevelMap[ent.Level]
	c.app.Logger().Log(il, strMsg)
	return nil
}

func (c *ZapToIris) Sync() error {
	return nil
}

func (c *ZapToIris) With([]zapcore.Field) zapcore.Core {
	return c
}

func getFullEncoder() zapcore.Encoder {
	cnf := zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		//TimeKey:        "T",
		//LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	encoder := zapcore.NewConsoleEncoder(cnf)
	return encoder
}
