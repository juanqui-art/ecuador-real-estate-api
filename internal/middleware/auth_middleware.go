package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"realty-core/internal/auth"
	"realty-core/internal/logging"
)

// AuthMiddleware provides JWT authentication middleware
type AuthMiddleware struct {
	jwtManager   *auth.JWTManager
	authManager  *auth.AuthorizationManager
	logger       *logging.Logger
	skipPaths    map[string]bool
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwtManager *auth.JWTManager, authManager *auth.AuthorizationManager) *AuthMiddleware {
	// Define paths that skip authentication
	skipPaths := map[string]bool{
		"/api/health":                    true,
		"/api/health/ready":              true,
		"/api/health/live":               true,
		"/api/health/detailed":           true,
		"/api/metrics":                   true,
		"/api/monitoring/metrics":        true,
		"/api/monitoring/prometheus":     true,
		"/api/auth/login":                true,
		"/api/auth/refresh":              true,
		"/api/properties":                true, // Public property listing
		"/api/properties/filter":         true, // Public property search
		"/api/properties/search/ranked":  true, // Public search
		"/api/properties/search/suggestions": true, // Public suggestions
		"/":                              true, // API documentation
	}

	return &AuthMiddleware{
		jwtManager:  jwtManager,
		authManager: authManager,
		logger:      logging.GetGlobalLogger(),
		skipPaths:   skipPaths,
	}
}

// contextKey is used for context values
type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	EmailKey    contextKey = "email"
	RoleKey     contextKey = "role"
	AgencyIDKey contextKey = "agency_id"
)

// Authenticate provides basic JWT authentication
func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip authentication for certain paths
		if am.skipPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		// Skip authentication for GET requests to public endpoints
		if r.Method == http.MethodGet && am.isPublicReadEndpoint(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Extract token from header
		authHeader := r.Header.Get("Authorization")
		token := auth.ExtractTokenFromHeader(authHeader)

		if token == "" {
			am.handleAuthError(w, "missing or invalid authorization header", http.StatusUnauthorized)
			return
		}

		// Validate token
		tokenInfo := am.jwtManager.ParseTokenInfo(token)
		if !tokenInfo.IsValid {
			am.handleAuthError(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add user info to context
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserIDKey, tokenInfo.UserID)
		ctx = context.WithValue(ctx, EmailKey, tokenInfo.Email)
		ctx = context.WithValue(ctx, RoleKey, tokenInfo.Role)
		ctx = context.WithValue(ctx, AgencyIDKey, tokenInfo.AgencyID)

		// Log authentication
		if am.logger != nil {
			am.logger.Info("User authenticated", map[string]interface{}{
				"user_id": tokenInfo.UserID,
				"role":    tokenInfo.Role,
				"method":  r.Method,
				"path":    r.URL.Path,
			})
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole creates middleware that requires a specific role
func (am *AuthMiddleware) RequireRole(requiredRole auth.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := GetUserRole(r.Context())
			if userRole == "" {
				am.handleAuthError(w, "authentication required", http.StatusUnauthorized)
				return
			}

			role, err := auth.ValidateRole(userRole)
			if err != nil {
				am.handleAuthError(w, "invalid user role", http.StatusForbidden)
				return
			}

			// Check if user has required role or higher
			if role != requiredRole && !auth.IsHigherRole(role, requiredRole) {
				am.handleAuthError(w, "insufficient privileges", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequirePermission creates middleware that requires a specific permission
func (am *AuthMiddleware) RequirePermission(permission auth.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := GetUserRole(r.Context())
			if userRole == "" {
				am.handleAuthError(w, "authentication required", http.StatusUnauthorized)
				return
			}

			role, err := auth.ValidateRole(userRole)
			if err != nil {
				am.handleAuthError(w, "invalid user role", http.StatusForbidden)
				return
			}

			// Check permission
			if !am.authManager.HasPermission(role, permission) {
				am.handleAuthError(w, "insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireResourceAccess creates middleware that checks resource-specific access
func (am *AuthMiddleware) RequireResourceAccess(permission auth.Permission, resourceIDExtractor func(*http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := GetUserRole(r.Context())
			userID := GetUserID(r.Context())
			agencyID := GetAgencyID(r.Context())

			if userRole == "" || userID == "" {
				am.handleAuthError(w, "authentication required", http.StatusUnauthorized)
				return
			}

			role, err := auth.ValidateRole(userRole)
			if err != nil {
				am.handleAuthError(w, "invalid user role", http.StatusForbidden)
				return
			}

			// Extract resource ID
			resourceID := resourceIDExtractor(r)

			// Check resource access
			context := &auth.ResourceContext{
				UserID:     userID,
				AgencyID:   agencyID,
				ResourceID: resourceID,
			}

			if !am.authManager.CanAccessResource(role, permission, context) {
				am.handleAuthError(w, "insufficient permissions for this resource", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// AdminOnly creates middleware that requires admin role
func (am *AuthMiddleware) AdminOnly() func(http.Handler) http.Handler {
	return am.RequireRole(auth.RoleAdmin)
}

// AgencyOrHigher creates middleware that requires agency role or higher
func (am *AuthMiddleware) AgencyOrHigher() func(http.Handler) http.Handler {
	return am.RequireRole(auth.RoleAgency)
}

// isPublicReadEndpoint checks if an endpoint allows public read access
func (am *AuthMiddleware) isPublicReadEndpoint(path string) bool {
	publicPaths := []string{
		"/api/properties/",
		"/api/images/",
		"/api/agencies",
		"/api/agencies/",
	}

	for _, publicPath := range publicPaths {
		if strings.HasPrefix(path, publicPath) {
			return true
		}
	}

	return false
}

// handleAuthError handles authentication/authorization errors
func (am *AuthMiddleware) handleAuthError(w http.ResponseWriter, message string, statusCode int) {
	if am.logger != nil {
		am.logger.SecurityEvent(
			"Authentication/Authorization Failed",
			"",
			message,
			map[string]interface{}{
				"status_code": statusCode,
			},
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"error":   "Authentication/Authorization Error",
		"message": message,
		"code":    getErrorCode(statusCode),
	}

	json.NewEncoder(w).Encode(response)
}

// getErrorCode returns error code based on status
func getErrorCode(statusCode int) string {
	switch statusCode {
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	default:
		return "AUTH_ERROR"
	}
}

// Helper functions to extract user info from context

// GetUserID extracts user ID from request context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// GetUserEmail extracts user email from request context
func GetUserEmail(ctx context.Context) string {
	if email, ok := ctx.Value(EmailKey).(string); ok {
		return email
	}
	return ""
}

// GetUserRole extracts user role from request context
func GetUserRole(ctx context.Context) string {
	if role, ok := ctx.Value(RoleKey).(string); ok {
		return role
	}
	return ""
}

// GetAgencyID extracts agency ID from request context
func GetAgencyID(ctx context.Context) string {
	if agencyID, ok := ctx.Value(AgencyIDKey).(string); ok {
		return agencyID
	}
	return ""
}

// ExtractResourceID helper functions for common patterns

// ExtractPropertyID extracts property ID from URL path
func ExtractPropertyID(r *http.Request) string {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	
	// Look for property ID in common patterns
	for i, part := range parts {
		if part == "properties" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	
	return ""
}

// ExtractUserID extracts user ID from URL path
func ExtractUserID(r *http.Request) string {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	
	for i, part := range parts {
		if part == "users" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	
	return ""
}

// ExtractAgencyID extracts agency ID from URL path
func ExtractAgencyID(r *http.Request) string {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	
	for i, part := range parts {
		if part == "agencies" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	
	return ""
}

// ExtractImageID extracts image ID from URL path
func ExtractImageID(r *http.Request) string {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	
	for i, part := range parts {
		if part == "images" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	
	return ""
}