package cache

import (
	"log/slog"
	"time"
)

var logger = slog.With("service", "cache")

type CacheStorage interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, opts *CacheOptions) error
	Flush() error
}

type CacheOptions struct {
	ttl time.Duration
}

// interval in seconds, 0 means no ttl
func WithTTL(ttl int) *CacheOptions {
	return &CacheOptions{time.Duration(ttl) * time.Second}
}

var storageInstance CacheStorage

func SetStorage(s CacheStorage) error {
	if storageInstance != nil {
		return storageInstance.Flush()
	}

	storageInstance = s
	return nil
}

func Get(key string) ([]byte, error) {
	return storageInstance.Get(key)
}

func Set(key string, value []byte, opts *CacheOptions) error {
	return storageInstance.Set(key, value, opts)
}

func Flush() error {
	return storageInstance.Flush()
}
