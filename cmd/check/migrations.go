package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
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

	// Connect to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Check connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Query migration status
	rows, err := db.Query(`
		SELECT version, dirty
		FROM schema_migrations
	`)
	if err != nil {
		log.Fatalf("Failed to query migration status: %v", err)
	}
	defer rows.Close()

	fmt.Println("\nMigration Status:")
	fmt.Println("================")
	fmt.Printf("%-10s %-10s\n", "Version", "Dirty")
	fmt.Println("-------------------")

	// Iterate through the rows
	for rows.Next() {
		var version int64
		var dirty bool

		if err := rows.Scan(&version, &dirty); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}

		fmt.Printf("%-10d %-10t\n", version, dirty)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating rows: %v", err)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}