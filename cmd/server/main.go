package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/home/unixify/feature/ai_assisted/internal/config"
	"github.com/home/unixify/feature/ai_assisted/internal/database"
	"github.com/home/unixify/feature/ai_assisted/internal/handlers"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	// Templates directory - we'll look in both current directory and root
	templateDirs := []string{
		"/web/templates",   // When run in container with mounted volume
		"./web/templates",  // When run directly
	}

	// Find the first existing template directory
	var templateDir string
	for _, dir := range templateDirs {
		if _, err := os.Stat(dir); err == nil {
			templateDir = dir
			break
		}
	}

	if templateDir == "" {
		log.Fatalf("Could not find templates directory")
	}

	log.Printf("Using templates directory: %s", templateDir)

	// Create handlers
	h := handlers.NewHandler(db, templateDir)

	// Static file directories - we'll look in both current directory and root
	staticDirs := []string{
		"/web/static",   // When run in container with mounted volume
		"./web/static",  // When run directly
	}

	// Find the first existing static directory
	var staticDir string
	for _, dir := range staticDirs {
		if _, err := os.Stat(dir); err == nil {
			staticDir = dir
			break
		}
	}

	if staticDir == "" {
		log.Fatalf("Could not find static directory")
	}

	log.Printf("Using static directory: %s", staticDir)

	// Static files
	fileServer := http.FileServer(http.Dir(staticDir))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// Register routes
	h.RegisterHomeRoutes(r)
	h.RegisterAccountRoutes(r)
	h.RegisterGroupRoutes(r)

	// Create server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}