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

	// Create service
	propertyService := service.NewPropertyService(propertyRepo, imageRepo)

	// Create handler
	propertyHandler := handlers.NewPropertyHandler(propertyService)

	// Configure routes
	router := configureRoutes(propertyHandler)

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

// configureRoutes configures all application routes
func configureRoutes(propertyHandler *handlers.PropertyHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// Property routes
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

	// Property by slug route
	mux.HandleFunc("/api/properties/slug/", propertyHandler.GetPropertyBySlug)

	// Filter properties route
	mux.HandleFunc("/api/properties/filter", propertyHandler.FilterProperties)

	// Statistics route
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
			"message": "Real Estate Properties API",
			"version": "1.0.0",
			"endpoints": {
				"GET /api/health": "Service status",
				"GET /api/properties": "List all properties",
				"POST /api/properties": "Create new property",
				"GET /api/properties/{id}": "Get property by ID",
				"GET /api/properties/slug/{slug}": "Get property by SEO slug",
				"PUT /api/properties/{id}": "Update property",
				"DELETE /api/properties/{id}": "Delete property",
				"GET /api/properties/filter": "Filter properties (query params: province, min_price, max_price, q)",
				"GET /api/properties/statistics": "Get statistics"
			},
			"seo": {
				"description": "SEO slugs are automatically generated from title",
				"format": "normalized-title-{id}",
				"example": "/api/properties/slug/beautiful-house-samborondon-abcd1234"
			}
		}`))
	})

	return mux
}
