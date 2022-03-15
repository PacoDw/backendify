package cache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

// Cache represents the main struct for the cache data, it is a wrapper of the go-cache library.
type Cache struct {
	*gocache.Cache
}

// LoadOrStore returns the existing value for the key if present. Otherwise, it stores
// and returns the given value. The loaded result is true if the value was loaded, false if stored.
func (c *Cache) StoreOrLoad(key string, value interface{}) (actual interface{}, loaded bool) {
	// Get the string associated with the key "foo" from the cache
	v, expiration, found := c.GetWithExpiration(key)
	if found {
		// verify if this key was expired if so then remove every key is currently expired.
		if expiration.After(time.Now()) {
			c.DeleteExpired()
		}

		return v, true
	}

	// if the key was not found then add it and add the same default time added (24hrs)
	c.Set(key, value, gocache.DefaultExpiration)

	return value, false
}

// ChainStoreOrLoad makes the same that StoreOrLoad but the difference is this one return
// the Cache instance.
// NOTE: Consider that it doesn't return any parameter to tell you that the data has been saved or loaded.
func (c *Cache) ChainStoreOrLoad(key string, value interface{}) *Cache {
	c.StoreOrLoad(key, value)

	return c
}

// New creates a new cache.
func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	return &Cache{
		gocache.New(defaultExpiration, cleanupInterval),
	}
}
