package redis

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"github.com/WingGao/go-utils"
	"encoding/gob"
	"bytes"
)

const (
	REDIS_UNIQUE_ID_KEY = "REDIS_UNIQUE_ID_KEY"
)

// 一般我们在一个系统里面使用redis
// 所以该Client下的基本命令都会自动追加Prefix
type RedisClient struct {
	Config utils.RedisConf
	pool   *redis.Pool
}

var MainClient *RedisClient

//获取附带Prefix的完整key
func (c *RedisClient) FullKey(key string) string {
	return c.Config.Prefix + key
}

func (c *RedisClient) newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			var opts = make([]redis.DialOption, 0)
			if c.Config.Password != "" {
				opts = append(opts, redis.DialPassword(c.Config.Password))
			}
			conn, err := redis.Dial("tcp", c.Config.Addr, opts...)
			if err != nil {
				return nil, err
			}

			if c.Config.Database != "" {
				_, err = conn.Do("SELECT", c.Config.Database)
			}

			return conn, err
		},
	}
}
func (c *RedisClient) Conn() (redis.Conn, error) {
	con := c.pool.Get()
	return con, con.Err()
}

func (c *RedisClient) Do(commandName string, args ...interface{}) (interface{}, error) {
	return c.pool.Get().Do(commandName, args...)
}

func (c *RedisClient) Del(key string) (interface{}, error) {
	key = c.FullKey(key)
	return c.Do("DEL", key)
}

// https://redis.io/commands/set
func (c *RedisClient) Set(key string, value interface{}, opts ...interface{}) (interface{}, error) {
	key = c.FullKey(key)
	return c.pool.Get().Do("SET", append([]interface{}{key, value}, opts...)...)
}

func (c *RedisClient) Incr(key string) (int64, error) {
	key = c.FullKey(key)
	out, err := redis.Int64(c.pool.Get().Do("INCR", key))
	return out, err
}
func (c *RedisClient) IncrBy(key string, increment int) (int64, error) {
	key = c.FullKey(key)
	return redis.Int64(c.pool.Get().Do("INCRBY", key, increment))
}

func (c *RedisClient) Get(key string) (interface{}, error) {
	key = c.FullKey(key)
	return c.pool.Get().Do("GET", key)
}

func (c *RedisClient) GetInt(key string, def int) (int, error) {
	key = c.FullKey(key)
	out, err := redis.Int(c.Get(key))
	if err != nil {
		return def, err
	}
	return out, err
}

func (c *RedisClient) GetInt64(key string, def int64) (int64, error) {
	key = c.FullKey(key)
	out, err := redis.Int64(c.Get(key))
	if err != nil {
		return def, err
	}
	return out, err
}

func (c *RedisClient) GetBytes(key string, def []byte) ([]byte, error) {
	key = c.FullKey(key)
	out, err := redis.Bytes(c.Get(key))
	if err != nil {
		return def, err
	}
	return out, err
}

//gob.Register
func (c *RedisClient) SetGlob(key string, ptr interface{}) (error) {
	key = c.FullKey(key)
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(ptr)
	if err != nil {
		return err
	}
	_, err = c.Set(key, buf.Bytes())
	return err
}
func (c *RedisClient) GetGlob(key string, out interface{}) (error) {
	key = c.FullKey(key)
	bs, err := redis.Bytes(c.Get(key))
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(bs)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(out)
	return err
}

// SADD
func (c *RedisClient) Csadd(key string, members ...interface{}) (added bool, err error) {
	key = c.FullKey(key)
	added, err = redis.Bool(c.Do("SADD", append([]interface{}{key}, members...)...))
	return
}

// SMEMBERS
//
func (c *RedisClient) Csmembers(key string, out interface{}) (err error) {
	key = c.FullKey(key)
	rep, err1 := redis.Values(c.Do("SMEMBERS", key))
	if err1 != nil {
		return err1
	}
	err = redis.ScanSlice(rep, out)
	return
}

// SISMEMBER
func (c *RedisClient) Csismember(key string, item interface{}) (ok bool, err error) {
	key = c.FullKey(key)
	rep, err1 := redis.Bool(c.Do("SISMEMBER", key, item))
	if err1 != nil {
		return false, err1
	}
	return rep, nil
}

// EXPIRE
func (c *RedisClient) Expire(key string, second int) (err error) {
	key = c.FullKey(key)
	_, err = c.Do("EXPIRE", key, second)
	return
}

func LoadClient(conf utils.RedisConf) error {
	if MainClient == nil {
		MainClient = &RedisClient{Config: conf}
		MainClient.pool = MainClient.newPool()
	}
	_, err := MainClient.Conn()
	return err
}
