package cache

import (
	"log"
	"time"

	"github.com/jonfriesen/redis-proxy/cache/lrucache"
	"github.com/jonfriesen/redis-proxy/storage"
)

type Cache struct {
	dataSource *storage.Storage
	lruCache   *lrucache.LRUCache
}

func New(mkeys int32, mage time.Duration, ds *storage.Storage) *Cache {
	lc := lrucache.New(mkeys, mage)

	return &Cache{
		dataSource: ds,
		lruCache:   lc,
	}
}

func (c *Cache) Get(key string) (string, error) {

	c.lruCache.Lock()
	defer c.lruCache.Unlock()
	v, err := c.lruCache.Get(key)
	if err == lrucache.ErrNotFound {
		log.Printf("Lookup not in cache %v", key)

		v, err = (*c.dataSource).Get(key)
		if err == storage.ErrNotFound {
			log.Printf("Lookup not in storage %v", key)
			return "", storage.ErrNotFound
		}

		if v != "" {
			log.Println("Pushing key-value pair into cache")
			err = c.lruCache.Push(key, v)
		}
	}
	if err == lrucache.ErrUnlocked {
		log.Fatalf("Critical Error: %v", err)
	}

	return v, nil
}
