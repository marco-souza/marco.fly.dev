package cache

import (
	"fmt"
	"time"
)

type MemCache struct {
	storage map[string][]byte
}

func NewMemCache() MemCache {
	return MemCache{
		storage: make(map[string][]byte),
	}
}

func (c MemCache) Get(key string) ([]byte, error) {
	value, ok := c.storage[key]
	if !ok {
		return nil, fmt.Errorf("cache miss for key: %s", key)
	}

	return value, nil
}

func (c MemCache) Set(key string, value []byte, opts *CacheOptions) error {
	if opts != nil && opts.ttl != 0 {
		go (func() {
			time.Sleep(opts.ttl)
			delete(c.storage, key)
			logger.Printf("key `%s` has been deleted\n", key)
		})()
	}

	c.storage[key] = value
	logger.Printf("cached key `%s`\n", key)
	return nil
}

func (c MemCache) Flush() error {
	for key := range c.storage {
		delete(c.storage, key)
	}

	return nil
}
