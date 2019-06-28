package redis

import (
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
	return strings.Split(rep.ShardInfo, " "), err
}
