package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// Claims represents JWT claims with user information
type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	AgencyID string `json:"agency_id,omitempty"`
	jwt.RegisteredClaims
}

// RefreshClaims represents refresh token claims
type RefreshClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token operations
type JWTManager struct {
	secretKey        []byte
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
	issuer           string
	blacklistedTokens map[string]bool // In production, use Redis
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey string, accessTTL, refreshTTL time.Duration, issuer string) *JWTManager {
	return &JWTManager{
		secretKey:        []byte(secretKey),
		accessTokenTTL:   accessTTL,
		refreshTokenTTL:  refreshTTL,
		issuer:           issuer,
		blacklistedTokens: make(map[string]bool),
	}
}

// GenerateTokenPair creates access and refresh tokens for a user
func (j *JWTManager) GenerateTokenPair(userID, email, role, agencyID string) (*TokenPair, error) {
	now := time.Now()
	
	// Create access token claims
	accessClaims := &Claims{
		UserID:   userID,
		Email:    email,
		Role:     role,
		AgencyID: agencyID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.issuer,
			Subject:   userID,
		},
	}
	
	// Create refresh token claims
	refreshClaims := &RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.issuer,
			Subject:   userID,
		},
	}
	
	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(j.secretKey)
	if err != nil {
		return nil, err
	}
	
	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(j.secretKey)
	if err != nil {
		return nil, err
	}
	
	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(j.accessTokenTTL.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// ValidateAccessToken validates and parses an access token
func (j *JWTManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	// Check if token is blacklisted
	if j.blacklistedTokens[tokenString] {
		return nil, errors.New("token is blacklisted")
	}
	
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secretKey, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Check if token is expired
		if claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, errors.New("token is expired")
		}
		return claims, nil
	}
	
	return nil, errors.New("invalid token")
}

// ValidateRefreshToken validates and parses a refresh token
func (j *JWTManager) ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
	// Check if token is blacklisted
	if j.blacklistedTokens[tokenString] {
		return nil, errors.New("refresh token is blacklisted")
	}
	
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secretKey, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		// Check if token is expired
		if claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, errors.New("refresh token is expired")
		}
		return claims, nil
	}
	
	return nil, errors.New("invalid refresh token")
}

// RefreshAccessToken creates a new access token using a valid refresh token
func (j *JWTManager) RefreshAccessToken(refreshTokenString string, email, role, agencyID string) (*TokenPair, error) {
	// Validate refresh token
	refreshClaims, err := j.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, err
	}
	
	// Generate new token pair
	return j.GenerateTokenPair(refreshClaims.UserID, email, role, agencyID)
}

// BlacklistToken adds a token to the blacklist
func (j *JWTManager) BlacklistToken(tokenString string) {
	j.blacklistedTokens[tokenString] = true
}

// BlacklistRefreshToken adds a refresh token to the blacklist
func (j *JWTManager) BlacklistRefreshToken(tokenString string) {
	j.blacklistedTokens[tokenString] = true
}

// IsTokenBlacklisted checks if a token is blacklisted
func (j *JWTManager) IsTokenBlacklisted(tokenString string) bool {
	return j.blacklistedTokens[tokenString]
}

// ExtractTokenFromHeader extracts token from Authorization header
func ExtractTokenFromHeader(authHeader string) string {
	if authHeader == "" {
		return ""
	}
	
	// Bearer token format: "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	
	return parts[1]
}

// TokenInfo contains parsed token information for middleware
type TokenInfo struct {
	UserID   string
	Email    string
	Role     string
	AgencyID string
	IsValid  bool
}

// ParseTokenInfo parses token and returns user information
func (j *JWTManager) ParseTokenInfo(tokenString string) *TokenInfo {
	claims, err := j.ValidateAccessToken(tokenString)
	if err != nil {
		return &TokenInfo{IsValid: false}
	}
	
	return &TokenInfo{
		UserID:   claims.UserID,
		Email:    claims.Email,
		Role:     claims.Role,
		AgencyID: claims.AgencyID,
		IsValid:  true,
	}
}

// CleanupExpiredTokens removes expired tokens from blacklist
// In production, this would be handled by Redis TTL
func (j *JWTManager) CleanupExpiredTokens() {
	now := time.Now()
	
	for tokenString := range j.blacklistedTokens {
		// Parse token to check expiration
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return j.secretKey, nil
		})
		
		if err != nil || !token.Valid {
			delete(j.blacklistedTokens, tokenString)
			continue
		}
		
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if exp, ok := claims["exp"].(float64); ok {
				expirationTime := time.Unix(int64(exp), 0)
				if expirationTime.Before(now) {
					delete(j.blacklistedTokens, tokenString)
				}
			}
		}
	}
}

// GetBlacklistedTokensCount returns the number of blacklisted tokens
func (j *JWTManager) GetBlacklistedTokensCount() int {
	return len(j.blacklistedTokens)
}

