package redis

import (
	stdContext "context"
	uredis "github.com/WingGao/go-utils/redis"
	sredis "github.com/kataras/iris/v12/sessions/sessiondb/redis"
	"time"
)


// GoRedisClient is the interface which both
// go-redis' Client and Cluster Client implements.
type GoRedisClient interface {
	uredis.RedisClient
}

// GoRedisDriver implements the Sessions Database Driver
// for the go-redis redis driver. See driver.go file.
type GoRedisDriver struct {
	// Both Client and ClusterClient implements this interface.
	Client uredis.RedisClient
}

var defaultContext = stdContext.Background()


// Connect initializes the redis Client.
func (r *GoRedisDriver) Connect(c sredis.Config) error {
	return nil
}

// PingPong sends a ping message and reports whether
// the PONG message received successfully.
func (r *GoRedisDriver) PingPong() (bool, error) {
	pong, err := r.Client.Ping(defaultContext).Result()
	return pong == "PONG", err
}

// CloseConnection terminates the underline redis connection.
func (r *GoRedisDriver) CloseConnection() error {
	return r.Client.Close()
}

// Set stores a "value" based on the session's "key".
// The value should be type of []byte, so unmarshal can happen.
func (r *GoRedisDriver) Set(sid, key string, value interface{}) error {
	return r.Client.HSet(defaultContext, sid, key, value).Err()
}

// Get returns the associated value of the session's given "key".
func (r *GoRedisDriver) Get(sid, key string) (interface{}, error) {
	return r.Client.HGet(defaultContext, sid, key).Bytes()
}

// Exists reports whether a session exists or not.
func (r *GoRedisDriver) Exists(sid string) bool {
	n, err := r.Client.Exists(defaultContext, sid).Result()
	if err != nil {
		return false
	}

	return n > 0
}

// TTL returns any TTL value of the session.
func (r *GoRedisDriver) TTL(sid string) time.Duration {
	dur, err := r.Client.TTL(defaultContext, sid).Result()
	if err != nil {
		return 0
	}

	return dur
}

// UpdateTTL sets expiration duration of the session.
func (r *GoRedisDriver) UpdateTTL(sid string, newLifetime time.Duration) error {
	_, err := r.Client.Expire(defaultContext, sid, newLifetime).Result()
	return err
}

// GetAll returns all the key values under the session.
func (r *GoRedisDriver) GetAll(sid string) (map[string]string, error) {
	return r.Client.HGetAll(defaultContext, sid).Result()
}

// GetKeys returns all keys under the session.
func (r *GoRedisDriver) GetKeys(sid string) ([]string, error) {
	return r.Client.HKeys(defaultContext, sid).Result()
}

// Len returns the total length of key-values of the session.
func (r *GoRedisDriver) Len(sid string) int {
	return int(r.Client.HLen(defaultContext, sid).Val())
}

// Delete removes a value from the redis store.
func (r *GoRedisDriver) Delete(sid, key string) error {
	if key == "" {
		return r.Client.Del(defaultContext, sid).Err()
	}
	return r.Client.HDel(defaultContext, sid, key).Err()
}
