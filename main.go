package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/cyriljohn147/carpooling/api"
	db "github.com/cyriljohn147/carpooling/db/sqlc"
	"github.com/cyriljohn147/carpooling/util"
)

func main() {
	// Load configuration
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer conn.Close()

	// Test database connection
	if err := conn.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Database connection established")

	// Run database migrations
	runDBMigrations(config.MigrationURL, config.DBSource)

	// Create store
	store := db.New(conn)

	// Create server
    server, err := api.NewServer(config, store, conn)
    if err != nil {
    	log.Fatal("Failed to create server:", err)
    }

    // Change server address to avoid port conflict
    config.ServerAddress = "0.0.0.0:8081"

    // Start server in a goroutine
    go func() {
    	if err := server.Start(config.ServerAddress); err != nil && err != http.ErrServerClosed {
    		log.Fatal("Failed to start server:", err)
    	}
    }()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	select {
	case <-ctx.Done():
		log.Println("Timeout of 10 seconds reached")
	default:
		log.Println("Server exiting")
	}
}

func runDBMigrations(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("Cannot create migration instance:", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run migrate up:", err)
	}

	log.Println("Database migration completed successfully")
}
