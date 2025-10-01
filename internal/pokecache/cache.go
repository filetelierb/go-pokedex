package pokecache

import (
	"sync"
	"time"
)


type cacheEntry struct {
	createdAt time.Time
	val []byte
}
type Cache struct {
	mu sync.Mutex
	cacheEntries map[string]cacheEntry
}

func (cache *Cache) Add(key string, value  []byte){
	cache.mu.Lock()
	defer cache.mu.Unlock()
	newCacheEntry := cacheEntry {
		createdAt: time.Now(),
		val: value,
	}
	cache.cacheEntries[key] = newCacheEntry
}

func (cache *Cache) Get(key string) ([]byte, bool){
	cache.mu.Lock()
	defer cache.mu.Unlock()
	value, ok := cache.cacheEntries[key]
	if !ok {
		return nil, ok
	}
	return value.val, ok
}

func (cache *Cache) reapLoop(t time.Time) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	for key, value := range cache.cacheEntries {
		if value.createdAt.Before(t){
			delete(cache.cacheEntries,key)
		}
	}
}


func NewCache(interval time.Duration) *Cache{
	cache := Cache{
		mu: sync.Mutex{},
		cacheEntries: make(map[string]cacheEntry),
	}
	ticker := time.NewTicker(interval)
	go func(){
        for t := range ticker.C {
            cache.reapLoop(t)
        }
		/*for {
			select {
			case t:= <-ticker.C:
				cache.reapLoop(t)
			}
		}*/
	}()

	return &cache
	
}