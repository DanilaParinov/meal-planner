package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
	"meal-planner/internal/config"
)

// Database encapsulates the database connection
type Database struct {
	db *sql.DB
}

// New creates a new database connection
func New(cfg *config.Config) (*Database, error) {
	dbURL := cfg.GetDatabaseURL()

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify the connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✓ Connected to database")

	d := &Database{db: conn}

	// Run migrations
	if err := d.runMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return d, nil
}

// runMigrations executes all SQL migrations from the migrations/ folder
func (d *Database) runMigrations() error {
	// Read migrations from the folder
	files, err := os.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			filePath := filepath.Join("migrations", f.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", filePath, err)
			}

			// Execute SQL
			if _, err := d.db.Exec(string(content)); err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", filePath, err)
			}

			log.Printf("✓ Executed migration: %s", f.Name())
		}
	}

	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// GetDB returns the underlying connection (for queries)
func (d *Database) GetDB() *sql.DB {
	return d.db
}
