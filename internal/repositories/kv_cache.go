package repositories

import (
	"sync"
)

type SimpleKVCache struct {
	mu    sync.RWMutex
	cache map[string]any
}

func NewSimpleKVCache() *SimpleKVCache {
	return &SimpleKVCache{
		cache: map[string]any{},
	}
}

func (c *SimpleKVCache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	cachedData, exists := c.cache[key]
	return cachedData, exists
}

func (c *SimpleKVCache) Put(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = value
}

func (c *SimpleKVCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}

func (c *SimpleKVCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = map[string]any{}
}
