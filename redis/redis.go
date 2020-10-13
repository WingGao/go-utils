package redis

import (
	"context"
	"fmt"
	"github.com/emirpasic/gods/lists/arraylist"
	gredis "github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"time"
)

const (
	REDIS_UNIQUE_ID_KEY = "REDIS_UNIQUE_ID_KEY"
)

type RedisClient interface {
	gredis.UniversalClient
	ISingleRedis
	FullKey(key string) string
	GetConfig() RedisConf
	ExpireSecond(key string, second int) (bool, error)
	SetGlob(key string, ptr interface{}, opt *Option) error
	GetGlob(key string, out interface{}) error
	DelAll(keyPatter string) (count uint64, err error)
	Batch(keyPatter string, batchSize int, act func(keys []string) error) (count uint64, err error)
}

type ISingleRedis interface {
	Ctx() context.Context
	CtxIncr(key string) *gredis.IntCmd
	CtxGet(key string) *gredis.StringCmd
	CtxSet(key string, value interface{}, expiration time.Duration) *gredis.StatusCmd
	CtxDel(keys ...string) *gredis.IntCmd
	CtxExists(keys ...string) *gredis.IntCmd
	CtxExpire( key string, expiration time.Duration) *gredis.BoolCmd
	CtxSetNX(key string, value interface{}, expiration time.Duration) *gredis.BoolCmd
	CtxMGet(keys ...string) *gredis.SliceCmd


	CtxScan( cursor uint64, match string, count int64) *gredis.ScanCmd

	CtxSAdd(key string, members ...interface{}) *gredis.IntCmd
	CtxSMembers( key string) *gredis.StringSliceCmd
	CtxSPopN( key string, count int64) *gredis.StringSliceCmd
}

type Option struct {
	ExpireSecond int
}

func (m Option) ToInterface() []interface{} {
	arr := arraylist.New()
	if m.ExpireSecond > 0 {
		arr.Add("EX", m.ExpireSecond)
	}
	return arr.Values()
}

var MainClient RedisClient

//
//// 一般我们在一个系统里面使用redis
//// 所以该Client下的基本命令都会自动追加Prefix
//type RedisClientOne struct {
//	Config utils.RedisConf
//	pool   *redis.Pool
//	cc     *gredis.ClusterClient
//}
//
//
//
//
////获取附带Prefix的完整key
//func (c *RedisClientOne) FullKey(key string) string {
//	return c.Config.Prefix + key
//}
//
//func (c *RedisClientOne) newPool() *redis.Pool {
//	return &redis.Pool{
//		MaxIdle:     3,
//		IdleTimeout: 240 * time.Second,
//		Dial: func() (redis.Conn, error) {
//			var opts = make([]redis.DialOption, 0)
//			if c.Config.Password != "" {
//				opts = append(opts, redis.DialPassword(c.Config.Password))
//			}
//			conn, err := redis.Dial("tcp", c.Config.Addr, opts...)
//			if err != nil {
//				return nil, err
//			}
//
//			if c.Config.Database > 0 {
//				_, err = conn.Do("SELECT", c.Config.Database)
//			}
//
//			return conn, err
//		},
//	}
//}
//func (c *RedisClientOne) useCluster() {
//	var cc = gredis.NewClusterClient(&gredis.ClusterOptions{
//		Addrs: c.Config.Shards,
//	})
//	c.cc = cc
//}
//func (c *RedisClientOne) Conn() (redis.Conn, error) {
//	con := c.pool.Get()
//	return con, con.Err()
//}
//
//func (c *RedisClientOne) Do(commandName string, args ...interface{}) (interface{}, error) {
//	return c.pool.Get().Do(commandName, args...)
//}
//
//func (c *RedisClientOne) Del(key string) (interface{}, error) {
//	key = c.FullKey(key)
//	return c.Do("DEL", key)
//}
//
//// https://redis.io/commands/set
//// EX seconds -- Set the specified expire time, in seconds.
//func (c *RedisClientOne) Set(key string, value interface{}, opts ...interface{}) (interface{}, error) {
//	key = c.FullKey(key)
//	return c.pool.Get().Do("SET", append([]interface{}{key, value}, opts...)...)
//}
//
//func (c *RedisClientOne) Incr(key string) (int64, error) {
//	key = c.FullKey(key)
//	out, err := redis.Int64(c.pool.Get().Do("INCR", key))
//	return out, err
//}
//func (c *RedisClientOne) IncrBy(key string, increment int) (int64, error) {
//	key = c.FullKey(key)
//	return redis.Int64(c.pool.Get().Do("INCRBY", key, increment))
//}
//
//func (c *RedisClientOne) Get(key string) (interface{}, error) {
//	key = c.FullKey(key)
//	return c.pool.Get().Do("GET", key)
//}
//
//func (c *RedisClientOne) GetString(key string, def string) (string, error) {
//	out, err := redis.String(c.Get(key))
//	if err != nil {
//		return def, err
//	}
//	return out, err
//}
//
//func (c *RedisClientOne) GetInt(key string, def int) (int, error) {
//	//key = c.FullKey(key)
//	out, err := redis.Int(c.Get(key))
//	if err != nil {
//		return def, err
//	}
//	return out, err
//}
//
//func (c *RedisClientOne) GetInt64(key string, def int64) (int64, error) {
//	//key = c.FullKey(key)
//	out, err := redis.Int64(c.Get(key))
//	if err != nil {
//		return def, err
//	}
//	return out, err
//}
//
//func (c *RedisClientOne) GetUint32(key string, def uint32) (uint32, error) {
//	out, err := redis.Uint64(c.Get(key))
//	if err != nil {
//		return def, err
//	}
//	return uint32(out), err
//}
//
//func (c *RedisClientOne) GetUint64(key string, def uint64) (uint64, error) {
//	out, err := redis.Uint64(c.Get(key))
//	if err != nil {
//		return def, err
//	}
//	return out, err
//}
//
//func (c *RedisClientOne) GetBytes(key string, def []byte) ([]byte, error) {
//	//key = c.FullKey(key)
//	out, err := redis.Bytes(c.Get(key))
//	if err != nil {
//		return def, err
//	}
//	return out, err
//}
//
////gob.Register
//func (c *RedisClientOne) SetGlob(key string, ptr interface{}, opt *Option) (error) {
//	return SetGlob(c, key, ptr, opt)
//}
//
//func (c *RedisClientOne) GetGlob(key string, out interface{}) (error) {
//	return GetGlob(c, key, out)
//}
//
//func (c *RedisClientOne) SetJson(key string, ptr interface{}, opt *Option) (error) {
//	b, err := jsoniter.Marshal(ptr)
//	if err != nil {
//		return err
//	}
//	_, err = c.Set(key, b, opt.ToInterface()...)
//	return err
//}
//
//func (c *RedisClientOne) GetJson(key string, out interface{}) (error) {
//	bs, err := redis.Bytes(c.Get(key))
//	if err != nil {
//		return err
//	}
//	err = jsoniter.Unmarshal(bs, out)
//	return err
//}
//
//// SADD
//func (c *RedisClientOne) Csadd(key string, members ...interface{}) (added bool, err error) {
//	key = c.FullKey(key)
//	added, err = redis.Bool(c.Do("SADD", append([]interface{}{key}, members...)...))
//	return
//}
//
//// SMEMBERS
////
//func (c *RedisClientOne) SMembersMap(key string) (map[string]struct{}, error) {
//	panic("not implied")
//}
//
//// SISMEMBER
//func (c *RedisClientOne) Csismember(key string, item interface{}) (ok bool, err error) {
//	key = c.FullKey(key)
//	rep, err1 := redis.Bool(c.Do("SISMEMBER", key, item))
//	if err1 != nil {
//		return false, err1
//	}
//	return rep, nil
//}
//
//// EXPIRE
//func (c *RedisClientOne) Expire(key string, second int) (err error) {
//	key = c.FullKey(key)
//	_, err = c.Do("EXPIRE", key, second)
//	return
//}
//
//// EXISTS
//func (c *RedisClientOne) Exists(key string) bool {
//	key = c.FullKey(key)
//	res, err := redis.Int(c.Do("EXISTS", key))
//	return res == 1 && err == nil
//}
//
//// KEYS, 补全前缀
//func (c *RedisClientOne) Keys(pattern string) (keys []string, err error) {
//	pattern = c.FullKey(pattern)
//	keys, err = redis.Strings(c.Do("KEYS", pattern))
//	return
//}
//
//func (c *RedisClientOne) Info() (string, error) {
//	return redis.String(c.Do("INFO"))
//}
//
//func (c *RedisClientOne) Ping() error {
//	panic("implement me")
//}

func LoadClient(conf RedisConf) (err error) {
	if MainClient == nil {
		MainClient, err = NewRedisClient(conf)
	}
	return
}

func NewRedisClient(conf RedisConf) (c RedisClient, err error) {
	uoption := &gredis.UniversalOptions{}
	if conf.CacheCloudUrl != "" {
		fmt.Println("redis use cachecloud")
		uoption.Addrs, err = GetCacheCloud(conf.CacheCloudUrl)
		conf.Addr = strings.Join(uoption.Addrs, ",")
	} else if len(conf.Shards) > 0 {
		uoption.Addrs = conf.Shards
	} else {
		uoption.Addrs = []string{conf.Addr}
	}
	uc := gredis.NewUniversalClient(uoption)
	c = &RedisUniversalClient{UniversalClient: uc, Config: conf,ctx: context.Background()}
	uc.AddHook(&rhook{client: c.(*RedisUniversalClient)})

	//TODO ping
	//_, err = client.Conn()
	return
}

//gob.Register
//目前使用json代替
func SetGlob(c RedisClient, key string, ptr interface{}, opt *Option) error {
	//key = c.FullKey(key)
	//var buf bytes.Buffer
	buf, err := jsoniter.Marshal(ptr)
	//enc := gob.NewEncoder(&buf)
	//err := enc.Encode(ptr)
	if err != nil {
		return err
	}
	if opt != nil {
		_, err = c.CtxSet(key, buf, time.Duration(opt.ExpireSecond)*time.Second).Result()
	} else {

		_, err = c.CtxSet(key, buf, 0).Result()
	}
	return err
}

func GetGlob(c RedisClient, key string, out interface{}) error {
	//key = c.FullKey(key)
	bs, err := c.CtxGet(key).Bytes()
	if err != nil {
		return err
	}
	err = jsoniter.Unmarshal(bs, out)
	//buf := bytes.NewBuffer(bs)
	//dec := gob.NewDecoder(buf)
	//err = dec.Decode(out)
	return err
}
