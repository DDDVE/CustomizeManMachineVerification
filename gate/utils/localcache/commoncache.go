package localcache

import (
	"sync"
	"time"
)

var localCache *commomCache = &commomCache{}

type commomCache struct {
	m sync.Map
	// sync.Mutex
}

type Entry struct {
	Expire time.Time
	Value  interface{}
}

func (c *commomCache) Set(key string, value interface{}, expiration time.Duration) {
	entry := &Entry{
		Expire: time.Now().Add(expiration),
		Value:  value,
	}
	c.m.Store(key, entry)
}

func (c *commomCache) Get(key string) (value interface{}, expired bool) {
	v, ok := c.m.Load(key)
	e, ok := v.(*Entry)
	if !ok {
		return nil, false
	}
	return e.Value, e.Expire.Before(time.Now())
}

// func (c *commomCache) GetAndSet(key string, handler func(old *Entry) (new *Entry)) (new *Entry) {
// 	c.Lock()
// 	defer c.Unlock()
// 	v, ok := c.m.Load(key)
// 	val, ok := v.(*Entry)
// 	if !ok {
// 		return nil
// 	}
// 	newval := handler(val)
// 	c.m.Store(key, newval)
// 	return newval
// }
