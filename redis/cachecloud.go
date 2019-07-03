package redis

import (
	"github.com/WingGao/go-utils/wlog"
	"github.com/json-iterator/go"
	"github.com/ungerik/go-dry"
	"strings"
)

type cacheCouldRep struct {
	ShardInfo string `json:"shardInfo"`
}

func GetCacheCloud(url string) (addrs []string, err error) {
	j, err1 := dry.FileGetString(url)
	if err1 != nil {
		return addrs, err1
	}
	rep := &cacheCouldRep{}
	jsoniter.UnmarshalFromString(j, rep)
	wlog.S().Infof("Redis CacheCloud %s", j)
	return strings.Split(rep.ShardInfo, " "), err
}
