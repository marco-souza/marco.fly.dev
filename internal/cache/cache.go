package cache

import "log"

var (
	logger = log.New(log.Writer(), "cache: ", log.Flags())
)

type CacheStorage interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, opts *CacheOptions) error
	Flush() error
}

type CacheOptions struct {
	ttl int
}

// interval in seconds, 0 means no ttl
func WithTTL(ttl int) *CacheOptions {
	return &CacheOptions{ttl: ttl}
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
	logger.Printf("getting key %s from cache\n", storageInstance)
	return storageInstance.Get(key)
}

func Set(key string, value []byte, opts *CacheOptions) error {
	return storageInstance.Set(key, value, opts)
}

func Flush() error {
	return storageInstance.Flush()
}
