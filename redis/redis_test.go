package redis

import (
	"testing"
	"os"
	"github.com/WingGao/go-utils"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"encoding/gob"
)

func TestMain(m *testing.M) {
	testConf, _ := utils.LoadConfig(os.Getenv("NXPT_GO_CONF"))
	LoadClient(testConf.Redis)
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

