package monitoring

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// MetricsCollector collects and aggregates application metrics
type MetricsCollector struct {
	mutex sync.RWMutex
	
	// HTTP metrics
	httpRequests    map[string]*Counter
	httpDurations   map[string]*Histogram
	httpErrors      map[string]*Counter
	
	// Database metrics
	dbConnections   *Gauge
	dbQueries       *Counter
	dbQueryDuration *Histogram
	dbErrors        *Counter
	
	// Cache metrics
	cacheHits       *Counter
	cacheMisses     *Counter
	cacheEvictions  *Counter
	cacheSize       *Gauge
	
	// Business metrics
	propertiesCount *Gauge
	imagesCount     *Gauge
	usersCount      *Gauge
	agenciesCount   *Gauge
	
	// System metrics
	systemMemory    *Gauge
	systemCPU       *Gauge
	goroutines      *Gauge
	
	// Custom metrics
	customCounters   map[string]*Counter
	customGauges     map[string]*Gauge
	customHistograms map[string]*Histogram
	
	startTime time.Time
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		httpRequests:     make(map[string]*Counter),
		httpDurations:    make(map[string]*Histogram),
		httpErrors:       make(map[string]*Counter),
		customCounters:   make(map[string]*Counter),
		customGauges:     make(map[string]*Gauge),
		customHistograms: make(map[string]*Histogram),
		
		// Initialize standard metrics
		dbConnections:   NewGauge("db_connections", "Number of active database connections"),
		dbQueries:       NewCounter("db_queries_total", "Total number of database queries"),
		dbQueryDuration: NewHistogram("db_query_duration", "Database query duration in milliseconds"),
		dbErrors:        NewCounter("db_errors_total", "Total number of database errors"),
		
		cacheHits:      NewCounter("cache_hits_total", "Total number of cache hits"),
		cacheMisses:    NewCounter("cache_misses_total", "Total number of cache misses"),
		cacheEvictions: NewCounter("cache_evictions_total", "Total number of cache evictions"),
		cacheSize:      NewGauge("cache_size", "Current cache size"),
		
		propertiesCount: NewGauge("properties_count", "Current number of properties"),
		imagesCount:     NewGauge("images_count", "Current number of images"),
		usersCount:      NewGauge("users_count", "Current number of users"),
		agenciesCount:   NewGauge("agencies_count", "Current number of agencies"),
		
		systemMemory: NewGauge("system_memory_bytes", "System memory usage in bytes"),
		systemCPU:    NewGauge("system_cpu_percent", "System CPU usage percentage"),
		goroutines:   NewGauge("goroutines", "Number of goroutines"),
		
		startTime: time.Now(),
	}
}

// Counter represents a monotonically increasing counter
type Counter struct {
	name        string
	description string
	value       int64
	mutex       sync.RWMutex
}

// NewCounter creates a new counter
func NewCounter(name, description string) *Counter {
	return &Counter{
		name:        name,
		description: description,
		value:       0,
	}
}

// Inc increments the counter by 1
func (c *Counter) Inc() {
	c.Add(1)
}

// Add adds the given value to the counter
func (c *Counter) Add(value int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.value += value
}

// Get returns the current counter value
func (c *Counter) Get() int64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.value
}

// Reset resets the counter to zero
func (c *Counter) Reset() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.value = 0
}

// Gauge represents a metric that can go up and down
type Gauge struct {
	name        string
	description string
	value       float64
	mutex       sync.RWMutex
}

// NewGauge creates a new gauge
func NewGauge(name, description string) *Gauge {
	return &Gauge{
		name:        name,
		description: description,
		value:       0,
	}
}

// Set sets the gauge to the given value
func (g *Gauge) Set(value float64) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.value = value
}

// Add adds the given value to the gauge
func (g *Gauge) Add(value float64) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.value += value
}

// Sub subtracts the given value from the gauge
func (g *Gauge) Sub(value float64) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.value -= value
}

// Get returns the current gauge value
func (g *Gauge) Get() float64 {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	return g.value
}

// Histogram represents a metric for measuring distributions
type Histogram struct {
	name        string
	description string
	buckets     []float64
	counts      []int64
	sum         float64
	count       int64
	mutex       sync.RWMutex
}

// NewHistogram creates a new histogram with default buckets
func NewHistogram(name, description string) *Histogram {
	buckets := []float64{
		0.5, 1, 2.5, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000,
	}
	
	return &Histogram{
		name:        name,
		description: description,
		buckets:     buckets,
		counts:      make([]int64, len(buckets)+1), // +1 for infinity bucket
	}
}

// Observe adds an observation to the histogram
func (h *Histogram) Observe(value float64) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	h.sum += value
	h.count++
	
	// Find the bucket for this value
	for i, bucket := range h.buckets {
		if value <= bucket {
			h.counts[i]++
			return
		}
	}
	
	// Value is greater than all buckets (infinity bucket)
	h.counts[len(h.buckets)]++
}

// GetQuantile returns the approximate quantile value
func (h *Histogram) GetQuantile(quantile float64) float64 {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	if h.count == 0 {
		return 0
	}
	
	targetCount := float64(h.count) * quantile
	runningCount := int64(0)
	
	for i, count := range h.counts {
		runningCount += count
		if float64(runningCount) >= targetCount {
			if i == len(h.buckets) {
				return h.buckets[len(h.buckets)-1] // Return largest bucket
			}
			return h.buckets[i]
		}
	}
	
	return 0
}

// GetMean returns the mean of all observations
func (h *Histogram) GetMean() float64 {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	if h.count == 0 {
		return 0
	}
	
	return h.sum / float64(h.count)
}

// GetCount returns the total number of observations
func (h *Histogram) GetCount() int64 {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.count
}

// GetSum returns the sum of all observations
func (h *Histogram) GetSum() float64 {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.sum
}

// HTTP Metrics Methods

// RecordHTTPRequest records HTTP request metrics
func (m *MetricsCollector) RecordHTTPRequest(method, path string, statusCode int, duration time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	// Create metric keys
	requestKey := fmt.Sprintf("%s_%s", method, path)
	errorKey := fmt.Sprintf("%s_%s_%d", method, path, statusCode)
	
	// Initialize metrics if they don't exist
	if _, exists := m.httpRequests[requestKey]; !exists {
		m.httpRequests[requestKey] = NewCounter(
			fmt.Sprintf("http_requests_total_%s_%s", method, path),
			fmt.Sprintf("Total HTTP requests for %s %s", method, path),
		)
	}
	
	if _, exists := m.httpDurations[requestKey]; !exists {
		m.httpDurations[requestKey] = NewHistogram(
			fmt.Sprintf("http_request_duration_%s_%s", method, path),
			fmt.Sprintf("HTTP request duration for %s %s", method, path),
		)
	}
	
	// Record metrics
	m.httpRequests[requestKey].Inc()
	m.httpDurations[requestKey].Observe(float64(duration.Milliseconds()))
	
	// Record errors (4xx and 5xx status codes)
	if statusCode >= 400 {
		if _, exists := m.httpErrors[errorKey]; !exists {
			m.httpErrors[errorKey] = NewCounter(
				fmt.Sprintf("http_errors_total_%s_%s_%d", method, path, statusCode),
				fmt.Sprintf("Total HTTP errors for %s %s with status %d", method, path, statusCode),
			)
		}
		m.httpErrors[errorKey].Inc()
	}
}

// Database Metrics Methods

// RecordDBConnection records database connection metrics
func (m *MetricsCollector) RecordDBConnection(activeConnections int) {
	m.dbConnections.Set(float64(activeConnections))
}

// RecordDBQuery records database query metrics
func (m *MetricsCollector) RecordDBQuery(duration time.Duration, isError bool) {
	m.dbQueries.Inc()
	m.dbQueryDuration.Observe(float64(duration.Milliseconds()))
	
	if isError {
		m.dbErrors.Inc()
	}
}

// Cache Metrics Methods

// RecordCacheHit records a cache hit
func (m *MetricsCollector) RecordCacheHit() {
	m.cacheHits.Inc()
}

// RecordCacheMiss records a cache miss
func (m *MetricsCollector) RecordCacheMiss() {
	m.cacheMisses.Inc()
}

// RecordCacheEviction records a cache eviction
func (m *MetricsCollector) RecordCacheEviction() {
	m.cacheEvictions.Inc()
}

// RecordCacheSize records the current cache size
func (m *MetricsCollector) RecordCacheSize(size int) {
	m.cacheSize.Set(float64(size))
}

// Business Metrics Methods

// UpdateEntityCounts updates business entity counts
func (m *MetricsCollector) UpdateEntityCounts(properties, images, users, agencies int) {
	m.propertiesCount.Set(float64(properties))
	m.imagesCount.Set(float64(images))
	m.usersCount.Set(float64(users))
	m.agenciesCount.Set(float64(agencies))
}

// System Metrics Methods

// UpdateSystemMetrics updates system-level metrics
func (m *MetricsCollector) UpdateSystemMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	m.systemMemory.Set(float64(memStats.Alloc))
	m.goroutines.Set(float64(runtime.NumGoroutine()))
	
	// CPU usage would require additional system calls
	// For now, we'll use a placeholder
	m.systemCPU.Set(0) // TODO: Implement actual CPU monitoring
}

// Custom Metrics Methods

// GetOrCreateCounter gets or creates a custom counter
func (m *MetricsCollector) GetOrCreateCounter(name, description string) *Counter {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if counter, exists := m.customCounters[name]; exists {
		return counter
	}
	
	counter := NewCounter(name, description)
	m.customCounters[name] = counter
	return counter
}

// GetOrCreateGauge gets or creates a custom gauge
func (m *MetricsCollector) GetOrCreateGauge(name, description string) *Gauge {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if gauge, exists := m.customGauges[name]; exists {
		return gauge
	}
	
	gauge := NewGauge(name, description)
	m.customGauges[name] = gauge
	return gauge
}

// GetOrCreateHistogram gets or creates a custom histogram
func (m *MetricsCollector) GetOrCreateHistogram(name, description string) *Histogram {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if histogram, exists := m.customHistograms[name]; exists {
		return histogram
	}
	
	histogram := NewHistogram(name, description)
	m.customHistograms[name] = histogram
	return histogram
}

// Reporting Methods

// GetMetricsSnapshot returns a snapshot of all metrics
func (m *MetricsCollector) GetMetricsSnapshot() MetricsSnapshot {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	snapshot := MetricsSnapshot{
		Timestamp: time.Now(),
		Uptime:    time.Since(m.startTime),
		HTTP:      make(map[string]HTTPMetric),
		Database: DatabaseMetrics{
			Connections:     m.dbConnections.Get(),
			Queries:         m.dbQueries.Get(),
			QueryDuration:   m.dbQueryDuration.GetMean(),
			Errors:          m.dbErrors.Get(),
		},
		Cache: CacheMetrics{
			Hits:       m.cacheHits.Get(),
			Misses:     m.cacheMisses.Get(),
			Evictions:  m.cacheEvictions.Get(),
			Size:       int64(m.cacheSize.Get()),
			HitRate:    calculateHitRate(m.cacheHits.Get(), m.cacheMisses.Get()),
		},
		Business: BusinessMetrics{
			Properties: int64(m.propertiesCount.Get()),
			Images:     int64(m.imagesCount.Get()),
			Users:      int64(m.usersCount.Get()),
			Agencies:   int64(m.agenciesCount.Get()),
		},
		System: SystemMetrics{
			Memory:     int64(m.systemMemory.Get()),
			CPU:        m.systemCPU.Get(),
			Goroutines: int(m.goroutines.Get()),
		},
	}
	
	// Collect HTTP metrics
	for key, counter := range m.httpRequests {
		if duration, exists := m.httpDurations[key]; exists {
			snapshot.HTTP[key] = HTTPMetric{
				Requests:      counter.Get(),
				AvgDuration:   duration.GetMean(),
				P95Duration:   duration.GetQuantile(0.95),
				P99Duration:   duration.GetQuantile(0.99),
			}
		}
	}
	
	return snapshot
}

// MetricsSnapshot contains a snapshot of all metrics at a point in time
type MetricsSnapshot struct {
	Timestamp time.Time               `json:"timestamp"`
	Uptime    time.Duration           `json:"uptime"`
	HTTP      map[string]HTTPMetric   `json:"http"`
	Database  DatabaseMetrics         `json:"database"`
	Cache     CacheMetrics            `json:"cache"`
	Business  BusinessMetrics         `json:"business"`
	System    SystemMetrics           `json:"system"`
}

// HTTPMetric contains HTTP-related metrics
type HTTPMetric struct {
	Requests    int64   `json:"requests"`
	AvgDuration float64 `json:"avg_duration_ms"`
	P95Duration float64 `json:"p95_duration_ms"`
	P99Duration float64 `json:"p99_duration_ms"`
}

// DatabaseMetrics contains database-related metrics
type DatabaseMetrics struct {
	Connections   float64 `json:"connections"`
	Queries       int64   `json:"queries"`
	QueryDuration float64 `json:"query_duration_ms"`
	Errors        int64   `json:"errors"`
}

// CacheMetrics contains cache-related metrics
type CacheMetrics struct {
	Hits      int64   `json:"hits"`
	Misses    int64   `json:"misses"`
	Evictions int64   `json:"evictions"`
	Size      int64   `json:"size"`
	HitRate   float64 `json:"hit_rate"`
}

// BusinessMetrics contains business-related metrics
type BusinessMetrics struct {
	Properties int64 `json:"properties"`
	Images     int64 `json:"images"`
	Users      int64 `json:"users"`
	Agencies   int64 `json:"agencies"`
}

// SystemMetrics contains system-related metrics
type SystemMetrics struct {
	Memory     int64   `json:"memory_bytes"`
	CPU        float64 `json:"cpu_percent"`
	Goroutines int     `json:"goroutines"`
}

// calculateHitRate calculates cache hit rate
func calculateHitRate(hits, misses int64) float64 {
	total := hits + misses
	if total == 0 {
		return 0
	}
	return float64(hits) / float64(total) * 100
}

// Global metrics collector instance
var globalMetrics *MetricsCollector

// InitializeMetrics initializes the global metrics collector
func InitializeMetrics() {
	globalMetrics = NewMetricsCollector()
}

// GetGlobalMetrics returns the global metrics collector
func GetGlobalMetrics() *MetricsCollector {
	return globalMetrics
}