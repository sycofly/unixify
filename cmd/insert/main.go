package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

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

	// Insert a test people group
	createdAt := time.Now()
	result, err := db.Exec(`
		INSERT INTO groups (groupname, unixgid, description, type, created_by, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`, "test_group", 9999, "Test people group", "people", "manual_insert", createdAt)
	
	if err != nil {
		log.Fatalf("Failed to insert test group: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Failed to get rows affected: %v", err)
	}

	fmt.Printf("Inserted %d rows\n", rows)

	// Now query to verify it was inserted correctly
	var id int
	var groupname, description, typ string
	var createdBy sql.NullString
	var gID int
	var dbCreatedAt sql.NullTime

	err = db.QueryRow(`
		SELECT id, groupname, unixgid, description, type, created_by, created_at
		FROM groups
		WHERE groupname = 'test_group'
	`).Scan(&id, &groupname, &gID, &description, &typ, &createdBy, &dbCreatedAt)

	if err != nil {
		log.Fatalf("Failed to query inserted row: %v", err)
	}

	createdByStr := ""
	if createdBy.Valid {
		createdByStr = createdBy.String
	}

	createdAtStr := ""
	if dbCreatedAt.Valid {
		createdAtStr = dbCreatedAt.Time.Format("2006-01-02 15:04:05")
	}

	fmt.Println("\nInserted Record:")
	fmt.Println("===============")
	fmt.Printf("ID: %d\n", id)
	fmt.Printf("Group Name: %s\n", groupname)
	fmt.Printf("G_ID: %d\n", gID)
	fmt.Printf("Description: %s\n", description)
	fmt.Printf("Type: %s\n", typ)
	fmt.Printf("Created By: %s\n", createdByStr)
	fmt.Printf("Created At: %s\n", createdAtStr)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}