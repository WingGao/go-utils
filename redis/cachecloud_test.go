package redis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCacheCloud(t *testing.T) {
	addrs, err := GetCacheCloud("http://fat-cache.ppdaicorp.com/cache/client/redis/cluster/10182.json?clientVersion=1.0")
	assert.Len(t, addrs, 6)
	assert.NoError(t, err)
}
