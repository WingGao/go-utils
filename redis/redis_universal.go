package redis

import (
	"context"
	"fmt"
	"github.com/WingGao/go-utils"
	gredis "github.com/go-redis/redis"
	"time"
)

var (
	CMD_WITH_KEY = map[string]bool{
		"decr": true,
		"decrby": true,
		"expire": true,
		"expireat": true,
		"incr": true,
		"incrby": true,
		"get": true,
		"set": true,
	}
)

type RedisUniversalClient struct {
	gredis.UniversalClient
	Config utils.RedisConf
}

func (c *RedisUniversalClient) ExpireSecond(key string, second int) (bool, error) {
	return c.Expire(key, time.Duration(second)*time.Second).Result()
}

func (c *RedisUniversalClient) GetConfig() utils.RedisConf {
	return c.Config
}
func (c *RedisUniversalClient) FullKey(key string) string {
	return c.Config.Prefix + key
}

type rhook struct {
	client *RedisUniversalClient
}

func (hk rhook) BeforeProcess(ctx context.Context, cmd gredis.Cmder) (context.Context, error) {
	// 更改名称
	cmdName := cmd.Name()
	if v, ok := CMD_WITH_KEY[cmdName]; ok {
		if v {
			args := cmd.Args()
			key := args[1].(string)
			args[1] = hk.client.FullKey(key)
		} else {
			// pass
		}
	} else {
		panic(fmt.Sprintf("redis command [%s] not checked", cmdName))
	}
	return ctx, nil
}

func (rhook) AfterProcess(ctx context.Context, cmd gredis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (rhook) BeforeProcessPipeline(ctx context.Context, cmds []gredis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (rhook) AfterProcessPipeline(ctx context.Context, cmds []gredis.Cmder) (context.Context, error) {
	return ctx, nil
}

//func (c *RedisClientCl) do(commandName string, args ...interface{}) (*gredis.Cmd) {
//	return c.cc.Do(append([]interface{}{commandName}, args...)...)
//}
//func (c *RedisClientCl) Cmdable() (gredis.Cmdable) {
//	return c.cc
//}
//
//func (c *RedisClientCl) cmdArg1(cmd, key string) *gredis.Cmd {
//	key = c.FullKey(key)
//	return c.do(cmd, key)
//}
//func (c *RedisClientCl) cmdArg2(cmd, key string, arg interface{}) *gredis.Cmd {
//	key = c.FullKey(key)
//	return c.do(cmd, key, arg)
//}
//
//func (c *RedisClientCl) Set(key string, value interface{}, opts ...interface{}) (interface{}, error) {
//	key = c.FullKey(key)
//	return c.Do("SET", append([]interface{}{key, value}, opts...)...)
//}
//func (c *RedisClientCl) Get(key string) (interface{}, error) {
//	return c.cmdArg1("GET", key).Result()
//}
//func (c *RedisClientCl) Del(key string) (interface{}, error) {
//	return c.cmdArg1("DEL", key).Result()
//}
//func (c *RedisClientCl) Expire(key string, second int) (err error) {
//	return c.cmdArg2("EXPIRE", key, second).Err()
//}
//func (c *RedisClientCl) Incr(key string) (int64, error) {
//	return c.cmdArg1("INCR", key).Int64()
//}
//func (c *RedisClientCl) IncrBy(key string, increment int) (int64, error) {
//	return c.cmdArg2("INCRBY", key, increment).Int64()
//}
//
//func (c *RedisClientCl) Keys(pattern string) (keys []string, err error) {
//	return c.Cmdable().Keys(c.FullKey(pattern)).Result()
//}
//
//// sets
//func (c *RedisClientCl) SMembersMap(key string) (map[string]struct{}, error) {
//	return c.Cmdable().SMembersMap(c.FullKey(key)).Result()
//}
//
//func (c *RedisClientCl) Ping() (err error) {
//	_, err = c.Do("PING")
//	return
//}
func (c *RedisUniversalClient) SetGlob(key string, ptr interface{}, opt *Option) (error) {
	return SetGlob(c, key, ptr, opt)
}

func (c *RedisUniversalClient) GetGlob(key string, out interface{}) (error) {
	return GetGlob(c, key, out)
}

//
//func NewRedisClientCl(conf utils.RedisConf) *RedisClientCl {
//	var client = &RedisClientCl{
//		Config: conf,
//	}
//	client.cc = gredis.NewClusterClient(&gredis.ClusterOptions{
//		Addrs: conf.Shards,
//	})
//	return client
//}