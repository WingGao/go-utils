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

type RedisClient struct {
	Config utils.RedisConf
	pool   *redis.Pool
}

var MainClient *RedisClient

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

func (c *RedisClient) Do(commandName string, args ...interface{}) (interface{}, error) {
	return c.pool.Get().Do(commandName, args...)
}

func (c *RedisClient) Del(key string) (interface{}, error) {
	return c.Do("DEL", key)
}

// https://redis.io/commands/set
func (c *RedisClient) Set(key string, value interface{}, opts ...interface{}) (interface{}, error) {
	return c.pool.Get().Do("SET", append([]interface{}{key, value}, opts...)...)
}

func (c *RedisClient) Incr(key string) (int64, error) {
	out, err := redis.Int64(c.pool.Get().Do("INCR", key))
	return out, err
}
func (c *RedisClient) IncrBy(key string, increment int) (int64, error) {
	return redis.Int64(c.pool.Get().Do("INCRBY", key, increment))
}

func (c *RedisClient) Get(key string) (interface{}, error) {
	return c.pool.Get().Do("GET", key)
}

func (c *RedisClient) GetInt(key string, def int) (int, error) {
	out, err := redis.Int(c.Get(key))
	if err != nil {
		return def, err
	}
	return out, err
}

func (c *RedisClient) GetInt64(key string, def int64) (int64, error) {
	out, err := redis.Int64(c.Get(key))
	if err != nil {
		return def, err
	}
	return out, err
}

func (c *RedisClient) GetBytes(key string, def []byte) ([]byte, error) {
	out, err := redis.Bytes(c.Get(key))
	if err != nil {
		return def, err
	}
	return out, err
}
//gob.Register
func (c *RedisClient) SetGlob(key string, ptr interface{}) (error) {
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
	bs, err := redis.Bytes(c.Get(key))
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(bs)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(out)
	return err
}

func LoadClient(conf utils.RedisConf) {
	if MainClient == nil {
		MainClient = &RedisClient{Config: conf}
		MainClient.pool = MainClient.newPool()
	}
}
