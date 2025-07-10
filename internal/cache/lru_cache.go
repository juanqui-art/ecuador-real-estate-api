package cache

import (
	"sync"
	"time"
)

// LRUNode represents a node in the doubly linked list
type LRUNode struct {
	Key       string
	Value     interface{}
	Size      int64
	Timestamp time.Time
	Prev      *LRUNode
	Next      *LRUNode
}

// LRUCache implements a thread-safe LRU cache with size limits
type LRUCache struct {
	capacity    int           // Maximum number of items
	maxSize     int64         // Maximum total size in bytes
	currentSize int64         // Current total size in bytes
	cache       map[string]*LRUNode
	head        *LRUNode      // Most recently used
	tail        *LRUNode      // Least recently used
	mutex       sync.RWMutex
	hits        int64
	misses      int64
	evictions   int64
	ttl         time.Duration // Time to live for entries
}

// NewLRUCache creates a new LRU cache with specified capacity and size limits
func NewLRUCache(capacity int, maxSizeBytes int64, ttl time.Duration) *LRUCache {
	if capacity <= 0 {
		capacity = 1000
	}
	if maxSizeBytes <= 0 {
		maxSizeBytes = 100 * 1024 * 1024 // 100MB default
	}
	
	cache := &LRUCache{
		capacity:    capacity,
		maxSize:     maxSizeBytes,
		currentSize: 0,
		cache:       make(map[string]*LRUNode),
		hits:        0,
		misses:      0,
		evictions:   0,
		ttl:         ttl,
	}
	
	// Initialize dummy head and tail nodes
	cache.head = &LRUNode{}
	cache.tail = &LRUNode{}
	cache.head.Next = cache.tail
	cache.tail.Prev = cache.head
	
	return cache
}

// Get retrieves a value from the cache
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	if node, exists := c.cache[key]; exists {
		// Check if entry has expired
		if c.ttl > 0 && time.Since(node.Timestamp) > c.ttl {
			c.removeNode(node)
			delete(c.cache, key)
			c.misses++
			return nil, false
		}
		
		// Move to front (most recently used)
		c.moveToFront(node)
		c.hits++
		return node.Value, true
	}
	
	c.misses++
	return nil, false
}

// Set adds or updates a value in the cache
func (c *LRUCache) Set(key string, value interface{}, size int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	// Don't store items that exceed max size on their own
	if size > c.maxSize {
		// If it's an update of an existing key, remove the old entry
		if node, exists := c.cache[key]; exists {
			c.currentSize -= node.Size
			c.removeNode(node)
			delete(c.cache, key)
		}
		return
	}
	
	if node, exists := c.cache[key]; exists {
		// Update existing node
		oldSize := node.Size
		node.Value = value
		node.Size = size
		node.Timestamp = time.Now()
		
		c.currentSize = c.currentSize - oldSize + size
		c.moveToFront(node)
		
		// Check if we need to evict due to size limit
		c.evictIfNeeded()
		return
	}
	
	// Create new node
	newNode := &LRUNode{
		Key:       key,
		Value:     value,
		Size:      size,
		Timestamp: time.Now(),
	}
	
	// Add to cache
	c.cache[key] = newNode
	c.currentSize += size
	c.addToFront(newNode)
	
	// Check if we need to evict
	c.evictIfNeeded()
}

// Delete removes a key from the cache
func (c *LRUCache) Delete(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	if node, exists := c.cache[key]; exists {
		c.currentSize -= node.Size
		c.removeNode(node)
		delete(c.cache, key)
		return true
	}
	
	return false
}

// Clear removes all entries from the cache
func (c *LRUCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.cache = make(map[string]*LRUNode)
	c.currentSize = 0
	c.head.Next = c.tail
	c.tail.Prev = c.head
	c.hits = 0
	c.misses = 0
	c.evictions = 0
}

// Size returns the current number of items in the cache
func (c *LRUCache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.cache)
}

// CurrentSize returns the current total size in bytes
func (c *LRUCache) CurrentSize() int64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.currentSize
}

// Stats returns cache statistics
func (c *LRUCache) Stats() CacheStats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	total := c.hits + c.misses
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(c.hits) / float64(total) * 100
	}
	
	return CacheStats{
		Hits:        c.hits,
		Misses:      c.misses,
		Evictions:   c.evictions,
		HitRate:     hitRate,
		Size:        len(c.cache),
		Capacity:    c.capacity,
		CurrentSize: c.currentSize,
		MaxSize:     c.maxSize,
	}
}

// Keys returns all keys in the cache (for debugging)
func (c *LRUCache) Keys() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	keys := make([]string, 0, len(c.cache))
	for key := range c.cache {
		keys = append(keys, key)
	}
	return keys
}

// CleanupExpired removes expired entries
func (c *LRUCache) CleanupExpired() int {
	if c.ttl <= 0 {
		return 0
	}
	
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	removed := 0
	cutoff := time.Now().Add(-c.ttl)
	
	// Start from tail (least recently used) and remove expired entries
	current := c.tail.Prev
	for current != c.head {
		if current.Timestamp.Before(cutoff) {
			prev := current.Prev
			c.currentSize -= current.Size
			c.removeNode(current)
			delete(c.cache, current.Key)
			removed++
			current = prev
		} else {
			// Since we're going from LRU to MRU, if this one isn't expired,
			// the rest won't be either
			break
		}
	}
	
	return removed
}

// evictIfNeeded removes entries if capacity or size limits are exceeded
func (c *LRUCache) evictIfNeeded() {
	// Evict by capacity
	for len(c.cache) > c.capacity {
		c.evictLRU()
	}
	
	// Evict by size
	for c.currentSize > c.maxSize {
		c.evictLRU()
	}
}

// evictLRU removes the least recently used item
func (c *LRUCache) evictLRU() {
	if c.tail.Prev == c.head {
		return // Cache is empty
	}
	
	lru := c.tail.Prev
	c.currentSize -= lru.Size
	c.removeNode(lru)
	delete(c.cache, lru.Key)
	c.evictions++
}

// moveToFront moves a node to the front of the list
func (c *LRUCache) moveToFront(node *LRUNode) {
	c.removeNode(node)
	c.addToFront(node)
}

// addToFront adds a node to the front of the list
func (c *LRUCache) addToFront(node *LRUNode) {
	node.Prev = c.head
	node.Next = c.head.Next
	c.head.Next.Prev = node
	c.head.Next = node
}

// removeNode removes a node from the list
func (c *LRUCache) removeNode(node *LRUNode) {
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
}

// CacheStats represents cache statistics
type CacheStats struct {
	Hits        int64   `json:"hits"`
	Misses      int64   `json:"misses"`
	Evictions   int64   `json:"evictions"`
	HitRate     float64 `json:"hit_rate"`
	Size        int     `json:"size"`
	Capacity    int     `json:"capacity"`
	CurrentSize int64   `json:"current_size"`
	MaxSize     int64   `json:"max_size"`
}

// CacheInterface defines the interface for cache operations
type CacheInterface interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, size int64)
	Delete(key string) bool
	Clear()
	Size() int
	CurrentSize() int64
	Stats() CacheStats
	Keys() []string
	CleanupExpired() int
}