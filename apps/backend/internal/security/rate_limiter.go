package security

import (
	"sync"
	"time"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	buckets map[string]*TokenBucket
	mutex   sync.RWMutex
	
	// Configuration
	maxRequests    int           // Maximum requests per window
	windowDuration time.Duration // Time window duration
	cleanupTicker  *time.Ticker  // Cleanup timer
}

// TokenBucket represents a token bucket for rate limiting
type TokenBucket struct {
	tokens       int       // Current number of tokens
	maxTokens    int       // Maximum number of tokens
	refillRate   int       // Tokens added per refill interval
	lastRefill   time.Time // Last refill time
	mutex        sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxRequests int, windowDuration time.Duration) *RateLimiter {
	rl := &RateLimiter{
		buckets:        make(map[string]*TokenBucket),
		maxRequests:    maxRequests,
		windowDuration: windowDuration,
	}
	
	// Start cleanup routine
	rl.cleanupTicker = time.NewTicker(windowDuration)
	go rl.cleanupRoutine()
	
	return rl
}

// Allow checks if a request is allowed for the given identifier
func (rl *RateLimiter) Allow(identifier string) bool {
	rl.mutex.Lock()
	bucket, exists := rl.buckets[identifier]
	if !exists {
		bucket = &TokenBucket{
			tokens:     rl.maxRequests - 1, // Allow this request
			maxTokens:  rl.maxRequests,
			refillRate: rl.maxRequests,
			lastRefill: time.Now(),
		}
		rl.buckets[identifier] = bucket
		rl.mutex.Unlock()
		return true
	}
	rl.mutex.Unlock()
	
	return bucket.takeToken()
}

// takeToken attempts to take a token from the bucket
func (tb *TokenBucket) takeToken() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	
	// Refill tokens based on time elapsed
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	if elapsed > time.Minute {
		// Reset tokens after a minute
		tb.tokens = tb.maxTokens
		tb.lastRefill = now
	}
	
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	
	return false
}

// cleanupRoutine removes old buckets to prevent memory leaks
func (rl *RateLimiter) cleanupRoutine() {
	for range rl.cleanupTicker.C {
		rl.cleanup()
	}
}

// cleanup removes buckets that haven't been used recently
func (rl *RateLimiter) cleanup() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	cutoff := time.Now().Add(-rl.windowDuration * 2)
	
	for identifier, bucket := range rl.buckets {
		bucket.mutex.Lock()
		if bucket.lastRefill.Before(cutoff) {
			delete(rl.buckets, identifier)
		}
		bucket.mutex.Unlock()
	}
}

// Stop stops the rate limiter and cleanup routine
func (rl *RateLimiter) Stop() {
	if rl.cleanupTicker != nil {
		rl.cleanupTicker.Stop()
	}
}

// GetStats returns rate limiter statistics
func (rl *RateLimiter) GetStats() RateLimiterStats {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()
	
	return RateLimiterStats{
		ActiveBuckets:  len(rl.buckets),
		MaxRequests:    rl.maxRequests,
		WindowDuration: rl.windowDuration,
	}
}

// RateLimiterStats contains rate limiter statistics
type RateLimiterStats struct {
	ActiveBuckets  int           `json:"active_buckets"`
	MaxRequests    int           `json:"max_requests"`
	WindowDuration time.Duration `json:"window_duration"`
}

// AdaptiveRateLimiter implements an adaptive rate limiter that adjusts based on server load
type AdaptiveRateLimiter struct {
	baseLimiter    *RateLimiter
	adaptiveConfig AdaptiveConfig
	loadMeter      *LoadMeter
}

// AdaptiveConfig contains configuration for adaptive rate limiting
type AdaptiveConfig struct {
	BaseMaxRequests    int           // Base maximum requests
	MaxMaxRequests     int           // Maximum allowed requests during low load
	MinMaxRequests     int           // Minimum allowed requests during high load
	AdaptationInterval time.Duration // How often to adjust limits
}

// LoadMeter measures server load for adaptive rate limiting
type LoadMeter struct {
	currentLoad   float64   // Current load percentage (0-100)
	lastUpdate    time.Time
	mutex         sync.RWMutex
	
	// Metrics
	requestCount  int64
	errorCount    int64
	responseTime  time.Duration
}

// NewAdaptiveRateLimiter creates a new adaptive rate limiter
func NewAdaptiveRateLimiter(config AdaptiveConfig, windowDuration time.Duration) *AdaptiveRateLimiter {
	return &AdaptiveRateLimiter{
		baseLimiter: NewRateLimiter(config.BaseMaxRequests, windowDuration),
		adaptiveConfig: config,
		loadMeter: &LoadMeter{
			currentLoad: 0.0,
			lastUpdate:  time.Now(),
		},
	}
}

// Allow checks if a request is allowed with adaptive limiting
func (arl *AdaptiveRateLimiter) Allow(identifier string) bool {
	// Update rate limits based on current load
	arl.updateLimits()
	
	return arl.baseLimiter.Allow(identifier)
}

// updateLimits adjusts rate limits based on current server load
func (arl *AdaptiveRateLimiter) updateLimits() {
	arl.loadMeter.mutex.RLock()
	load := arl.loadMeter.currentLoad
	lastUpdate := arl.loadMeter.lastUpdate
	arl.loadMeter.mutex.RUnlock()
	
	// Only update if enough time has passed
	if time.Since(lastUpdate) < arl.adaptiveConfig.AdaptationInterval {
		return
	}
	
	// Calculate new max requests based on load
	var newMaxRequests int
	if load < 30 { // Low load
		newMaxRequests = arl.adaptiveConfig.MaxMaxRequests
	} else if load > 80 { // High load
		newMaxRequests = arl.adaptiveConfig.MinMaxRequests
	} else { // Medium load - interpolate
		factor := (80 - load) / 50 // 0-1 scale
		newMaxRequests = arl.adaptiveConfig.MinMaxRequests + 
			int(factor * float64(arl.adaptiveConfig.MaxMaxRequests - arl.adaptiveConfig.MinMaxRequests))
	}
	
	// Update the base limiter's max requests
	arl.baseLimiter.maxRequests = newMaxRequests
	
	arl.loadMeter.mutex.Lock()
	arl.loadMeter.lastUpdate = time.Now()
	arl.loadMeter.mutex.Unlock()
}

// UpdateLoad updates the current server load measurement
func (arl *AdaptiveRateLimiter) UpdateLoad(load float64) {
	arl.loadMeter.mutex.Lock()
	defer arl.loadMeter.mutex.Unlock()
	
	arl.loadMeter.currentLoad = load
	arl.loadMeter.lastUpdate = time.Now()
}

// RecordRequest records request metrics for load calculation
func (arl *AdaptiveRateLimiter) RecordRequest(responseTime time.Duration, isError bool) {
	arl.loadMeter.mutex.Lock()
	defer arl.loadMeter.mutex.Unlock()
	
	arl.loadMeter.requestCount++
	arl.loadMeter.responseTime = responseTime
	
	if isError {
		arl.loadMeter.errorCount++
	}
	
	// Calculate load based on response time and error rate
	errorRate := float64(arl.loadMeter.errorCount) / float64(arl.loadMeter.requestCount)
	responseTimeMs := float64(responseTime.Milliseconds())
	
	// Simple load calculation (can be made more sophisticated)
	load := (responseTimeMs / 1000) * 50 + errorRate * 100
	if load > 100 {
		load = 100
	}
	
	arl.loadMeter.currentLoad = load
}

// GetLoadStats returns current load statistics
func (arl *AdaptiveRateLimiter) GetLoadStats() LoadStats {
	arl.loadMeter.mutex.RLock()
	defer arl.loadMeter.mutex.RUnlock()
	
	errorRate := float64(0)
	if arl.loadMeter.requestCount > 0 {
		errorRate = float64(arl.loadMeter.errorCount) / float64(arl.loadMeter.requestCount) * 100
	}
	
	return LoadStats{
		CurrentLoad:    arl.loadMeter.currentLoad,
		RequestCount:   arl.loadMeter.requestCount,
		ErrorCount:     arl.loadMeter.errorCount,
		ErrorRate:      errorRate,
		ResponseTime:   arl.loadMeter.responseTime,
		MaxRequests:    arl.baseLimiter.maxRequests,
	}
}

// LoadStats contains load measurement statistics
type LoadStats struct {
	CurrentLoad  float64       `json:"current_load"`
	RequestCount int64         `json:"request_count"`
	ErrorCount   int64         `json:"error_count"`
	ErrorRate    float64       `json:"error_rate"`
	ResponseTime time.Duration `json:"response_time"`
	MaxRequests  int           `json:"max_requests"`
}

// Stop stops the adaptive rate limiter
func (arl *AdaptiveRateLimiter) Stop() {
	arl.baseLimiter.Stop()
}

// SecurityMetrics tracks security-related metrics
type SecurityMetrics struct {
	mutex             sync.RWMutex
	blockedRequests   int64
	suspiciousIPs     map[string]int64
	threatAttempts    map[string]int64
	lastReset         time.Time
	resetInterval     time.Duration
}

// NewSecurityMetrics creates a new security metrics tracker
func NewSecurityMetrics(resetInterval time.Duration) *SecurityMetrics {
	sm := &SecurityMetrics{
		suspiciousIPs:  make(map[string]int64),
		threatAttempts: make(map[string]int64),
		lastReset:      time.Now(),
		resetInterval:  resetInterval,
	}
	
	// Start reset routine
	go sm.resetRoutine()
	
	return sm
}

// RecordBlockedRequest records a blocked request
func (sm *SecurityMetrics) RecordBlockedRequest(ip string, reason string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	sm.blockedRequests++
	sm.suspiciousIPs[ip]++
	sm.threatAttempts[reason]++
}

// GetMetrics returns current security metrics
func (sm *SecurityMetrics) GetMetrics() SecurityMetricsSnapshot {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	
	// Copy maps to avoid race conditions
	suspiciousIPs := make(map[string]int64)
	threatAttempts := make(map[string]int64)
	
	for ip, count := range sm.suspiciousIPs {
		suspiciousIPs[ip] = count
	}
	
	for threat, count := range sm.threatAttempts {
		threatAttempts[threat] = count
	}
	
	return SecurityMetricsSnapshot{
		BlockedRequests: sm.blockedRequests,
		SuspiciousIPs:   suspiciousIPs,
		ThreatAttempts:  threatAttempts,
		LastReset:       sm.lastReset,
	}
}

// SecurityMetricsSnapshot contains a snapshot of security metrics
type SecurityMetricsSnapshot struct {
	BlockedRequests int64            `json:"blocked_requests"`
	SuspiciousIPs   map[string]int64 `json:"suspicious_ips"`
	ThreatAttempts  map[string]int64 `json:"threat_attempts"`
	LastReset       time.Time        `json:"last_reset"`
}

// resetRoutine periodically resets metrics to prevent unbounded growth
func (sm *SecurityMetrics) resetRoutine() {
	ticker := time.NewTicker(sm.resetInterval)
	defer ticker.Stop()
	
	for range ticker.C {
		sm.reset()
	}
}

// reset clears metrics counters
func (sm *SecurityMetrics) reset() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	sm.blockedRequests = 0
	sm.suspiciousIPs = make(map[string]int64)
	sm.threatAttempts = make(map[string]int64)
	sm.lastReset = time.Now()
}