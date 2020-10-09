package rapollo

import (
	"context"
	"github.com/WingGao/go-utils/wlog"
	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	gredis "github.com/go-redis/redis/v8"
)

/**
一个类似于apollo的配置同步工具，基于redis
*/

type Rapollo struct {
	KeyPrefix   string
	RedisClient gredis.UniversalClient
	listenerMap map[string]*sll.List
	ctx context.Context
}
type UpdateHandler func(payload string)

func (c *Rapollo) FullChannelKey(key string) string {
	return c.KeyPrefix + ":" + key
}
func (c *Rapollo) FullValueKey(key string) string {
	return c.KeyPrefix + ":_value:" + key
}

func (c *Rapollo) RegisterKey(key string, onUpdate UpdateHandler) {
	if key == "_value" { //不允许
		return
	}
	fk := c.FullChannelKey(key)
	hs := c.listenerMap[fk]
	if hs == nil {
		hs = sll.New()
		c.listenerMap[fk] = hs
	}
	hs.Append(onUpdate)
	// 首次调用
	res := c.RedisClient.Get(c.ctx,c.FullValueKey(key))
	if res.Err() == nil {
		onUpdate(res.Val())
	}
}

func (c *Rapollo) UnregisterKey(key string, onUpdate UpdateHandler) {
	fk := c.FullChannelKey(key)
	hs := c.listenerMap[fk]
	if hs != nil {
		if i := hs.IndexOf(onUpdate); i >= 0 {
			hs.Remove(i)
		}
	}
}

func (c *Rapollo) Pub(key string, value string) {
	fk := c.FullChannelKey(key)
	vk := c.FullValueKey(key)
	c.RedisClient.Publish(c.ctx,fk, value)
	c.RedisClient.Set(c.ctx,vk, value, 0)
}

// 完整的key
func New(keyPrefix string, redisc gredis.UniversalClient) *Rapollo {
	p := &Rapollo{KeyPrefix: keyPrefix, listenerMap: make(map[string]*sll.List),
		RedisClient: redisc,ctx: context.Background()}
	// 开始监听
	ps := p.RedisClient.PSubscribe(p.ctx,keyPrefix + ":*")
	go func() {
		for {
			if msgi, err := ps.ReceiveMessage(p.ctx); err == nil {
				wlog.S().Debugf("receive %#v", msgi)
				hs := p.listenerMap[msgi.Channel]
				if hs != nil {
					hs.Each(func(index int, value interface{}) {
						value.(UpdateHandler)(msgi.Payload)
					})
				}
			} else {
				wlog.S().IfError(err)
			}
		}
	}()
	return p
}
