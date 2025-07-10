# Realty Debug & Performance

Debug and optimize performance for: $ARGUMENTS

## Context - System Performance
Monitor and optimize:
- **Database queries:** PostgreSQL performance
- **Cache performance:** Hit rates and memory usage
- **Image processing:** Resize and compression times
- **API response times:** Endpoint latency
- **Memory usage:** Go heap and garbage collection

## Debug Strategies:
1. **Database performance:**
   - Analyze slow queries with EXPLAIN
   - Monitor connection pool usage
   - Check index usage
   - Identify N+1 query problems

2. **Cache optimization:**
   - Monitor hit/miss rates
   - Analyze memory usage
   - Check eviction patterns
   - Optimize cache keys

3. **API performance:**
   - Profile endpoint response times
   - Monitor concurrent requests
   - Check error rates
   - Analyze payload sizes

## Performance tools:
```go
// Database query profiling
func (r *Repository) debugQuery(query string, args ...interface{}) {
    start := time.Now()
    defer func() {
        log.Printf("Query took %v: %s", time.Since(start), query)
    }()
    // Execute query
}

// Cache performance monitoring
func (c *Cache) logStats() {
    stats := c.Stats()
    log.Printf("Cache stats: hits=%d, misses=%d, hit_rate=%.2f%%", 
        stats.Hits, stats.Misses, stats.HitRate)
}
```

## Common debug scenarios:
- "slow property search queries"
- "cache miss rate too high"
- "image processing memory leak"
- "API endpoint timeouts"
- "database connection pool exhaustion"

## Monitoring setup:
- pprof profiling endpoints
- Prometheus metrics
- Grafana dashboards
- Log aggregation
- Error tracking

## Output format:
- Performance analysis
- Optimization recommendations
- Monitoring setup
- Debug tools
- Profiling results