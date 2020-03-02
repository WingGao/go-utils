package rate

import (
	"github.com/WingGao/go-utils/redis"
	"time"
)

// 用redis来限制频率,true表示可以进行操作，并计数
// 在单位时间内`unit`,最多进行`maxOp`次操作
func RedisRate(key string, dur time.Duration, maxOp int) bool {
	if v, e := redis.MainClient.Incr(key).Result(); e != nil || v > int64(maxOp) {
		return false
	} else if v == 1 { //初次设置过期时间
		redis.MainClient.Expire(key, dur)
	}
	return true
}
