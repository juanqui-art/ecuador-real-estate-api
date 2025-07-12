package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"realty-core/internal/auth"
	"realty-core/internal/cache"
	"realty-core/internal/config"
	"realty-core/internal/handlers"
	"realty-core/internal/logging"
	"realty-core/internal/middleware"
	"realty-core/internal/monitoring"
	"realty-core/internal/processors"
	"realty-core/internal/repository"
	"realty-core/internal/service"
	"realty-core/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	// Parse command line flags
	healthCheck := flag.Bool("health-check", false, "Run health check and exit")
	flag.Parse()

	// If health check flag is provided, run health check and exit
	if *healthCheck {
		err := performHealthCheck()
		if err != nil {
			log.Printf("Health check failed: %v", err)
			os.Exit(1)
		}
		log.Println("Health check passed")
		os.Exit(0)
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system variables")
	}

	// Load configuration
	cfg := config.LoadConfig()
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Initialize structured logging
	logger := logging.NewLogger(logging.Config{
		Level:       cfg.Logging.Level,
		ServiceName: cfg.Logging.ServiceName,
		Version:     cfg.Logging.Version,
		Format:      cfg.Logging.Format,
	})
	
	// Set global logger
	logging.SetGlobalLogger(logger)
	
	// Initialize monitoring
	monitoring.InitializeMetrics()
	monitoring.InitializeAlerts()
	
	logger.Info("Starting real estate properties server", map[string]interface{}{
		"version":     cfg.Logging.Version,
		"environment": cfg.Server.Environment,
		"port":        cfg.Server.Port,
		"log_level":   cfg.Logging.Level.String(),
	})

	// Connect to database
	logger.Info("Connecting to PostgreSQL database")
	db, err := repository.ConnectDatabase(cfg.Database.URL)
	if err != nil {
		logger.Fatal("Failed to connect to database", err)
	}
	defer db.Close()
	
	// Configure database connection pool
	maxOpen, maxIdle, maxLifetime, maxIdleTime := cfg.GetDatabaseConnectionPoolConfig()
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(maxLifetime)
	db.SetConnMaxIdleTime(maxIdleTime)
	
	logger.Info("Database connection established", map[string]interface{}{
		"max_open_conns":     maxOpen,
		"max_idle_conns":     maxIdle,
		"conn_max_lifetime":  maxLifetime.String(),
		"conn_max_idle_time": maxIdleTime.String(),
	})

	// Create repositories
	propertyRepo := repository.NewPostgreSQLPropertyRepository(db)
	imageRepo := repository.NewPostgreSQLImageRepository(db)
	userRepo := repository.NewUserRepository(db)
	agencyRepo := repository.NewAgencyRepository(db)

	// Create image system dependencies
	imageStorage, err := storage.NewLocalImageStorage(
		cfg.Image.StoragePath, 
		cfg.GetImageStorageURL(), 
		cfg.GetMaxUploadSizeBytes(),
	)
	if err != nil {
		logger.Fatal("Failed to create image storage", err)
	}
	
	imageProcessor := processors.NewImageProcessor(cfg.Image.MaxWidth, cfg.Image.MaxHeight)
	
	imageCacheConfig := cache.ImageCacheConfig{
		Enabled:         cfg.Cache.Enabled,
		Capacity:        cfg.Cache.Capacity,
		MaxSizeBytes:    cfg.Cache.MaxSizeBytes,
		TTL:             cfg.Cache.TTL,
		CleanupInterval: cfg.Cache.CleanupInterval,
	}
	imageCache := cache.NewImageCache(imageCacheConfig)
	
	logger.Info("Image system initialized", map[string]interface{}{
		"storage_path":        cfg.Image.StoragePath,
		"max_upload_size_mb":  cfg.Security.MaxUploadSizeMB,
		"cache_enabled":       cfg.Cache.Enabled,
		"cache_capacity":      cfg.Cache.Capacity,
		"cache_size_mb":       cfg.Cache.MaxSizeBytes / (1024 * 1024),
	})

	// Create services with structured logging
	propertyService := service.NewPropertyService(propertyRepo, imageRepo)
	imageService := service.NewImageService(imageRepo, propertyRepo, imageStorage, imageProcessor, imageCache)
	userService := service.NewUserService(userRepo, agencyRepo, log.Default()) // TODO: Update to use structured logger
	agencyService := service.NewAgencyService(agencyRepo, userRepo, log.Default()) // TODO: Update to use structured logger
	permissionService := service.NewPermissionService()
	
	logger.Info("Services initialized successfully")

	// Create JWT manager
	jwtManager := auth.NewJWTManager(
		cfg.JWT.SecretKey,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
		cfg.JWT.Issuer,
	)
	
	// Create authorization manager
	authManager := auth.NewAuthorizationManager()
	
	// Create authentication middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtManager, authManager)
	
	// Create security middleware
	securityMiddleware := middleware.NewSecurityMiddleware()
	
	// Create monitoring handlers
	metricsCollector := monitoring.GetGlobalMetrics()
	alertManager := monitoring.GetGlobalAlertManager()
	monitoringHandler := handlers.NewMonitoringHandler(metricsCollector, alertManager)
	
	// Create handlers
	propertyHandler := handlers.NewPropertyHandler(propertyService)
	imageHandler := handlers.NewImageHandler(imageService)
	userHandler := handlers.NewUserHandlerSimple(userService, permissionService, log.Default()) // TODO: Update to use structured logger
	agencyHandler := handlers.NewAgencyHandlerSimple(agencyService, permissionService, log.Default()) // TODO: Update to use structured logger
	paginationHandler := handlers.NewPaginationHandlerSimple(propertyService, imageService, userService, agencyService, log.Default()) // TODO: Update to use structured logger
	healthHandler := handlers.NewHealthHandler(db, propertyRepo, imageRepo, userRepo, agencyRepo, imageCache, propertyService)
	securityHandler := handlers.NewSecurityHandler(securityMiddleware)
	authHandler := handlers.NewAuthHandlers(userService, jwtManager)
	
	logger.Info("Handlers initialized successfully")

	// Configure routes
	router := configureRoutesWithAuthentication(propertyHandler, imageHandler, userHandler, agencyHandler, paginationHandler, healthHandler, securityHandler, monitoringHandler, authHandler, authMiddleware)
	
	// Apply middleware stack (order matters - monitoring and security first)
	finalHandler := middleware.ErrorLoggingMiddleware(
		middleware.AlertEvaluationMiddleware(
			middleware.MonitoringMiddleware(
				middleware.PerformanceMonitoringMiddleware(
					securityMiddleware.SecurityHeadersMiddleware(
						securityMiddleware.RateLimitMiddleware(
							securityMiddleware.IPValidationMiddleware(
								securityMiddleware.InputValidationMiddleware(
									middleware.SecurityLoggingMiddleware(
										middleware.LoggingMiddleware(
											middleware.PerformanceLogger(router),
										),
									),
								),
							),
						),
					),
				),
			),
		),
	)
	
	logger.Info("Middleware stack configured", map[string]interface{}{
		"middlewares": []string{"ErrorLogging", "AlertEvaluation", "Monitoring", "PerformanceMonitoring", "SecurityHeaders", "RateLimit", "IPValidation", "InputValidation", "SecurityLogging", "RequestLogging", "Performance"},
	})

	// Configure HTTP server with settings from config
	server := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        addCORSMiddleware(finalHandler, cfg.Server.CORSOrigins),
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}
	
	logger.Info("HTTP server configured", map[string]interface{}{
		"port":             cfg.Server.Port,
		"read_timeout":     cfg.Server.ReadTimeout.String(),
		"write_timeout":    cfg.Server.WriteTimeout.String(),
		"idle_timeout":     cfg.Server.IdleTimeout.String(),
		"max_header_bytes": cfg.Server.MaxHeaderBytes,
	})

	logger.Info("Server startup complete", map[string]interface{}{
		"endpoints": map[string]string{
			"documentation": fmt.Sprintf("http://localhost:%s/", cfg.Server.Port),
			"health_basic":  fmt.Sprintf("http://localhost:%s/api/health", cfg.Server.Port),
			"health_detailed": fmt.Sprintf("http://localhost:%s/api/health/detailed", cfg.Server.Port),
			"metrics":       fmt.Sprintf("http://localhost:%s/api/metrics", cfg.Server.Port),
		},
	})

	// Start server
	logger.Info("Starting HTTP server", map[string]interface{}{
		"address": server.Addr,
	})
	
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Failed to start HTTP server", err)
	}
}

// getVariable gets an environment variable or returns a default value
func getVariable(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// addCORSMiddleware adds CORS headers based on configuration
func addCORSMiddleware(next http.Handler, allowedOrigins []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configure CORS headers
		origin := r.Header.Get("Origin")
		if origin != "" && isAllowedOrigin(origin, allowedOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else if len(allowedOrigins) == 1 && allowedOrigins[0] == "*" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Execute next handler
		next.ServeHTTP(w, r)
	})
}

// isAllowedOrigin checks if the origin is in the allowed list
func isAllowedOrigin(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}



// configureRoutesWithAuthentication configures routes with authentication, authorization, and all existing systems
func configureRoutesWithAuthentication(propertyHandler *handlers.PropertyHandler, imageHandler *handlers.ImageHandler, userHandler *handlers.UserHandlerSimple, agencyHandler *handlers.AgencyHandlerSimple, paginationHandler *handlers.PaginationHandlerSimple, healthHandler *handlers.HealthHandler, securityHandler *handlers.SecurityHandler, monitoringHandler *handlers.MonitoringHandler, authHandler *handlers.AuthHandlers, authMiddleware *middleware.AuthMiddleware) *http.ServeMux {
	mux := http.NewServeMux()

	// ============== AUTHENTICATION ROUTES ==============
	// Authentication routes (public - no middleware)
	mux.HandleFunc("/api/auth/login", authHandler.LoginHandler)
	mux.HandleFunc("/api/auth/refresh", authHandler.RefreshTokenHandler)
	
	// Protected auth routes
	mux.Handle("/api/auth/logout", authMiddleware.Authenticate(http.HandlerFunc(authHandler.LogoutHandler)))
	mux.Handle("/api/auth/validate", authMiddleware.Authenticate(http.HandlerFunc(authHandler.ValidateTokenHandler)))
	mux.Handle("/api/auth/change-password", authMiddleware.Authenticate(http.HandlerFunc(authHandler.ChangePasswordHandler)))

	// ============== PROPERTY ROUTES ==============
	
	// Basic CRUD (GET is public, POST requires authentication)
	mux.HandleFunc("/api/properties", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Public endpoint - no authentication required
			propertyHandler.ListProperties(w, r)
		case http.MethodPost:
			// Requires property creation permission
			authMiddleware.RequirePermission(auth.PermissionPropertyCreate)(
				http.HandlerFunc(propertyHandler.CreateProperty),
			).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Property by slug (must be before /api/properties/{id})
	mux.HandleFunc("/api/properties/slug/", propertyHandler.GetPropertyBySlug)
	
	// Basic filtering and search
	mux.HandleFunc("/api/properties/filter", propertyHandler.FilterProperties)
	mux.HandleFunc("/api/properties/statistics", propertyHandler.GetStatistics)
	
	// Advanced search endpoints
	mux.HandleFunc("/api/properties/search/ranked", propertyHandler.SearchRanked)
	mux.HandleFunc("/api/properties/search/suggestions", propertyHandler.SearchSuggestions)
	
	// Advanced search with POST for complex filters
	mux.HandleFunc("/api/properties/search/advanced", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			propertyHandler.AdvancedSearch(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Property images management (specific routes first)
	mux.HandleFunc("/api/properties/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.Contains(path, "/images") {
			if strings.HasSuffix(path, "/images") && r.Method == http.MethodGet {
				// Public endpoint - view property images
				imageHandler.GetImagesByProperty(w, r)
			} else if strings.Contains(path, "/images/reorder") && r.Method == http.MethodPost {
				// Requires image update permission with resource access
				authMiddleware.RequireResourceAccess(auth.PermissionImageUpdate, middleware.ExtractPropertyID)(
					http.HandlerFunc(imageHandler.ReorderImages),
				).ServeHTTP(w, r)
			} else if strings.Contains(path, "/images/main") {
				if r.Method == http.MethodGet {
					// Public endpoint - view main image
					imageHandler.GetMainImage(w, r)
				} else if r.Method == http.MethodPost {
					// Requires image update permission with resource access
					authMiddleware.RequireResourceAccess(auth.PermissionImageUpdate, middleware.ExtractPropertyID)(
						http.HandlerFunc(imageHandler.SetMainImage),
					).ServeHTTP(w, r)
				}
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		} else {
			// Regular property operations
			switch r.Method {
			case http.MethodGet:
				// Public endpoint - view property details
				propertyHandler.GetProperty(w, r)
			case http.MethodPut:
				// Requires property update permission with resource access
				authMiddleware.RequireResourceAccess(auth.PermissionPropertyUpdate, middleware.ExtractPropertyID)(
					http.HandlerFunc(propertyHandler.UpdateProperty),
				).ServeHTTP(w, r)
			case http.MethodDelete:
				// Requires property delete permission with resource access
				authMiddleware.RequireResourceAccess(auth.PermissionPropertyDelete, middleware.ExtractPropertyID)(
					http.HandlerFunc(propertyHandler.DeleteProperty),
				).ServeHTTP(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})

	// ============== IMAGE ROUTES ==============

	// Image upload and basic operations
	mux.HandleFunc("/api/images", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// Requires image upload permission
			authMiddleware.RequirePermission(auth.PermissionImageUpload)(
				http.HandlerFunc(imageHandler.UploadImage),
			).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Image statistics and maintenance (specific routes first)
	mux.HandleFunc("/api/images/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			imageHandler.GetImageStats(w, r)
		}
	})

	mux.HandleFunc("/api/images/cleanup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			imageHandler.CleanupTempFiles(w, r)
		}
	})

	mux.HandleFunc("/api/images/cache/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			imageHandler.GetCacheStats(w, r)
		}
	})

	// Individual image routes (must be after specific routes)
	mux.HandleFunc("/api/images/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.Contains(path, "/variant") && r.Method == http.MethodGet {
			imageHandler.GetImageVariant(w, r)
		} else if strings.Contains(path, "/thumbnail") && r.Method == http.MethodGet {
			imageHandler.GetThumbnail(w, r)
		} else {
			// Regular image operations
			switch r.Method {
			case http.MethodGet:
				imageHandler.GetImage(w, r)
			case http.MethodPut:
				imageHandler.UpdateImageMetadata(w, r)
			case http.MethodDelete:
				imageHandler.DeleteImage(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})

	// ============== USER ROUTES ==============

	// User management routes (protected)
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Requires user list permission
			authMiddleware.RequirePermission(auth.PermissionUserList)(
				http.HandlerFunc(userHandler.SearchUsers),
			).ServeHTTP(w, r)
		case http.MethodPost:
			// Requires user creation permission - typically admin or agency
			authMiddleware.RequirePermission(auth.PermissionUserCreate)(
				http.HandlerFunc(userHandler.CreateUser),
			).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// User statistics and dashboard (specific routes first)
	mux.Handle("/api/users/statistics", authMiddleware.RequirePermission(auth.PermissionSystemAnalytics)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				userHandler.GetUserStatistics(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}),
	))

	mux.Handle("/api/users/dashboard", authMiddleware.Authenticate(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				userHandler.GetUserDashboard(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}),
	))

	// Users by role
	mux.Handle("/api/users/role/", authMiddleware.RequirePermission(auth.PermissionUserList)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				userHandler.GetUsersByRole(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}),
	))

	// Individual user routes (must be after specific routes)
	mux.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.Contains(path, "/password") && r.Method == http.MethodPost {
			// Password change via new auth handler
			authMiddleware.Authenticate(http.HandlerFunc(authHandler.ChangePasswordHandler)).ServeHTTP(w, r)
		} else {
			// Regular user operations
			switch r.Method {
			case http.MethodGet:
				// Requires user read permission with resource access
				authMiddleware.RequireResourceAccess(auth.PermissionUserRead, middleware.ExtractUserID)(
					http.HandlerFunc(userHandler.GetUser),
				).ServeHTTP(w, r)
			case http.MethodPut:
				// Requires user update permission with resource access
				authMiddleware.RequireResourceAccess(auth.PermissionUserUpdate, middleware.ExtractUserID)(
					http.HandlerFunc(userHandler.UpdateUser),
				).ServeHTTP(w, r)
			case http.MethodDelete:
				// Requires user delete permission with resource access
				authMiddleware.RequireResourceAccess(auth.PermissionUserDelete, middleware.ExtractUserID)(
					http.HandlerFunc(userHandler.DeleteUser),
				).ServeHTTP(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})

	// ============== AGENCY ROUTES ==============

	// Agency management routes
	mux.HandleFunc("/api/agencies", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Public endpoint - search agencies (for buyers/owners)
			agencyHandler.SearchAgencies(w, r)
		case http.MethodPost:
			// Requires agency creation permission - typically admin only
			authMiddleware.RequirePermission(auth.PermissionAgencyCreate)(
				http.HandlerFunc(agencyHandler.CreateAgency),
			).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Agency statistics (specific routes first)
	mux.Handle("/api/agencies/statistics", authMiddleware.RequirePermission(auth.PermissionSystemAnalytics)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				agencyHandler.GetAgencyStatistics(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}),
	))

	// Public endpoint - active agencies listing
	mux.HandleFunc("/api/agencies/active", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			agencyHandler.GetActiveAgencies(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Agencies by service area
	mux.HandleFunc("/api/agencies/service-area/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			agencyHandler.GetAgenciesByServiceArea(w, r)
		}
	})

	// Agencies by specialty
	mux.HandleFunc("/api/agencies/specialty/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			agencyHandler.GetAgenciesBySpecialty(w, r)
		}
	})

	// Individual agency routes (must be after specific routes)
	mux.HandleFunc("/api/agencies/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.Contains(path, "/agents") && r.Method == http.MethodGet {
			agencyHandler.GetAgencyAgents(w, r)
		} else if strings.Contains(path, "/license") && r.Method == http.MethodPost {
			agencyHandler.SetAgencyLicense(w, r)
		} else if strings.Contains(path, "/performance") && r.Method == http.MethodGet {
			agencyHandler.GetAgencyPerformance(w, r)
		} else {
			// Regular agency operations
			switch r.Method {
			case http.MethodGet:
				agencyHandler.GetAgency(w, r)
			case http.MethodPut:
				agencyHandler.UpdateAgency(w, r)
			case http.MethodDelete:
				agencyHandler.DeleteAgency(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})

	// ============== PAGINATION ROUTES ==============

	// Paginated entity endpoints
	mux.HandleFunc("/api/pagination/properties", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			paginationHandler.GetPaginatedProperties(w, r)
		}
	})

	mux.HandleFunc("/api/pagination/images", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			paginationHandler.GetPaginatedImages(w, r)
		}
	})

	mux.HandleFunc("/api/pagination/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			paginationHandler.GetPaginatedUsers(w, r)
		}
	})

	mux.HandleFunc("/api/pagination/agencies", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			paginationHandler.GetPaginatedAgencies(w, r)
		}
	})

	mux.HandleFunc("/api/pagination/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			paginationHandler.GetPaginatedSearch(w, r)
		}
	})

	mux.HandleFunc("/api/pagination/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			paginationHandler.GetPaginationStats(w, r)
		}
	})

	mux.HandleFunc("/api/pagination/advanced", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			paginationHandler.HandleAdvancedPagination(w, r)
		}
	})

	// ============== HEALTH CHECK AND MONITORING ROUTES ==============
	
	// Basic health check (for simple load balancers)
	mux.HandleFunc("/api/health", healthHandler.BasicHealthCheck)
	
	// Detailed health check (for monitoring systems)
	mux.HandleFunc("/api/health/detailed", healthHandler.DetailedHealthCheck)
	
	// Kubernetes readiness probe
	mux.HandleFunc("/api/health/ready", healthHandler.ReadinessCheck)
	
	// Kubernetes liveness probe
	mux.HandleFunc("/api/health/live", healthHandler.LivenessCheck)
	
	// Cache health and statistics
	mux.HandleFunc("/api/health/cache", healthHandler.CacheHealth)
	
	// Prometheus metrics endpoint
	mux.HandleFunc("/api/metrics", healthHandler.MetricsEndpoint)
	
	// ============== SECURITY ROUTES ==============
	
	// Security monitoring and metrics
	mux.HandleFunc("/api/security/metrics", securityHandler.SecurityMetrics)
	mux.HandleFunc("/api/security/status", securityHandler.SecurityStatus)
	mux.HandleFunc("/api/security/threats", securityHandler.ThreatIntelligence)
	
	// Security validation endpoints
	mux.HandleFunc("/api/security/validate/password", securityHandler.ValidatePassword)
	mux.HandleFunc("/api/security/validate/input", securityHandler.ValidateInput)
	
	// ============== MONITORING ROUTES ==============
	
	// Monitoring and metrics endpoints
	mux.HandleFunc("/api/monitoring/metrics", monitoringHandler.GetMetrics)
	mux.HandleFunc("/api/monitoring/prometheus", monitoringHandler.GetPrometheusMetrics)
	mux.HandleFunc("/api/monitoring/dashboard", monitoringHandler.GetDashboard)
	
	// Alert management endpoints
	mux.HandleFunc("/api/monitoring/alerts", monitoringHandler.GetAlerts)
	mux.HandleFunc("/api/monitoring/alerts/history", monitoringHandler.GetAlertHistory)
	mux.HandleFunc("/api/monitoring/rules", monitoringHandler.GetAlertRules)
	mux.HandleFunc("/api/monitoring/rules/", monitoringHandler.UpdateAlertRule)

	// Root route with API documentation
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"message": "Real Estate Properties API - Complete System with Pagination",
			"version": "1.9.0",
			"status": "Full-featured system with advanced pagination operational",
			"features": [
				"Property Management (CRUD)",
				"Image Management (Upload, Processing, Cache)",
				"User Management & Authentication",
				"Agency Management & Performance",
				"Advanced Pagination System",
				"Cross-entity Search",
				"Role-based Access Control (Admin, Agent, Buyer, Owner, Agency)",
				"Search & Filtering", 
				"PostgreSQL Database",
				"LRU Image Cache",
				"Image Processing (Resize, Thumbnails)",
				"Password Hashing (bcrypt)",
				"RUC Validation (Ecuador)",
				"SEO-friendly URLs"
			],
			"endpoints": {
				"properties": {
					"basic": ["GET,POST /api/properties", "GET,PUT,DELETE /api/properties/{id}", "GET /api/properties/slug/{slug}"],
					"search": ["GET /api/properties/filter", "GET /api/properties/search/ranked", "GET /api/properties/search/suggestions", "POST /api/properties/search/advanced"],
					"analytics": ["GET /api/properties/statistics"]
				},
				"images": {
					"basic": ["POST /api/images", "GET,PUT,DELETE /api/images/{id}"],
					"property_images": ["GET /api/properties/{id}/images", "POST /api/properties/{id}/images/reorder"],
					"main_image": ["GET,POST /api/properties/{id}/images/main"],
					"processing": ["GET /api/images/{id}/variant", "GET /api/images/{id}/thumbnail"],
					"stats": ["GET /api/images/stats", "GET /api/images/cache/stats"],
					"maintenance": ["POST /api/images/cleanup"]
				},
				"users": {
					"auth": ["POST /api/auth/login"],
					"basic": ["GET,POST /api/users", "GET,PUT,DELETE /api/users/{id}"],
					"management": ["POST /api/users/{id}/password", "GET /api/users/role/{role}"],
					"analytics": ["GET /api/users/statistics", "GET /api/users/dashboard"]
				},
				"agencies": {
					"basic": ["GET,POST /api/agencies", "GET,PUT,DELETE /api/agencies/{id}"],
					"search": ["GET /api/agencies/active", "GET /api/agencies/service-area/{area}", "GET /api/agencies/specialty/{specialty}"],
					"management": ["GET /api/agencies/{id}/agents", "POST /api/agencies/{id}/license"],
					"analytics": ["GET /api/agencies/statistics", "GET /api/agencies/{id}/performance"]
				},
				"pagination": {
					"entities": ["GET /api/pagination/properties", "GET /api/pagination/images", "GET /api/pagination/users", "GET /api/pagination/agencies"],
					"search": ["GET /api/pagination/search"],
					"advanced": ["POST /api/pagination/advanced"],
					"stats": ["GET /api/pagination/stats"]
				},
				"health": {
					"basic": ["GET /api/health"],
					"detailed": ["GET /api/health/detailed"],
					"kubernetes": ["GET /api/health/ready", "GET /api/health/live"],
					"cache": ["GET /api/health/cache"],
					"metrics": ["GET /api/metrics"]
				},
				"security": {
					"monitoring": ["GET /api/security/metrics", "GET /api/security/status", "GET /api/security/threats"],
					"validation": ["POST /api/security/validate/password", "POST /api/security/validate/input"]
				},
				"monitoring": {
					"metrics": ["GET /api/monitoring/metrics", "GET /api/monitoring/prometheus", "GET /api/monitoring/dashboard"],
					"alerts": ["GET /api/monitoring/alerts", "GET /api/monitoring/alerts/history", "GET /api/monitoring/rules", "PUT /api/monitoring/rules/{rule}"]
				}
			},
			"total_endpoints": 68,
			"database": "PostgreSQL",
			"pagination_features": [
				"Configurable page size (max 100)",
				"Multiple sort fields",
				"Cross-entity search",
				"Advanced filtering",
				"Pagination metadata",
				"Performance optimized"
			],
			"note": "Full-featured real estate management system with advanced pagination ready for production"
		}`))
	})

	return mux
}

// performHealthCheck checks if the application is healthy
func performHealthCheck() error {
	cfg := config.LoadConfig()
	url := fmt.Sprintf("http://localhost:%s/api/health", cfg.Server.Port)
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to reach health endpoint: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned status %d", resp.StatusCode)
	}
	
	return nil
}

