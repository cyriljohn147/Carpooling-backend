package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	db "github.com/cyriljohn147/carpooling/db/sqlc"
	"github.com/cyriljohn147/carpooling/util"
)

// Server serves HTTP requests for the carpooling service
type Server struct {
	config util.Config
	store  *db.Queries
	db     *sql.DB
	router *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing
func NewServer(config util.Config, store *db.Queries, database *sql.DB) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
		db:     database,
	}

	// Setup router
	server.setupRouter()

	return server, nil
}

// setupRouter sets up the HTTP router with middleware and routes
func (server *Server) setupRouter() {
	// Create router
	router := gin.Default()

	// Setup CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"}, // Add your frontend URLs
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Add logging middleware
	router.Use(gin.Logger())

	// Add recovery middleware
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", server.healthCheck)

	// API info endpoint
	router.GET("/", server.apiInfo)

	// Register API routes
	server.registerRoutes(router)

	server.router = router
}

// registerRoutes registers all API routes
func (server *Server) registerRoutes(router *gin.Engine) {
	// Register user routes
	RegisterUserRoutes(router, server.store)

	// Register trip routes
	RegisterTripRoutes(router, server.store, server.db)

	// Register balance routes
	RegisterBalanceRoutes(router, server.store, server.db)

	// You can add more route groups here as needed
	// RegisterOtherRoutes(router, server.store, server.db)
}

// healthCheck handles health check requests
func (server *Server) healthCheck(c *gin.Context) {
	// Test database connection
	if err := server.db.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unhealthy",
			"message": "Database connection failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"message":   "Carpooling API is running",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
	})
}

// apiInfo provides basic API information
func (server *Server) apiInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":        "Carpooling API",
		"version":     "1.0.0",
		"description": "REST API for carpooling trip management",
		"endpoints": gin.H{
			"health":     "/health",
			"users":      "/users",
			"trips":      "/trips",
			"balances":   "/balances",
		},
		"timestamp": time.Now().UTC(),
	})
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	log.Printf("Starting carpooling server on %s", address)
	log.Printf("API endpoints available at:")
	log.Printf("  Health check: GET %s/health", address)
	log.Printf("  API info: GET %s/", address)
	log.Printf("  Users: %s/users", address)
	log.Printf("  Trips: %s/trips", address)
	log.Printf("  Balances: %s/balances", address)

	return server.router.Run(address)
}

// Shutdown gracefully shuts down the server
func (server *Server) Shutdown() error {
	log.Println("Shutting down server...")
	if server.db != nil {
		return server.db.Close()
	}
	return nil
}