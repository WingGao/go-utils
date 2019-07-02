package redis

import (
	"github.com/WingGao/go-utils/redis"
	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/kataras/golog"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/sessions/sessiondb/redis/service"
	"runtime"
	"time"
)

// Database the redis back-end session database for the sessions.
type Database struct {
	redis  redis.RedisClient
	config service.Config
}

var _ sessions.Database = (*Database)(nil)

// New returns a new redis database.
func New(rd redis.RedisClient, cfg ...service.Config) *Database {
	db := &Database{redis: rd, config: cfg[0]}
	pong, err := db.redis.Ping().Result()
	if err != nil {
		golog.Debugf("error connecting to redis: %v", err)
		return nil
	}
	golog.Debugf("redis ping ==> %s", pong)
	runtime.SetFinalizer(db, closeDB)
	return db
}

// Config returns the configuration for the redis server bridge, you can change them.
func (db *Database) Config() *service.Config {
	return &db.config
}
func (db *Database) ttl(key string) (seconds int64, hasExpiration bool, found bool) {
	redisVal, err := db.redis.TTL(db.config.Prefix + key).Result()
	if err != nil {
		return -2, false, false
	}
	seconds = int64(redisVal.Seconds())
	// if -1 means the key has unlimited life time.
	hasExpiration = seconds > -1
	// if -2 means key does not exist.
	//found = !(c.Err() != nil || seconds == -2)
	found = !(seconds == -2)
	return
}

// Acquire receives a session's lifetime from the database,
// if the return value is LifeTime{} then the session manager sets the life time based on the expiration duration lives in configuration.
func (db *Database) Acquire(sid string, expires time.Duration) sessions.LifeTime {
	seconds, hasExpiration, found := db.ttl(sid)
	if !found {
		// not found, create an entry with ttl and return an empty lifetime, session manager will do its job.
		if err := db.redis.Set(sid, sid, time.Duration(expires.Seconds())*time.Second); err != nil {
			golog.Debug(err)
		}

		return sessions.LifeTime{} // session manager will handle the rest.
	}

	if !hasExpiration {
		return sessions.LifeTime{}

	}

	return sessions.LifeTime{Time: time.Now().Add(time.Duration(seconds) * time.Second)}
}

// https://redis.io/commands/expire#refreshing-expires
func (db *Database) getKeys(prefix string) (*sll.List, error) {
	iter := db.redis.Scan(0, db.config.Prefix+prefix+"*", 9999999999).Iterator()
	l := sll.New()
	for iter.Next() {
		l.Add(iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return l, nil
}
func (db *Database) UpdateTTLMany(prefix string, newSecondsLifeTime int64) error {
	keys, err := db.getKeys(prefix)
	if err != nil {
		return err
	}
	for i := 0; i < keys.Size(); i++ {
		k, _ := keys.Get(i)
		if _, err = db.redis.ExpireSecond(k.(string), int(newSecondsLifeTime)); err != nil {
			return err
		}
	}
	return nil
}

// OnUpdateExpiration will re-set the database's session's entry ttl.
// https://redis.io/commands/expire#refreshing-expires
func (db *Database) OnUpdateExpiration(sid string, newExpires time.Duration) error {
	return db.UpdateTTLMany(sid, int64(newExpires.Seconds()))
}

const delim = ":"

func makeKey(sid, key string) string {
	return sid + delim + key
}

// Set sets a key value of a specific session.
// Ignore the "immutable".
func (db *Database) Set(sid string, lifetime sessions.LifeTime, key string, value interface{}, immutable bool) {
	valueBytes, err := sessions.DefaultTranscoder.Marshal(value)
	if err != nil {
		golog.Error(err)
		return
	}

	if err = db.redis.Set(makeKey(sid, key), valueBytes, lifetime.DurationUntilExpiration()).Err(); err != nil {
		golog.Debug(err)
	}
}

// Get retrieves a session value based on the key.
func (db *Database) Get(sid string, key string) (value interface{}) {
	db.get(makeKey(sid, key), &value)
	return
}

func (db *Database) get(key string, outPtr interface{}) {
	data, err := db.redis.Get(key).Bytes()
	if err != nil {
		// not found.
		return
	}

	if err = sessions.DefaultTranscoder.Unmarshal(data, outPtr); err != nil {
		golog.Debugf("unable to unmarshal value of key: '%s': %v", key, err)
	}
}

func (db *Database) keys(sid string) *sll.List {
	keys, err := db.getKeys(sid + delim)
	if err != nil {
		golog.Debugf("unable to get all redis keys of session '%s': %v", sid, err)
		return nil
	}

	return keys
}

// Visit loops through all session keys and values.
func (db *Database) Visit(sid string, cb func(key string, value interface{})) {
	keys := db.keys(sid)
	keys.Each(func(index int, key interface{}) {
		var value interface{} // new value each time, we don't know what user will do in "cb".
		db.get(key.(string), &value)
		cb(key.(string), value)
	})
}

// Len returns the length of the session's entries (keys).
func (db *Database) Len(sid string) (n int) {
	return db.keys(sid).Size()
}

// Delete removes a session key value based on its key.
func (db *Database) Delete(sid string, key string) (deleted bool) {
	err := db.redis.Del(makeKey(sid, key)).Err()
	if err != nil {
		golog.Error(err)
	}
	return err == nil
}

// Clear removes all session key values but it keeps the session entry.
func (db *Database) Clear(sid string) {
	keys := db.keys(sid)
	keys.Each(func(index int, value interface{}) {
		key := value.(string)
		if err := db.redis.Del(key); err != nil {
			golog.Debugf("unable to delete session '%s' value of key: '%s': %v", sid, key, err)
		}
	})
}

// Release destroys the session, it clears and removes the session entry,
// session manager will create a new session ID on the next request after this call.
func (db *Database) Release(sid string) {
	// clear all $sid-$key.
	db.Clear(sid)
	// and remove the $sid.
	db.redis.Del(sid)
}

// Close terminates the redis connection.
func (db *Database) Close() error {
	return closeDB(db)
}

func closeDB(db *Database) error {
	//return db.redis.Close()
	return nil
}
