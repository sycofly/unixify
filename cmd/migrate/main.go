package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	// Parse command line arguments
	var direction string
	var forceVersion int
	flag.StringVar(&direction, "direction", "up", "Migration direction (up or down)")
	flag.IntVar(&forceVersion, "force", -1, "Force database version (use with caution)")
	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get database connection string from environment variables
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "postgres")
	dbName := getEnvOrDefault("DB_NAME", "unixify")
	dbSSLMode := getEnvOrDefault("DB_SSLMODE", "disable")

	// Create database URL
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	// Create a new migrate instance
	m, err := migrate.New("file://db/migrations", dbURL)
	if err != nil {
		log.Fatalf("Migration failed to initialize: %v", err)
	}

	// Handle force version mode
	if forceVersion >= 0 {
		err := m.Force(forceVersion)
		if err != nil {
			log.Fatalf("Failed to force version to %d: %v", forceVersion, err)
		}
		log.Printf("Successfully forced database version to %d", forceVersion)
		return
	}

	// Run migrations
	if direction == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migration up completed successfully")
	} else if direction == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Migration down completed successfully")
	} else {
		log.Fatalf("Invalid direction: %s. Use 'up' or 'down'", direction)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}