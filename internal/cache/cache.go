package cache

import (
	"fmt"
	"log"
	"time"
)

type CacheStorage interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}

type MemCache struct {
	storage map[string][]byte
}

var (
	logger  = log.New(log.Writer(), "cache: ", log.Flags())
	storage = make(map[string][]byte)
)

func NewCache() *MemCache {
	logger.Println("New cache created")

	return &MemCache{
		storage: storage,
	}
}

func (c *MemCache) Get(key string) ([]byte, error) {
	value, ok := c.storage[key]
	if !ok {
		return nil, fmt.Errorf("cache miss for key: %s", key)
	}

	return value, nil
}

type CacheOptions struct {
	ttl int
}

// interval in seconds, 0 means no ttl
func WithTTL(ttl int) *CacheOptions {
	return &CacheOptions{ttl: ttl}
}

func (c *MemCache) Set(key string, value []byte, opts *CacheOptions) error {
	if opts != nil && opts.ttl != 0 {
		go (func() {
			time.Sleep(time.Duration(opts.ttl) * time.Second)
			delete(c.storage, key)
			logger.Printf("key %s has been deleted from cache\n", key)
		})()
	}

	c.storage[key] = value
	return nil
}
