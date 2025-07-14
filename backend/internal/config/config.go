package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"realty-core/internal/logging"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Cache    CacheConfig
	Logging  LoggingConfig
	Security SecurityConfig
	Image    ImageConfig
	JWT      JWTConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	MaxHeaderBytes  int
	CORSOrigins     []string
	Environment     string // development, staging, production
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// CacheConfig holds caching configuration
type CacheConfig struct {
	Enabled         bool
	Capacity        int
	MaxSizeBytes    int64
	TTL             time.Duration
	CleanupInterval time.Duration
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level       logging.LogLevel
	Format      string // json, text
	ServiceName string
	Version     string
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	JWTSecret           string
	JWTExpiration       time.Duration
	BCryptCost          int
	RateLimitPerMinute  int
	MaxUploadSizeMB     int
	AllowedImageTypes   []string
}

// ImageConfig holds image processing configuration
type ImageConfig struct {
	StoragePath     string
	MaxWidth        int
	MaxHeight       int
	Quality         int
	ThumbnailSizes  []int
	AllowedFormats  []string
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	SecretKey        string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	Issuer           string
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			ReadTimeout:     getEnvDuration("READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getEnvDuration("WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:     getEnvDuration("IDLE_TIMEOUT", 120*time.Second),
			MaxHeaderBytes:  getEnvInt("MAX_HEADER_BYTES", 1<<20), // 1MB
			CORSOrigins:     getEnvList("CORS_ALLOWED_ORIGINS", []string{"*"}),
			Environment:     getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", "postgresql://juanquizhpi@localhost:5433/inmobiliaria_db?sslmode=disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
		},
		Cache: CacheConfig{
			Enabled:         getEnvBool("CACHE_ENABLED", true),
			Capacity:        getEnvInt("CACHE_CAPACITY", 1000),
			MaxSizeBytes:    int64(getEnvInt("CACHE_SIZE_MB", 100)) * 1024 * 1024,
			TTL:             getEnvDuration("CACHE_TTL", 24*time.Hour),
			CleanupInterval: getEnvDuration("CACHE_CLEANUP_INTERVAL", 10*time.Minute),
		},
		Logging: LoggingConfig{
			Level:       logging.ParseLogLevel(getEnv("LOG_LEVEL", "INFO")),
			Format:      getEnv("LOG_FORMAT", "json"),
			ServiceName: getEnv("SERVICE_NAME", "realty-core"),
			Version:     getEnv("SERVICE_VERSION", "1.9.0"),
		},
		Security: SecurityConfig{
			JWTSecret:           getEnv("JWT_SECRET", "default-secret-change-in-production"),
			JWTExpiration:       getEnvDuration("JWT_EXPIRATION", 24*time.Hour),
			BCryptCost:          getEnvInt("BCRYPT_COST", 12),
			RateLimitPerMinute:  getEnvInt("RATE_LIMIT_PER_MINUTE", 100),
			MaxUploadSizeMB:     getEnvInt("MAX_UPLOAD_SIZE_MB", 10),
			AllowedImageTypes:   getEnvList("ALLOWED_IMAGE_TYPES", []string{"image/jpeg", "image/png", "image/webp"}),
		},
		Image: ImageConfig{
			StoragePath:    getEnv("IMAGE_STORAGE_PATH", "uploads/images"),
			MaxWidth:       getEnvInt("IMAGE_MAX_WIDTH", 3000),
			MaxHeight:      getEnvInt("IMAGE_MAX_HEIGHT", 2000),
			Quality:        getEnvInt("IMAGE_QUALITY", 85),
			ThumbnailSizes: getEnvIntList("THUMBNAIL_SIZES", []int{150, 300, 600}),
			AllowedFormats: getEnvList("ALLOWED_IMAGE_FORMATS", []string{"jpeg", "jpg", "png", "webp"}),
		},
		JWT: JWTConfig{
			SecretKey:        getEnv("JWT_SECRET_KEY", "realty-core-jwt-secret-key-change-in-production-2025"),
			AccessTokenTTL:   getEnvDuration("JWT_ACCESS_TOKEN_TTL", 15*time.Minute),
			RefreshTokenTTL:  getEnvDuration("JWT_REFRESH_TOKEN_TTL", 7*24*time.Hour), // 7 days
			Issuer:           getEnv("JWT_ISSUER", "realty-core-api"),
		},
	}
}

// Helper functions for environment variable parsing

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvList(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

func getEnvIntList(key string, defaultValue []int) []int {
	if value := os.Getenv(key); value != "" {
		strValues := strings.Split(value, ",")
		intValues := make([]int, 0, len(strValues))
		for _, strValue := range strValues {
			if intValue, err := strconv.Atoi(strings.TrimSpace(strValue)); err == nil {
				intValues = append(intValues, intValue)
			}
		}
		if len(intValues) > 0 {
			return intValues
		}
	}
	return defaultValue
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return strings.ToLower(c.Server.Environment) == "production"
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.Server.Environment) == "development"
}

// IsStaging returns true if running in staging environment
func (c *Config) IsStaging() bool {
	return strings.ToLower(c.Server.Environment) == "staging"
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Add validation logic here
	// For example, check required fields, validate formats, etc.
	
	if c.Database.URL == "" {
		return &ConfigError{Field: "DATABASE_URL", Message: "Database URL is required"}
	}
	
	if c.Security.JWTSecret == "default-secret-change-in-production" && c.IsProduction() {
		return &ConfigError{Field: "JWT_SECRET", Message: "JWT secret must be changed in production"}
	}
	
	if c.Server.Port == "" {
		return &ConfigError{Field: "PORT", Message: "Server port is required"}
	}
	
	return nil
}

// ConfigError represents a configuration error
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return "Config error in " + e.Field + ": " + e.Message
}

// GetMaxUploadSizeBytes returns the maximum upload size in bytes
func (c *Config) GetMaxUploadSizeBytes() int64 {
	return int64(c.Security.MaxUploadSizeMB) * 1024 * 1024
}

// GetImageStorageURL returns the public URL path for images
func (c *Config) GetImageStorageURL() string {
	return "/uploads/images"
}

// GetDatabaseConnectionPoolConfig returns database connection pool configuration
func (c *Config) GetDatabaseConnectionPoolConfig() (maxOpen, maxIdle int, maxLifetime, maxIdleTime time.Duration) {
	return c.Database.MaxOpenConns, c.Database.MaxIdleConns, 
		   c.Database.ConnMaxLifetime, c.Database.ConnMaxIdleTime
}