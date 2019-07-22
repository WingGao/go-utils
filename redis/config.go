package redis

type RedisConf struct {
	Addr          string   //host:port 127.0.0.1:6379
	Shards        []string //host:port 127.0.0.1:6379
	Password      string
	Database      int
	UniqueIdKey   string
	Prefix        string
	CacheCloudUrl string
}
