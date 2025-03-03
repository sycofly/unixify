package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/home/unixify/internal/config"
)

func main() {
	// Load application configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to the database
	dsn := cfg.Database.GetDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Check connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Connected to database successfully")

	// Add created_by column to groups table if it doesn't exist
	_, err = db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT FROM information_schema.columns 
				WHERE table_name = 'groups' AND column_name = 'created_by'
			) THEN
				ALTER TABLE groups ADD COLUMN created_by TEXT;
				COMMENT ON COLUMN groups.created_by IS 'Username of the person who created this group';
			END IF;
		END
		$$;
	`)
	if err != nil {
		log.Fatalf("Failed to add created_by column: %v", err)
	}

	fmt.Println("Successfully added created_by column to groups table")
}