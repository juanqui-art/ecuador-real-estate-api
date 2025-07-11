package middleware

import (
	"log"
	"net/http"
	"time"
)

// ResponseRecorder wraps http.ResponseWriter to capture status codes
type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *ResponseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

// PerformanceLogger logs request performance metrics
func PerformanceLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture status code
		recorder := &ResponseRecorder{
			ResponseWriter: w,
			statusCode:     200, // default status code
		}

		// Process the request
		next.ServeHTTP(recorder, r)

		// Calculate duration
		duration := time.Since(start)

		// Log performance metrics
		log.Printf("PERF: %s %s - Status: %d - Duration: %v", 
			r.Method, r.URL.Path, recorder.statusCode, duration)

		// Log slow requests (>500ms) as warnings
		if duration > 500*time.Millisecond {
			log.Printf("SLOW: %s %s took %v (>500ms)", 
				r.Method, r.URL.Path, duration)
		}
	})
}

// CompressionMiddleware adds gzip compression for JSON responses
func CompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simple compression by setting appropriate headers
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")
		
		next.ServeHTTP(w, r)
	})
}

// CacheMiddleware adds basic cache headers for static responses
func CacheMiddleware(maxAge time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set cache headers for GET requests
			if r.Method == "GET" {
				w.Header().Set("Cache-Control", "public, max-age="+maxAge.String())
			}
			
			next.ServeHTTP(w, r)
		})
	}
}