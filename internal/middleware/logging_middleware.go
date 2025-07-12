package middleware

import (
	"net/http"
	"time"
	"strings"

	"realty-core/internal/logging"
)

// LoggingMiddleware provides structured HTTP request logging
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create a response writer that captures the status code
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		// Execute the next handler
		next.ServeHTTP(recorder, r)
		
		// Calculate duration
		duration := time.Since(start)
		
		// Extract user agent and remote address
		userAgent := r.Header.Get("User-Agent")
		remoteAddr := getRemoteAddr(r)
		
		// Log the request
		logger := logging.GetGlobalLogger()
		if logger != nil {
			fields := map[string]interface{}{
				"bytes_written": recorder.bytesWritten,
				"protocol":      r.Proto,
			}
			
			// Add query parameters if present
			if r.URL.RawQuery != "" {
				fields["query_params"] = r.URL.RawQuery
			}
			
			// Add content length if present
			if r.ContentLength > 0 {
				fields["content_length"] = r.ContentLength
			}
			
			// Add referer if present
			if referer := r.Header.Get("Referer"); referer != "" {
				fields["referer"] = referer
			}
			
			logger.HTTPRequest(
				r.Method,
				r.URL.Path,
				recorder.statusCode,
				duration,
				userAgent,
				remoteAddr,
				fields,
			)
		}
	})
}

// SecurityLoggingMiddleware logs security-related events
func SecurityLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetGlobalLogger()
		
		// Check for suspicious patterns
		if logger != nil {
			checkSuspiciousActivity(r, logger)
		}
		
		next.ServeHTTP(w, r)
	})
}

// ErrorLoggingMiddleware logs errors and panics
func ErrorLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger := logging.GetGlobalLogger()
				if logger != nil {
					fields := map[string]interface{}{
						"method":      r.Method,
						"url":         r.URL.Path,
						"remote_addr": getRemoteAddr(r),
						"user_agent":  r.Header.Get("User-Agent"),
						"panic":       err,
					}
					
					logger.Error("HTTP handler panic", nil, fields)
				}
				
				// Return 500 error
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}

// responseRecorder captures response information for logging
type responseRecorder struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

// WriteHeader captures the status code
func (rr *responseRecorder) WriteHeader(statusCode int) {
	rr.statusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}

// Write captures the number of bytes written
func (rr *responseRecorder) Write(data []byte) (int, error) {
	n, err := rr.ResponseWriter.Write(data)
	rr.bytesWritten += n
	return n, err
}

// getRemoteAddr extracts the real client IP address
func getRemoteAddr(r *http.Request) string {
	// Check for common proxy headers
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, get the first one
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

// checkSuspiciousActivity checks for common attack patterns
func checkSuspiciousActivity(r *http.Request, logger *logging.Logger) {
	userAgent := r.Header.Get("User-Agent")
	url := r.URL.Path
	query := r.URL.RawQuery
	
	// Check for SQL injection patterns
	suspiciousPatterns := []string{
		"union", "select", "insert", "delete", "update", "drop",
		"script", "javascript:", "onerror", "onload",
		"../", "..\\", "/etc/passwd", "/proc/",
		"<script", "</script>", "eval(", "alert(",
	}
	
	content := strings.ToLower(url + " " + query + " " + userAgent)
	
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(content, pattern) {
			logger.SecurityEvent(
				"Suspicious Request Pattern",
				"", // No user ID available at this level
				"Pattern detected: "+pattern,
				map[string]interface{}{
					"method":      r.Method,
					"url":         url,
					"query":       query,
					"user_agent":  userAgent,
					"remote_addr": getRemoteAddr(r),
					"pattern":     pattern,
				},
			)
			break
		}
	}
	
	// Check for unusual user agents
	if userAgent == "" {
		logger.SecurityEvent(
			"Missing User Agent",
			"",
			"Request without User-Agent header",
			map[string]interface{}{
				"method":      r.Method,
				"url":         url,
				"remote_addr": getRemoteAddr(r),
			},
		)
	}
	
	// Check for excessive query parameters (possible DoS attempt)
	if len(r.URL.Query()) > 50 {
		logger.SecurityEvent(
			"Excessive Query Parameters",
			"",
			"Request with unusually high number of query parameters",
			map[string]interface{}{
				"method":           r.Method,
				"url":              url,
				"param_count":      len(r.URL.Query()),
				"remote_addr":      getRemoteAddr(r),
			},
		)
	}
	
	// Check for unusually long URLs (possible buffer overflow attempt)
	if len(r.URL.RequestURI()) > 2048 {
		logger.SecurityEvent(
			"Unusually Long URL",
			"",
			"Request with unusually long URL",
			map[string]interface{}{
				"method":      r.Method,
				"url_length":  len(r.URL.RequestURI()),
				"remote_addr": getRemoteAddr(r),
			},
		)
	}
}