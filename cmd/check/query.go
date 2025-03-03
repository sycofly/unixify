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

	// Query all groups with no limit
	rows, err := db.Query(`
		SELECT id, groupname, unixgid, description, type, created_by, created_at
		FROM groups
		ORDER BY id
	`)
	if err != nil {
		log.Fatalf("Failed to query groups: %v", err)
	}
	defer rows.Close()

	fmt.Println("\nAll Groups:")
	fmt.Println("===========")
	fmt.Printf("%-5s %-15s %-10s %-30s %-10s %-15s %-20s\n", "ID", "Group Name", "Unix GID", "Description", "Type", "Created By", "Created At")
	fmt.Println("---------------------------------------------------------------------------------------------------------")

	// Iterate through the rows
	for rows.Next() {
		var id int
		var groupname, description, typ string
		var createdBy sql.NullString
		var gID int
		var createdAt sql.NullTime

		if err := rows.Scan(&id, &groupname, &gID, &description, &typ, &createdBy, &createdAt); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}

		createdByStr := ""
		if createdBy.Valid {
			createdByStr = createdBy.String
		}

		createdAtStr := ""
		if createdAt.Valid {
			createdAtStr = createdAt.Time.Format("2006-01-02 15:04:05")
		}

		fmt.Printf("%-5d %-15s %-10d %-30s %-10s %-15s %-20s\n", 
			id, groupname, gID, description, typ, createdByStr, createdAtStr)
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