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

	// Query all accounts
	rows, err := db.Query(`
		SELECT a.id, a.username, a.unixuid, a.type, a.created_at, g.groupname, g.unixgid, 
		       a.firstname, a.surname
		FROM accounts a
		LEFT JOIN groups g ON a.primary_group_id = g.id
		WHERE a.type = 'people'
		ORDER BY a.unixuid
	`)
	if err != nil {
		log.Fatalf("Failed to query accounts: %v", err)
	}
	defer rows.Close()

	fmt.Println("\nAll People Accounts:")
	fmt.Println("===================")
	fmt.Printf("%-5s %-12s %-8s %-7s %-20s %-15s %-8s %-15s %-15s\n", 
		"ID", "Username", "Unix UID", "Type", "Created At", "Primary Group", "Unix GID", "First Name", "Last Name")
	fmt.Println("-------------------------------------------------------------------------------------------------------")

	// Iterate through the rows
	for rows.Next() {
		var id int
		var username, accountType string
		var uid int
		var createdAt sql.NullTime
		var groupname sql.NullString
		var gid sql.NullInt64
		var firstname, surname sql.NullString

		if err := rows.Scan(&id, &username, &uid, &accountType, &createdAt, &groupname, &gid, &firstname, &surname); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}

		createdAtStr := ""
		if createdAt.Valid {
			createdAtStr = createdAt.Time.Format("2006-01-02 15:04:05")
		}

		groupStr := "None"
		if groupname.Valid {
			groupStr = groupname.String
		}

		gidStr := ""
		if gid.Valid {
			gidStr = fmt.Sprintf("%d", gid.Int64)
		}

		firstnameStr := ""
		if firstname.Valid {
			firstnameStr = firstname.String
		}

		surnameStr := ""
		if surname.Valid {
			surnameStr = surname.String
		}

		fmt.Printf("%-5d %-12s %-8d %-7s %-20s %-15s %-8s %-15s %-15s\n",
			id, username, uid, accountType, createdAtStr, groupStr, gidStr, firstnameStr, surnameStr)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating rows: %v", err)
	}
	
	// Query account memberships
	fmt.Println("\nAccount Memberships:")
	fmt.Println("===================")
	fmt.Printf("%-15s %-10s %-15s %-10s\n", "Username", "UID", "Group", "Unix GID")
	fmt.Println("----------------------------------------------")
	
	memberRows, err := db.Query(`
		SELECT a.username, a.unixuid, g.groupname, g.unixgid
		FROM account_groups ag
		JOIN accounts a ON ag.account_id = a.id
		JOIN groups g ON ag.group_id = g.id
		WHERE a.type = 'people'
		ORDER BY a.username, g.groupname
	`)
	if err != nil {
		log.Fatalf("Failed to query memberships: %v", err)
	}
	defer memberRows.Close()
	
	// Iterate through the membership rows
	for memberRows.Next() {
		var username string
		var uid int
		var groupname string
		var gid int
		
		if err := memberRows.Scan(&username, &uid, &groupname, &gid); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		
		fmt.Printf("%-15s %-10d %-15s %-10d\n", username, uid, groupname, gid)
	}
	
	// Check for errors from iterating over rows
	if err := memberRows.Err(); err != nil {
		log.Fatalf("Error iterating membership rows: %v", err)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}