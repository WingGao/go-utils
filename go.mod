module github.com/WingGao/go-utils

go 1.12

require (
	github.com/ReneKroon/ttlcache v1.6.0 // indirect
	github.com/RichardKnop/machinery v1.6.2
	github.com/Shopify/sarama v1.26.4
	github.com/WingGao/errors v0.0.0
	github.com/ajg/form v1.5.1 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.566
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/bwmarrin/discordgo v0.20.2 // indirect
	github.com/chanxuehong/wechat v0.0.0-20190521093015-fafb751f9916
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/docker/docker v1.13.1
	github.com/elastic/go-elasticsearch/v7 v7.8.0
	github.com/elazarl/goproxy v0.0.0-20190421051319-9d40249d3c2f // indirect
	github.com/elazarl/goproxy/ext v0.0.0-20190421051319-9d40249d3c2f // indirect
	github.com/emirpasic/gods v1.12.0
	github.com/fasthttp-contrib/websocket v0.0.0-20160511215533-1f3b11f56072 // indirect
	github.com/fatih/structs v1.1.0
	github.com/fortytw2/leaktest v1.3.0 // indirect
	github.com/go-errors/errors v1.0.1
	github.com/go-playground/form/v4 v4.1.1
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/iancoleman/strcase v0.0.0-20190422225806-e506e3ef7365
	github.com/imdario/mergo v0.3.8
	github.com/imkira/go-interpol v1.1.0 // indirect
	github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3
	github.com/jinzhu/gorm v1.9.12
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/json-iterator/go v1.1.9
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/kataras/golog v0.0.10
	github.com/kataras/iris/v12 v12.1.4
	github.com/klauspost/cpuid v1.2.1 // indirect
	github.com/lib/pq v1.3.0 // indirect
	github.com/lucas-clemente/quic-go v0.14.1 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/mholt/certmagic v0.9.3 // indirect
	github.com/miekg/dns v1.1.27 // indirect
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/nats-io/nats-server/v2 v2.1.4 // indirect
	github.com/nlopes/slack v0.6.1-0.20191106133607-d06c2a2b3249 // indirect
	github.com/olivere/elastic/v7 v7.0.20
	github.com/parnurzeal/gorequest v0.2.16
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rs/xid v1.2.1
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/sony/sonyflake v0.0.0-20181109022403-6d5bd6181009
	github.com/stretchr/testify v1.5.1
	github.com/t-tiger/gorm-bulk-insert/v2 v2.0.1
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/teris-io/shortid v0.0.0-20171029131806-771a37caa5cf
	github.com/thoas/go-funk v0.4.0
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200122045848-3419fae592fc // indirect
	github.com/ungerik/go-dry v0.0.0-20180411133923-654ae31114c8
	github.com/valyala/fasthttp v1.4.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190809123943-df4f5c81cb3b // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0 // indirect
	github.com/yudai/gojsondiff v1.0.0 // indirect
	github.com/yudai/golcs v0.0.0-20170316035057-ecda9a501e82 // indirect
	github.com/yudai/pp v2.0.1+incompatible // indirect
	go.mongodb.org/mongo-driver v1.1.2
	go.uber.org/zap v1.13.0
	golang.org/x/net v0.0.0-20200222125558-5a598a2470a0
	google.golang.org/grpc v1.26.0 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/jcmturner/goidentity.v3 v3.0.0 // indirect
	gopkg.in/yaml.v2 v2.2.8
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace (
	github.com/RichardKnop/machinery => ../../RichardKnop/machinery
	github.com/WingGao/errors => ../../WingGao/errors
	github.com/chanxuehong/wechat => ../../chanxuehong/wechat
	github.com/fatih/structs => ../../fatih/structs
	github.com/go-errors/errors => ../../go-errors/errors
	github.com/jinzhu/copier => ../../jinzhu/copier
	github.com/jinzhu/gorm => ../../jinzhu/gorm
	github.com/kataras/iris/v12 => ../../kataras/iris
	github.com/t-tiger/gorm-bulk-insert/v2 => ../../t-tiger/gorm-bulk-insert
)
