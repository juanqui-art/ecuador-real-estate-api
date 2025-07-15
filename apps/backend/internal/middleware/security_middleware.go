package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"realty-core/internal/logging"
	"realty-core/internal/security"
)

// SecurityMiddleware provides comprehensive security protection
type SecurityMiddleware struct {
	rateLimiter     *security.AdaptiveRateLimiter
	inputValidator  *security.InputValidator
	ipValidator     *security.IPValidator
	securityMetrics *security.SecurityMetrics
	logger          *logging.Logger
}

// NewSecurityMiddleware creates a new security middleware
func NewSecurityMiddleware() *SecurityMiddleware {
	adaptiveConfig := security.AdaptiveConfig{
		BaseMaxRequests:    100,
		MaxMaxRequests:     200,
		MinMaxRequests:     20,
		AdaptationInterval: 30 * time.Second,
	}

	return &SecurityMiddleware{
		rateLimiter:     security.NewAdaptiveRateLimiter(adaptiveConfig, time.Minute),
		inputValidator:  security.NewInputValidator(),
		ipValidator:     security.NewIPValidator(),
		securityMetrics: security.NewSecurityMetrics(time.Hour),
		logger:          logging.GetGlobalLogger(),
	}
}

// RateLimitMiddleware applies rate limiting with adaptive behavior
func (sm *SecurityMiddleware) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Get client identifier (IP address)
		clientIP := getClientIP(r)
		
		// Check rate limit
		if !sm.rateLimiter.Allow(clientIP) {
			sm.handleRateLimitExceeded(w, r, clientIP)
			return
		}
		
		// Create response recorder to capture status
		recorder := &securityResponseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		// Process request
		next.ServeHTTP(recorder, r)
		
		// Record request metrics
		duration := time.Since(start)
		isError := recorder.statusCode >= 400
		sm.rateLimiter.RecordRequest(duration, isError)
	})
}

// InputValidationMiddleware validates input for security threats
func (sm *SecurityMiddleware) InputValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip validation for GET requests (no body to validate)
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}
		
		// Validate URL parameters
		if !sm.validateURLParameters(w, r) {
			return
		}
		
		// Validate headers
		if !sm.validateHeaders(w, r) {
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// IPValidationMiddleware validates IP addresses and blocks suspicious ones
func (sm *SecurityMiddleware) IPValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r)
		
		// Validate IP address
		isValid, reason := sm.ipValidator.ValidateIP(clientIP)
		if !isValid {
			sm.handleBlockedIP(w, r, clientIP, reason)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// SecurityHeadersMiddleware adds security headers to responses
func (sm *SecurityMiddleware) SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		// Remove server information
		w.Header().Del("Server")
		w.Header().Del("X-Powered-By")
		
		next.ServeHTTP(w, r)
	})
}

// validateURLParameters validates URL parameters for security threats
func (sm *SecurityMiddleware) validateURLParameters(w http.ResponseWriter, r *http.Request) bool {
	for param, values := range r.URL.Query() {
		for _, value := range values {
			result := sm.inputValidator.ValidateInput(value, "url_param_"+param)
			if !result.IsValid {
				sm.handleSecurityThreat(w, r, "url_parameter_validation", result.Threats)
				return false
			}
		}
	}
	return true
}

// validateHeaders validates HTTP headers for security threats
func (sm *SecurityMiddleware) validateHeaders(w http.ResponseWriter, r *http.Request) bool {
	// Headers to validate
	headersToValidate := []string{
		"User-Agent",
		"Referer",
		"X-Forwarded-For",
		"X-Real-IP",
	}
	
	for _, headerName := range headersToValidate {
		headerValue := r.Header.Get(headerName)
		if headerValue != "" {
			result := sm.inputValidator.ValidateInput(headerValue, "header_"+strings.ToLower(headerName))
			if !result.IsValid {
				sm.handleSecurityThreat(w, r, "header_validation", result.Threats)
				return false
			}
		}
	}
	
	// Check for suspicious header patterns
	if sm.hasSuspiciousHeaders(r) {
		sm.handleSecurityThreat(w, r, "suspicious_headers", []security.ThreatInfo{
			{
				Type:        "suspicious_headers",
				Severity:    "medium",
				Description: "Suspicious header patterns detected",
				Pattern:     "header_analysis",
			},
		})
		return false
	}
	
	return true
}

// hasSuspiciousHeaders checks for suspicious header patterns
func (sm *SecurityMiddleware) hasSuspiciousHeaders(r *http.Request) bool {
	// Check for missing User-Agent (common in automated attacks)
	if r.Header.Get("User-Agent") == "" {
		return true
	}
	
	// Check for suspicious User-Agent patterns
	userAgent := strings.ToLower(r.Header.Get("User-Agent"))
	suspiciousPatterns := []string{
		"sqlmap", "nikto", "nmap", "masscan", "nessus",
		"burp", "owasp", "w3af", "dirbuster", "gobuster",
		"python-requests", "curl", "wget", "lwp-trivial",
	}
	
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(userAgent, pattern) {
			return true
		}
	}
	
	// Check for excessive headers (possible attack)
	if len(r.Header) > 50 {
		return true
	}
	
	// Check for unusually long header values
	for _, values := range r.Header {
		for _, value := range values {
			if len(value) > 8192 { // 8KB limit
				return true
			}
		}
	}
	
	return false
}

// handleRateLimitExceeded handles rate limit exceeded scenarios
func (sm *SecurityMiddleware) handleRateLimitExceeded(w http.ResponseWriter, r *http.Request, clientIP string) {
	sm.securityMetrics.RecordBlockedRequest(clientIP, "rate_limit_exceeded")
	
	if sm.logger != nil {
		sm.logger.SecurityEvent(
			"Rate Limit Exceeded",
			"",
			"Client exceeded rate limit",
			map[string]interface{}{
				"client_ip": clientIP,
				"method":    r.Method,
				"url":       r.URL.Path,
			},
		)
	}
	
	w.Header().Set("Retry-After", "60")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	
	response := map[string]interface{}{
		"error":   "Rate limit exceeded",
		"message": "Too many requests. Please try again later.",
		"code":    "RATE_LIMIT_EXCEEDED",
	}
	
	json.NewEncoder(w).Encode(response)
}

// handleBlockedIP handles blocked IP scenarios
func (sm *SecurityMiddleware) handleBlockedIP(w http.ResponseWriter, r *http.Request, clientIP, reason string) {
	sm.securityMetrics.RecordBlockedRequest(clientIP, "blocked_ip_"+reason)
	
	if sm.logger != nil {
		sm.logger.SecurityEvent(
			"Blocked IP Address",
			"",
			"Request from blocked IP address",
			map[string]interface{}{
				"client_ip": clientIP,
				"reason":    reason,
				"method":    r.Method,
				"url":       r.URL.Path,
			},
		)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	
	response := map[string]interface{}{
		"error":   "Access denied",
		"message": "Your IP address has been blocked",
		"code":    "IP_BLOCKED",
	}
	
	json.NewEncoder(w).Encode(response)
}

// handleSecurityThreat handles detected security threats
func (sm *SecurityMiddleware) handleSecurityThreat(w http.ResponseWriter, r *http.Request, threatType string, threats []security.ThreatInfo) {
	clientIP := getClientIP(r)
	sm.securityMetrics.RecordBlockedRequest(clientIP, threatType)
	
	if sm.logger != nil {
		sm.logger.SecurityEvent(
			"Security Threat Detected",
			"",
			"Malicious input pattern detected",
			map[string]interface{}{
				"client_ip":    clientIP,
				"threat_type":  threatType,
				"threats":      threats,
				"method":       r.Method,
				"url":          r.URL.Path,
				"user_agent":   r.Header.Get("User-Agent"),
			},
		)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	
	response := map[string]interface{}{
		"error":   "Security violation",
		"message": "Request contains potentially malicious content",
		"code":    "SECURITY_VIOLATION",
	}
	
	json.NewEncoder(w).Encode(response)
}

// GetSecurityMetrics returns current security metrics
func (sm *SecurityMiddleware) GetSecurityMetrics() security.SecurityMetricsSnapshot {
	return sm.securityMetrics.GetMetrics()
}

// GetLoadStats returns current load statistics
func (sm *SecurityMiddleware) GetLoadStats() security.LoadStats {
	return sm.rateLimiter.GetLoadStats()
}

// Stop stops the security middleware and cleanup routines
func (sm *SecurityMiddleware) Stop() {
	sm.rateLimiter.Stop()
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check for common proxy headers
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	
	if xcf := r.Header.Get("CF-Connecting-IP"); xcf != "" {
		return xcf
	}
	
	// Fallback to RemoteAddr
	addr := r.RemoteAddr
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[:idx]
	}
	
	return addr
}

// securityResponseRecorder captures response information for security middleware
type securityResponseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *securityResponseRecorder) WriteHeader(statusCode int) {
	rr.statusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}