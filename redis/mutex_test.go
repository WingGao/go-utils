package redis

import (
	"testing"
	"github.com/docker/docker/pkg/testutil/assert"
)

func TestNewRedisMutex(t *testing.T) {
	key := "test_redis_mutex"
	mut := NewRedisMutex(key, 1*60)
	mut.Lock()
	aval, _ := MainClient.Get(key).Result()
	assert.Equal(t, int(1), aval)
}
