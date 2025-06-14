package cache

import (
	"sync"
	"time"
)

// ImageURLCache provides caching for presigned URLs
type ImageURLCache struct {
	cache    map[string]string
	expires  map[string]time.Time
	duration time.Duration
	mu       sync.RWMutex
}

// NewImageURLCache creates a new URL cache with the specified TTL
func NewImageURLCache(duration time.Duration) *ImageURLCache {
	return &ImageURLCache{
		cache:    make(map[string]string),
		expires:  make(map[string]time.Time),
		duration: duration,
	}
}

// Get retrieves a URL from the cache
func (c *ImageURLCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	url, exists := c.cache[key]
	if !exists {
		return "", false
	}

	// Check if the URL has expired
	if expiry, ok := c.expires[key]; ok {
		if time.Now().After(expiry) {
			// URL has expired, remove it from cache
			// We need to unlock and relock with a write lock
			c.mu.RUnlock()
			c.mu.Lock()
			delete(c.cache, key)
			delete(c.expires, key)
			c.mu.Unlock()
			c.mu.RLock()
			return "", false
		}
	}

	return url, true
}

// Set adds a URL to the cache
func (c *ImageURLCache) Set(key string, url string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = url
	c.expires[key] = time.Now().Add(c.duration)
}

// Delete removes a URL from the cache
func (c *ImageURLCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)
	delete(c.expires, key)
}

// Clear removes all URLs from the cache
func (c *ImageURLCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]string)
	c.expires = make(map[string]time.Time)
}
