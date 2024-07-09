package cache

import (
	"fmt"
	"time"
)

type MemCache struct {
	Storage map[string][]byte
}

func NewMemCache() MemCache {
	return MemCache{
		Storage: make(map[string][]byte),
	}
}

func (c MemCache) Get(key string) ([]byte, error) {
	value, ok := c.Storage[key]
	if !ok {
		return nil, fmt.Errorf("cache miss for key: %s", key)
	}

	return value, nil
}

func (c MemCache) Set(key string, value []byte, opts *CacheOptions) error {
	if opts != nil && opts.ttl != 0 {
		go (func() {
			time.Sleep(opts.ttl)
			delete(c.Storage, key)
			logger.Printf("key `%s` has been deleted\n", key)
		})()
	}

	c.Storage[key] = value
	logger.Printf("cacheed key `%s`\n", key)
	return nil
}

func (c MemCache) Flush() error {
	for key := range c.Storage {
		delete(c.Storage, key)
	}

	return nil
}
