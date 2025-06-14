package cache

import (
	"sync"
	"time"
)

// ImageURLCache defines the interface for a cache of image URLs
type ImageURLCache interface {
	Get(key string) (string, bool)
	Set(key, value string)
	Delete(key string)
}

// MemoryImageURLCache implements ImageURLCache using in-memory storage
type MemoryImageURLCache struct {
	cache      map[string]*cacheEntry
	expiration time.Duration
	mutex      sync.RWMutex
}

type cacheEntry struct {
	value      string
	expiration time.Time
}

// NewImageURLCache creates a new image URL cache with the given expiration time
func NewImageURLCache(expiration time.Duration) *MemoryImageURLCache {
	return &MemoryImageURLCache{
		cache:      make(map[string]*cacheEntry),
		expiration: expiration,
	}
}

// Get retrieves a value from the cache
func (c *MemoryImageURLCache) Get(key string) (string, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, found := c.cache[key]
	if !found {
		return "", false
	}

	// Check if entry has expired
	if time.Now().After(entry.expiration) {
		delete(c.cache, key)
		return "", false
	}

	return entry.value, true
}

// Set adds a value to the cache
func (c *MemoryImageURLCache) Set(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache[key] = &cacheEntry{
		value:      value,
		expiration: time.Now().Add(c.expiration),
	}
}

// Delete removes a value from the cache
func (c *MemoryImageURLCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.cache, key)
}
