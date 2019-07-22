package wlog

import (
	"fmt"
	"github.com/Shopify/sarama"
	ucore "github.com/WingGao/go-utils/core"
	"github.com/go-errors/errors"
	"github.com/ungerik/go-dry"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

func NewZapToKafka(producer sarama.AsyncProducer, topic, appid, logName string) *ZapToKafka {
	return &ZapToKafka{producer: producer, topic: topic, appid: appid, name: dry.RealNetIP(),
		logName: logName}
}

type ZapToKafka struct {
	appid    string
	topic    string
	name     string
	logName  string
	producer sarama.AsyncProducer
}

func (c *ZapToKafka) Enabled(zapcore.Level) bool {
	return true
}

func (c *ZapToKafka) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return ce.AddCore(ent, c)
}

func (c *ZapToKafka) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	e := errors.Wrap(error(nil), 3)
	var linef errors.StackFrame
	for _, l := range e.Callers() {
		linef = errors.NewStackFrame(l)
		//fmt.Println(linef)
		if !strings.HasPrefix(linef.Package, "go.uber.org") && !strings.Contains(linef.File, "mod/go.uber.org") {
			break
		}
	}
	//linef := errors.NewStackFrame(e.Callers()[0])
	strMsg := fmt.Sprintf("%s:%d\n    %s", linef.File, linef.LineNumber, ent.Message)
	now := time.Now()
	lmsg := &logMessage{
		AppId:         c.appid,
		ThreadName:    c.name,
		LogName:       c.logName,
		Level:         ent.Level.String(),
		LayoutMessage: strMsg,
		TimeStamp:     now.Unix() * 1000,
	}
	js := ucore.JsonMarshalToString(lmsg)
	msg := &sarama.ProducerMessage{
		Topic:     c.topic,
		Value:     sarama.StringEncoder(js),
		Timestamp: now,
	}
	c.producer.Input() <- msg
	return nil
}

func (c *ZapToKafka) Sync() error {
	return nil
}

func (c *ZapToKafka) With([]zapcore.Field) zapcore.Core {
	return c
}

type logMessage struct {
	AppId         string            `json:"appId"`
	ThreadName    string            `json:"threadName"`
	LogName       string            `json:"logName"`
	Level         string            `json:"level"`
	LayoutMessage string            `json:"layoutMessage"`
	TimeStamp     int64             `json:"timeStamp"`
	Tags          map[string]string `json:"tags"`
}
