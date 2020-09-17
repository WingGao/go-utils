package elasticutil

import (
	"github.com/elastic/go-elasticsearch/v7"
)

type EsModel struct {
	Id     string
	Client *elasticsearch.Client
	// 指向父的指针
	parent interface{} `bson:"-"`
}
