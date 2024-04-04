package cache

import (
	"sync"
)

type Cache struct {
	mux   sync.RWMutex
	store map[string]interface{}
}

func NewCache() *Cache {
	return &Cache{mux: sync.RWMutex{},
		store: make(map[string]interface{}),
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.mux.Lock()
	c.store[key] = value
	c.mux.Unlock()
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mux.RLock()
	result, ok := c.store[key]
	c.mux.RUnlock()
	return result, ok
}

func (c *Cache) Len() int {
	return len(c.store)
}
