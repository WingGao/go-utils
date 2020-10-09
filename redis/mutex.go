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
	if r, _ := MainClient.CtxIncr(m.name).Result(); r > 1 {
		return false
	}
	b, _ := MainClient.ExpireSecond(m.name, m.timeout)
	return b
}

func (m *RedisMutex) Unlock() {
	MainClient.CtxDel(m.name)
}
