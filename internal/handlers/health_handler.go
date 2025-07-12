package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"realty-core/internal/cache"
	"realty-core/internal/repository"
	"realty-core/internal/service"
)

// HealthHandler handles health check and monitoring endpoints
type HealthHandler struct {
	db           *sql.DB
	propertyRepo repository.PropertyRepository
	imageRepo    repository.ImageRepository
	userRepo     *repository.UserRepository
	agencyRepo   *repository.AgencyRepository
	imageCache   cache.ImageCacheInterface
	propertyService *service.PropertyService
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(
	db *sql.DB,
	propertyRepo repository.PropertyRepository,
	imageRepo repository.ImageRepository,
	userRepo *repository.UserRepository,
	agencyRepo *repository.AgencyRepository,
	imageCache cache.ImageCacheInterface,
	propertyService *service.PropertyService,
) *HealthHandler {
	return &HealthHandler{
		db:              db,
		propertyRepo:    propertyRepo,
		imageRepo:       imageRepo,
		userRepo:        userRepo,
		agencyRepo:      agencyRepo,
		imageCache:      imageCache,
		propertyService: propertyService,
	}
}

// HealthStatus represents the overall health status
type HealthStatus struct {
	Status      string                 `json:"status"`
	Timestamp   time.Time              `json:"timestamp"`
	Version     string                 `json:"version"`
	Uptime      time.Duration          `json:"uptime"`
	Services    map[string]ServiceHealth `json:"services"`
	System      SystemHealth           `json:"system"`
}

// ServiceHealth represents the health of individual services
type ServiceHealth struct {
	Status      string        `json:"status"`
	ResponseTime time.Duration `json:"response_time_ms"`
	Message     string        `json:"message,omitempty"`
	LastChecked time.Time     `json:"last_checked"`
}

// SystemHealth represents system-level health metrics
type SystemHealth struct {
	Memory    MemoryStats `json:"memory"`
	Runtime   RuntimeStats `json:"runtime"`
	Goroutines int         `json:"goroutines"`
}

// MemoryStats represents memory usage statistics
type MemoryStats struct {
	Alloc      uint64 `json:"alloc_bytes"`
	TotalAlloc uint64 `json:"total_alloc_bytes"`
	Sys        uint64 `json:"sys_bytes"`
	NumGC      uint32 `json:"num_gc"`
}

// RuntimeStats represents Go runtime statistics
type RuntimeStats struct {
	Version   string `json:"version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	NumCPU    int    `json:"num_cpu"`
	GOMAXPROCS int   `json:"gomaxprocs"`
}

// CacheHealthStatus represents cache health and statistics
type CacheHealthStatus struct {
	ImageCache    CacheServiceHealth `json:"image_cache"`
	PropertyCache CacheServiceHealth `json:"property_cache"`
}

// CacheServiceHealth represents individual cache service health
type CacheServiceHealth struct {
	Enabled   bool    `json:"enabled"`
	HitRate   float64 `json:"hit_rate_percent"`
	Size      int     `json:"size"`
	Capacity  int     `json:"capacity"`
	Usage     float64 `json:"usage_percent"`
}

var startTime = time.Now()

// BasicHealthCheck provides a simple health check endpoint
func (h *HealthHandler) BasicHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Simple response for basic health check
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"service":   "realty-core",
		"version":   "1.9.0",
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DetailedHealthCheck provides comprehensive health information
func (h *HealthHandler) DetailedHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	startCheck := time.Now()
	
	// Check all services
	services := make(map[string]ServiceHealth)
	
	// Database health
	dbHealth := h.checkDatabaseHealth()
	services["database"] = dbHealth
	
	// Repository health (sample some operations)
	repoHealth := h.checkRepositoryHealth()
	services["repositories"] = repoHealth
	
	// Cache health
	cacheHealth := h.checkCacheHealth()
	services["cache"] = cacheHealth
	
	// Determine overall status
	overallStatus := "healthy"
	for _, service := range services {
		if service.Status != "healthy" {
			overallStatus = "degraded"
			break
		}
	}
	
	// Get system metrics
	systemHealth := h.getSystemHealth()
	
	// Build response
	healthStatus := HealthStatus{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Version:   "1.9.0",
		Uptime:    time.Since(startTime),
		Services:  services,
		System:    systemHealth,
	}
	
	// Set appropriate HTTP status
	statusCode := http.StatusOK
	if overallStatus != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}
	
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(healthStatus)
	
	// Log health check duration if it's slow
	checkDuration := time.Since(startCheck)
	if checkDuration > 1*time.Second {
		// Could log this as a warning
		_ = checkDuration
	}
}

// ReadinessCheck checks if the service is ready to serve traffic
func (h *HealthHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Check critical dependencies
	ready := true
	checks := make(map[string]bool)
	
	// Database connectivity
	if err := h.db.Ping(); err != nil {
		ready = false
		checks["database"] = false
	} else {
		checks["database"] = true
	}
	
	// Could add more readiness checks here
	// - Configuration validation
	// - External service dependencies
	// - Required data initialization
	
	response := map[string]interface{}{
		"ready":     ready,
		"timestamp": time.Now(),
		"checks":    checks,
	}
	
	statusCode := http.StatusOK
	if !ready {
		statusCode = http.StatusServiceUnavailable
	}
	
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// LivenessCheck checks if the service is alive (for Kubernetes)
func (h *HealthHandler) LivenessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Simple liveness check - if we can respond, we're alive
	response := map[string]interface{}{
		"alive":     true,
		"timestamp": time.Now(),
		"uptime":    time.Since(startTime).String(),
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CacheHealth returns detailed cache health and statistics
func (h *HealthHandler) CacheHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	cacheHealth := CacheHealthStatus{}
	
	// Image cache health
	if h.imageCache != nil && h.imageCache.IsEnabled() {
		imageStats := h.imageCache.Stats()
		cacheHealth.ImageCache = CacheServiceHealth{
			Enabled:  true,
			HitRate:  imageStats.HitRate,
			Size:     imageStats.Size,
			Capacity: imageStats.Capacity,
			Usage:    float64(imageStats.Size) / float64(imageStats.Capacity) * 100,
		}
	} else {
		cacheHealth.ImageCache = CacheServiceHealth{
			Enabled: false,
		}
	}
	
	// Property cache health
	if h.propertyService != nil {
		propertyStats := h.propertyService.GetCacheStats()
		totalCapacity := 1000 // Default capacity from cache config
		cacheHealth.PropertyCache = CacheServiceHealth{
			Enabled:  true,
			HitRate:  propertyStats.HitRate,
			Size:     propertyStats.Size,
			Capacity: totalCapacity,
			Usage:    float64(propertyStats.Size) / float64(totalCapacity) * 100,
		}
	} else {
		cacheHealth.PropertyCache = CacheServiceHealth{
			Enabled: false,
		}
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cacheHealth)
}

// MetricsEndpoint provides Prometheus-style metrics
func (h *HealthHandler) MetricsEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	// Basic metrics in Prometheus format
	metrics := []string{
		"# HELP realty_core_uptime_seconds Total uptime in seconds",
		"# TYPE realty_core_uptime_seconds counter",
		fmt.Sprintf("realty_core_uptime_seconds %.2f", time.Since(startTime).Seconds()),
		"",
		"# HELP realty_core_memory_alloc_bytes Current memory allocation in bytes",
		"# TYPE realty_core_memory_alloc_bytes gauge",
		fmt.Sprintf("realty_core_memory_alloc_bytes %d", m.Alloc),
		"",
		"# HELP realty_core_goroutines Current number of goroutines",
		"# TYPE realty_core_goroutines gauge",
		fmt.Sprintf("realty_core_goroutines %d", runtime.NumGoroutine()),
		"",
	}
	
	// Add cache metrics if available
	if h.imageCache != nil && h.imageCache.IsEnabled() {
		stats := h.imageCache.Stats()
		metrics = append(metrics,
			"# HELP realty_core_image_cache_hits_total Total image cache hits",
			"# TYPE realty_core_image_cache_hits_total counter",
			fmt.Sprintf("realty_core_image_cache_hits_total %d", stats.Hits),
			"",
			"# HELP realty_core_image_cache_misses_total Total image cache misses",
			"# TYPE realty_core_image_cache_misses_total counter",
			fmt.Sprintf("realty_core_image_cache_misses_total %d", stats.Misses),
			"",
		)
	}
	
	if h.propertyService != nil {
		propertyStats := h.propertyService.GetCacheStats()
		metrics = append(metrics,
			"# HELP realty_core_property_cache_hits_total Total property cache hits",
			"# TYPE realty_core_property_cache_hits_total counter",
			fmt.Sprintf("realty_core_property_cache_hits_total %d", propertyStats.Hits),
			"",
			"# HELP realty_core_property_cache_search_hits_total Total property search cache hits",
			"# TYPE realty_core_property_cache_search_hits_total counter",
			fmt.Sprintf("realty_core_property_cache_search_hits_total %d", propertyStats.SearchHits),
			"",
		)
	}
	
	for _, metric := range metrics {
		w.Write([]byte(metric + "\n"))
	}
}

// Helper methods

func (h *HealthHandler) checkDatabaseHealth() ServiceHealth {
	start := time.Now()
	
	if err := h.db.Ping(); err != nil {
		return ServiceHealth{
			Status:       "unhealthy",
			ResponseTime: time.Since(start),
			Message:      fmt.Sprintf("Database ping failed: %v", err),
			LastChecked:  time.Now(),
		}
	}
	
	return ServiceHealth{
		Status:       "healthy",
		ResponseTime: time.Since(start),
		LastChecked:  time.Now(),
	}
}

func (h *HealthHandler) checkRepositoryHealth() ServiceHealth {
	start := time.Now()
	
	// Try a simple operation on each repository
	// This is a basic check - could be enhanced with more specific tests
	
	// Test property repository
	if h.propertyRepo != nil {
		// Simple count operation
		if _, err := h.propertyRepo.GetAll(); err != nil {
			return ServiceHealth{
				Status:       "unhealthy",
				ResponseTime: time.Since(start),
				Message:      fmt.Sprintf("Property repository check failed: %v", err),
				LastChecked:  time.Now(),
			}
		}
	}
	
	return ServiceHealth{
		Status:       "healthy",
		ResponseTime: time.Since(start),
		LastChecked:  time.Now(),
	}
}

func (h *HealthHandler) checkCacheHealth() ServiceHealth {
	start := time.Now()
	
	// Check if caches are responding
	healthy := true
	var message string
	
	if h.imageCache != nil && h.imageCache.IsEnabled() {
		// Test image cache
		stats := h.imageCache.Stats()
		if stats.HitRate < 0 { // Basic sanity check
			healthy = false
			message = "Image cache returning invalid stats"
		}
	}
	
	if h.propertyService != nil {
		// Test property cache
		stats := h.propertyService.GetCacheStats()
		if stats.HitRate < 0 { // Basic sanity check
			healthy = false
			message = "Property cache returning invalid stats"
		}
	}
	
	status := "healthy"
	if !healthy {
		status = "degraded"
	}
	
	return ServiceHealth{
		Status:       status,
		ResponseTime: time.Since(start),
		Message:      message,
		LastChecked:  time.Now(),
	}
}

func (h *HealthHandler) getSystemHealth() SystemHealth {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return SystemHealth{
		Memory: MemoryStats{
			Alloc:      m.Alloc,
			TotalAlloc: m.TotalAlloc,
			Sys:        m.Sys,
			NumGC:      m.NumGC,
		},
		Runtime: RuntimeStats{
			Version:    runtime.Version(),
			OS:         runtime.GOOS,
			Arch:       runtime.GOARCH,
			NumCPU:     runtime.NumCPU(),
			GOMAXPROCS: runtime.GOMAXPROCS(0),
		},
		Goroutines: runtime.NumGoroutine(),
	}
}