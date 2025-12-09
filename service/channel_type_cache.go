package service

import (
	"sync"
	"time"

	"github.com/QuantumNous/new-api/model"
)

// ChannelTypeCache provides a thread-safe cache for channel type counts
type ChannelTypeCache struct {
	counts     map[int64]int64
	mutex      sync.RWMutex
	lastUpdate time.Time
	ttl        time.Duration
}

var channelTypeCache *ChannelTypeCache
var channelTypeCacheOnce sync.Once

// GetChannelTypeCache returns the singleton instance of ChannelTypeCache
func GetChannelTypeCache() *ChannelTypeCache {
	channelTypeCacheOnce.Do(func() {
		channelTypeCache = &ChannelTypeCache{
			counts: make(map[int64]int64),
			ttl:    60 * time.Second, // Default TTL: 60 seconds
		}
	})
	return channelTypeCache
}

// Get returns the cached type counts if valid, otherwise returns nil and false
func (c *ChannelTypeCache) Get() (map[int64]int64, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if !c.isValidLocked() {
		return nil, false
	}

	// Return a copy to prevent external modification
	result := make(map[int64]int64, len(c.counts))
	for k, v := range c.counts {
		result[k] = v
	}
	return result, true
}

// Set updates the cache with new type counts
func (c *ChannelTypeCache) Set(counts map[int64]int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.counts = make(map[int64]int64, len(counts))
	for k, v := range counts {
		c.counts[k] = v
	}
	c.lastUpdate = time.Now()
}

// Invalidate clears the cache
func (c *ChannelTypeCache) Invalidate() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.counts = make(map[int64]int64)
	c.lastUpdate = time.Time{} // Zero time
}

// IsValid checks if the cache is still valid (not expired)
func (c *ChannelTypeCache) IsValid() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.isValidLocked()
}

// isValidLocked checks validity without acquiring lock (caller must hold lock)
func (c *ChannelTypeCache) isValidLocked() bool {
	if c.lastUpdate.IsZero() {
		return false
	}
	return time.Since(c.lastUpdate) < c.ttl
}

// SetTTL updates the cache TTL
func (c *ChannelTypeCache) SetTTL(ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.ttl = ttl
}

// GetOrFetch returns cached counts if valid, otherwise fetches from database
func (c *ChannelTypeCache) GetOrFetch() (map[int64]int64, error) {
	// Try to get from cache first
	if counts, ok := c.Get(); ok {
		return counts, nil
	}

	// Fetch from database
	counts, err := model.CountChannelsGroupByType()
	if err != nil {
		return nil, err
	}

	// Update cache
	c.Set(counts)
	return counts, nil
}

// InvalidateChannelTypeCache is a convenience function to invalidate the cache
func InvalidateChannelTypeCache() {
	GetChannelTypeCache().Invalidate()
}
