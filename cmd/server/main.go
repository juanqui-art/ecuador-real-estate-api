package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"realty-core/internal/cache"
	"realty-core/internal/handlers"
	"realty-core/internal/middleware"
	"realty-core/internal/processors"
	"realty-core/internal/repository"
	"realty-core/internal/service"
	"realty-core/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system variables")
	}

	// Get configuration from environment variables
	databaseURL := getVariable("DATABASE_URL", "postgresql://juanquizhpi@localhost:5433/inmobiliaria_db?sslmode=disable")
	port := getVariable("PORT", "8080")
	logLevel := getVariable("LOG_LEVEL", "info")

	log.Printf("Starting real estate properties server...")
	log.Printf("Log Level: %s", logLevel)
	log.Printf("Port: %s", port)

	// Connect to database
	log.Println("Connecting to PostgreSQL...")
	db, err := repository.ConnectDatabase(databaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Create repositories
	propertyRepo := repository.NewPostgreSQLPropertyRepository(db)
	imageRepo := repository.NewPostgreSQLImageRepository(db)
	userRepo := repository.NewUserRepository(db)
	agencyRepo := repository.NewAgencyRepository(db)

	// Create image system dependencies
	imageStorage, err := storage.NewLocalImageStorage("uploads/images", "/uploads/images", 10*1024*1024) // 10MB max file size
	if err != nil {
		log.Fatalf("Error creating image storage: %v", err)
	}
	imageProcessor := processors.NewImageProcessor(3000, 2000)
	imageCacheConfig := cache.ImageCacheConfig{
		Enabled:         true,
		Capacity:        1000,
		MaxSizeBytes:    100 * 1024 * 1024, // 100MB
		TTL:             24 * time.Hour,
		CleanupInterval: 10 * time.Minute,
	}
	imageCache := cache.NewImageCache(imageCacheConfig)

	// Create services
	propertyService := service.NewPropertyService(propertyRepo, imageRepo)
	imageService := service.NewImageService(imageRepo, propertyRepo, imageStorage, imageProcessor, imageCache)
	userService := service.NewUserService(userRepo, agencyRepo, log.Default())
	agencyService := service.NewAgencyService(agencyRepo, userRepo, log.Default())
	permissionService := service.NewPermissionService()

	// Create handlers
	propertyHandler := handlers.NewPropertyHandler(propertyService)
	imageHandler := handlers.NewImageHandler(imageService)
	userHandler := handlers.NewUserHandlerSimple(userService, permissionService, log.Default())
	agencyHandler := handlers.NewAgencyHandlerSimple(agencyService, permissionService, log.Default())
	paginationHandler := handlers.NewPaginationHandlerSimple(propertyService, imageService, userService, agencyService, log.Default())

	// Configure routes
	router := configureRoutesWithPagination(propertyHandler, imageHandler, userHandler, agencyHandler, paginationHandler)
	
	// Apply performance monitoring middleware
	finalHandler := middleware.PerformanceLogger(router)

	// Configure HTTP server with optimized settings
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      addMiddleware(finalHandler),
		ReadTimeout:  10 * time.Second,  // Shorter read timeout for better responsiveness
		WriteTimeout: 30 * time.Second,  // Longer write timeout for large responses
		IdleTimeout:  120 * time.Second, // Longer idle timeout for persistent connections
		MaxHeaderBytes: 1 << 20,         // 1MB max header size
	}

	log.Printf("🚀 Server started at http://localhost:%s", port)
	log.Printf("📚 Documentation available at http://localhost:%s/", port)
	log.Printf("🏥 Health check at http://localhost:%s/api/health", port)

	// Start server
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// getVariable gets an environment variable or returns a default value
func getVariable(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// addMiddleware adds common middleware to all routes
func addMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Logging middleware
		start := time.Now()

		// Configure CORS headers for development
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Execute next handler
		next.ServeHTTP(w, r)

		// Log the request
		duration := time.Since(start)
		log.Printf("%s %s - %s", r.Method, r.URL.Path, duration)
	})
}



// configureRoutesWithPagination configures routes with all systems including pagination
func configureRoutesWithPagination(propertyHandler *handlers.PropertyHandler, imageHandler *handlers.ImageHandler, userHandler *handlers.UserHandlerSimple, agencyHandler *handlers.AgencyHandlerSimple, paginationHandler *handlers.PaginationHandlerSimple) *http.ServeMux {
	mux := http.NewServeMux()

	// ============== PROPERTY ROUTES ==============
	
	// Basic CRUD
	mux.HandleFunc("/api/properties", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			propertyHandler.ListProperties(w, r)
		case http.MethodPost:
			propertyHandler.CreateProperty(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Property by slug (must be before /api/properties/{id})
	mux.HandleFunc("/api/properties/slug/", propertyHandler.GetPropertyBySlug)
	
	// Basic filtering and search
	mux.HandleFunc("/api/properties/filter", propertyHandler.FilterProperties)
	mux.HandleFunc("/api/properties/statistics", propertyHandler.GetStatistics)

	// Property images management (specific routes first)
	mux.HandleFunc("/api/properties/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.Contains(path, "/images") {
			if strings.HasSuffix(path, "/images") && r.Method == http.MethodGet {
				imageHandler.GetImagesByProperty(w, r)
			} else if strings.Contains(path, "/images/reorder") && r.Method == http.MethodPost {
				imageHandler.ReorderImages(w, r)
			} else if strings.Contains(path, "/images/main") {
				if r.Method == http.MethodGet {
					imageHandler.GetMainImage(w, r)
				} else if r.Method == http.MethodPost {
					imageHandler.SetMainImage(w, r)
				}
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		} else {
			// Regular property operations
			switch r.Method {
			case http.MethodGet:
				propertyHandler.GetProperty(w, r)
			case http.MethodPut:
				propertyHandler.UpdateProperty(w, r)
			case http.MethodDelete:
				propertyHandler.DeleteProperty(w, r)
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
			imageHandler.UploadImage(w, r)
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

	// ============== USER/AUTH ROUTES ==============

	// Authentication routes
	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			userHandler.Login(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// User management routes
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.SearchUsers(w, r)
		case http.MethodPost:
			userHandler.CreateUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// User statistics and dashboard (specific routes first)
	mux.HandleFunc("/api/users/statistics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			userHandler.GetUserStatistics(w, r)
		}
	})

	mux.HandleFunc("/api/users/dashboard", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			userHandler.GetUserDashboard(w, r)
		}
	})

	// Users by role
	mux.HandleFunc("/api/users/role/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			userHandler.GetUsersByRole(w, r)
		}
	})

	// Individual user routes (must be after specific routes)
	mux.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.Contains(path, "/password") && r.Method == http.MethodPost {
			userHandler.ChangePassword(w, r)
		} else {
			// Regular user operations
			switch r.Method {
			case http.MethodGet:
				userHandler.GetUser(w, r)
			case http.MethodPut:
				userHandler.UpdateUser(w, r)
			case http.MethodDelete:
				userHandler.DeleteUser(w, r)
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
			agencyHandler.SearchAgencies(w, r)
		case http.MethodPost:
			agencyHandler.CreateAgency(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Agency statistics (specific routes first)
	mux.HandleFunc("/api/agencies/statistics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			agencyHandler.GetAgencyStatistics(w, r)
		}
	})

	mux.HandleFunc("/api/agencies/active", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			agencyHandler.GetActiveAgencies(w, r)
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

	// Health check route
	mux.HandleFunc("/api/health", propertyHandler.HealthCheck)

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
					"search": ["GET /api/properties/filter"],
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
				"system": ["GET /api/health"]
			},
			"total_endpoints": 51,
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

