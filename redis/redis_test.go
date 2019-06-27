package redis

import (
	"encoding/gob"
	"encoding/json"
	"github.com/WingGao/go-utils"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	clientCl *RedisClientCl
)

func TestMain(m *testing.M) {
	utils.DefaultConfig.Redis = utils.RedisConf{
		Shards: []string{"10.114.31.202:6423", "10.114.31.210:6434", "10.114.31.210:6435", "10.114.31.211:6456", "10.114.31.211:6457", "10.114.31.202:6424"},
		Prefix: "test:",
	}
	clientCl = NewRedisClientCl(utils.DefaultConfig.Redis)
	os.Exit(m.Run())
}

func TestRedisClient_GetInt(t *testing.T) {
	key := "test"
	MainClient.Set(key, 123)
	out, err := MainClient.GetInt(key, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(out)
	out, err = MainClient.GetInt("test_no", 233)
	t.Log(out, err)
}

func TestRedisClient_Set(t *testing.T) {
	out, err := MainClient.Set("test_ttl", "haha", "EX", 30)
	t.Log(out, err)
	t.Log(tset("test_ttl", "haha", 30, 40))
}

func tset(key string, value interface{}, opts ...interface{}) []interface{} {
	a := append([]interface{}{key, value}, opts...)
	return a
}

type Sa struct {
	FieldA int32
	fieldB int32
	FieldC interface{}
}
type sb struct {
	FieldA int32
}

func TestRedisClientGlob(t *testing.T) {
	gob.Register(Sa{})
	gob.Register(sb{})
	a := Sa{FieldA: 1, fieldB: 2, FieldC: sb{FieldA: 3}}
	testKey := xid.New().String()
	err := MainClient.SetGlob(testKey, &a)
	assert.NoError(t, err)
	b := &Sa{}
	err2 := MainClient.GetGlob(testKey, b)
	assert.NoError(t, err2)
	assert.Equal(t, a.FieldA, b.FieldA)
	assert.Equal(t, int32(3), b.FieldC.(sb).FieldA)
}

func TestRedisClient_Incr(t *testing.T) {
	tkey := "12345"
	v, err := MainClient.Incr(tkey)
	assert.Equal(t, int64(1), v)
	assert.NoError(t, err)
}

func TestRedisClient_GetUint64(t *testing.T) {
	var val uint64 = 185135722552891230
	key := xid.New().String()
	MainClient.Set(key, val)
	getVal, _ := MainClient.GetUint64(key, 0)
	assert.Equal(t, val, getVal)
	MainClient.Del(key)
}

func TestZadd(t *testing.T) {
	key := xid.New().String()
	val := struct {
		Val uint64
	}{Val: 185135722552891243}
	msg, _ := json.Marshal(val)
	MainClient.Do("ZADD", key, 1, msg)
	items, err := redis.ByteSlices(MainClient.Do("ZRANGEBYSCORE", key, 0, 2))
	assert.NoError(t, err)
	t.Log(string(items[0]))
	MainClient.Del(key)
}
