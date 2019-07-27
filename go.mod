module github.com/WingGao/go-utils

go 1.12

require (
	github.com/RichardKnop/machinery v1.6.2
	github.com/Shopify/sarama v1.23.0
	github.com/ajg/form v1.5.1 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v0.0.0-20190430075129-62f3cb8727f4
	github.com/chanxuehong/wechat v0.0.0-20190521093015-fafb751f9916
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/docker/docker v1.13.1
	github.com/elazarl/goproxy v0.0.0-20190421051319-9d40249d3c2f // indirect
	github.com/elazarl/goproxy/ext v0.0.0-20190421051319-9d40249d3c2f // indirect
	github.com/emirpasic/gods v1.12.0
	github.com/fatih/structs v1.1.0
	github.com/fortytw2/leaktest v1.3.0 // indirect
	github.com/gavv/monotime v0.0.0-20190418164738-30dba4353424 // indirect
	github.com/globalsign/mgo v0.0.0
	github.com/go-check/check v1.0.0-20180628173108-788fd7840127 // indirect
	github.com/go-errors/errors v1.0.1
	github.com/go-playground/form v3.1.4+incompatible
	github.com/go-redis/redis v0.0.0-20190626123411-a28bb0bd25c8
	github.com/go-sql-driver/mysql v1.4.1
	github.com/imdario/mergo v0.3.7
	github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3
	github.com/jinzhu/gorm v1.9.5
	github.com/json-iterator/go v1.1.6
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/kataras/golog v0.0.0-20190624001437-99c81de45f40
	github.com/kataras/iris v11.1.1+incompatible
	github.com/klauspost/compress v1.5.0 // indirect
	github.com/klauspost/cpuid v1.2.1 // indirect
	github.com/mailru/easyjson v0.0.0-20190403194419-1ea4449da983 // indirect
	github.com/micro/go-micro v1.7.0
	github.com/olivere/elastic v6.2.17+incompatible
	github.com/parnurzeal/gorequest v0.2.15
	github.com/qiniu/api.v7 v7.2.5+incompatible
	github.com/qiniu/x v7.0.8+incompatible // indirect
	github.com/rs/xid v1.2.1
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/goconvey v0.0.0-20190330032615-68dc04aab96a // indirect
	github.com/sony/sonyflake v0.0.0-20181109022403-6d5bd6181009
	github.com/stretchr/testify v1.3.0
	github.com/teris-io/shortid v0.0.0-20171029131806-771a37caa5cf
	github.com/thoas/go-funk v0.4.0
	github.com/ungerik/go-dry v0.0.0-20180411133923-654ae31114c8
	github.com/xeipuuv/gojsonschema v1.1.0 // indirect
	github.com/yudai/pp v2.0.1+incompatible // indirect
	go.mongodb.org/mongo-driver v1.0.0
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20190621222207-cc06ce4a13d4 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/ini.v1 v1.42.0 // indirect
	gopkg.in/jcmturner/goidentity.v3 v3.0.0 // indirect
	gopkg.in/yaml.v2 v2.2.2
	qiniupkg.com/x v7.0.8+incompatible // indirect
)

replace (
	github.com/globalsign/mgo v0.0.0 => github.com/WingGao/mgo v0.0.0-20190502114913-db5d70d36ad5
	github.com/go-errors/errors => ../../go-errors/errors
	github.com/kataras/iris => ../../kataras/iris
)
