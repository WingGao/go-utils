module github.com/WingGao/go-utils

go 1.12

require (
	github.com/RichardKnop/machinery v1.6.2
	github.com/Shopify/sarama v1.23.0
	github.com/WingGao/errors v0.0.0-00010101000000-000000000000
	github.com/ajg/form v1.5.1 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v0.0.0-20190808125512-07798873deee
	github.com/chanxuehong/wechat v0.0.0-20190521093015-fafb751f9916
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/docker/docker v1.13.1
	github.com/elazarl/goproxy v0.0.0-20190421051319-9d40249d3c2f // indirect
	github.com/elazarl/goproxy/ext v0.0.0-20190421051319-9d40249d3c2f // indirect
	github.com/emirpasic/gods v1.12.0
	github.com/fasthttp-contrib/websocket v0.0.0-20160511215533-1f3b11f56072 // indirect
	github.com/fatih/structs v1.1.0
	github.com/fortytw2/leaktest v1.3.0 // indirect
	github.com/globalsign/mgo v0.0.0
	github.com/go-errors/errors v1.0.1
	github.com/go-playground/form v3.1.4+incompatible
	github.com/go-redis/redis/v7 v7.0.0-beta.4
	github.com/go-sql-driver/mysql v1.4.1
	github.com/iancoleman/strcase v0.0.0-20190422225806-e506e3ef7365
	github.com/imdario/mergo v0.3.8
	github.com/imkira/go-interpol v1.1.0 // indirect
	github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3
	github.com/jinzhu/gorm v1.9.11
	github.com/json-iterator/go v1.1.9
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/kataras/golog v0.0.10
	github.com/kataras/iris/v12 v12.1.4
	github.com/klauspost/cpuid v1.2.1 // indirect
	github.com/mailru/easyjson v0.0.0-20190403194419-1ea4449da983 // indirect
	github.com/micro/go-micro v1.7.0
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/olivere/elastic v6.2.17+incompatible
	github.com/parnurzeal/gorequest v0.2.15
	github.com/pierrec/lz4 v2.0.5+incompatible // indirect
	github.com/qiniu/api.v7 v7.2.5+incompatible
	github.com/qiniu/x v7.0.8+incompatible // indirect
	github.com/rs/xid v1.2.1
	github.com/sirupsen/logrus v1.4.2
	github.com/sony/sonyflake v0.0.0-20181109022403-6d5bd6181009
	github.com/stretchr/testify v1.4.0
	github.com/t-tiger/gorm-bulk-insert v0.0.0-00010101000000-000000000000
	github.com/teris-io/shortid v0.0.0-20171029131806-771a37caa5cf
	github.com/thoas/go-funk v0.4.0
	github.com/ungerik/go-dry v0.0.0-20180411133923-654ae31114c8
	github.com/valyala/fasthttp v1.4.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190809123943-df4f5c81cb3b // indirect
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0 // indirect
	github.com/yudai/gojsondiff v1.0.0 // indirect
	github.com/yudai/golcs v0.0.0-20170316035057-ecda9a501e82 // indirect
	github.com/yudai/pp v2.0.1+incompatible // indirect
	go.mongodb.org/mongo-driver v1.1.2
	go.uber.org/zap v1.10.0
	gopkg.in/jcmturner/goidentity.v3 v3.0.0 // indirect
	gopkg.in/yaml.v2 v2.2.4
	qiniupkg.com/x v7.0.8+incompatible // indirect
)

replace (
	github.com/RichardKnop/machinery => ../../RichardKnop/machinery
	github.com/WingGao/errors => ../../WingGao/errors
	github.com/chanxuehong/wechat => ../../chanxuehong/wechat
	github.com/fatih/structs => ../../fatih/structs
	github.com/globalsign/mgo => ../../globalsign/mgo
	github.com/go-errors/errors => ../../go-errors/errors
	github.com/jinzhu/copier => ../../jinzhu/copier
	github.com/jinzhu/gorm => ../../jinzhu/gorm
	github.com/kataras/iris/v12 => ../../kataras/iris
	github.com/micro/go-micro => ../../micro/go-micro
	github.com/t-tiger/gorm-bulk-insert => ../../t-tiger/gorm-bulk-insert
)
