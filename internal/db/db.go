package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	dbInstance *sql.DB
	once       sync.Once
)

// GetDB returns a singleton instance of the database connection
func GetDB() *sql.DB {
	once.Do(func() {
		// Create data directory if it doesn't exist
		dataDir := getDataDir()
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			log.Fatalf("Failed to create data directory: %v", err)
		}

		dbPath := filepath.Join(dataDir, "property_management.db")
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}

		// Test connection
		if err := db.Ping(); err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		// Initialize the database schema
		if err := initSchema(db); err != nil {
			log.Fatalf("Failed to initialize database schema: %v", err)
		}

		dbInstance = db
	})

	return dbInstance
}

// Close closes the database connection
func Close() {
	if dbInstance != nil {
		dbInstance.Close()
	}
}

// Initialize database schema
func initSchema(db *sql.DB) error {
	// Create houses table
	housesSchema := `
	CREATE TABLE IF NOT EXISTS houses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		street TEXT NOT NULL,
		number TEXT NOT NULL,
		country TEXT NOT NULL,
		zip_code TEXT NOT NULL,
		city TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(housesSchema)
	if err != nil {
		return err
	}

	// Create apartments table
	apartmentsSchema := `
	CREATE TABLE IF NOT EXISTS apartments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		house_id INTEGER NOT NULL,
		size REAL NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (house_id) REFERENCES houses(id)
	);`

	_, err = db.Exec(apartmentsSchema)
	return err
}

// getDataDir returns the path to the data directory
func getDataDir() string {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	// Create application-specific data directory
	dataDir := filepath.Join(homeDir, ".property-management")
	return dataDir
}
