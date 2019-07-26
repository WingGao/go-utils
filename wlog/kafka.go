package wlog

import (
	"github.com/Shopify/sarama"
	ucore "github.com/WingGao/go-utils/ucore"
	"github.com/ungerik/go-dry"
	"go.uber.org/zap/zapcore"
	"time"
)

func NewZapToKafka(producer sarama.AsyncProducer, topic, appid, logName string) *ZapToKafka {
	alogger := &ZapToKafka{producer: producer, topic: topic, appid: appid, name: dry.RealNetIP(),
		logName: logName}

	alogger.encoder = getFullEncoder()
	return alogger
}

type ZapToKafka struct {
	appid    string
	topic    string
	name     string
	logName  string
	producer sarama.AsyncProducer
	encoder  zapcore.Encoder
}

func (c *ZapToKafka) Enabled(zapcore.Level) bool {
	return true
}

func (c *ZapToKafka) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return ce.AddCore(ent, c)
}

func (c *ZapToKafka) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	/*
		{/Users/ppd-03020144/Projs/go-web/src/github.com/WingGao/go-utils/wlog/kafka.go 36 (*ZapToKafka).Write github.com/WingGao/go-utils/wlog 79087703}
		{/Users/ppd-03020144/go/pkg/mod/go.uber.org/zap@v1.10.0/zapcore/entry.go 215 (*CheckedEntry).Write go.uber.org/zap/zapcore 77411577}
		{/Users/ppd-03020144/go/pkg/mod/go.uber.org/zap@v1.10.0/sugar.go 234 (*SugaredLogger).log go.uber.org/zap 77519217}
		{/Users/ppd-03020144/go/pkg/mod/go.uber.org/zap@v1.10.0/sugar.go 138 GetCacheCloud github.com/WingGao/go-utils/redis 79969412}
		{/Users/ppd-03020144/go/pkg/mod/go.uber.org/zap@v1.10.0/sugar.go 138 (*SugaredLogger).Infof go.uber.org/zap 79969327}
		{/Users/ppd-03020144/Projs/go-web/src/github.com/WingGao/go-utils/redis/redis.go 262 NewRedisClient github.com/WingGao/go-utils/redis 79970883}
		{/Users/ppd-03020144/Projs/go-web/src/github.com/WingGao/go-utils/redis/redis.go 253 LoadClient github.com/WingGao/go-utils/redis 87070099}
		{/Users/ppd-03020144/Projs/go-web/src/github.com/WingGao/go-utils/redis/redis.go 252 LoadClient github.com/WingGao/go-utils/redis 87061781}
		{/Users/ppd-03020144/Projs/go-web/wingao.net/webproj/mcmd/main.go 25 main main 87090345}
		{/usr/local/Cellar/go/1.12.5/libexec/src/runtime/proc.go 200 main runtime 67310428}
		{/usr/local/Cellar/go/1.12.5/libexec/src/runtime/asm_amd64.s 1337 goexit runtime 67498673}
	*/
	bf, err := c.encoder.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}
	strMsg := bf.String()
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
	//fmt.Printf("send %#v\n", msg)
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
