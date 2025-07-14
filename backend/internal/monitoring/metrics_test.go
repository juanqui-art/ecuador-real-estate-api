package monitoring

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMetricsCollector(t *testing.T) {
	collector := NewMetricsCollector()
	
	assert.NotNil(t, collector)
	assert.NotNil(t, collector.httpRequests)
	assert.NotNil(t, collector.dbConnections)
	assert.NotNil(t, collector.cacheHits)
	assert.True(t, collector.startTime.Before(time.Now()))
}

func TestCounter(t *testing.T) {
	counter := NewCounter("test_counter", "Test counter")
	
	// Test initial value
	assert.Equal(t, int64(0), counter.Get())
	
	// Test increment
	counter.Inc()
	assert.Equal(t, int64(1), counter.Get())
	
	// Test add
	counter.Add(5)
	assert.Equal(t, int64(6), counter.Get())
	
	// Test reset
	counter.Reset()
	assert.Equal(t, int64(0), counter.Get())
}

func TestGauge(t *testing.T) {
	gauge := NewGauge("test_gauge", "Test gauge")
	
	// Test initial value
	assert.Equal(t, float64(0), gauge.Get())
	
	// Test set
	gauge.Set(10.5)
	assert.Equal(t, 10.5, gauge.Get())
	
	// Test add
	gauge.Add(2.5)
	assert.Equal(t, 13.0, gauge.Get())
	
	// Test subtract
	gauge.Sub(3.0)
	assert.Equal(t, 10.0, gauge.Get())
}

func TestHistogram(t *testing.T) {
	histogram := NewHistogram("test_histogram", "Test histogram")
	
	// Test initial values
	assert.Equal(t, int64(0), histogram.GetCount())
	assert.Equal(t, float64(0), histogram.GetSum())
	assert.Equal(t, float64(0), histogram.GetMean())
	
	// Add observations
	histogram.Observe(1.0)
	histogram.Observe(2.0)
	histogram.Observe(3.0)
	
	// Test count and sum
	assert.Equal(t, int64(3), histogram.GetCount())
	assert.Equal(t, 6.0, histogram.GetSum())
	assert.Equal(t, 2.0, histogram.GetMean())
	
	// Test quantiles
	p50 := histogram.GetQuantile(0.5)
	assert.Greater(t, p50, 0.0)
	
	p95 := histogram.GetQuantile(0.95)
	assert.Greater(t, p95, 0.0)
}

func TestMetricsCollector_RecordHTTPRequest(t *testing.T) {
	collector := NewMetricsCollector()
	
	// Record some HTTP requests
	collector.RecordHTTPRequest("GET", "/api/test", 200, 100*time.Millisecond)
	collector.RecordHTTPRequest("GET", "/api/test", 404, 50*time.Millisecond)
	collector.RecordHTTPRequest("POST", "/api/test", 201, 200*time.Millisecond)
	
	snapshot := collector.GetMetricsSnapshot()
	
	// Check HTTP metrics exist
	assert.Greater(t, len(snapshot.HTTP), 0)
	
	// Check specific metrics
	getMetric, exists := snapshot.HTTP["GET_/api/test"]
	assert.True(t, exists)
	assert.Equal(t, int64(2), getMetric.Requests)
	assert.Greater(t, getMetric.AvgDuration, 0.0)
}

func TestMetricsCollector_DatabaseMetrics(t *testing.T) {
	collector := NewMetricsCollector()
	
	// Record database operations
	collector.RecordDBConnection(5)
	collector.RecordDBQuery(50*time.Millisecond, false)
	collector.RecordDBQuery(100*time.Millisecond, true)
	
	snapshot := collector.GetMetricsSnapshot()
	
	// Check database metrics
	assert.Equal(t, float64(5), snapshot.Database.Connections)
	assert.Equal(t, int64(2), snapshot.Database.Queries)
	assert.Equal(t, int64(1), snapshot.Database.Errors)
	assert.Greater(t, snapshot.Database.QueryDuration, 0.0)
}

func TestMetricsCollector_CacheMetrics(t *testing.T) {
	collector := NewMetricsCollector()
	
	// Record cache operations
	collector.RecordCacheHit()
	collector.RecordCacheHit()
	collector.RecordCacheMiss()
	collector.RecordCacheEviction()
	collector.RecordCacheSize(100)
	
	snapshot := collector.GetMetricsSnapshot()
	
	// Check cache metrics
	assert.Equal(t, int64(2), snapshot.Cache.Hits)
	assert.Equal(t, int64(1), snapshot.Cache.Misses)
	assert.Equal(t, int64(1), snapshot.Cache.Evictions)
	assert.Equal(t, int64(100), snapshot.Cache.Size)
	assert.Equal(t, float64(2)/float64(3)*100, snapshot.Cache.HitRate)
}

func TestMetricsCollector_BusinessMetrics(t *testing.T) {
	collector := NewMetricsCollector()
	
	// Update business metrics
	collector.UpdateEntityCounts(50, 200, 25, 10)
	
	snapshot := collector.GetMetricsSnapshot()
	
	// Check business metrics
	assert.Equal(t, int64(50), snapshot.Business.Properties)
	assert.Equal(t, int64(200), snapshot.Business.Images)
	assert.Equal(t, int64(25), snapshot.Business.Users)
	assert.Equal(t, int64(10), snapshot.Business.Agencies)
}

func TestMetricsCollector_CustomMetrics(t *testing.T) {
	collector := NewMetricsCollector()
	
	// Create custom metrics
	customCounter := collector.GetOrCreateCounter("custom_counter", "Custom counter")
	customGauge := collector.GetOrCreateGauge("custom_gauge", "Custom gauge")
	customHistogram := collector.GetOrCreateHistogram("custom_histogram", "Custom histogram")
	
	// Test that metrics are created and reused
	assert.NotNil(t, customCounter)
	assert.NotNil(t, customGauge)
	assert.NotNil(t, customHistogram)
	
	// Test reuse
	sameCounter := collector.GetOrCreateCounter("custom_counter", "Custom counter")
	assert.Equal(t, customCounter, sameCounter)
}

func TestMetricsCollector_SystemMetrics(t *testing.T) {
	collector := NewMetricsCollector()
	
	// Update system metrics
	collector.UpdateSystemMetrics()
	
	snapshot := collector.GetMetricsSnapshot()
	
	// Check system metrics are populated
	assert.Greater(t, snapshot.System.Memory, int64(0))
	assert.Greater(t, snapshot.System.Goroutines, 0)
	assert.NotZero(t, snapshot.Uptime)
}

func TestCalculateHitRate(t *testing.T) {
	testCases := []struct {
		hits     int64
		misses   int64
		expected float64
	}{
		{80, 20, 80.0},
		{0, 0, 0.0},
		{100, 0, 100.0},
		{0, 100, 0.0},
		{1, 1, 50.0},
	}
	
	for _, tc := range testCases {
		result := calculateHitRate(tc.hits, tc.misses)
		assert.Equal(t, tc.expected, result)
	}
}

func TestMetricsSnapshot(t *testing.T) {
	collector := NewMetricsCollector()
	
	// Add some data
	collector.RecordHTTPRequest("GET", "/test", 200, 100*time.Millisecond)
	collector.RecordCacheHit()
	collector.UpdateEntityCounts(10, 20, 5, 2)
	
	snapshot := collector.GetMetricsSnapshot()
	
	// Check snapshot structure
	assert.NotZero(t, snapshot.Timestamp)
	assert.Greater(t, snapshot.Uptime, time.Duration(0))
	assert.NotNil(t, snapshot.HTTP)
	assert.NotNil(t, snapshot.Database)
	assert.NotNil(t, snapshot.Cache)
	assert.NotNil(t, snapshot.Business)
	assert.NotNil(t, snapshot.System)
}

// Benchmark tests
func BenchmarkCounter_Inc(b *testing.B) {
	counter := NewCounter("bench_counter", "Benchmark counter")
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Inc()
		}
	})
}

func BenchmarkGauge_Set(b *testing.B) {
	gauge := NewGauge("bench_gauge", "Benchmark gauge")
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gauge.Set(float64(b.N))
		}
	})
}

func BenchmarkHistogram_Observe(b *testing.B) {
	histogram := NewHistogram("bench_histogram", "Benchmark histogram")
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			histogram.Observe(float64(b.N % 1000))
		}
	})
}

func BenchmarkMetricsCollector_RecordHTTPRequest(b *testing.B) {
	collector := NewMetricsCollector()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector.RecordHTTPRequest("GET", "/api/test", 200, time.Millisecond)
	}
}

func BenchmarkMetricsCollector_GetSnapshot(b *testing.B) {
	collector := NewMetricsCollector()
	
	// Add some data
	for i := 0; i < 100; i++ {
		collector.RecordHTTPRequest("GET", "/api/test", 200, time.Millisecond)
		collector.RecordCacheHit()
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector.GetMetricsSnapshot()
	}
}