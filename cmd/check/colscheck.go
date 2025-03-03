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

	// Query table columns for accounts table
	rows, err := db.Query(`
		SELECT column_name, data_type, column_default 
		FROM information_schema.columns 
		WHERE table_name = 'accounts'
		ORDER BY ordinal_position
	`)
	if err != nil {
		log.Fatalf("Failed to query columns: %v", err)
	}
	defer rows.Close()

	fmt.Println("Columns in accounts table:")
	fmt.Println("=======================")
	fmt.Printf("%-20s %-20s %-30s\n", "Column Name", "Data Type", "Default Value")
	fmt.Println("----------------------------------------------------------")

	// Iterate through the rows
	for rows.Next() {
		var name, dataType, defaultValue sql.NullString
		if err := rows.Scan(&name, &dataType, &defaultValue); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}

		defaultVal := ""
		if defaultValue.Valid {
			defaultVal = defaultValue.String
		}

		fmt.Printf("%-20s %-20s %-30s\n", name.String, dataType.String, defaultVal)
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