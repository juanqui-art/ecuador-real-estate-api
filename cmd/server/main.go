package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"realty-core/internal/repository"
	"realty-core/internal/service"
	"realty-core/internal/handlers"

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

	// Create services
	propertyService := service.NewPropertyService(propertyRepo, imageRepo)

	// Create handlers
	propertyHandler := handlers.NewPropertyHandler(propertyService)
	// TODO: Enable when dependencies are fixed
	// imageHandler := handlers.NewImageHandler(imageService)
	// userHandler := handlers.NewUserHandler(userService, permissionService, log.Default())
	// agencyHandler := handlers.NewAgencyHandler(agencyService, permissionService, log.Default())

	// Configure routes
	router := configureRoutesBasic(propertyHandler)

	// Configure HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      addMiddleware(router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
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

// configureRoutesBasic configures basic property routes that work
func configureRoutesBasic(propertyHandler *handlers.PropertyHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// ============== BASIC PROPERTY ROUTES ==============
	
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

	// Individual property routes
	mux.HandleFunc("/api/properties/", func(w http.ResponseWriter, r *http.Request) {
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
	})

	// Property by slug
	mux.HandleFunc("/api/properties/slug/", propertyHandler.GetPropertyBySlug)

	// Basic filtering and search
	mux.HandleFunc("/api/properties/filter", propertyHandler.FilterProperties)
	mux.HandleFunc("/api/properties/statistics", propertyHandler.GetStatistics)

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
			"message": "Real Estate Properties API - Basic System",
			"version": "1.5.0",
			"status": "Basic endpoints functional - Advanced features in development",
			"features": [
				"Property Management (CRUD)",
				"Basic Search & Filtering", 
				"PostgreSQL Database",
				"SEO-friendly URLs"
			],
			"endpoints": {
				"properties": {
					"basic": ["GET,POST /api/properties", "GET,PUT,DELETE /api/properties/{id}", "GET /api/properties/slug/{slug}"],
					"search": ["GET /api/properties/filter"],
					"analytics": ["GET /api/properties/statistics"]
				},
				"system": ["GET /api/health"]
			},
			"total_endpoints": 6,
			"database": "PostgreSQL",
			"note": "Extended features (images, users, agencies, FTS, pagination) coming soon"
		}`))
	})

	return mux
}
