package cache

type CacheService struct {
	storage CacheStorage
}

func New() *CacheService {
	return &CacheService{
		storage: NewMemCache(),
	}
}

func (cs *CacheService) Start() error {
	if err := SetStorage(cs.storage); err != nil {
		return err
	}
	return nil
}

func (cs *CacheService) Stop() error {
	return cs.storage.Flush()
}
