# Realty Cache Optimization

Optimize caching system for: $ARGUMENTS

## Context - Current Cache Implementation
We have a complete LRU cache system:
- **LRU Cache Core:** Thread-safe with O(1) operations
- **Image Cache:** Specialized for thumbnails and variants
- **Statistics:** Hit/miss rates, memory usage tracking
- **Configuration:** Capacity, size limits, TTL, cleanup intervals

## Cache Architecture:
```go
type LRUCache struct {
    capacity    int           // Max items
    maxSize     int64         // Max memory
    currentSize int64         // Current memory usage
    cache       map[string]*LRUNode
    head        *LRUNode      // Most recently used
    tail        *LRUNode      // Least recently used
    mutex       sync.RWMutex  // Thread safety
    ttl         time.Duration // Time to live
}

type ImageCache struct {
    lru         *LRUCache
    enabled     bool
    stats       tracking for thumbnails/variants
}
```

## Cache Operations:
1. **Performance optimization:**
   - Analyze hit/miss rates
   - Adjust cache size and TTL
   - Implement cache warming strategies
   - Optimize eviction policies

2. **Cache strategies:**
   - **Property search results:** Cache frequent queries
   - **Image variants:** Cache thumbnails and resized images
   - **Popular properties:** Cache most viewed properties
   - **Statistics:** Cache aggregated data

3. **Cache invalidation:**
   - Invalidate on property updates
   - Invalidate related images on property changes
   - Implement cache tags for grouped invalidation
   - Schedule periodic cache cleanup

## Configuration patterns:
```go
// Development
config := ImageCacheConfig{
    Enabled:         true,
    Capacity:        500,
    MaxSizeBytes:    50 * 1024 * 1024, // 50MB
    TTL:             30 * time.Minute,
    CleanupInterval: 5 * time.Minute,
}

// Production
config := ImageCacheConfig{
    Enabled:         true,
    Capacity:        2000,
    MaxSizeBytes:    500 * 1024 * 1024, // 500MB
    TTL:             2 * time.Hour,
    CleanupInterval: 30 * time.Minute,
}
```

## Cache key strategies:
- **Properties:** `property:{id}`, `property_search:{hash}`
- **Images:** `image:{id}_thumbnail_{size}`, `image:{id}_variant_{w}x{h}`
- **Statistics:** `stats:properties:{date}`, `stats:images:{hour}`
- **Search:** `search:{query_hash}_{filters_hash}`

## Performance monitoring:
```go
type CacheStats struct {
    Hits        int64   `json:"hits"`
    Misses      int64   `json:"misses"`
    Evictions   int64   `json:"evictions"`
    HitRate     float64 `json:"hit_rate"`
    Size        int     `json:"size"`
    CurrentSize int64   `json:"current_size"`
}
```

## Cache warming strategies:
- Preload popular properties on startup
- Background warming of search results
- Proactive thumbnail generation
- Cache popular property combinations

## Memory management:
- Monitor memory usage with cache.CurrentSize()
- Implement size-based eviction
- Use efficient data structures
- Regular cleanup of expired entries

## Common optimization scenarios:
- **High miss rate:** Increase cache size or adjust TTL
- **Memory pressure:** Implement better eviction policies
- **Slow queries:** Add caching layer for database results
- **Image processing:** Cache processed variants

## Integration points:
- **Property Service:** Cache search results and property details
- **Image Service:** Cache thumbnails and variants
- **FTS Service:** Cache search results and suggestions
- **Statistics:** Cache aggregated data

## Testing cache performance:
```go
func BenchmarkCacheOperations(b *testing.B) {
    cache := NewImageCache(config)
    
    // Benchmark Set operations
    b.Run("Set", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            cache.Set(fmt.Sprintf("key%d", i), data, "image/jpeg")
        }
    })
    
    // Benchmark Get operations
    b.Run("Get", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            cache.Get(fmt.Sprintf("key%d", i%1000))
        }
    })
}
```

## Cache invalidation patterns:
```go
// Property updated - invalidate related caches
func (s *PropertyService) UpdateProperty(id string, updates ...) error {
    if err := s.repo.Update(id, updates); err != nil {
        return err
    }
    
    // Invalidate property cache
    s.cache.Delete(fmt.Sprintf("property:%s", id))
    
    // Invalidate search results that might include this property
    s.cache.InvalidatePattern("search:*")
    
    return nil
}
```

## Common use cases:
- "optimize cache for property search results"
- "implement cache warming for popular properties"
- "add cache invalidation for property updates"
- "benchmark cache performance under load"
- "reduce cache memory usage"
- "improve cache hit rate for images"

## Monitoring and alerting:
- Track cache hit rates
- Monitor memory usage
- Alert on high miss rates
- Dashboard for cache statistics

## Environment-specific configurations:
- **Development:** Small cache, short TTL for testing
- **Staging:** Production-like cache for performance testing
- **Production:** Optimized cache settings for maximum performance

## Output format:
- Cache configuration optimizations
- Performance benchmarks
- Cache invalidation strategies
- Monitoring and alerting setup
- Memory usage optimizations
- Hit rate improvement recommendations