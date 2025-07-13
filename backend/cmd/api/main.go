// Package main contains the main entrypoint for the application
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pet-mgt/backend/internal/config"
	"pet-mgt/backend/internal/routes"
	"pet-mgt/backend/internal/store"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadCfg()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	// Create main context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := store.NewSupabaseService(cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// Test database connection
	if err := db.Ping(ctx); err != nil {
		log.Printf("database ping failed: %v", err)
	} else {
		log.Println("database connection established")
	}

	r := routes.SetupRouter(cfg, db)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("starting server on port %s", cfg.Port)
		log.Printf("test connection at http://localhost:%s/ping", cfg.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("shutting down server...")

	// Cancel the main context to signal all goroutines to stop
	cancel()

	// Give the server 30 seconds to finish handling existing requests
	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	} else {
		log.Println("server gracefully stopped")
	}

	// Add any additional cleanup here (database connections, etc.)
	if err := db.Close(); err != nil {
		log.Printf("error closing database: %v", err)
	} else {
		log.Println("database connection closed")
	}

	log.Println("cleanup completed")
}
