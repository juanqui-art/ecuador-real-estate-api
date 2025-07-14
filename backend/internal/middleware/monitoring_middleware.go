package middleware

import (
	"net/http"
	"strings"
	"time"

	"realty-core/internal/monitoring"
)

// MonitoringMiddleware provides metrics collection for HTTP requests
func MonitoringMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create response recorder to capture status code
		recorder := &monitoringResponseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		// Execute next handler
		next.ServeHTTP(recorder, r)
		
		// Record metrics
		duration := time.Since(start)
		path := sanitizePath(r.URL.Path)
		
		// Get global metrics collector
		metrics := monitoring.GetGlobalMetrics()
		if metrics != nil {
			metrics.RecordHTTPRequest(r.Method, path, recorder.statusCode, duration)
		}
	})
}

// PerformanceMonitoringMiddleware provides detailed performance monitoring
func PerformanceMonitoringMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create response recorder
		recorder := &monitoringResponseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			bytesWritten:   0,
		}
		
		// Execute next handler
		next.ServeHTTP(recorder, r)
		
		// Record detailed metrics
		duration := time.Since(start)
		path := sanitizePath(r.URL.Path)
		
		metrics := monitoring.GetGlobalMetrics()
		if metrics != nil {
			// Record HTTP request
			metrics.RecordHTTPRequest(r.Method, path, recorder.statusCode, duration)
			
			// Record custom performance metrics
			performanceCounter := metrics.GetOrCreateCounter(
				"http_performance_requests_total",
				"Total requests for performance monitoring",
			)
			performanceCounter.Inc()
			
			// Record response size if significant
			if recorder.bytesWritten > 0 {
				responseSizeHistogram := metrics.GetOrCreateHistogram(
					"http_response_size_bytes",
					"HTTP response size in bytes",
				)
				responseSizeHistogram.Observe(float64(recorder.bytesWritten))
			}
			
			// Record slow requests (> 1 second)
			if duration > time.Second {
				slowRequestCounter := metrics.GetOrCreateCounter(
					"http_slow_requests_total",
					"Total slow HTTP requests (>1s)",
				)
				slowRequestCounter.Inc()
			}
		}
	})
}

// AlertEvaluationMiddleware periodically evaluates alert rules
func AlertEvaluationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Execute next handler first
		next.ServeHTTP(w, r)
		
		// Periodically evaluate alerts (not on every request to avoid overhead)
		// This is a simplified approach - in production, you'd use a separate goroutine
		if shouldEvaluateAlerts() {
			go evaluateAlerts()
		}
	})
}

// monitoringResponseRecorder captures response information for monitoring
type monitoringResponseRecorder struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (mrr *monitoringResponseRecorder) WriteHeader(statusCode int) {
	mrr.statusCode = statusCode
	mrr.ResponseWriter.WriteHeader(statusCode)
}

func (mrr *monitoringResponseRecorder) Write(data []byte) (int, error) {
	n, err := mrr.ResponseWriter.Write(data)
	mrr.bytesWritten += n
	return n, err
}

// sanitizePath sanitizes URL paths for metrics collection
func sanitizePath(path string) string {
	// Remove query parameters
	if idx := strings.Index(path, "?"); idx != -1 {
		path = path[:idx]
	}
	
	// Replace dynamic segments with placeholders
	segments := strings.Split(path, "/")
	for i, segment := range segments {
		if segment == "" {
			continue
		}
		
		// Replace UUIDs and IDs with placeholder
		if isID(segment) {
			segments[i] = "{id}"
		}
	}
	
	sanitized := strings.Join(segments, "/")
	
	// Ensure it starts with /
	if !strings.HasPrefix(sanitized, "/") {
		sanitized = "/" + sanitized
	}
	
	return sanitized
}

// isID checks if a segment looks like an ID (UUID, numeric, etc.)
func isID(segment string) bool {
	if len(segment) == 0 {
		return false
	}
	
	// Check if it's all digits (numeric ID)
	allDigits := true
	for _, r := range segment {
		if r < '0' || r > '9' {
			allDigits = false
			break
		}
	}
	
	if allDigits && len(segment) > 0 {
		return true
	}
	
	// Check if it looks like a UUID (36 chars with hyphens in specific positions)
	if len(segment) == 36 && segment[8] == '-' && segment[13] == '-' && segment[18] == '-' && segment[23] == '-' {
		return true
	}
	
	// Check if it's a hash-like string (long alphanumeric)
	if len(segment) > 20 {
		alphanumeric := true
		for _, r := range segment {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
				alphanumeric = false
				break
			}
		}
		if alphanumeric {
			return true
		}
	}
	
	return false
}

// shouldEvaluateAlerts determines if alerts should be evaluated
// This is a simple rate limiting mechanism
var lastAlertEvaluation time.Time

func shouldEvaluateAlerts() bool {
	now := time.Now()
	if now.Sub(lastAlertEvaluation) > 30*time.Second { // Evaluate every 30 seconds
		lastAlertEvaluation = now
		return true
	}
	return false
}

// evaluateAlerts evaluates alert rules against current metrics
func evaluateAlerts() {
	metrics := monitoring.GetGlobalMetrics()
	alertManager := monitoring.GetGlobalAlertManager()
	
	if metrics != nil && alertManager != nil {
		// Update system metrics before evaluation
		metrics.UpdateSystemMetrics()
		
		// Get current metrics snapshot
		snapshot := metrics.GetMetricsSnapshot()
		
		// Evaluate alert rules
		alertManager.EvaluateRules(&snapshot)
	}
}

// DatabaseMonitoringMiddleware monitors database-related operations
func DatabaseMonitoringMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This middleware would be used to monitor database operations
		// For now, it just passes through to the next handler
		
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		
		// If this was a database-heavy operation, we could record specific metrics
		metrics := monitoring.GetGlobalMetrics()
		if metrics != nil && isDatabaseOperation(r.URL.Path) {
			dbOpCounter := metrics.GetOrCreateCounter(
				"database_operations_total",
				"Total database operations",
			)
			dbOpCounter.Inc()
			
			dbOpDuration := metrics.GetOrCreateHistogram(
				"database_operation_duration",
				"Database operation duration",
			)
			dbOpDuration.Observe(float64(duration.Milliseconds()))
		}
	})
}

// isDatabaseOperation checks if the request path indicates a database operation
func isDatabaseOperation(path string) bool {
	// Simple heuristic - operations that typically involve database queries
	dbPaths := []string{
		"/api/properties",
		"/api/users",
		"/api/agencies",
		"/api/images",
		"/api/pagination",
	}
	
	for _, dbPath := range dbPaths {
		if strings.HasPrefix(path, dbPath) {
			return true
		}
	}
	
	return false
}

// CacheMonitoringMiddleware monitors cache operations
func CacheMonitoringMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create a custom response recorder that can detect cache hits
		recorder := &cacheMonitoringRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		
		// Record cache-related metrics
		metrics := monitoring.GetGlobalMetrics()
		if metrics != nil {
			// If response was very fast (< 10ms), it might be a cache hit
			if duration < 10*time.Millisecond && isCacheableRequest(r) {
				metrics.RecordCacheHit()
			} else if isCacheableRequest(r) {
				metrics.RecordCacheMiss()
			}
		}
	})
}

// cacheMonitoringRecorder monitors cache-related response information
type cacheMonitoringRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (cmr *cacheMonitoringRecorder) WriteHeader(statusCode int) {
	cmr.statusCode = statusCode
	cmr.ResponseWriter.WriteHeader(statusCode)
}

// isCacheableRequest determines if a request is cacheable
func isCacheableRequest(r *http.Request) bool {
	// Only GET requests are typically cacheable
	if r.Method != http.MethodGet {
		return false
	}
	
	// Check for cacheable paths
	cacheablePaths := []string{
		"/api/properties",
		"/api/images",
		"/api/users",
		"/api/agencies",
	}
	
	for _, path := range cacheablePaths {
		if strings.HasPrefix(r.URL.Path, path) {
			return true
		}
	}
	
	return false
}