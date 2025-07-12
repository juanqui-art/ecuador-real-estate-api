package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"realty-core/internal/auth"
	"realty-core/internal/domain"
	"realty-core/internal/logging"
	"realty-core/internal/middleware"
	"realty-core/internal/service"
)

// AuthHandlers handles authentication endpoints
type AuthHandlers struct {
	userService *service.UserServiceSimple
	jwtManager  *auth.JWTManager
	logger      *logging.Logger
}

// NewAuthHandlers creates a new auth handlers instance
func NewAuthHandlers(userService *service.UserServiceSimple, jwtManager *auth.JWTManager) *AuthHandlers {
	return &AuthHandlers{
		userService: userService,
		jwtManager:  jwtManager,
		logger:      logging.GetGlobalLogger(),
	}
}

// LoginRequest represents login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse represents login response
type LoginResponse struct {
	User         *domain.User     `json:"user"`
	TokenPair    *auth.TokenPair  `json:"tokens"`
	ExpiresAt    string          `json:"expires_at"`
	Message      string          `json:"message"`
}

// RefreshRequest represents refresh token request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// LogoutRequest represents logout request
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token,omitempty"`
}

// LoginHandler handles user login
func (ah *AuthHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ah.handleError(w, "Invalid request format", http.StatusBadRequest, err)
		return
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		ah.handleError(w, "Email and password are required", http.StatusBadRequest, nil)
		return
	}

	// Clean email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Authenticate user
	user, err := ah.userService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		// Log security event
		if ah.logger != nil {
			ah.logger.SecurityEvent(
				"Login Failed",
				req.Email,
				"Invalid credentials",
				map[string]interface{}{
					"email": req.Email,
					"ip":    getClientIP(r),
				},
			)
		}
		ah.handleError(w, "Invalid email or password", http.StatusUnauthorized, err)
		return
	}

	// Generate JWT token pair
	agencyID := ""
	if user.AgencyID != nil {
		agencyID = *user.AgencyID
	}
	
	tokenPair, err := ah.jwtManager.GenerateTokenPair(
		user.ID,
		user.Email,
		string(user.Role),
		agencyID,
	)
	if err != nil {
		ah.handleError(w, "Failed to generate authentication tokens", http.StatusInternalServerError, err)
		return
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(15 * time.Minute).Format(time.RFC3339)

	// Log successful login
	if ah.logger != nil {
		ah.logger.Info("User logged in successfully", map[string]interface{}{
			"user_id": user.ID,
			"email":   user.Email,
			"role":    user.Role,
			"ip":      getClientIP(r),
		})
	}

	// Response
	response := LoginResponse{
		User:      user,
		TokenPair: tokenPair,
		ExpiresAt: expiresAt,
		Message:   "Login successful",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RefreshTokenHandler handles token refresh
func (ah *AuthHandlers) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ah.handleError(w, "Invalid request format", http.StatusBadRequest, err)
		return
	}

	// Basic validation for refresh token
	if req.RefreshToken == "" {
		ah.handleError(w, "Refresh token is required", http.StatusBadRequest, nil)
		return
	}

	// Validate refresh token
	refreshClaims, err := ah.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		ah.handleError(w, "Invalid or expired refresh token", http.StatusUnauthorized, err)
		return
	}

	// Get user info
	user, err := ah.userService.GetUser(refreshClaims.UserID)
	if err != nil {
		ah.handleError(w, "User not found", http.StatusNotFound, err)
		return
	}

	// Check if user is active
	if !user.Active {
		ah.handleError(w, "User account is disabled", http.StatusForbidden, nil)
		return
	}

	// Generate new token pair
	agencyID := ""
	if user.AgencyID != nil {
		agencyID = *user.AgencyID
	}
	
	tokenPair, err := ah.jwtManager.RefreshAccessToken(
		req.RefreshToken,
		user.Email,
		string(user.Role),
		agencyID,
	)
	if err != nil {
		ah.handleError(w, "Failed to refresh tokens", http.StatusInternalServerError, err)
		return
	}

	// Blacklist old refresh token
	ah.jwtManager.BlacklistRefreshToken(req.RefreshToken)

	// Calculate expiration time
	expiresAt := time.Now().Add(15 * time.Minute).Format(time.RFC3339)

	// Log token refresh
	if ah.logger != nil {
		ah.logger.Info("Token refreshed", map[string]interface{}{
			"user_id": user.ID,
			"email":   user.Email,
			"ip":      getClientIP(r),
		})
	}

	// Response
	response := LoginResponse{
		User:      user,
		TokenPair: tokenPair,
		ExpiresAt: expiresAt,
		Message:   "Token refreshed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// LogoutHandler handles user logout
func (ah *AuthHandlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get current access token
	authHeader := r.Header.Get("Authorization")
	accessToken := auth.ExtractTokenFromHeader(authHeader)

	// Get refresh token from request body (optional)
	var req LogoutRequest
	json.NewDecoder(r.Body).Decode(&req)

	// Get user ID from context
	userID := middleware.GetUserID(r.Context())
	email := middleware.GetUserEmail(r.Context())

	// Blacklist access token
	if accessToken != "" {
		ah.jwtManager.BlacklistToken(accessToken)
	}

	// Blacklist refresh token if provided
	if req.RefreshToken != "" {
		ah.jwtManager.BlacklistRefreshToken(req.RefreshToken)
	}

	// Log successful logout
	if ah.logger != nil {
		ah.logger.Info("User logged out", map[string]interface{}{
			"user_id": userID,
			"email":   email,
			"ip":      getClientIP(r),
		})
	}

	// Response
	response := map[string]interface{}{
		"message": "Logout successful",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ValidateTokenHandler validates current token (health check for auth)
func (ah *AuthHandlers) ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	email := middleware.GetUserEmail(r.Context())
	role := middleware.GetUserRole(r.Context())
	agencyID := middleware.GetAgencyID(r.Context())

	// Get user info
	user, err := ah.userService.GetUser(userID)
	if err != nil {
		ah.handleError(w, "User not found", http.StatusNotFound, err)
		return
	}

	// Check if user is still active
	if !user.Active {
		ah.handleError(w, "User account is disabled", http.StatusForbidden, nil)
		return
	}

	// Response with user info
	response := map[string]interface{}{
		"valid": true,
		"user": map[string]interface{}{
			"id":        userID,
			"email":     email,
			"role":      role,
			"agency_id": agencyID,
			"is_active": user.Active,
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ChangePasswordHandler handles password change
func (ah *AuthHandlers) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		ah.handleError(w, "Authentication required", http.StatusUnauthorized, nil)
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ah.handleError(w, "Invalid request format", http.StatusBadRequest, err)
		return
	}

	// Basic validation for change password
	if req.CurrentPassword == "" || req.NewPassword == "" {
		ah.handleError(w, "Current password and new password are required", http.StatusBadRequest, nil)
		return
	}

	// Change password
	err := ah.userService.ChangePassword(userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		if strings.Contains(err.Error(), "current password") {
			ah.handleError(w, "Current password is incorrect", http.StatusBadRequest, err)
			return
		}
		ah.handleError(w, "Failed to change password", http.StatusInternalServerError, err)
		return
	}

	// Log password change
	if ah.logger != nil {
		ah.logger.SecurityEvent(
			"Password Changed",
			middleware.GetUserEmail(r.Context()),
			"User changed password",
			map[string]interface{}{
				"user_id": userID,
				"ip":      getClientIP(r),
			},
		)
	}

	// Response
	response := map[string]interface{}{
		"message": "Password changed successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Helper functions

func (ah *AuthHandlers) handleError(w http.ResponseWriter, message string, statusCode int, err error) {
	if ah.logger != nil && err != nil {
		ah.logger.Error("Auth handler error", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"error":   "Authentication Error",
		"message": message,
		"code":    getErrorCode(statusCode),
	}

	json.NewEncoder(w).Encode(response)
}

// Helper function removed - using domain.User directly

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to remote address
	return r.RemoteAddr
}

func getErrorCode(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusInternalServerError:
		return "INTERNAL_ERROR"
	default:
		return "AUTH_ERROR"
	}
}