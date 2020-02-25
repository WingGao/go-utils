package redis

import (
	"context"
	"fmt"
	gredis "github.com/go-redis/redis/v7"
	"strings"
	"time"
)

var (
	CMD_WITH_KEY = map[string]int{
		"decr":     1,
		"decrby":   1,
		"expire":   1,
		"expireat": 1,
		"incr":     1,
		"incrby":   1,
		"del":      1,
		"get":      1,
		"set":      1,
		"sadd":     1,
		"ttl":      1,
		"scan":     0,
		"ping":     -1,
	}
)

type RedisUniversalClient struct {
	gredis.UniversalClient
	Config RedisConf
}

func (c *RedisUniversalClient) ExpireSecond(key string, second int) (bool, error) {
	return c.Expire(key, time.Duration(second)*time.Second).Result()
}

func (c *RedisUniversalClient) GetConfig() RedisConf {
	return c.Config
}

// 约定，不能使用与c.Config.Prefix 一样的开头
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
		args := cmd.Args()
		if v > 0 {
			key := args[v].(string)
			if !strings.HasPrefix(key, hk.client.Config.Prefix) {
				args[v] = hk.client.FullKey(key)
			}
		} else {
			switch cmdName {
			case "scan":
				if len(args) > 2 && args[2] == "match" {
					args[3] = hk.client.FullKey(args[3].(string))
				}
			}
		}
		//wlog.S().Debugf("%#v", cmd)
	} else {
		panic(fmt.Sprintf("redis command [%s] not checked", cmdName))
	}
	return ctx, nil
}
func (hk rhook) AfterProcess(ctx context.Context, cmd gredis.Cmder) error {
	return nil
}

func (rhook) BeforeProcessPipeline(ctx context.Context, cmds []gredis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (rhook) AfterProcessPipeline(ctx context.Context, cmds []gredis.Cmder) (error) {
	return nil
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

func (c *RedisUniversalClient) DelAll(keyPatter string) (count uint64, err error) {
	return c.Batch(keyPatter, 300, func(keys []string) error {
		// keys已经
		// CROSSSLOT Keys in request don't hash to the same slot
		for _,k := range keys{
			err = c.Del(k).Err()
			if err != nil{
				return err
			}
		}
		return nil
	})
}

// 批量操作 act如果返回err则停止遍历
// 出来的keys已经FullKey了
func (c *RedisUniversalClient) Batch(keyPatter string, batchSize int, act func(keys []string) error) (count uint64, err error) {
	var cursor uint64 = 0
	var keys []string
	for {
		keys, cursor = c.Scan(cursor, keyPatter, int64(batchSize)).Val()
		lenKey := uint64(len(keys))
		if lenKey > 0 {
			if err = act(keys); err != nil {
				return
			}
		}
		count += lenKey
		if cursor <= 0 {
			break
		}
	}
	return
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
