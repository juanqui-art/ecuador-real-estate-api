package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLRUCache(t *testing.T) {
	tests := []struct {
		name         string
		capacity     int
		maxSize      int64
		ttl          time.Duration
		wantCapacity int
		wantMaxSize  int64
	}{
		{
			name:         "valid parameters",
			capacity:     100,
			maxSize:      1024 * 1024,
			ttl:          1 * time.Hour,
			wantCapacity: 100,
			wantMaxSize:  1024 * 1024,
		},
		{
			name:         "zero capacity should use default",
			capacity:     0,
			maxSize:      1024 * 1024,
			ttl:          1 * time.Hour,
			wantCapacity: 1000,
			wantMaxSize:  1024 * 1024,
		},
		{
			name:         "negative capacity should use default",
			capacity:     -100,
			maxSize:      1024 * 1024,
			ttl:          1 * time.Hour,
			wantCapacity: 1000,
			wantMaxSize:  1024 * 1024,
		},
		{
			name:         "zero max size should use default",
			capacity:     100,
			maxSize:      0,
			ttl:          1 * time.Hour,
			wantCapacity: 100,
			wantMaxSize:  100 * 1024 * 1024,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLRUCache(tt.capacity, tt.maxSize, tt.ttl)

			assert.NotNil(t, cache)
			assert.Equal(t, tt.wantCapacity, cache.capacity)
			assert.Equal(t, tt.wantMaxSize, cache.maxSize)
			assert.Equal(t, tt.ttl, cache.ttl)
			assert.Equal(t, int64(0), cache.currentSize)
			assert.Equal(t, 0, cache.Size())
			assert.NotNil(t, cache.head)
			assert.NotNil(t, cache.tail)
			assert.Equal(t, cache.tail, cache.head.Next)
			assert.Equal(t, cache.head, cache.tail.Prev)
		})
	}
}

func TestLRUCache_SetAndGet(t *testing.T) {
	cache := NewLRUCache(3, 1024, 1*time.Hour)

	// Test setting and getting values
	cache.Set("key1", "value1", 10)
	cache.Set("key2", "value2", 20)
	cache.Set("key3", "value3", 30)

	// Test getting existing values
	val1, found1 := cache.Get("key1")
	assert.True(t, found1)
	assert.Equal(t, "value1", val1)

	val2, found2 := cache.Get("key2")
	assert.True(t, found2)
	assert.Equal(t, "value2", val2)

	val3, found3 := cache.Get("key3")
	assert.True(t, found3)
	assert.Equal(t, "value3", val3)

	// Test getting non-existent value
	val4, found4 := cache.Get("key4")
	assert.False(t, found4)
	assert.Nil(t, val4)

	// Check cache size
	assert.Equal(t, 3, cache.Size())
	assert.Equal(t, int64(60), cache.CurrentSize())
}

func TestLRUCache_CapacityEviction(t *testing.T) {
	cache := NewLRUCache(2, 1024, 1*time.Hour)

	// Fill cache to capacity
	cache.Set("key1", "value1", 10)
	cache.Set("key2", "value2", 20)
	assert.Equal(t, 2, cache.Size())

	// Add third item, should evict least recently used
	cache.Set("key3", "value3", 30)
	assert.Equal(t, 2, cache.Size())

	// key1 should be evicted
	_, found1 := cache.Get("key1")
	assert.False(t, found1)

	// key2 and key3 should still exist
	_, found2 := cache.Get("key2")
	assert.True(t, found2)

	_, found3 := cache.Get("key3")
	assert.True(t, found3)
}

func TestLRUCache_SizeEviction(t *testing.T) {
	cache := NewLRUCache(10, 50, 1*time.Hour) // 50 bytes max

	// Fill cache to size limit
	cache.Set("key1", "value1", 20)
	cache.Set("key2", "value2", 20)
	assert.Equal(t, 2, cache.Size())
	assert.Equal(t, int64(40), cache.CurrentSize())

	// Add item that exceeds size limit
	cache.Set("key3", "value3", 30)
	assert.Equal(t, 2, cache.Size())

	// key1 should be evicted due to size limit
	_, found1 := cache.Get("key1")
	assert.False(t, found1)

	// key2 and key3 should still exist
	_, found2 := cache.Get("key2")
	assert.True(t, found2)

	_, found3 := cache.Get("key3")
	assert.True(t, found3)
}

func TestLRUCache_Update(t *testing.T) {
	cache := NewLRUCache(3, 1024, 1*time.Hour)

	// Set initial value
	cache.Set("key1", "value1", 10)
	assert.Equal(t, int64(10), cache.CurrentSize())

	// Update with different size
	cache.Set("key1", "updated_value1", 20)
	assert.Equal(t, 1, cache.Size())
	assert.Equal(t, int64(20), cache.CurrentSize())

	// Verify updated value
	val, found := cache.Get("key1")
	assert.True(t, found)
	assert.Equal(t, "updated_value1", val)
}

func TestLRUCache_Delete(t *testing.T) {
	cache := NewLRUCache(3, 1024, 1*time.Hour)

	// Set values
	cache.Set("key1", "value1", 10)
	cache.Set("key2", "value2", 20)
	assert.Equal(t, 2, cache.Size())
	assert.Equal(t, int64(30), cache.CurrentSize())

	// Delete existing key
	deleted1 := cache.Delete("key1")
	assert.True(t, deleted1)
	assert.Equal(t, 1, cache.Size())
	assert.Equal(t, int64(20), cache.CurrentSize())

	// Verify deletion
	_, found := cache.Get("key1")
	assert.False(t, found)

	// Delete non-existent key
	deleted2 := cache.Delete("key3")
	assert.False(t, deleted2)
	assert.Equal(t, 1, cache.Size())
}

func TestLRUCache_Clear(t *testing.T) {
	cache := NewLRUCache(3, 1024, 1*time.Hour)

	// Set values
	cache.Set("key1", "value1", 10)
	cache.Set("key2", "value2", 20)
	assert.Equal(t, 2, cache.Size())
	assert.Equal(t, int64(30), cache.CurrentSize())

	// Clear cache
	cache.Clear()
	assert.Equal(t, 0, cache.Size())
	assert.Equal(t, int64(0), cache.CurrentSize())

	// Verify values are gone
	_, found1 := cache.Get("key1")
	assert.False(t, found1)

	_, found2 := cache.Get("key2")
	assert.False(t, found2)
}

func TestLRUCache_LRUOrder(t *testing.T) {
	cache := NewLRUCache(3, 1024, 1*time.Hour)

	// Add items
	cache.Set("key1", "value1", 10)
	cache.Set("key2", "value2", 20)
	cache.Set("key3", "value3", 30)

	// Access key1 to make it recently used
	cache.Get("key1")

	// Add fourth item, should evict key2 (least recently used)
	cache.Set("key4", "value4", 40)

	// key2 should be evicted
	_, found2 := cache.Get("key2")
	assert.False(t, found2)

	// key1, key3, key4 should still exist
	_, found1 := cache.Get("key1")
	assert.True(t, found1)

	_, found3 := cache.Get("key3")
	assert.True(t, found3)

	_, found4 := cache.Get("key4")
	assert.True(t, found4)
}

func TestLRUCache_TTLExpiration(t *testing.T) {
	cache := NewLRUCache(10, 1024, 100*time.Millisecond)

	// Set a value
	cache.Set("key1", "value1", 10)
	
	// Immediately get it (should work)
	val, found := cache.Get("key1")
	assert.True(t, found)
	assert.Equal(t, "value1", val)

	// Wait for TTL to expire
	time.Sleep(150 * time.Millisecond)

	// Try to get expired value
	val, found = cache.Get("key1")
	assert.False(t, found)
	assert.Nil(t, val)

	// Cache should be empty now
	assert.Equal(t, 0, cache.Size())
}

func TestLRUCache_CleanupExpired(t *testing.T) {
	cache := NewLRUCache(10, 1024, 100*time.Millisecond)

	// Set values
	cache.Set("key1", "value1", 10)
	cache.Set("key2", "value2", 20)
	cache.Set("key3", "value3", 30)
	assert.Equal(t, 3, cache.Size())

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Cleanup expired entries
	removed := cache.CleanupExpired()
	assert.Equal(t, 3, removed)
	assert.Equal(t, 0, cache.Size())
}

func TestLRUCache_CleanupExpired_NoTTL(t *testing.T) {
	cache := NewLRUCache(10, 1024, 0) // No TTL

	// Set values
	cache.Set("key1", "value1", 10)
	cache.Set("key2", "value2", 20)

	// Cleanup should remove nothing when TTL is 0
	removed := cache.CleanupExpired()
	assert.Equal(t, 0, removed)
	assert.Equal(t, 2, cache.Size())
}

func TestLRUCache_Stats(t *testing.T) {
	cache := NewLRUCache(10, 1024, 1*time.Hour)

	// Initial stats
	stats := cache.Stats()
	assert.Equal(t, int64(0), stats.Hits)
	assert.Equal(t, int64(0), stats.Misses)
	assert.Equal(t, int64(0), stats.Evictions)
	assert.Equal(t, float64(0), stats.HitRate)
	assert.Equal(t, 0, stats.Size)
	assert.Equal(t, 10, stats.Capacity)
	assert.Equal(t, int64(0), stats.CurrentSize)
	assert.Equal(t, int64(1024), stats.MaxSize)

	// Add some data and test hits/misses
	cache.Set("key1", "value1", 10)
	cache.Get("key1") // Hit
	cache.Get("key2") // Miss

	stats = cache.Stats()
	assert.Equal(t, int64(1), stats.Hits)
	assert.Equal(t, int64(1), stats.Misses)
	assert.Equal(t, float64(50), stats.HitRate)
	assert.Equal(t, 1, stats.Size)
	assert.Equal(t, int64(10), stats.CurrentSize)
}

func TestLRUCache_Keys(t *testing.T) {
	cache := NewLRUCache(10, 1024, 1*time.Hour)

	// Empty cache
	keys := cache.Keys()
	assert.Empty(t, keys)

	// Add some keys
	cache.Set("key1", "value1", 10)
	cache.Set("key2", "value2", 20)
	cache.Set("key3", "value3", 30)

	keys = cache.Keys()
	assert.Len(t, keys, 3)
	assert.Contains(t, keys, "key1")
	assert.Contains(t, keys, "key2")
	assert.Contains(t, keys, "key3")
}

func TestLRUCache_ConcurrentAccess(t *testing.T) {
	cache := NewLRUCache(100, 10240, 1*time.Hour)

	// Test concurrent writes
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				value := fmt.Sprintf("value-%d-%d", id, j)
				cache.Set(key, value, 10)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify cache state
	assert.True(t, cache.Size() > 0)
	assert.True(t, cache.CurrentSize() > 0)
	stats := cache.Stats()
	assert.Equal(t, int64(0), stats.Hits) // No gets were performed
	assert.Equal(t, int64(0), stats.Misses)
}

func TestLRUCache_EdgeCases(t *testing.T) {
	cache := NewLRUCache(1, 100, 1*time.Hour)

	// Test with capacity 1
	cache.Set("key1", "value1", 10)
	assert.Equal(t, 1, cache.Size())

	cache.Set("key2", "value2", 20)
	assert.Equal(t, 1, cache.Size())

	// key1 should be evicted
	_, found1 := cache.Get("key1")
	assert.False(t, found1)

	_, found2 := cache.Get("key2")
	assert.True(t, found2)

	// Test with size limit 1
	cache = NewLRUCache(10, 1, 1*time.Hour)
	cache.Set("key1", "value1", 1)
	assert.Equal(t, 1, cache.Size())

	cache.Set("key2", "value2", 1)
	assert.Equal(t, 1, cache.Size())

	// key1 should be evicted due to size limit
	_, found1 = cache.Get("key1")
	assert.False(t, found1)

	_, found2 = cache.Get("key2")
	assert.True(t, found2)
}

func TestLRUCache_LargeValues(t *testing.T) {
	cache := NewLRUCache(10, 1024, 1*time.Hour)

	// First add some smaller items
	cache.Set("small1", "value1", 100)
	cache.Set("small2", "value2", 200)
	assert.Equal(t, 2, cache.Size())

	// Test with large value that exceeds size limit
	largeValue := make([]byte, 2000)
	cache.Set("large", largeValue, 2000)

	// Large value should not be stored because it exceeds maxSize
	// Small items should still be there
	assert.Equal(t, 2, cache.Size())

	val, found := cache.Get("large")
	assert.False(t, found)
	assert.Nil(t, val)
	
	// Small items should still be accessible
	val1, found1 := cache.Get("small1")
	assert.True(t, found1)
	assert.Equal(t, "value1", val1)
}

func BenchmarkLRUCache_Set(b *testing.B) {
	cache := NewLRUCache(1000, 10240, 1*time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%1000)
		value := fmt.Sprintf("value-%d", i)
		cache.Set(key, value, 10)
	}
}

func BenchmarkLRUCache_Get(b *testing.B) {
	cache := NewLRUCache(1000, 10240, 1*time.Hour)

	// Pre-populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		cache.Set(key, value, 10)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%1000)
		cache.Get(key)
	}
}

func BenchmarkLRUCache_Mixed(b *testing.B) {
	cache := NewLRUCache(1000, 10240, 1*time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%1000)
		if i%2 == 0 {
			value := fmt.Sprintf("value-%d", i)
			cache.Set(key, value, 10)
		} else {
			cache.Get(key)
		}
	}
}