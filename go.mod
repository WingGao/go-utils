module github.com/WingGao/go-utils

go 1.12

require (
	cloud.google.com/go v0.50.0 // indirect
	github.com/RichardKnop/machinery v1.6.2
	github.com/Shopify/sarama v1.26.4
	github.com/WingGao/errors v0.0.0
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.566
	github.com/chanxuehong/wechat v0.0.0-20190521093015-fafb751f9916
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/docker/docker v1.13.1
	github.com/elastic/go-elasticsearch/v7 v7.8.0
	github.com/elazarl/goproxy v0.0.0-20190421051319-9d40249d3c2f // indirect
	github.com/elazarl/goproxy/ext v0.0.0-20190421051319-9d40249d3c2f // indirect
	github.com/emirpasic/gods v1.12.0
	github.com/fatih/structs v1.1.0
	github.com/go-errors/errors v1.0.1
	github.com/go-playground/form/v4 v4.1.1
	github.com/go-redis/redis/v8 v8.11.3
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/iancoleman/strcase v0.0.0-20190422225806-e506e3ef7365
	github.com/imdario/mergo v0.3.8
	github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3
	github.com/jinzhu/gorm v1.9.12
	github.com/json-iterator/go v1.1.12
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/kataras/golog v0.1.7
	github.com/kataras/iris/v12 v12.1.4
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/olivere/elastic/v7 v7.0.20
	github.com/parnurzeal/gorequest v0.2.16
	github.com/rs/xid v1.2.1
	github.com/shirou/gopsutil/v3 v3.22.1
	github.com/sirupsen/logrus v1.4.2
	github.com/sony/sonyflake v0.0.0-20181109022403-6d5bd6181009
	github.com/stretchr/testify v1.7.0
	github.com/t-tiger/gorm-bulk-insert/v2 v2.0.1
	github.com/teris-io/shortid v0.0.0-20171029131806-771a37caa5cf
	github.com/thoas/go-funk v0.4.0
	github.com/ungerik/go-dry v0.0.0-20180411133923-654ae31114c8
	github.com/xeipuuv/gojsonpointer v0.0.0-20190809123943-df4f5c81cb3b // indirect
	go.mongodb.org/mongo-driver v1.1.2
	go.uber.org/zap v1.13.0
	golang.org/x/net v0.0.0-20210913180222-943fd674d43e
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320
	google.golang.org/grpc v1.26.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
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
