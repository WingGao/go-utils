package redis

type RedisMutex struct {
	name    string
	timeout int
}

// 简陋redis锁
// name: 锁的key
// timeout: 过期秒数
func NewRedisMutex(name string, timeout int) *RedisMutex {
	m := &RedisMutex{name: name, timeout: timeout}
	return m
}

func (m *RedisMutex) Lock() bool {
	if r, _ := MainClient.Incr(m.name); r > 1 {
		return false
	}
	MainClient.Expire(m.name, m.timeout)
	return true
}

func (m *RedisMutex) Unlock() {
	MainClient.Del(m.name)
}
