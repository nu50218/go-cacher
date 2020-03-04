package cacher

import (
	"sync"
	"sync/atomic"
	"time"
)

type Cacher interface {
	Load() interface{}
}

func New(expires time.Duration, get func() interface{}) Cacher {
	return &basicCacher{
		val:     atomic.Value{},
		expires: expires,
		get:     get,
	}
}

type basicCacher struct {
	val     atomic.Value
	expires time.Duration
	update  time.Time
	mut     sync.Mutex
	get     func() interface{}
}

func (c *basicCacher) Load() interface{} {
	c.mut.Lock()
	t := time.Now()
	if t.Sub(c.update) > c.expires {
		c.update = t
		c.val.Store(c.get())
	}
	c.mut.Unlock()
	return c.val.Load()
}
